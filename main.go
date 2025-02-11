package main

import (
	"fmt"

	"github.com/miekg/dns"
)

type dnsHandler struct{}

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

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		fmt.Printf("Received query: %s\n", question.Name)
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
