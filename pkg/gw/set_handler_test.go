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
	immugwclient "github.com/codenotary/immugw/pkg/client"

	"github.com/codenotary/immudb/pkg/client/clienttest"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func testSetHandler(t *testing.T, mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) {
	prefixPattern := "SetHandler - Test case: %s"
	method := "POST"
	path := "/db/defaultdb/set"
	for _, tc := range setHandlerTestCases(mux, client, opts) {
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			tc.setHandler.Set(res, req, defaultTestParams)
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

type setHandlerTestCase struct {
	name       string
	setHandler SetHandler
	payload    string
	testFunc   func(*testing.T, string, int, map[string]interface{})
}

func setHandlerTestCases(mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) []setHandlerTestCase {
	rt := newDefaultRuntime()
	json := json.DefaultJSON()
	sh := NewSetHandler(mux, client, rt, json)
	icd, _ := immuclient.NewImmuClient(immuclient.DefaultOptions())

	setAllWErr := func(context.Context, *schema.SetRequest) (*schema.TxHeader, error) {
		return nil, errors.New("set error")
	}

	validKey := base64.StdEncoding.EncodeToString([]byte("setKey1"))
	validValue := base64.StdEncoding.EncodeToString([]byte("setValue1"))
	validPayload := fmt.Sprintf(
		`{
  "KVs": [
    {
        "key": "%s",
        "value": "%s"
    }
  ]
}`,
		validKey,
		validValue,
	)

	return []setHandlerTestCase{
		{
			"Sending correct request",
			sh,
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
				requireResponseFields(t, testCase, []string{"id"}, body)
			},
		},
		/*		{
				"Missing value field",
				sh,
				fmt.Sprintf(
					`{
					  "KVs": [
						{
							"key": "%s"
						}
					  ]
					}`,
					validKey,
				),
				func(t *testing.T, testCase string, status int, body map[string]interface{}) {
					requireResponseStatus(t, testCase, http.StatusOK, status)
					requireResponseFields(t, testCase, []string{"index"}, body)
				},
			},*/
		{
			"Sending incorrect json field",
			sh,
			fmt.Sprintf(
				`{
						  "KVs": [
							{
								"keyX": "%s",
								"value": "%s"
							}
						  ]
						}`,
				validKey,
				validValue,
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "illegal arguments"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Sending plain text instead of base64 encoded",
			sh,
			`{
			"KVs": [
					{
						"key": "key",
						"value": "val"
					}
				]
				}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected :=
					map[string]interface{}{"error": "illegal base64 data at input byte 0"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Missing key field",
			sh,
			`{
				  "KVs": [{}]
				}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "illegal arguments"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"AnnotateContext error",
			NewSetHandler(mux, client, newTestRuntimeWithAnnotateContextErr(), json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "annotate context error"}, body)
			},
		},
		{
			"Set error",
			NewSetHandler(mux, immugwclient.NewMockClient(&clienttest.ImmuClientMock{ImmuClient: icd, SetAllF: setAllWErr}, opts), rt, json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "set error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewSetHandler(mux, client, rt, newTestJSONWithMarshalErr()),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "JSON marshal error"}, body)
			},
		},
	}
}
