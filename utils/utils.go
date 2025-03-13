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
