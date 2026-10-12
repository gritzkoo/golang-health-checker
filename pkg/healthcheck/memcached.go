package healthcheck

import (
	"context"

	"github.com/bradfitz/gomemcache/memcache"
)

// memcachedChecker
type memcachedChecker struct {
	name       string
	server     string
	client     *memcache.Client
	isOptional bool
}

func (m *memcachedChecker) Check(ctx context.Context) (Integration, error) {
	var response = Integration{
		Name:       m.name,
		Kind:       "Memcached",
		Status:     true,
		IsOptional: m.isOptional,
		Url:        m.server,
	}
	if err := m.client.Ping(); err != nil {
		response.Status = false
		response.Error = err.Error()
	}
	m.client.Close()
	return response, nil
}

// NewMemcachedChecker func return a new IntegrationItem with memcached configurations
func NewMemcachedChecker(name string, isOptional bool, server string) IntegrationItem {
	return &memcachedChecker{
		name:       name,
		server:     server,
		isOptional: isOptional,
		client:     memcache.New(server),
	}
}
