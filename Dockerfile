# 编译阶段
FROM golang:1.18 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .

# 部署阶段
FROM scratch

WORKDIR /app

COPY --from=builder /app/pure-dns ./pure-dns

# 以后上tls-dns需要证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert

EXPOSE 53

ENTRYPOINT ["./pure-dns"]