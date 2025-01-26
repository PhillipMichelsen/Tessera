package helper

import (
	"fmt"
	"strconv"
	"time"
)

// FieldInfo holds information about a required field: its name and expected type.
type FieldInfo struct {
	Name string
	Type string
}

// ValidateAndExtract checks if required fields exist in the data map with the correct types,
// parses them as needed, and returns a map of extracted and validated fields.
func ValidateAndExtract(data map[string]interface{}, fields []FieldInfo) (map[string]interface{}, error) {
	validatedFields := make(map[string]interface{})

	for _, field := range fields {
		value, exists := data[field.Name]
		if !exists {
			return nil, fmt.Errorf("missing field: %s", field.Name)
		}

		// Parse and type-check based on expected type
		switch field.Type {
		case "string":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for field %s: expected string", field.Name)
			}
			validatedFields[field.Name] = strValue

		case "float64":
			// Allow parsing of strings as floats for flexibility
			switch v := value.(type) {
			case string:
				floatValue, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse float for field %s: %v", field.Name, err)
				}
				validatedFields[field.Name] = floatValue
			case float64:
				validatedFields[field.Name] = v
			default:
				return nil, fmt.Errorf("invalid type for field %s: expected float64", field.Name)
			}

		case "int":
			// Allow parsing of strings as integers for flexibility, int64 is used to handle large numbers
			switch v := value.(type) {
			case string:
				intValue, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse int for field %s: %v", field.Name, err)
				}
				validatedFields[field.Name] = intValue
			case int:
				validatedFields[field.Name] = int64(v)
			case float64:
				validatedFields[field.Name] = int64(v)
			default:
				return nil, fmt.Errorf("invalid type for field %s: expected int", field.Name)

			}

		case "bool":
			// Allow parsing of strings as booleans for flexibility
			switch v := value.(type) {
			case string:
				boolValue, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("failed to parse bool for field %s: %v", field.Name, err)
				}
				validatedFields[field.Name] = boolValue
			case bool:
				validatedFields[field.Name] = v
			default:
				return nil, fmt.Errorf("invalid type for field %s: expected bool", field.Name)
			}

		case "time":
			// Convert millisecond timestamps (float64) to time.Time
			switch v := value.(type) {
			case float64:
				validatedFields[field.Name] = time.UnixMilli(int64(v)).UTC()
			default:
				return nil, fmt.Errorf("invalid type for field %s: expected timestamp in milliseconds as float64", field.Name)
			}

		case "map[string]interface{}":
			mapValue, ok := value.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid type for field %s: expected map[string]interface{}", field.Name)
			}
			validatedFields[field.Name] = mapValue

		case "[]interface{}":
			sliceValue, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid type for field %s: expected []interface{}", field.Name)
			}
			validatedFields[field.Name] = sliceValue

		default:
			return nil, fmt.Errorf("unsupported type for field %s", field.Name)
		}
	}

	return validatedFields, nil
}
