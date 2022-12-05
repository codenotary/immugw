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

	"github.com/codenotary/immudb/pkg/api/schema"
	immugwclient "github.com/codenotary/immugw/pkg/client"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifiedSetHandler ...
type VerifiedSetHandler interface {
	VerifiedSet(w http.ResponseWriter, req *http.Request, pathParams map[string]string)
}

type verifiedSetHandler struct {
	mux     *runtime.ServeMux
	client  immugwclient.Client
	runtime Runtime
	json    json.JSON
}

// NewVerifiedSetHandler ...
func NewVerifiedSetHandler(mux *runtime.ServeMux, client immugwclient.Client, rt Runtime, json json.JSON) VerifiedSetHandler {
	return &verifiedSetHandler{
		mux:     mux,
		client:  client,
		runtime: rt,
		json:    json,
	}
}

func (h *verifiedSetHandler) VerifiedSet(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
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

	var protoReq schema.VerifiableSetRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", berr))
		return
	}
	if err = inboundMarshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "%v", err))
		return
	}
	if protoReq.SetRequest == nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Error(codes.InvalidArgument, "incorrect JSON payload"))
		return
	}

	if len(protoReq.SetRequest.KVs) > 1 {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Error(codes.InvalidArgument, "verifiedSet accept only one key value pair"))
		return
	}

	if len(protoReq.SetRequest.KVs) == 0 {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Error(codes.InvalidArgument, "verifiedSet accept at least one key value pair"))
		return
	}

	msg, err := client.VerifiedSet(rctx, protoReq.SetRequest.KVs[0].Key, protoReq.SetRequest.KVs[0].Value)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, mapSdkError(err))
		return
	}

	ctx = h.runtime.NewServerMetadataContext(ctx, metadata)
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
