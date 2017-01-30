package main

import (
	"fmt"
	"github.com/hashicorp/mdns"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const port = 8000

var host string
var server *mdns.Server
var alive bool = true

func main() {
	// turn off logging
	devNull, _ := os.Open(os.DevNull)
	log.SetOutput(devNull)
	server = createServer()
	defer server.Shutdown()
	query()
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
			if entry.Host[0:len(entry.Host)-1] != host &&
				entry.Port == port {
				pipe(entry.AddrV4)
			}
		}
	}()

	for alive {
		mdns.Query(&mdns.QueryParam{
			Service: "_zerocat._tcp",
			Entries: entriesCh,
		})
		time.Sleep(5)
	}
	close(entriesCh)
}

// Spawn two goroutines, one to listen and read the connection,
// and another to dial and write to the connection
func pipe(addr net.IP) {
	alive = false
	go listenAndRead()
	go dialAndWrite(addr)
}

func listenAndRead() {
	laddr := fmt.Sprintf("0.0.0.0:%d", port)
	listener, _ := net.Listen("tcp", laddr)
	conn, err := listener.Accept()
	if err != nil && conn != nil {
		io.Copy(os.Stdout, conn)
		conn.Close()
	} else {
		fmt.Printf("listening connection to %v failed: %v", laddr, err)
	}
}

func dialAndWrite(addr net.IP) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil && conn != nil {
		io.Copy(conn, os.Stdin)
		conn.Close()
	} else {
		fmt.Printf("dialing connection to %v failed: %v", addr, err)
	}
}
