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

	if strings.Contains(string(out), sarNotEnabledOut) || strings.Contains(string(out), cmdNotFound) {

		return false
	}

	return true
}

// Function to parse sar and iostat command output
// metricList -> metrics to be collected
// instanceID -> VolumeID/Disk Id
// metricMap	-> metric to command output column mapping
// out 		-> command output
// returnMap	-> metric to value map to be returned
func (c *MetricCli) parseCommandOutput(metricList []string, returnMap map[string]string, instanceID string, metricMap map[string]int, out string) {

	tableRows := strings.Split(string(out), "\n")

	// TODO(Prakash):re-visit the below logic when we add disk metrics support
	// LVM stores the created volume with -- instead of -, so we need to adjust the input instance ID
	instanceID = strings.Replace(instanceID, "-", "--", -1)

	for _, row := range tableRows {

		if strings.Contains(row, instanceID) {
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
			for _, metric := range metricList {
				val := columns[metricMap[metric]]
				returnMap[metric] = val
				returnMap["InstanceName"] = columns[metricMap["InstanceID"]]

			}
		}

	}
}

// CollectMetrics function is to call the cli for metrics collection. This will be invoked  by lvm metric driver
// metricList	-> metrics to be collected
// instanceID	-> for which instance to be collected
// returnMap	-> metrics to value map
func (cli *MetricCli) CollectMetrics(metricList []string, instanceID string) (map[string]string, error) {

	returnMap := make(map[string]string)
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
		metricMap["IOPS"] = 2
		metricMap["ReadThroughput"] = 3
		metricMap["WriteThroughput"] = 4
		metricMap["ResponseTime"] = 7
		metricMap["ServiceTime"] = 8
		metricMap["UtilizationPercentage"] = 9
		//call parser
		cli.parseCommandOutput(metricList, returnMap, instanceID, metricMap, out)
	} else {
		cmd := []string{"env", "LC_ALL=C", "iostat", "-N"}

		out, err := cli.execute(cmd...)

		if strings.Contains(string(out), cmdNotFound) {

			log.Errorf("iostat is not vaoilable: cmd.Run() failed with %s\n", err)
			err = nil
		} else if err != nil {
			log.Errorf("cmd.Run() failed with %s\n", err)
			return nil, err
		}
		metricMap := make(map[string]int)
		// iostat command output mapping
		metricMap["IOPS"] = 1
		metricMap["ReadThroughput"] = 2
		metricMap["WriteThroughput"] = 3
		cli.parseCommandOutput(metricList, returnMap, instanceID, metricMap, out)

	}
	return returnMap, err
}
