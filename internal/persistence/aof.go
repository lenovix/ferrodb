package persistence

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"ferrodb/internal/storage"
)

type AOF struct {
	file *os.File
}

func OpenAOF(path string) (*AOF, error) {
	// pastikan directory ada
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

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
		line := scanner.Text()
		if line == "" {
			continue
		}
		apply(line)
	}
	return scanner.Err()
}

func (a *AOF) Rewrite(snapshot map[int]map[string]storage.Item) error {
	tmpPath := a.file.Name() + ".tmp"

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(tmpFile)

	for db, kv := range snapshot {
		for key, item := range kv {
			// SET
			if _, err := writer.WriteString(
				fmt.Sprintf("SET %d %s %s\n", db, key, item.Value),
			); err != nil {
				return err
			}

			// EXPIREAT (absolute)
			if item.ExpireAt > 0 {
				if _, err := writer.WriteString(
					fmt.Sprintf("EXPIREAT %d %s %d\n", db, key, item.ExpireAt),
				); err != nil {
					return err
				}
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	// atomic replace
	oldPath := a.file.Name()
	if err := a.file.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, oldPath); err != nil {
		return err
	}

	a.file, err = os.OpenFile(oldPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	return err
}

func (a *AOF) Sync() error {
	return a.file.Sync()
}

func (a *AOF) Close() error {
	return a.file.Close()
}
