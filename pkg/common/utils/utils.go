package utils

import "fmt"

func Bool(val interface{}) *bool {
	if val == nil {
		return nil
	}
	if val, ok := val.(bool); ok {
		return &val
	} else {
		return nil
	}
}

func String(val interface{}) *string {
	if val == nil {
		return nil
	}
	if str, ok := val.(string); ok {
		return &str
	} else {
		str = fmt.Sprintf("%s", val)
		return &str
	}
}
func Int(val interface{}) *int {
	if val == nil {
		return nil
	}
	if str, ok := val.(int); ok {
		return &str
	} else {
		return nil
	}
}
