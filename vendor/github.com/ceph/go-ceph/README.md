# go-ceph - Go bindings for Ceph APIs

[![Build Status](https://travis-ci.org/ceph/go-ceph.svg)](https://travis-ci.org/ceph/go-ceph) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/ceph/go-ceph) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/ceph/go-ceph/master/LICENSE)

## Installation

    go get github.com/ceph/go-ceph

The native RADOS library and development headers are expected to be installed.

On debian systems (apt):
```sh
libcephfs-dev librbd-dev librados-dev
```

On rpm based systems (dnf, yum, etc):
```sh
libcephfs-devel librbd-devel librados-devel
```

## Documentation

Detailed documentation is available at
<http://godoc.org/github.com/ceph/go-ceph>.

### Connecting to a cluster

Connect to a Ceph cluster using a configuration file located in the default
search paths.

```go
conn, _ := rados.NewConn()
conn.ReadDefaultConfigFile()
conn.Connect()
```

A connection can be shutdown by calling the `Shutdown` method on the
connection object (e.g. `conn.Shutdown()`). There are also other methods for
configuring the connection. Specific configuration options can be set:

```go
conn.SetConfigOption("log_file", "/dev/null")
```

and command line options can also be used using the `ParseCmdLineArgs` method.

```go
args := []string{ "--mon-host", "1.1.1.1" }
err := conn.ParseCmdLineArgs(args)
```

For other configuration options see the full documentation.

### Object I/O

Object in RADOS can be written to and read from with through an interface very
similar to a standard file I/O interface:

```go
// open a pool handle
ioctx, err := conn.OpenIOContext("mypool")

// write some data
bytesIn := []byte("input data")
err = ioctx.Write("obj", bytesIn, 0)

// read the data back out
bytesOut := make([]byte, len(bytesIn))
_, err := ioctx.Read("obj", bytesOut, 0)

if !bytes.Equal(bytesIn, bytesOut) {
    fmt.Println("Output is not input!")
}
```

### Pool maintenance

The list of pools in a cluster can be retreived using the `ListPools` method
on the connection object. On a new cluster the following code snippet:

```go
pools, _ := conn.ListPools()
fmt.Println(pools)
```

will produce the output `[data metadata rbd]`, along with any other pools that
might exist in your cluster. Pools can also be created and destroyed. The
following creates a new, empty pool with default settings.

```go
conn.MakePool("new_pool")
```

Deleting a pool is also easy. Call `DeletePool(name string)` on a connection object to
delete a pool with the given name. The following will delete the pool named
`new_pool` and remove all of the pool's data.

```go
conn.DeletePool("new_pool")
```

# Development

```
docker run --rm -it --net=host \
  --device /dev/fuse --cap-add SYS_ADMIN --security-opt apparmor:unconfined \
  -v ${PWD}:/go/src/github.com/ceph/go-ceph:z \
  -v /home/nwatkins/src/ceph/build:/home/nwatkins/src/ceph/build:z \
  -e CEPH_CONF=/home/nwatkins/src/ceph/build/ceph.conf \
  ceph-golang
```

Run against a `vstart.sh` cluster without installing Ceph:

```
export CGO_CPPFLAGS="-I/ceph/src/include"
export CGO_LDFLAGS="-L/ceph/build/lib"
go build
```

## Contributing

Contributions are welcome & greatly appreciated, every little bit helps. Make code changes via Github pull requests:

- Fork the repo and create a topic branch for every feature/fix. Avoid
  making changes directly on master branch.
- All incoming features should be accompanied with tests.
- Make sure that you run `go fmt` before submitting a change
  set. Alternatively the Makefile has a flag for this, so you can call
  `make fmt` as well.
- The integration tests can be run in a docker container, for this run:

```
make test-docker
```
