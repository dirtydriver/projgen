package utils

import (
	"errors"
	"strings"
)

// RemoveDuplicates returns a new slice containing unique elements from the input slice.
func RemoveDuplicates(list []string) []string {
	uniqCheck := make(map[string]bool)
	var uniqElements []string

	for _, element := range list {
		if uniqCheck[element] {
			continue
		}
		uniqElements = append(uniqElements, element)
		uniqCheck[element] = true
	}
	return uniqElements
}

// CheckMissingKeys verifies that all required keys in the list exist in the given map.
// It returns an error listing any missing keys, or nil if all keys are present.
func CheckMissingKeys(m map[string]interface{}, list []string) error {
	var missing []string
	for _, key := range list {
		if _, found := m[key]; !found {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return errors.New("missing keys: " + strings.Join(missing, ", "))
	}
	return nil
}

func setNestedValues(m map[string]interface{}, path []string, value interface{}) {
	for i := 0; i < len(path)-1; i++ {
		k := path[i]

		if _, exists := m[k]; !exists {
			m[k] = make(map[string]interface{})
		}

		if submap, exists := m[k].(map[string]interface{}); exists {
			m = submap
		} else {
			newMap := make(map[string]interface{})
			m[k] = newMap
			m = newMap
		}
	}
	m[path[len(path)-1]] = value
}

// ApplyOverrides applies a list of key-value overrides to a map using dot notation for nested keys.
// Each override should be in the format "key.subkey=value". For example:
//   - "name=John" sets m["name"] = "John"
//   - "config.port=8080" sets m["config"]["port"] = "8080"
// If intermediate nested maps don't exist, they will be created automatically.
func ApplyOverrides(m map[string]interface{}, overrides []string) {
	for _, value := range overrides {

		if kv := strings.SplitN(value, "=", 2); len(kv) == 2 {
			keyPath := strings.Split(kv[0], ".")
			setNestedValues(m, keyPath, kv[1])
		}
	}
}
