# 指定镜像源
FROM golang:1.23-alpine
LABEL authors="Vdebu"

# 设置国内镜像源用于go的依赖下载
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn
# 设置工作目录
WORKDIR /

# 复制依赖清单
COPY go.mod go.sum ./

# 下载所有的依赖项
RUN go mod download

# 将项目所有文件拷贝到工作目录
COPY . .

# 编译可执行文件(使用makefile)
RUN go build -o=./bin/api ./cmd/api

# 声明服务器运行时对外暴漏的端口
EXPOSE 3939

# 运行API
CMD ["sh", "-c", "echo 'Starting backend...' && ./bin/api -docker"]