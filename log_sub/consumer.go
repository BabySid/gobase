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

	watcher gobase.FileWatcher
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

			newFiles := c.getNextFile()
			fileArr := make([]string, len(newFiles))

			for i, f := range newFiles {
				fileArr[i] = f.Name
			}

			nxt, err := c.waitFileChanges(fileArr)
			if err != nil {
				c.sendLine("", err)
				return
			}

			c.Logger.Infof("waitFileChanges. nextFile=%s", nxt)
			for _, f := range newFiles {
				if f.Name == nxt {
					err = c.openFile(f.Name)
					gobase.True(err == nil)
					c.curDateTimeLogMeta.cur = f.Ts
				}
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
	c.watcher = gobase.NewPollingFileWatcher()
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

// waitFileChanges return next or continue to read from current file
// but now, we only support next
func (c *Consumer) waitFileChanges(newFiles []string) (string, error) {
	pos, err := c.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	if c.watcher == nil {
		c.setFileWatcher()
	}

	events := c.watcher.ChangeEvents(c.file.Name(), pos, newFiles)

	select {
	case data := <-events.Created:
		c.Logger.Debug("got events.Created %+v", data)
		return data.NxtFile, nil
	case data := <-events.Modified:
		c.Logger.Debug("got events.Modified %+v", data)
		return data.NxtFile, nil
	case data := <-events.Deleted:
		c.Logger.Debug("got events.Deleted %+v", data)
		return data.NxtFile, nil
	case data := <-events.Truncated:
		c.Logger.Debug("got events.Truncated %+v", data)
		return data.NxtFile, nil
	}
}

func (c *Consumer) Close() {
	if c.file == nil {
		return
	}

	_ = c.file.Close()
	c.file = nil
}
