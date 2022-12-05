/*
Copyright 2022 CodeNotary, Inc. All rights reserved.

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

package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/stretchr/testify/assert"
)

const (
	validVersion = 1
	anything     = 0
)

func segments(path string) (components []string, verb string) {
	if path == "" {
		return nil, ""
	}
	components = strings.Split(path, "/")
	l := len(components)
	c := components[l-1]
	if idx := strings.LastIndex(c, ":"); idx >= 0 {
		components[l-1], verb = c[:idx], c[idx+1:]
	}
	return components, verb
}

func TestBinding(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
		path string
		verb string
		want map[string]string
	}{
		// Pattern_ImmuService_CreateDatabaseV2_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"db", "databaseName", "create", "v2"}, "", runtime.AssumeColonVerbOpt(true)))
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpPush), 0,
				int(utilities.OpCapture), 1,
				int(utilities.OpLitPush), 2,
				int(utilities.OpLitPush), 3,
			},
			pool: []string{"db", "dbname", "create", "v2"},
			path: "db/testdb/create/v2",
			want: map[string]string{
				"dbname": "testdb",
			},
		},

		// Pattern_ImmuService_Count_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"db", "databaseName", "count", "prefix"}, "", runtime.AssumeColonVerbOpt(true)))
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpPush), 0,
				int(utilities.OpCapture), 1,
				int(utilities.OpLitPush), 2,
				int(utilities.OpLitPush), 3,
				int(utilities.OpPush), 0,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 3,
			},
			pool: []string{"db", "dbname", "count", "prefix"},
			path: "db/testdb/count/prefix/1",
			want: map[string]string{
				"prefix": "1",
				"dbname": "testdb",
			},
		},

		// Pattern_ImmuService_Get_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"db", "databaseName", "get", "key"}, "", runtime.AssumeColonVerbOpt(true)))
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpPush), 0,
				int(utilities.OpCapture), 1,
				int(utilities.OpLitPush), 2,
				int(utilities.OpLitPush), 3,
				int(utilities.OpPush), 0,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 3,
			},
			pool: []string{"db", "dbname", "get", "key"},
			path: "db/testdb/get/key/1",
			want: map[string]string{
				"key":    "1",
				"dbname": "testdb",
			},
		},

		// Pattern_ImmuService_TxById_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 1}, []string{"db", "databaseName", "tx"}, "", runtime.AssumeColonVerbOpt(true)))
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpPush), 0,
				int(utilities.OpCapture), 1,
				int(utilities.OpLitPush), 2,
				int(utilities.OpPush), 0,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 2,
			},
			pool: []string{"db", "dbname", "tx"},
			path: "db/testdb/tx/1",
			want: map[string]string{
				"tx":     "1",
				"dbname": "testdb",
			},
		},

		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpPush), 0,
				int(utilities.OpCapture), 1,
				int(utilities.OpLitPush), 2,
				int(utilities.OpLitPush), 3,
			},
			pool: []string{"db", "dbname", "verified", "set"},
			path: "db/my-bucket/verified/set",
			want: map[string]string{
				"dbname": "my-bucket",
			},
		},
	} {
		pat, err := runtime.NewPattern(validVersion, spec.ops, spec.pool, spec.verb)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, spec.verb, err)
			continue
		}

		components, verb := segments(spec.path)
		got, err := pat.Match(components, verb)
		if err != nil {
			t.Errorf("pat.Match(%q) failed with %v; want success; pattern = (%v, %q)", spec.path, err, spec.ops, spec.pool)
		}
		if !reflect.DeepEqual(got, spec.want) {
			t.Errorf("pat.Match(%q) = %q; want %q; pattern = (%v, %q)", spec.path, got, spec.want, spec.ops, spec.pool)
		}
	}
}

func TestCommonPatterns(t *testing.T) {
	for _, spec := range []struct {
		pattern runtime.Pattern
		path    string
		verb    string
		want    map[string]string
		wantErr bool
	}{
		{
			pattern: Pattern_ImmuService_Set_0,
			path:    "db/testdb/set",
			want: map[string]string{
				"databaseName": "testdb",
			},
		},
		{
			pattern: Pattern_ImmuService_Get_0,
			path:    "db/testdb/get/key/key1",
			want: map[string]string{
				"databaseName": "testdb",
				"key":          "key1",
			},
		},
		{
			pattern: Pattern_ImmuService_CreateDatabaseV2_0,
			path:    "db/testdb/create/v2",
			want: map[string]string{
				"databaseName": "testdb",
			},
		},
		{
			pattern: Pattern_ImmuService_TxById_0,
			path:    "db/testdb/tx/1",
			want: map[string]string{
				"databaseName": "testdb",
				"tx":           "1",
			},
		},
		{
			pattern: Pattern_ImmuService_Count_0,
			path:    "db/testdb/count/prefix/1",
			want: map[string]string{
				"databaseName": "testdb",
				"prefix":       "1",
			},
		},
		{
			pattern: Pattern_ImmuService_Count_0,
			path:    "db/testdb/count/prefix/1/abc",
			wantErr: true,
		},
	} {
		pat := spec.pattern
		components, verb := segments(spec.path)
		got, err := pat.Match(components, verb)
		if spec.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoErrorf(t, err, "pat.Match(%q) failed with %v; want success;", spec.path, err)
			if !reflect.DeepEqual(got, spec.want) {
				t.Errorf("pat.Match(%q) = %q; want %q;", spec.path, got, spec.want)
			}
		}
	}
}
