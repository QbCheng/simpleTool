package simpleImageStitching

import (
	"errors"
	"os"
)

func SaveFile(fileName string, data []string) (*os.File, error) {
	fs, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	for i := range data {
		n, err := fs.WriteString(data[i])
		if err != nil {
			return nil, err
		}
		if n != len(data[i]) {
			return nil, errors.New(" write length diff")
		}
		_, err = fs.WriteString("\n")
		if err != nil {
			return nil, err
		}
	}
	err = fs.Sync()
	if err != nil {
		return nil, err
	}
	return fs, nil
}
