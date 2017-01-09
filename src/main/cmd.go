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
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Create volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "show") {
				result, err := showVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Show volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "list") {
				result, err := listVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("List volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "update") {
				result, err := updateVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Update volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "delete") {
				result, err := deleteVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Delete volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "mount") {
				result, err := mountVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Mount volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
				}
			}
			if strings.Contains(cmd, "unmount") {
				result, err := unmountVolume(cmd)

				if err != nil {
					fmt.Println(err)
				} else {
					if result == "" {
						fmt.Println("Unmount volume failed!")
					} else {
						fmt.Printf("%v\n", result)
					}
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
			fmt.Println("Input Error, please input y/n!")
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

		return volumes.Create(resourceType, name, size)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
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

		return volumes.Show(resourceType, volID)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
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

	return volumes.List(resourceType, allowDetails)
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

		return volumes.Update(resourceType, volID, name)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
	}
}

func deleteVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") {
		input := make([]string, 4, 8)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[3], "=")
		resourceType := input[2]
		volID := idSlice[1]

		return volumes.Delete(resourceType, volID)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
	}
}

func mountVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "host") &&
		strings.Contains(cmd, "mountpoint") {
		input := make([]string, 6, 12)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[3], "=")
		var hostSlice []string = strings.Split(input[4], "=")
		var mountSlice []string = strings.Split(input[5], "=")
		resourceType := input[2]
		volID := idSlice[1]
		host := hostSlice[1]
		mountpoint := mountSlice[1]

		return volumes.Mount(resourceType, volID, host, mountpoint)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
	}
}

func unmountVolume(cmd string) (string, error) {
	if strings.Contains(cmd, "id") && strings.Contains(cmd, "attachment") {
		input := make([]string, 5, 10)
		input = strings.Split(cmd, "--")
		var idSlice []string = strings.Split(input[3], "=")
		var attachSlice []string = strings.Split(input[4], "=")
		resourceType := input[2]
		volID := idSlice[1]
		attachment := attachSlice[1]

		return volumes.Unmount(resourceType, volID, attachment)

	} else {
		err := fmt.Errorf("Input Error: %s", cmd)
		return "", err
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
