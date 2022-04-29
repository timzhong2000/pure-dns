# Pure-DNS

## 优点

### 多协议

支持提供tcp和udp的dns服务，并且支持dns over udp/tcp/tcp-tls/https(doh)/quic/sdns(DNSCrypt) 上游服务。

### 加速响应

从多个上游选择最快的一个响应。

### 易调试

所有dns请求都会返回一个dns.provider的TXT记录，携带了本次请求选择的上游和本次请求的耗时

### 体积小

安装upx后，使用 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" . && upx -9 pure-dns` 命令编译并压缩，可执行文件小于2MB。

## 配置

配置文件路径 `/etc/pure-dns/setting.json`

```json
{
  "net": "udp",                     // "udp" | "tcp"
  "listen": "0.0.0.0:53",           // "<ip>:<port>"
  "timeout": 2000,                  // 2000ms
  "upstreams": [
    {
      "net": "udp",                 // "udp" | "tcp" | "tcp-tls" | "https" | "quic" | "sdns"
      "address": "8.8.8.8:53",
      "mode": "dnsproxy"            // optional(default hybrid). when using hybrid mode the upstream of udp/tcp/tcp-tls is provided by miekg/dns
    },
    {
      "net": "tcp",
      "address": "8.8.8.8:53",
      "mode": "hybrid"
    },
    {
      "net": "tcp-tls",
      "address": "1.1.1.1:853"
    },
    {
      "net": "https",
      "address": "dns.adguard.com/dns-query"
    },
    {
      "net": "quic",
      "address": "dns.adguard.com:853"
    },
    {
      "net": "sdns",
      "address": "AQIAAAAAAAAAFDE3Ni4xMDMuMTMwLjEzMDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20"
    }
  ]
}
```

## 未来的新功能
- [ ] 统计功能
- [ ] 支持提供tls服务
- [ ] 对于不可信dns上游限制可返回的ip段

## 不会添加的功能
- dns缓存和域名分流。因为本项目只是个单纯的dns转发器，开多个pure-dns实例然后使用dnsmasq作为缓存和分流是更好的选择。