// Package easyssh provides a simple implementation of some SSH protocol
// features in Go. You can simply run a command on a remote server or get a file
// even simpler than native console SSH client. You don't need to think about
// Dials, sessions, defers, or public keys... Let easyssh think about it!
package easyssh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var defaultTimeout = 60 * time.Second

type (
	// MakeConfig Contains main authority information.
	// User field should be a name of user on remote server (ex. john in ssh john@example.com).
	// Server field should be a remote machine address (ex. example.com in ssh john@example.com)
	// Key is a path to private key on your local machine.
	// Port is SSH server port on remote machine.
	// Note: easyssh looking for private key in user's home directory (ex. /home/john + Key).
	// Then ensure your Key begins from '/' (ex. /.ssh/id_rsa)
	MakeConfig struct {
		User     string
		Server   string
		Key      string
		KeyPath  string
		Port     string
		Password string
		Timeout  time.Duration
		Proxy    DefaultConfig
	}

	// DefaultConfig for ssh proxy config
	DefaultConfig struct {
		User     string
		Server   string
		Key      string
		KeyPath  string
		Port     string
		Password string
		Timeout  time.Duration
	}
)

// returns ssh.Signer from user you running app home path + cutted key path.
// (ex. pubkey,err := getKeyFile("/.ssh/id_rsa") )
func getKeyFile(keypath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

// returns *ssh.ClientConfig and io.Closer.
// if io.Closer is not nil, io.Closer.Close() should be called when
// *ssh.ClientConfig is no longer used.
func getSSHConfig(config DefaultConfig) (*ssh.ClientConfig, io.Closer) {
	var sshAgent io.Closer

	// auths holds the detected ssh auth methods
	auths := []ssh.AuthMethod{}

	// figure out what auths are requested, what is supported
	if config.Password != "" {
		auths = append(auths, ssh.Password(config.Password))
	}
	if config.KeyPath != "" {
		if pubkey, err := getKeyFile(config.KeyPath); err != nil {
			log.Printf("getKeyFile error: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(pubkey))
		}
	}

	if config.Key != "" {
		if signer, err := ssh.ParsePrivateKey([]byte(config.Key)); err != nil {
			log.Printf("ssh.ParsePrivateKey: %v\n", err)
		} else {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
	}

	return &ssh.ClientConfig{
		Timeout:         config.Timeout,
		User:            config.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, sshAgent
}

// Connect to remote server using MakeConfig struct and returns *ssh.Session
func (ssh_conf *MakeConfig) Connect() (*ssh.Session, error) {
	var client *ssh.Client
	var err error

	targetConfig, closer := getSSHConfig(DefaultConfig{
		User:     ssh_conf.User,
		Key:      ssh_conf.Key,
		KeyPath:  ssh_conf.KeyPath,
		Password: ssh_conf.Password,
		Timeout:  ssh_conf.Timeout,
	})
	if closer != nil {
		defer closer.Close()
	}

	// Enable proxy command
	if ssh_conf.Proxy.Server != "" {
		proxyConfig, closer := getSSHConfig(DefaultConfig{
			User:     ssh_conf.Proxy.User,
			Key:      ssh_conf.Proxy.Key,
			KeyPath:  ssh_conf.Proxy.KeyPath,
			Password: ssh_conf.Proxy.Password,
			Timeout:  ssh_conf.Proxy.Timeout,
		})
		if closer != nil {
			defer closer.Close()
		}

		proxyClient, err := ssh.Dial("tcp", net.JoinHostPort(ssh_conf.Proxy.Server, ssh_conf.Proxy.Port), proxyConfig)
		if err != nil {
			return nil, err
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(ssh_conf.Server, ssh_conf.Port))
		if err != nil {
			return nil, err
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(ssh_conf.Server, ssh_conf.Port), targetConfig)
		if err != nil {
			return nil, err
		}

		client = ssh.NewClient(ncc, chans, reqs)
	} else {
		client, err = ssh.Dial("tcp", net.JoinHostPort(ssh_conf.Server, ssh_conf.Port), targetConfig)
		if err != nil {
			return nil, err
		}
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Stream returns one channel that combines the stdout and stderr of the command
// as it is run on the remote machine, and another that sends true when the
// command is done. The sessions and channels will then be closed.
func (ssh_conf *MakeConfig) Stream(command string, timeout ...time.Duration) (<-chan string, <-chan string, <-chan bool, <-chan error, error) {
	// continuously send the command's output over the channel
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan bool)
	errChan := make(chan error)

	// connect to remote host
	session, err := ssh_conf.Connect()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}
	// defer session.Close()
	// connect to both outputs (they are of type io.Reader)
	outReader, err := session.StdoutPipe()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}
	errReader, err := session.StderrPipe()
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}
	err = session.Start(command)
	if err != nil {
		return stdoutChan, stderrChan, doneChan, errChan, err
	}

	// combine outputs, create a line-by-line scanner
	stdoutReader := io.MultiReader(outReader)
	stderrReader := io.MultiReader(errReader)
	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrScanner := bufio.NewScanner(stderrReader)

	go func(stdoutScanner, stderrScanner *bufio.Scanner, stdoutChan, stderrChan chan string, doneChan chan bool, errChan chan error) {
		defer close(stdoutChan)
		defer close(stderrChan)
		defer close(doneChan)
		defer close(errChan)
		defer session.Close()

		// default timeout value
		executeTimeout := defaultTimeout
		if len(timeout) > 0 {
			executeTimeout = timeout[0]
		}
		timeoutChan := time.After(executeTimeout)
		res := make(chan struct{}, 1)
		var resWg sync.WaitGroup
		resWg.Add(2)

		go func() {
			for stdoutScanner.Scan() {
				stdoutChan <- stdoutScanner.Text()
			}
			resWg.Done()
		}()

		go func() {
			for stderrScanner.Scan() {
				stderrChan <- stderrScanner.Text()
			}
			resWg.Done()
		}()

		go func() {
			resWg.Wait()
			// close all of our open resources
			res <- struct{}{}
		}()

		select {
		case <-res:
			errChan <- session.Wait()
			doneChan <- true
		case <-timeoutChan:
			stderrChan <- "Run Command Timeout!"
			errChan <- nil
			doneChan <- false
		}
	}(stdoutScanner, stderrScanner, stdoutChan, stderrChan, doneChan, errChan)

	return stdoutChan, stderrChan, doneChan, errChan, err
}

// Run command on remote machine and returns its stdout as a string
func (ssh_conf *MakeConfig) Run(command string, timeout ...time.Duration) (outStr string, errStr string, isTimeout bool, err error) {
	stdoutChan, stderrChan, doneChan, errChan, err := ssh_conf.Stream(command, timeout...)
	if err != nil {
		return outStr, errStr, isTimeout, err
	}
	// read from the output channel until the done signal is passed
loop:
	for {
		select {
		case isTimeout = <-doneChan:
			break loop
		case outline := <-stdoutChan:
			if outline != "" {
				outStr += outline + "\n"
			}
		case errline := <-stderrChan:
			if errline != "" {
				errStr += errline + "\n"
			}
		case err = <-errChan:
		}
	}
	// return the concatenation of all signals from the output channel
	return outStr, errStr, isTimeout, err
}

// Scp uploads sourceFile to remote machine like native scp console app.
func (ssh_conf *MakeConfig) Scp(sourceFile string, etargetFile string) error {
	session, err := ssh_conf.Connect()

	if err != nil {
		return err
	}
	defer session.Close()

	targetFile := filepath.Base(etargetFile)

	src, srcErr := os.Open(sourceFile)

	if srcErr != nil {
		return srcErr
	}

	srcStat, statErr := src.Stat()

	if statErr != nil {
		return statErr
	}

	go func() {
		w, err := session.StdinPipe()

		if err != nil {
			return
		}
		defer w.Close()

		fmt.Fprintln(w, "C0644", srcStat.Size(), targetFile)

		if srcStat.Size() > 0 {
			io.Copy(w, src)
			fmt.Fprint(w, "\x00")
		} else {
			fmt.Fprint(w, "\x00")
		}
	}()

	return session.Run(fmt.Sprintf("scp -tr %s", etargetFile))
}
