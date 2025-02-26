#依赖镜像
FROM golang:1.21-alpine

#作者信息
MAINTAINER "cy"

# 配置模块代理
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#工作目录
WORKDIR /opt
ADD .  /opt

#在Docker工作目录下执行命令
RUN go build -o main ./main.go

#暴露端口
EXPOSE 8091

#执行项目的命令
CMD ["/opt/main"]