# zerocat

Zerocat uses zeroconf networking (mDNS) to create a 1-to-1 pipe between two hosts on the same local network.

Note on naming: the "cat" part of the name is following the example of [netcat](http://nc110.sourceforge.net/),
the TCP/IP swiss army knife. Whereas netcat requires a port and host name or an IP address, zerocat can establish a connection automatically.

NOTE: **Not a replacement for netcat**, it only works with TCP connections on the same LAN, and doesn't allow 
specifying a port.

## Problem this solves
Alice and Bob are on the same local network, want to share data, but don't want to bother with IP addresses or hostnames. 

### File sharing example
Alice has a large file she wants to send to Bob, who is on the same local network.
Alice runs Linux and Bob runs macOS, so AirDrop is out. Sending the file over Slack 
doesn't work, because the file is too large. But Alice and Bob are on a fast wireless 
network together and would benefit from sending it locally.

Bob types `zerocat > big_ol_file.dat`, it just sits there waiting.

Then Alice types:

```
zerocat < big_ol_file.dat
```

Once Bob's command returns to the prompt, the transfer is done.

`zerocat` uses multicast DNS to find Bob's computer and open up a connection between them. 


## How it works

When the first instance of `zerocat` starts running, it starts an mDNS server listening for queries with the name `_zerocat._tcp` and is ready to respond with a SRV record containing the local machines' IP and a port number. It then queries the local network for `_zerocat._tcp`, since this is the first instance, nothing happens, and the program sits and waits.

When the second instance of `zerocat` on the local network starts running, it queries the local network for `_zerocat._tcp`, then the first instance responds with a SRV record, this is enough information to create a TCP socket to the first machine. At this point the mDNS server is shutdown and a socket is opened between the two machines, if data is available to be read, it is, and if there is data waiting to be written, it is. When one side reads an EOF, the connection is closed.
