# Golang Health Checker

[![Build Status](https://travis-ci.org/gritzkoo/golang-health-checker.svg?branch=master)](https://travis-ci.org/gritzkoo/golang-health-checker)

Simple Golang package to simplify applications based in Go, to trace the healthy of the pods

## How to install

```sh
go mod init github.com/my/repo
go get github.com/gritzkoo/golang-health-checker
```

## How to use

There are 2 ways of call this package, __*simple*__ and __*detailed*__ as the example below

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
  e.GET("/health-check/simple", func(c echo.Context) error {
    return c.JSON(http.StatusOK, healthcheck.HealthCheckerSimple())
  })
  // example of detailed call
  e.GET("/health-check/detailed", func(c echo.Context) error {
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
          Port: "11213",
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

## How it works

You can *__git clone__* this repo and try

```sh
docker-compose build
docker-compose up app
```

And using you browser you can call:
* [Simple Call](http://localhost:8888/health-check/simple)
* [Detailed Call](http://localhost:8888/health-check/detailed)

And to run tests

```sh
docker-compose run test
```
