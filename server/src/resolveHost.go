package main

import (
	"time"

	"github.com/miekg/dns"
)

func resolveHost(host string) string {
	if len(*dnsServerPtr) == 0 {
		return host
	}

	ch := make(chan string)
	go func() {
		m := new(dns.Msg)
		m.SetQuestion(host+".", dns.TypeA)

		c := new(dns.Client)
		in, _, err := c.Exchange(m, *dnsServerPtr+":53")

		if err != nil {
			ch <- host
			return
		}

		if len(in.Answer) == 0 {
			ch <- host
			return
		}

		if t, ok := in.Answer[0].(*dns.A); ok {
			ch <- t.A.String()
			return
		}
	}()

	select {
	case answer := <-ch:
		return answer
	case <-time.After(200 * time.Millisecond):
		return host
	}
}
