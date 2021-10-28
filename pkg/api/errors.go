/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

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

package api

import (
	"context"
	"github.com/codenotary/immudb/pkg/client/errors"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var (
	StatusErrMalformedPayload = status.Error(codes.InvalidArgument, "malformed payload")
)

func DefaultGWErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := err.(errors.ImmuError); ok {
		w.Header().Del("Trailer")
		contentType := marshaler.ContentType()
		w.Header().Set("Content-Type", contentType)
		var st int
		switch e.Code() {
		case errors.CodInvalidDatabaseName:
			st = http.StatusNotFound
		default:
			st = http.StatusInternalServerError
		}
		w.WriteHeader(st)
		if e.Error() != "" {
			je := map[string]string{"error": e.Error()}
			j, _ := json.DefaultJSON().Marshal(je)
			w.Write(j)
		}
		return
	}
	runtime.DefaultHTTPProtoErrorHandler(ctx, mux, marshaler, w, r, err)
}
