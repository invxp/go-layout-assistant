package daemon

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
工具包
Daemonize让进程在后台执行（不需要nohup）
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

//Daemonize 以后台方式启动应用(重新启动一个自己)
func Daemonize(isDaemonize bool, argFlag, daemonizeFlag string) {
	if len(os.Args) > 2 && os.Args[1] == "-"+daemonizeFlag {
		log.Printf("application %d run in daemonize: %v, ppid: %d\n", os.Getpid(), os.Args, os.Getppid())
		return
	}

	if !isDaemonize {
		log.Printf("application %d run in frontend: %v, ppid: %d\n", os.Getpid(), os.Args, os.Getppid())
		return
	}

	runArg := []string{"-" + daemonizeFlag, strconv.Itoa(os.Getpid())}
	daemonizePos := 0
	for i := 1; i < len(os.Args); i++ {
		if daemonizePos > 0 {
			daemonizePos = 0
			val := strings.ToLower(os.Args[i])
			if val == "t" || val == "f" || val == "1" || val == "0" || val == "true" || val == "false" {
				continue
			}
		}

		if os.Args[i] == "-"+daemonizeFlag {
			i++
			continue
		}

		if os.Args[i] == "-"+argFlag {
			daemonizePos = i
			continue
		}

		runArg = append(runArg, os.Args[i])
	}

	log.Printf("application %d prepare to run in daemonize: %s\n", os.Getpid(), runArg)

	cmd := exec.Command(os.Args[0], runArg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = fullPath()

	err := cmd.Start()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("application %d exit for daemonize: %s, release: %v\n", os.Getpid(), runArg, cmd.Process.Release())

	os.Exit(0)
}
