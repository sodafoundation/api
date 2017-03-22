// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
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

package objectstorage

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/util"
)

var zeroByte = &([]byte{}) //pointer to empty []byte

//ListContainers calls the OpenStack list containers API using
//previously obtained token.
//"limit" and "marker" corresponds to the API's "limit" and "marker".
//"url" can be regular storage or cdn-enabled storage URL.
//It returns []byte which then needs to be unmarshalled to decode the JSON.
func ListContainers(session *openstack.Session, limit int64, marker, url string) ([]byte, error) {
	return ListObjects(session, limit, marker, "", "", "", url)
}

//GetAccountMeta calls the OpenStack retrieve account metadata API using
//previously obtained token.
func GetAccountMeta(session *openstack.Session, url string) (http.Header, error) {
	return GetObjectMeta(session, url)
}

//DeleteContainer calls the OpenStack delete container API using
//previously obtained token.
func DeleteContainer(session *openstack.Session, url string) error {
	return DeleteObject(session, url)
}

//GetContainerMeta calls the OpenStack retrieve object metadata API
//using previously obtained token.
//url can be regular storage or CDN-enabled storage URL.
func GetContainerMeta(session *openstack.Session, url string) (http.Header, error) {
	return GetObjectMeta(session, url)
}

//SetContainerMeta calls the OpenStack API to create / update meta data
//for container using previously obtained token.
//url can be regular storage or CDN-enabled storage URL.
func SetContainerMeta(session *openstack.Session, url string, headers http.Header) (err error) {
	return SetObjectMeta(session, url, headers)
}

//PutContainer calls the OpenStack API to create / update
//container using previously obtained token.
func PutContainer(session *openstack.Session, url string, headers http.Header) error {
	return PutObject(session, zeroByte, url, headers)
}

//ListObjects calls the OpenStack list object API using previously
//obtained token. "Limit", "marker", "prefix", "path", "delim" corresponds
//to the API's "limit", "marker", "prefix", "path", and "delimiter".
func ListObjects(session *openstack.Session, limit int64,
	marker, prefix, path, delim, conURL string) ([]byte, error) {
	var query url.Values = url.Values{}
	query.Add("format", "json")
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(limit, 10))
	}
	if marker != "" {
		query.Add("marker", url.QueryEscape(marker))
	}
	if prefix != "" {
		query.Add("prefix", url.QueryEscape(prefix))
	}
	if path != "" {
		query.Add("path", url.QueryEscape(path))
	}
	if delim != "" {
		query.Add("delimiter", url.QueryEscape(delim))
	}
	resp, err := session.Get(conURL, &query, nil)
	if err != nil {
		return nil, err
	}
	if err = util.CheckHTTPResponseStatusCode(resp); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

//PutObject calls the OpenStack create object API using previously
//obtained token.
//url can be regular storage or CDN-enabled storage URL.
func PutObject(session *openstack.Session, fContent *[]byte, url string, headers http.Header) (err error) {
	resp, err := session.Put(url, nil, &headers, fContent)
	if err != nil {
		return err
	}
	return util.CheckHTTPResponseStatusCode(resp)
}

//CopyObject calls the OpenStack copy object API using previously obtained
//token.  Note from API doc: "The destination container must exist before
//attempting the copy."
func CopyObject(session *openstack.Session, srcURL, destURL string) (err error) {
	var headers http.Header = http.Header{}
	headers.Add("Destination", destURL)
	resp, err := session.Request("COPY", srcURL, nil, &headers, zeroByte)
	if err != nil {
		return err
	}
	return util.CheckHTTPResponseStatusCode(resp)
}

//DeleteObject calls the OpenStack delete object API using
//previously obtained token.
//
//Note from API doc: "A DELETE to a versioned object removes the current version
//of the object and replaces it with the next-most current version, moving it
//from the non-current container to the current." .. "If you want to completely
//remove an object and you have five total versions of it, you must DELETE it
//five times."
func DeleteObject(session *openstack.Session, url string) (err error) {
	resp, err := session.Delete(url, nil, nil)
	if err != nil {
		return err
	}
	return util.CheckHTTPResponseStatusCode(resp)
}

//SetObjectMeta calls the OpenStack API to create/update meta data for
//object using previously obtained token.
func SetObjectMeta(session *openstack.Session, url string, headers http.Header) (err error) {
	// headers.Add("X-Auth-Token", token)
	resp, err := session.Post(url, nil, &headers, zeroByte)
	if err != nil {
		return err
	}
	return util.CheckHTTPResponseStatusCode(resp)
}

//GetObjectMeta calls the OpenStack retrieve object metadata API using
//previously obtained token.
func GetObjectMeta(session *openstack.Session, url string) (http.Header, error) {
	resp, err := session.Head(url, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp.Header, util.CheckHTTPResponseStatusCode(resp)
}

//GetObject calls the OpenStack retrieve object API using previously
//obtained token. It returns http.Header, object / file content downloaded
//from the server, and err.
//
//Since this implementation of GetObject retrieves header info, it
//effectively executes GetObjectMeta also in addition to getting the
//object content.
func GetObject(session *openstack.Session, url string) (http.Header, []byte, error) {
	resp, err := session.Get(url, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	if err = util.CheckHTTPResponseStatusCode(resp); err != nil {
		return nil, nil, err
	}
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	}
	resp.Body.Close()
	return resp.Header, body, nil
}
