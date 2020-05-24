package error

import (
	"strings"
)

type Aggregated struct {
	messages []string
}

func (e *Aggregated) Add(m string) {
	e.messages = append(e.messages, m)
}

func (e *Aggregated) Error() string {
	return strings.Join(e.messages, ",")
}

func (e *Aggregated) WithErrors() bool {
	return 0 != len(e.messages)
}
