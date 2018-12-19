# 使用方式
- 1,Linux操作系统发行版 centos 7.x
- 2,复制两个service文件到/usr/lib/systemd/system
- 3,启动master: systemctl enable crond-master.service && systemctl start crond-master.service 
- 4,启动worker: systemctl enable crond-worker.service && systemctl start crond-worker.service
- 5,查看启动状态: systemctl status crond-master.service && systemctl status crond-worker.service 