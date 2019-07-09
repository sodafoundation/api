// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package ceph

import (
	"strconv"
	"time"

	"github.com/opensds/opensds/pkg/model"
)

// Supported metrics
var data = `
resources:
  - resource: pool
    metrics:
      - pool_used_bytes
      - pool_raw_used_bytes
      - pool_available_bytes
      - pool_objects_total
      - pool_dirty_objects_total
      - pool_read_total
      - pool_read_bytes_total
      - pool_write_total
      - pool_write_bytes_total
  - resource: cluster
    metrics:
      - cluster_capacity_bytes
      - cluster_used_bytes
      - cluster_available_bytes
      - cluster_objects
  - resource: osd
    metrics:
      - osd_perf_commit_latency
      - osd_perf_apply_latency
      - osd_crush_weight
      - osd_depth
      - osd_reweight
      - osd_bytes
      - osd_used_bytes
      - osd_avail_bytes
      - osd_utilization
      - osd_variance
      - osd_pgs
      - osd_total_bytes
      - osd_total_used_bytes
      - osd_total_avail_bytes
      - osd_average_utilization
  - resource: health
    metrics:
      - health_status
      - total_pgs
      - active_pgs
      - scrubbing_pgs
      - deep_scrubbing_pgs
      - recovering_pgs
      - recovery_wait_pgs
      - backfilling_pgs
      - forced_recovery_pgs
      - forced_backfill_pgs
      - down_pgs
      - slow_requests          
      - degraded_pgs
      - stuck_degraded_pgs
      - unclean_pgs
      - stuck_unclean_pgs
      - undersized_pgs
      - stuck_undersized_pgs
      - stale_pgs
      - stuck_stale_pgs
      - peering_pgs
      - degraded_objects
      - misplaced_objects
      - osdmap_flag_full
      - osdmap_flag_pauserd
      - osdmap_flag_pausewr
      - osdmap_flag_noup
      - osdmap_flag_nodown
      - osdmap_flag_noin
      - osdmap_flag_noout
      - osdmap_flag_nobackfill
      - osdmap_flag_norecover
      - osdmap_flag_norebalance
      - osdmap_flag_noscrub
      - osdmap_flag_nodeep_scrub
      - osdmap_flag_notieragent
      - osds_down
      - osds_up
      - osds_in
      - osds
      - pgs_remapped
      - recovery_io_bytes
      - recovery_io_keys
      - recovery_io_objects
      - client_io_read_bytes
      - client_io_write_bytes
      - client_io_ops
      - client_io_read_ops
      - client_io_write_ops
      - cache_flush_io_bytes
      - cache_evict_io_bytes
      - cache_promote_io_ops
  - resource: monitor
    metrics:
      - name
      - kb_total
      - kb_used
      - kb_avail
      - avail_percent
      - bytes_total
      - bytes_sst
      - bytes_log
      - bytes_misc
      - skew
      - latency
      - quorum
  - resource: volume
    metrics:
      - volume_name
      - volume_size_bytes
      - volume_objects
      - volume_objects_size_bytes
`

type MetricDriver struct {
	cli *MetricCli
}

func getCurrentUnixTimestamp() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}

func (d *MetricDriver) CollectMetrics() ([]*model.MetricSpec, error) {

	metricMap, instance, err := d.cli.CollectMetrics()
	var tempMetricArray []*model.MetricSpec
	for i := 0; i < len(metricMap); i++ {
		val, _ := strconv.ParseFloat(metricMap[i].Value, 64)
		associatorMap := make(map[string]string)
		for k := range metricMap[i].Const_Label {
			associatorMap[k] = metricMap[i].Const_Label[k]
		}
		if metricMap[i].Var_Label != nil {
			for k := range metricMap[i].Var_Label {
				associatorMap[k] = metricMap[i].Var_Label[k]
			}
		}
		metricValue := &model.Metric{
			Value:     val,
			Timestamp: getCurrentUnixTimestamp(),
		}
		metricValues := make([]*model.Metric, 0)
		metricValues = append(metricValues, metricValue)
		metric := &model.MetricSpec{
			InstanceID:   instance[0],
			InstanceName: instance[1],
			Job:          "ceph",
			Labels:       associatorMap,
			Component:    metricMap[i].Component,
			Name:         metricMap[i].Name,
			Unit:         metricMap[i].Unit,
			AggrType:     metricMap[i].AggrType,
			MetricValues: metricValues,
		}
		tempMetricArray = append(tempMetricArray, metric)
	}
	metricArray := tempMetricArray
	return metricArray, err
}

func (d *MetricDriver) Setup() error {
	cli, err := NewMetricCli()
	if err != nil {
		return err
	}
	d.cli = cli
	return nil
}

func (d *MetricDriver) Teardown() error {
	d.cli.conn.Shutdown()
	return nil
}
