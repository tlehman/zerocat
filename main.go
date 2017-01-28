package main

import (
	"fmt"
	"github.com/hashicorp/mdns"
	"io"
	"os"
)

func main() {
	// turn off logging
	// devNull, _ := os.Open(os.DevNull)
	// log.SetOutput(devNull)
	server := createServer()
	defer server.Shutdown()
	//performLookup()
	// send input to output

	io.Copy(os.Stdout, os.Stdin)
}

func createServer() *mdns.Server {
	host, _ := os.Hostname()
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

func performLookup() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry.Name)
		}
	}()

	// Start the lookup
	mdns.Query(&mdns.QueryParam{
		Service: "_zerocat._tcp",
		Entries: entriesCh,
	})
	close(entriesCh)
}
