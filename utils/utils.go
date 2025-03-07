package utils

import (
	"reflect"
	"slices"
)

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

func ValidateParameters(ParamList, TemplateParamList *[]string) bool {

	slices.Sort(*ParamList)
	slices.Sort(*TemplateParamList)

	return reflect.DeepEqual(*ParamList, *TemplateParamList)
}

func CollectMissingParameters(ParamList, TemplateParamList *[]string) []string {

	var missingParams []string
	for _, param := range *TemplateParamList {
		if !slices.Contains(*ParamList, param) {

			missingParams = append(missingParams, param)
		}

	}

	return missingParams

}
