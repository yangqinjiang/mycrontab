package lib

import (
	"encoding/json"
	"io/ioutil"
	logs "github.com/sirupsen/logrus"
)

//程序配置
type Config struct {
	//api server
	ApiPort         int `json:"apiPort"`
	ApiReadTimeout  int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`

	//etcd
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`

	WebRoot string `json:"webroot"`

	MongodbUri            string `json:"mongodbUri"`
	MongodbConnectTimeout int    `json:"mongodbConnectTimeout"`
	MongodbUsername       string `json:"mongodbUsername"`
	MongodbPassword       string `json:"mongodbPassword"`

	LogsProduction bool `json:"LogsProduction"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	logs.Info("InitConfig 读取配置文件")
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
	return

}
