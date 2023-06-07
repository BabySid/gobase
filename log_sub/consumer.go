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

type Line struct {
	Text string
	Err  error
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

	logger gobase.Logger
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
	nxtFileName        string

	watcher gobase.FileWatcher
}

const (
	defaultBufSize = 1024
	maxReadSize    = 1024 * 1024
)

func NewConsumer(config Config) (*Consumer, error) {
	if config.logger == nil {
		config.logger = &gobase.StdErrLogger{}
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

	startTime, err := time.Parse(config.DateTimeLogLayout.Layout, config.Location.FileName)
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

		watcher: nil,
	}

	c.logger.Infof("Consumer.curDateTimeLogMeta: %+v", meta)

	go c.startConsume()

	return &c, nil
}

func (c *Consumer) startConsume() {
	err := c.openFile()
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

			err = c.waitFileChanges()
			if err != nil {
				c.sendLine("", err)
				return
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

func (c *Consumer) setFileWatcher() {
	c.watcher = gobase.NewPollingFileWatcher(c.file.Name(), c.nxtFileName)
}

func (c *Consumer) openReader() {
	c.reader = bufio.NewReaderSize(c.file, maxReadSize)
}

func (c *Consumer) setNextFile() {
	multi := time.Duration(1)
	if c.curDateTimeLogMeta.step == daily {
		multi = 24
	}
	c.nxtFileName = c.curDateTimeLogMeta.cur.Add(time.Hour * multi).Format(c.DateTimeLogLayout.Layout)
}

func (c *Consumer) openFile() error {
	fName := ""
	if c.file == nil { // first time
		fName = c.curDateTimeLogMeta.cur.Format(c.DateTimeLogLayout.Layout)
	} else {
		fName = c.nxtFileName
	}

	file, err := os.Open(fName)
	if err != nil {
		return err
	}

	if c.file != nil {
		c.Close()
	}

	c.file = file
	c.setNextFile()

	c.setFileWatcher()
	c.openReader()

	c.logger.Infof("openFile(%s) successful", fName)
	return nil
}

// waitFileChanges return next or continue to read from current file
// but now, we only support next
func (c *Consumer) waitFileChanges() error {
	pos, err := c.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	events, err := c.watcher.ChangeEvents(pos)
	if err != nil {
		return err
	}

	select {
	case <-events.Created:
		c.logger.Debug("got events.Created")
		return c.openFile()
	case <-events.Modified:
		c.logger.Debug("got events.Modified")
		return nil
	case <-events.Deleted:
		c.logger.Debug("got events.Deleted")
		return c.openFile()
	case <-events.Truncated:
		c.logger.Debug("got events.Truncated")
		return c.openFile()
	}
}

func (c *Consumer) Close() {
	if c.file == nil {
		return
	}

	_ = c.file.Close()
	c.file = nil
}
