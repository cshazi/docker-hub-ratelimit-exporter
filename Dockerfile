FROM golang:1.15.5 AS builder
WORKDIR /src
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ratelimit-exporter .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/ratelimit-exporter .
CMD ["./ratelimit-exporter"]

EXPOSE 8080
