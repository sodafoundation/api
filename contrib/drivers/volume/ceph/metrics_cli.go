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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ceph/go-ceph/rados"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils/exec"
)

type MetricCli struct {
	conn         Conn
	RootExecuter exec.Executer
}

type Conn interface {
	ReadDefaultConfigFile() error
	Connect() error
	GetFSID() (fsid string, err error)
	Shutdown()
	MonCommand([]byte) ([]byte, string, error)
}

// Verify that *rados.Conn implements Conn correctly.
var _ Conn = &rados.Conn{}

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
		exec.NewRootExecuter(),
	}, nil
}

type CephMetricStats struct {
	Name        string
	Value       string
	Unit        string
	Const_Label map[string]string
	AggrType    string
	Var_Label   map[string]string
	Component   string
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

type cephClusterStats struct {
	Stats struct {
		TotalBytes      json.Number `json:"total_bytes"`
		TotalUsedBytes  json.Number `json:"total_used_bytes"`
		TotalAvailBytes json.Number `json:"total_avail_bytes"`
		TotalObjects    json.Number `json:"total_objects"`
	} `json:"stats"`
}

type cephPerfStat struct {
	PerfInfo []struct {
		ID    json.Number `json:"id"`
		Stats struct {
			CommitLatency json.Number `json:"commit_latency_ms"`
			ApplyLatency  json.Number `json:"apply_latency_ms"`
		} `json:"perf_stats"`
	} `json:"osd_perf_infos"`
}

type cephOSDDF struct {
	OSDNodes []struct {
		Name        json.Number `json:"name"`
		CrushWeight json.Number `json:"crush_weight"`
		Depth       json.Number `json:"depth"`
		Reweight    json.Number `json:"reweight"`
		KB          json.Number `json:"kb"`
		UsedKB      json.Number `json:"kb_used"`
		AvailKB     json.Number `json:"kb_avail"`
		Utilization json.Number `json:"utilization"`
		Variance    json.Number `json:"var"`
		Pgs         json.Number `json:"pgs"`
	} `json:"nodes"`

	Summary struct {
		TotalKB      json.Number `json:"total_kb"`
		TotalUsedKB  json.Number `json:"total_kb_used"`
		TotalAvailKB json.Number `json:"total_kb_avail"`
		AverageUtil  json.Number `json:"average_utilization"`
	} `json:"summary"`
}

type cephOSDDump struct {
	OSDs []struct {
		OSD json.Number `json:"osd"`
		Up  json.Number `json:"up"`
		In  json.Number `json:"in"`
	} `json:"osds"`
}

type cephHealthStats struct {
	Health struct {
		Summary []struct {
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
		} `json:"summary"`
		OverallStatus string `json:"overall_status"`
		Status        string `json:"status"`
		Checks        map[string]struct {
			Severity string `json:"severity"`
			Summary  struct {
				Message string `json:"message"`
			} `json:"summary"`
		} `json:"checks"`
	} `json:"health"`
	OSDMap struct {
		OSDMap struct {
			NumOSDs        json.Number `json:"num_osds"`
			NumUpOSDs      json.Number `json:"num_up_osds"`
			NumInOSDs      json.Number `json:"num_in_osds"`
			NumRemappedPGs json.Number `json:"num_remapped_pgs"`
		} `json:"osdmap"`
	} `json:"osdmap"`
	PGMap struct {
		NumPGs                  json.Number `json:"num_pgs"`
		WriteOpPerSec           json.Number `json:"write_op_per_sec"`
		ReadOpPerSec            json.Number `json:"read_op_per_sec"`
		WriteBytePerSec         json.Number `json:"write_bytes_sec"`
		ReadBytePerSec          json.Number `json:"read_bytes_sec"`
		RecoveringObjectsPerSec json.Number `json:"recovering_objects_per_sec"`
		RecoveringBytePerSec    json.Number `json:"recovering_bytes_per_sec"`
		RecoveringKeysPerSec    json.Number `json:"recovering_keys_per_sec"`
		CacheFlushBytePerSec    json.Number `json:"flush_bytes_sec"`
		CacheEvictBytePerSec    json.Number `json:"evict_bytes_sec"`
		CachePromoteOpPerSec    json.Number `json:"promote_op_per_sec"`
		DegradedObjects         json.Number `json:"degraded_objects"`
		MisplacedObjects        json.Number `json:"misplaced_objects"`
		PGsByState              []struct {
			Count  float64 `json:"count"`
			States string  `json:"state_name"`
		} `json:"pgs_by_state"`
	} `json:"pgmap"`
}

type cephMonitorStats struct {
	Health struct {
		Health struct {
			HealthServices []struct {
				Mons []struct {
					Name         string      `json:"name"`
					KBTotal      json.Number `json:"kb_total"`
					KBUsed       json.Number `json:"kb_used"`
					KBAvail      json.Number `json:"kb_avail"`
					AvailPercent json.Number `json:"avail_percent"`
					StoreStats   struct {
						BytesTotal json.Number `json:"bytes_total"`
						BytesSST   json.Number `json:"bytes_sst"`
						BytesLog   json.Number `json:"bytes_log"`
						BytesMisc  json.Number `json:"bytes_misc"`
					} `json:"store_stats"`
				} `json:"mons"`
			} `json:"health_services"`
		} `json:"health"`
		TimeChecks struct {
			Mons []struct {
				Name    string      `json:"name"`
				Skew    json.Number `json:"skew"`
				Latency json.Number `json:"latency"`
			} `json:"mons"`
		} `json:"timechecks"`
	} `json:"health"`
	Quorum []int `json:"quorum"`
}

type cephTimeSyncStatus struct {
	TimeChecks map[string]struct {
		Health  string      `json:"health"`
		Latency json.Number `json:"latency"`
		Skew    json.Number `json:"skew"`
	} `json:"time_skew_status"`
}

func (cli *MetricCli) CollectPoolMetrics() ([]CephMetricStats, error) {
	returnMap := []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "df",
		"detail": "detail",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}

	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph df detail")
		return nil, err
	}

	pool_stats := &cephPoolStats{}
	if err := json.Unmarshal(buf, pool_stats); err != nil {
		log.Errorf("unmarshal error: %v", err)
		return nil, err
	}

	for _, pool := range pool_stats.Pools {
		var_label := make(map[string]string)
		var_label["pool"] = pool.Name
		returnMap = append(returnMap, CephMetricStats{
			"used",
			pool.Stats.BytesUsed.String(),
			"bytes", const_label,
			"",
			var_label,
			"pool"})

		returnMap = append(returnMap, CephMetricStats{
			"raw_used",
			pool.Stats.RawBytesUsed.String(),
			"bytes", const_label,
			"",
			var_label,
			"pool"})

		returnMap = append(returnMap, CephMetricStats{
			"available",
			pool.Stats.MaxAvail.String(),
			"bytes",
			const_label,
			"",
			var_label,
			"pool"})

		returnMap = append(returnMap, CephMetricStats{
			"objects",
			pool.Stats.Objects.String(),
			"",
			const_label,
			"total",
			var_label,
			"pool"})

		returnMap = append(returnMap, CephMetricStats{
			"dirty_objects",
			pool.Stats.DirtyObjects.String(),
			"",
			const_label,
			"total",
			var_label,
			"pool"})

		returnMap = append(returnMap, CephMetricStats{
			"read", pool.Stats.ReadIO.String(),
			"",
			const_label,
			"total",
			var_label, "pool"})

		returnMap = append(returnMap, CephMetricStats{
			"read",
			pool.Stats.ReadBytes.String(),
			"bytes",
			const_label,
			"total",
			var_label, "pool"})

		returnMap = append(returnMap, CephMetricStats{
			"write",
			pool.Stats.WriteIO.String(),
			"", const_label,
			"total",
			var_label, "pool"})

		returnMap = append(returnMap, CephMetricStats{
			"write",
			pool.Stats.WriteBytes.String(),
			"bytes",
			const_label,
			"total",
			var_label, "pool"})
	}
	return returnMap, nil
}

