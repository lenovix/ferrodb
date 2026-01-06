package persistence

import (
	"bufio"
	"os"
)

type AOF struct {
	file *os.File
}

func OpenAOF(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &AOF{file: file}, nil
}

func (a *AOF) Write(command string) error {
	_, err := a.file.WriteString(command + "\n")
	return err
}

func (a *AOF) Replay(apply func(string)) error {
	file, err := os.Open(a.file.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		apply(scanner.Text())
	}
	return scanner.Err()
}

func (a *AOF) Close() error {
	return a.file.Close()
}
