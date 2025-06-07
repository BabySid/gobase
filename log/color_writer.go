package log

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/BabySid/gobase"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

type colorWriter struct {
	writer      io.Writer
	shouldColor bool
	isJSON      bool
	re          *regexp.Regexp
}

func newColorWriter(w io.Writer, isJSON bool) *colorWriter {
	isTerminal := isTerminal(w) // 检测是否为终端
	c := &colorWriter{
		writer:      w,
		shouldColor: isTerminal && !isJSON, // JSON 格式不添加颜色
		isJSON:      isJSON,
	}
	if c.shouldColor {
		// 匹配 slog TextHandler 的输出格式
		// e.g. time=2025-06-07T14:15:32.307+08:00 level=DEBUG source=/path/to/code.go:123 msg="this is a msg" key=value
		c.re = regexp.MustCompile(`^time=([^\s]+)\s+level=(\w+)(?:\s+source=([\S+]+))?\s+(.*)$`)
	}
	return c
}

func (cw *colorWriter) Write(p []byte) (n int, err error) {
	if !cw.shouldColor {
		return cw.writer.Write(p) // 直接输出，不添加颜色
	}

	// 为 TextHandler 输出添加颜色
	coloredOutput := cw.addColorToLogLine(string(p))
	return cw.writer.Write([]byte(coloredOutput))
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

var (
	timeColor     = termenv.ANSIBrightBlack
	logLevelColor = map[string]termenv.Color{
		"TRACE": termenv.ANSIBrightBlack,
		"DEBUG": termenv.ANSICyan,
		"INFO":  termenv.ANSIGreen,
		"WARN":  termenv.ANSIYellow,
		"ERROR": termenv.ANSIRed,
	}
	sourceColor = termenv.RGBColor("#6C7B95")
	msgColor    = termenv.ANSIBrightWhite
	// attrColor   = termenv.RGBColor("#8A7CA8")
)

// 为日志行添加颜色
func (cw *colorWriter) addColorToLogLine(line string) string {
	matches := cw.re.FindStringSubmatch(strings.TrimSpace(line))

	if len(matches) < 4 {
		return line // 格式不匹配，原样返回
	}

	dt, _ := gobase.ParseTimestamp(matches[1])
	ts := termenv.String(gobase.FormatTimeStamp(dt.Unix())).Foreground(timeColor) // 时间戳
	lvl := termenv.String(matches[2]).Foreground(logLevelColor[matches[2]])       // 日志级别
	src := termenv.String(matches[3]).Foreground(sourceColor)                     // 源文件（可能为空）
	msg := termenv.String(matches[4]).Foreground(msgColor)                        // 消息内容

	return fmt.Sprintf("%s %s %s %s\n", ts, lvl, src, msg)
}
