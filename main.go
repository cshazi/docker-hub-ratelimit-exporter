package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	listeningAddr = ""
	metricsPath   = ""
	authUrl       = ""
	rateLimitUrl  = ""
)

func main() {
	flag.StringVar(&listeningAddr,
		"addr",
		"0.0.0.0:8080",
		"HTTP Server address")
	flag.StringVar(&metricsPath,
		"metrics-path",
		"/metrics",
		"Metrics URL path")
	flag.StringVar(&authUrl,
		"auth-url",
		"https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull",
		"Docker Hub auth URL")
	flag.StringVar(&rateLimitUrl,
		"rate-limit-url",
		"https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest",
		"Docker Hub rate limit URL")
	flag.Parse()

	log.Println("Starting Dockerhub exporter.")
	log.Println("Listening on:", "'"+listeningAddr+"'")

	http.HandleFunc(metricsPath, rateLimitQueryFunc)
	http.ListenAndServe(listeningAddr, nil)
}

func rateLimitQueryFunc(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("HEAD", rateLimitUrl, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+getToken(client))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	limitStr := resp.Header.Get("RateLimit-Limit")
	remainingStr := resp.Header.Get("RateLimit-Remaining")

	writeTo(w, getLimit(limitStr), getLimit(remainingStr), getInterval(limitStr))
}

func getLimit(header string) int {
	limit, err := strconv.Atoi(strings.Split(header, ";w=")[0])
	if err != nil {
		log.Println(err)
	}
	return limit
}

func getToken(client *http.Client) string {
	resp, err := client.Get(authUrl)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
	}

	return data["token"].(string)
}

func getInterval(header string) int {
	interval, err := strconv.Atoi(strings.Split(header, ";w=")[1])
	if err != nil {
		log.Println(err)
	}
	return interval
}

func writeTo(w io.Writer, limit int, remaining int, interval int) (int64, error) {
	buf := &bytes.Buffer{}

	// Dockerhub Rate Rate
	buf.WriteString(fmt.Sprintf("# HELP %s %s\n", "dockerhub_ratelimit_rate_limit", "Dockerhub Rate Limit"))
	buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", "dockerhub_ratelimit_rate_limit", "gauge"))
	buf.WriteString(fmt.Sprintf("%s %d\n", "dockerhub_ratelimit_rate_limit", limit))

	buf.WriteString(fmt.Sprintf("# HELP %s %s\n", "dockerhub_ratelimit_rate_remaining", "Dockerhub Rate Remaining"))
	buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", "dockerhub_ratelimit_rate_remaining", "gauge"))
	buf.WriteString(fmt.Sprintf("%s %d\n", "dockerhub_ratelimit_rate_remaining", remaining))

	buf.WriteString(fmt.Sprintf("# HELP %s %s\n", "dockerhub_ratelimit_rate_interval", "Dockerhub Rate Interval"))
	buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", "dockerhub_ratelimit_rate_interval", "gauge"))
	buf.WriteString(fmt.Sprintf("%s %d\n", "dockerhub_ratelimit_rate_interval", interval))

	io.Copy(w, buf)

	return 0, nil
}
