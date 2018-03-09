Here is a tutorial guiding users and new contributors to get familiar with [OpenSDS](https://github.com/opensds/opensds) by installing a simple local cluster and managing lvm device.

## Prepare
Before you start, please make sure you have all stuffs below ready:
- Ubuntu environment (suggest v16.04+).
- More than 30GB remaining disk.
- Make sure have access to the Internet.
- Some tools (```git```, ```docker```) prepared.

## Step by Step Installation
### Bootstrap
Firstly, you need to download [bootstrap](https://github.com/opensds/opensds/blob/development/script/cluster/bootstrap.sh) script and run it locally with root access.
```shell
curl -sSL https://raw.githubusercontent.com/opensds/opensds/master/script/cluster/bootstrap.sh | sudo bash
```
If there is no error report, you'll have all dependency packages installed.

### Run all services in one command!
Don't be surprised, you could do it in one command:
```
cd $GOPATH/src/github.com/opensds/opensds && script/devsds/install.sh
```

## Testing
### Config osdsctl tool.
```
sudo cp build/out/bin/osdsctl /usr/local/bin

export OPENSDS_ENDPOINT=http://127.0.0.1:50040
osdsctl pool list
```

### Create default profile.
```
osdsctl profile create '{"name": "default", "description": "default policy"}'
```

### Create a volume.
```
osdsctl volume create 1 --name=test-001
```

### List all volumes.
```
osdsctl volume list
```

### Delete the volume.
```
osdsctl volume delete <your_volume_id>
```

## Uninstall the local cluster
It's also cool to uninstall the cluster in one command:
```
cd $GOPATH/src/github.com/opensds/opensds && script/devsds/uninstall.sh
```

Hope you could enjoy it, and more suggestions are welcomed!
