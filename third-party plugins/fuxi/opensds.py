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
import json

CONF = cfg.CONF
cinder_conf = CONF.cinder

# Volume states
UNKNOWN = consts.UNKNOWN
NOT_ATTACH = consts.NOT_ATTACH
ATTACH_TO_THIS = consts.ATTACH_TO_THIS
ATTACH_TO_OTHER = consts.ATTACH_TO_OTHER

OPENSTACK = 'openstack'
OSBRICK = 'osbrick'

volume_connector_conf = {
    OPENSTACK: 'fuxi.connector.cloudconnector.openstack.CinderConnector',
    OSBRICK: 'fuxi.connector.osbrickconnector.CinderConnector'}

LOG = logging.getLogger(__name__)

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
		info = self._etcd.watch(self._url, index=resp.etcd_index,
								timeout=self.TMOUT)
		return json.loads(info._prev_node.value)

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

	def create(self, docker_volume_name, volume_opts):
		if not volume_opts:
			volume_opts = {"size": 1}

		cinder_volume, state = self._get_docker_volume(docker_volume_name)
		if state == NOT_ATTACH:
			LOG.warning(_LW("The volume {0} {1} already exists and attached "
					"to this server").format(docker_volume_name,
								cinder_volume))
		else:
			return self.opensdsclient.create(docker_volume_name,
							volume_opts["size"])

	def delete(self, docker_volume_name):
		cinder_volume, state = self._get_docker_volume(docker_volume_name)
		LOG.info(_LI("Get docker volume {0} {1} with state "
			"{2}").format(docker_volume_name, cinder_volume, state))

        	if cinder_volume is not None:
        		self.opensdsclient.create(docker_volume_name, volume_opts["size"])

	def list(self):
        	LOG.info(_LI("Start to retrieve all docker volumes from OpenSDS"))
        	vols = self.opensdsclient.list()
        	docker_volumes = []
        	try:
        	    for vol in vols:
            		mountpoint = self._get_mountpoint(vol["name"])
            		docker_vol = {'Name': vol["name"],
            				'Mountpoint': mountpoint}
            		docker_volumes.append(docker_vol)
        	except cinder_exception.ClientException as e:
        	    LOG.error(_LE("Retrieve volume list failed. Error: {0}").format(e))
        	    raise

        	LOG.info(_LI("Retrieve docker volumes {0} from OpenSDS "
        	             "successfully").format(docker_volumes))
        	return docker_volumes

	def _get_docker_volume(self, docker_volume_name):
	        LOG.info(_LI("Retrieve docker volume {0} from "
	                     "OpenSDS").format(docker_volume_name))

	        for vol in self.opensdsclient.list():
	        	if vol["name"] == docker_volume_name:
	        		return vol, NOT_ATTACH
	        	else:
        			return None, UNKNOWN

	def show(self, docker_volume_name):

    		vols = self.opensdsclient.list()
    		for vol in vols:
    			if vol["name"] == docker_volume_name:
    				return {"Name": docker_volume_name,
    					"Mountpoing": ""}
      	  		else:        
            			msg = _LW("Can't find this volume '{0}' in "
                		      "OpenSDS").format(docker_volume_name)
            		LOG.warning(msg)
            		raise exceptions.NotFound(msg)

	def mount(self, docker_volume_name):
        	return "/mnt/docker/test002"

	def unmount(self, docker_volume_name):
        	return

	def check_exist(self, docker_volume_name):
    		return True
        	_, state = self._get_docker_volume(docker_volume_name)
        	LOG.info(_LI("Get docker volume {0} with state "
        	             "{1}").format(docker_volume_name, state))

        	if state == UNKNOWN:
        	    return False
        	return True
