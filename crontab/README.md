#库安装,使用go mode 管理依赖模块
`export GO111MODULE=on` / `set GO111MODULE=on`

###  goproxy 是一个开源项目，当用户请求一个依赖库时，如果它发现本地没有这份代码就会自动请求源，然后缓存到本地，用户就可以从 goproxy.io 请求到数据

`export GOPROXY=https://goproxy.io`
``
## 
# 课后练习
- web增加任务超时配置项
- worker支持超时检查
- 任务超时在日志中得到体现

# master选主
- 启动后抢占etcd乐观锁/cron/master
- 抢到锁的,成为leader并持续续租

# 任务失败告警
- worker任务失败向etcd的/cron/warn/{job}标识告警
- leader master 监听/cron/warn目录变化
- etcd 性能不高, 可用队列优化

# 热更新
- https://segmentfault.com/a/1190000008487440
# 程序配置方案
- https://github.com/koding/multiconfig