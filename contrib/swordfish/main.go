/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	model "github.com/opensds/opensds/contrib/swordfish/proto"
	"github.com/opensds/opensds/contrib/swordfish/provider"
	"github.com/opensds/opensds/contrib/swordfish/utils"

	"github.com/gorilla/mux"
)

var (
	port int

	StorageServicePrefix = "/redfish/v1/StorageServices(1)"
	ClassOfServicePrefix = StorageServicePrefix + "/ClassesOfService"

	DataProtectionLineOfServicePrefix        = StorageServicePrefix + "/Links/DataProtectionLoSCapabilities"
	DataSecurityLineOfServicePrefix          = StorageServicePrefix + "/Links/DataSecurityLoSCapabilities"
	DataStorageLineOfServicePrefix           = StorageServicePrefix + "/Links/DataStorageLoSCapabilities"
	IOConnectivityLineOfServicePrefix        = StorageServicePrefix + "/Links/IOConnectivityLoSCapabilities"
	IOPerformanceLineOfServicePrefix         = StorageServicePrefix + "/Links/IOPerformanceLoSCapabilities"
	GetSupportedDataProtectionLinesOfService = DataProtectionLineOfServicePrefix + "/SupportedDataProtectionLinesOfService"
	GetSupportedDataSecurityLinesOfService   = DataSecurityLineOfServicePrefix + "/SupportedDataSecurityLinesOfService"
	GetSupportedDataStorageLinesOfService    = DataStorageLineOfServicePrefix + "/SupportedDataStorageLinesOfService"
	GetSupportedIOConnectivityLinesOfService = IOConnectivityLineOfServicePrefix + "/SupportedIOConnectivityLinesOfService"
	GetSupportedIOPerformanceLinesOfService  = IOPerformanceLineOfServicePrefix + "/SupportedIOPerformanceLinesOfService"
)

func init() {
	flag.IntVar(&port, "port", 8081, "use '--port' option to specify the port to listen on")
	flag.Parse()
}

type server struct {
	prov *provider.Provider
}

// CreateHandler creates Broker HTTP handler based on an implementation
// of a controller.Controller interface.
func CreateHandler(p *provider.Provider) http.Handler {
	s := server{
		prov: p,
	}

	var router = mux.NewRouter()

	router.HandleFunc(ClassOfServicePrefix, s.createClassOfService).Methods("POST")
	router.HandleFunc(ClassOfServicePrefix, s.getClassesOfService).Methods("GET")

	router.HandleFunc(DataProtectionLineOfServicePrefix, s.createDataProtectionLineOfService).Methods("PATCH")
	router.HandleFunc(DataSecurityLineOfServicePrefix, s.createDataSecurityLineOfService).Methods("PATCH")
	router.HandleFunc(DataStorageLineOfServicePrefix, s.createDataStorageLineOfService).Methods("PATCH")
	router.HandleFunc(IOConnectivityLineOfServicePrefix, s.createIOConnectivityLineOfService).Methods("PATCH")
	router.HandleFunc(IOPerformanceLineOfServicePrefix, s.createIOPerformanceLineOfService).Methods("PATCH")

	router.HandleFunc(GetSupportedDataProtectionLinesOfService, s.getSupportedDataProtectionLinesOfService).Methods("GET")
	router.HandleFunc(GetSupportedDataSecurityLinesOfService, s.getSupportedDataSecurityLinesOfService).Methods("GET")
	router.HandleFunc(GetSupportedDataStorageLinesOfService, s.getSupportedDataStorageLinesOfService).Methods("GET")
	router.HandleFunc(GetSupportedIOConnectivityLinesOfService, s.getSupportedIOConnectivityLinesOfService).Methods("GET")
	router.HandleFunc(GetSupportedIOPerformanceLinesOfService, s.getSupportedIOPerformanceLinesOfService).Methods("GET")

	return router
}

func (s *server) createClassOfService(w http.ResponseWriter, r *http.Request) {
	var cos model.ClassOfService

	if err := util.BodyToObject(r, &cos); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	result, err := s.prov.CreateClassOfService(&cos)
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("Location", result)
	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) getClassesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetClassesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusOK, result)
	return
}

func (s *server) createDataProtectionLineOfService(w http.ResponseWriter, r *http.Request) {
	var req model.DataProtectionLoSCapabilities
	var err error

	if err = util.BodyToObject(r, &req); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = s.prov.CreateDataProtectionLineOfService(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) createDataSecurityLineOfService(w http.ResponseWriter, r *http.Request) {
	var req model.DataSecurityLoSCapabilities
	var err error

	if err = util.BodyToObject(r, &req); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = s.prov.CreateDataSecurityLineOfService(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) createDataStorageLineOfService(w http.ResponseWriter, r *http.Request) {
	var req model.DataStorageLoSCapabilities
	var err error

	if err = util.BodyToObject(r, &req); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = s.prov.CreateDataStorageLineOfService(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) createIOConnectivityLineOfService(w http.ResponseWriter, r *http.Request) {
	var req model.IOConnectivityLoSCapabilities
	var err error

	if err = util.BodyToObject(r, &req); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = s.prov.CreateIOConnectivityLineOfService(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) createIOPerformanceLineOfService(w http.ResponseWriter, r *http.Request) {
	var req model.IOPerformanceLoSCapabilities
	var err error

	if err = util.BodyToObject(r, &req); err != nil {
		log.Printf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	if err = s.prov.CreateIOPerformanceLineOfService(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	util.WriteResponse(w, http.StatusCreated, nil)

	return
}

func (s *server) getSupportedDataProtectionLinesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetSupportedDataProtectionLinesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("ETag", "123-a")
	util.WriteResponse(w, http.StatusCreated, result)

	return
}

func (s *server) getSupportedDataSecurityLinesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetSupportedDataSecurityLinesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("ETag", "123-a")
	util.WriteResponse(w, http.StatusOK, result)

	return
}

func (s *server) getSupportedDataStorageLinesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetSupportedDataStorageLinesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("ETag", "123-a")
	util.WriteResponse(w, http.StatusOK, result)

	return
}

func (s *server) getSupportedIOConnectivityLinesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetSupportedIOConnectivityLinesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("ETag", "123-a")
	util.WriteResponse(w, http.StatusOK, result)

	return
}

func (s *server) getSupportedIOPerformanceLinesOfService(w http.ResponseWriter, r *http.Request) {
	result, err := s.prov.GetSupportedIOPerformanceLinesOfService()
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}

	w.Header().Set("ETag", "123-a")
	util.WriteResponse(w, http.StatusOK, result)

	return
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.Handle("/", CreateHandler(provider.NewProvider()))

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
