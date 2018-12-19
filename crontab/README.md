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
