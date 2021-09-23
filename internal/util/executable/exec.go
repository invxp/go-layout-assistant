package executable

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"time"
)

/*
工具包
执行Shell命令
*/

//Exec 执行一个命令
func Exec(command, workDir string, timeoutSecond uint) (stdout string, stderr string, err error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	var cmd *exec.Cmd

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Dir = workDir

	e := cmd.Run()

	return outBuf.String(), errBuf.String(), e
}
