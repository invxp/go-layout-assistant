package io

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
工具包
IO/系统路径相关
*/

//FullPath 获取当前进程路径
func FullPath() string {
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

//CurrentExecutableName 获取当前进程名称
func CurrentExecutableName() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
	}
	return file[strings.LastIndex(file, string(os.PathSeparator))+len(string(os.PathSeparator)):]
}
