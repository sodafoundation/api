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
	"github.com/ceph/go-ceph/rados"
	log "github.com/golang/glog"
)

type MetricCli struct {
	conn *rados.Conn
}

func NewMetricCli() (*MetricCli, error) {

	conn, err := rados.NewConn()
	if err != nil {
		log.Error("when connecting to rados:", err)
		return nil, err
	}

	err = conn.ReadDefaultConfigFile()
	if err != nil {
		log.Error("file ReadDefaultConfigFile can't read", err)
		return nil, err
	}

	err = conn.Connect()
	if err != nil {
		log.Error("when connecting to ceph cluster:", err)
		return nil, err
	}

	return &MetricCli{
		conn,
	}, nil
}

type CephMetricStats struct {
	Name        string
	Value       string
	Unit        string
	Const_Label string
	AggrType    string
	Var_Label   string
}

type cephPoolStats struct {
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

func (cli *MetricCli) CollectMetrics(metricList []string, instanceID string) ([]CephMetricStats, error) {
	returnMap := []CephMetricStats{}
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "df",
		"detail": "detail",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
	}

	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
	}

	pool_stats := &cephPoolStats{}
	if err := json.Unmarshal(buf, pool_stats); err != nil {
		log.Errorf("unmarshal error: %v", err)
	}

	for _, pool := range pool_stats.Pools {

		for _, element := range metricList {
			switch element {
			case "pool_used_bytes":
				returnMap = append(returnMap, CephMetricStats{
					"used",
					pool.Stats.BytesUsed.String(),
					"bytes", "ceph",
					"",
					pool.Name})

			case "pool_raw_used_bytes":
				returnMap = append(returnMap, CephMetricStats{
					"raw_used",
					pool.Stats.RawBytesUsed.String(),
					"bytes", "ceph",
					"",
					pool.Name})

			case "pool_available_bytes":
				returnMap = append(returnMap, CephMetricStats{
					"available",
					pool.Stats.MaxAvail.String(),
					"bytes",
					"ceph",
					"",
					pool.Name})

			case "pool_objects_total":
				returnMap = append(returnMap, CephMetricStats{
					"objects",
					pool.Stats.Objects.String(),
					"",
					"ceph",
					"total",
					pool.Name})

			case "pool_dirty_objects_total":
				returnMap = append(returnMap, CephMetricStats{
					"dirty_objects",
					pool.Stats.DirtyObjects.String(),
					"",
					"ceph",
					"total",
					pool.Name})

			case "pool_read_total":
				returnMap = append(returnMap, CephMetricStats{
					"read", pool.Stats.ReadIO.String(),
					"",
					"ceph",
					"total",
					pool.Name})

			case "pool_read_bytes_total":
				returnMap = append(returnMap, CephMetricStats{
					"read",
					pool.Stats.ReadBytes.String(),
					"bytes",
					"ceph",
					"total",
					pool.Name})

			case "pool_write_total":
				returnMap = append(returnMap, CephMetricStats{
					"write",
					pool.Stats.WriteIO.String(),
					"", "ceph",
					"total",
					pool.Name})

			case "pool_write_bytes_total":
				returnMap = append(returnMap, CephMetricStats{
					"write",
					pool.Stats.WriteBytes.String(),
					"bytes",
					"ceph",
					"total",
					pool.Name})
			}
		}
	}
	return returnMap, nil
}
