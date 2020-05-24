package execution

type InvalidCreateData struct {
	error string
}

func (e InvalidCreateData) Error() string {
	return e.error
}

func NewInvalidCreateData(error string) InvalidCreateData {
	return InvalidCreateData{error: error}
}

type InvalidCreateExecution struct {
	error string
}

func NewInvalidCreateExecution(error string) *InvalidCreateExecution {
	return &InvalidCreateExecution{error: error}
}

func (i InvalidCreateExecution) Error() string {
	return i.error
}
