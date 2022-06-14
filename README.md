# Pure-DNS

## 优点

### 多协议

支持提供tcp和udp的dns服务，并且支持dns over udp/tcp/tcp-tls/https(doh)/quic/sdns(DNSCrypt) 上游服务。

### 黑名单模式
支持设定cidr格式的ipv4和ipv6黑名单，避免网络提供商返回的明显错误的dns应答。

### 加速响应

从多个上游选择最快的一个响应。

### 易调试

所有dns请求都会返回一个dns.provider的TXT记录，携带了本次请求选择的上游和本次请求的耗时

### 体积小

安装upx后，使用 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" . && upx -9 pure-dns` 命令编译并压缩，可执行文件小于3MB。

## 用法

### Linux

1. 在右侧 `Releases` 根据系统和cpu架构下载编译好的二进制文件
2. 解压并且复制到 `/usr/bin` 下
``` bash
tar -xzf ./pure-dns_0.1.7_linux_amd64.tar.gz && sudo mv pure-dns /usr/bin/pure-dns && sudo chmod +x /usr/bin/pure-dns
```
3. 参考下面配置（要删除//注释的部分），创建或修改 `/etc/pure-dns/setting.json` 配置文件
4. 启动服务器
``` bash
$ pure-dns
2022/04/30 00:00:00 Starting DNS server on udp://0.0.0.0:5353
```

### Windows
1. 在右侧`Releases`根据系统和cpu架构下载编译好的二进制文件，并且解压
2. 参考下面配置（要删除//注释的部分），在 `pure-dns.exe` 相同目录下创建 `setting.json`
3. 在当前目录打开命令行界面，执行 `./pure-dns.exe`

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
  ],
  "blackList": [
    "192.168.16.0/24",    // 必须是CIDR格式
    "fe80::0/128"         // 也支持ipv6
  ]
}
```

## 未来的新功能
- [ ] 统计功能
- [x] 支持提供tls服务
- [x] 对于不可信dns上游限制可返回的ip段

## 不会添加的功能
- dns缓存和域名分流。因为本项目只是个单纯的dns转发器，开多个pure-dns实例然后使用dnsmasq作为缓存和分流是更好的选择。