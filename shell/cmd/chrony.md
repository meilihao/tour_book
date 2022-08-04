# chrony
## 配置上游ntp
删除`/etc/chrony.conf`中原先的`server xxx iburst`, 再加入自己的`server xxx iburst`, 再执行`systemctl restart chronyd`即可.