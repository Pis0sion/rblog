# 打包依赖阶段使用golang作为基础镜像
FROM golang:alpine3.15 as builder

# 启用go module
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

ADD . ./

# 指定OS等，并go build
RUN go build -ldflags="-w -s" -o rblog cmd/rblog/rblog.go

# 由于我不止依赖二进制文件，还依赖config文件夹下的配置文件
# 所以我将这些文件放到了publish文件夹
RUN mkdir publish && cp rblog publish && \
    cp -r config publish

# 运行阶段指定scratch作为基础镜像
FROM alpine

WORKDIR /app

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/publish .

# 指定运行时环境变量
ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["./rblog"]
