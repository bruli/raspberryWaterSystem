package rain

import "fmt"

type Rain struct {
	isRain bool
	value  uint16
}

func New(isRain bool, value uint16) Rain {
	return Rain{isRain: isRain, value: value}
}

func (r Rain) IsRain() bool {
	return r.isRain
}

func (r Rain) Value() uint16 {
	return r.value
}

type Reader struct {
	rep Repository
}

func NewReader(rep Repository) *Reader {
	return &Reader{rep: rep}
}

func (r *Reader) Read() (Rain, error) {
	rai, err := r.rep.Get()
	if err != nil {
		return Rain{}, fmt.Errorf("failed to read rain data: %w", err)
	}

	return rai, err
}
