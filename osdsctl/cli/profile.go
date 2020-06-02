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
	"encoding/json"
	"os"

	"github.com/sodafoundation/api/pkg/model"
	"github.com/spf13/cobra"
)

var profileCommand = &cobra.Command{
	Use:   "profile",
	Short: "manage OpenSDS profile resources",
	Run:   profileAction,
}

var profileCreateCommand = &cobra.Command{
	Use: "create <profile-info>",
	Example: "osdsctl profile create '{\"name\": \"default_block\", \"description\": \"default policy\", \"storageType\": \"block\"}'" + "\n" +
		"osdsctl profile create '{\"name\": \"default_file\", \"description\": \"default policy\", \"storageType\": \"file\", \"provisioningProperties\":{\"ioConnectivity\": {\"accessProtocol\": \"NFS\"},\"DataStorage\":{\"StorageAccessCapability\":[\"Read\",\"Write\",\"Execute\"]}}}'" +
		"\n" +
		"\n" +
		"The example of more supported \"profile-info\" parameters:\n" +
		"\n" +
		"'{" +
		"\"name\": \"File_Profile\"," +
		"\"storageType\": \"block\"," +
		"\"description\": \"string\"," +
		"\"provisioningProperties\": {" +
		"\"dataStorage\": {" +
		"\"recoveryTimeObjective\": 10," +
		"\"provisioningPolicy\": \"Thick\"," +
		"\"compression\": false," +
		"\"deduplication\": false," +
		"\"characterCodeSet\": \"ASCII\"," +
		"\"maxFileNameLengthBytes\": 255," +
		"\"storageAccessCapability\": [\"Read\"] " +
		"}," +
		"\"ioConnectivity\": {" +
		"\"accessProtocol\": \"iscsi\"," +
		"\"maxIOPS\": 150," +
		"\"minIOPS\": 50," +
		"\"maxBWS\": 5," +
		"\"minBWS\": 1," +
		"\"latency\": 1" +
		"}" +
		"}," +
		"\"replicationProperties\": {" +
		"\"dataProtection\": {" +
		"\"isIsolated\": true," +
		"\"minLifetime\": \"P3Y6M4DT12H30M55\"," +
		"\"RecoveryGeographicObjective\": \"datacenter\"," +
		"\"RecoveryPointObjectiveTime\": \"P3Y6M4DT12H30M5S\"," +
		"\"RecoveryTimeObjective\": \"offline\"," +
		"\"ReplicaType\": \"snapshot\"" +
		"}," +
		"\"replicaInfos\": {" +
		"\"replicaUpdateMode\": \"Active\"," +
		"\"replcationBandwidth\": 5," +
		"\"replicationPeriod\": \"P3Y6M4DT12H30M5S\"," +
		"\"consistencyEnalbed\": true " +
		"}" +
		"}," +
		"\"snapshotProperties\": {" +
		"\"schedule\": {" +
		"\"datetime\": \"2019-09-07T07:02:35.389\"," +
		"\"occurrence\": \"Daily\"" +
		"}," +
		"\"retention\": {" +
		"\"duration\": 15," +
		"\"number\": 10" +
		"}," +
		"\"topology\": {" +
		"\"bucket\": \"string\"" +
		"	}" +
		"}," +
		"\"dataProtectionProperties\": {" +
		"\"dataProtection\": {" +
		"\"isIsolated\": true," +
		"\"minLifetime\": \"P3Y6M4DT12H30M5S\"," +
		"\"RecoveryGeographicObjective\": \"datacenter\"," +
		"\"RecoveryPointObjectiveTime\": \"P3Y6M4DT12H30M5S\"," +
		"\"RecoveryTimeObjective\": \"offline\"," +
		"\"ReplicaType\": \"snapshot\"" +
		"}," +
		"\"consistencyEnalbed\": true " +
		"}," +
		"\"customProperties\": {" +
		"\"key1\": \"value1\"," +
		"\"key2\": false, " +
		"\"key3\": { " +
		"\"key31\": \"value31\"" +
		"}" +
		"}" +
		"}'",

	Short: "create a new profile resource",
	Run:   profileCreateAction,
}

var profileShowCommand = &cobra.Command{
	Use:   "show <profile id>",
	Short: "show information of specified profile",
	Run:   profileShowAction,
}

var profileListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all profile resources",
	Example: "osdsctl profile list --description \"test\"\n" +
		"osdsctl profile list --id def32d39-78e2-47f3-9e2f-8c43a8b9ee3a\n" +
		"osdsctl profile list --limit 2\n" +
		"osdsctl profile list --name test\n" +
		"osdsctl profile list --offset 2\n" +
		"osdsctl profile list --sortDir desc\n" +
		"osdsctl profile list --sortKey id\n" +
		"osdsctl profile list --storageType block\n",
	Run: profileListAction,
}

var profileDeleteCommand = &cobra.Command{
	Use:   "delete <profile id>",
	Short: "delete a specified profile resource",
	Run:   profileDeleteAction,
}

var (
	profLimit       string
	profOffset      string
	profSortDir     string
	profSortKey     string
	profId          string
	profName        string
	profDescription string
	profStorageType string
	profInfo        string
)

func init() {
	profileListCommand.Flags().StringVarP(&profLimit, "limit", "", "50", "the number of ertries displayed per page")
	profileListCommand.Flags().StringVarP(&profOffset, "offset", "", "0", "all requested data offsets")
	profileListCommand.Flags().StringVarP(&profSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	profileListCommand.Flags().StringVarP(&profSortKey, "sortKey", "", "id", "the sort key of all requested data. supports id(default), name, description")
	profileListCommand.Flags().StringVarP(&profId, "id", "", "", "list profile by id")
	profileListCommand.Flags().StringVarP(&profName, "name", "", "", "list profile by name")
	profileListCommand.Flags().StringVarP(&profDescription, "description", "", "", "list profile by description")
	profileListCommand.Flags().StringVarP(&profStorageType, "storageType", "", "", "list profile by storage type")

	profileCommand.AddCommand(profileCreateCommand)
	profileCreateCommand.Flags().Lookup("profile-info")
	profileCommand.AddCommand(profileShowCommand)
	profileCommand.AddCommand(profileListCommand)
	profileCommand.AddCommand(profileDeleteCommand)
}

func profileAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var profileFormatters = FormatterList{"ProvisioningProperties": JsonFormatter, "ReplicationProperties": JsonFormatter,
	"SnapshotProperties": JsonFormatter, "DataProtectionProperties": JsonFormatter, "CustomProperties": JsonFormatter}

func profileCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	prf := &model.ProfileSpec{}
	if err := json.Unmarshal([]byte(args[0]), prf); err != nil {
		Errorln(err)
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.CreateProfile(prf)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "Name", "Description", "StorageType", "ProvisioningProperties",
		"ReplicationProperties", "SnapshotProperties", "DataProtectionProperties", "CustomProperties"}
	PrintDict(resp, keys, profileFormatters)
}

func profileShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetProfile(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "StorageType", "ProvisioningProperties",
		"ReplicationProperties", "SnapshotProperties", "DataProtectionProperties", "CustomProperties"}
	PrintDict(resp, keys, profileFormatters)
}

func profileListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	var opts = map[string]string{"limit": profLimit, "offset": profOffset, "sortDir": profSortDir,
		"sortKey": profSortKey, "Id": profId,
		"Name": profName, "Description": profDescription, "StorageType": profStorageType}

	resp, err := client.ListProfiles(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "StorageType"}
	PrintList(resp, keys, FormatterList{})
}

func profileDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	err := client.DeleteProfile(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}
