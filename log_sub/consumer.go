package log_sub

import (
	"bufio"
	"errors"
	"github.com/BabySid/gobase"
	"io"
	"os"
	"strings"
	"time"
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
	Layout string // e.g. 2006010215.log
}

type Config struct {
	Location          *SeekInfo // if nil it will consumer log from cur time
	DateTimeLogLayout *DateTimeLayout

	Logger gobase.Logger
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
		config.Logger = &gobase.StdErrLogger{}
	}

	if config.DateTimeLogLayout == nil {
		return nil, errors.New("invalid config")
	}

	now := time.Now()
	if config.Location == nil {
		config.Location = &SeekInfo{
			FileName: now.Format(config.DateTimeLogLayout.Layout),
			Offset:   0,
			Whence:   0,
		}
	}
	step := verifyLogStep(config.DateTimeLogLayout.Layout, config.Location.FileName)

	meta := dateTimeLog{
		cur:  now,
		step: step,
	}

	startTime, err := time.ParseInLocation(config.DateTimeLogLayout.Layout, config.Location.FileName, time.Local)
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

	c.Logger.Infof("Consumer.curDateTimeLogMeta: cur=%s, step=%d", gobase.FormatTimeStamp(meta.cur.Unix()), step)

	go c.startConsume()

	return &c, nil
}

func (c *Consumer) startConsume() {
	err := c.openFile(c.curDateTimeLogMeta.cur.Format(c.DateTimeLogLayout.Layout))
	gobase.True(err == nil)

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

			c.Logger.Infof("waitFileChanges. nextFile=%s", nxt)
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
	c.Lines <- &Line{Text: line, Err: err}
}

func (c *Consumer) readLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		// Note ReadString "returns the data read before the error" in
		// case of an error, including EOF, so we return it as is. The
		// caller is expected to process it if err is EOF.
		return line, err
	}

	line = strings.TrimRight(line, "\n")

	return line, err
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
		name := nxt.Format(c.DateTimeLogLayout.Layout)
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

	c.Logger.Infof("openFile(%s) successful", c.file.Name())
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
		c.Logger.Infof("getNextFile %v", gobase.AbbreviateArray(newFiles))

		var newFile nxtFile
		for _, f := range newFiles {
			_, err := os.Stat(f.Name)
			if err == nil {
				newFile = f
				c.Logger.Infof("newFile(%s) from getNextFile exist", newFile)
				break
			}
		}

		latest, err := os.Stat(c.file.Name())
		if err != nil {
			c.Logger.Warnf("os.Stat(%s) encounter an error=%v", c.file.Name(), err)
			return nxtFile{}, err
		}

		if size > 0 && size < latest.Size() {
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
