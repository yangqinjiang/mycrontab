# 使用这个镜像编译代码
FROM golang:1.12.7

WORKDIR /go/src/app
COPY . .
#下载代码
RUN go get -d -v ./...

RUN go install -v ./...

CMD ["app"]
# 运行
#$ docker build -t my-golang-app .
#$ docker images|grep my-golang-app                       
# 输出 my-golang-app / latest / 43c9e3802f10 /  16 seconds ago /   816MB
#$ docker run -it --rm --name my-running-app my-golang-app
# 输出 `hello world`