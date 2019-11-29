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
This module implements a entry into the OpenSDS service.

*/

package cli

import (
	"os"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var hostCommand = &cobra.Command{
	Use:   "host",
	Short: "manage hosts in the cluster",
	Run:   hostAction,
}

var hostCreateCommand = &cobra.Command{
	Use:   "create <name>",
	Short: "create a host in the cluster",
	Run:   hostCreateAction,
}

var hostShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a host in the cluster",
	Run:   hostShowAction,
}

var hostListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all hosts in the cluster",
	Run:   hostListAction,
}

var hostDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a host in the cluster",
	Run:   hostDeleteAction,
}

var hostUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a host in the cluster",
	Run:   hostUpdateAction,
}

var hostInitiatorCommand = &cobra.Command{
	Use:   "initiator",
	Short: "manage initiators of host in the cluster",
	Run:   hostInitiatorAction,
}

var hostAddInitiatorCommand = &cobra.Command{
	Use:   "add <host id> <port name> <port protocol> ",
	Short: "add/update an initiator into a host in the cluster",
	Run:   hostAddInitiatorAction,
}

var hostRemoveInitiatorCommand = &cobra.Command{
	Use:   "remove <host id> <port name>",
	Short: "remove an initiator from a host in the cluster",
	Run:   hostRemoveInitiatorAction,
}

var (
	accessMode        string
	hostName          string
	osType            string
	ip                string
	availabilityZones []string
)

var (
	hostFormatters = FormatterList{"Initiators": JsonFormatter}
	keysForDetail  = KeyList{"Id", "HostName", "OsType", "IP", "Port", "AccessMode", "Username",
		"AvailabilityZones", "Initiators", "CreatedAt", "UpdatedAt"}
	keysForSummary = KeyList{"Id", "HostName", "OsType", "IP", "AccessMode", "AvailabilityZones"}
)

func init() {

	hostCommand.AddCommand(hostCreateCommand)
	hostCommand.AddCommand(hostDeleteCommand)
	hostCommand.AddCommand(hostShowCommand)
	hostCommand.AddCommand(hostListCommand)
	hostCommand.AddCommand(hostUpdateCommand)

	hostCreateCommand.Flags().StringVarP(&accessMode, "accessMode", "", "agentless", "the access mode of host, including: agentless, agent")
	hostCreateCommand.Flags().StringVarP(&osType, "osType", "", "linux", "the os type of host, includding: linux, windows")
	hostCreateCommand.Flags().StringVarP(&ip, "ip", "", "", "the IP address for access the host")
	hostCreateCommand.Flags().StringSliceVarP(&availabilityZones, "availabilityZones", "", []string{"default"}, "the array of availability zones which host belongs to")

	hostUpdateCommand.Flags().StringVarP(&accessMode, "accessMode", "", "agentless", "the access mode of host, including: agentless, agent")
	hostUpdateCommand.Flags().StringVarP(&hostName, "hostName", "", "", "the host name of host")
	hostUpdateCommand.Flags().StringVarP(&osType, "osType", "", "linux", "the os type of host, includding: linux, windows")
	hostUpdateCommand.Flags().StringVarP(&ip, "ip", "", "", "the IP address for access the host")
	hostUpdateCommand.Flags().StringSliceVarP(&availabilityZones, "availabilityZones", "", []string{"default"}, "the array of availability zones which host belongs to")

	hostInitiatorCommand.AddCommand(hostAddInitiatorCommand)
	hostInitiatorCommand.AddCommand(hostRemoveInitiatorCommand)
	hostCommand.AddCommand(hostInitiatorCommand)
}

func hostAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func hostInitiatorAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func hostCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	host := &model.HostSpec{
		AccessMode:        accessMode,
		HostName:          args[0],
		OsType:            osType,
		IP:                ip,
		AvailabilityZones: availabilityZones,
	}

	resp, err := client.CreateHost(host)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintDict(resp, keysForDetail, hostFormatters)
}

func hostShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetHost(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintDict(resp, keysForDetail, hostFormatters)
}

func hostListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	var opts = map[string]string{"hostName": hostName}
	resp, err := client.ListHosts(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintList(resp, keysForSummary, hostFormatters)
}

func hostDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	err := client.DeleteHost(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func hostUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	host := &model.HostSpec{
		AccessMode:        accessMode,
		HostName:          hostName,
		OsType:            osType,
		IP:                ip,
		AvailabilityZones: availabilityZones,
	}

	resp, err := client.UpdateHost(args[0], host)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintDict(resp, keysForDetail, hostFormatters)
}

func hostAddInitiatorAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 3)
	tmpHost, err := client.GetHost(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	var initiators []*model.Initiator
	for _, e := range tmpHost.Initiators {
		if args[1] == e.PortName {
			continue
		}
		initiators = append(initiators, e)
	}
	initiators = append(initiators, &model.Initiator{
		PortName: args[1],
		Protocol: args[2],
	})

	host := &model.HostSpec{
		Initiators: initiators,
	}

	resp, err := client.UpdateHost(args[0], host)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintDict(resp, keysForDetail, hostFormatters)
}

func hostRemoveInitiatorAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	tmpHost, err := client.GetHost(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	var initiators []*model.Initiator
	for _, e := range tmpHost.Initiators {
		if args[1] == e.PortName {
			continue
		}
		initiators = append(initiators, e)
	}

	host := &model.HostSpec{
		Initiators: initiators,
	}

	resp, err := client.UpdateHost(args[0], host)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	PrintDict(resp, keysForDetail, hostFormatters)
}
