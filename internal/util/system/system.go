package system

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

/*
工具包
系统函数相关
*/

//Hostname 获取主机名
func Hostname() string {
	if hostname, err := os.Hostname(); err != nil {
		return fmt.Sprintf("unknown-%v", err)
	} else {
		return hostname
	}
}

//UniqueID 计算UUID(Google)
func UniqueID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