func (cli *MetricCli) CollectClusterMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "df",
		"detail": "detail",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}
	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph df detail")
		return nil, err
	}
	cluster_stats := &cephClusterStats{}
	if err := json.Unmarshal(buf, cluster_stats); err != nil {
		log.Fatalf("unmarshal error: %v", err)
		return nil, err
	}
	returnMap = append(returnMap,
		CephMetricStats{
			"capacity",
			cluster_stats.Stats.TotalBytes.String(),
			"bytes",
			const_label,
			"",
			nil,
			"cluster"},
		CephMetricStats{
			"available",
			cluster_stats.Stats.TotalAvailBytes.String(),
			"bytes",
			const_label,
			"",
			nil,
			"cluster"},
		CephMetricStats{
			"used",
			cluster_stats.Stats.TotalUsedBytes.String(),
			"bytes",
			const_label,
			"",
			nil,
			"cluster"},
		CephMetricStats{
			"objects",
			cluster_stats.Stats.TotalObjects.String(),
			"",
			const_label,
			"", nil, "cluster"},
	)
	return returnMap, nil
}

func (cli *MetricCli) CollectPerfMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd perf",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}
	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph osd perf")
		return nil, err
	}
	osdPerf := &cephPerfStat{}
	if err := json.Unmarshal(buf, osdPerf); err != nil {
		log.Errorf("unmarshal failed")
		return nil, err
	}
	for _, perfStat := range osdPerf.PerfInfo {
		var_label := make(map[string]string)
		osdID, err := perfStat.ID.Int64()
		if err != nil {
			log.Errorf("when collecting ceph cluster metrics")
			return nil, err
		}
		var_label["osd"] = fmt.Sprintf("osd.%v", osdID)
		returnMap = append(returnMap,
			CephMetricStats{
				"perf_commit_latency",
				perfStat.Stats.CommitLatency.String(),
				"ms",
				const_label,
				"",
				var_label, ""},
			CephMetricStats{
				"perf_apply_latency",
				perfStat.Stats.ApplyLatency.String(),
				"ms",
				const_label,
				"",
				var_label, ""})
	}
	return returnMap, nil
}

