# cat Dockerfile
# docker build --network host -t ccapi:v0.1 .
FROM alpine:latest

# Set timezone to Shanghai
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone

# Create a directory named 'apps' in the root directory of the container
WORKDIR /apps

# Mount container directories
# VOLUME ["/apps/conf"]

# Copy the executable file 'ccapi' from the current directory
COPY ./apps/ccapi /apps/ccapi

# Copy configuration files to the container
# COPY conf/config.toml /apps/conf/config.toml

# Set encoding
# ENV LANG C.UTF-8
ENV TZ Asia/Shanghai
ENV PORT 8080
ENV GIN_MODE release

# Expose ports
EXPOSE 8080

# Command to run the Golang program
ENTRYPOINT ["/apps/ccapi"]
