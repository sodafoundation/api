# easyssh-proxy

[![GoDoc](https://godoc.org/github.com/appleboy/easyssh-proxy?status.svg)](https://godoc.org/github.com/appleboy/easyssh-proxy) 
[![Build Status](https://cloud.drone.io/api/badges/appleboy/easyssh-proxy/status.svg)](https://cloud.drone.io/appleboy/easyssh-proxy)
[![codecov](https://codecov.io/gh/appleboy/easyssh-proxy/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/easyssh-proxy) 
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/easyssh-proxy)](https://goreportcard.com/report/github.com/appleboy/easyssh-proxy) 
[![Sourcegraph](https://sourcegraph.com/github.com/appleboy/easyssh-proxy/-/badge.svg)](https://sourcegraph.com/github.com/appleboy/easyssh-proxy?badge)
[![Release](https://github-release-version.herokuapp.com/github/appleboy/easyssh-proxy/release.svg?style=flat)](https://github.com/appleboy/easyssh-proxy/releases/latest)

easyssh-proxy provides a simple implementation of some SSH protocol features in Go.

## Feature

This project is forked from [easyssh](https://github.com/hypersleep/easyssh) but add some features as the following.

* [x] Support plain text of user private key.
* [x] Support key path of user private key.
* [x] Support Timeout for the TCP connection to establish.
* [x] Support SSH ProxyCommand.

```
     +--------+       +----------+      +-----------+
     | Laptop | <-->  | Jumphost | <--> | FooServer |
     +--------+       +----------+      +-----------+

                         OR

     +--------+       +----------+      +-----------+
     | Laptop | <-->  | Firewall | <--> | FooServer |
     +--------+       +----------+      +-----------+
     192.168.1.5       121.1.2.3         10.10.29.68
```

## Usage:

You can see `ssh`, `scp`, `ProxyCommand` on `examples` folder.

### ssh

See [example/ssh/ssh.go](./example/ssh/ssh.go)

[embedmd]:# (example/ssh/ssh.go go)
```go
package main

import (
	"fmt"
	"time"

	"github.com/appleboy/easyssh-proxy"
)

func main() {
	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:   "appleboy",
		Server: "example.com",
		// Optional key or Password without either we try to contact your agent SOCKET
		//Password: "password",
		// Paste your source content of private key
		// Key: `-----BEGIN RSA PRIVATE KEY-----
		// MIIEpAIBAAKCAQEA4e2D/qPN08pzTac+a8ZmlP1ziJOXk45CynMPtva0rtK/RB26
		// 7XC9wlRna4b3Ln8ew3q1ZcBjXwD4ppbTlmwAfQIaZTGJUgQbdsO9YA==
		// -----END RSA PRIVATE KEY-----
		// `,
		KeyPath: "/Users/username/.ssh/id_rsa",
		Port:    "22",
		Timeout: 60 * time.Second,
	}

	// Call Run method with command you want to run on remote server.
	stdout, stderr, done, err := ssh.Run("ls -al", 60)
	// Handle errors
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	} else {
		fmt.Println("don is :", done, "stdout is :", stdout, ";   stderr is :", stderr)
	}

}
```

### scp

See [example/scp/scp.go](./example/scp/scp.go)

[embedmd]:# (example/scp/scp.go go)
```go
package main

import (
	"fmt"

	"github.com/appleboy/easyssh-proxy"
)

func main() {
	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     "appleboy",
		Server:   "example.com",
		Password: "123qwe",
		Port:     "22",
	}

	// Call Scp method with file you want to upload to remote server.
	// Please make sure the `tmp` floder exists.
	err := ssh.Scp("/root/source.csv", "/tmp/target.csv")

	// Handle errors
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	} else {
		fmt.Println("success")
	}
}
```

### SSH ProxyCommand

See [example/proxy/proxy.go](./example/proxy/proxy.go)

[embedmd]:# (example/proxy/proxy.go go /\tssh :=/ /\t}$/)
```go
	ssh := &easyssh.MakeConfig{
		User:    "drone-scp",
		Server:  "localhost",
		Port:    "22",
		KeyPath: "./tests/.ssh/id_rsa",
		Proxy: easyssh.DefaultConfig{
			User:    "drone-scp",
			Server:  "localhost",
			Port:    "22",
			KeyPath: "./tests/.ssh/id_rsa",
		},
	}
```

### SSH Stream Log

See [example/stream/stream.go](./example/stream/stream.go)

[embedmd]:# (example/stream/stream.go go /func/ /^}$/)
```go
func main() {
	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		Server:  "localhost",
		User:    "drone-scp",
		KeyPath: "./tests/.ssh/id_rsa",
		Port:    "22",
		Timeout: 60 * time.Second,
	}

	// Call Run method with command you want to run on remote server.
	stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream("for i in {1..5}; do echo ${i}; sleep 1; done; exit 2;", 60)
	// Handle errors
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	} else {
		// read from the output channel until the done signal is passed
		isTimeout := true
	loop:
		for {
			select {
			case isTimeout = <-doneChan:
				break loop
			case outline := <-stdoutChan:
				fmt.Println("out:", outline)
			case errline := <-stderrChan:
				fmt.Println("err:", errline)
			case err = <-errChan:
			}
		}

		// get exit code or command error.
		if err != nil {
			fmt.Println("err: " + err.Error())
		}

		// command time out
		if !isTimeout {
			fmt.Println("Error: command timeout")
		}
	}
}
```
