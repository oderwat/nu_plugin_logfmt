package main

import (
	"fmt"
	"github.com/ainvaltin/nu-plugin"
)

func asValue(v any) (_ nu.Value, err error) {
	switch in := v.(type) {
	case uint64, float64, bool, string, []byte:
		return nu.Value{Value: in}, nil
	case []any:
		lst := make([]nu.Value, len(in))
		for i := 0; i < len(in); i++ {
			if lst[i], err = asValue(in[i]); err != nil {
				return nu.Value{}, err
			}
		}
		return nu.Value{Value: lst}, nil
	case map[string]any:
		rec := nu.Record{}
		for k, v := range in {
			if rec[k], err = asValue(v); err != nil {
				return nu.Value{}, err
			}
		}
		return nu.Value{Value: rec}, nil
	default:
		return nu.Value{}, fmt.Errorf("unsupported value type %T", in)
	}
}
