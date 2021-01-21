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

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	"github.com/google/shlex"
)

type DockerExec struct {
	image   string
	timeout time.Duration
}

func NewDockerExec(image string, timeout time.Duration) *DockerExec {
	return &DockerExec{image: image, timeout: timeout}
}

func (e *DockerExec) Exec(code []byte, arg string, env interface{}) (ExecResult, error) {
	// TODO: Handle env if we are to revive Docker
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
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
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
