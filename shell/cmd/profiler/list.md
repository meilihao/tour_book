# list
-[Red Hat Enterprise Linux8监控和管理系统状态和性能](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/8/html/monitoring_and_managing_system_status_and_performance)

## perf
- perf tools : 综合性能概要分析工具
- ftrace : 追踪kernel的函数调用
- systemtap
- [tracy](https://mp.weixin.qq.com/s/VWdMizmmVlf-7Bsd0AkRzQ)

## 网络
- nuttcp : 带宽吞吐量测试工具
- netperf : 带宽测试工具
- iperf3 : servers间的带宽测速

	```bash
	# --- 服务端执行
	iperf -s
	# --- 客户端执行
	iperf -c <server_ip>
	```
## 内存
- memtester
- memtest86+

## windows
PAL(Performance Analysis of logs) 是一个分析perfmon计数器日志的实用工具.