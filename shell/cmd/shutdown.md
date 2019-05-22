# shutdown

## 描述

关机/重启系统

原理: 通知所有用户即将关机的信息, 并禁止新用户登录. 到达指定的关机时间后根据其接收的参数, 发送请求给系统的init进程, 将系统运行级调整到参数所指定的状态(0, 关机; 6, 重启)

### 类型
- init
- halt
- poweroff
- reboot

> 在 systemd中, shutdown,half,poweroff,reboot均指向`/bin/systemctl`

## 参数

## 例
```sh
$ shutdown -h now # 立即关机
$ shutdown -h +1 # 1m后关机
$ shutdown -r now # 立即重启
$ shutdown -r 11:00 # 在11:00时重启系统
```
