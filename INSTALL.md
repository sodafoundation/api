#### Build

1. export GOPATH=$HOME/gopath

   export PATH=$HOME/gopath/bin:$PATH
   
   mkdir -p $HOME/gopath/src/github.com/opensds/
   
   cd $HOME/gopath/src/github.com/opensds
   
2. git clone https://github.com/opensds/opensds.git $HOME/gopath/src/github.com/opensds/

3. cd opensds (import necessary packages)

   go get github.com/opensds/opensds/cmd/osdsctl

   go get github.com/opensds/opensds/cmd/osdslet

   go get github.com/opensds/opensds/cmd/osdsdock
   
4. cd cmd/sdslet

   go build
   
5. cd cmd/sdsctl

   go buld

6. cd cmd/sdsdock

   go build
   
7. cp cmd/sdslet/sdslet /usr/local/bin

   cp cmd/sdsctl/sdsctl /usr/local/bin

   cp cmd/sdsdock/sdsdock /usr/local/bin

7. vim examples/config.json (config backend storage credential information)

   vim examples/dock_route.json (config dock route table in controller module)

   vim examples/dock_node.json (config dock node in dock module)

   sudo mkdir /etc/opensds

   sudo cp examples/*.json /etc/opensds/

8. sudo mkdir /var/log/opensds (create OpenSDS logging directory)

#### Run

* Start **sdsdock** with root access (for logging purpose)

```sh
sudo sdsdock //suppose the user has copied the compiled binary to /usr/local/bin
```

* Start **sdslet** with root access (for logging purpose)

```sh
sudo sdslet //suppose the user has copied the compiled binary to /usr/local/bin
```

* Run **sdsctl** for operations you want to perform. 

```sh
sdsctl --help //see what you can do with opensds
```

Currently sdsctl supports all the basic Cinder/Manila operations, for example if you want to 
create a 1GB volume from a Dell-EMC VMAX, which is connected to the OpenSDS underlay infra - 
OpenStack Cinder via its in-tree vmax cinder driver, using OpenSDS for an easy access:

```sh
sdsctl volume create 1 -n cinder-vmax-volume -b cinder
```
Viola !
