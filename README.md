# golang-health-checker

[![test](https://github.com/gritzkoo/golang-health-checker/actions/workflows/test.yaml/badge.svg)](https://github.com/gritzkoo/golang-health-checker/actions/workflows/test.yaml)
[![build](https://github.com/gritzkoo/golang-health-checker/actions/workflows/build.yaml/badge.svg)](https://github.com/gritzkoo/golang-health-checker/actions/workflows/build.yaml)
[![Coverage Status](https://coveralls.io/repos/github/gritzkoo/golang-health-checker/badge.svg?branch=master)](https://coveralls.io/github/gritzkoo/golang-health-checker?branch=master)
![views](https://raw.githubusercontent.com/gritzkoo/golang-health-checker/traffic/traffic-golang-health-checker/views.svg)
![views per week](https://raw.githubusercontent.com/gritzkoo/golang-health-checker/traffic/traffic-golang-health-checker/views_per_week.svg)
![clones](https://raw.githubusercontent.com/gritzkoo/golang-health-checker/traffic/traffic-golang-health-checker/clones.svg)
![clones per week](https://raw.githubusercontent.com/gritzkoo/golang-health-checker/traffic/traffic-golang-health-checker/clones_per_week.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/gritzkoo/golang-health-checker/pkg/healthcheck.svg)](https://pkg.go.dev/github.com/gritzkoo/golang-health-checker/pkg/healthcheck)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gritzkoo/golang-health-checker)
![GitHub repo size](https://img.shields.io/github/repo-size/gritzkoo/golang-health-checker)
![GitHub](https://img.shields.io/github/license/gritzkoo/golang-health-checker)
![GitHub issues](https://img.shields.io/github/issues/gritzkoo/golang-health-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/gritzkoo/golang-health-checker)](https://goreportcard.com/report/github.com/gritzkoo/golang-health-checker)

A lightweight, interface-driven health check library for Go applications. Supports Kubernetes liveness and readiness probes with context cancellation, concurrent checks, and optional dependencies.

## Features

- **Liveness Probe**: Simple "is the app running?" check
- **Readiness Probe**: Detailed check of all integration dependencies
- **Context-Aware**: Full support for timeouts and cancellation
- **Concurrent Checks**: All integrations validated in parallel
- **Optional Dependencies**: Mark integrations as non-critical (readiness doesn't fail if optional dependencies are down)
- **Interface-Driven**: Easy to extend with custom integrations
- **Type-Safe**: No magic strings or config structs—everything is explicit
- **100% Test Coverage**: All code paths covered

## Supported Integrations

- ✅ Redis (via `redis/go-redis/v9`)
- ✅ Memcached (via `bradfitz/gomemcache`)
- ✅ HTTP/HTTPS (any HTTP endpoint)
- ✅ Custom Functions (implement your own checks)

---

## Installation

If you're starting a new Go project:

```sh
go mod init github.com/my/repo
```

Then add this package:

```sh
go get github.com/gritzkoo/golang-health-checker
```

---

## Quick Start

### Basic Usage with Gin

```go
package main

import (
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gritzkoo/golang-health-checker/pkg/healthcheck"
	"github.com/redis/go-redis/v9"
)

func main() {
	router := gin.Default()

	// Liveness Probe: is the app running?
	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, healthcheck.Liveness())
	})

	// Readiness Probe: are all dependencies up?
	router.GET("/health/ready", func(c *gin.Context) {
		config := healthcheck.RedinessConfig{
			Name:    "My App",
			Version: "1.0.0",
			Integrations: []healthcheck.IntegrationItem{
				healthcheck.NewRedisChecker("redis-main", false, &redis.Options{
					Addr: "localhost:6379",
				}),
				healthcheck.NewHttpChecker(healthcheck.HttpCheckerParams{
					Name:    "GitHub Status",
					Url:     "https://github.com/status",
					Timeout: 5 * time.Second,
				}),
				healthcheck.NewCustomChecker("db-migration", false, func(ctx context.Context) error {
					// Custom logic: check database, run queries, etc.
					return checkDatabaseMigration(ctx)
				}),
			},
		}

		response := healthcheck.Readiness(c.Request.Context(), config)
		if response.Status {
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusServiceUnavailable, response)
		}
	})

	router.Run(":8888")
}

func checkDatabaseMigration(ctx context.Context) error {
	// Your custom check logic here
	return nil
}
```

### Response Format

**Liveness Response:**

```json
{
  "status": "fully functional"
}
```

**Readiness Response:**

```json
{
  "name": "My App",
  "status": true,
  "version": "1.0.0",
  "date": "2026-10-12T10:15:30Z",
  "duration": 0.245,
  "integrations": [
    {
      "name": "redis-main",
      "kind": "Redis",
      "status": true,
      "is_optional": false,
      "response_time": 0.0015,
      "url": "localhost:6379",
      "error": ""
    },
    {
      "name": "GitHub Status",
      "kind": "Http",
      "status": true,
      "is_optional": false,
      "response_time": 0.234,
      "url": "https://github.com/status",
      "error": ""
    },
    {
      "name": "db-migration",
      "kind": "Custom Checker",
      "status": true,
      "is_optional": false,
      "response_time": 0.008,
      "url": "",
      "error": ""
    }
  ]
}
```

---

## API Reference

### Functions

#### `Liveness() LivenessResponse`

Returns a simple "fully functional" status. Used for Kubernetes liveness probes to verify the app is running.

#### `Readiness(ctx context.Context, config RedinessConfig) ReadinessResponse`

Performs concurrent checks on all integrations and returns detailed status. Used for Kubernetes readiness probes.

**Parameters:**

- `ctx`: Context for timeouts and cancellation
- `config`: Configuration with app metadata and integrations

---

### Configuration

#### `RedinessConfig`

```go
type RedinessConfig struct {
	Name         string              // Application name
	Version      string              // Application version
	Integrations []IntegrationItem   // List of checks to perform
}
```

---

### Building Integrations

All integrations implement the `IntegrationItem` interface:

```go
type IntegrationItem interface {
	Check(ctx context.Context) (Integration, error)
}
```

#### Redis Checker

```go
checker := healthcheck.NewRedisChecker(
	"redis-cache",           // name
	false,                   // isOptional (set true to not fail readiness if down)
	&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	},
)
```

#### HTTP Checker

```go
checker := healthcheck.NewHttpChecker(healthcheck.HttpCheckerParams{
	Name:       "api-gateway",
	Url:        "https://api.example.com/health",
	Timeout:    10 * time.Second,              // Default: 10 seconds
	Headers:    map[string]string{
		"Authorization": "Bearer token",
	},
	IsOptional: false,
})
```

#### Memcached Checker

```go
checker := healthcheck.NewMemcachedChecker(
	"cache-layer",  // name
	false,          // isOptional
	"localhost:11211",
)
```

#### Custom Checker

```go
checker := healthcheck.NewCustomChecker(
	"database",
	false,  // isOptional
	func(ctx context.Context) error {
		// Custom validation logic
		db := getDatabaseConnection()
		if err := db.PingContext(ctx); err != nil {
			return err
		}
		return nil
	},
)
```

---

## Optional Dependencies

Mark any integration as optional so failures don't fail the overall readiness check:

```go
config := healthcheck.RedinessConfig{
	Integrations: []healthcheck.IntegrationItem{
		// Critical dependency
		healthcheck.NewRedisChecker("redis-cache", false, opts),

		// Optional dependency
		healthcheck.NewHttpChecker(healthcheck.HttpCheckerParams{
			Name:       "Analytics Service",
			Url:        "https://analytics.example.com/health",
			IsOptional: true,  // ← If this fails, readiness still returns true
		}),
	},
}
```

**Result when optional checks fail:**

- The integration is marked `status: false` with an error message
- The root `ReadinessResponse.Status` remains `true`
- HTTP status code: `200 OK` (not `503`)

---

## Kubernetes Integration

### Liveness Probe

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
        - name: app
          image: my-app:latest
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8888
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 2
            failureThreshold: 3
```

### Readiness Probe

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
        - name: app
          image: my-app:latest
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8888
            initialDelaySeconds: 10
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 2
```

---

## Extending with Custom Integrations

Implement the `IntegrationItem` interface to add custom checks:

```go
type DatabaseHealthChecker struct {
	name string
	db   *sql.DB
}

func (d *DatabaseHealthChecker) Check(ctx context.Context) (healthcheck.Integration, error) {
	response := healthcheck.Integration{
		Name:   d.name,
		Kind:   "PostgreSQL",
		Status: true,
		URL:    "database.internal:5432",
	}

	if err := d.db.PingContext(ctx); err != nil {
		return response, err
	}

	return response, nil
}

// Then use it:
checker := &DatabaseHealthChecker{
	name: "postgres-primary",
	db:   db,
}

config.Integrations = append(config.Integrations, checker)
```

---

## Concurrency Model

All integrations are checked **concurrently** using goroutines:

```
Start
  ├─ Check Redis (goroutine 1)      ~1ms
  ├─ Check HTTP (goroutine 2)       ~250ms
  ├─ Check Memcached (goroutine 3)  ~2ms
  └─ Check Custom (goroutine 4)     ~50ms

  Total Time: ~250ms (max of all, not sum)
  vs. Sequential: ~303ms
```

Each goroutine:

- Respects the parent context timeout
- Reports individual response times
- Can fail independently (optional integrations don't fail the whole response)

---

## Error Handling

Errors are captured per-integration and included in the response:

```json
{
  "name": "redis-cache",
  "kind": "Redis",
  "status": false,
  "error": "dial tcp localhost:6379: connect: connection refused",
  "response_time": 0.015
}
```

The main readiness status is `false` only if a **non-optional** integration fails.

---

## Testing

Run the full test suite:

```sh
make test
```

With race condition detection:

```sh
go test -race ./...
```

View coverage:

```sh
make coverage
```

All code paths are covered (100% coverage).

---

## Performance Notes

- Goroutine overhead is minimal (~2KB per check)
- Connections are reused (Redis client persists across checks)
- HTTP response bodies are properly closed (no file descriptor leaks)
- Context cancellation works at the integration level
- Typical health check (4 integrations): **250-500ms**

---

## Migration from v1

If you're upgrading from v1, here are the key changes:

| v1                               | v2                                          |
| -------------------------------- | ------------------------------------------- |
| `HealthCheckerSimple()`          | `Liveness()`                                |
| `HealthCheckerDetailed(config)`  | `Readiness(ctx, config)`                    |
| `ApplicationConfig`              | `RedinessConfig`                            |
| `ApplicationHealthDetailed`      | `ReadinessResponse`                         |
| `ApplicationHealthSimple`        | `LivenessResponse`                          |
| Config struct with magic strings | Interface-driven (type-safe)                |
| No context support               | Full context support                        |
| No optional dependencies         | `IsOptional` field                          |
| Manual checker instantiation     | Builder functions: `NewRedisChecker()` etc. |

**Old API:**

```go
config := ApplicationConfig{
	Integrations: []IntegrationConfig{
		{
			Type: "redis",
			Name: "cache",
			Host: "localhost",
			Port: "6379",
		},
	},
}
result := HealthCheckerDetailed(config)
```

**New API:**

```go
config := RedinessConfig{
	Integrations: []IntegrationItem{
		NewRedisChecker("cache", false, &redis.Options{
			Addr: "localhost:6379",
		}),
	},
}
result := Readiness(ctx, config)
```

---

## Lightweight Version

Looking for a version with zero external dependencies? Check out [golang-health-checker-lw](https://github.com/gritzkoo/golang-health-checker-lw).

---

## Contributing

Issues and PRs welcome. Please ensure:

- All tests pass: `go test -race ./...`
- Code coverage remains at 100%
- Go code follows `gofmt` standards

---

## License

MIT License. See `LICENSE` file for details.

---

## Support

- **Bugs:** [GitHub Issues](https://github.com/gritzkoo/golang-health-checker/issues)
- **Questions:** Open a discussion or issue
- **Docs:** [pkg.go.dev](https://pkg.go.dev/github.com/gritzkoo/golang-health-checker/pkg/healthcheck)
