package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadFileAsString(filename string) (string, error) {
	fullPath := filepath.Join("/Users/k.song/src/golang/bigip-ctlr/test/configmaps", filename)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %v", err)
		}
		return "", fmt.Errorf("error checking file: %v", err)
	}
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(data), nil
}
