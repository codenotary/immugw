/*
Copyright 2019-2020 vChain, Inc.

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

package gw

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/codenotary/immugw/pkg/api"
	immugwclient "github.com/codenotary/immugw/pkg/client"

	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifiedGetHandler ...
type VerifiedSQLGetHandler interface {
	VerifiedSQLGetHandler(w http.ResponseWriter, req *http.Request, pathParams map[string]string)
}

type verifiedSQLGetHandler struct {
	mux     *runtime.ServeMux
	client  immugwclient.Client
	runtime Runtime
	json    json.JSON
	sync.RWMutex
}

// NewVerifiedSQLGetHandler ...
func NewVerifiedSQLGetHandler(mux *runtime.ServeMux, client immugwclient.Client, rt Runtime, json json.JSON) VerifiedSQLGetHandler {
	return &verifiedSQLGetHandler{
		mux:     mux,
		client:  client,
		runtime: rt,
		json:    json,
	}
}

func (h *verifiedSQLGetHandler) VerifiedSQLGetHandler(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	inboundMarshaler, outboundMarshaler := h.runtime.MarshalerForRequest(h.mux, req)
	rctx, err := h.runtime.AnnotateContext(ctx, h.mux, req)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	databasename, ok := pathParams["databaseName"]
	if !ok {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "missing parameter %s", "databaseName"))
		return
	}
	client, err := h.client.For(databasename)
	if err != nil {
		h.runtime.HTTPError(rctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	var protoReq api.VerifyRowSQLRequest
	var metadata runtime.ServerMetadata

	//protoReq.PkValues = make([]*schema.SQLValue, 0)
	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		h.runtime.HTTPError(rctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", berr))
		return
	}
	if err = inboundMarshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		h.runtime.HTTPError(rctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", err))
		return
	}

	err = client.VerifyRow(rctx, protoReq.Row, protoReq.Table, protoReq.PkValues)
	ctx = h.runtime.NewServerMetadataContext(ctx, metadata)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, mapSdkError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
