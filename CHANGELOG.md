# Golang Health Checker CHANGE LOG

## V2.0.0 2026-10-12

**BREAKING CHANGES**

🚨 This package has several changes and you need to take some actions before update to this new version 🚨

### Functions and structs changes

| kind     | Before                                          | After                                                 | effects              |
| -------- | ----------------------------------------------- | ----------------------------------------------------- | -------------------- |
| `func`   | HealthCheckerSimple()                           | Liveness()                                            | `renamed`            |
| `func`   | HealthCheckerDetailed(config ApplicationConfig) | Readiness(ctx context.Context, config RedinessConfig) | `renamed + context`  |
| `struct` | ApplicationHealthSimple                         | LivenessResponse                                      | `renamed`            |
| `struct` | ApplicationHealthDetailed                       | ReadinessResponse                                     | `renamed`            |
| `struct` | ApplicationConfig                               | RedinessConfig                                        | `renamed`            |
| `struct` | Integration                                     | Integration                                           | `added IsOptional`🚨 |
| `struct` | Auth                                            | ❌                                                    | `removed`            |
| `struct` | IntegrationConfig                               | ❌                                                    | `removed`            |
| `struct` | HTTPHeader                                      | ❌                                                    | `removed`            |
| `const`  | Redis                                           | ❌                                                    | `removed`            |
| `const`  | Memcached                                       | ❌                                                    | `removed`            |
| `const`  | Web                                             | ❌                                                    | `removed`            |
| `const`  | Custom                                          | ❌                                                    | `removed`            |

**Notice** Now the `Readiness` function needs `ctx context.Context` 🚨

### Key Architectural Changes

#### 1. Interface-Driven Design (No More Magic Strings)

**Before v2:** Config structs with magic string discriminators

```go
type IntegrationConfig struct {
    Type    string  // "redis" | "memcached" | "web" | "custom"
    Name    string
    Host    string
    Port    string
    DB      int
    // ... 10 more fields scattered around
    Handle  func() error
}
```

**After v2:** Type-safe interface implementations

```go
type IntegrationItem interface {
    Check(ctx context.Context) (Integration, error)
}

// Concrete types:
type redisChecker struct { ... }
type httpChecker struct { ... }
type memcachedChecker struct { ... }
type customChecker struct { ... }

// Constructor functions:
NewRedisChecker(name string, isOptional bool, opt *redis.Options) IntegrationItem
NewHttpChecker(params HttpCheckerParams) IntegrationItem
NewMemcachedChecker(name string, isOptional bool, server string) IntegrationItem
NewCustomChecker(name string, isOptional bool, handler CustomCheckerHandler) IntegrationItem
```

#### 2. Context Support

**Before:** No context, no timeouts, no cancellation

```go
func HealthCheckerDetailed(config ApplicationConfig) ApplicationHealthDetailed
```

**After:** Full context awareness

```go
func Readiness(ctx context.Context, config RedinessConfig) ReadinessResponse
```

Each integration respects:

- HTTP request timeouts
- Kubernetes probe timeouts
- Custom context cancellation

#### 3. Optional Dependencies

**New Feature:** Mark integrations as non-critical

```go
// If this fails, the readiness check still returns true
NewHttpChecker(HttpCheckerParams{
    Name:       "Analytics",
    Url:        "https://analytics.example.com/health",
    IsOptional: true,  // ← New field
})
```

The `Integration` struct now includes:

```go
type Integration struct {
    // ... existing fields
    IsOptional bool `json:"is_optional"`  // New field
}
```

#### 4. Fixed Resource Leaks

- ✅ HTTP response bodies now closed properly (`defer r.Body.Close()`)
- ✅ Redis connections reused (not created/destroyed per check)
- ✅ No goroutine leaks from cancellation

#### 5. Dependency Updates

