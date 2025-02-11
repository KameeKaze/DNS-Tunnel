# DNS-Tunnel

## How It Works

The client encodes the data and splits it into small chunks. Each chunk is hex encoded and appended as a subdomain to the configured domain (e.g., deadbeef.example.com). Then the query is sent to the DNS server controlled by the attacker. The server extracts and decodes the data from received queries, while still resolving all queries using Google's DNS server (8.8.8.8).

## Installation
### server
```sh
git clone https://github.com/KameeKaze/DNS-Tunnel.git
cd DNS-Tunnel
go build -o dns-tunnel
./dns-tunnel
```
## Usage

```
./dns-tunnel -d example.com -f output.txt
```

## Options
```
  -d string
        domain used to exfiltrate data (default "example.com")
  -f string
        file used to store incomming data
```

## TODO

- encrypt data
- write client