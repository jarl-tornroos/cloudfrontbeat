package workers

import (
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
	"fmt"
)

// GeoLogData contains geo data for one ip address
type GeoLogData struct {
	CountryCode   string  `workers:"geoip.country_code"`
	Country       string  `workers:"geoip.country"`
	Region        string  `workers:"geoip.region"`
	City          string  `workers:"geoip.city"`
	ContinentCode string  `workers:"geoip.continent_code"`
	Continent     string  `workers:"geoip.continent"`
	Latitude      float64 `workers:"geoip.latitude"`
	Longitude     float64 `workers:"geoip.longitude"`
	Location      string  `workers:"geoip.location"`
}

// GeoLog
type GeoLog struct {
	provider cflib.Geo
}

// SetIP for the address we want geo data from
func (m *GeoLog) SetIP(ip string) error {
	return m.provider.SetIP(ip)
}

// GetGeoData return formatted geo information
func (m *GeoLog) GetGeoData() GeoLogData {

	latitude := m.provider.Latitude()
	longitude := m.provider.Longitude()

	return GeoLogData{
		CountryCode:   m.provider.CountryCode(),
		Country:       m.provider.CountryName(),
		Region:        m.provider.Region(),
		City:          m.provider.City(),
		ContinentCode: m.provider.ContinentCode(),
		Continent:     m.provider.Continent(),
		Latitude:      latitude,
		Longitude:     longitude,
		Location: fmt.Sprintf("%f,%f", latitude, longitude),
	}
}
