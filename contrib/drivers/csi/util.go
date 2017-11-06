package csi

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (

	// NameSpace for CSI
	NameSpace = "csi"

	// CSIEndPoint environment variable name
	CSIEndPoint = "CSI_ENDPOINT"
)

// getProtoandAdd return protocal and address
func getProtoandAdd(target string) (string, string) {
	reg := `(?i)^((?:(?:tcp|udp|ip)[46]?)|` + `(?:unix(?:gram|packet)?))://(.+)$`
	t := regexp.MustCompile(reg).FindStringSubmatch(target)
	return t[1], t[2]
}

// GetCSIEndPoint from environment variable
func GetCSIEndPoint() (string, error) {
	// example: CSI_ENDPOINT=unix://path/to/unix/domain/socket.sock
	csiEndPoint := os.Getenv(CSIEndPoint)
	csiEndPoint = strings.TrimSpace(csiEndPoint)

	if csiEndPoint == "" {
		err := errors.New("CSIEndPoint is empty")
		log.Fatalf("%v", err)
		return csiEndPoint, err
	}

	return csiEndPoint, nil
}

// GetCSIEndPointListener from endpoint
func GetCSIEndPointListener() (net.Listener, error) {
	target, err := GetCSIEndPoint()
	if err != nil {
		return nil, err
	}

	proto, addr := getProtoandAdd(target)
	return net.Listen(proto, addr)
}

// GetCSIClientConn from endpoint
func GetCSIClientConn() (*grpc.ClientConn, error) {
	// Get parameters for grpc
	ctx := context.Background()
	target, err := GetCSIEndPoint()
	if err != nil {
		return nil, err
	}

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(
			func(target string, timeout time.Duration) (net.Conn, error) {
				proto, addr := getProtoandAdd(target)
				return net.DialTimeout(proto, addr, timeout)
			}),
	}

	// Set up a connection to the server
	return grpc.DialContext(ctx, target, dialOpts...)
}
