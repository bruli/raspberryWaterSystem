package relay

type ZoneRelayRepository struct {
	r *relays
}

func (z ZoneRelayRepository) Get() []string {
	var p []string
	for pin := range *z.r {
		p = append(p, pin)
	}

	return p
}

func NewZoneRelayRepository() *ZoneRelayRepository {
	return &ZoneRelayRepository{r: getRelays()}
}
