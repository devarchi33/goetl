package entities

import "errors"

func checkRequirement(data map[string]string, requiredKeys ...string) error {
	for _, key := range requiredKeys {
		if v, ok := data[key]; !ok || len(v) == 0 {
			return errors.New(key + " is required")
		}
	}
	return nil
}
