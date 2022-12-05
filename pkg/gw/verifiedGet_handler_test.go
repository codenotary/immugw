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

func testSafeGetHandler(t *testing.T, mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) {
	prefixPattern := "SafeGetHandler - Test case: %s"
	method := "POST"
	path := "/db/defaultdb/verified/get"
	for _, tc := range safeGetHandlerTestCases(mux, client, opts) {
		handlerFunc := func(res http.ResponseWriter, req *http.Request) {
			tc.safeGetHandler.VerifiedGet(res, req, defaultTestParams)
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

type safeGetHandlerTestCase struct {
	name           string
	safeGetHandler VerifiedGetHandler
	payload        string
	testFunc       func(*testing.T, string, int, map[string]interface{})
}

func safeGetHandlerTestCases(mux *runtime.ServeMux, client immugwclient.Client, opts *immuclient.Options) []safeGetHandlerTestCase {
	rt := newDefaultRuntime()
	defaultJSON := json.DefaultJSON()
	sgh := NewVerifiedGetHandler(mux, client, rt, defaultJSON)
	//icd := client.DefaultClient()
	verifiedGetWErr := func(context.Context, []byte, uint64) (*schema.Entry, error) {
		return nil, errors.New("verified get error")
	}
	validKey := base64.StdEncoding.EncodeToString([]byte("setKey1"))
	validPayload := fmt.Sprintf(`{
  "keyRequest": {
    "key": "%s"
  }
}`, validKey)

	return []safeGetHandlerTestCase{
		{
			"Sending correct request",
			sgh,
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusOK, status)
				requireResponseFields(
					t, testCase, []string{"tx", "key", "value"}, body)
			},
		},
		{
			"Sending incorrect json field",
			sgh,
			fmt.Sprintf(`{
  "keyRequest": {
    "keyX": "%s"
  }
}`, validKey),
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "illegal arguments"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Sending plain text instead of base64 encoded",
			sgh,
			`{
  "keyRequest": {
    "key": "setKey1"
  }
}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected :=
					map[string]interface{}{"error": "illegal base64 data at input byte 4"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"Missing key field",
			sgh,
			`{
  "keyRequest": {}
}`,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusBadRequest, status)
				expected := map[string]interface{}{"error": "illegal arguments"}
				requireResponseFieldsEqual(t, testCase, expected, body)
			},
		},
		{
			"AnnotateContext error",
			NewVerifiedGetHandler(mux, client, newTestRuntimeWithAnnotateContextErr(), defaultJSON),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "annotate context error"}, body)
			},
		},
		{
			"VerifiedGet error",
			NewVerifiedGetHandler(mux, immugwclient.NewMockClient(&clienttest.ImmuClientMock{VerifiedGetAtF: verifiedGetWErr}, opts), rt, defaultJSON),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "verified get error"}, body)
			},
		},
		{
			"JSON marshal error",
			NewVerifiedGetHandler(mux, client, rt, newTestJSONWithMarshalErr()),
			validPayload,
			func(t *testing.T, testCase string, status int, body map[string]interface{}) {
				requireResponseStatus(t, testCase, http.StatusInternalServerError, status)
				requireResponseFieldsEqual(
					t, testCase, map[string]interface{}{"error": "JSON marshal error"}, body)
			},
		},
	}
}
