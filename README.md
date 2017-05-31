# autoiss-go
automatic generate config.json for shadowsocks-go/shadowsocks-local cmd according to iss  website.
inspired by  github.com/ystyle/autoss-go

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
        server address (default "ss.ishadowx.com")
```

## Reference 参考
* Chrome浏览器设置 [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega/releases)
* GFW配置[gfwlist](https://github.com/gfwlist/gfwlist)
