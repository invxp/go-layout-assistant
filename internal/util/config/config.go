package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

//MustLoad 加载Config文件
func MustLoad(fileName string, config interface{}) {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(fullPath(), fileName)
	}

	_, err := toml.DecodeFile(fileName, config)

	if err != nil {
		log.Panic(err)
	}
}
