export GO111MODULE=on
export GOPROXY=https://goproxy.io
echo "building......"
go build . && echo "------> build ok <------"
ls -l
echo "run ./worker"