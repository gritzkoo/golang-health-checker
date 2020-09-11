package healthcheck

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis"
)

// used to track internally if one of the integrated services fails
var mainStatus = true

// HealthCheckerSimple performs a simple check of the application
func HealthCheckerSimple() ApplicationHealthSimple {
	return ApplicationHealthSimple{
		Status: "fully functional",
	}
}

// HealthCheckerDetailed perform a check for every integration informed
func HealthCheckerDetailed(config ApplicationConfig) ApplicationHealthDetailed {
	var integrations []Integration
	start := time.Now()
	for _, v := range config.Integrations {
		switch v.Type {
		case Redis:
			temp := checkRedisClient(v)
			integrations = append(integrations, temp)
			break
		case Memcached:
			temp := checkMemcachedClient(v)
			integrations = append(integrations, temp)
			break
		case Web:
			temp := checkWebServiceClient(v)
			integrations = append(integrations, temp)
			break
		default:
			fmt.Println("Configuration error, type unsuported:", v)
			break
		}
	}
	return ApplicationHealthDetailed{
		Status:       mainStatus,
		Name:         config.Name,
		Version:      config.Version,
		Date:         time.UnixDate,
		Duration:     time.Now().Sub(start).Seconds(),
		Integrations: integrations,
	}
}

func checkRedisClient(config IntegrationConfig) Integration {
	var host = validateHost(config)
	var DB = 0
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
	if err != nil {
		mainStatus = false
		fmt.Println(err)
	}
	return Integration{
		Name:         config.Name,
		Kind:         RedisIntegration,
		Status:       response == "PONG",
		ResponseTime: elapsed.Seconds(),
		URL:          host,
		Error:        err,
	}
}

func checkMemcachedClient(config IntegrationConfig) Integration {
	var host = validateHost(config)
	mcClient := memcache.New(host)
	start := time.Now()
	err := mcClient.Ping()
	elapsed := time.Now().Sub(start)
	if err != nil {
		mainStatus = false
		fmt.Println(err)
	}
	return Integration{
		Name:         config.Name,
		Kind:         MemcachedIntegration,
		Status:       err == nil,
		ResponseTime: elapsed.Seconds(),
		URL:          host,
		Error:        err,
	}
}

func checkWebServiceClient(config IntegrationConfig) Integration {
	var host = validateHost(config)
	var timeout = 10
	var myStatus = true
	if config.TimeOut > 0 {
		timeout = config.TimeOut
	}
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	start := time.Now()
	request, err := http.NewRequest("GET", host, nil)
	if err != nil {
		mainStatus = false
		myStatus = false
		fmt.Println(err)
	} else {
		if len(config.Headers) > 0 {
			for _, v := range config.Headers {
				request.Header.Add(v.Key, v.Value)
			}
		}
		response, err := client.Do(request)
		if err != nil || response.StatusCode != 200 {
			myStatus = false
			mainStatus = false
			fmt.Println(err)
		}
	}
	elapsed := time.Now().Sub(start)
	return Integration{
		Name:         config.Name,
		Kind:         WebServiceIntegration,
		Status:       myStatus,
		ResponseTime: elapsed.Seconds(),
		URL:          host,
		Error:        err,
	}
}

func validateHost(config IntegrationConfig) string {
	var host = config.Host
	if config.Port != "" {
		host = host + ":" + config.Port
	}
	return host
}
