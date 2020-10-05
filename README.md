# golang-health-checker

![test](https://github.com/gritzkoo/golang-health-checker/workflows/test/badge.svg?branch=master)
[![Build Status](https://travis-ci.org/gritzkoo/golang-health-checker.svg?branch=master)](https://travis-ci.org/gritzkoo/golang-health-checker)
[![Coverage Status](https://coveralls.io/repos/github/gritzkoo/golang-health-checker/badge.svg?branch=master)](https://coveralls.io/github/gritzkoo/golang-health-checker?branch=master)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gritzkoo/golang-health-checker)
![GitHub repo size](https://img.shields.io/github/repo-size/gritzkoo/golang-health-checker)
![GitHub](https://img.shields.io/github/license/gritzkoo/golang-health-checker)
![GitHub issues](https://img.shields.io/github/issues/gritzkoo/golang-health-checker)

A simple package to allow you to track your application healthy providing two ways of checking:

*__Simple__*: will return a "fully functional" string and with this, you can check if your application is online and responding without any integration check

*__Detailed__*: will return a detailed status for any integration configuration informed on the integrations just like in the examples below

## How to install

If you are just starting a Go projetct you must start a go.mod file like below

```sh
go mod init github.com/my/repo
```

Or else, you already has a started project, just run the command below

```sh
go get github.com/gritzkoo/golang-health-checker
```

## How to use

In this example, we will use the Echo web server to show how to import and use *Simple* and *Detailed* calls.

If you want check the full options in configurations, look this [IntegrationConfig struct](https://github.com/gritzkoo/golang-health-checker/blob/master/pkg/healthcheck/structs.go#L45-L54)

### Available integrations

- [x] Redis
- [x] Memcached
- [x] Web integration (https)

```go
package main

import (
  "net/http"

  "github.com/gritzkoo/golang-health-checker/pkg/healthcheck"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

func main() {
  // all the content below is just an example
  // Echo instance
  e := echo.New()
  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  // example of simple call
  e.GET("/health-check/liveness", func(c echo.Context) error {
    return c.JSON(http.StatusOK, healthcheck.HealthCheckerSimple())
  })
  // example of detailed call
  e.GET("/health-check/readiness", func(c echo.Context) error {
    // define all integrations of your application with type healthcheck.ApplicationConfig
    myApplicationConfig := healthcheck.ApplicationConfig{ // check the full list of available props in structs.go
      Name:    "You APP Name", // optional prop
      Version: "V1.0.0",       // optional prop
      Integrations: []healthcheck.IntegrationConfig{ // mandatory prop
        {
          Type: healthcheck.Redis, // this prop will determine the kind of check, the list of types available in structs.go
          Name: "redis-user-db",   // the name of you integration to display in response
          Host: "redis",           // you can pass host:port and omit Port attribute
          Port: "6379",
          DB:   0, // default value is 0
        }, {
          Type: healthcheck.Memcached, // this prop will determine the kind of check, the list of types available in structs.go
          Name: "Memcached server",    // the name of you integration to display in response
          Host: "memcache",            // you can pass host:port and omit Port attribute
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
        },
      },
    }
    return c.JSON(http.StatusOK, healthcheck.HealthCheckerDetailed(myApplicationConfig))
  })
  // Start server
  e.Logger.Fatal(e.Start(":8888"))
}

```

This simple call will return a JSON as below

```json
{
  "status": "fully functional"
}
```

And detailed call will return a JSON as below

```json
{
  "name": "You APP Name",
  "status": true, # here is the main status of your application when one of the integrations fails.. false will return
  "version": "V1.0.0",
  "date": "Mon Jan _2 15:04:05 MST 2006",
  "Duration": 0.53102304,
  "integrations": [
    {
      "name": "redis-user-db",
      "kind": "Redis DB",
      "status": true,
      "response_time": 0.001160881,
      "url": "localhost:6379"
    },
    {
      "name": "Memcached server",
      "kind": "Memcached DB",
      "status": true,
      "response_time": 0.036013866,
      "url": "localhost:11211"
    },
    {
      "name": "Github Integration",
      "kind": "Web service API",
      "status": true,
      "response_time": 0.493425975,
      "url": "https://github.com/status"
    }
  ]
}
```

## Kubernetes liveness and readiness probing

And then, you could call this endpoints manually to see your application health, but, if you are using modern kubernetes deployment, you can config your chart to check your application with the setup below:

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: liveness-http
spec:
  containers:
  - name: liveness
    image: 'go' #your application image
    args:
    - /server
    livenessProbe:
      httpGet:
        path: /health-check/liveness
        port: 80
        httpHeaders:
        - name: Custom-Header
          value: Awesome
      initialDelaySeconds: 3
      periodSeconds: 3
  - name: readiness
    image: 'go' #your application image
    args:
    - /server
    readinessProbe:
      httpGet:
        path: /health-check/readiness
        port: 80
        httpHeaders:
        - name: Custom-Header
          value: Awesome
      initialDelaySeconds: 3
      periodSeconds: 3
```
