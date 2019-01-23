package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"
)

// CommandResult defines the execute command's result
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// Run is used to execute command with timeout
func Run(timeout int, bin string, args ...string) (*CommandResult, error) {
	var (
		err      error
		exitCode int
	)
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

	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return nil, err
	}

	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return nil, err
	}

	return &CommandResult{
		Stdout:   string(stdout),
		Stderr:   string(stderr),
		ExitCode: exitCode}, nil
}
