package relay

import (
	"io/ioutil"
	"os"
)

type gioPins []*GpioPin

func (g *gioPins) add(gp *GpioPin) {
	*g = append(*g, gp)
}

func (g *gioPins) getPin(p string) *GpioPin {
	for _, j := range *g {
		if p == j.Name {
			return j
		}
	}
	return nil
}

func newGioPins() *gioPins {
	g := gioPins{}
	for _, j := range *getRelays() {
		g.add(NewGpioPin(j))
	}
	return &g
}

type GpioPin struct {
	Name string
}

func NewGpioPin(name string) *GpioPin {
	g := &GpioPin{Name: name}
	g.build()
	return g
}

func (r *GpioPin) build() {
	filename := r.filename()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// export gpio pin
		_ = ioutil.WriteFile("/sys/class/gpio/export", []byte(r.Name), 0666)
	}
}

func (r *GpioPin) filename() string {
	return "/sys/class/gpio/gpio" + r.Name
}
func (r *GpioPin) write(where, what string) *GpioPin {
	filename := r.filename() + "/" + where
	_ = ioutil.WriteFile(filename, []byte(what), 0666)
	return r
}
func (r *GpioPin) Output() *GpioPin {
	return r.write("direction", "out")
}
func (r *GpioPin) High() *GpioPin {
	return r.write("value", "1")
}
func (r *GpioPin) Low() *GpioPin {
	return r.write("value", "0")
}
