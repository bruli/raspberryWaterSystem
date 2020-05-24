package server

type executionData struct {
	seconds uint8
	zone    string
}

func newExecutionData(seconds uint8, zone string) *executionData {
	return &executionData{seconds: seconds, zone: zone}
}
