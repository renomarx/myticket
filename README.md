Myticket
==========

See the [exercise subject](https://gist.github.com/n-e/917fae053f543667d2e77798dde6d5cd)

Requirements
-----------

- Unix system (tested on Linux, should work on macOS but not tested), compatible with GNU `make` command
- docker & docker-compose

Run
----

```bash
make dev
```

Then access to:
- Metrics (Prometheus): http://localhost:9098/metrics
- logs: `docker logs -f myticket_dev`
- Postgres DB: accessible on `localhost:5432`
- Test ping:
```bash
curl -v http://localhost:9098/ping
```

Test
----

```bash
make tests
```



Benchmarks
----------

Using [hey](https://github.com/rakyll/hey)

```bash
hey 'http://localhost:9098/ticket' --data 'TODO'
```

Benchmark on my very sloowwww machine:
cpu `Intel(R) Core(TM) m3-6Y30 CPU @ 0.90GHz`
```
# TODO

```




Architecture / technical choices
---------

- Based on [a golang standard project layout](https://github.com/golang-standards/project-layout): a standard layout that adapts pretty well to every kind of golang project

- Using [Hexagonal architecture](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3): in my opinion, this architecture is good for maintainability without being too strict

- `pkg` vs `internal`: everything is in `pkg` directory because nothing aims to be private, but I didn't particularly paid attention
about packages being easily exportable

- Prometheus metrics: I like to have metrics on my services, and I like to use prometheus for that, particularly for services exposing a http API.

- Using [sqlx](https://github.com/jmoiron/sqlx): because the most used ORM in go is gorm and I don't like it
because it doesn't respect go idioms (methods chaining and bad errors handling). I found myself pretty quickly limited by ORMs too, so I usually
prefer not using one (even if it saves some time at the beginning of a project)

- Using [dbmate](https://github.com/amacneil/dbmate), because it's a good tool to handle migrations, especially useful when your company have more than
one main coding language (like nodejs and go): you can use the same tool in different projects
