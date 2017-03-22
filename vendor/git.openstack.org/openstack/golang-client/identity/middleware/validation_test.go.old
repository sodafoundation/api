// +build !unit

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
	"fmt"
	"testing"
	"time"

	"git.openstack.org/openstack/golang-client/openstack"
)

const (
	cacheKeyPath = "/tmp/signing.PEM"
)

func getAuthOpts() openstack.AuthOpts {
	return openstack.AuthOpts{
		AuthUrl:     "http://localhost:5000/v2.0",
		ProjectName: "demo",
		Username:    "demo",
		Password:    "demo",
	}
}

func TestVerifyLocalPKIZ(t *testing.T) {
	token := `PKIZ_eJxdVduSqjoQfecrzvvUruEijj6GWwRNGBQC4U3AgXDTGVQuX78jzpy9z7EqVdJ091qru9P8-sV_mglt_I-ODo-HXwKybVMB1NmeY7u4pxh4qX7mtsnXdQ1VOvDMwfDBTsvrvKhyLfaQCXKzz3P0qeXVp1BUDK57UdN6mHNjotm2YSkdoyEu4whPdnnOT6MjJvK1Tpi9tDcdO3J7FpIqki1ReDjYpd0jn16REfCDA-RXS-yDJfLzG_bNHk25bLM-txuySGE9xhFiLgMMlWAQMBN7fODH8gbXP0vIOA-udWY7HbBj6M2OKRwKKpNzIg9VHNnzuwySGw3VJwNOsUyUZ9YYkm-0fZ3KeDxGmngM17eHX9KSLtFnJjzY6WIOILhME9OW1H8jIogGXCLV9R1GZVxjIxjiMhfxxOX4cY1CU0E-EF0jZ8JMp9XKuRYcNWN2Z-sOR5DqH1sSrtjH4YFsSTGsxUiuq3TkBbU6JiCzGo6lLWOlPtBgJXlS3PpV7RFLI9noZHNgi-u0jS-8DgaN9kUC1-0jeci7JcShWqVwfUna_fRjpOH-egxVElT9LOm42Yupcb7v5OFOZas7wvWUGecJ-0gS8LiWEVv0c-Fap47l-j73W3emDFptfJgl9Rkv2jHE_yZ9FrquvucAVxjihvrBSENPdn1voBNhbkhqPKUqLrMKTUhxDYvZrZQ96vbDXHhQj-ThkjTXpwRJfCCKp0j7LmJdcWfxGO3VJ6paJGHwlCYTVUiV_T1pAvZB5kAplcn4PRvz_6SxrlzGQxJPJN2_pd4TSGZg4X_Il1h_ONsTLp2Kyp6MIZVjvxJx6Sm01Gpc1txuL1wDl3SyJ8HV5yKNfEzrZ2-vb3_T9ENypQ1n1XIAfpH-YlTEUOIMmLOODv_Vxsf5mozO-lGsJCQib2uRQZO_03jymlO3qiSs5_EWkP5H0yxltNW4Cca4oWocBiK_iSItzYXr7wvXBws8FRWSAxU1TkH51AqcSoNCXGCjYLGRirFh87OvKQ_AvtfTho7I8HgSr0c6n8xIWiNqawVfOICausA3Th8ZWoS8rtc9ahDPg2bvkGAyXQR6CKTA1EFvkVCd4sgbjBIgLcdEAynSCLkJvFUi8ha9AebgrQEISfjsZ4ZJkUbnBKDv3zMlU3Z8ofBL1icwuFF5feWrTRMQ6KAe6vAArBr0ng58gB-77mfVcYYWAAe-G_WCnrqRmp_Wockkpwx39CV0BZvhXLZ2b2R4lWm1ircIOjp6Z8Dp9Oh4wtsV7PL0vicvm3fn86gV7eErk3R7a42bdukHggOOh6F6Wb1btCPZdoViT_psEaSAeY6rh4ssjr6KF602S_08sZ7mUA6i5uYaBdRongkrkOpKkOUsqqRop1b5eXz9APS1BRejWYr4a2dU4qbQGnyzbu5YFiT3P5rUNUv1w98ub4L0vowcxdJEiymnfLOXqvtp97rb7_F60jbs6y2kfgmvqvylti8Ef8rvO3OZv3pfLBET1r0IoyeudBj24IK2NLSNKhmUj9gRr_gmS6feGRdLaVGSLdhWa2Le3tCFdxePi1twmNz3ggnz98nExp9v1W82DTmf`
	validator := NewValidator(getAuthOpts(), token, cacheKeyPath, 6)
	access, err := validator.Validate()
	if err != nil {
		t.Fatal(err)
	}
	user := access.Access.User.(map[string]interface{})
	fmt.Println(user["roles"])
	project := access.Access.Token.Project
	if project.Name != "demo" {
		t.Fail()
	}
}

