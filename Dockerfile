# cat Dockerfile
# docker build -t ccapi:v0.1 .
# 表示依赖 alpine 最新版
FROM alpine:latest

# 设置时区为上海
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone

# 在容器根目录 创建一个 apps 目录
WORKDIR /apps

# 挂载容器目录
# VOLUME ["/apps/conf"]

# 拷贝当前目录下 go_docker_demo1 可以执行文件
COPY ./apps/ccapi /apps/ccapi

# 拷贝配置文件到容器中
# COPY conf/config.toml /apps/conf/config.toml

# 设置编码
# ENV LANG C.UTF-8
ENV TZ Asia/Shanghai
ENV PORT 8080
ENV GIN_MODE release

# 暴露端口
EXPOSE 8080

# 运行golang程序的命令
ENTRYPOINT ["/apps/ccapi"]
#ENTRYPOINT ["/bin/bash"]
