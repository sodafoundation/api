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
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/opensds/opensds/pkg/utils/exec"

	"github.com/opensds/opensds/pkg/model"
)

var pool_Label map[string]string = map[string]string{"cluster": "ceph", "pool": "rbd"}
var osd_label map[string]string = map[string]string{"cluster": "ceph", "osd": "osd.0"}
var cluster_label map[string]string = map[string]string{"cluster": "ceph"}
var health_label map[string]string = map[string]string{"cluster": "ceph"}
var volume_label map[string]string = map[string]string{"cluster": "ceph"}

var expected_data map[string]CephMetricStats = map[string]CephMetricStats{
	"pool_used_bytes":          {"used", "859", "bytes", nil, "", pool_Label, "pool"},
	"pool_raw_used_bytes":      {"raw_used", "859", "bytes", nil, "", pool_Label, "pool"},
	"pool_available_bytes":     {"available", "469501706240", "bytes", nil, "", pool_Label, "pool"},
	"pool_objects_total":       {"objects", "14", "", nil, "total", pool_Label, "pool"},
	"pool_dirty_objects_total": {"dirty_objects", "14", "", nil, "total", pool_Label, "pool"},
	"pool_read_total":          {"read", "145", "", nil, "total", pool_Label, "pool"},
	"pool_read_bytes_total":    {"read", "918304", "bytes", nil, "total", pool_Label, "pool"},
	"pool_write_total":         {"write", "1057", "", nil, "total", pool_Label, "pool"},
	"pool_write_bytes_total":   {"write", "16384", "bytes", nil, "total", pool_Label, "pool"},
	"cluster_capacity_bytes":   {"capacity", "494462976000", "bytes", nil, "", cluster_label, "cluster"},
	"cluster_used_bytes":       {"used", "238116864", "bytes", nil, "", cluster_label, "cluster"},
	"cluster_available_bytes":  {"available", "494224859136", "bytes", nil, "", cluster_label, "cluster"},
	"cluster_objects":          {"objects", "14", "", nil, "", cluster_label, "cluster"},
	"perf_commit_latency_ms":   {"perf_commit_latency", "0", "ms", nil, "", osd_label, ""},
	"perf_apply_latency_ms":    {"perf_apply_latency", "0", "ms", nil, "", osd_label, ""},
	"osd_crush_weight":         {"crush_weight", ":0", "", nil, "", osd_label, "osd"},
	"osd_depth":                {"depth", "2", "", nil, "", osd_label, "osd"},
	"osd_reweight":             {"reweight", "1.000000", "", nil, "", osd_label, "osd"},
	"osd_bytes":                {"bytes", "15717356000", "", nil, "", osd_label, "osd"},
	"osd_bytes_used":           {"bytes_used", "114624000", "", nil, "", osd_label, "osd"},
	"osd_bytes_avail":          {"bytes_avail", "15602732000", "", nil, "", osd_label, "osd"},
	"osd_utilization":          {"utilization", "0.729283", "", nil, "", osd_label, "osd"},
	"osd_var":                  {"var", "1.000000", "", nil, "", osd_label, "osd"},
	"osd_pgs":                  {"pgs", "102", "", nil, "", osd_label, "osd"},
	"osd_total_bytes":          {"total", "15717356000", "bytes", nil, "", health_label, "osd"},
	"osd_total_used_bytes":     {"total_used", "114624000", "bytes", nil, "", health_label, "osd"},
	"osd_total_avail_bytes":    {"total_avail", "15602732000", "bytes", nil, "", health_label, "osd"},
	"osd_average_utilization":  {"average_utilization", "0.729283", "", nil, "", health_label, "osd"},
	"osd":                      {"osd", "0", "", nil, "", osd_label, ""},
	"osd_up":                   {"osd_up", "1", "", nil, "", osd_label, ""},
	"osd_in":                   {"osd_in", "1", "", nil, "", osd_label, ""},
	"client_io_write_ops":      {"io_write", "0", "ops", nil, "", health_label, "client"},
	"client_io_read_bytes":     {"io_read", "0", "bytes", nil, "", health_label, "client"},
	"client_io_read_ops":       {"io_read", "0", "ops", nil, "", health_label, "client"},
	"client_io_write_bytes":    {"io_write", "0", "bytes", nil, "", health_label, "client"},
	"cache_flush_io_bytes":     {"cache_flush_io", "0", "bytes", nil, "", health_label, ""},
	"cache_evict_io_bytes":     {"cache_evict_io", "0", "bytes", nil, "", health_label, ""},
	"cache_promote_io_ops":     {"cache_promote_io", "0", "ops", nil, "", health_label, ""},
	"degraded_objects":         {"degraded_objects", "0", "", nil, "", health_label, ""},
	"misplaced_objects":        {"misplaced_objects", "0", "", nil, "", health_label, ""},
	"osds":                     {"osds", "1", "", nil, "", health_label, ""}, "osds_up": {"osds_up", "1", "", nil, "", health_label, ""}, "osds_in": {"osds_in", "1", "", nil, "", health_label, ""},
	"pgs_remapped": {"pgs_remapped", "0", "", nil, "", health_label, ""}, "total_pgs": {"total_pgs", "102", "", nil, "", health_label, ""}, "volume_name": {"name", "opensds-4c5cb264-50d1-4bfd-a663-dface9b669c9", "", nil, "", volume_label, "volume"},
	"volume_size_bytes":        {"size", "1073741824", "bytes", nil, "", volume_label, "volume"},
	"volume_objects":           {"objects", "1024", "", nil, "", volume_label, "volume"},
	"volume_object_size_bytes": {"object_size", "1048576", "bytes", nil, "", volume_label, "volume"}}