func TestVerifyLocalPKI(t *testing.T) {
	token := `MIIE3AYJKoZIhvcNAQcCoIIEzTCCBMkCAQExDTALBglghkgBZQMEAgEwggMqBgkqhkiG9w0BBwGgggMbBIIDF3siYWNjZXNzIjogeyJ0b2tlbiI6IHsiaXNzdWVkX2F0IjogIjIwMTYtMDUtMDNUMTk6NTE6MTIuNTM1NzQzIiwgImV4cGlyZXMiOiAiMjAxNi0wNS0wNFQxOTo1MToxMloiLCAiaWQiOiAicGxhY2Vob2xkZXIiLCAidGVuYW50IjogeyJjb3MiOiAiZGV2IiwgImRlc2NyaXB0aW9uIjogbnVsbCwgImVuYWJsZWQiOiB0cnVlLCAiaWQiOiAiMGMxNjM5OTJiY2NlNDUxZjg0NzEwMTZlMWE3MTA0ODgiLCAidnBjIjogImRldiIsICJuYW1lIjogImRlbW8ifSwgImF1ZGl0X2lkcyI6IFsiS01WV3F1U3NSYkdaTEc1Q0E2YnE2ZyJdfSwgInNlcnZpY2VDYXRhbG9nIjogW3siZW5kcG9pbnRzIjogW3siYWRtaW5VUkwiOiAiaHR0cDovL2xvY2FsaG9zdDozNTM1Ny92Mi4wIiwgInJlZ2lvbiI6ICJzdGFnZSIsICJwdWJsaWNVUkwiOiAiIiwgImlkIjogIjNkNGNmYTUyYWQ2OTQxYzViOWVlNzc5NjdkMzM3ODFiIn1dLCAiZW5kcG9pbnRzX2xpbmtzIjogW10sICJ0eXBlIjogImlkZW50aXR5IiwgIm5hbWUiOiAia2V5c3RvbmUifV0sICJ1c2VyIjogeyJ1c2VybmFtZSI6ICJkZW1vIiwgInJvbGVzX2xpbmtzIjogW10sICJpZCI6ICIzNjJkY2Q2NGY2ZTk0NjQ3YjBlNjlkY2I4ODNjYzIzOCIsICJyb2xlcyI6IFt7Im5hbWUiOiAiTWVtYmVyIn0sIHsibmFtZSI6ICJhZG1pbiJ9XSwgIm5hbWUiOiAiZGVtbyJ9LCAibWV0YWRhdGEiOiB7ImlzX2FkbWluIjogMCwgInJvbGVzIjogWyI5ZmUyZmY5ZWU0Mzg0YjE4OTRhOTA4NzhkM2U5MmJhYiIsICJmMWNhNDhiZDc0ZDI0ZDRlYTRhNTQwYmYyMDQ0YjQwMCJdfX19MYIBhTCCAYECAQEwXDBXMQswCQYDVQQGEwJVUzEOMAwGA1UECAwFVW5zZXQxDjAMBgNVBAcMBVVuc2V0MQ4wDAYDVQQKDAVVbnNldDEYMBYGA1UEAwwPd3d3LmV4YW1wbGUuY29tAgEBMAsGCWCGSAFlAwQCATANBgkqhkiG9w0BAQEFAASCAQBKR9e1K9TYOanIHpoxpCwKjnFY1Ue66+GbKVr956TLA+d2Q82IS-vJpmGRdxZh0t05knoErJEwjaq2XtA2subfIPnX6zm34y6Q1f8AJXDUowWX8YeeyRs548oCdaHoE1ak81jGOzYjMhZc-kljUlEDE4ejlO4wkxCnagDiA7uaRJgmSzB2kuuKZeeMxhlTe78tkoco3a1gZCGjGsUuEzbH5HU6RSugI5uxUGyMW0PS2j4K2+BBq2Uk-nHX0pIb513NOoDZztVq6ZuYx3KPIe-h29IMzoqL9OcZ4JH49ehzlDlTw8otu8wS8JUaIv7HNnGgJbCbsUmQOPvWju89rB3k`
	validator := NewValidator(getAuthOpts(), token, cacheKeyPath, 6)
	access, err := validator.Validate()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(access)
	project := access.Access.Token.Project
	if project.Name != "demo" {
		t.Fail()
	}
}

func TestVerifyLocalUUID(t *testing.T) {
	token := `399789012f4fbedc63c55396f59654d6`
	validator := NewValidator(getAuthOpts(), token, cacheKeyPath, 120)
	access, err := validator.Validate()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(access)
	project := access.Access.Token.Project
	if project.Name != "demo" {
		t.Fail()
	}
}

// Should be half traffic sending to keystone, not every time
func TestCache(t *testing.T) {
	for i := 0; i < 4; i++ {
		TestVerifyLocalPKIZ(t)
		TestVerifyLocalPKI(t)
		time.Sleep(3 * time.Second)
	}
}
