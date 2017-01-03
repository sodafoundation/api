// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

//
// Request Tests
//

type hfunc http.HandlerFunc

type payload struct {
	Foo string
}

var reqTests = []struct {
	method  string
	params  bool
	payload bool
}{
	{"GET", true, false},
	{"POST", false, true},
	{"PUT", false, true},
	{"DELETE", false, false},
}

type pair struct {
	r  Request
	hf hfunc
}

func paramHandler(t *testing.T, p url.Values, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		q := req.URL.Query()
		for k, _ := range p {
			if !assert.Equal(t, p[k], q[k]) {
				msg := "Bad query params: " + q.Encode()
				t.Error(msg)
				return
			}
		}
	}
}

func payloadHandler(t *testing.T, p payload, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		if req.ContentLength <= 0 {
			t.Error("Content-Length must be greater than 0.")
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("Bad content type")
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error("Body is nil")
			return
		}
		var s payload
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Error("JSON Unmarshal failed: ", err)
			return
		}
		if s != p {
			t.Error("Bad request body")
			return
		}
	}
}

func methodHandler(t *testing.T, method string, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		if req.Method != method {
			t.Error("Incorrect method, got ", req.Method, " expected ", method)
		}
	}
}

func headerHandler(t *testing.T, h http.Header, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		for key, _ := range h {
			expected := h.Get(key)
			actual := req.Header.Get(key)
			if expected != actual {
				t.Error("Missing/bad header")
			}
			return
		}
	}
}

func TestRequest(t *testing.T) {
	// NOTE:  Do we really need to test different combinations for different
	// HTTP methods?
	pairs := []pair{}
	for _, test := range reqTests {
		baseReq := Request{
			Method: test.method,
		}
		allReq := baseReq // allRR has all supported attribues for this verb
		var allHF hfunc   // allHF is combination of all relevant handlers
		//
		// Generate a random key/value pair
		//
		key, err := randutil.AlphaString(8)
		if err != nil {
			t.Error(err)
		}
		value, err := randutil.AlphaString(8)
		if err != nil {
			t.Error(err)
		}
		//
		// Method
		//
		r := baseReq
		f := methodHandler(t, test.method, nil)
		allHF = methodHandler(t, test.method, allHF)
		pairs = append(pairs, pair{r, f})
		//
		// Header
		//
		h := http.Header{}
		h.Add(key, value)
		r = baseReq
		r.Header = &h
		allReq.Header = &h
		f = headerHandler(t, h, nil)
		allHF = headerHandler(t, h, allHF)
		pairs = append(pairs, pair{r, f})
		//
		// Params
		//
		if test.params {
			p := Params{key: value}.AsUrlValues()
			f := paramHandler(t, p, nil)
			allHF = paramHandler(t, p, allHF)
			r = baseReq
			r.Params = &p
			allReq.Params = &p
			pairs = append(pairs, pair{r, f})
		}
		//
		// Payload
		//
		if test.payload {
			p := payload{value}
			f = payloadHandler(t, p, nil)
			allHF = payloadHandler(t, p, allHF)
			r = baseReq
			r.Payload = p
			allReq.Payload = p
			pairs = append(pairs, pair{r, f})
		}
		//
		// All
		//
		pairs = append(pairs, pair{allReq, allHF})
	}
	for _, p := range pairs {
		srv := httptest.NewServer(http.HandlerFunc(p.hf))
		defer srv.Close()
		//
		// Good request
		//
		p.r.Url = "http://" + srv.Listener.Addr().String()
		_, err := Send(&p.r)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidTLS(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(HandleInvalidTLS))
	defer srv.Close()
	// The first request, which is supposed to fail, will print something similar to
	// "20:45:27 server.go:2161: http: TLS handshake error from 127.0.0.1:56293: remote error: bad certificate" to the console.
	// NOTE: Is this something that should be capture and silently ignored?
	s := Session{}
	r := Request{
		Url:    "https://" + srv.Listener.Addr().String(),
		Method: "GET",
	}
	_, err := s.Send(&r)
	if err == nil {
		t.Fatal("Invalid TLS without custom Transport object. The request should have errored out!")
	}

	s2 := Session{}
	r2 := Request{
		Url:    "https://" + srv.Listener.Addr().String(),
		Method: "GET",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp2, err2 := s2.Send(&r2)
	if err2 != nil {
		t.Fatal(err2)
	}
	if resp2.Status() != http.StatusOK {
		t.Fatalf("Expected status %d but got %v\n", http.StatusOK, resp2.Status())
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGetBasicAuth))
	defer srv.Close()
	s := Session{}
	r := Request{
		Url:      "http://" + srv.Listener.Addr().String(),
		Method:   "GET",
		Userinfo: url.UserPassword("jtkirk", "Beam me up, Scotty!"),
	}
	resp, err := s.Send(&r)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status() != 200 {
		t.Fatalf("Expected status 200 but got %v\n", resp.Status())
	}
}

func TestBasicUrlAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGetBasicAuth))
	defer srv.Close()
	s := Session{}
	testUrl, _ := url.Parse("http://" + srv.Listener.Addr().String())
	testUrl.User = url.UserPassword("jtkirk", "Beam me up, Scotty!")
	r := Request{
		Url:    testUrl.String(),
		Method: "GET",
	}
	resp, err := s.Send(&r)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status() != 200 {
		t.Fatalf("Expected status 200 but got %v\n", resp.Status())
	}
}

//
// TODO: Response Tests
//

func TestErrMsg(t *testing.T) {}

func TestStatus(t *testing.T) {}

func TestUnmarshall(t *testing.T) {}

// func TestUnmarshallFail() {}

func HandleInvalidTLS(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HandleGetBasicAuth(w http.ResponseWriter, req *http.Request) {
	authRegex := regexp.MustCompile(`[Bb]asic (?P<encoded>\S+)`)
	str := req.Header.Get("Authorization")
	matches := authRegex.FindStringSubmatch(str)
	if len(matches) != 2 {
		msg := "Regex doesn't match"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	encoded := matches[1]
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		msg := "Base64 decode failed"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		msg := "String split failed"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	username := parts[0]
	password := parts[1]
	if username != "jtkirk" || password != "Beam me up, Scotty!" {
		code := http.StatusUnauthorized
		text := http.StatusText(code)
		http.Error(w, text, code)
		return
	}
	w.WriteHeader(200)
}
