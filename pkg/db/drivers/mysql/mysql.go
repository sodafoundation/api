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
This module implements the mysql database operation of data structure
defined in api module.

*/

package mysql

import (
	"database/sql"
	"sync"
)

var c = &client{}

// Init
func Init(driver, credential string) *client {
	// driver equals "mysql"
	db, err := sql.Open(driver, credential)
	if err != nil {
		db.Close()
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	if err = db.Ping(); err != nil {
		db.Close()
		panic(err) // proper error handling instead of panic in your app
	}

	c.cli = db
	return c
}

type client struct {
	cli  *sql.DB
	lock sync.Mutex
}
