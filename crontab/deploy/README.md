# 使用方式
- 1,Linux操作系统发行版 centos 7.x
- 2,复制两个service文件到/usr/lib/systemd/system
- systemctl daemon-reload
- 3,启动master: systemctl enable go-crond-master.service && systemctl start go-crond-master.service 
- 4,启动worker: systemctl enable go-crond-worker.service && systemctl start go-crond-worker.service
- 5,查看启动状态: systemctl status go-crond-master.service && systemctl status go-crond-worker.service 