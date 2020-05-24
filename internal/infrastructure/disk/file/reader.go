package file

import "io/ioutil"

//go:generate moq -out reader_mock.go . reader
type reader interface {
	read() ([]byte, error)
}

type fileReader struct {
	filePath string
	fileChecker
}

func (f *fileReader) read() ([]byte, error) {
	if err := f.fileChecker.check(f.filePath); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(f.filePath)
}

func newFileReader(filePath string) reader {
	return &fileReader{filePath: filePath}
}
