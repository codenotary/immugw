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
	"net/http"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/clienttest"
	immugwclient "github.com/codenotary/immugw/pkg/client"

	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func testHistoryHandler(t *testing.T, mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) {
	prefixPattern := "HistoryHandler - Test case: %s"
	method := "POST"
	for _, tc := range historyHandlerTestCases(mux, client, opts) {
		path := "/db/defaultdb/history/"
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			tc.historyHandler.History(res, req, defaultTestParams)
		}
		err := testHandler(
			t,
			fmt.Sprintf(prefixPattern, tc.name),
			method,
			path,
			tc.payload,
			handlerFunc,
			tc.testFunc,
		)
		require.NoError(t, err)
	}
}

type historyHandlerTestCase struct {
	name           string
	historyHandler HistoryHandler
	payload        string
	testFunc       func(*testing.T, string, int, map[string]interface{})
}

func historyHandlerTestCases(mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) []historyHandlerTestCase {
	rt := newDefaultRuntime()
	defaultJSON := json.DefaultJSON()
	hh := NewHistoryHandler(mux, client, rt, defaultJSON)
	icd, _ := immuclient.NewImmuClient(immuclient.DefaultOptions())
	historyWErr := func(context.Context, *schema.HistoryRequest) (*schema.Entries, error) {
		return nil, errors.New("history error")
	}
	validPayload := fmt.Sprintf(
		`{
			"key": "%s"
		}`,
		base64.StdEncoding.EncodeToString([]byte("setKey1")),
	)

	return []historyHandlerTestCase{
		{
			"Sending correct request",
			hh,
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
				items, ok := body["entries"]
				require.True(t, ok, "%sfield \"entries\" not found in response %v", testCase, body)
				notEmptyMsg := "%sexpected more than on item in response %v"
				require.True(t, len(items.([]interface{})) > 0, notEmptyMsg, testCase, body)
			},
		},
		{
			"Sending correct request with non-existent key",
			hh,
			fmt.Sprintf(
				`{
			"key": "%s"
		}`, base64.StdEncoding.EncodeToString([]byte("historyNonExistentKey1"))),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusNotFound, status)
			},
		},
		{
			"Missing key path param",
			hh,
			`{}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusNotFound, status)
			},
		},
		{
			"Sending plain text instead of base64 encoded",
			hh,
			`{
				"key": "setKey1"
			}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expectedError :=
					"illegal base64 data at input byte 4"
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": expectedError}, body)
			},
		},
		{
			"AnnotateContext error",
			NewHistoryHandler(mux, client, newTestRuntimeWithAnnotateContextErr(), defaultJSON),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "annotate context error"}, body)
			},
		},
		{
			"History error",
			NewHistoryHandler(mux, immugwclient.NewMockClient(&clienttest.ImmuClientMock{ImmuClient: icd, HistoryF: historyWErr}, opts), rt, defaultJSON),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "history error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewHistoryHandler(mux, client, rt, newTestJSONWithMarshalErr()),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "JSON marshal error"}, body)
			},
		},
	}
}
