package cflib

// Geo interface for geo providers
type Geo interface {
	CloseDb()
	SetIP(string) error
	CountryCode() string
	CountryName() string
	Region() string
	City() string
	Latitude() float64
	Longitude() float64
	ContinentCode() string
	Continent() string
}
