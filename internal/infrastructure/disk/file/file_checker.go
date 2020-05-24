package file

import (
	"io/ioutil"
	"os"
)

type fileChecker struct {
}

func (f *fileChecker) check(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			if err := ioutil.WriteFile(file, nil, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
