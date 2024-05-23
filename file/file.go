package file

import (
	"os"
	"path/filepath"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Touch(file string) error {
	if Exist(file) {
		return nil
	}

	if err := MakeDirByFile(file); err != nil {
		return err
	}
	create, err := os.Create(file)
	defer create.Close()
	if err != nil {
		return err
	}
	return nil
}

func OpenOrCreate(file string) (*os.File, error) {
	if err := MakeDirByFile(file); err != nil {
		return nil, err
	}
	return os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
}

func MakeDirByFile(file string) error {
	if !Exist(file) {
		if dir, _ := filepath.Split(file); !Exist(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

func ReWrite(file string) (*os.File, error) {
	if err := MakeDirByFile(file); err != nil {
		return nil, err
	}
	return os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
}

func ReWriteString(file, content string) error {
	fileObj, err := ReWrite(file)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	if _, err = fileObj.WriteString(content); err != nil {
		return err
	}
	return nil
}
