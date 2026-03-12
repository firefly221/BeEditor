package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path       string
	IsModified bool
}

func (f *File) LoadFile() ([][]rune, error) {
	if strings.TrimSpace(f.Path) == "" {
		return nil, errors.New("empty path")
	}

	data, err := os.ReadFile(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			f.IsModified = false
			return [][]rune{{}}, nil
		}
		return nil, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")

	parts := strings.Split(text, "\n")

	lines := make([][]rune, 0, len(parts))
	for _, p := range parts {
		lines = append(lines, []rune(p))
	}

	f.IsModified = false
	return lines, nil
}

func (f *File) SaveFile(lines [][]rune) error {
	if strings.TrimSpace(f.Path) == "" {
		return errors.New("empty path")
	}

	dir := filepath.Dir(f.Path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	var b strings.Builder
	for i, line := range lines {
		b.WriteString(string(line))
		if i != len(lines)-1 {
			b.WriteByte('\n')
		}
	}

	if err := os.WriteFile(f.Path, []byte(b.String()), 0o644); err != nil {
		return err
	}

	f.IsModified = false
	return nil
}

func (f *File) CreateFile() error {
	if strings.TrimSpace(f.Path) == "" {
		return errors.New("empty path")
	}

	dir := filepath.Dir(f.Path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	fd, err := os.OpenFile(f.Path, os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return fd.Close()
}
