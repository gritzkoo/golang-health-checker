package healthcheck

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSimpleChecker(t *testing.T) {
	result := HealthCheckerSimple().Status

	assert.Equal(t, result, "fully functional")
}

type detailedListProvider struct {
	Type     string
	Host     string
	Port     string
	Expected bool
}

var detailedDataProvider = []detailedListProvider{
	{
		Type:     Redis,
		Host:     "redis",
		Port:     "6379",
		Expected: true,
	}, {
		Type:     Redis,
		Host:     "redis",
		Port:     "63",
		Expected: false,
	}, {
		Type:     Memcached,
		Host:     "memcache",
		Port:     "11211",
		Expected: true,
	}, {
		Type:     Memcached,
		Host:     "memcache",
		Port:     "11",
		Expected: false,
	}, {
		Type:     Web,
		Host:     "https://github.com/status",
		Expected: true,
	}, {
		Type:     Web,
		Host:     "https://google.com/status",
		Expected: false,
	},
}

func TestDetailedChecker(t *testing.T) {

	for _, v := range detailedDataProvider {
		config := ApplicationConfig{
			Name:    "test",
			Version: "test",
			Integrations: []IntegrationConfig{
				{
					Type: v.Type,
					Host: v.Host,
					Port: v.Port,
				},
			},
		}

		result := HealthCheckerDetailed(config)

		assert.IsEqual(result.Status, v)
	}
}
