package file

type repository struct {
	writer writer
	reader reader
}

func newRepository(file string) *repository {
	return &repository{writer: newFileWriter(file), reader: newFileReader(file)}
}
