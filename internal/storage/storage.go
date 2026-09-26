package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"uuid"
)

const (
	FilesDirPath   = "files/"
	codeAccessDir  = 0700
	codeAccessFile = 0600
)

type Files struct{}

type Saver interface {
	Save(fileData []byte) (string, error)
}

type Cleaner interface {
	RemoveOld() error
}

func NewStorage() *Files {
	return &Files{}
}

func (s *Files) Save(fileData []byte) (string, error) {
	time := time.Now()
	filepathDir := fmt.Sprintf("%v/%v/%v/%v/%v/", FilesDirPath, time.Year(), int(time.Month()), time.Day(), time.Hour())

	if err := os.MkdirAll(filepathDir, codeAccessDir); err != nil {
		return "", err
	}

	filepath := filepathDir + uuid.New().String()
	if err := os.WriteFile(filepath, fileData, codeAccessFile); err != nil {
		return "", err
	}

	return filepath, nil
}

func (s *Files) RemoveOld() error {
	root := FilesDirPath
	removeTime := time.Now().Add(-15 * 24 * time.Hour)
	removeYear := removeTime.Year()
	removeMonth := int(removeTime.Month())
	removeDay := removeTime.Day()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		splittedDirPath := strings.Split(d.Name(), "/")

		year, err := strconv.Atoi(splittedDirPath[1])
		if err != nil {
			return err
		}
		if year > removeYear {
			if err := os.RemoveAll(d.Name()); err != nil {
				return err
			}

			return filepath.SkipDir
		}

		month, err := strconv.Atoi(splittedDirPath[2])
		if err != nil {
			return err
		}
		if month > removeMonth {
			if err := os.RemoveAll(d.Name()); err != nil {
				return err
			}

			return filepath.SkipDir
		}

		day, err := strconv.Atoi(splittedDirPath[3])
		if err != nil {
			return err
		}
		if day > removeDay {
			if err := os.RemoveAll(d.Name()); err != nil {
				return err
			}

			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
