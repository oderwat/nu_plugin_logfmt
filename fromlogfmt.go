package main

import (
	"context"
	"fmt"
	"io"

	"github.com/ainvaltin/nu-plugin"

	"github.com/oderwat/nu_plugin_logfmt/logfmt"
)

func fromLogFmt() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:                 "from logfmt",
			Category:             "Formats",
			SearchTerms:          []string{"logfmt", "slog", "logging"},
			InputOutputTypes:     [][]string{{"String", "Any"}},
			Usage:                `Convert from 'logfmt' format to Nushell Value.`,
			AllowMissingExamples: true,
		},
		Examples: nu.Examples{
			{
				Description: `Convert a logfmt string to a Nu record`,
				Example:     `'msg="Test message" level=info esc="» Say \"Hello\"' | from logfmt`,
				Result: &nu.Value{Value: nu.Record{
					"level": nu.Value{Value: "info"},
					"msg":   nu.Value{Value: "Test message"},
					"esc":   nu.Value{Value: "» Say \"Hello\""},
				}},
			},
		},
		OnRun: fromLogFmtHandler,
	}
}

func fromLogFmtHandler(ctx context.Context, call *nu.ExecCommand) error {
	switch in := call.Input.(type) {
	case nil:
		return nil
	case nu.Value:
		var buf []byte
		switch data := in.Value.(type) {
		case []byte:
			buf = data
		case string:
			buf = []byte(data)
		default:
			return fmt.Errorf("unsupported input value type %T", data)
		}
		fields := logfmt.Decode(string(buf))
		rv, err := asValue(fields)
		if err != nil {
			return fmt.Errorf("converting to Value: %w", err)
		}
		return call.ReturnValue(ctx, rv)
	case io.Reader:
		// decoder wants io.ReadSeeker so we need to read to buf.
		// could read just enough that the decoder can detect the
		// format and stream the rest?
		buf, err := io.ReadAll(in)
		if err != nil {
			return fmt.Errorf("reding input: %w", err)
		}
		fields := logfmt.Decode(string(buf))
		rv, err := asValue(fields)
		if err != nil {
			return fmt.Errorf("converting to Value: %w", err)
		}
		return call.ReturnValue(ctx, rv)
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}
