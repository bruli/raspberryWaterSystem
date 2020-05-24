package temperature

import (
	"github.com/d2r2/go-dht"
	"github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("dht",
	logger.InfoLevel,
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Read() (temp float32, hum float32, err error) {
	defer logger.FinalizeLogger()
	sensorType := dht.DHT11
	pin := 4
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	if err != nil {
		return 0.00, 0.00, err
	}
	return temperature, humidity, nil
}
