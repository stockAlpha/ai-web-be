# 基础镜像，基于golang的alpine镜像构建--编译阶段
FROM golang AS builder
# 作者
MAINTAINER stock
# 全局工作目录
WORKDIR /go/stock-web-be
# 把运行Dockerfile文件的当前目录所有文件复制到目标目录
COPY . /go/stock-web-be
# 环境变量
#  用于代理下载go项目依赖的包
RUN go env -w GOPRIVATE=github.com/stockAlpha
RUN git config --global url."https://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@github.com/".insteadOf "https://github.com/"
# swagger重新生成
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN sh swag.gen
# 编译，关闭CGO，防止编译后的文件有动态链接，而alpine镜像里有些c库没有，直接没有文件的错误
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go


# 使用alpine这个轻量级镜像为基础镜像--运行阶段
FROM alpine AS runner
# 全局工作目录
WORKDIR /go/stock-web-be
# 复制编译阶段编译出来的运行文件到目标目录
COPY --from=builder /go/stock-web-be/main .
# 复制编译阶段里的config文件夹到目标目录
COPY --from=builder /go/stock-web-be/conf ./conf
# 拷贝所有静态文件
COPY --from=builder /go/stock-web-be/disk ./disk
# 需暴露的端口
RUN if [ -z "$PORT" ]; then \
        export PORT=8080; \
    fi
EXPOSE $PORT
# 设置环境变量
RUN if [ -z "ENV" ]; then \
        export ENV="prod"; \
    fi
ENV ENV=prod
# docker run命令触发的真实命令(相当于直接运行编译后的可运行文件)
ENTRYPOINT ["./main"]