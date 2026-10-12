package healthcheck

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisChecker struct {
	client     *redis.Client
	name       string
	isOptional bool
}

func (r *redisChecker) Check(ctx context.Context) (Integration, error) {
	var response = Integration{
		Name:       r.name,
		Kind:       "Redis",
		Status:     true,
		IsOptional: r.isOptional,
		Url:        r.client.Options().Addr,
	}
	redisCtx, redisCtxCancel := context.WithTimeout(ctx, time.Second)
	_, err := r.client.Ping(redisCtx).Result()
	redisCtxCancel()

	if err != nil {
		return response, err
	}

	return response, nil
}

// NewRedisChecker func return a new IntegrationItem with redis configurations
func NewRedisChecker(name string, isOptional bool, opt *redis.Options) IntegrationItem {
	return &redisChecker{
		client:     redis.NewClient(opt),
		name:       name,
		isOptional: isOptional,
	}
}
