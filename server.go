package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kkyr/fig"
	"github.com/miekg/dns"
)

type server struct {
	Net       string     `fig:"net" default:"udp"`
	Listen    string     `fig:"listen" default:"0.0.0.0:53"`
	Upstreams []upstream `fig:"upstreams" default:"[]"`
}

func MakeServer() (server server) {
	fig.Load(&server,
		fig.File("setting.json"),
		fig.Dirs("/etc/pure-dns", "."),
	)
	return
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
		go func(upstream upstream) {
			if ok, res, rtt := upstream.Resolve(req); ok {
				identRR, _ := dns.NewRR(fmt.Sprintf("%s TXT %s://%s ttl:%s", "dns.provider", upstream.Net, upstream.Address, rtt.String()))
				res.Answer = append(res.Answer, identRR)
				c <- res
			}
		}(item)
	}
	select {
	case result := <-c:
		result.SetReply(req)
		return true, result
	case <-time.After(time.Duration(time.Second * 5)):
		emptyMsg := dns.Msg{}
		emptyMsg.SetReply(req)
		return false, &emptyMsg
	}
}
