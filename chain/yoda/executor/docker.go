package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/google/shlex"
)

type DockerExec struct {
	image string
}

var testProgram []byte = []byte("#!/usr/bin/env python3\nimport sys\nprint(sys.argv[1])")

func NewDockerExec(image string) *DockerExec {
	exec := &DockerExec{image: image}
	res, err := exec.Exec(5*time.Second, testProgram, "TEST_ARG")
	if err != nil {
		panic(fmt.Sprintf("NewDockerExec: failed to run test program: %s", err.Error()))
	}
	if res.Code != 0 {
		panic(fmt.Sprintf("NewDockerExec: test program returned nonzero code: %d", res.Code))
	}
	if string(res.Output) != "TEST_ARG\n" {
		panic(fmt.Sprintf("NewDockerExec: test program returned wrong output: %s", res.Output))
	}
	return exec
}

func (e *DockerExec) Exec(timeout time.Duration, code []byte, arg string) (ExecResult, error) {
	dir, err := ioutil.TempDir("/tmp", "executor")
	if err != nil {
		return ExecResult{}, err
	}
	defer os.RemoveAll(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "exec"), code, 0777)
	if err != nil {
		return ExecResult{}, err
	}
	name := filepath.Base(dir)
	args, err := shlex.Split(arg)
	if err != nil {
		return ExecResult{}, err
	}
	dockerArgs := append([]string{
		"run", "--rm",
		"-v", dir + ":/scratch:ro",
		"--name", name,
		e.image,
		"/scratch/exec",
	}, args...)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		exec.Command("docker", "kill", name).Start()
		return ExecResult{}, ErrExecutionimeout
	}
	exitCode := uint32(0)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = uint32(exitError.ExitCode())
		} else {
			return ExecResult{}, err
		}
	}
	output, err := ioutil.ReadAll(io.LimitReader(&buf, types.MaxDataSize))
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Output: output, Code: exitCode}, nil
}
