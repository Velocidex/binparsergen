package binparsergen

import (
	"reflect"
	"sort"
)

func InString(hay []string, needle string) bool {
	for _, x := range hay {
		if x == needle {
			return true
		}
	}

	return false
}

func SortedKeys(any interface{}) []string {
	if reflect.TypeOf(any).Kind() != reflect.Map {
		return nil
	}

	result := []string{}
	for _, k := range reflect.ValueOf(any).MapKeys() {
		result = append(result, k.Interface().(string))
	}

	sort.Strings(result)

	return result
}

func SortedIntKeys(any interface{}) []int {
	if reflect.TypeOf(any).Kind() != reflect.Map {
		return nil
	}

	result := []int{}
	for _, k := range reflect.ValueOf(any).MapKeys() {
		result = append(result, k.Interface().(int))
	}

	sort.Ints(result)

	return result
}
