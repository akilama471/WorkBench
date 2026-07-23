package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func OpenLogFile(logsDir string) (*os.File, error) {
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory %s: %w", logsDir, err)
	}

	logPath := filepath.Join(logsDir, "devbox.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", logPath, err)
	}

	return f, nil
}

func ReadLastLines(path string, maxLines int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > maxLines*2 {
			lines = lines[len(lines)-maxLines:]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read log file %s: %w", path, err)
	}

	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	return lines, nil
}

func ReadLastBytes(path string, maxBytes int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	size := stat.Size()
	if size == 0 {
		return nil, nil
	}

	offset := size - maxBytes
	if offset < 0 {
		offset = 0
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek file %s: %w", path, err)
	}

	data := make([]byte, size-offset)
	n, err := f.Read(data)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return data[:n], nil
}
