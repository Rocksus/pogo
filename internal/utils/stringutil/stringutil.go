package stringutil

import "strings"

func GetFirstWord(str string) string {
	words := strings.Fields(str)
	if len(words) == 0 {
		return ""
	}
	return words[0]
}
