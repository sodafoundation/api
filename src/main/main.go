// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements a entry into the OpenSDS system.

*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"api/databases"
	"api/fileSystems"
	"api/volumes"
)

func main() {
	// Start the system loop.
	for {
		var cmd string
		log.Printf("Please input your action:")
		fmt.Scanln(&cmd)

		if strings.Contains(cmd, "volume") {
			if strings.Contains(cmd, "create") {
				result, err := createVolume(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "show") {
				result, err := showVolume(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "list") {
				result, err := listVolume(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "update") {
				result, err := updateVolume(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "delete") {
				result, err := deleteVolume(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
		}
		if strings.Contains(cmd, "db") {
			if strings.Contains(cmd, "create") {
				result, err := createDatabase(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "show") {
				result, err := showDatabase(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "list") {
				result, err := listDatabase()

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "update") {
				result, err := updateDatabase(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "delete") {
				result, err := deleteDatabase(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
		}
		if strings.Contains(cmd, "fs") {
			if strings.Contains(cmd, "create") {
				result, err := createFileSystem(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "show") {
				result, err := showFileSystem(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "list") {
				result, err := listFileSystem()

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "update") {
				result, err := updateFileSystem(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
			if strings.Contains(cmd, "delete") {
				result, err := deleteFileSystem(cmd)

				if err != nil {
					fmt.Println("Failure!")
				} else {
					fmt.Printf("%v\n", result)
				}
			}
		}
		var rsp string
		log.Printf("Do you want to do it again?(y/n)")
		fmt.Scanln(&rsp)

		if rsp != "y" && rsp != "n" {
			fmt.Println("Input Error!")
			break
		} else {
			if rsp == "n" {
				break
			}
		}
	}
	fmt.Println("Good Bye!")
}

func createVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "name") && strings.Contains(cmd, "size") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		nameSlice := make([]string, 2, 4)
		sizeSlice := make([]string, 2, 4)
		nameSlice = strings.Split(input[3], "=")
		sizeSlice = strings.Split(input[4], "=")
		resourceType := input[2]
		name := nameSlice[1]
		size, _ := strconv.Atoi(sizeSlice[1])

		result, err := volumes.Create(resourceType, name, size)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func showVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") {
		input := make([]string, 4, 8)
		input = strings.Split(cmd, "--")
		idSlice := make([]string, 2, 4)
		idSlice = strings.Split(input[3], "=")
		resourceType := input[2]
		volID := idSlice[1]

		result, err := volumes.Show(resourceType, volID)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func listVolume(cmd string) (string, error) {
	input := make([]string, 3, 6)
	input = strings.Split(cmd, "--")
	resourceType := input[2]
	var allowDetails bool = false
	if strings.Contains(cmd, "detail") {
		allowDetails = true
	}

	result, err := volumes.List(resourceType, allowDetails)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func updateVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[3], "=")
		var nameSlice []string = strings.Split(input[4], "=")
		resourceType := input[2]
		volID := idSlice[1]
		name := nameSlice[1]

		result, err := volumes.Update(resourceType, volID, name)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func deleteVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") {
		input := make([]string, 4, 8)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[3], "=")
		resourceType := input[2]
		volID := idSlice[1]

		result, err := volumes.Delete(resourceType, volID)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func createDatabase(cmd string) (string, error) {
	if strings.Contains(cmd, "name") && strings.Contains(cmd, "size") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		nameSlice := make([]string, 2, 4)
		sizeSlice := make([]string, 2, 4)
		nameSlice = strings.Split(input[2], "=")
		sizeSlice = strings.Split(input[3], "=")
		name := nameSlice[1]
		size, _ := strconv.Atoi(sizeSlice[1])

		result, err := databases.Create(name, size)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func showDatabase(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		idSlice := make([]string, 2, 4)
		nameSlice := make([]string, 2, 4)
		idSlice = strings.Split(input[2], "=")
		nameSlice = strings.Split(input[3], "=")
		id, _ := strconv.Atoi(idSlice[1])
		name := nameSlice[1]

		result, err := databases.Show(id, name)

		if err != nil {
			return result, err
		} else {
			return result, nil
		}
	} else {
		return "Input Error", nil
	}
}

func listDatabase() (string, error) {
	result, err := databases.List()

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func updateDatabase(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") && strings.Contains(cmd, "size") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[2], "=")
		var sizeSlice []string = strings.Split(input[4], "=")
		var nameSlice []string = strings.Split(input[3], "=")
		id, _ := strconv.Atoi(idSlice[1])
		size, _ := strconv.Atoi(sizeSlice[1])
		name := nameSlice[1]

		result, err := databases.Update(id, size, name)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func deleteDatabase(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") && strings.Contains(cmd, "cascade") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[2], "=")
		var nameSlice []string = strings.Split(input[3], "=")
		var cascadeSlice []string = strings.Split(input[4], "=")
		id, _ := strconv.Atoi(idSlice[1])
		name := nameSlice[1]
		cascade, _ := strconv.ParseBool(cascadeSlice[1])

		result, err := databases.Delete(id, name, cascade)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func createFileSystem(cmd string) (string, error) {
	if strings.Contains(cmd, "name") && strings.Contains(cmd, "size") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		nameSlice := make([]string, 2, 4)
		sizeSlice := make([]string, 2, 4)
		nameSlice = strings.Split(input[2], "=")
		sizeSlice = strings.Split(input[3], "=")
		name := nameSlice[1]
		size, _ := strconv.Atoi(sizeSlice[1])

		result, err := fileSystems.Create(name, size)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func showFileSystem(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		idSlice := make([]string, 2, 4)
		nameSlice := make([]string, 2, 4)
		idSlice = strings.Split(input[2], "=")
		nameSlice = strings.Split(input[3], "=")
		id, _ := strconv.Atoi(idSlice[1])
		name := nameSlice[1]

		result, err := fileSystems.Show(id, name)

		if err != nil {
			return result, err
		} else {
			return result, nil
		}
	} else {
		return "Input Error", nil
	}
}

func listFileSystem() (string, error) {
	result, err := fileSystems.List()

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func updateFileSystem(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") && strings.Contains(cmd, "size") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[2], "=")
		var sizeSlice []string = strings.Split(input[4], "=")
		var nameSlice []string = strings.Split(input[3], "=")
		id, _ := strconv.Atoi(idSlice[1])
		size, _ := strconv.Atoi(sizeSlice[1])
		name := nameSlice[1]

		result, err := fileSystems.Update(id, size, name)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}

func deleteFileSystem(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "name") && strings.Contains(cmd, "cascade") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[2], "=")
		var nameSlice []string = strings.Split(input[3], "=")
		var cascadeSlice []string = strings.Split(input[4], "=")
		id, _ := strconv.Atoi(idSlice[1])
		name := nameSlice[1]
		cascade, _ := strconv.ParseBool(cascadeSlice[1])

		result, err := fileSystems.Delete(id, name, cascade)

		if err != nil {
			return "Error", err
		} else {
			return result, nil
		}
	} else {
		return "Input Error!", nil
	}
}
