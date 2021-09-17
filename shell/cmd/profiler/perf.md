# perf
参考:
- [Linux性能分析工具Perf简介](https://segmentfault.com/a/1190000021465563)
- [利用perf剖析Linux应用程序](https://blog.gmem.cc/perf)
- [Linux性能优化实战学习笔记：第四十九讲](https://www.cnblogs.com/luoahong/p/11577395.html)

## perf report
Children/Self: 如果在record时收集了调用链, 则Overhead可以在Children、Self两个列中显示. Children显示子代函数的样本计数、Self显示函数自己的样本计数.

## FAQ
### 软中断ksoftirqd/n 占用CPU 过高排查
```bash
# perf top # 全局查看
# perf top -p 45558 # 查看进程pid=45558的性能信息
# perf record -a -g -p 9 -- sleep 30 # 9为高ksoftirqd的pid
# perf report
```