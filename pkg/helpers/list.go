package helpers

import "strings"

func ListToCommaString(labels []string) string {
	if len(labels) == 0 {
		return "None"
	}
	return strings.Join(labels, ", ")
}
