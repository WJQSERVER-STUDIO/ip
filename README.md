# ip

![GitHub Release](https://img.shields.io/github/v/release/WJQSERVER-STUDIO/ip?display_name=tag&style=flat)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/WJQSERVER-STUDIO/ip)
![pull](https://img.shields.io/docker/pulls/wjqserver/ip.svg)
![size](https://img.shields.io/docker/image-size/wjqserver/ip)

用于获取IP相关信息的API程序

## Demo

[演示站](https://ip.1888866.xyz)

## 部署

```
docker run -d -p 8980:80 -v ./ipinfo/db:/data/ipinfo/db -v ./ipinfo/log:/data/ipinfo/log --restart always wjqserver/ip:latest
```

需导入mmdb数据库(ipinfo)才可正常使用自有API
