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

package lvm

import (
	"regexp"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils/exec"
)

const (
	sarNotEnabledOut = "Please check if data collecting is enabled"
	cmdNotFound      = "No such file or directory"
	sarNotFound      = "Command 'sar' not found"
	iostatNotFound   = "Command 'iostat' not found"
)

type MetricCli struct {
	// Command executer
	BaseExecuter exec.Executer
	// Command Root executer
	RootExecuter exec.Executer
}

func NewMetricCli() (*MetricCli, error) {
	return &MetricCli{
		BaseExecuter: exec.NewBaseExecuter(),
		RootExecuter: exec.NewRootExecuter(),
	}, nil
}

func (c *MetricCli) execute(cmd ...string) (string, error) {
	return c.RootExecuter.Run(cmd[0], cmd[1:]...)
}

func isSarEnabled(out string) bool {

	if strings.Contains(string(out), sarNotEnabledOut) || strings.Contains(string(out), cmdNotFound) || strings.Contains(string(out), sarNotFound) {

		return false
	}

	return true
}

// Function to parse sar and iostat command output
// metricList -> metrics to be collected
// metricMap	-> metric to command output column mapping
// out 		-> command output
// returnMap	-> metric to value map to be returned
func (c *MetricCli) parseCommandOutput(metricList []string, returnMap map[string]map[string]string, labelMap map[string]map[string]string, metricMap map[string]int, out string) {

	tableRows := strings.Split(string(out), "\n")
	for _, row := range tableRows[3:] {
		if row != "" {
			tokens := regexp.MustCompile(" ")
			cols := tokens.Split(row, -1)
			// remove all empty space
			var columns = make([]string, 0, 0)
			for _, v := range cols {
				if v != "" {
					columns = append(columns, v)
				}
			}
			// map the values
			aVolMap := make(map[string]string)
			aLabMap := make(map[string]string)
			for _, metric := range metricList {
				val := columns[metricMap[metric]]
				aVolMap[metric] = val
				aVolMap["InstanceName"] = columns[metricMap["InstanceID"]]
				aLabMap["device"] = columns[metricMap["InstanceID"]]
			}
			returnMap[columns[1]] = aVolMap
			labelMap[columns[1]] = aLabMap
		}
	}

}

// CollectMetrics function is to call the cli for metrics collection. This will be invoked  by lvm metric driver
// metricList	-> metrics to be collected
// returnMap	-> metrics to value map
func (cli *MetricCli) CollectMetrics(metricList []string) ( /*returnMAp*/ map[string]map[string]string /*labelMap*/, map[string]map[string]string, error) {

	returnMap := make(map[string]map[string]string)
	labelMap := make(map[string]map[string]string)
	var err error

	cmd := []string{"env", "LC_ALL=C", "sar", "-dp", "1", "1"}
	out, err := cli.execute(cmd...)
	if err != nil {
		log.Errorf("cmd.Run() failed with %s\n", err)
		err = nil

	}
	//check whether sar collection is enabled ?
	//If not use iostat command
	if isSarEnabled(out) {
		// sar command output mapping
		metricMap := make(map[string]int)
		metricMap["InstanceID"] = 1
		metricMap["iops"] = 2
		metricMap["read_throughput"] = 3
		metricMap["write_throughput"] = 4
		metricMap["response_time"] = 7
		metricMap["service_time"] = 8
		metricMap["utilization"] = 9
		//call parser
		cli.parseCommandOutput(metricList, returnMap, labelMap, metricMap, out)
	} else {
		cmd := []string{"env", "LC_ALL=C", "iostat", "-N"}
		out, err := cli.execute(cmd...)
		if strings.Contains(string(out), cmdNotFound) || strings.Contains(string(out), iostatNotFound) {
			log.Errorf("iostat is not available: cmd.Run() failed with %s\n", err)
			return nil, nil, err
		} else if err != nil {
			log.Errorf("cmd.Run() failed with %s\n", err)
			return nil, nil, err
		}
		metricMap := make(map[string]int)
		// iostat command output mapping
		metricMap["iops"] = 1
		metricMap["read_throughput"] = 2
		metricMap["write_throughput"] = 3
		cli.parseCommandOutput(metricList, returnMap, labelMap, metricMap, out)

	}
	return returnMap, labelMap, err
}

// Discover LVM volumes
func (c *MetricCli) DiscoverVolumes() ([]string, []string, error) {
	cmd := []string{"env", "LC_ALL=C", "lvs"}
	out, err := c.execute(cmd...)
	tableRows := strings.Split(string(out), "\n")
	var volumes []string
	var vgs []string
	for _, row := range tableRows[1:] {
		if row != "" {
			tokens := regexp.MustCompile(" ")
			cols := tokens.Split(row, -1)
			var columns = make([]string, 0, 0)
			for _, v := range cols {
				if v != "" {
					columns = append(columns, v)
				}
			}
			volumes = append(volumes, columns[0])
			vgs = append(vgs, columns[1])

		}
	}
	return volumes, vgs, err
}

// Discover LVM Disks
func (c *MetricCli) DiscoverDisks() ([]string, error) {
	cmd := []string{"env", "LC_ALL=C", "pvs"}
	out, err := c.execute(cmd...)
	tableRows := strings.Split(string(out), "\n")
	var volumes []string
	for _, row := range tableRows[1:] {
		if row != "" {
			tokens := regexp.MustCompile(" ")
			cols := tokens.Split(row, -1)
			var columns = make([]string, 0, 0)
			for _, v := range cols {
				if v != "" {
					columns = append(columns, v)
				}
			}
			volumes = append(volumes, columns[0])

		}
	}
	return volumes, err
}
