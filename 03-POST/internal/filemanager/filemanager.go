package filemanager

import "os"

type FileManager struct {
	FilePath string
	Data     []byte
}

func NewFileManager(path string) *FileManager {
	return &FileManager{FilePath: path, Data: nil}
}

func (fm *FileManager) ReadFile() error {
	data, err := os.ReadFile(fm.FilePath)
	if err != nil {
		return err
	}

	fm.Data = data

	return nil
}
