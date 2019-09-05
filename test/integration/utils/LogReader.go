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
package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ReadFile(fileName string) {
	// Open file and create scanner on top of it
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	// Could also use a custom function of SplitFunc type
	scanner.Split(bufio.ScanWords)

	// Scan for next token.
	success := scanner.Scan()
	if success == false {
		// False on error or EOF. Check error
		err = scanner.Err()
		if err == nil {
			log.Println("Scan completed and reached EOF")
		} else {
			log.Fatal(err)
		}
	}

	// Get data from scan with Bytes() or Text()
	fmt.Println("First word found:", scanner.Text())

	// Call scanner.Scan() again to find next token
}
func ReadAndFindTextInFile(fileName string, text string) bool {
	found := false
	/* ioutil.ReadFile returns []byte, error */
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	/* ... omitted error check..and please add ... */
	/* find index of newline */
	file := string(data)
	line := 0
	/* func Split(s, sep string) []string */
	temp := strings.Split(file, "\n")

	for _, item := range temp {
		fmt.Println("[", line, "]\t", item)
		if strings.Contains(item, text) {
			fmt.Println("Found given text in line [", line, "]\t")
			found = true
		}
		line++
	}
	return found
}