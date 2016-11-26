package main

func IfNotEmpty(prefix, value string) string {
	if len(value) > 0 {
		return prefix + value
	}
	return ""
}
