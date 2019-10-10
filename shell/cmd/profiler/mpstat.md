# mpstat
与vmstat类似, 但能够给出针对单个处理器的监控结果. 对调试SMP(symmetric multiprocessing, 对称多处理机技术)的软件极为有用, 因为它能显示出使用多处理器的系统效率的高低.

部分列:
- %guest : 若干cpu运行一个虚拟处理器所花费的百分比
- %gnice : 若干cpu运行一个nice的虚拟机所花费的百分比

## example
```sh
$ mpstat -P ALL
$ mpstat -P 0 1 2 # 在cpu 0上每个1s采样一次, 共采样2次
```
