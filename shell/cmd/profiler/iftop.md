# iftop
查看主机的网络带宽使用情况, 如果未指定该接口, 则显示所有网络流量, 并按主机对显示当前带宽使用情况表.

参数说明:
- TX：发送流量
- RX：接收流量
- TOTAL：总流量
- cum：运行iftop到目前时间的总流量
- peak：流量峰值
- rates：分别表示过去 2s 10s 40s 的平均流量

## example
```bash
iftop -i eth0 -n # 监控某网卡
iftop -i eth0 -n -P # 同时显示是什么服务
```