package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/miekg/dns"
)

type upstream struct {
	Net            string `fig:"net" default:"udp"`
	Address        string `fig:"address"`
	SkipCertVerify bool   `fig:"skipCertVerify" default:"false"`
}

func (upstream *upstream) Resolve(req *dns.Msg) (ok bool, res *dns.Msg, rtt time.Duration) {
	var client dns.Client
	if upstream.Net == "tcp-tls" {
		client = dns.Client{Net: upstream.Net, Timeout: 5 * time.Second, TLSConfig: &tls.Config{InsecureSkipVerify: upstream.SkipCertVerify}}
	} else {
		client = dns.Client{Net: upstream.Net, Timeout: 5 * time.Second}
	}
	if result, rtt, err := client.Exchange(req, upstream.Address); err != nil {
		log.Printf("[error]\tresolve: %v\tupstream: %v\treason: \"%v\"", req.Question[0].Name, upstream.Address, err.Error())
		return false, nil, time.Microsecond * 0
	} else {
		log.Printf("[success]\tresolve: %v\tupstream: %v\trtt: %v", req.Question[0].Name, upstream.Address, rtt)
		return true, result, rtt
	}
}
