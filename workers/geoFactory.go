package workers

import (
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
	"github.com/jarl-tornroos/cloudfrontbeat/config"
	"fmt"
)

// GetGeoProvider is a factory for geo providers
func GetGeoProvider(config *config.Config) (cflib.Geo, error) {
	var err error
	var geo cflib.Geo

	switch config.GeoManager {
	case "maxmind":
		geo, err = cflib.NewMaxMind(config.MaxMindDb)
	default:
		err = fmt.Errorf("No geo provider %s implemented", config.GeoManager)
	}

	return geo, err
}
