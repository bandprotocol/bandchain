package cli

import (
	"bufio"
	"os"
)

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return []byte{}, err
	}

	size := stats.Size()
	output := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(output)

	if err != nil {
		return []byte{}, err
	}

	return output, nil
}
