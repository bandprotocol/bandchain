package byteexec

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
)

var (
	fileMode = os.FileMode(0744)
)

func writeFile(executable []byte) (string, string, error) {
	dir, err := ioutil.TempDir("/tmp", "temp")
	if err != nil {
		return "", "", err
	}
	filename := filepath.Join(dir, "exec")
	err = ioutil.WriteFile(filename, executable, fileMode)
	if err != nil {
		return "", "", err
	}
	filename, err = filepath.Abs(filename)
	if err != nil {
		return "", "", err
	}
	return dir, filename, nil
}

// RunOnLocal spawns a new subprocess and runs the given executable. NOT SAFE!
func RunOnLocal(executable []byte, timeOut time.Duration, arg string) ([]byte, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return nil, err
	}

	dir, filename, err := writeFile(executable)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dir) // clean up

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)

	defer cancel()

	return exec.CommandContext(ctx, filename, args...).Output()
}

// RunOnDocker runs the given executable in a new docker container.
func RunOnDocker(executable []byte, timeOut time.Duration, arg string) ([]byte, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return nil, err
	}

	dir, filename, err := writeFile(executable)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	rawID, err := exec.Command(
		"docker", "run", "-d", "--rm", "band-provider", "sleep", fmt.Sprintf("%d", int(timeOut.Seconds())),
	).Output()
	if err != nil {
		return nil, err
	}

	containerID := strings.TrimSpace(string(rawID))
	defer exec.Command("docker", "stop", containerID).Output()

	_, err = exec.Command(
		"docker", "cp", filename, fmt.Sprintf("%s:/exec", containerID),
	).Output()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	newArgs := append([]string{"exec", containerID, "./exec"}, args...)

	return exec.CommandContext(ctx, "docker", newArgs...).Output()
}
