module Worker

go 1.12

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190804053845-51ab0e2deafa

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4

replace golang.org/x/net => github.com/golang/net v0.0.0-20190724013045-ca1201d0de80

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.22.1

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190805222050-c5a2fd39b72a

replace honnef.co/go/tools => github.com/dominikh/go-tools v0.0.1-2019.2.2

replace golang.org/x/mod => github.com/golang/mod v0.1.0

require (
	github.com/astaxie/beego v1.12.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/yangqinjiang/mycrontab v0.0.0-20190806045921-c9e9e6b7d3d7
)
