package main

import (
	"context"
	"fmt"
	"github.com/ainvaltin/nu-plugin"
	"github.com/oderwat/nu_plugin_logfmt/logfmt"
)

func toLogFmt() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:                 "to logfmt",
			Category:             "Formats",
			SearchTerms:          []string{"logfmt", "slog", "logging"},
			InputOutputTypes:     [][]string{{"Any", "String"}},
			Usage:                `Convert Nushell Value to 'logfmt' format.`,
			AllowMissingExamples: true,
		},
		Examples: nu.Examples{
			{
				Description: `Convert an Nu record to a logmt string`,
				Example:     `{ "msg": "Hello World!", "Lang": { "Go": true, "Rust": false } } | to logfmt`,
				Result:      &nu.Value{Value: `msg="Hello World!" Lang.Go=true Lang.Rust=false`},
			},
		},
		OnRun: toLogFmtHandler,
	}
}

func toLogFmtHandler(ctx context.Context, call *nu.ExecCommand) error {
	switch in := call.Input.(type) {
	case nil:
		return nil
	case nu.Value:
		v, err := toPlistValue(in)
		if err != nil {
			return err
		}
		return call.ReturnValue(ctx, v)
	case <-chan nu.Value:
		out, err := call.ReturnListStream(ctx)
		if err != nil {
			return err
		}
		defer close(out)
		for v := range in {
			v, err := toPlistValue(v)
			if err != nil {
				return err
			}
			out <- v
		}
		return nil
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}

func toPlistValue(v nu.Value) (nu.Value, error) {
	var buf []byte
	var data map[string]any
	data = fromValue(v).(map[string]interface{})
	buf = []byte(logfmt.Encode(data))
	return nu.Value{Value: string(buf)}, nil
}

func fromValue(v nu.Value) any {
	switch vt := v.Value.(type) {
	case []nu.Value:
		lst := make([]any, len(vt))
		for i := 0; i < len(vt); i++ {
			lst[i] = fromValue(vt[i])
		}
		return lst
	case nu.Record:
		rec := map[string]any{}
		for k, v := range vt {
			rec[k] = fromValue(v)
		}
		return rec
	}
	return v.Value
}
