package relay

type relays map[string]string

func getRelays() *relays {
	r := relays{}
	r["1"] = "18"
	r["2"] = "24"
	r["3"] = "23"
	r["4"] = "25"
	r["5"] = "17"
	r["6"] = "27"

	return &r
}

func GetPins() []string {
	var p []string
	for pin := range *getRelays() {
		p = append(p, pin)
	}

	return p
}
func (r relays) getPin(i string) string {
	return r[i]
}