func (cli *MetricCli) CollectOsddfMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd df",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}
	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph osd df")
		return nil, err
	}
	osddf := &cephOSDDF{}
	if err := json.Unmarshal(buf, osddf); err != nil {
		log.Errorf("unmarshal failed")
		return nil, err
	}
	for _, osd_df := range osddf.OSDNodes {
		var_label := make(map[string]string)
		var_label["osd"] = osd_df.Name.String()
		returnMap = append(returnMap,
			CephMetricStats{
				"crush_weight",
				osd_df.CrushWeight.String(),
				"", const_label,
				"",
				var_label, "osd"})
		returnMap = append(returnMap,
			CephMetricStats{
				"depth",
				osd_df.Depth.String(),
				"", const_label,
				"", var_label, "osd"})

		returnMap = append(returnMap,
			CephMetricStats{
				"reweight",
				osd_df.Reweight.String(),
				"", const_label,
				"", var_label, "osd"})

		osd_df_kb, _ := osd_df.KB.Float64()
		osd_df_bytes := fmt.Sprint(osd_df_kb * 1e3)
		returnMap = append(returnMap,
			CephMetricStats{
				"bytes",
				osd_df_bytes,
				"", const_label,
				"", var_label, "osd"})
		osd_df_kb_used, _ := osd_df.UsedKB.Float64()
		osd_df_bytes_used := fmt.Sprint(osd_df_kb_used * 1e3)
		returnMap = append(returnMap,
			CephMetricStats{
				"bytes_used",
				osd_df_bytes_used,
				"", const_label,
				"", var_label, "osd"})
		osd_df_kb_avail, _ := osd_df.AvailKB.Float64()
		osd_df_bytes_avail := fmt.Sprint(osd_df_kb_avail * 1e3)
		returnMap = append(returnMap,
			CephMetricStats{
				"bytes_avail",
				osd_df_bytes_avail,
				"", const_label,
				"", var_label, "osd"})
		returnMap = append(returnMap,
			CephMetricStats{
				"utilization",
				osd_df.Utilization.String(),
				"", const_label,
				"", var_label, "osd"})
		returnMap = append(returnMap,
			CephMetricStats{
				"var",
				osd_df.Variance.String(),
				"", const_label,
				"", var_label, "osd"})
		returnMap = append(returnMap,
			CephMetricStats{
				"pgs",
				osd_df.Pgs.String(),
				"", const_label,
				"", var_label, "osd"})
	}
	total_kb, _ := osddf.Summary.TotalKB.Float64()
	total_bytes := fmt.Sprint(total_kb * 1e3)
	returnMap = append(returnMap, CephMetricStats{
		"total",
		total_bytes,
		"bytes",
		const_label,
		"",
		nil, "osd"})

	total_used_kb, _ := osddf.Summary.TotalUsedKB.Float64()
	total_used_bytes := fmt.Sprint(total_used_kb * 1e3)
	returnMap = append(returnMap, CephMetricStats{
		"total_used",
		total_used_bytes,
		"bytes",
		const_label,
		"",
		nil, "osd"})

	total_avail_kb, _ := osddf.Summary.TotalAvailKB.Float64()
	total_avail_bytes := fmt.Sprint(total_avail_kb * 1e3)
	returnMap = append(returnMap, CephMetricStats{
		"total_avail",
		total_avail_bytes,
		"bytes",
		const_label,
		"",
		nil, "osd"})
	returnMap = append(returnMap, CephMetricStats{
		"average_utilization",
		osddf.Summary.AverageUtil.String(),
		"",
		const_label,
		"",
		nil, "osd"})

	return returnMap, nil
}

