# gc
## [gc log](https://godoc.org/runtime)
参考:
- [GODEBUG之gctrace解析](http://cbsheng.github.io/posts/godebug%E4%B9%8Bgctrace%E8%A7%A3%E6%9E%90/)

调试GC打开: `GODEBUG=gctrace=1`, log:
```
scvg: 0 MB released
scvg: inuse: 3, idle: 59, sys: 63, released: 58, consumed: 4 (MB)
gc 252 @4316.062s 0%: 0.013+2.9+0.050 ms clock, 0.10+0.23/5.4/12+0.40 ms cpu, 16->17->8 MB, 17 MB goal, 8 P
```
scvg: 将gctrace设置为大于0的任何值时使得gc在内存释放回系统时输出的摘要信息. 将内存释放回系统的过程称为清理(scavenging).
scvg解读:
- inuse: 在使用或在部分使用的span大小
- idle: 等待清理的span大小
- sys: 系统映射的内存
- released: 释放回系统的内存
- consumed: 从系统申请的内存

gc解读:
- gc 252 ： 这是第252次gc
- @4316.062s ： 这次gc的markTermination阶段完成后距离runtime启动到现在的时间, 即程序启动以来的秒数
- 0% ：到目前为止，gc的标记工作（包括两次mark阶段的STW和并发标记）所用的CPU时间占总CPU的百分比
- 0.013+2.9+0.050 ms clock ：按顺序分成三部分，0.013表示mark阶段的STW时间（单P的）; 2.9表示并发标记用的时间（所有P的）; 0.050表示markTermination阶段的STW时间（单P的）
- 0.10+0.23⁄5.4⁄12+0.40 ms cpu ：gc占用cpu时间,按顺序分成三部分，0.10表示整个进程在mark阶段STW停顿时间(0.013*8)；0.23⁄5.4/12有三块信息，0.23是mutator assists占用的时间，5.4是dedicated mark workers+fractional mark worker占用的时间，12是idle mark workers占用的时间. 这三块时间加起来会接近2.9*8；0.40 ms表示整个进程在markTermination阶段STW停顿时间(0.050 * 8), 8是P的个数.
- 16->17->8 MB ：按顺序分成三部分，16表示开始mark阶段前的heap_live大小；17表示开始markTermination阶段前的heap_live大小；8表示live heap的大小
- 17 MB goal：表示下一次触发GC的内存占用阀值是17MB，等于8MB * 2，向上取整
- 8 P ：本次gc共有多少个P