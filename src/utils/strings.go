package utils

import (
	"hash/fnv"
	"strconv"
	"strings"
)

func HashString(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(h.Sum32()))
}

func TruncateString(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

func ObfuscateString(s string, length int) string {
	dotsToAppend := len(s)
	result := ""
	if len(s) > length {
		result = s[:length]
		dotsToAppend = len(s) - length
	}
	result += strings.Repeat("*", dotsToAppend)
	return result
}
