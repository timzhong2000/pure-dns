package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	dnsproxy "github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
)

type upstream struct {
	/*
		support list:
		- udp
		- tcp
		- tcp-tls
		- https
		- quic
		- sdns
	*/
	Net     string `fig:"net" default:"udp"`
	Address string `fig:"address"`
	/*
	 mode:
	 - hybrid: use miekg/dns for udp/tcp/tls and use dnsproxy for https/quic/sdns
	 - dnsproxy: use dnsproxy only
	*/
	Mode string `fig:"mode" default:"hybrid"`
}

func (upstream *upstream) getLegacyExchanger() func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
	client := dns.Client{Net: upstream.Net, Timeout: 5 * time.Second}
	return client.Exchange
}

func convertToDnsProxyNetType(net string) string {
	if net == "tcp-tls" {
		return "tls"
	}
	return net
}

func (upstream *upstream) getDnsProxyExchanger() func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
	return func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
		var client dnsproxy.Upstream
		startTime := time.Now()
		if client, err = dnsproxy.AddressToUpstream(fmt.Sprintf("%v://%v", convertToDnsProxyNetType(upstream.Net), upstream.Address), &dnsproxy.Options{Timeout: 5 * time.Second}); err != nil {
			return nil, time.Since(startTime), err
		}
		if r, err = client.Exchange(m); err != nil {
			return nil, time.Since(startTime), err
		}
		return r, time.Since(startTime), nil
	}
}

func (upstream *upstream) GetExchanger() func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
	if upstream.Mode == "dnsproxy" {
		return upstream.getDnsProxyExchanger()
	} else if upstream.Mode == "hybrid" {
		switch upstream.Net {
		case "":
			fallthrough
		case "udp":
			fallthrough
		case "tcp":
			fallthrough
		case "tcp-tls":
			return upstream.getLegacyExchanger()
		case "https":
			fallthrough
		case "quic":
			fallthrough
		case "sdns":
			return upstream.getDnsProxyExchanger()
		default:
			log.Printf("wrong net type %v. expect one of the following type: \"udp\" \"tcp\" \"tcp-tls\" \"https\" \"quic\" \"sdns\"", upstream.Net)
			panic("wrong net type")
		}
	} else {
		log.Printf("wrong mode %v. expect hybrid or dnsproxy", upstream.Mode)
		panic("wrong mode")
	}
}

func getAnswer(raw string) string {
	return strings.Split(raw, "\t")[4]
}

func (upstream *upstream) Resolve(req *dns.Msg) (ok bool, res *dns.Msg, rtt time.Duration) {
	if result, rtt, err := upstream.GetExchanger()(req, upstream.Address); err != nil {
		log.Printf("[error]\tresolve: %v\tupstream: %v://%v\treason: \"%v\"", req.Question[0].Name, upstream.Net, upstream.Address, err.Error())
		return false, nil, time.Microsecond * 0
	} else {
		log.Printf("[success]\tresolve: %v\tresult: %v\trtt: %v\tupstream: %v://%v", req.Question[0].Name, getAnswer(result.Answer[0].String()), rtt, upstream.Net, upstream.Address)
		return true, result, rtt
	}
}
