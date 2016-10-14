# firewall-cmd

## 描述

firewalld的字符界面管理工具.

## 选项:
- --zone=NAME          # 指定 zone
- --permanent          # 永久修改，--reload 后生效
- --timeout=seconds    # 持续效果，到期后自动移除，用于调试，不能与 --permanent 同时使用

## 参考
- [CentOS 7 下使用 Firewall](https://havee.me/linux/2015-01/using-firewalls-on-centos-7.html)

## 例

### 添加端口
```
// 永久打开一个端口
firewall-cmd --permanent --zone=public --add-port=80/tcp
// 永久打开某项服务
firewall-cmd --permanent --zone=public --add-service=http
// 需重新加载防火墙
firewall-cmd --reload
```

### 检查是否生效
```
firewall-cmd --zone=public --query-port=80/tcp
```

### 列出所有的开放端口
```
firewall-cmd --list-all
```

### 删除端口
```
// 永久关闭一个端口
firewall-cmd --zone=public --remove-port=80/tcp
// 永久关闭某项服务
firewall-cmd --zone=public --remove-service=http
```

### 进行端口转发
```
firewall-cmd --permanent --add-forward-port=port=80:proto=tcp:toport=8080:toaddr=192.0.2.55
```

### 查看
```
// 查看已被激活的 Zone 信息
firewall-cmd --get-active-zones
// 查看指定接口的 Zone 信息
firewall-cmd --get-zone-of-interface=eth0
// 查看指定级别的接口
firewall-cmd --zone=public --list-interfaces
// 查看指定级别的所有信息,比如public
firewall-cmd --zone=public --list-all
// 查看所有级别被允许的信息
firewall-cmd --get-service
// 查看重启后所有 Zones 级别中被允许的服务，即永久放行的服务
firewall-cmd --get-service --permanent
```

## 扩展

### 运行、停止、禁用firewalld
```
# systemctl start firewalld
# systemctl status firewalld 或者 firewall-cmd –state
# systemctl disable firewalld
# systemctl stop firewalld
```
