package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// DefaultTimeOut represents the default timeout of running command.
const DefaultTimeOut = 3600

// CommandResult defines the execute command's result
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// Run is used to execute command with timeout
func Run(timeout int, bin string, args ...string) (*CommandResult, error) {
	var (
		err       error
		exitCode  int
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)

	if timeout <= 0 {
		timeout = DefaultTimeOut
	}
	cmd := exec.Command(bin, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	defer stderrPipe.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	done := make(chan error)
	go func() {
		if _, err := bufio.NewReader(stdoutPipe).WriteTo(&stdoutBuf); err != nil {
			done <- err
		}
		if _, err := bufio.NewReader(stderrPipe).WriteTo(&stderrBuf); err != nil {
			done <- err
		}
		done <- cmd.Wait()
	}()

	select {
	case err = <-done:
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("run timeout, [%s] %v", bin, args)
	}

	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			exitCode = int(exit.ProcessState.Sys().(syscall.WaitStatus) / 256)
		}
	}

	logrus.Debugf("success to run [ %s %s ], result [%s %s]",
		bin, strings.Join(args, " "), stdoutBuf.String(), stderrBuf.String())

	return &CommandResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
	}, nil
}
