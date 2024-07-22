package log_sub

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BabySid/gobase"
	mylog "github.com/BabySid/gobase/log"
)

type LineMeta struct {
	FileName string
}

type Line struct {
	Text string
	Err  error
	Meta LineMeta
}

type SeekInfo struct {
	FileName string
	Offset   int64
	Whence   int
}

type DateTimeLayout struct {
	FilePath string // e.g. /path/to/20230815.log
	Layout   string // e.g. 20060102.log
}

func (dt *DateTimeLayout) FormatFile(t time.Time) string {
	return filepath.Join(dt.filePath(), t.Format(dt.Layout))
}

func (dt *DateTimeLayout) filePath() string {
	return filepath.Dir(dt.FilePath)
}

type Config struct {
	Location          *SeekInfo // if nil it will consumer log from cur time
	DateTimeLogLayout *DateTimeLayout

	Logger mylog.Logger
}

const (
	hourly = iota
	daily
)

type dateTimeLog struct {
	cur  time.Time
	step int
}

type Consumer struct {
	Lines chan *Line
	Config

	curDateTimeLogMeta dateTimeLog
	file               *os.File
	reader             *bufio.Reader
}

const (
	defaultBufSize = 1024
	maxReadSize    = 1024 * 1024
)

func NewConsumer(config Config) (*Consumer, error) {
	if config.Logger == nil {
		config.Logger = mylog.NewSLogger(mylog.WithOutFile(mylog.StdErr))
	}

	if config.DateTimeLogLayout == nil {
		return nil, errors.New("invalid config")
	}

	now := time.Now()
	if config.Location == nil {
		config.Location = &SeekInfo{
			FileName: config.DateTimeLogLayout.FormatFile(now),
			Offset:   0,
			Whence:   0,
		}
	}
	step := verifyLogStep(*config.DateTimeLogLayout, config.Location.FileName)

	meta := dateTimeLog{
		cur:  now,
		step: step,
	}

	fName := filepath.Base(config.Location.FileName)
	startTime, err := time.ParseInLocation(config.DateTimeLogLayout.Layout, fName, time.Local)
	if err != nil {
		return nil, err
	}
	meta.cur = startTime

	c := Consumer{
		Lines:  make(chan *Line, defaultBufSize),
		Config: config,

		curDateTimeLogMeta: meta,
		file:               nil,
		reader:             nil,
	}

	c.Logger.Info("Consumer.curDateTimeLogMeta", slog.String("cur", gobase.FormatTimeStamp(meta.cur.Unix())), slog.Int("step", step))

	go c.startConsume()

	return &c, nil
}

func (c *Consumer) startConsume() {
	file := filepath.Join(
		c.DateTimeLogLayout.filePath(),
		c.curDateTimeLogMeta.cur.Format(c.DateTimeLogLayout.Layout))

	err := c.openFile(file)
	gobase.TrueF(err == nil, "openFile(%s) failed. err=%v", file, err)

	if c.Location != nil {
		_, err = c.file.Seek(c.Location.Offset, c.Location.Whence)
		if err != nil {
			c.sendLine("", err)
			return
		}
	}

	for {
		line, err := c.readLine()
		if err == nil {
			c.sendLine(line, nil)
		} else if err == io.EOF {
			if line != "" {
				c.sendLine(line, nil)
			}

			nxt, err := c.waitNxtFile()
			if err != nil {
				c.sendLine("", err)
				return
			}

			c.Logger.Info("waitFileChanges", slog.String("nextFile", nxt.Name))
			if nxt.Name != c.file.Name() {
				err = c.openFile(nxt.Name)
				gobase.True(err == nil)
				c.curDateTimeLogMeta.cur = nxt.Ts
			}

		} else {
			c.sendLine(line, err)
			return
		}
	}
}

func (c *Consumer) Tell() (*SeekInfo, error) {
	if c.file == nil {
		return nil, nil
	}

	offset, err := c.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	info := SeekInfo{
		FileName: c.file.Name(),
		Offset:   offset,
		Whence:   0,
	}

	return &info, nil
}

func (c *Consumer) sendLine(line string, err error) {
	c.Lines <- &Line{Text: line, Err: err, Meta: LineMeta{FileName: c.file.Name()}}
}

// readLine read a line unless meet a '\n' or some error except io.EOF
func (c *Consumer) readLine() (string, error) {
	var line string
	for {
		str, err := c.reader.ReadString('\n')
		line += str
		if err != nil {
			// Note ReadString "returns the data read before the error" in
			// case of an error, including EOF, so we return it as is. The
			// caller is expected to process it if err is EOF.
			if err == io.EOF && len(line) > 0 {
				if !strings.HasSuffix(line, "\n") {
					continue
				}
			}
			return line, err
		}
		break
	}

	line = strings.TrimRight(line, "\n")

	return line, nil
}

func (c *Consumer) openReader() {
	c.reader = bufio.NewReaderSize(c.file, maxReadSize)
}

type nxtFile struct {
	Name string
	Ts   time.Time
}

func (c *Consumer) getNextFile() []nxtFile {
	multi := time.Duration(1)
	if c.curDateTimeLogMeta.step == daily {
		multi = 24
	}

	f := make([]nxtFile, 0)

	nxt := c.curDateTimeLogMeta.cur
	for {
		nxt = nxt.Add(time.Hour * multi)
		if nxt.After(time.Now()) {
			break
		}
		name := c.DateTimeLogLayout.FormatFile(nxt)
		f = append(f, nxtFile{Name: name, Ts: nxt})
	}

	return f
}

func (c *Consumer) openFile(fName string) error {
	file, err := os.Open(fName)
	if err != nil {
		return err
	}

	if c.file != nil {
		c.Close()
	}

	c.file = file

	c.openReader()

	c.Logger.Info("openFile successful", slog.String("fileName", c.file.Name()))
	return nil
}

func (c *Consumer) waitNxtFile() (nxtFile, error) {
	size, err := c.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nxtFile{}, err
	}

	for {
		time.Sleep(time.Second)

		newFiles := c.getNextFile()
		c.Logger.Info("getNextFile", slog.Any("name", gobase.AbbreviateArray(newFiles)))

		var newFile nxtFile
		for _, f := range newFiles {
			_, err := os.Stat(f.Name)
			if err == nil {
				newFile = f
				c.Logger.Info("newFile from getNextFile exist", slog.String("name", newFile.Name))
				break
			}
		}

		latest, err := os.Stat(c.file.Name())
		if err != nil {
			c.Logger.Warn("os.Stat encounter an error", slog.String("name", c.file.Name()), slog.Any("err", err))
			return nxtFile{}, err
		}

		if size >= 0 && size < latest.Size() {
			return nxtFile{
				Name: c.file.Name(),
				Ts:   c.curDateTimeLogMeta.cur,
			}, nil
		}

		if newFile.Name != "" {
			return newFile, nil
		}
	}
}

func (c *Consumer) Close() {
	if c.file == nil {
		return
	}

	_ = c.file.Close()
	c.file = nil
}
