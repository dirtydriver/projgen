package utils

import (
	"reflect"
	"slices"
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

// ValidateParameters checks if two parameter lists contain the same elements (order-independent).
func ValidateParameters(ParamList, TemplateParamList *[]string) bool {

	slices.Sort(*ParamList)
	slices.Sort(*TemplateParamList)

	return reflect.DeepEqual(*ParamList, *TemplateParamList)
}

// CollectMissingParameters returns a slice of parameters that are in TemplateParamList but not in ParamList.
func CollectMissingParameters(ParamList, TemplateParamList *[]string) []string {

	var missingParams []string
	for _, param := range *TemplateParamList {
		if !slices.Contains(*ParamList, param) {

			missingParams = append(missingParams, param)
		}

	}

	return missingParams

}
