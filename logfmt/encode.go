package logfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Encode(v any) string {
	var parts []string
	encodeValue("", v, &parts)
	return strings.Join(parts, " ")
}

func encodeValue(prefix string, v any, parts *[]string) {
	switch val := v.(type) {
	case map[string]any:
		encodeMap(val, prefix, parts)
	case []any:
		encodeSlice(val, prefix, parts)
	default:
		if prefix != "" {
			*parts = append(*parts, prefix+"="+encodeScalar(v))
		}
	}
}

func encodeMap(m map[string]any, prefix string, parts *[]string) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		encodeValue(fullKey, v, parts)
	}
}

func encodeSlice(slice []any, prefix string, parts *[]string) {
	for i, v := range slice {
		indexKey := prefix
		if prefix == "" {
			indexKey = "[" + strconv.Itoa(i) + "]"
		} else {
			indexKey = prefix + ".[" + strconv.Itoa(i) + "]"
		}
		encodeValue(indexKey, v, parts)
	}
}

func encodeScalar(v any) string {
	if v == nil {
		return "null"
	}

	s := toString(v)
	if needsQuoting(s) {
		return `"` + escapeString(s) + `"`
	}
	return s
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func needsQuoting(s string) bool {
	return strings.ContainsAny(s, " \t\n\r\"\\")
}

func escapeString(s string) string {
	var escaped strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			escaped.WriteString(`\"`)
		case '\\':
			escaped.WriteString(`\\`)
		case '\n':
			escaped.WriteString(`\n`)
		case '\r':
			escaped.WriteString(`\r`)
		case '\t':
			escaped.WriteString(`\t`)
		default:
			escaped.WriteRune(r)
		}
	}
	return escaped.String()
}
