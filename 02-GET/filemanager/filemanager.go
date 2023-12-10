package filemanager

import "os"

type FileManager struct {
	FilePath string
}

func NewFileManager(path string) *FileManager {
	return &FileManager{FilePath: path}
}

func (fm *FileManager) ReadFile() ([]byte, error) {
	data, err := os.ReadFile(fm.FilePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
