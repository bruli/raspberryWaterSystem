package disk

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func readYamlFile(path string, data interface{}) error {
	if err := checkFile(path); err != nil {
		return fmt.Errorf("failed to check %s file", path)
	}
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if err = yaml.Unmarshal(fileData, data); err != nil {
		return err
	}
	return nil
}

func readJsonFile(path string, data interface{}) error {
	if err := checkFile(path); err != nil {
		return fmt.Errorf("failed to check %s file", path)
	}
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if err = json.Unmarshal(fileData, data); err != nil {
		return err
	}
	return nil
}

func checkFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.WriteFile(path, nil, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeYamlFile(path string, data interface{}) error {
	dataFile, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return os.WriteFile(path, dataFile, 0o755)
}

func writeJsonFile(path string, data interface{}) error {
	dataFile, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return os.WriteFile(path, dataFile, 0o755)
}
