# 使用 Go 镜像作为基础镜像
FROM harbor.wlj.com/demo/golang1.17 AS build

# 设置 GOPROXY
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 将项目文件复制到工作目录
COPY . .

# 编译 Go 项目
RUN go build -o k8s-go-demo .

# 使用轻量的 Alpine 镜像作为最终镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从第一个阶段复制二进制文件到最终镜像
COPY --from=build /app/k8s-go-demo .

# 设置二进制文件可执行权限
RUN chmod +x k8s-go-demo

# 启动应用程序
CMD ["./k8s-go-demo"]