package worker

import "fmt"

type ReadCurrentStatusError struct {
	err error
}

func (r ReadCurrentStatusError) Error() string {
	return fmt.Sprintf("failed reading current status: %q", r.err.Error())
}

type FindProgramsError struct {
	err error
}

func (f FindProgramsError) Error() string {
	return fmt.Sprintf("failed finding programs: %q", f.err.Error())
}

type ExecuteDailyError struct {
	err error
}

func (e ExecuteDailyError) Error() string {
	return fmt.Sprintf("failed executing daily programs: %q", e.err.Error())
}

type ExecuteOddEvenError struct {
	err error
}

func (e ExecuteOddEvenError) Error() string {
	return fmt.Sprintf("failed executing odd/even programs: %q", e.err.Error())
}

type ExecuteWeeklyError struct {
	err error
}

func (e ExecuteWeeklyError) Error() string {
	return fmt.Sprintf("failed executing weekly programs: %q", e.err.Error())
}

type ExecuteTemperatureError struct {
	err error
}

func (e ExecuteTemperatureError) Error() string {
	return fmt.Sprintf("failed executing temperature programs: %q", e.err.Error())
}
