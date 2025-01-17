# time
Go语言标准库提供的timer实质上是由Go runtime自行维护的,而不是操作系统级的定时器资源.

Go运行时启动了一个单独的goroutine,该goroutine执行了一个名为timerproc的函数,维护了一个“最小堆”. 该goroutine会被定期唤醒并读取堆顶的timer对象,执行该timer对象对应的函数(向timer.C中发送一条数据,触发定时器),执行完毕后就会从最小堆中移除该timer对象.

创建一个time.Timer实则就是在这个最小堆中添加一个timer对象实例,而调用timer.Stop方法则是从堆中删除对应的timer对象.