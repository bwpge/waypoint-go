package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	p := strings.Replace(path, "~/", "", 1)
	return filepath.Join(home, p)
}

func checkFile(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("file does not exist")
		} else {
			return fmt.Errorf("error reading file: %s", err.Error())
		}
	}
	if !stat.Mode().IsRegular() {
		return errors.New("path is not a file")
	}

	return nil
}

func fetch(url string) (string, error) {
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
