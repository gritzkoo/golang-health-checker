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
	Expected bool
	Config   IntegrationConfig
}

var detailedDataProvider = []detailedListProvider{
	{
		Expected: true,
		Config: IntegrationConfig{
			Type: Redis,
			Name: "go-test-redis",
			DB:   1,
			Host: "localhost",
			Port: "6379",
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Redis,
			Name: "go-test-redis",
			DB:   1,
			Host: "localhost",
			Port: "63",
		},
	}, {
		Expected: true,
		Config: IntegrationConfig{
			Type: Memcached,
			Name: "go-test-memcached",
			Host: "localhost",
			Port: "11211",
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Memcached,
			Name: "go-test-memcached",
			Host: "localhost",
			Port: "11",
		},
	}, {
		Expected: true,
		Config: IntegrationConfig{
			Type: Web,
			Name: "go-test-web",
			Host: "https://github.com/status",
			Headers: []HTTPHeader{
				{
					Key:   "Accept",
					Value: "application/json",
				},
			},
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Web,
			Name: "go-test-web",
			Host: "#@$*&",
			Headers: []HTTPHeader{
				{
					Key:   "Accept",
					Value: "application/json",
				},
			},
			TimeOut: 1000,
		},
	},
}

func TestDetailedChecker(t *testing.T) {

	for _, v := range detailedDataProvider {
		config := ApplicationConfig{
			Name:    "test",
			Version: "test",
			Integrations: []IntegrationConfig{
				v.Config,
			},
		}

		result := HealthCheckerDetailed(config)

		assert.IsEqual(result.Status, v.Expected)
	}
}
