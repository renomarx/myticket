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
- Test webhook:
```bash
./scripts/send_ticket.sh
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
hey -m POST -D 'scripts/ticket.txt' http://localhost:9098/ticket
```

Benchmark on my very sloowwww machine:
cpu `Intel(R) Core(TM) m3-6Y30 CPU @ 0.90GHz`
```
Summary:
  Total:	0.7807 secs
  Slowest:	0.7164 secs
  Fastest:	0.0050 secs
  Average:	0.1598 secs
  Requests/sec:	256.1792


Response time histogram:
  0.005 [1]	|
  0.076 [82]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.147 [41]	|■■■■■■■■■■■■■■■■■■■■
  0.218 [12]	|■■■■■■
  0.290 [23]	|■■■■■■■■■■■
  0.361 [15]	|■■■■■■■
  0.432 [10]	|■■■■■
  0.503 [9]	|■■■■
  0.574 [2]	|■
  0.645 [3]	|■
  0.716 [2]	|■


Latency distribution:
  10% in 0.0089 secs
  25% in 0.0314 secs
  50% in 0.0949 secs
  75% in 0.2624 secs
  90% in 0.3990 secs
  95% in 0.4953 secs
  99% in 0.6552 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0038 secs, 0.0050 secs, 0.7164 secs
  DNS-lookup:	0.0029 secs, 0.0000 secs, 0.0369 secs
  req write:	0.0003 secs, 0.0000 secs, 0.0048 secs
  resp wait:	0.1548 secs, 0.0048 secs, 0.6853 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0002 secs

Status code distribution:
  [200]	200 responses


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

- On performances: more things could be done to improve performances, but I assume it should be enough for our needs. At scale, we could eventually have many instances of this service running in parallel with a HAProxy in front dispatching requests among the instances. We would however need postgres running on a high-performances machine (because a distributed postrgres would not improve performances, although it could improve safety and resiliency)
