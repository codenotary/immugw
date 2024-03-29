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
	"net/http"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"

	immugwclient "github.com/codenotary/immugw/pkg/client"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UseDatabaseHandler ...
type UseDatabaseHandler interface {
	UseDatabase(w http.ResponseWriter, req *http.Request, pathParams map[string]string)
}

type useDatabaseHandler struct {
	mux     *runtime.ServeMux
	client  immugwclient.Client
	runtime Runtime
	json    json.JSON
}

// NewUseDatabaseHandler ...
func NewUseDatabaseHandler(mux *runtime.ServeMux, client immugwclient.Client, rt Runtime, json json.JSON) UseDatabaseHandler {
	return &useDatabaseHandler{
		mux:     mux,
		client:  client,
		runtime: rt,
		json:    json,
	}
}

func (h *useDatabaseHandler) UseDatabase(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	var (
		err    error
		client immuclient.ImmuClient
	)

	_, outboundMarshaler := h.runtime.MarshalerForRequest(h.mux, req)
	rctx, err := h.runtime.AnnotateContext(ctx, h.mux, req)
	if err != nil {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, err)
		return
	}
	var metadata runtime.ServerMetadata

	var (
		databasename string
		ok           bool
	)

	databasename, ok = pathParams["databaseName"]
	if !ok {
		h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, status.Errorf(codes.InvalidArgument, "missing parameter %s", "key"))
		return
	}

	client, err = h.client.For(databasename)
	if err == immugwclient.ErrDatabaseNotFound {
		var err error
		client, err = h.client.Add(databasename)
		if err != nil {
			h.runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, req, mapSdkError(err))
			return
		}
	}

	msg, err := client.UseDatabase(rctx, &schema.Database{DatabaseName: databasename})
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
