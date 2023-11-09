package utils

// getStringValue извлекает строковое значение из map[string]interface{}
func GetStringValue(data map[string]interface{}, key string) string {
	value, ok := data[key].(string)
	if !ok {
		return ""
	}
	return value
}