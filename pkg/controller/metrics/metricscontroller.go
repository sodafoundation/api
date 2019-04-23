// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"encoding/json"
	"fmt"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Controller is an interface for exposing some operations of metric controllers.
type Controller interface {
	GetLatestMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error)
	GetInstantMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error)
	GetRangeMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error)
	SetDock(dockInfo *model.DockSpec)
}

// NewController method creates a controller structure and expose its pointer.
func NewController() Controller {
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
	Name     string `json:"__name__"`
	Device   string `json:"device"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
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

// latest+instant metrics structs begin
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

// latest+instant metrics structs end

func (c *controller) GetLatestMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error) {

	var metrics []model.MetricSpec
	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query?query=" + opt.MetricName)
	if err != nil {
		fmt.Printf("The HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv InstantMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		fmt.Println(err0)

		metrics := make([]model.MetricSpec, len(fv.Data.Result))

		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for i, res := range fv.Data.Result {
			metrics[i].InstanceID = res.Metric.Instance + res.Metric.Device
			metrics[i].Name = res.Metric.Name
			metrics[i].MetricValues = make([]model.Metric, len(res.Value))
			for j, v := range res.Value {
				switch v.(type) {
				case string:
					metrics[i].MetricValues[j].Value, _ = strconv.ParseFloat(v.(string), 64)
				case float64:
					secs := int64(v.(float64))
					metrics[i].MetricValues[j].Timestamp = secs
				/*case []interface{}:
				for i, u := range vv {
					fmt.Println(i, u)
				}*/
				default:
					fmt.Println(v, "is of a type I don't know how to handle")
				}
				/*metrics[i].CollectedAt = t
				metrics[i].Value = v*/
			}
		}

		fmt.Println("metrics struct is ", metrics)
		bArr, _ := json.Marshal(metrics)
		fmt.Println("metrics response json is ", string(bArr))

		if err != nil {
			fmt.Println(err)
		}

	}
	return &metrics, err
}

func (c *controller) GetInstantMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error) {

	var metrics []model.MetricSpec
	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query?query=" + opt.MetricName + "&time=" + opt.StartTime)
	if err != nil {
		fmt.Printf("The HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv InstantMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		fmt.Println(err0)

		metrics := make([]model.MetricSpec, len(fv.Data.Result))

		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for i, res := range fv.Data.Result {
			metrics[i].InstanceID = res.Metric.Instance + res.Metric.Device
			metrics[i].Name = res.Metric.Name
			metrics[i].MetricValues = make([]model.Metric, len(res.Value))
			for j, v := range res.Value {
				switch v.(type) {
				case string:
					metrics[i].MetricValues[j].Value, _ = strconv.ParseFloat(v.(string), 64)
				case float64:
					secs := int64(v.(float64))
					metrics[i].MetricValues[j].Timestamp = secs
				/*case []interface{}:
				for i, u := range vv {
					fmt.Println(i, u)
				}*/
				default:
					fmt.Println(v, "is of a type I don't know how to handle")
				}
				/*metrics[i].CollectedAt = t
				metrics[i].Value = v*/
			}
		}

		fmt.Println("metrics struct is ", metrics)
		bArr, _ := json.Marshal(metrics)
		fmt.Println("metrics response json is ", string(bArr))

		if err != nil {
			fmt.Println(err)
		}

	}
	return &metrics, err
}

func (c *controller) GetRangeMetrics(opt *pb.GetMetricsOpts) (*[]model.MetricSpec, error) {

	var metrics []model.MetricSpec
	// make a call to Prometheus, convert the response to our format, return
	response, err := http.Get("http://localhost:9090/api/v1/query_range?query=" + opt.MetricName + "&start=" + opt.StartTime + "&end=" + opt.EndTime + "&step=30")
	if err != nil {
		fmt.Printf("The HTTP query request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		// unmarshal the JSON response into a struct (generated using the JSON, using this https://mholt.github.io/json-to-go/
		var fv RangeMetricReponseFromPrometheus
		err0 := json.Unmarshal(data, &fv)
		fmt.Println(err0)

		metrics := make([]model.MetricSpec, len(fv.Data.Result))

		// now convert to our repsonse struct, so we can marshal it and send out the JSON
		for i, res := range fv.Data.Result {
			metrics[i].InstanceID = res.Metric.Instance + res.Metric.Device
			metrics[i].Name = res.Metric.Name
			metrics[i].MetricValues = make([]model.Metric, len(res.Values))
			for j := 0; j < len(res.Values); j++ {
				for _, v := range res.Values[j] {
					switch v.(type) {
					case string:
						metrics[i].MetricValues[j].Value, _ = strconv.ParseFloat(v.(string), 64)
					case float64:
						secs := int64(v.(float64))
						metrics[i].MetricValues[j].Timestamp = secs
					/*case []interface{}:
					for i, u := range vv {
						fmt.Println(i, u)
					}*/
					default:
						fmt.Println(v, "is of a type I don't know how to handle")
					}
					/*metrics[i].CollectedAt = t
					metrics[i].Value = v*/
				}
			}
		}

		fmt.Println("metrics struct is ", metrics)
		bArr, _ := json.Marshal(metrics)
		fmt.Println("metrics response json is ", string(bArr))

		if err != nil {
			fmt.Println(err)
		}
		return &metrics, err

	}
	return &metrics, nil
}

func (c *controller) SetDock(dockInfo *model.DockSpec) {
	c.DockInfo = dockInfo
}
