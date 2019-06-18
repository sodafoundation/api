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

package main

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/opensds/opensds/pkg/utils/pwd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	confFile = "./pwdEncrypter.yaml"
)

type tool struct {
	PwdEncrypter string `yaml:"PwdEncrypter,omitempty"`
}

var encrypterCommand = &cobra.Command{
	Use:   "pwdEncrypter <password>",
	Short: "password encryption tool",
	Run:   encrypter,
}

func encrypter(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Usage()
		os.Exit(0)
	}

	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	// Initialize configuration file
	pwdEncrypter, err := loadConf(confFile)
	if err != nil {
		fmt.Println("Encrypt password error:", err)
		os.Exit(1)
	}

	// Encrypt the password
	encrypterTool := pwd.NewPwdEncrypter(pwdEncrypter.PwdEncrypter)
	plaintext, err := encrypterTool.Encrypter(args[0])
	if err != nil {
		fmt.Println("Encrypt password error:", err)
		os.Exit(1)
	}

	fmt.Println(plaintext)
}

func loadConf(f string) (*tool, error) {
	conf := &tool{}
	confYaml, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("Read config yaml file (%s) failed, reason:(%v)", f, err)
	}
	if err = yaml.Unmarshal(confYaml, conf); err != nil {
		return nil, fmt.Errorf("Parse error: %v", err)
	}
	return conf, nil
}

func main() {
	if err := encrypterCommand.Execute(); err != nil {
		fmt.Println("Encrypt password error:", err)
		os.Exit(1)
	}
}
