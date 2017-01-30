package main

import (
	"fmt"
	"github.com/hashicorp/mdns"
	"log"
	"os"
	"time"
)

var host string

func main() {
	// turn off logging
	devNull, _ := os.Open(os.DevNull)
	log.SetOutput(devNull)
	server := createServer()
	defer server.Shutdown()
	query()
	//io.Copy(os.Stdout, os.Stdin)
}

func createServer() *mdns.Server {
	host, _ = os.Hostname()
	info := []string{"Zero configuration netcat"}
	// TODO: Use https://github.com/phayes/freeport
	service, err := mdns.NewMDNSService(host,
		"_zerocat._tcp", "", "",
		8000, nil, info,
	)
	if err != nil {
		fmt.Printf("error creating new mDNS service:  %v", err)
	}
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	return server
}

func query() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			if entry.Host[0:len(entry.Host)-1] != host {
				fmt.Printf("Got new entry: %v, %v\n", entry.AddrV4, entry.Port)
			}
		}
	}()

	for {
		// Start the lookup
		mdns.Query(&mdns.QueryParam{
			Service: "_zerocat._tcp",
			Entries: entriesCh,
		})
		time.Sleep(5)
	}
	close(entriesCh)
}
