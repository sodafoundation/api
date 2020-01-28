#!/bin/bash

import os
import subprocess, shlex
from cStringIO import StringIO
import time

class Hotpot(object):

    def __init__(self):
        self.stderr = None
        self.stdout = None

    def retry_check(self, function_to_retry):
        """
        Retries function passed in 5 times with a 1 second delay inbetween
        each try. Function should be passed in with lambda to call the
        function instead of just get its result.
        """
        check_delay = 1
        check_max_attempts = 5

        for retry_attempt in xrange(check_max_attempts):
            result = function_to_retry()
            if result == "available":
                return result

            print("Attempt {0} of 5".format(retry_attempt + 1))
            if retry_attempt < 4:
                print("Waiting {0} seconds before "
                    "rechecking".format(check_delay))
                time.sleep(check_delay)

        return result

    def get_index(self, output, arg):
        strip_list = []
        for line in output:
            if '|' in line:
                list = line.split('|')
                for item in list:
                    strip_list.append(item.strip())
                break
        return strip_list.index(arg)

    def get_status(self, command):
        strip_list = []
        stdout, stderr = self.run_command(command)
        output = StringIO(stdout)
        for line in output:
            if 'Status' not in line and '|' in line:
                linelist = line.split('|')
                for item in linelist:
                    strip_list.append(item.strip())

        return strip_list[self.get_index(StringIO(stdout), 'Status')]

    def get_id(self, stdout):
        strip_list = []
        output = StringIO(stdout)
        for line in output:
            if 'Id' in line and '|' in line:
                linelist = line.split('|')
                for item in linelist:
                    strip_list.append(item.strip())
        id = strip_list[2]
        return id

    def run_command(self, command):
        # path = "/root/gopath/src/github.com/opensds/opensds/build/out/bin/"
        print("Command: {0}".format(command))
        args = shlex.split(command)
        print("Args: {0}".format(args))
        process = subprocess.Popen(args,
              stdout=subprocess.PIPE,
              stderr=subprocess.PIPE)
        stdout, stderr = process.communicate()
        if stderr != '':
            raise Exception(stderr)
        else:
            print("Stdout: {0}".format(stdout))
        return stdout, stderr

    def hello_world(self):
        print "BASIC HOTPOT TESTING!"

    def pool_list(self):
        self.run_command('build/out/bin/osdsctl pool list')

    def profile_list(self):
        return self.run_command('osdsctl profile list')

    def dock_list(self):
        self.run_command('build/out/bin/osdsctl dock list')

    def volume_list(self):
        self.run_command('osdsctl volume list')

    def fileshare_list(self):
        self.run_command('osdsctl fileshare list')

    def version_list(self):
        self.run_command('osdsctl version list')

    def volume_create(self):
        stdout, stderr = self.run_command('build/out/bin/osdsctl volume create 1 --name=vol1')
        return self.get_id(stdout)

    def check_volume_status_available(self, volume_id):
        time.sleep(5)
        command = 'osdsctl volume list' + ' ' + '--id' + ' ' + volume_id
        status = self.retry_check(lambda: self.get_status(command))
        if status != "available":
            raise Exception("The volume status is: ",status)

    def volume_delete(self, volume_id):
        command = 'osdsctl volume delete' + ' ' + volume_id
        self.run_command(command)

    def fileshare_create(self):
        stdout, stderr = self.run_command('osdsctl fileshare create 2 --name=randomfile')
        return self.get_id(stdout)

    def check_fileshare_status_available(self, fileshare_id):
        time.sleep(5)
        command = 'osdsctl fileshare list' + ' ' + '--id' + ' ' + fileshare_id
        status = self.retry_check(lambda: self.get_status(command))
        if status != "available":
            raise Exception("The fileshare status is: ", status)
    def fileshare_create_acl(self):
        stdout, stderr = self.run_command('osdsctl fileshare create 2 --name=randomfile')
        return self.get_id(stdout)

    def fileshare_delete(self, fileshare_id):
        command = 'osdsctl fileshare delete' + ' ' + fileshare_id
        self.run_command(command)

    def profile_create_block(self):
        stdout, stderr = self.run_command('osdsctl profile create '"'"'{"name": "adefault_block_test", "description": "default policy", "storageType": "block"}'"'"'')
        return self.get_id(stdout)

    def profile_delete_block(self, block_profile_id):
        command = 'osdsctl profile delete' + ' ' + block_profile_id
        self.run_command(command)

    def profile_create_file(self):
        stdout, stderr = self.run_command('osdsctl profile create '"'"'{"name": "default_file_test", "description": "default policy", "storageType": "file", "provisioningProperties":{"ioConnectivity": {"accessProtocol": "NFS"},"DataStorage":{"StorageAccessCapability":["Read","Write","Execute"]}}}'"'"'')
        return self.get_id(stdout)

    def profile_delete_file(self, file_profile_id):
        command = 'osdsctl profile delete' + ' ' + file_profile_id
        self.run_command(command)
