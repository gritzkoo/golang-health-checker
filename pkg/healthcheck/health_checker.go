package healthcheck

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis"
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
		result    = ApplicationHealthDetailed{
			Name:         config.Name,
			Version:      config.Version,
			Status:       true,
			Date:         start.Format(time.RFC3339),
			Duration:     0,
			Integrations: []Integration{},
		}
	)
	wg.Add(len(config.Integrations))
	for _, v := range config.Integrations {
		switch v.Type {
		case Redis:
			go checkRedisClient(v, &result, &wg, checklist)
		case Memcached:
			go checkMemcachedClient(v, &result, &wg, checklist)
		case Web:
			go checkWebServiceClient(v, &result, &wg, checklist)
		case Custom:
			go CheckCustom(v, &result, &wg, checklist)
		default:
			go defaultAction(v, &result, &wg, checklist)
		}
	}
	go func() {
		wg.Wait()
		close(checklist)
		result.Duration = time.Since(start).Seconds()
	}()
	for chk := range checklist {
		result.Integrations = append(result.Integrations, chk)
	}
	return result
}

func checkRedisClient(config IntegrationConfig, result *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		start        = time.Now()
		myStatus     = true
		host         = validateHost(config)
		DB           = 0
		errorMessage = ""
	)
	if config.DB > 0 {
		DB = config.DB
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: config.Auth.Password, // no password set
		DB:       DB,                   // use default DB
	})
	response, err := rdb.Ping().Result()
	rdb.Close()
	if err != nil {
		myStatus = false
		result.Status = false
		errorMessage = fmt.Sprintf("response: %s error message: %s", response, err.Error())
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         RedisIntegration,
		Status:       myStatus,
		ResponseTime: time.Since(start).Seconds(),
		URL:          host,
		Error:        errorMessage,
	}
}

func checkMemcachedClient(config IntegrationConfig, result *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		start        = time.Now()
		myStatus     = true
		host         = validateHost(config)
		errorMessage = ""
	)
	mcClient := memcache.New(host)
	err := mcClient.Ping()
	if err != nil {
		myStatus = false
		result.Status = false
		errorMessage = err.Error()
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         MemcachedIntegration,
		Status:       myStatus,
		ResponseTime: time.Since(start).Seconds(),
		URL:          host,
		Error:        errorMessage,
	}
}

func checkWebServiceClient(config IntegrationConfig, result *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		host         = validateHost(config)
		timeout      = 10
		myStatus     = true
		start        = time.Now()
		errorMessage = ""
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
	if err != nil {
		myStatus = false
		result.Status = false
		errorMessage = err.Error()
	} else if response.StatusCode != 200 {
		myStatus = false
		result.Status = false
		errorMessage = fmt.Sprintf("Expected request status code 200 got %d", response.StatusCode)
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         WebServiceIntegration,
		Status:       myStatus,
		ResponseTime: time.Since(start).Seconds(),
		URL:          host,
		Error:        errorMessage,
	}
}

func CheckCustom(config IntegrationConfig, result *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	var (
		myStatus     = true
		start        = time.Now()
		host         = validateHost(config)
		errorMessage = ""
	)
	tmp := config.Handle()
	if tmp != nil {
		myStatus = false
		result.Status = false
		errorMessage = tmp.Error()
	}
	checklist <- Integration{
		Name:         config.Name,
		Kind:         CustomizedTestFunction,
		Status:       myStatus,
		ResponseTime: time.Since(start).Seconds(),
		URL:          host,
		Error:        errorMessage,
	}
}

func defaultAction(config IntegrationConfig, result *ApplicationHealthDetailed, wg *sync.WaitGroup, checklist chan Integration) {
	defer (*wg).Done()
	result.Status = false
	checklist <- Integration{
		Name:         config.Name,
		Kind:         config.Type,
		Status:       false,
		ResponseTime: 0,
		URL:          config.Host,
		Error:        fmt.Sprintf("unsuported type of:" + config.Type),
	}
}

func validateHost(config IntegrationConfig) string {
	var host = config.Host
	if config.Port != "" {
		host = host + ":" + config.Port
	}
	return host
}