func (cli *MetricCli) CollectOsddumpMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd dump",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}
	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph osd perf")
		return nil, err
	}
	osd_dump := &cephOSDDump{}
	if err := json.Unmarshal(buf, osd_dump); err != nil {
		log.Errorf("unmarshal failed")
		return nil, err
	}
	var_label := make(map[string]string)
	var_label["osd"] = fmt.Sprintf("osd.%s", osd_dump.OSDs[0].OSD.String())
	returnMap = append(returnMap,
		CephMetricStats{"osd",
			osd_dump.OSDs[0].OSD.String(),
			"",
			const_label,
			"",
			var_label, ""},
		CephMetricStats{
			"osd_up",
			osd_dump.OSDs[0].Up.String(),
			"",
			const_label,
			"",
			var_label, ""},
		CephMetricStats{
			"osd_in",
			osd_dump.OSDs[0].In.String(),
			"",
			const_label,
			"",
			var_label, ""})
	return returnMap, nil
}

func (cli *MetricCli) CollectHealthMetrics() ([]CephMetricStats, error) {
	returnMap := []CephMetricStats{}
	constlabel := make(map[string]string)
	constlabel["cluster"] = "ceph"
	health_cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "status",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}
	buff, _, err := cli.conn.MonCommand(health_cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph status")
		return nil, err
	}
	health_stats := &cephHealthStats{}
	if err := json.Unmarshal(buff, health_stats); err != nil {
		log.Fatalf("unmarshal error: %v", err)
		return nil, err
	}

	returnMap = append(returnMap, CephMetricStats{
		"io_write",
		health_stats.PGMap.WriteOpPerSec.String(),
		"ops", constlabel,
		"",
		nil,
		"client"})

	returnMap = append(returnMap, CephMetricStats{
		"io_read",
		health_stats.PGMap.ReadBytePerSec.String(),
		"bytes", constlabel,
		"",
		nil,
		"client"})

	returnMap = append(returnMap, CephMetricStats{
		"io_read",
		(health_stats.PGMap.ReadOpPerSec.String() + health_stats.PGMap.WriteOpPerSec.String()),
		"ops",
		constlabel,
		"",
		nil,
		"client"})
	returnMap = append(returnMap, CephMetricStats{
		"io_write",
		health_stats.PGMap.WriteBytePerSec.String(),
		"bytes",
		constlabel,
		"",
		nil,
		"client"})
	returnMap = append(returnMap, CephMetricStats{
		"cache_flush_io",
		health_stats.PGMap.CacheFlushBytePerSec.String(),
		"bytes",
		constlabel,
		"",
		nil,
		""})
	returnMap = append(returnMap, CephMetricStats{
		"cache_evict_io",
		health_stats.PGMap.CacheEvictBytePerSec.String(),
		"bytes",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"cache_promote_io",
		health_stats.PGMap.CachePromoteOpPerSec.String(),
		"ops",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"degraded_objects",
		health_stats.PGMap.DegradedObjects.String(),
		"", constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"misplaced_objects",
		health_stats.PGMap.MisplacedObjects.String(),
		"",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"osds",
		health_stats.OSDMap.OSDMap.NumOSDs.String(),
		"",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"osds_up",
		health_stats.OSDMap.OSDMap.NumUpOSDs.String(),
		"",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"osds_in",
		health_stats.OSDMap.OSDMap.NumInOSDs.String(),
		"",
		constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"pgs_remapped",
		health_stats.OSDMap.OSDMap.NumRemappedPGs.String(),
		"", constlabel,
		"",
		nil,
		""})

	returnMap = append(returnMap, CephMetricStats{
		"total_pgs",
		health_stats.PGMap.NumPGs.String(),
		"",
		constlabel,
		"",
		nil,
		""})
	return returnMap, nil
}

