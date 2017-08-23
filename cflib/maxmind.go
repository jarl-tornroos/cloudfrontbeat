package cflib

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"regexp"
	"errors"
)

// MaxMind geo provider
type MaxMind struct {
	db *geoip2.Reader
	record *geoip2.City
}

// Construct function for MaxMind data base.
// MaxMind database file is given as argument
func NewMaxMind(dbFile string) (*MaxMind, error) {
	var err error
	geo := &MaxMind{}
	geo.db, err = geoip2.Open(dbFile)
	return geo, err
}

// CloseDb
func (g *MaxMind) CloseDb()  {
	g.db.Close()
}

// SetIP for the address we want geo data from
func (g *MaxMind) SetIP(ip string) error {
	var err error

	// Check that we provided a real IP address
	realIp, _ := regexp.MatchString(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`, ip)
	if realIp {
		parsedIp := net.ParseIP(ip)
		g.record, err = g.db.City(parsedIp)
		return err
	} else {
		return errors.New(ip + " is not an IP address")
	}
}

func (g *MaxMind) CountryCode() string {
	return g.record.Country.IsoCode
}

func (g *MaxMind) CountryName() string {
	return g.record.Country.Names["en"]
}

func (g *MaxMind) Region() string {
	if len(g.record.Subdivisions) > 0 {
		return g.record.Subdivisions[0].Names["en"]
	}
	return ""
}

func (g *MaxMind) City() string {
	return g.record.City.Names["en"]
}

func (g *MaxMind) Latitude() float64 {
	return g.record.Location.Latitude
}

func (g *MaxMind) Longitude() float64 {
	return g.record.Location.Longitude
}

func (g *MaxMind) ContinentCode() string {
	return g.record.Continent.Code
}

func (g *MaxMind) Continent() string {
	return g.record.Continent.Names["en"]
}
