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