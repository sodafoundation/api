# Thin OpenSDS Installation Guide
## Prepare
Before you start, please make sure you have all the following:
- Ubuntu environment (suggest v16.04+).
- More than 30GB remaining disk.
- Make sure have access to the Internet.
- Some tools (`git`, `make` and `gcc`) prepared.

#### Bootstrap
First you need to download [bootstrap](https://github.com/opensds/opensds/blob/thin-opensds/install/devsds/bootstrap.sh) script and run it locally with root access.
```shell
curl -sSL https://raw.githubusercontent.com/opensds/opensds/thin-opensds/install/devsds/bootstrap.sh | sudo bash
```
If there is no error reported, you have all dependency packages installed.

#### Authentication configuration
Because the default authentication strategy is `noauth`, so if you want to enable multi-tenants feature or want to use Dashboard, please set the field `OPENSDS_AUTH_STRATEGY=keystone` in local.conf file:
```shell
cd $GOPATH/src/github.com/opensds/opensds
vi install/devsds/local.conf
```

#### Run osdsapiserver and osdsdock services of OpenSDS using following
```
cd $GOPATH/src/github.com/opensds/opensds && install/devsds/install.sh
```
If everything goes well, you will get some connection messages at the console output:
```shell
Execute commands blow to set up ENVs which are needed by OpenSDS CLI:
------------------------------------------------------------------
export OPENSDS_AUTH_STRATEGY=keystone
export OPENSDS_ENDPOINT=http://localhost:50040
export OS_AUTH_URL=http://10.10.3.150/identity
export OS_USERNAME=admin
export OS_PASSWORD=opensds@123
export OS_TENANT_NAME=admin
export OS_PROJECT_NAME=admin
export OS_USER_DOMAIN_ID=default
------------------------------------------------------------------
Enjoy it !!
```

## Testing

### Testing OpenSDS(Hotpot) using CLI
#### Config osdsctl tool
```shell
sudo cp $GOPATH/src/github.com/opensds/opensds/build/out/bin/osdsctl /usr/local/bin
```

#### Set some environment variables
```shell
export OPENSDS_ENDPOINT=http://127.0.0.1:50040
export OPENSDS_AUTH_STRATEGY=noauth # Set the value to keystone for multi-tenants.
```

If you choose keystone for authentication strategy, you need to execute different commands for logging in as different roles:
* For admin role
```shell
export OPENSDS_AUTH_STRATEGY=keystone
export OPENSDS_ENDPOINT=http://localhost:50040
export OS_AUTH_URL=http://$HOST_IP/identity
export OS_USERNAME=admin
export OS_PASSWORD=opensds@123
export OS_TENANT_NAME=admin
export OS_PROJECT_NAME=admin
export OS_USER_DOMAIN_ID=default
```

#### Create a volume using create volume API request
##### Endpoint

```$xslt
http://localhost:50040/v1beta/{{ tenant ID }}/block/volumes
```

##### Request body (update poolName and poolID)
```
{
  "name": "test-1",
  "description": "volume create test through api",
  "size": 1,
  "availabilityZone": "default",
  "profileId": "",
  "poolId":"{{ pool ID }]",
  "metadata": {
  	"poolName": "opensds-volumes"
  },
  "snapshotFromCloud": false
}
```

#### List all volumes

```
osdsctl volume list
```

#### Delete the volume
```
osdsctl volume delete <your_volume_id>
```


You can install [CSI Plugin using Helm](https://github.com/opensds/opensds-installer/blob/master/charts/OpenSDS%20Installation%20using%20Helm.md)
## Uninstall the local cluster
### OpenSDS(Hotpot)
It's also nice to uninstall the cluster with one command:
```
cd $GOPATH/src/github.com/opensds/opensds && install/devsds/uninstall.sh
```

If you want to destroy the cluster, please run the command below instead:
```
cd $GOPATH/src/github.com/opensds/opensds && install/devsds/uninstall.sh -purge
```