- ❌ Removed: `github.com/go-redis/redis v6.15.9+incompatible` (ancient, 2019)
- ✅ Added: `github.com/redis/go-redis/v9 v9.23.0` (current, maintained)
- ✅ Updated indirect dependencies

#### 6. API Renames (Kubernetes Semantics)

| Old Name                    | New Name            | Meaning                     |
| --------------------------- | ------------------- | --------------------------- |
| `HealthCheckerSimple()`     | `Liveness()`        | Is the app process alive?   |
| `HealthCheckerDetailed()`   | `Readiness()`       | Are all dependencies ready? |
| `ApplicationHealthSimple`   | `LivenessResponse`  | Simple "running" status     |
| `ApplicationHealthDetailed` | `ReadinessResponse` | Detailed integration status |

---

## Migration Examples

### Function Rename: Liveness

**Before:**

```go
package main

import(
    "fmt"
    "github.com/gritzkoo/golang-health-checker"
)

func main() {
  result := healthcheck.HealthCheckerSimple()
  fmt.Println(result) // {"status": "fully functional"}
}
```

**After:**

```go
package main

import(
    "fmt"
    "github.com/gritzkoo/golang-health-checker/v2/pkg/healthcheck"
)

func main() {
  result := healthcheck.Liveness()
  fmt.Println(result) // {"status": "fully functional"}
}
```

### Function Rename + Context: Readiness

**Before:**

```go
import "github.com/gritzkoo/golang-health-checker"

config := ApplicationConfig{
    Name:    "MyApp",
    Version: "1.0.0",
    Integrations: []IntegrationConfig{
        {
            Type: "redis",
            Name: "cache",
            Host: "localhost",
            Port: "6379",
        },
    },
}

result := healthcheck.HealthCheckerDetailed(config)
```

**After:**

```go
import (
    "context"
    "github.com/gritzkoo/golang-health-checker/v2/pkg/healthcheck"
    "github.com/redis/go-redis/v9"
)

config := healthcheck.RedinessConfig{
    Name:    "MyApp",
    Version: "1.0.0",
    Integrations: []healthcheck.IntegrationItem{
        healthcheck.NewRedisChecker("cache", false, &redis.Options{
            Addr: "localhost:6379",
        }),
    },
}

result := healthcheck.Readiness(context.Background(), config)
```

### Redis Integration

**Before:**

```go
config := ApplicationConfig{
    Integrations: []IntegrationConfig{
        {
            Type: Redis,
            Name: "cache",
            Host: "localhost",
            Port: "6379",
            DB:   0,
            Auth: Auth{
                Password: "secret",
            },
        },
    },
}
```

**After:**

```go
import "github.com/redis/go-redis/v9"

config := RedinessConfig{
    Integrations: []IntegrationItem{
        NewRedisChecker("cache", false, &redis.Options{
            Addr:     "localhost:6379",
            Password: "secret",
            DB:       0,
        }),
    },
}
```

### HTTP Integration

**Before:**

```go
{
    Type: Web,
    Name: "external-api",
    Host: "https://api.example.com/health",
    TimeOut: 5,
    Headers: []HTTPHeader{
        {Key: "Authorization", Value: "Bearer token"},
    },
}
```

**After:**

```go
import "time"

NewHttpChecker(HttpCheckerParams{
    Name:    "external-api",
    Url:     "https://api.example.com/health",
    Timeout: 5 * time.Second,
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
    IsOptional: false,
})
```

### Memcached Integration

**Before:**

```go
{
    Type: Memcached,
    Name: "session-cache",
    Host: "localhost",
    Port: "11211",
}
```

**After:**

```go
NewMemcachedChecker("session-cache", false, "localhost:11211")
```

### Custom Integration

**Before:**

```go
{
    Type: Custom,
    Name: "db-check",
    Host: "postgres-primary",
    Handle: func() error {
        return db.Ping()
    },
}
```

**After:**

```go
NewCustomChecker("db-check", false, func(ctx context.Context) error {
    return db.PingContext(ctx)
})
```

