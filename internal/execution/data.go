package execution

import (
	"fmt"
	"time"
)

type Data struct {
	Hour  time.Time
	Zones []string
}

func NewData(hour time.Time, zones []string) (*Data, error) {
	if len(zones) == 0 {
		return nil, NewInvalidCreateData("zones cannot be empty")
	}
	for _, z := range zones {
		if "" == z {
			return nil, NewInvalidCreateData(fmt.Sprintf("invalid zone data: %s", z))
		}
	}
	return &Data{Hour: hour, Zones: zones}, nil
}
