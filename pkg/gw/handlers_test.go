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
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/tokenservice"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immudb/pkg/server/servertest"
	"github.com/codenotary/immugw/pkg/api"
	immugwclient "github.com/codenotary/immugw/pkg/client"
	"github.com/codenotary/immugw/pkg/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	defaultTestParams = map[string]string{"databaseName": "defaultdb"}
)

func TestGw(t *testing.T) {
	options := server.DefaultOptions().WithAuth(false).WithDir(t.TempDir())
	bs := servertest.NewBufconnServer(options)

	bs.Start()
	defer bs.Stop()

	defer func() {
		matches, _ := filepath.Glob("state-*")
		os.RemoveAll(options.Dir)
		for _, m := range matches {
			os.RemoveAll(m)
		}
	}()

	opts := client.DefaultOptions().WithDialOptions([]grpc.DialOption{grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure()}).WithAuth(false).WithDir(t.TempDir())
	client := immugwclient.New(opts)
	client.Add("defaultdb")

	mux := runtime.NewServeMux(runtime.WithProtoErrorHandler(runtime.DefaultHTTPError))

	testSafeSetHandler(t, mux, client, opts)
	testSetHandler(t, mux, client, opts)
	testSafeGetHandler(t, mux, client, opts)
	testHistoryHandler(t, mux, client, opts)
	testVerifiedSetReferenceHandler(t, mux, client, opts)
	testVerifiedZaddHandler(t, mux, client, opts)
}

func TestAuthGw(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithDir(t.TempDir())
	bs := servertest.NewBufconnServer(options)

	bs.Start()
	defer bs.Stop()

	defer func() {
		matches, _ := filepath.Glob("state-*")
		os.RemoveAll(options.Dir)
		for _, m := range matches {
			os.RemoveAll(m)
		}
	}()

	opts := immuclient.DefaultOptions().WithDialOptions([]grpc.DialOption{grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure()}).WithAuth(false).WithDir(t.TempDir())
	client := immugwclient.New(opts)
	immuClient, _ := client.Add("defaultdb")

	mux := runtime.NewServeMux(runtime.WithProtoErrorHandler(api.DefaultGWErrorHandler))

	ctx := context.TODO()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure(),
	}
	pr := &PasswordReader{
		Pass: []string{"immudb"},
	}
	ts := tokenservice.NewInmemoryTokenService()
	cliopt := immuclient.DefaultOptions().WithDialOptions(dialOptions).WithPasswordReader(pr)
	cliopt.PasswordReader = pr
	cliopt.DialOptions = dialOptions

	cli, _ := immuclient.NewImmuClient(cliopt)
	cli.WithTokenService(ts)

	lresp, err := cli.Login(ctx, []byte("immudb"), []byte("immudb"))
	if err != nil {
		t.Fatal(err)
	}

	md := metadata.Pairs("authorization", lresp.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	require.NoError(t, immuClient.HealthCheck(ctx))

	testUseDatabaseHandler(t, ctx, mux, client, opts)
}

func testHandler(
	t *testing.T,
	name string,
	method string,
	path string,
	body string,
	handlerFunc func(http.ResponseWriter, *http.Request),
	testFunc func(*testing.T, string, int, map[string]interface{}),
) error {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(w, req)
	testCase := fmt.Sprintf("%s - %s %s %s - ", name, method, path, body)
	respBytes := w.Body.Bytes()
	var respBody map[string]interface{}
	if err := json.DefaultJSON().Unmarshal(respBytes, &respBody); err != nil {
		return fmt.Errorf(
			"%s - error unmarshaling JSON from response %s", testCase, respBytes)
	}
	testFunc(t, testCase, w.Code, respBody)
	return nil
}

func newTestJSONWithMarshalErr() json.JSON {
	return &json.StandardJSON{
		MarshalF: func(v interface{}) ([]byte, error) {
			return nil, errors.New("JSON marshal error")
		},
	}
}

func requireResponseStatus(
	t *testing.T,
	testCase string,
	expected int,
	actual int,
) {
	require.Equal(
		t,
		expected,
		actual,
		"%sexpected HTTP status %d, actual %d", testCase, expected, actual)
}

func getMissingResponseFieldPattern(testCase string) string {
	return testCase + "\"%s\" field is missing from response %v"
}

func requireResponseFields(
	t *testing.T,
	testCase string,
	fields []string,
	body map[string]interface{},
) {
	missingPattern := getMissingResponseFieldPattern(testCase)
	for _, field := range fields {
		_, ok := body[field]
		require.True(t, ok, missingPattern, field, body)
	}
}

func requireResponseFieldsTrue(
	t *testing.T,
	testCase string,
	fields []string,
	body map[string]interface{},
) {
	missingPattern := getMissingResponseFieldPattern(testCase)
	isFalsePattern := testCase + "\"%s\" field is false in response %v"
	for _, field := range fields {
		fieldValue, ok := body[field]
		require.True(t, ok, missingPattern, field, body)
		require.True(t, fieldValue.(bool), isFalsePattern, field, body)
	}
}

func requireResponseFieldsEqual(
	t *testing.T,
	testCase string,
	fields map[string]interface{},
	body map[string]interface{},
) {
	missingPattern := getMissingResponseFieldPattern(testCase)
	notEqPattern := testCase +
		"expected response %v to have field \"%s\" = \"%v\", but actual field value is \"%v\""
	for field, expected := range fields {
		fieldValue, ok := body[field]
		require.True(t, ok, missingPattern, field, body)
		require.Equal(
			t, expected, fieldValue, notEqPattern, body, field, expected, fieldValue)
	}
}

type PasswordReader struct {
	Pass       []string
	callNumber int
}

func (pr *PasswordReader) Read(msg string) ([]byte, error) {
	if len(pr.Pass) <= pr.callNumber {
		log.Fatal("Application requested the password more times than number of passwords supplied")
	}
	pass := []byte(pr.Pass[pr.callNumber])
	pr.callNumber++
	return pass, nil
}
