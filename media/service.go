package media

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Service interface {
	PostMedia(string, io.Reader) error

	// I don't know if returning byte slice is a best way.
	GetMedia(string) ([]byte, error)
	DeleteMedia(string) error
}

type ServiceMiddleware func(Service) Service

type SimpleDirectoryService struct {
	dirName string
}

func (s SimpleDirectoryService) DeleteMedia(fileName string) error {
	filepath := filepath.Join(s.dirName, fileName)
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return ErrNotExist
	}
	return os.Remove(filepath)
}

func (s SimpleDirectoryService) PostMedia(fileName string, reader io.Reader) error {
	filePath := filepath.Join(s.dirName, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return ErrFileExists
	}
	writeFile, err := os.Create(filePath)
	if err != nil {
		return ErrIO
	}
	defer writeFile.Close()

	_, err = io.Copy(writeFile, reader)
	if err != nil {
		return ErrCopy
	}
	return nil
}

func (s SimpleDirectoryService) GetMedia(fileName string) ([]byte, error) {
	b, e := os.ReadFile(filepath.Join(s.dirName, fileName))
	if e != nil {
		return nil, ErrIO
	}
	return b, nil
}

func NewSimpleStatelessService() Service {
	return &SimpleDirectoryService{
		dirName: "files",
	}
}

var (
	ErrFileExists = errors.New("file already exists")
	ErrCopy       = errors.New("copy file error")
	ErrIO         = errors.New("error in file IO operation")
	ErrNotExist   = errors.New("no such file exist")
)
