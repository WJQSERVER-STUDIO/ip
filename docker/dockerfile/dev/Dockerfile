FROM alpine:latest AS builder

ARG USER=WJQSERVER-STUDIO
ARG REPO=ip
ARG APPLICATION=ip
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

# 拉取依赖
RUN apk add --no-cache wget curl

# 创建目录
RUN mkdir -p /data/www
RUN mkdir -p /data/${APPLICATION}/config 
RUN mkdir -p /data/${APPLICATION}/log
RUN mkdir -p /data/${APPLICATION}/db

# 后端
RUN VERSION=$(curl -s https://raw.githubusercontent.com/${USER}/${REPO}/main/DEV-VERSION) && \
    wget -O /data/${APPLICATION}/${APPLICATION} https://github.com/${USER}/${REPO}/releases/download/$VERSION/${APPLICATION}-${TARGETOS}-${TARGETARCH}

# 权限
RUN chmod +x /data/${APPLICATION}/${APPLICATION}

FROM alpine:latest

COPY --from=builder /data/www /data/www
COPY --from=builder /data/${APPLICATION} /data/${APPLICATION}
RUN mkdir -p /data/${APPLICATION}/db

# 权限
RUN chmod +x /data/${APPLICATION}/${APPLICATION}

CMD ["/data/ip/ip"]


