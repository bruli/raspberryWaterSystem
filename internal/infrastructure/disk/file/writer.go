package file

import "io/ioutil"

//go:generate moq -out writer_mock.go . writer
type writer interface {
	write(d []byte) error
}

type fileWriter struct {
	filePath string
	fileChecker
}

func (f *fileWriter) write(d []byte) error {
	if err := f.fileChecker.check(f.filePath); err != nil {
		return err
	}
	return ioutil.WriteFile(f.filePath, d, 0755)

}

func newFileWriter(filePath string) writer {
	return &fileWriter{filePath: filePath}
}
