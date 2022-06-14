package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/kkyr/fig"
	"github.com/miekg/dns"
	"github.com/yl2chen/cidranger"
)

type server struct {
	Net             string     `fig:"net" default:"udp"`
	Listen          string     `fig:"listen" default:"0.0.0.0:53"`
	Timeout         int        `fig:"timeout" default:"1000"`
	Upstreams       []upstream `fig:"upstreams" default:"[]"`
	BlackList       []string   `fig:"blackList" default:"[]"`
	blackListRanger cidranger.Ranger
}

func MakeServer() (ok bool, server server) {
	err := fig.Load(&server,
		fig.File("setting.json"),
		fig.Dirs("/etc/pure-dns", "."),
	)
	if err != nil {
		log.Print(err.Error())
		log.Printf("Load config failed!")
		return false, server
	}
	server.blackListRanger = cidranger.NewPCTrieRanger()
	for _, blockCidrString := range server.BlackList {
		if _, blockCidr, err := net.ParseCIDR(blockCidrString); err == nil {
			server.blackListRanger.Insert(cidranger.NewBasicRangerEntry(*blockCidr))
		} else {
			log.Printf("%s is not a CIDR format string. Please set blacklist string like 192.168.1.0/24", blockCidrString)
		}
	}
	return true, server
}

func (s *server) ListenAndServe() {
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		_, res := s.Resolve(req)
		w.WriteMsg(res)
	})
	server := &dns.Server{Addr: s.Listen, Net: s.Net}
	log.Printf("Starting DNS server on %s://%s", s.Net, s.Listen)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}

func (server *server) Resolve(req *dns.Msg) (ok bool, res *dns.Msg) {
	c := make(chan *dns.Msg)
	for _, item := range server.Upstreams {
		go server.resolve(item, req, c)
	}
	select {
	case result := <-c:
		result.SetReply(req)
		return true, result
	case <-time.After(time.Duration(server.Timeout) * time.Millisecond):
		emptyMsg := dns.Msg{}
		emptyMsg.SetReply(req)
		return false, &emptyMsg
	}
}

func (server *server) resolve(upstream upstream, req *dns.Msg, c chan *dns.Msg) {
	if ok, res, rtt := upstream.Resolve(req); ok {
		for _, answer := range res.Answer {
			if dnsARecord, ok := answer.(*dns.A); ok {
				if contains, _ := server.blackListRanger.Contains(dnsARecord.A); contains {
					log.Printf("[block]\tresolve: %s\trtt: %v\tanswer: %s\tupstream: %v://%v", dnsARecord.Hdr.Name, rtt, dnsARecord.A, upstream.Net, upstream.Address)
					return
				}
			}
			if dnsAAAARecord, ok := answer.(*dns.AAAA); ok {
				if contains, _ := server.blackListRanger.Contains(dnsAAAARecord.AAAA); contains {
					log.Printf("[block]\tresolve: %s\trtt: %v\tanswer: %s\tupstream: %v://%v", dnsAAAARecord.Hdr.Name, rtt, dnsAAAARecord.AAAA, upstream.Net, upstream.Address)
					return
				}
			}
		}
		identRR, _ := dns.NewRR(fmt.Sprintf("%s TXT %s://%s ttl:%s", "dns.provider", upstream.Net, upstream.Address, rtt.String()))
		res.Answer = append(res.Answer, identRR)
		c <- res
	}
}
