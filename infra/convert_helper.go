package infra

import "strconv"

// ConvertByteResult ...
func ConvertByteResult(source []map[string][]byte) []map[string]string {
	result := make([]map[string]string, 0)
	if source == nil || len(source) == 0 {
		return result
	}
	for _, sourceItem := range source {
		item := make(map[string]string, 0)
		if sourceItem == nil || len(sourceItem) == 0 {
			continue
		}
		for key, value := range sourceItem {
			item[key] = string(value)
		}
		result = append(result, item)
	}
	return result
}

// ConvertStringToInt string to int
func ConvertStringToInt(v string) int {
	result, _ := strconv.Atoi(v)
	return result
}
