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
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"net/http"
	"sync"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifiedTxByIdHandler ...
type VerifiedTxByIdHandler interface {
	VerifiedTxById(w http.ResponseWriter, req *http.Request, pathParams map[string]string)
}

type verifiedTxByIdHandler struct {
	mux     *runtime.ServeMux
	client  client.ImmuClient
	runtime Runtime
	json    json.JSON
	sync.RWMutex
}

// NewVerifiedTxById ...
func NewVerifiedTxByIdHandler(mux *runtime.ServeMux, client client.ImmuClient, rt Runtime, json json.JSON) VerifiedTxByIdHandler {
	return &verifiedTxByIdHandler{
		mux:     mux,
		client:  client,
		runtime: rt,
		json:    json,
	}
}

func (h *verifiedTxByIdHandler) VerifiedTxById(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	_, outboundMarshaler := h.runtime.MarshalerForRequest(h.mux, req)
	rctx, err := h.runtime.AnnotateContext(ctx, h.mux, req)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	var protoReq schema.VerifiableTxRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		_   = err
	)

	val, ok = pathParams["tx"]
	if !ok {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "missing parameter %s", "key"))
		return
	}

	protoReq.Tx, err = runtime.Uint64(val)

	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "tx", err))
		return
	}

	if err := req.ParseForm(); err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", err))
		return
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, &utilities.DoubleArray{Encoding: map[string]int{"tx": 0}, Base: []int{1, 1, 0}, Check: []int{0, 1, 2}}); err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", err))
	}

	msg, err := h.client.VerifiedTxByID(rctx, protoReq.Tx)
	ctx = h.runtime.NewServerMetadataContext(ctx, metadata)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	newData, err := h.json.Marshal(msg)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}
	if _, err := w.Write(newData); err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}
}
