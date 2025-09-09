package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := info.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	remain := fileSize - offset
	toCopy := remain
	if limit > 0 && limit < remain {
		toCopy = limit
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	bufSize := 1 * 1024
	buf := make([]byte, bufSize)
	var copied int64

	for copied < toCopy {
		readSize := bufSize
		left := toCopy - copied
		if left < int64(readSize) {
			readSize = int(left)
		}
		n, readErr := file.Read(buf[:readSize])
		if n > 0 {
			written, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			if written != n {
				return io.ErrShortWrite
			}
			copied += int64(written)

			printProgress(copied, toCopy)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	println()
	return nil
}

func printProgress(done, total int64) {
	if total == 0 {
		print("\r[----------] 100%")
		return
	}
	barWidth := 30
	pct := float64(done) / float64(total)
	filled := int(pct * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	bar := "[" + string(repeat('=', filled)) + ">" + string(repeat(' ', barWidth-filled-1)) + "]"
	percent := int(pct * 100)
	if percent > 100 {
		percent = 100
	}
	print("\r", bar, " ", percent, "%")
}

func repeat(char rune, count int) []rune {
	if count <= 0 {
		return []rune{}
	}
	s := make([]rune, count)
	for i := range s {
		s[i] = char
	}
	return s
}
