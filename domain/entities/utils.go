package entities

import "errors"

func checkRequirement(data map[string]string, requiredKeys ...string) error {
	for _, key := range requiredKeys {
		if _, ok := data[key]; !ok {
			return errors.New(key + " is required")
		}
	}
	return nil
}
