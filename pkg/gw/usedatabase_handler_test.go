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
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/clienttest"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func testUseDatabaseHandler(t *testing.T, ctx context.Context, mux *runtime.ServeMux, ic immuclient.ImmuClient) {
	prefixPattern := "UseDatabaseHandler - Test case: %s"
	method := "GET"
	md, _ := metadata.FromOutgoingContext(ctx)
	for _, tc := range useDatabaseHandlerTestCases(mux, ic, ctx) {
		path := "/db/use/" + tc.dbname
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			var pathParams map[string]string
			if tc.dbname != "" {
				pathParams = map[string]string{"databaseName": tc.dbname}
			}
			req.Header.Set("authorization", md.Get("authorization")[0])
			tc.useDatabaseHandler.UseDatabase(res, req, pathParams)
		}
		err := testHandler(
			t,
			fmt.Sprintf(prefixPattern, tc.name),
			method,
			path,
			"",
			handlerFunc,
			tc.testFunc,
		)
		require.NoError(t, err)
	}
}

type useDatabaseHandlerTestCase struct {
	name               string
	useDatabaseHandler UseDatabaseHandler
	dbname             string
	testFunc           func(*testing.T, string, int, map[string]interface{})
}

func useDatabaseHandlerTestCases(mux *runtime.ServeMux, ic immuclient.ImmuClient, ctx context.Context) []useDatabaseHandlerTestCase {
	rt := newDefaultRuntime()
	defaultJSON := json.DefaultJSON()
	hh := NewUseDatabaseHandler(mux, ic, rt, defaultJSON)
	icd, _ := client.NewImmuClient(client.DefaultOptions())
	useDatabaseWErr := func(ctx context.Context, d *schema.Database) (*schema.UseDatabaseReply, error) {
		return nil, errors.New("useDatabase error")
	}
	dbName := "dbtest"

	err := ic.CreateDatabase(ctx, &schema.DatabaseSettings{
		DatabaseName: dbName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return []useDatabaseHandlerTestCase{
		{
			"Sending correct request",
			hh,
			dbName,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
				_, ok := body["token"]
				require.True(t, ok, "%sfield \"token\" not found in response %v", testCase, body)
			},
		},
		{
			"Sending correct request with non-existent key",
			hh,
			base64.StdEncoding.EncodeToString([]byte("useDatabaseNonExistentKey1")),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusNotFound, status)
			},
		},
		{
			"Missing key path param",
			hh,
			"",
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"message": "missing parameter key"}, body)
			},
		},
		{
			"AnnotateContext error",
			NewUseDatabaseHandler(mux, ic, newTestRuntimeWithAnnotateContextErr(), defaultJSON),
			dbName,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"message": "annotate context error"}, body)
			},
		},
		{
			"UseDatabase error",
			NewUseDatabaseHandler(mux, &clienttest.ImmuClientMock{ImmuClient: icd, UseDatabaseF: useDatabaseWErr}, rt, defaultJSON),
			dbName,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"message": "useDatabase error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewUseDatabaseHandler(mux, ic, rt, newTestJSONWithMarshalErr()),
			dbName,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"message": "JSON marshal error"}, body)
			},
		},
	}
}
