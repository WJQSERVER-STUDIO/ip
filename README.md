# ip

![pull](https://img.shields.io/docker/pulls/wjqserver/ip.svg) ![size](https://img.shields.io/docker/image-size/wjqserver/ip)

## Demo

[演示站](https://ip.1888866.xyz)

## 部署

```
docker run -d -p 80:80 -v ./ipinfo/db:/data/ipinfo/db -v ./ipinfo/log:/data/ipinfost/log --restart always wjqserver/ip:latest
```

需导入mmdb数据库才可正常使用自有API
