package zone

//go:generate moq -out repository_mock.go . Repository
type Repository interface {
	GetZones() *Zones
	Save(z Zones) error
	Find(id string) *Zone
}
