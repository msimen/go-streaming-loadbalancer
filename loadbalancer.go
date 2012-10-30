package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type LoadbalancerResponse struct {
	Redirect    string `json:"redirect"`
	ClientIP    string `json:"client_ip"`
	CountryCode string `json:"country_code"`
	Timestamp   int64  `json:"timestamp"`
}

func Loadbalancer(w http.ResponseWriter, req *http.Request) {
	// Set some headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "go-streaming-loadbalancer")

	// Connect to Redis
	c := pool.Get()

	if c.Err() != nil {
		log.Print("Error: ", c.Err())

		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Could not connect to Redis server", 503)

		return
	}

	// Make sure that the Redis connection is closed
	defer c.Close()

	// Retrieve the redirect ip
	redirectIP, err := getRedirectIP(c)

	if err != nil {
		log.Print("Error: ", err)
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)

	cc := ""
	loc := gi.GetLocationByIP(host)

	if loc != nil {
		cc = loc.CountryCode
	}

	// Populate the loadbalancer response struct
	response := LoadbalancerResponse{
		Redirect:    redirectIP,
		Timestamp:   time.Now().Unix(),
		CountryCode: cc,
		ClientIP:    host,
	}

	// Marshal the JSON
	json, err := json.Marshal(response)

	if err != nil {
		fmt.Println("JSON Error:", err)
	}

	// Return the JSON to the client
	w.Write(json)
}
