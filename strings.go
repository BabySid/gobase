package gobase

import (
	"encoding/json"
	"unicode/utf8"
)

func CutUTF8(str string, start int, end int) string {
	if end <= utf8.RuneCountInString(str) {
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

		return aft
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
