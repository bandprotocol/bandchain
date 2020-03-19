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

	"github.com/levigross/grequests"
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
func RunOnDocker(executable []byte, sandboxMode bool, timeOut time.Duration, arg string) ([]byte, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return nil, err
	}

	dir, filename, err := writeFile(executable)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	commands := []string{"run"}
	if sandboxMode {
		commands = append(commands, "--runtime=runsc")
	}
	commands = append(
		commands, "-d", "--rm", "band-provider", "sleep", fmt.Sprintf("%d", int(timeOut.Seconds())),
	)
	rawID, err := exec.Command("docker", commands...).Output()

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

// RunOnAWSLambda runs the given executable on AMS Lambda platform.
func RunOnAWSLambda(executable []byte, timeOut time.Duration, arg string, executeEndPoint string) ([]byte, error) {
	resp, err := grequests.Post(
		executeEndPoint,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			JSON: map[string]string{
				"executable": string(executable),
				"calldata":   arg,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if resp.Ok != true {
		return nil, resp.Error
	}

	type result struct {
		Returncode int    `json:"returncode"`
		Stdout     string `json:"stdout"`
		Stderr     string `json:"stderr"`
	}

	r := result{}
	err = resp.JSON(&r)
	if err != nil {
		return nil, err
	}

	return []byte(r.Stdout), nil
}
