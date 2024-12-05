package logfmt

import (
	"strings"
	"unicode"
)

func Decode(input string) map[string]any {
	result := make(map[string]any)

	input = strings.TrimSpace(input)

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
			// Handle escaped characters
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
			// Start of an escape sequence
			if isInQuotedValue {
				isEscaped = true
			} else if !isInKey {
				value.WriteRune(char)
			}

		case char == '"':
			// Toggle quoted value state
			if isInQuotedValue {
				isInQuotedValue = false
			} else if isInKey || value.Len() == 0 {
				isInQuotedValue = true
			} else {
				value.WriteRune(char)
			}

		case char == '=' && !isInQuotedValue:
			// Transition from key to value
			isInKey = false

		case (unicode.IsSpace(char) && !isInQuotedValue && !isInKey):
			// End of a key-value pair
			setNestedValue(result, strings.TrimSpace(key.String()), value.String())
			key.Reset()
			value.Reset()
			isInKey = true

		default:
			// Accumulate characters
			if isInKey {
				key.WriteRune(char)
			} else {
				value.WriteRune(char)
			}
		}
	}

	// Add the last key-value pair if exists
	if key.Len() > 0 {
		setNestedValue(result, strings.TrimSpace(key.String()), value.String())
	}

	return result
}

// setNestedValue sets a value in a nested map structure
func setNestedValue(m map[string]any, key string, value string) {
	parts := strings.Split(key, ".")

	// Navigate or create nested maps
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := m[part]; !exists {
			m[part] = make(map[string]any)
		}

		// Type assert and update m to the nested map
		switch nested := m[part].(type) {
		case map[string]any:
			m = nested
		default:
			// If not a map, replace with a new map
			newNested := make(map[string]any)
			m[part] = newNested
			m = newNested
		}
	}

	// Set the final value
	m[parts[len(parts)-1]] = value
}
