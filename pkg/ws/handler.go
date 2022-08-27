package ws

import "context"

type (
	StatusFunc      func(ctx context.Context) (Status, error)
	WeatherFunc     func(ctx context.Context) (Weather, error)
	LogsFunc        func(ctx context.Context, number int) ([]Log, error)
	ExecuteZoneFunc func(ctx context.Context, zone string, seconds int) error
	Handlers        struct {
		GetStatus   StatusFunc
		Weather     WeatherFunc
		Logs        LogsFunc
		ExecuteZone ExecuteZoneFunc
	}
)
