// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// GetSize returns the size of a file
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// GetExt returns the file extension
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist checks if a file/dir exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission checks if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistMkDir checks if directory does not exist and creates it
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// MkDir creates a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Open opens a file with flags and specific permission
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// FormatFileSize formats byte size to human readable string
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		// return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2f B", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f TB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { // if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2f EB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
