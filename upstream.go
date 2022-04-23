package main

import (
	"log"
	"time"

	"github.com/miekg/dns"
)

type upstream struct {
	Net     string `fig:"net" default:"udp"`
	Address string `fig:"address"`
}

func (upstream *upstream) Resolve(req *dns.Msg) (ok bool, res *dns.Msg, rtt time.Duration) {
	client := dns.Client{Net: upstream.Net, Timeout: time.Second * 5}
	if result, rtt, err := client.Exchange(req, upstream.Address); err != nil {
		log.Printf("%v 解析失败", upstream.Address)
		return false, nil, time.Microsecond * 0
	} else {
		log.Printf("%v 解析成功 rtt:%v", upstream.Address, rtt)
		return true, result, rtt
	}
}