# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

import json
import os
import time

from cinderclient import exceptions as cinder_exception
from oslo_config import cfg
from oslo_log import log as logging
from oslo_utils import importutils
from oslo_utils import strutils

from fuxi.common import constants as consts
from fuxi.common import mount
from fuxi.common import state_monitor
from fuxi import exceptions
from fuxi.i18n import _, _LE, _LI, _LW
from fuxi import utils
from fuxi.volumeprovider import provider
import etcd

CONF = cfg.CONF
cinder_conf = CONF.cinder

# Volume states
UNKNOWN = consts.UNKNOWN
NOT_ATTACH = consts.NOT_ATTACH
ATTACH_TO_THIS = consts.ATTACH_TO_THIS
ATTACH_TO_OTHER = consts.ATTACH_TO_OTHER

OPENSTACK = 'openstack'
OSBRICK = 'osbrick'
OPENSDS = 'opensds'

volume_connector_conf = {
    OPENSTACK: 'fuxi.connector.cloudconnector.openstack.CinderConnector',
    OSBRICK: 'fuxi.connector.osbrickconnector.CinderConnector',
    OPENSDS: 'fuxi.connector.osbrickconnector.OpenSDSConnector'}

LOG = logging.getLogger(__name__)

def get_host_id():
    """Get a value that could represent this server."""
    host_id = None
    volume_connector = cinder_conf.volume_connector
    if volume_connector == OPENSTACK:
        host_id = utils.get_instance_uuid()
    elif volume_connector == OSBRICK:
        host_id = utils.get_hostname().lower()
    return host_id

class APIDictWrapper(object):
    """Simple wrapper for api dictionaries
    Some api calls return dictionaries.  This class provides identical
    behavior as APIResourceWrapper, except that it will also behave as a
    dictionary, in addition to attribute accesses.
    Attribute access is the preferred method of access, to be
    consistent with api resource objects from novaclient.
    """

    _apidict = {}  # Make sure _apidict is there even in __init__.

    def __init__(self, apidict):
        self._apidict = apidict

    def __getattribute__(self, attr):
        try:
            return object.__getattribute__(self, attr)
        except AttributeError:
            if attr not in self._apidict:
                raise
            return self._apidict[attr]

    def __getitem__(self, item):
        try:
            return getattr(self, item)
        except (AttributeError, TypeError) as e:
            # caller is expecting a KeyError
            raise KeyError(e)

    def __contains__(self, item):
        try:
            return hasattr(self, item)
        except TypeError:
            return False

    def get(self, item, default=None):
        try:
            return getattr(self, item)
        except (AttributeError, TypeError):
            return default

    def __repr__(self):
        return "<%s: %s>" % (self.__class__.__name__, self._apidict)

    def to_dict(self):
        return self._apidict

class GrpcApi(object):
	TMOUT = 10
	RESOURCE_TYPE = "cinder"
	def __init__(self):
		self._etcd = etcd.Client(host='127.0.0.1', port=2379)
		self._url = "opensds/api"

	def _joinargs(self, action, *args):
		val = [action, self.RESOURCE_TYPE]
		val.extend(args)
		return ",".join(map(lambda x: str(x), val))

	def _call(self, action, *args):
		val = self._joinargs(action, *args)
		resp = self._etcd.set(self._url, val)
		info = self._etcd.watch(self._url, index=resp.etcd_index + 1,
								timeout=self.TMOUT)
		try:
			ret = json.loads(info.value)
		except ValueError:
			return info.value
		else:
			return ret

	def create(self, name, size):
		return self._call("CreateVolume", name, size)

	def get(self, voluuid):
		return self._call("GetVolume", voluuid)

	def list(self, allow_detail = True):
		return self._call("GetAllVolumes", allow_detail)

	def update(self, voluuid, name):
		return self._call("UpdateVolume", voluuid, name)

	def delete(self, voluuid):
		return self._call("DeleteVolume", voluuid)

