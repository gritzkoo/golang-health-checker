package healthcheck

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
	"uuid"

	"github.com/redis/go-redis/v9"
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
	result := Liveness().Status

	if result != internalStatus {
		t.Errorf("HealthCheckerSimple() expect: %s got: %s", internalStatus, result)
	}
}

func TestDetailedChecker(t *testing.T) {
	testLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	tests := []struct {
		name string
		args []IntegrationItem
		want bool
	}{
		// http test block ===============================
		{
			name: "testing http integration with success",
			want: true,
			args: []IntegrationItem{
				NewHttpChecker(HttpCheckerParams{
					Name: "go-test-web1",
					Url:  "https://github.com/status",
					Headers: map[string]string{
						"Accept": "application/json",
					},
				}),
			},
		},
		{
			name: "testing http integration with timeout fail",
			want: false,
			args: []IntegrationItem{
				NewHttpChecker(HttpCheckerParams{
					Name:    "go-test-with timeout",
					Url:     "https://google.com/status",
					Timeout: time.Millisecond,
					Headers: map[string]string{
						"Accept": "application/json",
					},
				}),
			},
		},
		{
			name: "testing http integration with timeout fail but is optional",
			want: true,
			args: []IntegrationItem{
				NewHttpChecker(HttpCheckerParams{
					Name:       "go-test-with timeout",
					IsOptional: true,
					Url:        "https://google.com/status",
					Timeout:    time.Millisecond,
					Headers: map[string]string{
						"Accept": "application/json",
					},
				}),
			},
		},
		{
			name: "testing http integration with request error parse",
			want: false,
			args: []IntegrationItem{
				NewHttpChecker(HttpCheckerParams{
					Name: "go-test-with request fail",
					Url:  "://invalid-url",
				}),
			},
		},
		{
			name: "testing http integration with request error status code 404",
			want: false,
			args: []IntegrationItem{
				NewHttpChecker(HttpCheckerParams{
					Name:    "go-test-with request fail",
					Url:     "https://httpbin.org/status/404",
					Timeout: time.Minute,
					Headers: map[string]string{
						"Postman-Token": uuid.NewV4().String(),
						"Accept":        "application/json",
					},
				}),
			},
		},
		// redis test block ==============================
		{
			name: "testing redis integration with success",
			want: true,
			args: []IntegrationItem{
				NewRedisChecker("go-test-redis", false, &redis.Options{
					Addr: fmt.Sprintf("%s:%d", REDIS_HOST, 6379),
				}),
			},
		},
		{
			name: "testing redis integration with fail",
			want: false,
			args: []IntegrationItem{
				NewRedisChecker("go-test-redis", false, &redis.Options{
					Addr: fmt.Sprintf("%s:%d", REDIS_HOST, 63),
				}),
			},
		},
		{
			name: "testing redis integration with fail but is optional",
			want: true,
			args: []IntegrationItem{
				NewRedisChecker("go-test-redis", true, &redis.Options{
					Addr: fmt.Sprintf("%s:%d", REDIS_HOST, 63),
				}),
			},
		},
		// memcached test block ==============================
		{
			name: "testing memcached integration with success",
			want: true,
			args: []IntegrationItem{
				NewMemcachedChecker("go-test-memcached", false, fmt.Sprintf("%s:%d", MEMCACHE_HOST, 11211)),
			},
		},
		{
			name: "testing memcached integration with fail",
			want: false,
			args: []IntegrationItem{
				NewMemcachedChecker("go-test-memcached", false, fmt.Sprintf("%s:%d", MEMCACHE_HOST, 11)),
			},
		},
		{
			name: "testing memcached integration with fail but is optional",
			want: true,
			args: []IntegrationItem{
				NewMemcachedChecker("go-test-memcached", true, fmt.Sprintf("%s:%d", MEMCACHE_HOST, 11)),
			},
		},
		// custon test block ==============================
		{
			name: "testing custom integration with success",
			want: true,
			args: []IntegrationItem{
				NewCustomChecker("testing-custon-func-success", false, func(ctx context.Context) error {
					return nil
				}),
			},
		},
		{
			name: "testing custom integration with fail",
			want: false,
			args: []IntegrationItem{
				NewCustomChecker("testing-custon-func-success", false, func(ctx context.Context) error {
					return fmt.Errorf("error triggered by testing")
				}),
			},
		},
		{
			name: "testing custom integration with fail but is optional",
			want: true,
			args: []IntegrationItem{
				NewCustomChecker("testing-custon-func-success", true, func(ctx context.Context) error {
					return fmt.Errorf("error triggered by testing")
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RedinessConfig{
				Name:         "testing",
				Version:      "testing",
				Integrations: tt.args,
			}
			got := Readiness(t.Context(), config)

			testLogger.Debug(tt.name, slog.Any("result", got))
			if tt.want != got.Status {
				t.Errorf("HealthCheckerDetailed() Testing: %s expect: %v got: %v", tt.name, tt.want, got.Status)
			}
		})
	}
}