var expctdMetricList []string = []string{"pool_used_bytes", "pool_raw_used_bytes", "pool_available_bytes", "pool_objects_total", "pool_dirty_objects_total", "pool_read_total", "pool_read_bytes_total", "pool_write_total", "pool_write_bytes_total",
	"cluster_capacity_bytes", "cluster_available_bytes", "cluster_used_bytes", "cluster_objects", "perf_commit_latency_ms", "perf_apply_latency_ms", "osd_crush_weight", "osd_depth", "osd_reweight", "osd_bytes", "osd_bytes_used", "osd_bytes_avail", "osd_utilization",
	"osd_var", "osd_pgs", "osd_total_bytes", "osd_total_used_bytes", "osd_total_avail_bytes", "osd_average_utilization", "osd", "osd_up", "osd_in", "client_io_write_ops", "client_io_read_bytes", "client_io_read_ops", "client_io_write_bytes", "cache_flush_io_bytes", "cache_evict_io_bytes", "cache_promote_io_ops", "degraded_objects", "misplaced_objects", "osds", "osds_up", "osds_in", "pgs_remapped", "total_pgs", "volume_name", "volume_size_bytes", "volume_objects", "volume_object_size_bytes"}

var fakeResp map[string]*MetricFakeResp = map[string]*MetricFakeResp{`{"detail":"detail","format":"json","prefix":"df"}`: {[]byte(`{"stats":{"total_bytes":494462976000,"total_used_bytes":238116864,"total_avail_bytes":494224859136,"total_objects":14},"pools":[{"name":"rbd","id":1,"stats":{"kb_used":1,"bytes_used":859,"percent_used":0.00,"max_avail":469501706240,"objects":14,"quota_objects":0,"quota_bytes":0,"dirty":14,"rd":145,"rd_bytes":918304,"wr":1057,"wr_bytes":16384,"raw_bytes_used":859}}]}`), "", nil},
	`{"format":"json","prefix":"osd df"}`:           {[]byte(`{"nodes":[{"id":0,"device_class":"hdd","name":"osd.0","type":"osd","type_id":0,"crush_weight":0,"depth":2,"pool_weights":{},"reweight":1.000000,"kb":15717356,"kb_used":114624,"kb_avail":15602732,"utilization":0.729283,"var":1.000000,"pgs":102}],"stray":[],"summary":{"total_kb":15717356,"total_kb_used":114624,"total_kb_avail":15602732,"average_utilization":0.729283,"min_var":1.000000,"max_var":1.000000,"dev":0.000000}}`), "", nil},
	`{"format":"json","prefix":"osd dump"}`:         {[]byte(`{"epoch":19,"fsid":"282d4751-4f33-4186-b983-b51cc21a5a8e","created":"2019-05-07 11:49:02.459507","modified":"2019-05-14 16:25:56.992964","flags":"sortbitwise,recovery_deletes,purged_snapdirs","crush_version":3,"full_ratio":0.950000,"backfillfull_ratio":0.900000,"nearfull_ratio":0.850000,"cluster_snapshot":"","pool_max":3,"max_osd":1,"require_min_compat_client":"jewel","min_compat_client":"jewel","require_osd_release":"luminous","pools":[{"pool":1,"pool_name":"sample_pool","flags":1,"flags_names":"hashpspool","type":1,"size":3,"min_size":2,"crush_rule":0,"object_hash":2,"pg_num":1,"pg_placement_num":1,"crash_replay_interval":0,"last_change":"8","last_force_op_resend":"0","last_force_op_resend_preluminous":"0","auid":0,"snap_mode":"selfmanaged","snap_seq":0,"snap_epoch":0,"pool_snaps":[],"removed_snaps":"[]","quota_max_bytes":0,"quota_max_objects":0,"tiers":[],"tier_of":-1,"read_tier":-1,"write_tier":-1,"cache_mode":"none","target_max_bytes":0,"target_max_objects":0,"cache_target_dirty_ratio_micro":400000,"cache_target_dirty_high_ratio_micro":600000,"cache_target_full_ratio_micro":800000,"cache_min_flush_age":0,"cache_min_evict_age":0,"erasure_code_profile":"","hit_set_params":{"type":"none"},"hit_set_period":0,"hit_set_count":0,"use_gmt_hitset":true,"min_read_recency_for_promote":0,"min_write_recency_for_promote":0,"hit_set_grade_decay_rate":0,"hit_set_search_last_n":0,"grade_table":[],"stripe_width":0,"expected_num_objects":0,"fast_read":false,"options":{},"application_metadata":{}},{"pool":2,"pool_name":"rbd","flags":1,"flags_names":"hashpspool","type":1,"size":3,"min_size":2,"crush_rule":0,"object_hash":2,"pg_num":1,"pg_placement_num":1,"crash_replay_interval":0,"last_change":"11","last_force_op_resend":"0","last_force_op_resend_preluminous":"0","auid":0,"snap_mode":"selfmanaged","snap_seq":0,"snap_epoch":0,"pool_snaps":[],"removed_snaps":"[]","quota_max_bytes":0,"quota_max_objects":0,"tiers":[],"tier_of":-1,"read_tier":-1,"write_tier":-1,"cache_mode":"none","target_max_bytes":0,"target_max_objects":0,"cache_target_dirty_ratio_micro":400000,"cache_target_dirty_high_ratio_micro":600000,"cache_target_full_ratio_micro":800000,"cache_min_flush_age":0,"cache_min_evict_age":0,"erasure_code_profile":"","hit_set_params":{"type":"none"},"hit_set_period":0,"hit_set_count":0,"use_gmt_hitset":true,"min_read_recency_for_promote":0,"min_write_recency_for_promote":0,"hit_set_grade_decay_rate":0,"hit_set_search_last_n":0,"grade_table":[],"stripe_width":0,"expected_num_objects":0,"fast_read":false,"options":{},"application_metadata":{}},{"pool":3,"pool_name":"sapm_pools","flags":1,"flags_names":"hashpspool","type":1,"size":1,"min_size":1,"crush_rule":0,"object_hash":2,"pg_num":100,"pg_placement_num":100,"crash_replay_interval":0,"last_change":"19","last_force_op_resend":"0","last_force_op_resend_preluminous":"0","auid":0,"snap_mode":"selfmanaged","snap_seq":3,"snap_epoch":19,"pool_snaps":[],"removed_snaps":"[1~3]","quota_max_bytes":0,"quota_max_objects":0,"tiers":[],"tier_of":-1,"read_tier":-1,"write_tier":-1,"cache_mode":"none","target_max_bytes":0,"target_max_objects":0,"cache_target_dirty_ratio_micro":400000,"cache_target_dirty_high_ratio_micro":600000,"cache_target_full_ratio_micro":800000,"cache_min_flush_age":0,"cache_min_evict_age":0,"erasure_code_profile":"","hit_set_params":{"type":"none"},"hit_set_period":0,"hit_set_count":0,"use_gmt_hitset":true,"min_read_recency_for_promote":0,"min_write_recency_for_promote":0,"hit_set_grade_decay_rate":0,"hit_set_search_last_n":0,"grade_table":[],"stripe_width":0,"expected_num_objects":0,"fast_read":false,"options":{},"application_metadata":{}}],"osds":[{"osd":0,"uuid":"44559c14-fb71-4183-9e44-167e0e9c057a","up":1,"in":1,"weight":1.000000,"primary_affinity":1.000000,"last_clean_begin":0,"last_clean_end":0,"up_from":7,"up_thru":15,"down_at":6,"lost_at":0,"public_addr":"192.168.1.47:6800/1393","cluster_addr":"192.168.1.47:6801/1393","heartbeat_back_addr":"192.168.1.47:6802/1393","heartbeat_front_addr":"192.168.1.47:6803/1393","state":["exists","up"]}],"osd_xinfo":[{"osd":0,"down_stamp":"2019-05-10 18:23:18.070300","laggy_probability":0.000000,"laggy_interval":0,"features":4611087853746454523,"old_weight":0}],"pg_upmap":[],"pg_upmap_items":[],"pg_temp":[],"primary_temp":[],"blacklist":{},"erasure_code_profiles":{"default":{"k":"2","m":"1","plugin":"jerasure","technique":"reed_sol_van"}}}`), "", nil},
	`{"format":"json","prefix":"osd perf"}`:         {[]byte(`{"osd_perf_infos":[{"id":0,"perf_stats":{"commit_latency_ms":0,"apply_latency_ms":0}}]}`), "", nil},
	`{"format":"json","prefix":"status"}`:           {[]byte(`{"fsid":"282d4751-4f33-4186-b983-b51cc21a5a8e","health":{"checks":{"PG_AVAILABILITY":{"severity":"HEALTH_WARN","summary":{"message":"Reduced data availability: 2 pgs inactive"}},"PG_DEGRADED":{"severity":"HEALTH_WARN","summary":{"message":"Degraded data redundancy: 2 pgs undersized"}},"POOL_APP_NOT_ENABLED":{"severity":"HEALTH_WARN","summary":{"message":"application not enabled on 1 pool(s)"}}},"status":"HEALTH_WARN","summary":[{"severity":"HEALTH_WARN","summary":"'ceph health' JSON format has changed in luminous. If you see this your monitoring system is scraping the wrong fields. Disable this with 'mon health preluminous compat warning = false'"}],"overall_status":"HEALTH_WARN"},"election_epoch":5,"quorum":[0],"quorum_names":["openSDS-arpita"],"monmap":{"epoch":1,"fsid":"282d4751-4f33-4186-b983-b51cc21a5a8e","modified":"2019-05-07 11:49:01.502074","created":"2019-05-07 11:49:01.502074","features":{"persistent":["kraken","luminous"],"optional":[]},"mons":[{"rank":0,"name":"openSDS-arpita","addr":"192.168.1.47:6789/0","public_addr":"192.168.1.47:6789/0"}]},"osdmap":{"osdmap":{"epoch":19,"num_osds":1,"num_up_osds":1,"num_in_osds":1,"full":false,"nearfull":false,"num_remapped_pgs":0}},"pgmap":{"pgs_by_state":[{"state_name":"active+clean","count":100},{"state_name":"undersized+peered","count":2}],"num_pgs":102,"num_pools":3,"num_objects":8,"data_bytes":247,"bytes_used":117374976,"bytes_avail":15977197568,"bytes_total":16094572544,"inactive_pgs_ratio":0.019608},"fsmap":{"epoch":1,"by_rank":[]},"mgrmap":{"epoch":9,"active_gid":14097,"active_name":"openSDS-arpita","active_addr":"192.168.1.47:6804/1294","available":true,"standbys":[],"modules":["status"],"available_modules":["balancer","dashboard","influx","localpool","prometheus","restful","selftest","status","zabbix"],"services":{}},"servicemap":{"epoch":1,"modified":"0.000000","services":{}}}`), "", nil},
	`{"format":"json","prefix":"time-sync-status"}`: {[]byte(`{"timechecks":{"epoch":5,"round":0,"round_status":"finished"}}`), "", nil}}

