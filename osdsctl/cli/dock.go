// Copyright 2017 The OpenSDS Authors.
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
This module implements a entry into the OpenSDS service.

*/

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var dockCommand = &cobra.Command{
	Use:   "dock",
	Short: "manage OpenSDS dock resources",
	Run:   dockAction,
}

var dockShowCommand = &cobra.Command{
	Use:   "show <dock id>",
	Short: "show information of specified dock",
	Run:   dockShowAction,
}

var dockListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all dock resources",
	Run:   dockListAction,
}

var dockImportCommand = &cobra.Command{
	Use:   "import <storage type> <dock config file>",
	Short: "import dock to opensds",
	Run:   dockImportAction,
}

var (
	dockLimit       string
	dockOffset      string
	dockSortDir     string
	dockSortKey     string
	dockId          string
	dockName        string
	dockDescription string
	dockStatus      string
	dockStorageType string
	dockEndpoint    string
	dockDriverName  string
	dockPoolConfigFilePath string
)

func init() {
	dockListCommand.Flags().StringVarP(&dockLimit, "limit", "", "50", "the number of ertries displayed per page")
	dockListCommand.Flags().StringVarP(&dockOffset, "offset", "", "0", "all requested data offsets")
	dockListCommand.Flags().StringVarP(&dockSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	dockListCommand.Flags().StringVarP(&dockSortKey, "sortKey", "", "id", "the sort key of all requested data. supports id(default), name, status, endpoint, drivername, description")
	dockListCommand.Flags().StringVarP(&dockId, "id", "", "", "list docks by id")
	dockListCommand.Flags().StringVarP(&dockName, "name", "", "", "list docks by name")
	dockListCommand.Flags().StringVarP(&dockDescription, "description", "", "", "list docks by description")
	dockListCommand.Flags().StringVarP(&dockStatus, "status", "", "", "list docks by status")
	dockListCommand.Flags().StringVarP(&dockStorageType, "storageType", "", "", "list docks by storage type")
	dockListCommand.Flags().StringVarP(&dockEndpoint, "endpoint", "", "", "list docks by endpoint")
	dockListCommand.Flags().StringVarP(&dockDriverName, "driverName", "", "", "list docks by driver name")
	dockImportCommand.Flags().StringVarP(&dockPoolConfigFilePath, "poolConfig", "p", "", "indicate the pool config file")

	dockCommand.AddCommand(dockShowCommand)
	dockCommand.AddCommand(dockListCommand)
	dockCommand.AddCommand(dockImportCommand)
}

func dockAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var dockFormatters = FormatterList{"Metadata": JsonFormatter}

func dockShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetDock(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Endpoint", "DriverName",
		"Metadata"}
	PrintDict(resp, keys, dockFormatters)
}

func dockListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	var opts = map[string]string{"limit": dockLimit, "offset": dockOffset, "sortDir": dockSortDir,
		"sortKey": dockSortKey, "Id": dockId,
		"Name": dockName, "Description": dockDescription, "DriverName": dockDriverName,
		"Endpoint": dockEndpoint, "Status": dockStatus, "StorageType": dockStorageType}

	resp, err := client.ListDocks(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Endpoint", "DriverName"}
	PrintList(resp, keys, dockFormatters)
}

func dockImportAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	storageType := args[0]
	dockConfigFilePath := args[1]
	addBackendType(storageType)
	addBackendInfo(dockConfigFilePath)
	if dockPoolConfigFilePath != "" {
		importDriverConfig(storageType, dockPoolConfigFilePath)
	}
}

func addBackendType(backendType string) {
	config, err := ioutil.ReadFile(constants.OpensdsConfigPath)
	if err != nil {
		Errorln(constants.OpensdsConfigPath, " is not exist!")
		os.Exit(1)
	}
	lines := strings.Split(string(config), "\n")
	for i, line := range lines {
		if strings.Contains(line, "enabled_backends") {
			lines[i] = lines[i] + "," + backendType
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(constants.OpensdsConfigPath, []byte(output), 0644)
}

func addBackendInfo(backendConfig string) {
	dockConfig, err := ioutil.ReadFile(backendConfig)
	if err != nil {
		Errorln(backendConfig, " is not exist!")
		os.Exit(1)
	}
	importDockCmd := fmt.Sprintf("cat>>%s<<EOF\n%sEOF", constants.OpensdsConfigPath, string(dockConfig))
	_,err = exec.Command("bash", "-c",importDockCmd).CombinedOutput()
	if err != nil {
		fmt.Println("err is ", err)
	}
}

func importDriverConfig(storageType string, driverConfig string) {
	poolConfig, err := ioutil.ReadFile(driverConfig)
	driverPath := fmt.Sprintf("/etc/opensds/driver/%s.yaml", storageType)
	err = ioutil.WriteFile(driverPath, poolConfig, 0644)
	if err != nil {
		Errorln("Write file failed")
		os.Exit(1)
	}
}