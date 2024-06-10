FROM wjqserver/caddy:alpine

RUN mkdir -p /data/www
RUN mkdir -p /data/ipinfo/db
RUN mkdir -p /data/ipinfo/log