func (cli *MetricCli) CollectMonitorsMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"

	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "status",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}

	buf, _, err := cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph status")
		return nil, err
	}

	mon_stats := &cephMonitorStats{}
	if err := json.Unmarshal(buf, mon_stats); err != nil {
		log.Fatalf("unmarshal error: %v", err)
		return nil, err
	}

	for _, healthService := range mon_stats.Health.Health.HealthServices {
		for _, monstat := range healthService.Mons {
			var_label := make(map[string]string)
			var_label["monitor"] = monstat.Name
			kbTotal, _ := monstat.KBTotal.Float64()
			bytesTotal_val := fmt.Sprintf("%f", kbTotal*1e3)
			returnMap = append(returnMap, CephMetricStats{
				"capacity",
				bytesTotal_val,
				"bytes", const_label,
				"",
				var_label,
				""})
			kbUsed, _ := monstat.KBUsed.Float64()
			bytesUsed_val := fmt.Sprintf("%f", kbUsed*1e3)
			returnMap = append(returnMap, CephMetricStats{
				"used",
				bytesUsed_val,
				"bytes", const_label,
				"",
				var_label,
				""})
			kbAvail, _ := monstat.KBAvail.Float64()
			bytesAvail_val := fmt.Sprintf("%f", kbAvail*1e3)
			returnMap = append(returnMap, CephMetricStats{
				"avail",
				bytesAvail_val,
				"bytes", const_label,
				"",
				var_label,
				""})
			returnMap = append(returnMap, CephMetricStats{
				"avail_percent",
				monstat.AvailPercent.String(),
				"", const_label,
				"",
				var_label,
				""})
			returnMap = append(returnMap, CephMetricStats{
				"store_capacity",
				monstat.StoreStats.BytesTotal.String(),
				"bytes", const_label,
				"",
				var_label,
				""})
			returnMap = append(returnMap, CephMetricStats{
				"store_sst",
				monstat.StoreStats.BytesSST.String(),
				"", const_label,
				"bytes",
				var_label,
				""})
			returnMap = append(returnMap, CephMetricStats{
				"store_log",
				monstat.StoreStats.BytesLog.String(),
				"bytes", const_label,
				"",
				var_label,
				""})
			returnMap = append(returnMap, CephMetricStats{
				"store_misc",
				monstat.StoreStats.BytesMisc.String(),
				"bytes", const_label,
				"",
				var_label,
				""})
		}
	}

	cmd, err = json.Marshal(map[string]interface{}{
		"prefix": "time-sync-status",
		"format": "json",
	})
	if err != nil {
		log.Errorf("cmd failed with %s\n", err)
		return nil, err
	}

	buf, _, err = cli.conn.MonCommand(cmd)
	if err != nil {
		log.Errorf("unable to collect data from ceph time-sync-status")
		return nil, err
	}

	timeStats := &cephTimeSyncStatus{}
	if err := json.Unmarshal(buf, mon_stats); err != nil {
		log.Fatalf("unmarshal error: %v", err)
		return nil, err
	}

	for monNode, tstat := range timeStats.TimeChecks {
		var_label := make(map[string]string)
		var_label["monitor"] = monNode
		returnMap = append(returnMap, CephMetricStats{
			"clock_skew",
			tstat.Skew.String(),
			"seconds",
			const_label,
			"",
			var_label,
			""})
		returnMap = append(returnMap, CephMetricStats{
			"latency",
			tstat.Latency.String(),
			"seconds",
			const_label,
			"",
			var_label,
			""})
		returnMap = append(returnMap, CephMetricStats{
			"quorum_count",
			fmt.Sprintf("%v", mon_stats.Quorum),
			"", const_label,
			"",
			var_label,
			""})
	}
	return returnMap, nil
}

