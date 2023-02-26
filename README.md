# golang-health-checker

![test](https://github.com/gritzkoo/golang-health-checker/workflows/test/badge.svg?branch=master)
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

A simple package to allow you to track your application healthy providing two ways of checking:

*__Simple__*: will return a "fully functional" string and with this, you can check if your application is online and responding without any integration check

*__Detailed__*: will return a detailed status for any integration configuration informed on the integrations, just like in the examples below

___
>This package has a `lightweight` version with no extra dependencies. If you are looking to something more simple, please check [golnag-health-checker-lw on github](https://github.com/gritzkoo/golang-health-checker-lw "golang health checker lightweight") or [golang-health-checke-lw at go.pkg.dev](https://pkg.go.dev/github.com/gritzkoo/golang-health-checker-lw "golang health checker lightweight at go.pkg.dev")
___

## How to install

If you are just starting a Go project you must start a go.mod file like below

```sh
go mod init github.com/my/repo
```

Or else, if you already have a started project, just run the command below

```sh
go get github.com/gritzkoo/golang-health-checker
```

## How to use

In this example, we will use the Echo web server to show how to import and use *Simple* and *Detailed* calls.

If you want to check the full options in configurations, look at this [IntegrationConfig struct](https://github.com/gritzkoo/golang-health-checker/blob/master/pkg/healthcheck/structs.go#L45-L54)

### Available integrations

- [x] Redis
- [x] Memcached
- [x] Web integration (https)

```go
package main

import (
 "net/http"

 "github.com/gin-gonic/gin"
 "github.com/gritzkoo/golang-health-checker/pkg/healthcheck"
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
     Host: "redis",       // you can pass host:port and omit Port attribute
     Port: "6379",
     DB:   0, // default value is 0
    }, {
     Type: healthcheck.Memcached, // this prop will determine the kind of check, the list of types available in structs.go
     Name: "Memcached server",    // the name of you integration to display in response
     Host: "memcache",           // you can pass host:port and omit Port attribute
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
  "status": false, # here is the main status of your application when one of the integrations fails.. false will return
  "version": "V1.0.0",
  "date": "2021-08-27 08:18:06.762044096 -0300 -03 m=+24.943851850",
  "duration": 0.283596049,
  "integrations": [
    {
      "name": "Github Integration",
      "kind": "unknown",
      "status": false,
      "response_time": 0,
      "url": "https://github.com/status",
      "errors": "unsuported type of:unknown"
    },
    {
      "name": "Memcached server",
      "kind": "Memcached DB",
      "status": true,
      "response_time": 0.000419116,
      "url": "localhost:11211"
    },
    {
      "name": "redis-user-db",
      "kind": "Redis DB",
      "status": true,
      "response_time": 0.000845594,
      "url": "localhost:6379"
    },
    {
      "name": "Github Integration",
      "kind": "Web service API",
      "status": true,
      "response_time": 0.283513713,
      "url": "https://github.com/status"
    }
  ]
}
```

## Kubernetes liveness and readiness probing

And then, you could call these endpoints manually to see your application health, but, if you are using modern kubernetes deployment, you can config your chart to check your application with the setup below:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-golang-app
spec:
  selector:
    matchLabels:
      app: my-golang-app
  template:
    metadata:
      labels:
        app: my-golang-app
    spec:
      containers:
      - name: my-golang-app
        image: your-app-image:tag
        resources:
          request:
            cpu: 10m
            memory: 5Mi
          limits:
            cpu: 50m
            memory: 50Mi 
        livenessProbe:
          httpGet:
            path: /health-check/liveness
            port: 8888
            scheme: http
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 2
          successThreshold: 1
        readinessProbe:
          httpGet:
            path: /health-check/liveness
            port: 8888
            scheme: http
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 2
          successThreshold: 1
        ports:
        - containerPort: 8888
```
