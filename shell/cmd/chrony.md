# chrony
```sudo
sudo apt/yum remove ntp
sudo apt/yum install chrony
vim /etc/chrony/chrony.conf
sudo systemctl restart chronyd
sudo systemctl enable chronyd
chronyc sources # 查看时间源
chronyc tracking 使用 chronyc 查询同步状态
```

## 配置上游ntp
删除`/etc/chrony.conf`中原先的`server xxx iburst`, 再加入自己的`server xxx iburst`, 再执行`systemctl restart chronyd`即可.

> iburst表示启用快速同步

## FAQ
### `sudo systemctl enable chronyd`报`Failed to enable unit: Refusing to operate on alias name or linked unit file: chronyd.service`
`systemctl enable chrony`

### chrony.conf的pool和server区别
- server : 指定一个或多个具体的 NTP 服务器

    手动指定，每个服务器都独立列出, 可逐个设置选项，比如 minpoll, maxpoll

    内部固定服务器	（例如企业内部 NTP 服务器）
- pool :  指定一个 NTP 池（域名），解析后得到多个服务器

    由 DNS 返回的多个服务器组成，适合动态管理

    公共池（更可靠、自动管理）
### 未彻底关闭ntpd，双服务抢占导致时间疯狂跳变
ntpd和Chronyd绝对不能共存

```bash
# 停止并禁用ntpd服务
systemctl stop ntpd
systemctl disable ntpd
# 屏蔽开机自启、防止残留生效
systemctl mask ntpd
```

### 防火墙未放行123/UDP端口，同步成功假象误导排查
时钟同步协议只走123/UDP端口

```bash
# firewalld放行
firewall-cmd --add-service=ntp --permanent
firewall-cmd --reload
# iptables放行
iptables -A INPUT -p udp --dport 123 -j ACCEPT
```

### Chrony服务端未配置allow网段，内网客户端全部同步失败
ntpd默认允许所有内网设备接入同步, 但Chrony默认关闭所有外网、内网访问权限, 即默认不会对外提供时钟同步服务，所有客户端请求都会被拒绝, 需配置allow网段(`/etc/chrony.conf`)

### 服务器时间偏差过大，Chrony拒绝自动校正时间
NTP支持大幅度跳变校时，Chrony默认渐进式校时，需手动开启makestep强制校时

```
#偏差1秒以上、连续3次检测即强制校时
makestep 1.0 3
```

### 虚拟机/云主机同步硬件时钟，引发时间漂移加剧
虚拟机硬件时钟由宿主机统一管控，虚拟机内部无法直接读写硬件时钟，强行配置只会冲突报错(`can't sync hardware clock`)

```
# 注释原有硬件时钟配置
# hwclockfile /etc/adjtime
# 开启内核时钟同步适配虚拟机
rtcsync
```

### 双上游NTP源无权重区分，Chrony动态择优导致频繁切源、时间跳变
配置主源prefer权重，强制固定主次，杜绝无序切换。并且设置仅保留1个主源，其它的作为备用，避免出现频繁选举。同时，放宽择优条件，降低选主的敏感度

```bash
# 主时钟，固定优先锁定，长期作为唯一同步源
server 上层时钟1的IP地址 iburst prefer
# 备用时钟，仅主源故障、断连时自动接管
server 上层时钟2的IP地址 iburst
# 仅保留1个有效主源，另一台只做校验备份，不参与择优切换
maxsources 1
# 收紧择优距离，关闭激进自动换源逻辑
reselectdist 0.5
```