// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package signer implements signing and signature validation for opensds multi-cloud signer.
//
// Provides request signing for request that need to be signed with the Signature.
// Provides signature validation for request.
//
package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/backup/multicloud/credentials"
)

const (
	authHeaderPrefix  = "OPENSDS-HMAC-SHA256"
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

type Signature struct {
	Service string
	Region  string
	Request *http.Request
	Body    string
	Query   url.Values

	SignedHeaderValues http.Header

	credValues       credentials.Value
	requestDateTime  string
	requestDate      string
	requestPayload   string
	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
	authorization    string
}

// Signer provides sign requests that need to be signed with the Signature.
type Signer struct {
	// The authentication credentials the request will be signed against.
	Credentials *credentials.Credentials
}

// NewSigner returns a Signer pointer configured with the credentials and optional
// option values provided.
func NewSigner(credentials *credentials.Credentials, options ...func(*Signer)) *Signer {
	signer := &Signer{
		Credentials: credentials,
	}

	for _, option := range options {
		option(signer)
	}

	return signer
}

// Sign signs OpenSDS multi-cloud requests with the service, name, region,
// date time the request is signed at.
//
// Returns the signature or an error if signing the request failed.
func (signer Signer) Sign(req *http.Request, body string, service, region string, requestDateTime string, requestDate string, credentialStr string) (string, error) {

	sign := &Signature{
		Request:         req,
		requestDateTime: requestDateTime,
		requestDate:     requestDate,
		Query:           req.URL.Query(),
		Service:         service,
		Region:          region,
	}

	for key := range sign.Query {
		sort.Strings(sign.Query[key])
	}

	var err error
	sign.credValues, err = signer.Credentials.Get()
	sign.credentialString = credentialStr

	if err != nil {
		return "", err
	}

	if err := sign.build(); err != nil {
		return "", err
	}

	return sign.signature, nil
}

// ************* CREATE A SIGNATURE *************
func (sign *Signature) build() error {

	if err := sign.buildPayloadDigest(); err != nil {
		return err
	}
	unsignedHeaders := sign.Request.Header
	sign.buildCanonicalHeaders(unsignedHeaders)

	sign.buildCanonicalString() // depends on canonical headers / signed headers
	sign.buildStringToSign()    // depends on credential string
	sign.buildSignature()       // depends on string to sign

	return nil
}

// Build the canonical header list. Convert all header names to lowercase and
// remove leading spaces and trailing spaces. Convert sequential spaces in
// the header value to a single space.
func (sign *Signature) buildCanonicalHeaders(header http.Header) {
	var headers []string
	//TODO Add other Header parameters like content-type etc.
	for k, v := range header {

		canonicalKey := http.CanonicalHeaderKey(k)
		if canonicalKey != "X-Auth-Date" {
			continue
		}

		if sign.SignedHeaderValues == nil {
			sign.SignedHeaderValues = make(http.Header)
		}

		lowerCaseKey := strings.ToLower(k)
		if _, ok := sign.SignedHeaderValues[lowerCaseKey]; ok {
			sign.SignedHeaderValues[lowerCaseKey] = append(sign.SignedHeaderValues[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		sign.SignedHeaderValues[lowerCaseKey] = v
	}
	sort.Strings(headers)

	sign.signedHeaders = strings.Join(headers, ";")

	headerValues := make([]string, len(headers))
	for i, k := range headers {
		headerValues[i] = k + ":" +
			strings.Join(sign.SignedHeaderValues[k], ",")
	}
	stripExcessSpaces(headerValues)
	sign.canonicalHeaders = strings.Join(headerValues, "\n")
}

// ************* TASK 1: CREATE A CANONICAL REQUEST *************
func (sign *Signature) buildCanonicalString() {

	sign.Request.URL.RawQuery = strings.Replace(sign.Query.Encode(), "+", "%20", -1)

	canonicalURI := getURIPath(sign.Request.URL)

	sign.canonicalString = strings.Join([]string{
		// Step 1 is to define the verb (GET, POST, etc.) Already defined in the Request.
		sign.Request.Method,
		// Step 2: Create canonical URI--the part of the URI from domain to query
		canonicalURI,
		// Step 3: Create the canonical query string.Query string values must be URL-encoded
		// (space=%20). The parameters must be sorted by name.
		sign.Request.URL.RawQuery,
		// Step 4: Create the canonical headers and signed headers. Header names must be
		// trimmed and lowercase, and in alpha order \n.
		sign.canonicalHeaders + "\n",
		// Step 5: Create the list of signed headers. This lists the headers in the
		// canonical_headers list, delimited with ";" and in alpha order.
		sign.signedHeaders,
		// Step 6: Create payload hash (hash of the request body content). For GET requests,
		// the payload is an empty string ("").
		sign.requestPayload,
	}, "\n") //Step 7: Combine elements to create canonical request
	log.V(5).Infof("sign.Request.URL.RawQuery:%v", sign.Request.URL.RawQuery)
	log.V(5).Infof("sign.canonicalHeaders:%v", sign.canonicalHeaders)
	log.V(5).Infof("sign.signedHeaders:%v", sign.signedHeaders)
	log.V(5).Infof("sign.requestPayload:%v", sign.requestPayload)
	log.V(5).Infof("sign.canonicalString:%v", sign.canonicalString)
}

// ************* TASK 2: CREATE THE STRING TO SIGN *************
// Match the algorithm to the hashing algorithm, SHA-256
func (sign *Signature) buildStringToSign() {
	sign.stringToSign = strings.Join([]string{
		// Step 1: is to define the hashing Algorithm.
		authHeaderPrefix,
		// Step 2: Append the request date value which is specified with ISO8601 basic format
		// in the x-auth-date header in the format YYYYMMDD'T'HHMMSS'Z'.
		sign.requestDateTime,
		//Step 3: Append the credential scope value that includes the date, the region,service,
		// and a termination string ("sign_request") in lowercase characters.
		sign.credentialString,
		//Step 4: Append the hash of the canonical request created in Task 1
		hex.EncodeToString(makeSha256([]byte(sign.canonicalString))),
	}, "\n") //Step 5: Combine elements to create canonical request
	log.V(5).Infof("sign.stringToSign:%v", sign.stringToSign)
}

// ************* TASK 3: CALCULATE THE SIGNATURE *************
func (sign *Signature) buildSignature() {
	// Step 1:  Create the signing key, use the secret access key to create a series of
	// hash-based message authentication codes (HMACs).
	kSecret := sign.credValues.SecretAccessKey
	kDate := makeHmac([]byte("OPENSDS"+kSecret), []byte(sign.requestDate))
	kRegion := makeHmac(kDate, []byte(sign.Region))
	kService := makeHmac(kRegion, []byte(sign.Service))
	signingKey := makeHmac(kService, []byte("sign_request"))

	// Step 2: Calculate the signature using the signing key and the string to sign as inputs
	// to the keyed hash function. Convert the binary value to a hexadecimal representation.
	signature := makeHmac(signingKey, []byte(sign.stringToSign))
	sign.signature = hex.EncodeToString(signature)
	log.V(5).Infof("sign.signature:%v", sign.signature)
}

// Use a SHA256 hash to create a hashed value from the payload in the body of the HTTP request
// If the payload is empty, use an empty string as the input to the hash function.
func (sign *Signature) buildPayloadDigest() error {
	hash := ""

	if sign.Body == "" {
		hash = emptyStringSHA256
	} else {
		hash = hex.EncodeToString(makeSha256([]byte(sign.Body)))
	}

	sign.requestPayload = hash

	return nil
}

func makeHmac(key []byte, value []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(value)
	return hash.Sum(nil)
}

func makeSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// stripExcessSpaces will trim multiple side-by-side spaces.
func stripExcessSpaces(vals []string) {
	var j, k, l, m, spaces int
	for i, str := range vals {
		// Trim trailing spaces
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- {
		}

		// Trim leading spaces
		for k = 0; k < j && str[k] == ' '; k++ {
		}
		str = str[k : j+1]

		// Strip multiple spaces.
		j = strings.Index(str, "  ")
		if j < 0 {
			vals[i] = str
			continue
		}

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ {
			if buf[k] == ' ' {
				if spaces == 0 {
					buf[m] = buf[k]
					m++
				}
				spaces++
			} else {
				spaces = 0
				buf[m] = buf[k]
				m++
			}
		}

		vals[i] = string(buf[:m])
	}
}

func getURIPath(url *url.URL) string {
	var canonicalURI string

	if len(url.Opaque) > 0 {
		canonicalURI = "/" + strings.Join(strings.Split(url.Opaque, "/")[3:], "/")
	} else {
		canonicalURI = url.EscapedPath()
	}

	if len(canonicalURI) == 0 {
		canonicalURI = "/"
	}

	return canonicalURI
}
