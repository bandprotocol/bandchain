package executor

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TODO: Make this configurable
const MAX_OUTPUT_SIZE = 512

type DockerExec struct {
	image string
}

func NewDockerExec(image string) *DockerExec {
	return &DockerExec{image: image}
}

func (e *DockerExec) Exec(timeout time.Duration, code []byte, args ...string) (ExecResult, error) {
	dir, err := ioutil.TempDir("", "executor")
	if err != nil {
		return ExecResult{}, err
	}
	defer os.RemoveAll(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "exec"), code, 0777)
	if err != nil {
		return ExecResult{}, err
	}
	name := filepath.Base(dir)
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
			return ExecResult{}, nil
		}
	}
	output, err := ioutil.ReadAll(io.LimitReader(&buf, MAX_OUTPUT_SIZE))
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Output: output, Code: exitCode}, nil
}
