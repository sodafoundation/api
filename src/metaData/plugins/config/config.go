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
This module defines some structs to simulate configuration information of
storage resources.

*/

package config

type DbInfo struct {
	Context, Type, Ip, Port      string
	Usrname, Pwd, Dbname, Engine string
}

type FsInfo struct {
	Context, Type     string
	BlockSize, Device string
}

func (sqlBtree *DbInfo) GetSqlBtreeInfo() *DbInfo {
	sqlBtree.Context = "database"
	sqlBtree.Type = "mysql"
	sqlBtree.Ip = "172.17.0.2"
	sqlBtree.Port = "3306"
	sqlBtree.Usrname = "root"
	sqlBtree.Pwd = "199582"
	sqlBtree.Dbname = "db1"
	sqlBtree.Engine = "InnoDB"
	return sqlBtree
}

func (sqlHash *DbInfo) GetSqlHashInfo() *DbInfo {
	sqlHash.Context = "database"
	sqlHash.Type = "mysql"
	sqlHash.Ip = "172.17.0.2"
	sqlHash.Port = "3306"
	sqlHash.Usrname = "root"
	sqlHash.Pwd = "199582"
	sqlHash.Dbname = "db1"
	sqlHash.Engine = "Memory"
	return sqlHash
}

func (mongodb *DbInfo) GetMongodbInfo() *DbInfo {
	mongodb.Context = "database"
	mongodb.Type = "mongodb"
	mongodb.Ip = "172.17.0.3"
	mongodb.Port = "27017"
	mongodb.Usrname = "root"
	mongodb.Pwd = "199582"
	mongodb.Dbname = "db1"
	mongodb.Engine = "WiredTiger"
	return mongodb
}

func (rocksdb *DbInfo) GetRocksdbInfo() *DbInfo {
	rocksdb.Context = "database"
	rocksdb.Type = "rocksdb"
	rocksdb.Ip = "172.17.0.4"
	rocksdb.Port = "27017"
	rocksdb.Usrname = "root"
	rocksdb.Pwd = "199582"
	rocksdb.Dbname = "db1"
	rocksdb.Engine = "Rocksdb"
	return rocksdb
}

func (ext4 *FsInfo) GetExt4Info() *FsInfo {
	ext4.Context = "fileSystem"
	ext4.Type = "ext4"
	ext4.BlockSize = "4096"
	ext4.Device = "/dev/sda1"
	return ext4
}

func (btrfs *FsInfo) GetBtrfsInfo() *FsInfo {
	btrfs.Context = "fileSystem"
	btrfs.Type = "btrfs"
	btrfs.BlockSize = "4096"
	btrfs.Device = "/dev/sda2"
	return btrfs
}

func (xfs *FsInfo) GetXfsInfo() *FsInfo {
	xfs.Context = "fileSystem"
	xfs.Type = "xfs"
	xfs.BlockSize = "4096"
	xfs.Device = "/dev/sda3"
	return xfs
}
