openstack
=========

`openstack` is the API to an OpenStack cloud.

* `session.go` - A Session object that encapsulates the HTTP REST handler
  and authentication and logging

* `auth.go` - The basic authentication interface

* `auth-password.go` - Implements password authentication (v2 only at present)

* `auth-token.go` - The returned token objects
