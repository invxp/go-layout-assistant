package config

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

/*
工具包
读取toml配置
*/

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

//Load 加载Config文件
func Load(fileName string, config interface{}) (err error) {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(fullPath(), fileName)
	}
	_, err = toml.DecodeFile(fileName, config)
	return
}
