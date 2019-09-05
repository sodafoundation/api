// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	dialTimeout    = 10 * time.Second
	requestTimeout = 200 * time.Second
)

func GetValueByKeyFromDB(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"192.168.20.123:62379"},
	})
	if err != nil {
		log.Fatal(err)
		// can define the below one in constants file and import
		return "OPERATION_FAILED"
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	//GetallKeys(ctx, kv, key)
	GetValueByKey(ctx, kv, key)
	// GetSingleValueDemo(ctx, kv)
	// GetMultipleValuesWithPaginationDemo(ctx, kv)
	// WatchDemo(ctx, cli, kv)
	// LeaseDemo(ctx, cli, kv)
	return "Done"
}

func GetallKeys(ctx context.Context, kv clientv3.KV, key string) string {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(3),
	}

	gr, _ := kv.Get(ctx, "key", opts...)

	fmt.Println("--- First page ---")
	for _, item := range gr.Kvs {
		fmt.Println(string(item.Key), string(item.Value))
	}
	return "Done2"
}

func GetKeys(ctx context.Context, kv clientv3.KV, key string) string {
	fmt.Println("*** GetKeys()")
	// Delete all keys
	// kv.Delete(ctx, "key", clientv3.WithPrefix())
	gr, err := kv.Get(ctx, "v1beta/file/shares/e93b4c0934da416eb9c8d120c5d04d96")
	if err != nil {
		log.Fatal(err)
		// can define the below one in constants file and import
		return "OPERATION_FAILED"
	}

	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	return string(gr.Kvs[0].Value)
}

func GetValueByKey(ctx context.Context, kv clientv3.KV, key string) string {
	fmt.Println("*** GetValueByKey()")
	// Delete all keys
	// kv.Delete(ctx, "key", clientv3.WithPrefix())

	// // Insert a key value
	// pr, err := kv.Put(ctx, "key", "444")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rev := pr.Header.Revision

	// fmt.Println("Revision:", rev)

	// gr, err := kv.Get(ctx, "v1beta/file/shares/e93b4c0934da416eb9c8d120c5d04d96/f2ab9308-f208-40c6-bb1f-6fbfa8bf14b5")
	//gr, err := kv.Get(ctx, "v1beta/file/shares")
	// 09b1c6e4-9dac-46cc-bb09-54795a354a79
	//gr, err := kv.Get(ctx, "09b1c6e4-9dac-46cc-bb09-54795a354a79")

	gr, err := kv.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
		// can define the below one in constants file and import
		return "OPERATION_FAILED"
	}

	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	return string(gr.Kvs[0].Value)
	// // Modify the value of an existing key (create new revision)
	// kv.Put(ctx, "key", "555")

	// gr, _ = kv.Get(ctx, "key")
	// fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// // Get the value of the previous revision
	// gr, _ = kv.Get(ctx, "key", clientv3.WithRev(rev))
	// fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
}