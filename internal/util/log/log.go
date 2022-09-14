package log

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

/*
工具包
日志库
*/

type Log struct {
	log *log.Logger
	mu  *sync.Mutex
}

func fullPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		log.Panic(err)
	}

	return path[0:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator)
}

//New 新建一个日志
//fileName 				- 文件名称
//maxAgeHours 			- 最大生命周期(小时)
//maxFileSizeMegabytes  - 文件超过多大开始流转(MB)
func New(fileName string, maxAgeHours, maxFileSizeMegabytes uint) (*Log, error) {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(fullPath(), fileName)
	}

	writer, err := rotatelogs.New(
		fileName+".%Y-%m-%d_%H:%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(maxAgeHours)*time.Hour),
		rotatelogs.WithRotationSize(int64(maxFileSizeMegabytes)*1024*1024),
		rotatelogs.ForceNewFile(),
	)

	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Log{log.New(io.MultiWriter(logFile, writer, os.Stdout), "", log.LstdFlags|log.Lshortfile), &sync.Mutex{}}, err
}

//Printf 打印日志(线程安全)
func (l *Log) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Printf(format, v...)
}
