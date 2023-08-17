package apm

import (
	"bufio"
	"bytes"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var logger *Log

// NewLogger creates a new APM logger
func NewLogger(name string, limit int) (*Log, error) {
	if "" == name {
		return nil, errors.Errorf(constant.SystemInternalError, "Invalid apm logger name:[%s]", name)
	}
	dirPath := name[0:strings.LastIndex(name, "/")]
	if len(dirPath) > 0 && "." != dirPath {
		if err := os.MkdirAll(dirPath, os.ModePerm); nil != err {
			return nil, errors.Wrap(constant.SystemInternalError, err, 0)
		}
	}
	logFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	count, err := lineCounter(bufio.NewReader(logFile))
	if nil != err {
		return nil, err
	}
	loger := log.New(logFile, "", 0)
	return &Log{
		log:   loger,
		limit: uint(limit),
		name:  name,
		line:  uint(count),
	}, nil
}

// Log is a APM log controller
type Log struct {
	fnum  int
	name  string
	f     *os.File
	log   *log.Logger
	mu    sync.Mutex
	ln    sync.Mutex
	limit uint
	line  uint
}

func (l *Log) count() uint {
	l.ln.Lock()
	defer l.ln.Unlock()
	l.line++
	c := l.line
	return c
}
func (l *Log) writer(b []byte) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintf(os.Stderr, "failed to do APM log, error:%++v\n", e)
		}
	}()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Writer().Write(b)
	line := l.count()
	if line >= l.limit {
		return l.refile()
	}
	return nil
}

func (l *Log) refile() error {
	l.ln.Lock()
	defer l.ln.Unlock()
	err := l.f.Close()
	//if err != nil {
	//	return err
	//}
	err = os.Rename(l.name, fmt.Sprintf("%s.%d", l.name, l.fnum))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to rename file, error:%++v\n", err)
		//return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	err = os.Chmod(fmt.Sprintf("%s.%d", l.name, l.fnum), 0766)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to chmod file, error:%++v\n", err)
		//return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	l.fnum = (l.fnum + 1) % 10
	logFile, err := os.OpenFile(l.name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		fmt.Fprintf(os.Stderr, "failed to open file, error:%++v\n", err)
		return nil
		//return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	l.line = 0
	logger := log.New(logFile, "", 0)
	l.log = logger
	return nil
}

// APMlogf formats according to a format specifier and writes the resulting string to the APM file
func (l *Log) APMlogf(format string, v ...interface{}) error {
	l.writer([]byte(fmt.Sprintf(format+"\n", v...)))
	return nil
}

// APMlog writes the input log to the APM file
func (l *Log) APMlog(log string) error {
	if "" == log {
		return nil
	}
	l.writer([]byte(log + "\n"))
	return nil
}

func lineCounter(r io.Reader) (int, error) {
	var readSize int
	var err error
	var count int

	buf := make([]byte, 1024)

	for {
		readSize, err = r.Read(buf)
		if err != nil {
			break
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], '\n')
			if i == -1 || readSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
	}
	if readSize > 0 && count == 0 || count > 0 {
		count++
	}
	if err == io.EOF {
		return count, nil
	}

	return count, nil
}
