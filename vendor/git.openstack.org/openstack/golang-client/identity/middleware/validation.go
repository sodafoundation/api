// Copyright (c) 2016 eBay Inc.
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

package middleware

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"git.openstack.org/openstack/golang-client/openstack"
	"github.com/fullsailor/pkcs7"
)

const (
	PKI_ASN1_PREFIX = "MII"
	PKIZ_PREFIX     = "PKIZ_"
)

// Cache the token until it gets expired
var serviceTokenSession *openstack.Session

// Cache token revocation list until expired
var revocationListCache *revokedListCache

// NewValidator gets the credential for service account, token need to be validated,
// signing cert location (will store the cert from keystone if not there), and the
// revocation list cache duration (in seconds) and returns the validator.
func NewValidator(authOpts openstack.AuthOpts, token string, signingKeyPath string, revCacheSecs int) *Validator {
	return &Validator{
		SvcAuthOpts:          authOpts,
		CachedSigningKeyPath: signingKeyPath,
		TokenId:              token,
		RevCacheDuration:     time.Duration(revCacheSecs) * time.Second,
	}
}

// Validate does the local validation for PKI & PKIZ token and sends to keystone
// for other format tokens validation. It returns the extracted AuthToken struct
func (validator *Validator) Validate() (*openstack.AuthToken, error) {
	var token *openstack.AuthToken

	if strings.HasPrefix(validator.TokenId, PKIZ_PREFIX) ||
		strings.HasPrefix(validator.TokenId, PKI_ASN1_PREFIX) {
		// do local validation for PKI and PKIZ token
		if revocationListCache == nil || time.Now().Sub(revocationListCache.Time) > validator.RevCacheDuration {
			_, err := validator.getRevocationList()
			if err != nil {
				return nil, err
			}
		}
		access, err := validator.ValidateOffline()
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(access, &token); err != nil {
			return nil, err
		}

		// set the token ID
		token.Access.Token.ID = validator.TokenId

		hashedTokenId := fmt.Sprintf("%x", md5.Sum([]byte(validator.TokenId)))
		for _, rtoken := range revocationListCache.Revoked {
			if rtoken.ID == hashedTokenId {
				return nil, fmt.Errorf("token %s was revoked", hashedTokenId)
			}
		}
	} else {
		access, err := validator.ValidateRemote()
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(access, &token); err != nil {
			return nil, err
		}
	}

	// validation should fail if token is expired
	if token.Access.Token.Expires.Sub(time.Now()) < 0 {
		return nil, fmt.Errorf("token %s is expired", validator.TokenId)
	}
	return token, nil
}

// Validate the token locally without sending to keystone
// It take the token body and return the extracted access object as []byte
func (validator *Validator) ValidateOffline() ([]byte, error) {
	token := ""
	switch {
	case strings.HasPrefix(validator.TokenId, PKIZ_PREFIX):
		token = strings.TrimPrefix(validator.TokenId, PKIZ_PREFIX)
		decompressedToken, err := decompressToken(token)
		if err != nil {
			return nil, err
		}

		token = trimCMSFormat(decompressedToken)
	case strings.HasPrefix(validator.TokenId, PKI_ASN1_PREFIX):
		token = validator.TokenId

	default:
		return nil, errors.New("can not validate offline, it has to be sent to keystone")
	}

	decodedToken, err := base64DecodeFromCms(token)
	if err != nil {
		return nil, err
	}

	content, err := validator.checkSignature(decodedToken)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (validator *Validator) ValidateRemote() ([]byte, error) {
	resp, err := validator.reqTokenValidation()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// retry one more time using the new service token
		err = validator.renewToken()
		if err != nil {
			return nil, err
		}
		resp, err = validator.reqTokenValidation()
		if err != nil {
			return nil, err
		}
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}
	return rbody, nil
}

func decompressToken(token string) (string, error) {
	decToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	zr, err := zlib.NewReader(bytes.NewBuffer(decToken))
	if err != nil {
		return "", err
	}
	bb, err := ioutil.ReadAll(zr)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func base64DecodeFromCms(token string) ([]byte, error) {
	t := strings.Replace(token, "-", "/", -1)
	decToken, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return nil, err
	}
	return decToken, nil
}

