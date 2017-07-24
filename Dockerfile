FROM hub.c.163.com/library/golang:1.8.3

# 在容器内创建一个目录来存储我们的应用,使他成为工作目录
RUN mkdir -p /go/src/toon-web
RUN mkdir -p /go/src/golang.org
WORKDIR /go/src/toon-web

# 复制 toon-web目录到容器中
COPY . /go/src/toon-web

# 下载并安装第三方依赖到容器中
RUN git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net
RUN git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys
RUN git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools
RUN git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto
RUN ln -s $GOPATH/src/github.com/golang $GOPATH/src/golang.org/x


RUN go-wrapper download
RUN go-wrapper install

# 设置Port 环境变量
ENV PORT 1323

# 给主机暴漏 1323端口
EXPOSE 1323

# 告诉Docker启动容器运行的命令
CMD ["go-wrapper", "run"]
