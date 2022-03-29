package gobase

import (
	"fmt"
	"github.com/go-cmd/cmd"
)

const (
	SuccessCode = 1
)

var (
	charNeedQuota = map[rune]int{
		' ':  1,
		'&':  1,
		'(':  1,
		')':  1,
		'[':  1,
		']':  1,
		'{':  1,
		'}':  1,
		'^':  1,
		'=':  1,
		';':  1,
		'!':  1,
		'\'': 1,
		'+':  1,
		',':  1,
		'`':  1,
		'~':  1}
)

func ExecExplorer(params []string) error {
	c := cmd.NewCmd("explorer", params...)
	status := <-c.Start()
	if status.StartTs > 0 && (status.Exit != SuccessCode || !status.Complete) {
		return fmt.Errorf("%d", status.Exit)
	}
	if status.StartTs == 0 && status.Error != nil {
		return status.Error
	}
	return nil
}

func ExecApp(appFullPath string) error {
	path := QuotaPath(appFullPath)

	c := cmd.NewCmd("cmd", []string{"/c", path}...)
	_ = <-c.Start()
	return nil
}

func QuotaPath(fullPath string) string {
	path := ""
	for _, ch := range fullPath {
		if _, ok := charNeedQuota[ch]; ok {
			path += string('^') + string(ch)
		} else {
			path += string(ch)
		}
	}

	return path
}
