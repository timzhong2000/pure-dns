# Pure-DNS

## 优点

### 多协议

支持提供tcp和udp的dns服务，并且支持多上游配置

### 加速dns

从多个上游选择最快的一个响应。

### 可追溯

所有dns请求都会返回一个dns.provider的TXT记录，携带了本次请求选择的上游和本次请求的耗时

### 体积小

安装upx后，使用 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" . && upx -9 pure-dns` 命令编译并压缩，可执行文件小于2MB。

## 配置

配置文件路径 `/etc/pure-dns/setting.json`

```json
{
  "net": "tcp",             // "udp" | "tcp"
  "listen": "0.0.0.0:5353", // "<ip>:<port>"
  "upstreams": [
    {
      "net": "udp",
      "address": "119.29.29.29:53"
    },
    {
      "net": "tcp",
      "address": "8.8.8.8:53"
    }
  ]
}
```
