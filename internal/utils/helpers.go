package utils

import "strings"

func SplitSafe(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "|")
}
