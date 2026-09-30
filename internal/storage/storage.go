package storage

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	TimeToLiveFilesDirPath = "users/"
	ForeversDirPath        = "avatars/"
	codeAccessDir          = 0700
	codeAccessFile         = 0600
)

type File struct {
	storagePath string
}

type Storage interface {
	SaveFile(fileData []byte, saveAvatar bool) (string, error)
	DelFile(storagePath string) error
}

type Cleaner interface {
	RemoveOld() error
}

func NewStorage(storagePath string) *File {
	return &File{
		storagePath: storagePath,
	}
}

func (s *File) usersDir() string   { return filepath.Join(s.storagePath, TimeToLiveFilesDirPath) }
func (s *File) avatarsDir() string { return filepath.Join(s.storagePath, ForeversDirPath) }

func (s *File) SaveFile(fileData []byte, saveForever bool) (string, error) {
	time := time.Now()

	fileDirPathToSave := s.usersDir()
	if saveForever {
		fileDirPathToSave = s.avatarsDir()
	}

	filepathDir := fmt.Sprintf("%v/%v/%v/%v/%v/", fileDirPathToSave, time.Year(), int(time.Month()), time.Day(), time.Hour())

	if err := os.MkdirAll(filepathDir, codeAccessDir); err != nil {
		return "", err
	}

	filepath := filepathDir + uuid.New().String()
	if err := os.WriteFile(filepath, fileData, codeAccessFile); err != nil {
		return "", err
	}

	return filepath, nil
}

func (s *File) DelFile(storagePath string) error {
	return os.Remove(storagePath)
}

func (s *File) RemoveOld() error {
	defer func() {
		if r := recover(); r != nil {
			log.Default().Printf("Recovered in RemoveOld: %v", r)
		}
	}()

	cutoff := time.Now().Add(-15 * 24 * time.Hour)
	root := s.usersDir()

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return nil
		}
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) < 3 {
			return nil
		}

		year, err1 := strconv.Atoi(parts[0])
		month, err2 := strconv.Atoi(parts[1])
		day, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil
		}

		dirTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if dirTime.Before(cutoff) {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return nil
	})
}
