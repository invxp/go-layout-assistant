package main

import (
	"fmt"
	"github.com/invxp/go-layout-assistant/internal/util/convert"
	"github.com/invxp/go-layout-assistant/internal/util/io"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultServiceNameLower = "application"
	defaultServiceNameUpper = "Application"
	defaultModPath          = "github.com/invxp/go-layout-assistant"
	defaultInstallDir       = "/home/application"
)

var (
	blacklist = map[string]struct{}{
		".git":          {},
		".idea":         {},
		"logs":          {},
		"cmd/assistant": {}}

	whitelist = map[string]struct{}{
		"go.mod": {},
		"go.sum": {}}
)

func main() {
	var serviceName string
	var modPath string
	var installDir string

	fmt.Printf("---input your service name(app)*: ")
	_, _ = fmt.Scanln(&serviceName)
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		log.Printf("create service name must be exist")
		os.Exit(0)
	}

	fmt.Printf("---input your go mod path(github.com/invxp/fsm)*: ")
	_, _ = fmt.Scanln(&modPath)
	modPath = strings.TrimSpace(modPath)
	if modPath == "" {
		log.Printf("create mod path must be exist")
		os.Exit(0)
	}
	fmt.Printf("---input %s install dir(/home/%s): ", serviceName, serviceName)
	_, _ = fmt.Scanln(&installDir)
	if err := CreateService(serviceName, modPath, installDir); err != nil {
		log.Printf("create service: %s failed: %v", serviceName, err)
	} else {
		log.Printf("create service: %s success", serviceName)
	}

	if err := FixServiceCodes(serviceName, strings.Split(serviceName, "-")[len(strings.Split(serviceName, "-"))-1]); err != nil {
		log.Printf("fix service code: %s failed: %v.%v", serviceName, err, os.RemoveAll(serviceName))
	} else {
		log.Printf("fix service code: %s success: %v", serviceName, os.RemoveAll(filepath.Join(serviceName, serviceName)))
	}

	os.Exit(0)
}

//CreateService 克隆一个应用
func CreateService(serviceName, modPath, installDir string) error {
	if installDir == "" {
		installDir = filepath.Join("/home", serviceName)
	}

	err := os.RemoveAll(serviceName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(serviceName, 0644)
	if err != nil {
		return err
	}

	log.Printf("prepare create %s success, installDir: %s, modPath: %s", serviceName, installDir, modPath)

	fullPath := io.FullPath()

	return filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasPrefix(path, filepath.Join(fullPath, serviceName)) {
			log.Printf("skip self: %s from blacklist", path)
			return err
		}

		for name := range blacklist {
			if strings.HasPrefix(path, filepath.Join(fullPath, name)) {
				log.Printf("skip path: %s from blacklist", path)
				return err
			}
		}

		path = strings.ReplaceAll(path, fullPath, "")

		if dir := filepath.Dir(path); dir != "." {
			return replaceFile(path, serviceName, modPath, installDir, true)
		}

		if _, exists := whitelist[path]; exists {
			return replaceFile(path, serviceName, modPath, installDir, true)
		} else {
			log.Printf("skip file: %s from whitelist", path)
		}

		return err
	})
}

//replaceFile 克隆文件+内容
func replaceFile(originFilename, serviceName, modPath, installDir string, replaceContent bool) error {
	var write string

	upperServiceName := strings.ToUpper(serviceName[:1]) + serviceName[1:]

	if dir := filepath.Dir(originFilename); dir != "." {
		if replaceContent {
			dir = strings.ReplaceAll(dir, defaultServiceNameLower, serviceName)
		}
		if err := os.MkdirAll(filepath.Join(serviceName, dir), 0644); err != nil {
			return err
		}
	}

	read, err := ioutil.ReadFile(originFilename)
	if err != nil {
		return err
	}

	write = convert.ByteToString(read)

	if replaceContent {
		write = strings.ReplaceAll(convert.ByteToString(read), defaultModPath, modPath)
		write = strings.ReplaceAll(write, defaultInstallDir, installDir)
		write = strings.ReplaceAll(write, defaultServiceNameLower, serviceName)
		write = strings.ReplaceAll(write, defaultServiceNameUpper, upperServiceName)
		originFilename = strings.ReplaceAll(originFilename, defaultServiceNameLower, serviceName)
	}

	return errors.Wrapf(ioutil.WriteFile(filepath.Join(serviceName, originFilename), convert.StringToByte(write), 0644), "create file: %s", filepath.Join(serviceName, originFilename))
}

//FixServiceCodes 修复带'-'符号的工程代码
func FixServiceCodes(serviceName, fixedName string) error {
	if serviceName == fixedName {
		return nil
	}

	if err := os.MkdirAll(filepath.Join(serviceName, fixedName), 0644); err != nil {
		return err
	}

	return filepath.WalkDir(serviceName, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if path == filepath.Join(serviceName, "cmd", serviceName, "main.go") {
			return fixFile(path, path, serviceName, fixedName)
		}

		if strings.HasPrefix(path, filepath.Join(serviceName, serviceName)) {
			return fixFile(path, filepath.Join(serviceName, fixedName, strings.ReplaceAll(strings.Split(path, string(os.PathSeparator))[len(strings.Split(path, string(os.PathSeparator)))-1], serviceName, fixedName)), serviceName, fixedName)
		}

		return err
	})
}

func fixFile(originFilename, destFilename, serviceName, fixedName string) error {
	read, err := ioutil.ReadFile(originFilename)
	if err != nil {
		return err
	}

	upperServiceName := strings.ToUpper(serviceName[:1]) + serviceName[1:]
	upperFixedName := strings.ToUpper(fixedName[:1]) + fixedName[1:]

	write := convert.ByteToString(read)

	write = strings.ReplaceAll(convert.ByteToString(read), serviceName, fixedName)
	write = strings.ReplaceAll(write, upperServiceName, upperFixedName)

	return errors.Wrapf(ioutil.WriteFile(destFilename, convert.StringToByte(write), 0644), "create file: %s", filepath.Join(destFilename))
}