class OpenSDS(provider.Provider):
    volume_provider_type = 'opensds'

    def __init__(self):
        super(OpenSDS, self).__init__()
        self.opensdsclient = GrpcApi()

    def _get_connector(self):
        connector = cinder_conf.volume_connector
        if not connector or connector not in volume_connector_conf:
            msg = _LE("Must provide an valid volume connector")
            LOG.error(msg)
            raise exceptions.FuxiException(msg)
        return importutils.import_class(volume_connector_conf[connector])()

    def _get_docker_volume(self, docker_volume_name):
        LOG.info(_LI("Retrieve docker volume {0} from "
                     "OpenSDS").format(docker_volume_name))

        try:
            host_id = get_host_id()

            volume_connector = cinder_conf.volume_connector
            # search_opts = {'name': docker_volume_name,
            #                'metadata': {consts.VOLUME_FROM: CONF.volume_from}}
            for vol in self.opensdsclient.list():
				vol = APIDictWrapper(vol)
                LOG.info("dsl: vol=%s", vol.__dict__)
                if vol.name == docker_volume_name:
                    if vol.attachments:
                        for am in vol.attachments:
                            if volume_connector == OPENSTACK:
                                if am['server_id'] == host_id:
                                    return vol, ATTACH_TO_THIS
                            elif volume_connector == OSBRICK:
                                if (am['host_name'] or '').lower() == host_id:
                                    return vol, ATTACH_TO_THIS
                        return vol, ATTACH_TO_OTHER
                    else:
                        return vol, NOT_ATTACH
            return None, UNKNOWN
        except cinder_exception.ClientException as ex:
            LOG.error(_LE("Error happened while getting volume list "
                          "information from cinder. Error: {0}").format(ex))
            raise

    def _check_attached_to_this(self, cinder_volume):
        host_id = get_host_id()
        vol_conn = cinder_conf.volume_connector
        for am in cinder_volume.attachments:
            if vol_conn == OPENSTACK and am['server_id'] == host_id:
                return True
            elif vol_conn == OSBRICK and am['host_name'] \
                    and am['host_name'].lower() == host_id:
                return True
        return False

    def _create_volume(self, docker_volume_name, volume_opts):
        LOG.info(_LI("Start to create docker volume {0} from "
                     "OpenSDS").format(docker_volume_name))

        try:
            volume = self.opensdsclient.create(docker_volume_name, volume_opts['size'])
        except cinder_exception.ClientException as e:
            msg = _LE("Error happened when create an volume {0} from OpenSDS. "
                      "Error: {1}").format(docker_volume_name, e)
            LOG.error(msg)
            raise
	volume = APIDictWrapper(volume)
        time.sleep(5)

        LOG.info(_LI("Create docker volume {0} {1} from OpenSDS "
                     "successfully").format(docker_volume_name, volume))
        return volume

    def create(self, docker_volume_name, volume_opts):
        if not volume_opts:
            volume_opts = {"size": 1}

        connector = self._get_connector()
        cinder_volume, state = self._get_docker_volume(docker_volume_name)
        LOG.info(_LI("Get docker volume {0} {1} with state "
                     "{2}").format(docker_volume_name, cinder_volume, state))

        device_info = {}
        if state == ATTACH_TO_THIS:
            LOG.warning(_LW("The volume {0} {1} already exists and attached "
                            "to this server").format(docker_volume_name,
                                                     cinder_volume))
            device_info = {'path': connector.get_device_path(cinder_volume)}
        elif state == NOT_ATTACH:
            LOG.warning(_LW("The volume {0} {1} is already exists but not "
                            "attached").format(docker_volume_name,
                                               cinder_volume))
            device_info = connector.connect_volume(cinder_volume)
        elif state == ATTACH_TO_OTHER:
            if cinder_volume.multiattach:
                fstype = volume_opts.get('fstype', cinder_conf.fstype)
                vol_fstype = cinder_volume.metadata.get('fstype',
                                                        cinder_conf.fstype)
                if fstype != vol_fstype:
                    msg = _LE("Volume already exists with fstype: {0}, but "
                              "currently provided fstype is {1}, not "
                              "match").format(vol_fstype, fstype)
                    LOG.error(msg)
                    raise exceptions.FuxiException('FSType Not Match')
                device_info = connector.connect_volume(cinder_volume)
            else:
                msg = _LE("The volume {0} {1} is already attached to another "
                          "server").format(docker_volume_name, cinder_volume)
                LOG.error(msg)
                raise exceptions.FuxiException(msg)
        elif state == UNKNOWN:
            cinder_volume = self._create_volume(docker_volume_name,
                                                    volume_opts)
            device_info = connector.connect_volume(cinder_volume)

        return device_info

    def _delete_volume(self, volume):
        try:
            self.opensdsclient.delete(volume.id)
        except cinder_exception.NotFound:
            return
        except cinder_exception.ClientException as e:
            msg = _LE("Error happened when delete volume from OpenSDS. "
                      "Error: {0}").format(e)
            LOG.error(msg)
            raise

        """start_time = time.time()
        # Wait until the volume is not there or until the operation timeout
        while (time.time() - start_time < consts.DESTROY_VOLUME_TIMEOUT):
            try:
                self.opensdsclient.get(volume.id)
            except cinder_exception.NotFound:
                return
            time.sleep(consts.VOLUME_SCAN_TIME_DELAY)

        # If the volume is not deleted, raise an exception
        msg_ft = _LE("Timed out while waiting for volume. "
                     "Expected Volume: {0}, "
                     "Expected State: {1}, "
                     "Elapsed Time: {2}").format(volume,
                                                 None,
                                                 time.time() - start_time)
        raise exceptions.TimeoutException(msg_ft)"""

    def delete(self, docker_volume_name):
        cinder_volume, state = self._get_docker_volume(docker_volume_name)
        LOG.info(_LI("Get docker volume {0} {1} with state "
                     "{2}").format(docker_volume_name, cinder_volume, state))

        if state == ATTACH_TO_THIS:
            link_path = self._get_connector().get_device_path(cinder_volume)
            if not link_path or not os.path.exists(link_path):
                msg = _LE(
                    "Could not find device link path for volume {0} {1} "
                    "in host").format(docker_volume_name, cinder_volume)
                LOG.error(msg)
                raise exceptions.FuxiException(msg)

            devpath = os.path.realpath(link_path)
            if not os.path.exists(devpath):
                msg = _LE("Could not find device path for volume {0} {1} in "
                          "host").format(docker_volume_name, cinder_volume)
                LOG.error(msg)
                raise exceptions.FuxiException(msg)

            mounter = mount.Mounter()
            mps = mounter.get_mps_by_device(devpath)
            ref_count = len(mps)
            if ref_count > 0:
                mountpoint = self._get_mountpoint(docker_volume_name)
                if mountpoint in mps:
                    mounter.unmount(mountpoint)

                    self._clear_mountpoint(mountpoint)

                    # If this volume is still mounted on other mount point,
                    # then return.
                    if ref_count > 1:
                        return True
                else:
                    return True

            # Detach device from this server.
            self._get_connector().disconnect_volume(cinder_volume)

            available_volume = self.opensdsclient.get(cinder_volume.id)
            # If this volume is not used by other server any more,
            # than delete it from Cinder.
	    available_volume = APIDictWrapper(available_volume)
            if not available_volume.attachments:
                msg = _LW("No other servers still use this volume {0} "
                          "{1} any more, so delete it from OpenSDS"
                          "").format(docker_volume_name, cinder_volume)
                LOG.warning(msg)
                self._delete_volume(available_volume)
            return True
        elif state == NOT_ATTACH:
            self._delete_volume(cinder_volume)
            return True
        elif state == ATTACH_TO_OTHER:
            msg = _LW("Volume %s is still in use, could not delete it")
            LOG.warning(msg, cinder_volume)
            return True
        elif state == UNKNOWN:
            return False
        else:
            msg = _LE("Volume %(vol_name)s %(c_vol)s "
                      "state %(state)s is invalid")
            LOG.error(msg, {'vol_name': docker_volume_name,
                            'c_vol': cinder_volume,
                            'state': state})
            raise exceptions.NotMatchedState()

    def list(self):
        LOG.info(_LI("Start to retrieve all docker volumes from OpenSDS"))

        docker_volumes = []
        try:
            for vol in self.opensdsclient.list():
		# LOG.info(_LI("Retrieve docker volumes {0} from OpenSDS "
                #      "successfully").format(vol))
                docker_volume_name = vol['name']
                if not docker_volume_name:
                    continue

                mountpoint = self._get_mountpoint(docker_volume_name)
		vol = APIDictWrapper(vol)
                devpath = os.path.realpath(
                    self._get_connector().get_device_path(vol))
                mps = mount.Mounter().get_mps_by_device(devpath)
                mountpoint = mountpoint if mountpoint in mps else ''
                docker_vol = {'Name': docker_volume_name,
                              'Mountpoint': mountpoint}
                docker_volumes.append(docker_vol)
        except cinder_exception.ClientException as e:
            LOG.error(_LE("Retrieve volume list failed. Error: {0}").format(e))
            raise

        LOG.info(_LI("Retrieve docker volumes {0} from OpenSDS "
                     "successfully").format(docker_volumes))
        return docker_volumes

    def show(self, docker_volume_name):
        cinder_volume, state = self._get_docker_volume(docker_volume_name)
        LOG.info(_LI("Get docker volume {0} {1} with state "
                     "{2}").format(docker_volume_name, cinder_volume, state))

        if state == ATTACH_TO_THIS:
            devpath = os.path.realpath(
                self._get_connector().get_device_path(cinder_volume))
            mp = self._get_mountpoint(docker_volume_name)
            LOG.info("Expected devpath: {0} and mountpoint: {1} for volume: "
                     "{2} {3}".format(devpath, mp, docker_volume_name,
                                      cinder_volume))
            mounter = mount.Mounter()
            return {"Name": docker_volume_name,
                    "Mountpoint": mp if mp in mounter.get_mps_by_device(
                        	      devpath) else ''}
        elif state in (NOT_ATTACH, ATTACH_TO_OTHER):
            return {'Name': docker_volume_name, 'Mountpoint': ''}
        elif state == UNKNOWN:
            msg = _LW("Can't find this volume '{0}' in "
                      "OpenSDS").format(docker_volume_name)
            LOG.warning(msg)
            raise exceptions.NotFound(msg)
        else:
            msg = _LE("Volume '{0}' exists, but not attached to this volume,"
                      "and current state is {1}").format(docker_volume_name,
                                                         state)
            raise exceptions.NotMatchedState(msg)

    def mount(self, docker_volume_name):
        cinder_volume, state = self._get_docker_volume(docker_volume_name)
        LOG.info(_LI("Get docker volume {0} {1} with state "
                     "{2}").format(docker_volume_name, cinder_volume, state))

        connector = self._get_connector()
        if state == NOT_ATTACH:
            connector.connect_volume(cinder_volume)
        elif state == ATTACH_TO_OTHER:
            if cinder_volume.multiattach:
                connector.connect_volume(cinder_volume)
            else:
                msg = _("Volume {0} {1} is not shareable").format(
                    docker_volume_name, cinder_volume)
                raise exceptions.FuxiException(msg)
        elif state != ATTACH_TO_THIS:
            msg = _("Volume %(vol_name)s %(c_vol)s is not in correct state, "
                    "current state is %(state)s")
            LOG.error(msg, {'vol_name': docker_volume_name,
                            'c_vol': cinder_volume,
                            'state': state})
            raise exceptions.NotMatchedState()

        link_path = connector.get_device_path(cinder_volume)
        if not os.path.exists(link_path):
            LOG.warning(_LW("Could not find device link file, "
                            "so rebuild it"))
            connector.disconnect_volume(cinder_volume)
            connector.connect_volume(cinder_volume)

        devpath = os.path.realpath(link_path)
        if not devpath or not os.path.exists(devpath):
            msg = _("Can't find volume device path")
            LOG.error(msg)
            raise exceptions.FuxiException(msg)

        mountpoint = self._get_mountpoint(docker_volume_name)
        self._create_mountpoint(mountpoint)

        fstype = cinder_volume.metadata.get('fstype', cinder_conf.fstype)

        mount.do_mount(devpath, mountpoint, fstype)

        return mountpoint

    def unmount(self, docker_volume_name):
        return

    def check_exist(self, docker_volume_name):
        _, state = self._get_docker_volume(docker_volume_name)
        LOG.info(_LI("Get docker volume {0} with state "
                     "{1}").format(docker_volume_name, state))

        if state == UNKNOWN:
            return False
        return True
