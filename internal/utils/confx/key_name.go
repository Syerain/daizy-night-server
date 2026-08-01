package confx

import (
	"reflect"
	"strings"
)

// priority: mapstructure > yaml > json > toml > StructFieldName
func resolveKeyName(field reflect.StructField) string {
	priority := []string{"mapstructure", "yaml", "json", "toml"}
	for _, key := range priority {
		tag := field.Tag.Get(key)
		if tag == "" {
			continue
		}

		// Example: `mapstructure:"fieldName,omitempty"`
		name := strings.SplitN(tag, ",", 2)[0] // extract 'fieldName'
		if name != "" && name != "-" {
			return name
		}
	}
	// fallback to StructFieldName
	return field.Name
}