### Optional Dependencies

**New Feature: Mark non-critical integrations**

```go
config := RedinessConfig{
    Integrations: []IntegrationItem{
        // Critical: if this fails, readiness = false
        NewRedisChecker("main-cache", false, opts),

        // Optional: if this fails, readiness = true (but reports the failure)
        NewHttpChecker(HttpCheckerParams{
            Name:       "Analytics",
            Url:        "https://analytics.example.com/health",
            IsOptional: true,  // ← New
        }),
    },
}
```

---

## Bug Fixes

### Race Condition Fixed

- ✅ Fixed: Main thread now waits for all goroutines before draining results channel
- ✅ Removed: Confusing separate goroutine for `wg.Wait()` and channel close

**Before (race condition):**

```go
go func() {
    wg.Wait()
    close(checklist)
}()
result.Duration = time.Since(start).Seconds()  // ❌ Measured too early
for chk := range checklist { ... }
```

**After (correct):**

```go
wg.Wait()                                       // ✅ Wait for all checks
close(checklist)
result.Duration = time.Since(start).Seconds()  // ✅ Accurate timing
for chk := range checklist { ... }
```

### Data Race Fixed

- ✅ Removed: Concurrent writes to `result.Status` from multiple goroutines
- ✅ Fixed: Only main thread writes to `result.Status` after all goroutines complete

### Resource Leaks Fixed

- ✅ HTTP response bodies now properly closed
- ✅ Redis connections reused instead of created/destroyed per check
- ✅ Graceful shutdown with context cancellation

### URL Validation

- ✅ Fixed: `http.NewRequest` errors now handled properly (was silently ignored)
- ✅ Improved: Malformed URLs return explicit error instead of crashing

---

## Performance Improvements

| Metric                | v1           | v2           | Change                              |
| --------------------- | ------------ | ------------ | ----------------------------------- |
| HTTP timeout default  | 1 minute     | 10 seconds   | Better Kubernetes probe integration |
| Redis client          | v6 (ancient) | v9 (current) | Better performance and features     |
| Connection reuse      | None         | Persistent   | Faster repeated checks              |
| Response body cleanup | ❌ Leaked    | ✅ Closed    | No file descriptor exhaustion       |
| Concurrency safety    | ❌ Unsafe    | ✅ Safe      | Works under load                    |

---

## Test Coverage

- ✅ 100% code coverage (up from ~60%)
- ✅ Tests for optional integrations
- ✅ Structured logging in tests
- ✅ All paths exercise context cancellation

---

## Breaking Changes Summary

1. **Function signatures changed** - `Readiness()` now requires `context.Context`
2. **Config structs removed** - Use type-safe constructors instead
3. **Magic strings removed** - No more `Type: "redis"` constants
4. **Struct names changed** - `ApplicationHealthDetailed` → `ReadinessResponse`
5. **New required fields** - `Integration` now has `IsOptional` field

---

## Migration Checklist

- [ ] Update import paths to `github.com/gritzkoo/golang-health-checker/v2/pkg/healthcheck`
- [ ] Replace `HealthCheckerSimple()` with `Liveness()`
- [ ] Add `context.Context` parameter to `HealthCheckerDetailed()` → `Readiness()`
- [ ] Replace config structs with type-safe constructor functions
- [ ] Remove magic string constants (`Redis`, `Memcached`, `Web`, `Custom`)
- [ ] Update struct names in response handling
- [ ] Add `github.com/redis/go-redis/v9` to go.mod
- [ ] Remove old `github.com/go-redis/redis` if present
- [ ] Test with `-race` flag to verify concurrency safety
- [ ] Update Kubernetes probe configurations if needed

---

## What's Coming in Future Versions

- [ ] Connection pooling strategies for multiple concurrent checks
- [ ] Metrics export (Prometheus integration)
- [ ] Distributed tracing support (OpenTelemetry)
- [ ] Health check scheduling and caching
