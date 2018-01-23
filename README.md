# autoiss-go

automatic generate config.json for shadowsocks-go/shadowsocks-local cmd according to iss  website.
inspired by  github.com/ystyle/autoss-go

自动获取 https://global.ishadowx.com/ free shadowsock 配置，自动运行shadowsocks-local

## log 开发日志

* 01.23.2018 change site to [ishadowx](https://global.ishadowx.net)
* 11.26.2017 change site to [ishadowx](https://go.ishadowx.net)

## Install 安装

```shell
go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-local
go get github.com/corbamico/autoiss-go
autoiss-go
```

## Usage  使用

```shell
autoiss-go -s <server> -l <local-port> -n <index of servers>
  -d    print debug message
  -l int
        local socks5 proxy port (default 1080)
  -n int
        which shadowsocks server to use.
        (0:first one,-1:last one)
  -s string
        server address (default "go.ishadowx.net")
```

## Reference 参考

* Chrome浏览器设置 [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega/releases)
* GFW配置[gfwlist](https://github.com/gfwlist/gfwlist)
