# shutdown

## 描述

关机/重启系统

原理: 通知所有用户即将关机的信息, 并禁止新用户登录(冻结login命令). 到达指定的关机时间后根据其接收的参数, 发送请求给系统的init进程, 由init将系统运行级调整到参数所指定的状态(0, 关机; 6, 重启)

### 类型
- init
- halt
- poweroff
- reboot

> 在 systemd中, shutdown,half,poweroff,reboot均指向`/bin/systemctl`

## 参数
- -F : 重启时执行fsck
- h : 将关机, 功能上与half命令相当
- -k : 仅发送消息给所用用户, 但不真关机
- -r : 重启
- -c : 取消前一个shutdown命令
- t<second> : 多久后关机

## 例
```sh
$ shutdown -h now # 立即关机
$ shutdown -h +1 # 1m后关机
$ shutdown -r now # 立即重启
$ shutdown -r 11:00 # 在11:00时重启系统
```
