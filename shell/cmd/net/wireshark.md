# wireshark
## install
ref:
- [How to Install Wireshark on Ubuntu 22.04 ](https://www.cherryservers.com/blog/install-wireshark-ubuntu)

```bash
# ubuntu 20.04
sudo add-apt-repository universe # 24.04上不用这句
sudo apt install wireshark
# 其他源
sudo add-apt-repository ppa:wireshark-dev/stable -y
wireshark --version

tcpdump -nn -i eth0 icmp and host 172.16.25.161  -X -vvv -w icmp.pcap # 用wireshark打开icmp.pcap即可
```

## 分析数据包
- [tcpdump/wireshark 抓包及分析](https://arthurchiao.art/blog/tcpdump-practice-zh/)

filter:
- arp
- icmp
- ip.addr == 172.16.25.161
- ether host 80:f6:2e:ce:3f:00 : 过滤目标或源地址是xxx的数据包
- ether dst host 80:f6:2e:ce:3f:00 : 过滤目标地址是xxx的数据包
- ether src host 80:f6:2e:ce:3f:00 : 过滤源地址是xxx的数据包
- eth.addr==52:54:00:44:c6:e3 : 过滤目标或源地址是xxx的数据包
- eth.src : 过滤源地址是xxx的数据包
- eth.dst : 过滤目标地址是xxx的数据包
- not ssh and not tcp: 排除ssh, tcp协议

## FAQ
### 允许非root用户拦截数据包
`sudo usermod -aG wireshark <username>`