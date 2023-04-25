package gpio

import (
	"os"
)

type pin struct {
	name string
}

func newPin(pinNumber string) *pin {
	g := &pin{name: pinNumber}
	g.build()
	return g
}

func (r *pin) build() {
	filename := r.filename()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// export gpio pin
		_ = os.WriteFile("/sys/class/gpio/export", []byte(r.name), 0o666)
	}
}

func (r *pin) filename() string {
	return "/sys/class/gpio/gpio" + r.name
}

func (r *pin) write(where, what string) *pin {
	filename := r.filename() + "/" + where
	_ = os.WriteFile(filename, []byte(what), 0o666)
	return r
}

func (r *pin) output() *pin {
	return r.write("direction", "out")
}

func (r *pin) high() *pin {
	return r.write("value", "1")
}

func (r *pin) low() *pin {
	return r.write("value", "0")
}
