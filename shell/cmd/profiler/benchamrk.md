# benchmakr
tools:
- spec cup 2006 : cpu
- stream : 内存带宽
- fio : 模拟特定io工作负载
- bonnie++ : 测试硬盘驱动器性能/文件系统性能
- netpref : 测试网络各个方面的性能

	可用`pmstat -p ALL`查看cpu

## netpref
`netpref -t TCP_STREAM -H 10.0.0.161 -p 1234 -l 60`
参数:
- t testname: 指定测试类型, 支持TCP_STREAM, UDP_STREAM, TCP_RR, TCP_CRR, UDP_RR
- H : 被测试端运行的netserver的ip
- p : 被测试端运行的netserver的port
- l ： 指定测试的时间, 单位秒