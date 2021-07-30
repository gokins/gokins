FROM golang:1.16.6-alpine3.14 AS builder
# ENV GOPROXY=https://goproxy.cn,direct
# RUN apk add git build-base && git clone https://gitee.com/gokins/gokins.git /build
RUN apk add git build-base && git clone https://github.com/gokins/gokins.git /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gokins main.go


FROM alpine:latest AS final

ENV GOKINS_WORKPATH=/data/gokins

RUN apk --no-cache add openssl ca-certificates curl git wget \
    && rm -rf /var/cache/apk \
    && mkdir -p /app /data/gokins

COPY --from=builder /build/bin/gokins /app
WORKDIR /app
ENTRYPOINT ["/app/gokins"]