package component

import "os"

type FileCtrl struct {
	path string
}

func NewFileCtrl(path string) *FileCtrl {
	return &FileCtrl{path: path}
}
func (f *FileCtrl) AppendFile(content string) error {
	path := f.path
	file, err1 := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		return err1
	}
	defer file.Close()
	_, err2 := file.Write([]byte(content))
	if err2 != nil {
		return err2
	}
	return nil
}
func (f *FileCtrl) CoverFile(content string) error {
	path := f.path
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
