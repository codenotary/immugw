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
	"github.com/codenotary/immudb/pkg/api/schema"
	"net/http"
	"testing"

	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/clienttest"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func testSafeSetHandler(t *testing.T, mux *runtime.ServeMux, ic immuclient.ImmuClient) {
	prefixPattern := "SafeSetHandler - Test case: %s"
	method := "POST"
	path := "/db/verifiable/set"
	for _, tc := range safeSetHandlerTestCases(mux, ic) {
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			tc.verifiedSetHandler.VerifiedSet(res, req, nil)
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

type safeSetHandlerTestCase struct {
	name               string
	verifiedSetHandler VerifiedSetHandler
	payload            string
	testFunc           func(*testing.T, string, int, map[string]interface{})
}

func safeSetHandlerTestCases(mux *runtime.ServeMux, ic immuclient.ImmuClient) []safeSetHandlerTestCase {
	rt := newDefaultRuntime()
	json := json.DefaultJSON()
	ssh := NewVerifiedSetHandler(mux, ic, rt, json)
	icd := client.DefaultClient()
	safeSetWErr := func(context.Context, []byte, []byte) (*schema.TxMetadata, error) {
		return nil, errors.New("safeset error")
	}
	validKey := base64.StdEncoding.EncodeToString([]byte("safeSetKey1"))
	validValue := base64.StdEncoding.EncodeToString([]byte("safeSetValue1"))
	validPayload := fmt.Sprintf(
		"{\n  \"setRequest\": {\n    \"KVs\": [\n      {\n         \"key\": \"%s\",\n\t     \"value\": \"%s\"\n      }\n    ]\n  }\n}",
		validKey,
		validValue,
	)

	return []safeSetHandlerTestCase{
		{
			"Sending correct request",
			ssh,
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
			},
		},
		{
			"Missing value field",
			ssh,
			fmt.Sprintf("{\n  \"setRequest\": {\n    \"KVs\": [\n      {\n         \"key\": \"%s\"\n      }\n    ]\n  }\n}",
				validKey),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
			},
		},
		{
			"Sending incorrect json field",
			ssh,
			fmt.Sprintf(
				"{\"data\": {\"key\": \"%s\", \"value\": \"%s\"}}",
				validKey,
				validValue,
			),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "incorrect JSON payload"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Sending plain text instead of base64 encoded",
			ssh,
			`{"setRequest": {"KVs": [{"key": "key","value": "val"}]}}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "illegal base64 data at input byte 0"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Missing key field",
			ssh,
			`{"setRequest": {"KVs": []}}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "verifiedSet accept at least one key value pair"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"AnnotateContext error",
			NewVerifiedSetHandler(mux, ic, newTestRuntimeWithAnnotateContextErr(), json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "annotate context error"}, body)
			},
		},
		{
			"SafeSet error",
			NewVerifiedSetHandler(mux, &clienttest.ImmuClientMock{ImmuClient: icd, VerifiedSetF: safeSetWErr}, rt, json),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "safeset error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewVerifiedSetHandler(mux, ic, rt, newTestJSONWithMarshalErr()),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "JSON marshal error"}, body)
			},
		},
	}
}
