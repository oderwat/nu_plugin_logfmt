package logfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Encode(m map[string]any) string {
	var parts []string
	encodeMap(m, "", &parts)
	return strings.Join(parts, " ")
}

func encodeMap(m map[string]any, prefix string, parts *[]string) {
	for k, v := range m {
		// Construct full key name with potential prefix
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		// Handle slices
		switch val := v.(type) {
		case []any:
			encodeSlice(val, fullKey, parts)
		case map[string]any:
			encodeMap(val, fullKey, parts)
		default:
			// Encode the value, escaping as necessary
			encodedValue := encodeValue(v)
			*parts = append(*parts, fullKey+"="+encodedValue)
		}
	}
}

func encodeSlice(slice []any, prefix string, parts *[]string) {
	for i, v := range slice {
		// Create key with index as string
		indexKey := prefix + ".[" + strconv.Itoa(i) + "]"

		// Handle nested structures
		switch val := v.(type) {
		case map[string]any:
			encodeMap(val, indexKey, parts)
		case []any:
			encodeSlice(val, indexKey, parts)
		default:
			// Encode the value, escaping as necessary
			encodedValue := encodeValue(v)
			*parts = append(*parts, indexKey+"="+encodedValue)
		}
	}
}

func encodeValue(v any) string {
	if v == nil {
		return "null"
	}

	// Convert to string, escaping special characters
	s := toString(v)

	// Check if value needs quoting (contains spaces or special characters)
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
	// Quote if contains spaces, quotes, backslashes, or control characters
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
