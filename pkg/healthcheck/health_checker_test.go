package healthcheck

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func getenv(key string, fallback string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return fallback
}

var (
	REDIS_HOST    = getenv("REDIS_HOST", "localhost")
	MEMCACHE_HOST = getenv("MEMCACHE_HOST", "localhost")
)

func TestSimpleChecker(t *testing.T) {
	result := HealthCheckerSimple().Status

	assert.Equal(t, result, "fully functional")
}

type detailedListProvider struct {
	Expected bool
	FakeHttp bool
	Config   IntegrationConfig
}

var detailedDataProvider = []detailedListProvider{
	{
		Expected: true,
		Config: IntegrationConfig{
			Type: Redis,
			Name: "go-test-redis",
			DB:   1,
			Host: REDIS_HOST,
			Port: "6379",
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Redis,
			Name: "go-test-redis",
			DB:   1,
			Host: REDIS_HOST,
			Port: "63",
		},
	}, {
		Expected: true,
		Config: IntegrationConfig{
			Type: Memcached,
			Name: "go-test-memcached",
			Host: MEMCACHE_HOST,
			Port: "11211",
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Memcached,
			Name: "go-test-memcached",
			Host: MEMCACHE_HOST,
			Port: "11",
		},
	}, {
		Expected: true,
		Config: IntegrationConfig{
			Type: Web,
			Name: "go-test-web1",
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
			Type:    Web,
			Name:    "go-test-with no 200",
			Host:    "https://google.com/status",
			TimeOut: 1,
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type:    Web,
			Name:    "go-test-with-error",
			Host:    "tcp://jsfiddle.net",
			TimeOut: 1,
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: "unknow",
			Name: "go-test-unknow",
		},
	}, {
		Expected: true,
		Config: IntegrationConfig{
			Type: Custom,
			Name: "testing-custom-func-success",
			Handle: func() error {
				return nil
			},
		},
	}, {
		Expected: false,
		Config: IntegrationConfig{
			Type: Custom,
			Name: "testing-custom-func-with-error",
			Handle: func() error {
				return fmt.Errorf("error triggered by testing")
			},
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
		condition := v.Expected == result.Status
		printstring := "ok"
		if !condition {
			printstring = "nok"
		}
		fmt.Println("Running config:", v.Config.Name, " and result: ", printstring)
		assert.Equal(t, result.Status, v.Expected)
	}
}
