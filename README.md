# zerocat(1)

Zerocat uses zeroconf networking (mDNS) to create a 1-to-1 pipe between two hosts on the same local network.

Note on naming: the "cat" part of the name is following the example of [netcat](http://nc110.sourceforge.net/),
the TCP/IP swiss army knife. Whereas netcat requires a port and host name or an IP address, zerocat can establish 
a connection 

NOTE: **Not a replacement for netcat**, it only works with TCP connections on the same LAN, and doesn't allow 
specifying a port.

## Problem this solves
Alice and Bob are on the same local network, want to share data, but don't want to bother with IP addresses or hostnames. 

### File sharing example
Alice has a large file she wants to send to Bob, who is on the same local network.
Alice runs Linux and Bob runs macOS, so AirDrop is out. Sending the file over Slack 
doesn't work, because the file is too large. But Alice and Bob are on a fast wireless 
network together and would benefit from sending it locally.

Bob is chatting with Alice and types `zerocat > big_ol_file.dat`, it just sits there waiting.

Then Alice types:

```
zerocat < big_ol_file.dat
```

Once Bob's command returns to the prompt, the transfer is done.

`zerocat` uses multicast DNS to find Bob's computer and open up a connection between them. 