func (c *MetricCli) execute(cmd ...string) (string, error) {
	return c.RootExecuter.Run(cmd[0], cmd[1:]...)
}
func (cli *MetricCli) CollectVolumeMetrics() ([]CephMetricStats, error) {
	var returnMap []CephMetricStats
	returnMap = []CephMetricStats{}
	const_label := make(map[string]string)
	const_label["cluster"] = "ceph"
	cmd := []string{"env", "rbd", "ls"}
	out, err := cli.execute(cmd...)
	if err != nil {
		log.Errorf("cmd.Run() failed with %s\n", err)
		err = nil

	}
	result := strings.Split(string(out), "\n")

	for i := 0; i < len(result); i++ {
		if result[i] != "" {
			command := []string{"env", "rbd", "info", result[i], "--format", "json"}
			command_out, _ := cli.execute(command...)
			command_output := strings.Split(string(command_out), ",")
			var output []string
			for j := 0; j < (len(command_output)); j++ {
				result := strings.Split(command_output[j], ":")
				output = append(output, result[1])

			}

			returnMap = append(returnMap, CephMetricStats{"name",
				result[i],
				"", const_label,
				"", nil,
				"volume"})
			returnMap = append(returnMap, CephMetricStats{"size",
				output[1],
				"bytes", const_label,
				"", nil,
				"volume"})
			returnMap = append(returnMap, CephMetricStats{"objects",
				output[2],
				"", const_label,
				"", nil,
				"volume"})
			returnMap = append(returnMap, CephMetricStats{"object_size",
				output[4],
				"bytes", const_label,
				"", nil,
				"volume"})
		}
	}

	return returnMap, nil
}

func (cli *MetricCli) CollectMetrics() ([]CephMetricStats, []string, error) {
	returnMap := []CephMetricStats{}
	var instance []string
	instanceID, _ := cli.conn.GetFSID()
	instance = append(instance, instanceID)
	instanceName := fmt.Sprintf("%s%s", "opensds-ceph-", instance[0])
	instance = append(instance, instanceName)
	// Collecting Pool Metrics
	pool_metric, _ := cli.CollectPoolMetrics()
	for i := range pool_metric {
		returnMap = append(returnMap, pool_metric[i])
	}
	// Collects Cluster Metrics
	cluster_metric, _ := cli.CollectClusterMetrics()
	for i := range cluster_metric {
		returnMap = append(returnMap, cluster_metric[i])
	}
	// Collects Performance Metrics
	perf_metric, _ := cli.CollectPerfMetrics()
	for i := range perf_metric {
		returnMap = append(returnMap, perf_metric[i])
	}
	// Collects OSD Metrics
	osd_df_metric, _ := cli.CollectOsddfMetrics()
	for i := range osd_df_metric {
		returnMap = append(returnMap, osd_df_metric[i])
	}
	// Collects OSD Dump Metrics
	osd_dump_metric, _ := cli.CollectOsddumpMetrics()
	for i := range osd_dump_metric {
		returnMap = append(returnMap, osd_dump_metric[i])
	}
	// Collects Ceph Health Metrics
	health_metrics, _ := cli.CollectHealthMetrics()
	for i := range health_metrics {
		returnMap = append(returnMap, health_metrics[i])
	}
	// Collects Ceph Monitor Metrics
	monitor_metrics, _ := cli.CollectMonitorsMetrics()
	for i := range monitor_metrics {
		returnMap = append(returnMap, monitor_metrics[i])
	}
	// Collects Ceph Volume Metrics
	volume_metrics, _ := cli.CollectVolumeMetrics()
	for i := range volume_metrics {
		returnMap = append(returnMap, volume_metrics[i])
	}
	return returnMap, instance, nil
}
