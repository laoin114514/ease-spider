package component

import "os"

func AppendFile(path, content string) error {
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
func CoverFile(path, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
