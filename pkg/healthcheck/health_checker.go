package healthcheck

import (
	"net/http"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/openpgp/errors"
)

// HealthCheckerSimple performs a simple check of the application
func HealthCheckerSimple() ApplicationHealthSimple {
	return ApplicationHealthSimple{
		Status: "fully functional",
	}
}

// HealthCheckerDetailed perform a check for every integration informed
func HealthCheckerDetailed(config ApplicationConfig) ApplicationHealthDetailed {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		checklist = make(chan Integration, len(config.Integrations))
		rt        = ApplicationHealthDetailed{
			Name:         config.Name,
			Version:      config.Version,
			Status:       true,
			Date:         start.String(),
			Duration:     0,
			Integrations: []Integration{},
		}
	)
	for _, v := range config.Integrations {
		switch v.Type {
		case Redis:
			wg.Add(1)
			go checkRedisClient(v, &rt, &wg, checklist)
		case Memcached:
			wg.Add(1)
			go checkMemcachedClient(v, &rt, &wg, checklist)
		case Web:
			wg.Add(1)
			go checkWebServiceClient(v, &rt, &wg, checklist)
		default:
			wg.Add(1)
			go defaultAction(v, &rt, &wg, checklist)
		}
	}
	go func() {
		wg.Wait()
		close(checklist)
		rt.Duration = time.Now().Sub(start).Seconds()
	}()
	for chk := range checklist {
		rt.Integrations = append(rt.Integrations, chk)
	}
	return rt
}

func checkRedisClient(config IntegrationConfig, rt *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		host = validateHost(config)
		DB   = 0
	)
	if config.DB > 0 {
		DB = config.DB
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: config.Auth.Password, // no password set
		DB:       DB,                   // use default DB
	})
	start := time.Now()
	response, err := rdb.Ping().Result()
	elapsed := time.Now().Sub(start)
	rdb.Close()
	if err != nil || response != "PONG" {
		rt.Status = false
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         RedisIntegration,
		Status:       response == "PONG",
		ResponseTime: elapsed.Seconds(),
		URL:          host,
		Error:        err,
	}
}

func checkMemcachedClient(config IntegrationConfig, rt *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var host = validateHost(config)
	mcClient := memcache.New(host)
	start := time.Now()
	err := mcClient.Ping()
	elapsed := time.Now().Sub(start)
	localStatus := true
	if err != nil {
		localStatus = false
		rt.Status = false
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         MemcachedIntegration,
		Status:       localStatus,
		ResponseTime: elapsed.Seconds(),
		URL:          host,
		Error:        err,
	}
}

func checkWebServiceClient(config IntegrationConfig, rt *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		host     = validateHost(config)
		timeout  = 10
		myStatus = true
		start    = time.Now()
	)
	if config.TimeOut > 0 {
		timeout = config.TimeOut
	}
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	request, _ := http.NewRequest("GET", host, nil)

	if len(config.Headers) > 0 {
		for _, v := range config.Headers {
			request.Header.Add(v.Key, v.Value)
		}
	}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		myStatus = false
		rt.Status = false
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         WebServiceIntegration,
		Status:       myStatus,
		ResponseTime: time.Now().Sub(start).Seconds(),
		URL:          host,
		Error:        err,
	}
}

func defaultAction(config IntegrationConfig, rt *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	rt.Status = false
	checklist <- Integration{
		Name:         config.Name,
		Kind:         config.Type,
		Status:       false,
		ResponseTime: 0,
		URL:          config.Host,
		Error:        errors.UnsupportedError("unsuported type of:" + config.Type),
	}
}

func validateHost(config IntegrationConfig) string {
	var host = config.Host
	if config.Port != "" {
		host = host + ":" + config.Port
	}
	return host
}
