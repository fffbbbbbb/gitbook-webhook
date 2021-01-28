FROM golang:alpine AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o webhook

FROM alpine

RUN echo -e 'https://mirrors.aliyun.com/alpine/v3.12/main/\nhttps://mirrors.aliyun.com/alpine/v3.12/community/' > /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /build /app

ENTRYPOINT ["sh","-c","/app/webhook"]