var respMap map[string]*MetricFakeRep = map[string]*MetricFakeRep{"ls": {"opensds-4c5cb264-50d1-4bfd-a663-dface9b669c9", nil}, "info": {`{"name":"opensds-4c5cb264-50d1-4bfd-a663-dface9b669c9","size":1073741824,"objects":1024,"order":20,"object_size":1048576,"block_name_prefix":"rbd_data.1e5246b8b4567","format":2,"features":["layering"],"flags":[],"create_timestamp":"Wed Jun  5 12:45:23 2019"}`, nil}}

type MetricFakeconn struct {
	RespMap map[string]*MetricFakeResp
}

func NewMetricFakeconn(respMap map[string]*MetricFakeResp) Conn {
	return &MetricFakeconn{RespMap: fakeResp}
}

type MetricFakeResp struct {
	buf  []byte
	info string
	err  error
}

func (n *MetricFakeconn) ReadDefaultConfigFile() error {
	return nil
}

func (n *MetricFakeconn) Connect() error {
	return nil
}

func (n *MetricFakeconn) GetFSID() (fsid string, err error) {
	fake_fsid := "b987-654-321"
	return fake_fsid, nil
}
func (n *MetricFakeconn) MonCommand(arg []byte) ([]byte, string, error) {
	temp := string(arg)
	var buffer []byte
	if temp != "" {
		buffer = fakeResp[temp].buf
	}
	return buffer, "", nil
}

