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
autoiss-go -s <server> -l <local-port>
    server     - ishadowsocks server address, default "ss.ishadowx.com"
    local-port - sock5  proxy local-port,default 1080
```

## Reference 参考
* Chrome浏览器设置 [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega/releases)
* GFW配置[gfwlist](https://github.com/gfwlist/gfwlist)
