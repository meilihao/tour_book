# wireshark
图形化的网络协议分析工具. 捕获文件格式是libpcap格式.

## FAQ
### [filter expression](https://www.wireshark.org/docs/wsug_html_chunked/ChWorkBuildDisplayFilterSection.html)
```
ip.addr in {192.168.0.245 192.168.0.83}
ip.addr==192.168.0.245 || ip.addr==192.168.0.83
ip.src==192.168.0.245 || ip.dst==192.168.0.83
```
