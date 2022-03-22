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

	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/clienttest"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func testVerifiedZaddHandler(t *testing.T, mux *runtime.ServeMux, ic immuclient.ImmuClient) {
	prefixPattern := "VerifiedZaddHandler - Test case: %s"
	method := "POST"
	path := "/db/verified/zadd"
	for _, tc := range verifiedZaddHandlerTestCases(mux, ic) {
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			tc.verifiedZaddHandler.VerifiedZadd(res, req, nil)
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

type verifiedZaddHandlerTestCase struct {
	name                string
	verifiedZaddHandler VerifiedZaddHandler
	payload             string
	testFunc            func(*testing.T, string, int, map[string]interface{})
}

func verifiedZaddHandlerTestCases(mux *runtime.ServeMux, ic immuclient.ImmuClient) []verifiedZaddHandlerTestCase {
	rt := newDefaultRuntime()
	json := json.DefaultJSON()
	szh := NewVerifiedZaddHandler(mux, ic, rt, json)
	icd, _ := client.NewImmuClient(client.DefaultOptions())
	verifiedZaddWErr := func(context.Context, []byte, float64, []byte, uint64) (*schema.TxHeader, error) {
		return nil, errors.New("verifiedZadd error")
	}

	validSet := base64.StdEncoding.EncodeToString([]byte("verifiedZaddSet1"))
	validKey := base64.StdEncoding.EncodeToString([]byte("setKey1"))
	validPayload := fmt.Sprintf(
		`{
				  "zAddRequest": {
					"set": "%s",
					"score": %f,
					"key": "%s"
					}
				}`,
		validSet,
		1.0,
		validKey,
	)

	return []verifiedZaddHandlerTestCase{
		{
			"Sending correct request",
			szh,
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
			},
		},
		{
			"Sending request with non-existent key",
			szh,
			fmt.Sprintf(
				`{
				  "zAddRequest": {
					"set": "%s",
					"score": %f,
					"key": "%s"
					}
				}`,
				validSet,
				1.0,
				base64.StdEncoding.EncodeToString([]byte("verifiedZaddUnknownKey")),
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusNotFound, status)
			},
		},
		{
			"Sending request with incorrect JSON field",
			szh,
			fmt.Sprintf(
				`{
				  "zAddRequestsss": {
					"set": "%s",
					"score": %f,
					"key": "%s"
					}
				}`,
				validSet,
				1.0,
				validKey,
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "incorrect JSON payload"}, body)
			},
		},
		{
			"Missing key field",
			szh,
			fmt.Sprintf(
				`{
				  "zAddRequest": {
					"set": "%s",
					"score": %f
					}
				}`,
				validSet,
				1.0,
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "illegal arguments"}, body)
			},
		},
		{
			"Send plain text instead of base64 encoded",
			szh,
			fmt.Sprintf(
				`{
				  "zAddRequest": {
					"set": "%s",
					"score": %f,
					"key": "myFirstKey"
					}
				}`,
				validSet,
				1.0,
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "illegal base64 data at input byte 8"}, body)
			},
		},
		{
			"AnnotateContext error",
			NewVerifiedZaddHandler(mux, ic, newTestRuntimeWithAnnotateContextErr(), json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "annotate context error"}, body)
			},
		},
		{
			"VerifiedZadd error",
			NewVerifiedZaddHandler(mux, &clienttest.ImmuClientMock{ImmuClient: icd, VerifiedZAddF: verifiedZaddWErr}, rt, json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "verifiedZadd error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewVerifiedZaddHandler(mux, ic, rt, newTestJSONWithMarshalErr()),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "JSON marshal error"}, body)
			},
		},
	}
}
