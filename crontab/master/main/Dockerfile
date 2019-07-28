# 使用这个镜像编译代码
FROM golang:1.12.7

ARG buildno
ARG gitcommithash

RUN echo "Build number: $buildno"
RUN echo "Based on commit: $gitcommithash"
WORKDIR /go/src/app
COPY . .
#下载代码
RUN go get -d -v ./...

RUN go install -v ./...

CMD ["app"]