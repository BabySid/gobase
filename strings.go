package gobase

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

func CutUTF8(str string, start int, end int, suffix string) string {
	if end < utf8.RuneCountInString(str) {
		aft := ""
		i := 0
		for _, v := range str {
			if i >= start && i < end {
				aft += string(v)
			}
			if i >= end {
				break
			}
			i++
		}

		return aft + suffix
	}
	return str
}

func CompressJson(js string) (string, error) {
	var v interface{}
	err := json.Unmarshal([]byte(js), &v)
	if err != nil {
		return "", err
	}

	dst, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}

func PrettyPrintJson(js string, indent string) (string, error) {
	var v interface{}
	err := json.Unmarshal([]byte(js), &v)
	if err != nil {
		return "", err
	}

	dst, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}

func SplitAndTrimSpace(s string, sep string) []string {
	rs := strings.Split(s, sep)

	tmp := make([]string, 0)
	for _, v := range rs {
		m := strings.TrimSpace(v)
		if m != "" {
			tmp = append(tmp, m)
		}
	}

	return tmp
}

func AbbreviateArray[T any](arr []T) string {
	switch len(arr) {
	case 0:
		return "[]"
	case 1:
		return fmt.Sprintf("[%+v]", arr[0])
	case 2:
		return fmt.Sprintf("[%+v,%+v]", arr[0], arr[1])
	default:
		return fmt.Sprintf("[%+v,...,%+v]", arr[0], arr[len(arr)-1])
	}
}
