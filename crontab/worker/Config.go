package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//程序配置
type Config struct {

	//etcd
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`

	MongodbUri       string `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
	MongodbUsername string `json:"mongodbUsername"`
	MongodbPassword string `json:"mongodbPassword"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {

	//1,读取配置文件
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	conf := &Config{}
	//反序列化
	err = json.Unmarshal(content, conf)
	if err != nil {
		return
	}
	//赋值单例
	G_config = conf
	fmt.Println("配置文件:", G_config)
	return

}
