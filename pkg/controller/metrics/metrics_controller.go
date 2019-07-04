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

/*
This module implements a entry into the OpenSDS metrics controller service.

*/

package metrics

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/controller/metrics/adapters"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// Controller is an interface for exposing some operations of metric controllers.
type Controller interface {
	GetLatestMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error)
	GetInstantMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error)
	GetRangeMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error)
	CollectMetrics(opt *pb.CollectMetricsOpts) ([]*model.MetricSpec, error)
	SetDock(dockInfo *model.DockSpec)
	GetUrls() (*map[string]model.UrlDesc, error)
}

// NewController method creates a controller structure and expose its pointer.
func NewController() Controller {

	// start the metrics dispatcher
	adapters.StartDispatcher()

	return &controller{
		Client: client.NewClient(),
	}
}

type controller struct {
	client.Client
	DockInfo *model.DockSpec
}

// latest+instant metrics structs begin
type InstantMetricReponseFromPrometheus struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
type Metric struct {
	Name       string `json:"__name__"`
	Device     string `json:"device"`
	InstanceID string `json:"instanceID"`
	Job        string `json:"job"`
}
type Result struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}
type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

// latest+instant metrics structs end

// latest+range metrics structs begin
type RangeMetricReponseFromPrometheus struct {
	Status string    `json:"status"`
	Data   RangeData `json:"data"`
}
type RangeMetric struct {
	Name     string `json:"__name__"`
	Device   string `json:"device"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type RangeResult struct {
	Metric RangeMetric     `json:"metric"`
	Values [][]interface{} `json:"values"`
}
type RangeData struct {
	ResultType string        `json:"resultType"`
	Result     []RangeResult `json:"result"`
}

// latest+range metrics structs end

func (c *controller) GetLatestMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error) {

	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query?query=" + opt.MetricName)
	if err != nil {
		log.Errorf("the HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Infof("response data is %s", string(data))
		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv InstantMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		if err0 != nil {
			log.Errorf("unmarshell operation failed %s\n", err0)
		}
		var metrics []*model.MetricSpec
		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for _, res := range fv.Data.Result {

			metricValues := make([]*model.Metric, 0)
			metricValue := &model.Metric{}
			for _, v := range res.Value {

				switch v.(type) {
				case string:
					metricValue.Value, err = strconv.ParseFloat(v.(string), 64)
				case float64:
					secs := int64(v.(float64))
					metricValue.Timestamp = secs
				default:
					log.Info(v, "is of a type I don't know how to handle")
				}

			}
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{}
			metric.InstanceID = res.Metric.InstanceID
			metric.Name = res.Metric.Name
			metric.InstanceName = res.Metric.Device
			metric.MetricValues = metricValues
			metrics = append(metrics, metric)
		}

		bArr, _ := json.Marshal(metrics)
		log.Infof("metrics response json is %s", string(bArr))

		if err != nil {
			log.Error(err)
		}
		return metrics, err

	}
	return nil, err
}

func (c *controller) GetInstantMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error) {

	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query?query=" + opt.MetricName + "&time=" + opt.StartTime)
	if err != nil {
		log.Errorf("the HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Infof("response data is %s", string(data))
		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv InstantMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		if err0 != nil {
			log.Errorf("unmarshell operation failed %s\n", err0)
		}
		var metrics []*model.MetricSpec
		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for _, res := range fv.Data.Result {

			metricValues := make([]*model.Metric, 0)
			metricValue := &model.Metric{}
			for _, v := range res.Value {

				switch v.(type) {
				case string:
					metricValue.Value, err = strconv.ParseFloat(v.(string), 64)
				case float64:
					secs := int64(v.(float64))
					metricValue.Timestamp = secs
				default:
					log.Info(v, "is of a type I don't know how to handle")
				}

			}
			metricValues = append(metricValues, metricValue)
			metric := &model.MetricSpec{}
			metric.InstanceID = res.Metric.InstanceID
			metric.Name = res.Metric.Name
			metric.InstanceName = res.Metric.Device
			metric.MetricValues = metricValues
			metrics = append(metrics, metric)
		}

		bArr, _ := json.Marshal(metrics)
		log.Infof("metrics response json is %s", string(bArr))

		if err != nil {
			log.Error(err)
		}
		return metrics, err

	}
	return nil, err
}

func (c *controller) GetRangeMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error) {

	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query_range?query=" + opt.MetricName + "&start=" + opt.StartTime + "&end=" + opt.EndTime + "&step=30")
	if err != nil {
		log.Errorf("the HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Info(string(data))

		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv RangeMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		if err0 != nil {
			log.Errorf("unmarshell operation failed %s\n", err0)
		}
		var metrics []*model.MetricSpec
		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for _, res := range fv.Data.Result {

			metricValues := make([]*model.Metric, 0)
			for j := 0; j < len(res.Values); j++ {
				metricValue := &model.Metric{}
				for _, v := range res.Values[j] {
					switch v.(type) {
					case string:
						metricValue.Value, _ = strconv.ParseFloat(v.(string), 64)
					case float64:
						secs := int64(v.(float64))
						metricValue.Timestamp = secs
					default:
						log.Infof("%s is of a type I don't know how to handle", v)
					}

				}
				metricValues = append(metricValues, metricValue)

			}
			metric := &model.MetricSpec{}
			metric.InstanceID = res.Metric.Instance
			metric.Name = res.Metric.Name
			metric.InstanceName = res.Metric.Device
			metric.MetricValues = metricValues
			metrics = append(metrics, metric)
		}

		bArr, _ := json.Marshal(metrics)
		log.Infof("metrics response json is %s", string(bArr))

		if err != nil {
			log.Error(err)
		}
		return metrics, err

	}
	return nil, err
}

func (c *controller) SetDock(dockInfo *model.DockSpec) {
	c.DockInfo = dockInfo
}

func (c *controller) CollectMetrics(opt *pb.CollectMetricsOpts) ([]*model.MetricSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Errorf("error when connecting dock client: %s", err.Error())
		return nil, err
	}

	response, err := c.Client.CollectMetrics(context.Background(), opt)
	if err != nil {
		log.Errorf("collect metrics failed in metrics controller: %s", err.Error())
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to collect metrics in metrics controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	res := make([]*model.MetricSpec, 0)

	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), &res); err != nil {
		log.Errorf("collect metrics failed in metrics controller:%s", err.Error())
		return nil, err
	}

	// send the metrics to the registered adapters
	log.Infoln("send the metrics to the registered adapters")
	for _, metricSpecPtr := range res {
		logMetricSpec(metricSpecPtr)
		adapters.SendMetricToRegisteredSenders(metricSpecPtr)
	}

	return res, nil

}
func CheckServiceStatus(sName string) error {
	cmd := exec.Command("systemctl", "status", sName, "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("status of service %s is inactive err: %s\n", sName, err)
		return err
	}
	return nil
}

func (c *controller) GetUrls() (*map[string]model.UrlDesc, error) {
	res := map[string]model.UrlDesc{}
	flagGrafanaUrl := flag.Lookup("grafana-url")
	flagAlertMgrUrl := flag.Lookup("alertmgr-url")
	if CheckServiceStatus("grafana-server.service") == nil {
		res["Grafana"] = model.UrlDesc{Url: flagGrafanaUrl.Value.String(), Desc: "Open Grafana tool to visualize collected metrics"}
	}
	if CheckServiceStatus("alertmanager.service") == nil {
		res["Alert Manager"] = model.UrlDesc{Url: flagAlertMgrUrl.Value.String(), Desc: "Open Alert Manager to view and handle alerts sent by prometheus"}
	}

	return &res, nil

}

func logMetricSpec(spec *model.MetricSpec) {
	log.Infoln("Name - " + spec.Name)
	log.Infoln("Labels - ")
	for k, v := range spec.Labels {
		log.Infof("\t %s=%s \n", k, v)
	}
	log.Infof("\nAggrType - %s", spec.AggrType)
	log.Infof("\nComponent - %s", spec.Component)
	log.Infof("\nInstanceName - %s", spec.InstanceName)
	log.Infof("\nJob - %s", spec.Job)
	log.Infof("\nUnit - %s", spec.Unit)
	log.Infof("\nInstanceID - %s", spec.InstanceID)
	log.Infof("\nMetricValues - ")
	for _, ms := range spec.MetricValues {
		log.Infof("\n\t value = %v", ms.Value)
		log.Infof("\n\t ts = %v", ms.Timestamp)
	}
}
