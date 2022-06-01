package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gritzkoo/golang-health-checker/v2/pkg/healthcheck"
)

func main() {
	// all the content below is just an example
	// Gin instance
	e := gin.Default()

	// example of simple call
	e.GET("/health-check/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, healthcheck.HealthCheckerSimple())
	})
	// example of detailed call
	e.GET("/health-check/readiness", func(c *gin.Context) {
		// define all integrations of your application with type healthcheck.ApplicationConfig
		myApplicationConfig := healthcheck.ApplicationConfig{ // check the full list of available props in structs.go
			Name:    "You APP Name", // optional prop
			Version: "V1.0.0",       // optional prop
			Integrations: []healthcheck.IntegrationConfig{ // mandatory prop
				{
					Type: healthcheck.Redis, // this prop will determine the kind of check, the list of types available in structs.go
					Name: "redis-user-db",   // the name of you integration to display in response
					Host: "localhost",       // you can pass host:port and omit Port attribute
					Port: "6379",
					DB:   0, // default value is 0
				}, {
					Type: healthcheck.Memcached, // this prop will determine the kind of check, the list of types available in structs.go
					Name: "Memcached server",    // the name of you integration to display in response
					Host: "localhost",           // you can pass host:port and omit Port attribute
					Port: "11211",
				}, {
					Type:    healthcheck.Web,             // this prop will determine the kind of check, the list of types available in structs.go
					Name:    "Github Integration",        // the name of you integration to display in response
					Host:    "https://github.com/status", // you can pass host:port and omit Port attribute
					TimeOut: 5,                           // default value to web call is 10s
					Headers: []healthcheck.HTTPHeader{ // to customize headers to perform a GET request
						{
							Key:   "Accept",
							Value: "application/json",
						},
					},
				}, {
					Type:    "unknown",                   // this prop will determine the kind of check, the list of types available in structs.go
					Name:    "Github Integration",        // the name of you integration to display in response
					Host:    "https://github.com/status", // you can pass host:port and omit Port attribute
					TimeOut: 5,                           // default value to web call is 10s
					Headers: []healthcheck.HTTPHeader{ // to customize headers to perform a GET request
						{
							Key:   "Accept",
							Value: "application/json",
						},
					},
				},
			},
		}
		c.JSON(http.StatusOK, healthcheck.HealthCheckerDetailed(myApplicationConfig))
	})
	// Start server
	e.Run(":8888")
}
