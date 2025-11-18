package utils

import "strings"

func CacheKey(keys ...string) string {
	return strings.Join(keys, ":")
}

func QueueName(keys ...string) string {
	return strings.ToUpper(strings.Join(keys, "_"))
}

func LogKey(keys ...string) string {
	return strings.ToLower(strings.Join(keys, "."))
}
