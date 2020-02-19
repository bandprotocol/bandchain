package byteexec

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	fileMode   = os.FileMode(0744)
	initMutex  sync.Mutex
	maximumTry = 10
)

// Exec is a handle to an executable that can be used to create an exec.Cmd
// using the Command method. Exec is safe for concurrent use.
type Exec struct {
	Filename string
}

func writeFile(data []byte) (string, error) {
	filename := fmt.Sprintf("./temp%d", rand.Uint64())
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, fileMode)
	defer file.Close()

	if err != nil {
		if !os.IsExist(err) {
			return "", fmt.Errorf("Unexpected error opening %s: %s", filename, err)
		}
		return "", nil
	}

	_, err = file.Write(data)
	if err != nil {
		os.Remove(filename)
		return "", fmt.Errorf("Unable to write to file at %s: %s", filename, err)
	}
	return filename, nil
}

// New creates a new instace of Exec
func New(data []byte) (Exec, error) {
	// Use initMutex to synchronize file operations by this process
	initMutex.Lock()
	defer initMutex.Unlock()

	for i := 0; i < maximumTry; i++ {
		filename, err := writeFile(data)
		if err != nil {
			return Exec{}, err
		}
		if filename != "" {
			return newExec(filename)
		}
	}
	return Exec{}, fmt.Errorf("Cannot create new file")
}

// Command creates an exec.Cmd using the supplied args.
func (be *Exec) Command(args ...string) *exec.Cmd {
	return exec.Command(be.Filename, args...)
}

// RunOnDocker runs command in new docker container
func (be *Exec) RunOnDocker(timeOut time.Duration, args ...string) ([]byte, error) {
	rawID, err := exec.Command(
		"docker", "run", "-d", "--rm", "band-provider", "sleep", fmt.Sprintf("%d", int(timeOut.Seconds())),
	).Output()
	if err != nil {
		return []byte{}, err
	}
	containerID := strings.TrimSpace(string(rawID))
	defer exec.Command("docker", "stop", containerID).Output()

	_, err = exec.Command(
		"docker", "cp", be.Filename, fmt.Sprintf("%s:/exec.sh", containerID),
	).Output()
	if err != nil {
		return []byte{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	newArgs := append([]string{"exec", containerID, "./exec.sh"}, args...)

	return exec.CommandContext(ctx, "docker", newArgs...).Output()
}

// Close deletes temp file after used.
func (be *Exec) Close() {
	os.Remove(be.Filename)
}

func newExec(filename string) (Exec, error) {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return Exec{}, err
	}
	return Exec{Filename: absolutePath}, nil
}
