package main

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var Timeout = 5

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	_, err := cli.Put(ctx, "foo", "bar")
	cancel()
	if err != nil {
		fmt.Println(err)
	}
}
