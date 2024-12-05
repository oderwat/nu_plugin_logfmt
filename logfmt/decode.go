package logfmt

import (
	"strconv"
	"strings"
	"unicode"
)

type DecodeOptions struct {
	ParseTypes bool
}

func Decode(input string, opts ...DecodeOptions) any {
	parseTypes := true
	if len(opts) > 0 {
		parseTypes = opts[0].ParseTypes
	}

	var result any = make(map[string]any)
	isArrayRoot := true

	input = strings.TrimSpace(input)
	if input == "" {
		return result
	}

	rootMap := make(map[string]any)
	parseKeyValuePairs(input, rootMap, parseTypes)

	for key := range rootMap {
		parts := strings.Split(key, ".")
		if !isArrayIndex(parts[0]) {
			isArrayRoot = false
			break
		}
	}

	if isArrayRoot {
		return mapToSlice(rootMap)
	}
	return convertMapsToSlices(rootMap)
}

func parseKeyValuePairs(input string, result map[string]any, parseTypes bool) {
	var key strings.Builder
	var value strings.Builder
	var isInKey = true
	var isInQuotedValue = false
	var isEscaped = false

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		char := runes[i]

		switch {
		case isEscaped:
			switch char {
			case 'n':
				value.WriteRune('\n')
			case 'r':
				value.WriteRune('\r')
			case 't':
				value.WriteRune('\t')
			default:
				value.WriteRune(char)
			}
			isEscaped = false

		case char == '\\':
			if isInQuotedValue {
				isEscaped = true
			} else if !isInKey {
				value.WriteRune(char)
			}

		case char == '"':
			if isInQuotedValue {
				isInQuotedValue = false
			} else if isInKey || value.Len() == 0 {
				isInQuotedValue = true
			} else {
				value.WriteRune(char)
			}

		case char == '=' && !isInQuotedValue:
			isInKey = false

		case (unicode.IsSpace(char) && !isInQuotedValue && !isInKey):
			if key.Len() > 0 {
				var val any
				if parseTypes {
					val = parseValue(value.String())
				} else {
					val = value.String()
				}
				setNestedValue(result, strings.TrimSpace(key.String()), val)
				key.Reset()
				value.Reset()
				isInKey = true
			}

		default:
			if isInKey {
				key.WriteRune(char)
			} else {
				value.WriteRune(char)
			}
		}
	}

	if key.Len() > 0 {
		var val any
		if parseTypes {
			val = parseValue(value.String())
		} else {
			val = value.String()
		}
		setNestedValue(result, strings.TrimSpace(key.String()), val)
	}
}

func setNestedValue(m map[string]any, key string, value any) {
	parts := parsePath(key)
	current := m

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]any)
		}
		current = current[part].(map[string]any)
	}

	current[parts[len(parts)-1]] = value
}

func parsePath(path string) []string {
	var parts []string
	var current strings.Builder

	runes := []rune(path)
	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if char == '.' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}

		if char == '[' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			// Find the closing bracket
			for i < len(runes) && runes[i] != ']' {
				current.WriteRune(runes[i])
				i++
			}
			current.WriteRune(']')
			parts = append(parts, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(char)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func convertMapsToSlices(data map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range data {
		if nestedMap, ok := value.(map[string]any); ok {
			if isArrayMap(nestedMap) {
				result[key] = mapToSlice(nestedMap)
			} else {
				result[key] = convertMapsToSlices(nestedMap)
			}
		} else {
			result[key] = value
		}
	}

	return result
}

func isArrayMap(m map[string]any) bool {
	if len(m) == 0 {
		return false
	}

	for key := range m {
		if !isArrayIndex(key) {
			return false
		}
	}
	return true
}

func isArrayIndex(key string) bool {
	if !strings.HasPrefix(key, "[") || !strings.HasSuffix(key, "]") {
		return false
	}

	numStr := key[1 : len(key)-1]
	_, err := strconv.Atoi(numStr)
	return err == nil
}

func mapToSlice(m map[string]any) []any {
	if len(m) == 0 {
		return nil
	}

	maxIndex := -1
	for key := range m {
		index, _ := strconv.Atoi(key[1 : len(key)-1])
		if index > maxIndex {
			maxIndex = index
		}
	}

	result := make([]any, maxIndex+1)
	for key, value := range m {
		index, _ := strconv.Atoi(key[1 : len(key)-1])
		if nestedMap, ok := value.(map[string]any); ok {
			if isArrayMap(nestedMap) {
				result[index] = mapToSlice(nestedMap)
			} else {
				result[index] = convertMapsToSlices(nestedMap)
			}
		} else {
			result[index] = value
		}
	}

	return result
}

func parseValue(val string) any {
	// Remove quotes if present
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		return val[1 : len(val)-1]
	}

	if val == "null" {
		return nil
	}
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	return val
}
