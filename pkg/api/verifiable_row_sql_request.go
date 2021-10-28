package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"strconv"
)

type VerifiableRowSQLRequest struct {
	Row      *schema.Row
	Table    string
	PkValues []*schema.SQLValue
}

func (r *VerifiableRowSQLRequest) UnmarshalJSON(data []byte) (err error) {
	var obj map[string]*json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	var pkValues []*schema.SQLValue
	if v, ok := obj["pkValues"]; ok {
		pkValues, err = unmarshalSqlValues(*v)
		if err != nil {
			return err
		}
	}
	r.PkValues = pkValues

	var table string
	if err := json.Unmarshal(*obj["table"], &table); err != nil {
		return err
	}
	r.Table = table

	var values []*schema.SQLValue
	if v, ok := obj["values"]; ok {
		values, err = unmarshalSqlValues(*v)
		if err != nil {
			return err
		}
	}

	var row *schema.Row
	if err := json.Unmarshal(*obj["row"], &row); err != nil {
		return err
	}
	r.Row = row
	r.Row.Values = values
	return nil
}

func unmarshalSqlValues(data []byte) ([]*schema.SQLValue, error) {
	var rawValues []map[string]interface{}
	if err := json.Unmarshal(data, &rawValues); err != nil {
		return nil, err
	}

	var values []*schema.SQLValue
	for _, e := range rawValues {
		sqlVal := &schema.SQLValue{}
		f := func(key string) bool { _, ok := e[key]; return ok }

		switch {
		case f("s"):
			if v, ok := e["s"].(string); ok {
				sqlVal.Value = &schema.SQLValue_S{S: v}
			} else {
				return nil, StatusErrMalformedPayload
			}
		case f("b"):
			if v, ok := e["b"].(bool); ok {
				sqlVal.Value = &schema.SQLValue_B{B: v}
			} else {
				return nil, StatusErrMalformedPayload
			}
		case f("n"):
			if v, ok := e["n"].(string); ok {
				dv, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				sqlVal.Value = &schema.SQLValue_N{N: int64(dv)}
			} else {
				return nil, StatusErrMalformedPayload
			}
		case f("bs"):
			if v, ok := e["bs"].(string); ok {
				dv, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					return nil, err
				}
				sqlVal.Value = &schema.SQLValue_Bs{Bs: dv}
			} else {
				return nil, StatusErrMalformedPayload
			}
		case f("null"):
			sqlVal.Value = &schema.SQLValue_Null{}
		default:
			return nil, StatusErrMalformedPayload
		}
		values = append(values, sqlVal)
	}
	return values, nil
}