func (n *MetricFakeconn) Shutdown() {}

type MetricFakeExecuter struct {
	RespMap map[string]*MetricFakeRep
}

type MetricFakeRep struct {
	out string
	err error
}

func (f *MetricFakeExecuter) Run(name string, args ...string) (string, error) {
	var cmd = name
	if name == "env" {
		cmd = args[1]
	}
	v, ok := f.RespMap[cmd]
	if !ok {
		return "", fmt.Errorf("can't find specified op: %s", args[1])
	}
	return v.out, v.err
}
func NewMetricFakeExecuter(respMap map[string]*MetricFakeRep) exec.Executer {
	return &MetricFakeExecuter{RespMap: respMap}
}

func TestCollectMetrics(t *testing.T) {
	var md = &MetricDriver{}
	md.Setup()
	md.cli = &MetricCli{nil, nil}
	md.cli.conn = NewMetricFakeconn(fakeResp)
	md.cli.RootExecuter = NewMetricFakeExecuter(respMap)
	var tempMetricArray []*model.MetricSpec
	for _, element := range expctdMetricList {
		val, _ := strconv.ParseFloat(expected_data[element].Value, 64)
		expctdmetricValue := &model.Metric{
			Timestamp: 123456,
			Value:     val,
		}
		expctdMetricValues := make([]*model.Metric, 0)
		expctdMetricValues = append(expctdMetricValues, expctdmetricValue)
		metric := &model.MetricSpec{
			InstanceID:   "b987-654-321",
			InstanceName: "opensds-ceph-b987-654-321",
			Job:          "ceph",
			Labels:       expected_data[element].Var_Label,
			Component:    expected_data[element].Component,
			Name:         expected_data[element].Name,
			Unit:         expected_data[element].Unit,
			AggrType:     expected_data[element].AggrType,
			MetricValues: expctdMetricValues,
		}
		tempMetricArray = append(tempMetricArray, metric)
	}

	expectedMetrics := tempMetricArray
	retunMetrics, err := md.CollectMetrics()

	if err != nil {
		t.Error("failed to collect stats:", err)
	}
	// we can't use deep equal on metric spec objects as the timesatmp calulation is time.Now() in driver
	// validate equivalence of go weteach metricspec fields against expected except timestamp
	var b bool = true
	for i, m := range expectedMetrics {
		b = b && reflect.DeepEqual(m.InstanceName, retunMetrics[i].InstanceName)
		b = b && reflect.DeepEqual(m.InstanceID, retunMetrics[i].InstanceID)
		b = b && reflect.DeepEqual(m.Job, retunMetrics[i].Job)
		for k, _ := range m.Labels {
			b = b && reflect.DeepEqual(m.Labels[k], retunMetrics[i].Labels[k])
		}
		b = b && reflect.DeepEqual(m.Component, retunMetrics[i].Component)
		b = b && reflect.DeepEqual(m.Unit, retunMetrics[i].Unit)
		b = b && reflect.DeepEqual(m.AggrType, retunMetrics[i].AggrType)
		for j, v := range m.MetricValues {
			b = b && reflect.DeepEqual(v.Value, retunMetrics[i].MetricValues[j].Value)
		}
	}
	if !b {
		t.Errorf("expected metric spec")
		for _, p := range expectedMetrics {
			t.Logf("%+v\n", p)
			for _, v := range p.MetricValues {
				t.Logf("%+v\n", v)
			}
		}
		t.Errorf("returned metric spec")
		for _, p := range retunMetrics {
			t.Logf("%+v\n", p)
			for _, v := range p.MetricValues {
				t.Logf("%+v\n", v)
			}
		}
	}
}
