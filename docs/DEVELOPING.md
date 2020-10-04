# How it works

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

## The next steps to this package is adding compatibility with the integrations below:

* [x] Redis
* [x] Memcache
* [x] Web API calls
* [ ] Mongodb
* [ ] Mysql
* [ ] Postgres
* [ ] RabbitMQ

 And any other that in this very moment we just don't known
