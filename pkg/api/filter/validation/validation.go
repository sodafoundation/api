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

package validation

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	bctx "github.com/astaxie/beego/context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang/glog"
	myctx "github.com/sodafoundation/api/pkg/context"
)

// Factory returns a fiter function of api request
func Factory(filename string) beego.FilterFunc {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(filename)
	if err != nil {
		glog.Warningf("error loading %s api swagger file: %s", filename, err)
		return func(httpCtx *bctx.Context) {}
	}

	// Server is not required for finding route
	swagger.Servers = nil
	router := openapi3filter.NewRouter().WithSwagger(swagger)
	return func(httpCtx *bctx.Context) {
		req := httpCtx.Request
		route, pathParams, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			glog.Errorf("failed to find route from swagger: %s", err)
			myctx.HttpError(httpCtx, http.StatusBadRequest, "failed to find route %s:%s from swagger: %s", req.Method, req.URL, err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
		}
		if err := openapi3filter.ValidateRequest(context.Background(), requestValidationInput); err != nil {
			errMsg := ""
			switch e := err.(type) {
			case *openapi3filter.RequestError:
				// Retrieve first line of err message
				errMsg = strings.Split(e.Error(), "\n")[0]
			default:
				errMsg = fmt.Sprintf("%s", err)
			}
			glog.Errorf("invalid request: %s", errMsg)
			myctx.HttpError(httpCtx, http.StatusBadRequest, "%s", errMsg)
		}
	}
}
