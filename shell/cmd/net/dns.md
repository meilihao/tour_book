# dns

## FAQ
### opensuse 15.4配置dns
```bash
# /etc/sysconfig/network/config
...
NETCONFIG_DNS_STATIC_SERVERS="223.5.5.5 114.114.114.114 8.8.8.8"
...
# systemctl restart network
```