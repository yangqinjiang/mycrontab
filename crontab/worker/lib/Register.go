package lib

import (
	"context"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"github.com/coreos/etcd/clientv3"
	"net"
	"sync"
	"time"
	"github.com/yangqinjiang/mycrontab/worker/lib/config"
)

var (
	G_register   *Register
	onceRegister sync.Once
)

//注册节点到etcd /cron/worker/ip地址
type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string //本地IP
}

//注册到/cron/workers/ip
func (register *Register) keepOnLine() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse

		cancelCtx  context.Context
		cancelFunc context.CancelFunc
	)
	for {
		//注册路径
		regKey = common.JOB_WORKER_DIR + register.localIP

		cancelFunc = nil
		//10s 租约
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY //租约失败,不断重试
		}
		//自动续租
		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}
		//续租失败,可取消
		cancelCtx, cancelFunc = context.WithCancel(context.TODO())
		//注册到etcd,put不关心应答
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}
		//处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { //续租失败
					goto RETRY
				}

			}
		}

	RETRY: //租约失败,不断重试
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}

func InitRegistr() (err error) {
	onceRegister.Do(func() {

		//初始化配置
		//读取配置文件
		config := clientv3.Config{
			Endpoints:   config.G_config.EtcdEndpoints, //集群地址
			DialTimeout: time.Duration(config.G_config.EtcdDialTimeout) * time.Microsecond,
		}
		//建立连接
		client, err := clientv3.New(config)

		if err != nil {
			return
		}
		//得到Kv和Lease的API子集
		kv := clientv3.NewKV(client)
		lease := clientv3.NewLease(client)

		localIp, err := getLocalIP()

		if err != nil {
			return
		}
		//赋值单例
		G_register = &Register{
			client:  client,
			kv:      kv,
			lease:   lease,
			localIP: localIp,
		}
		//服务注册
		go G_register.keepOnLine()
	})
	return
}

//获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {

	var (
		addrs   []net.Addr
		ipNet   *net.IPNet //IP地址
		isIpNet bool
	)
	//获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	//取第一个非Io的网卡IP
	for _, addr := range addrs {
		//ipv4,ipv6
		//反解,这个网络是IP地址

		//过滤环回地址
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			//再跳过IPv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() //192.168.1.x
				return                   //找到ipv4, 直接返回
			}

		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUND //没有找到网卡
	return
}
