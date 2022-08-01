package disk

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func readFile(path string, data interface{}) error {
	if err := checkFile(path); err != nil {
		return fmt.Errorf("failed to check %s file", path)
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if err = yaml.Unmarshal(fileData, data); err != nil {
		return err
	}
	return nil
}

func checkFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = ioutil.WriteFile(path, nil, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeFile(path string, data interface{}) error {
	dataFile, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return ioutil.WriteFile(path, dataFile, 0o755)
}
