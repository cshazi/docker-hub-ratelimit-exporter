# Docker Hub Pull Rate Limits Prometheus exporter

Exposes metrics of rate limits and rate remaining pulls from the Docker Hub API, to a Prometheus compatible endpoint.
The exporter uses the anonymous user to determine the rate limit.

## Configuration

* `addr` Address on which to expose metrics. (default: "0.0.0.0:8080")
* `metrics-path` Path under which to expose metrics. (default: "metrics")
* `auth-url` Docker Hub auth URL (default: "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull")
* `rate-limit-url` Docker Hub rate limit URL (default: "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest")

## Run locally

```
go run main.go
```

## Build a docker image

```
docker build -t ratelimit-exporter . && docker run -it -p 8080:8080 ratelimit-exporter
```

## Metrics

Metrics will be made available on port 8080 by default. Below is an example of the metrics as exposed by this exporter. 

```
# HELP dockerhub_ratelimit_rate_limit Dockerhub Rate Limit
# TYPE dockerhub_ratelimit_rate_limit gauge
dockerhub_ratelimit_rate_limit 100
# HELP dockerhub_ratelimit_rate_remaining Dockerhub Rate Remaining
# TYPE dockerhub_ratelimit_rate_remaining gauge
dockerhub_ratelimit_rate_remaining 92
# HELP dockerhub_ratelimit_rate_interval Dockerhub Rate Interval
# TYPE dockerhub_ratelimit_rate_interval gauge
dockerhub_ratelimit_rate_interval 21600
```
