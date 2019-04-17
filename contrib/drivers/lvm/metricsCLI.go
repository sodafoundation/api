package lvm

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	myExec "os/exec"
)

// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

type MetricsCli struct {

}


func (c *MetricsCli) CollectMetrics(metricList []string,instanceID string) (map[string] float64,error) {

	cmd := myExec.Command("iostat", "-N")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Println("Finished: output is %s",string(out))
	metricMap:=make(map[string] int)
	metricMap["kB_read/s"]=2
	metricMap["kB_wrtn/s"]=3


	someSlice:=strings.Split(string(out),"\n")

	instanceID = strings.Replace(instanceID, "-", "--",-1)

	returnMap := make(map[string] float64)
	for _, element := range someSlice {
		// index is the index where we are
		// element is the element from someSlice for where we are

		if strings.Contains(element, instanceID) {
			//strings.Replace(element," ","",-1)
			//strings.FieldsFunc()
			tokens:=regexp.MustCompile(" .")
			stringSlice:=tokens.Split(element,-1)
			// remove all empty space
			var temparray = make([]string,0,0)
			for _,v:= range stringSlice{
				if v != ""{
					temparray=append(temparray, v)
				}
			}
			for _,metric := range metricList{
				val,_:=strconv.ParseFloat(temparray[metricMap[metric]],64)
				returnMap[metric]=val

			}

			fmt.Println(returnMap)

			//	kB_read/s    kB_wrtn/s
		}
	}

	//_,err:=c.execute("iostat")
	return returnMap,err
}



