// remove the customerized header and footer in PEM token
// -----BEGIN CMS-----
// -----END CMS-----
func trimCMSFormat(token string) string {
	token = strings.Trim(token, "\n")
	l := strings.Index(token, "\n")
	r := strings.LastIndex(token, "\n")
	return token[l:r]
}

// Get the signging certificate from local dir
// It will get the cert from keystone if cache file does not exist
func (validator *Validator) getSigningCert() (*x509.Certificate, error) {
	signPEM, err := ioutil.ReadFile(validator.CachedSigningKeyPath)
	if err != nil {
		resp, err := http.Get(validator.SvcAuthOpts.AuthUrl + "/certificates/signing")
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, errors.New("can not get signing cert")
		}
		signPEM, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// cache the file to location
		if err = ioutil.WriteFile(validator.CachedSigningKeyPath, []byte(signPEM), 0644); err != nil {
			log.Println("error caching signging cert")
		}
	}

	block, _ := pem.Decode(signPEM)
	if block == nil {
		return nil, errors.New("can not decode PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// check the signature of the token
func (validator *Validator) checkSignature(data []byte) ([]byte, error) {
	p7, err := pkcs7.Parse(data)
	if err != nil {
		return nil, err
	}

	if len(p7.Signers) != 1 {
		return nil, errors.New("should be only one signature found")
	}

	signer := p7.Signers[0]
	cert, err := validator.getSigningCert()
	if err != nil {
		return nil, err
	}

	err = cert.CheckSignature(x509.SHA256WithRSA, p7.Content, signer.EncryptedDigest)
	if err != nil {
		return nil, err
	}

	return p7.Content, nil
}

// getRevocationList get a list of revoked tokens
func (validator *Validator) getRevocationList() ([]openstack.Token, error) {
	// Get the service token to get the token revocation list
	resp, err := validator.reqRevocationList()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// try again when getting 401, it could be the cached token was revoked
		// or the token got expired
		validator.renewToken()
		resp, err = validator.reqRevocationList()
		if err != nil {
			return nil, err
		}
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	revokeCMSResp := revokeResp{}
	if err = json.Unmarshal(rbody, &revokeCMSResp); err != nil {
		return nil, err
	}

	revokeResp := trimCMSFormat(revokeCMSResp.Signed)
	decodedResp, err := base64DecodeFromCms(revokeResp)
	if err != nil {
		return nil, err
	}

	revoked, err := validator.checkSignature(decodedResp)
	if err != nil {
		return nil, err
	}
	revokedList := RevokedList{}
	if err = json.Unmarshal(revoked, &revokedList); err != nil {
		return nil, err
	}

	// update the revocation list cache
	revocationListCache = &revokedListCache{
		Revoked: revokedList.Revoked,
		Time:    time.Now(),
	}
	return revokedList.Revoked, nil
}

// reqRevocationList sends the GET request to keystone and get response back
func (validator *Validator) reqRevocationList() (*http.Response, error) {
	if serviceTokenSession == nil {
		validator.renewToken()
	}
	// Get token revocation list
	return serviceTokenSession.Request("GET", validator.SvcAuthOpts.AuthUrl+"/tokens/revoked", nil, nil, nil)
}

func (validator *Validator) reqTokenValidation() (*http.Response, error) {
	if serviceTokenSession == nil {
		validator.renewToken()
	}
	// Get token revocation list
	reqUrl := fmt.Sprintf("%s/tokens/%s", validator.SvcAuthOpts.AuthUrl, validator.TokenId)
	return serviceTokenSession.Request("GET", reqUrl, nil, nil, nil)
}

// renewToken gets the keystone token from service AuthOpts
func (validator *Validator) renewToken() error {
	// Get the service token to get the token revocation list
	auth, err := openstack.DoAuthRequest(validator.SvcAuthOpts)
	if err != nil {
		return err
	}

	// Make a new client with these creds
	serviceTokenSession, err = openstack.NewSession(nil, auth, nil)
	if err != nil {
		return err
	}
	return nil
}
