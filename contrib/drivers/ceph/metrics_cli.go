// Copyright (c) 2019 The OpenSDS Authors.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.
package ceph

import (
	"encoding/json"
	"fmt"
	"github.com/ceph/go-ceph/rados"
)

//
type cephClusterStats struct {
	Pools []struct {
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Stats struct {
			BytesUsed    json.Number `json:"bytes_used"`
			RawBytesUsed json.Number `json:"raw_bytes_used"`
			MaxAvail     json.Number `json:"max_avail"`
			Objects      json.Number `json:"objects"`
			DirtyObjects json.Number `json:"dirty"`
			ReadIO       json.Number `json:"rd"`
			ReadBytes    json.Number `json:"rd_bytes"`
			WriteIO      json.Number `json:"wr"`
			WriteBytes   json.Number `json:"wr_bytes"`
		} `json:"stats"`
	} `json:"pools"`
}

func CollectMetrics(metricList []string, instanceID string) (map[string]map[string]string, error) {

	returnMap := make(map[string]map[string]string)
	var err error
	conn, err := rados.NewConn()
	if err != nil {
		fmt.Println(err)
	}

	err = conn.ReadDefaultConfigFile()
	if err != nil {
		fmt.Println(err)
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println(err)
	}

	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "df",
		"detail": "detail",
		"format": "json",
	})
	if err != nil {
		// panic! because ideally in no world this hard-coded input
		// should fail.
		panic(err)
	}

	buf, _, err := conn.MonCommand(cmd)
	if err != nil {
	}
	st := &cephClusterStats{}
	if err := json.Unmarshal(buf, st); err != nil {
		fmt.Printf("error")
		//return
	}

	//fmt.Printf("Command Output: %v",st)
	for _, pool := range st.Pools {
		miniarray := make(map[string]string)
		//miniarray[]=
		// UsedBytes tracks the amount of bytes currently allocated for the pool
		miniarray["pool_used_bytes"] = pool.Stats.BytesUsed.String()
		// RawUsedBytes tracks the amount of raw bytes currently used for the pool.
		miniarray["pool_raw_used_bytes"] = pool.Stats.RawBytesUsed.String()
		// MaxAvail tracks the amount of bytes currently free for the pool
		miniarray["pool_available_bytes"] = pool.Stats.MaxAvail.String()
		// Objects shows the no. of RADOS objects created within the pool.
		miniarray["pool_objects_total"] = pool.Stats.Objects.String()
		// DirtyObjects shows the no. of RADOS dirty objects in a cache-tier pool,
		// this doesn't make sense in a regular pool, see:
		// http://lists.ceph.com/pipermail/ceph-users-ceph.com/2015-April/000557.html
		miniarray["pool_dirty_objects_total"] = pool.Stats.DirtyObjects.String()
		// ReadIO tracks the read IO calls made for the images within each pool.
		miniarray["pool_read_total"] = pool.Stats.ReadIO.String()
		// Readbytes tracks the read throughput made for the images within each pool.
		miniarray["pool_read_bytes_total"] = pool.Stats.ReadBytes.String()
		// WriteIO tracks the write IO calls made for the images within each pool.
		miniarray["pool_write_total"] = pool.Stats.WriteIO.String()
		// WriteBytes tracks the write throughput made for the images within each pool.
		miniarray["pool_write_bytes_total"] = pool.Stats.WriteBytes.String()
		//For Pool Label, we have used pool name
		returnMap[pool.Name] = miniarray
	}

	conn.Shutdown()

	return returnMap, nil
}
