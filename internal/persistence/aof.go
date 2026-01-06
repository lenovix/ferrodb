package persistence

import (
	"bufio"
	"fmt"
	"os"

	"ferrodb/internal/storage"
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

func (a *AOF) Rewrite(snapshot map[string]storage.Item) error {
	tmpPath := a.file.Name() + ".tmp"

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(tmpFile)

	for key, item := range snapshot {
		_, err := writer.WriteString(
			fmt.Sprintf("SET %s %s\n", key, item.Value),
		)
		if err != nil {
			return err
		}

		if item.ExpireAt > 0 {
			_, err = writer.WriteString(
				fmt.Sprintf("EXPIREAT %s %d\n", key, item.ExpireAt),
			)
			if err != nil {
				return err
			}
		}
	}

	writer.Flush()
	tmpFile.Sync()
	tmpFile.Close()

	a.file.Close()

	err = os.Rename(tmpPath, a.file.Name())
	if err != nil {
		return err
	}

	a.file, err = os.OpenFile(a.file.Name(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	return err
}
