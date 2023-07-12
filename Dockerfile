# 使用Go官方镜像作为基础镜像
FROM golang:1.20-alpine as builder

ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 将go mod下载的依赖复制到容器的缓存中
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将项目文件复制到容器内
COPY . .

# 编译Go代码生成可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o noname .

# 创建一个新的阶段，用于生成较小的最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# EXPOSE 8080

# 设置容器启动命令
CMD ["./noname"]
