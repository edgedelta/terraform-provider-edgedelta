package edgedelta

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJSON validates that a string is valid JSON
func validateJSON(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v == "" {
		return nil, nil
	}
	var js interface{}
	if err := json.Unmarshal([]byte(v), &js); err != nil {
		errs = append(errs, fmt.Errorf("%q must be valid JSON, got: %s, error: %v", key, v, err))
	}
	return warns, errs
}

// suppressEquivalentJSON is a DiffSuppressFunc that suppresses diffs for equivalent JSON
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "" {
		return true
	}
	if old == "" || new == "" {
		return false
	}

	var oldJSON, newJSON interface{}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return false
	}

	return reflect.DeepEqual(oldJSON, newJSON)
}

// stringSliceToInterface converts a []string to []interface{} for Terraform state
func stringSliceToInterface(s []string) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// interfaceSliceToStringSlice converts []interface{} from Terraform state to []string
func interfaceSliceToStringSlice(s []interface{}) []string {
	result := make([]string, len(s))
	for i, v := range s {
		if v != nil {
			result[i] = v.(string)
		}
	}
	return result
}

// jsonMapToString converts a map[string]interface{} to a JSON string
func jsonMapToString(m map[string]interface{}) (string, error) {
	if m == nil {
		return "", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal definition to JSON: %v", err)
	}
	return string(b), nil
}

// stringToJSONMap converts a JSON string to map[string]interface{}
func stringToJSONMap(s string) (map[string]interface{}, error) {
	if s == "" {
		return nil, nil
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return result, nil
}

// stringToJSONArray converts a JSON string to []map[string]interface{}
func stringToJSONArray(s string) ([]map[string]interface{}, error) {
	if s == "" {
		return nil, nil
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON array: %v", err)
	}
	return result, nil
}

// jsonArrayToString converts a []map[string]interface{} to a JSON string
func jsonArrayToString(arr []map[string]interface{}) (string, error) {
	if arr == nil {
		return "", nil
	}
	b, err := json.Marshal(arr)
	if err != nil {
		return "", fmt.Errorf("failed to marshal array to JSON: %v", err)
	}
	return string(b), nil
}
