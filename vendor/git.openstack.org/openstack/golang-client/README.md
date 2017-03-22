OpenStack Golang Client
=======================

NOTE(dtroyer) Apr 2015: This repo is under heavy revision as it is being revived.

`openstack/golang-client` is an implementation of [OpenStack]
(http://www.openstack.org/) API client in [Go language](http://golang.org).
The code follows OpenStack licensing and uses its CI infrastructure
for hosting.  It currently implements [Identity Service v2] 
(http://docs.openstack.org/api/openstack-identity-service/2.0/content/) 
and [Object Storage v1] 
(http://docs.openstack.org/api/openstack-object-storage/1.0/content/).

The initial focus is on building a solid core REST Session and OpenStack
authentication on which to build the usual API interfaces.  The architecture
if the `Session` and authentication is similar to that used in the current
Python Keystone client library: The `Session` object contains the HTTP
interface methods and an authentication object that provides access to
the auth token and service catalog.

Current State
-------------
Code maturity is considered experimental.

* The new Session object is functional and used by most of the code now.
* The examples work.
* The image tests work.
* The obejct store tests do not work.
* identity/v2/auth.go is now unused, will be kept around for a short time
  for easier reference.

Installation
------------
Use `go get git.openstack.org/openstack/golang-client`.  Or alternatively,
download or clone the repository.

The lib was developed and tested on go 1.3. No external dependencies, so far.

Examples
--------
The examples directory contains examples for using the SDK using
real world working code. Each example starts with a two digit number followed
by a name (e.g., `00-authentication.go`). If you have a `config.json` file in the
examples directory following the format of `config.json.dist` the example can be
executed using `go run [example name] setup.go`. Or, all the examples can be
executed running the script `run-all.sh` from the examples directory.

Testing
-------
There are two types of test files.  The `*_test.go` are standard
golang unit test files.  The examples can be run as integration tests.

The tests were written against the [OpenStack API specifications]
(http://docs.openstack.org/api/api-specs.html).
The integration test were successful against the following:

- [DevStack](http://devstack.org)

If you use another provider and successfully completed the tests, please email
the maintainer(s) so your service can be mentioned here.  Alternatively, if you
are a service provider and can arrange a free (temporary) account, a quick test
can be arranged.

License
-------
Apache v2.

Contributing
------------
The code repository utilizes the OpenStack CI infrastructure.
Please use the [recommended workflow]
(http://docs.openstack.org/infra/manual/developers.html#development-workflow).  If you are not a member yet,
please consider joining as an [OpenStack contributor]
(http://docs.openstack.org/infra/manual/developers.html).  If you have questions or
comments, you can email the maintainer(s).

Coding Style
------------
The source code is automatically formatted to follow `go fmt` by the [IDE]
(https://code.google.com/p/liteide/).  And where pragmatic, the source code
follows this general [coding style]
(http://slamet.neocities.org/coding-style.html).
