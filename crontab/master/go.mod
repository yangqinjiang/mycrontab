module github.com/yangqinjiang/mycrontab/master

go 1.12

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190804053845-51ab0e2deafa

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4

replace golang.org/x/net => github.com/golang/net v0.0.0-20190724013045-ca1201d0de80

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.22.1

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190805222050-c5a2fd39b72a

replace honnef.co/go/tools => github.com/dominikh/go-tools v0.0.1-2019.2.2

replace golang.org/x/mod => github.com/golang/mod v0.1.0

replace golang.org/x/text => github.com/golang/text v0.3.2

replace go.mongodb.org/mongo-driver => github.com/mongodb/mongo-go-driver v1.0.4

require (
	github.com/coreos/etcd v3.3.12+incompatible
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/sirupsen/logrus v1.4.2
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.0.1
	google.golang.org/grpc v0.0.0-00010101000000-000000000000 // indirect
)
