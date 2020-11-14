Docker Hub Pull Rate Limits Prometheus exporter
===============================================

Locally
-------

```
go run main.go
```

Dcoker build and run
--------------------

```
docker build -t ratelimit-exporter . && docker run -it -p 8080:8080 ratelimit-exporter
```
