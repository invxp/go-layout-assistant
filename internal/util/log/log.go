package log

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
工具包
日志库
*/

type Log struct {
	log *log.Logger
	mu  *sync.Mutex
}

func fullPath() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path[0:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator), err
}

//New 新建一个日志
//pathName 				- 日志目录
//fileName 				- 文件名称
//maxAgeHours 			- 最大生命周期(小时)
//maxFileSizeMegabytes  - 文件超过多大开始流转(MB)
func New(pathName, fileName string, maxAgeHours, maxFileSizeMegabytes uint) (*Log, error) {
	if !filepath.IsAbs(pathName) {
		path, err := fullPath()
		if err != nil {
			return nil, err
		}
		pathName = filepath.Join(path, pathName)
	}

	fullFilePathName := filepath.Join(pathName, fileName)
	writer, err := rotatelogs.New(
		fullFilePathName+".%Y-%m-%d_%H:%M",
		rotatelogs.WithLinkName(fullFilePathName),
		rotatelogs.WithMaxAge(time.Duration(maxAgeHours)*time.Hour),
		rotatelogs.WithRotationSize(int64(maxFileSizeMegabytes)*1024*1024),
		rotatelogs.ForceNewFile(),
	)

	err = os.MkdirAll(pathName, 0644)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(fullFilePathName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
