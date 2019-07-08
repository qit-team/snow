package logger

/**
 * 日志文件分割
 */
import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	rollDay = iota
	rollHour
)

const (
	defaultDayTimePattern  = "20060102"
	defaultHourTimePattern = "20060102-15"
)

type RollHook struct {
	dir          string
	name         string
	currFileTime string
	writer       *os.File
	timePattern  string
	lock         sync.Mutex
	logger       *logrus.Logger
}

func (rh *RollHook) openNewFile() (*os.File, error) {
	_, err := os.Stat(rh.dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rh.dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	newFileTime := time.Now().Format(rh.timePattern)
	newFileName := fmt.Sprintf("%s/%s.%s.log", rh.dir, rh.name, newFileTime)
	newWriter, err := os.OpenFile(newFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	rh.currFileTime = newFileTime

	return newWriter, nil
}

func NewRollHook(logger *logrus.Logger, dir, name string) (*RollHook, error) {
	rh := new(RollHook)
	rh.name = name
	rh.timePattern = defaultDayTimePattern
	rh.logger = logger
	rh.dir = dir

	writer, err := rh.openNewFile()
	if err != nil {
		return nil, err
	}
	rh.writer = writer
	logger.Out = writer

	return rh, nil
}

func (rh *RollHook) needRoll() bool {
	return rh.currFileTime != time.Now().Format(rh.timePattern)
}

func (rh *RollHook) roll() error {
	rh.lock.Lock()
	defer rh.lock.Unlock()

	if !rh.needRoll() {
		return nil
	}

	oldWriter := rh.writer
	newWriter, err := rh.openNewFile()
	if err != nil {
		return err
	}

	rh.writer = newWriter
	rh.logger.Out = newWriter

	err = oldWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rh *RollHook) SetRollType(rType int) {
	switch rType {
	case rollDay:
		rh.timePattern = defaultDayTimePattern
	case rollHour:
		rh.timePattern = defaultHourTimePattern
	}
}

func (rh *RollHook) Fire(entry *logrus.Entry) error {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	if rh.needRoll() {
		return rh.roll()
	}
	return nil
}

func (rh *RollHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}
