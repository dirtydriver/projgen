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
// It supports both simple key names and nested YAML paths using dot notation (e.g., 'project.name').
// It returns an error listing any missing keys, or nil if all keys are present.
func CheckMissingKeys(m map[string]interface{}, list []string) error {
	var missing []string
	for _, key := range list {
		if !hasNestedKey(m, strings.Split(key, ".")) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return errors.New("missing keys: " + strings.Join(missing, ", "))
	}
	return nil
}

// hasNestedKey checks if a nested key exists in the map using path segments.
// For example, for path ["project", "name"] it checks m["project"]["name"].
func hasNestedKey(m map[string]interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}

	val, exists := m[path[0]]
	if !exists {
		return false
	}

	if len(path) == 1 {
		return true
	}

	// For nested keys, check if the value is a map and recurse
	if nextMap, ok := val.(map[string]interface{}); ok {
		return hasNestedKey(nextMap, path[1:])
	}

	// Handle yaml.Node type which might come from yaml parsing
	if nextMap, ok := val.(map[interface{}]interface{}); ok {
		// Convert to map[string]interface{} for consistent handling
		strMap := make(map[string]interface{})
		for k, v := range nextMap {
			if ks, ok := k.(string); ok {
				strMap[ks] = v
			}
		}
		return hasNestedKey(strMap, path[1:])
	}

	return false
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
//
// If intermediate nested maps don't exist, they will be created automatically.
func ApplyOverrides(m map[string]interface{}, overrides []string) {
	for _, value := range overrides {

		if kv := strings.SplitN(value, "=", 2); len(kv) == 2 {
			keyPath := strings.Split(kv[0], ".")
			setNestedValues(m, keyPath, kv[1])
		}
	}
}
