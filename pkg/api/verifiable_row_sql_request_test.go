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

package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifiableRowSQLRequest_UnmarshalJSON(t *testing.T) {
	jsonPayload := []byte(`{
    "row": {
            "columns": [
                "(testdb1.mytable22.id)",
                "(testdb1.mytable22.amount)",
                "(testdb1.mytable22.total)",
                "(testdb1.mytable22.title)",
                "(testdb1.mytable22.content)",
                "(testdb1.mytable22.ispresent)"
            ]
            },
            "values": [
                {
                    "n": "1"
                },
                {
                    "n": "1000"
                },
                {
                    "null": null
                },
                {
                    "s": "title 1"
                },
                {
                    "bs": "YmxvYiBjb250ZW50"
                },
                {
                    "b": true
                }
            ],
        "table": "myTable22",
        "pkValues": [
            {
                "n": "1"
            }
        ]
      }`)

	var r VerifiableRowSQLRequest
	err := json.Unmarshal(jsonPayload, &r)
	require.NoError(t, err)
}
