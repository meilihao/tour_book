## base
- [Redis开发运维实践指南](https://www.gitbook.com/book/gnuhpc/redis-all-about)
- [你确定 Redis 是单线程的进程吗？](https://www.tuicool.com/articles/Q7NR326)

	Redis 单线程指的是「接收客户端请求->解析请求 ->进行数据读写等操作->发生数据给客户端」这个过程是由一个线程（主线程）来完成的.

	在 Redis 6.0 版本之后，也采用了多个 I/O 线程来处理网络请求 ， 这是因为随着网络硬件的性能提升，Redis 的性能瓶颈有时会出现在网络 I/O 的处理上. 所以为了提高网络请求处理的并行度，Redis 6.0 对于网络请求采用多线程来处理, 但是对于读写命令，Redis 仍然使用单线程来处理.

## 持久化
- [redis持久化机制](http://shanks.leanote.com/post/Untitled-55ca439338f41148cd000759-22)

## ha
- [keepalived实现redis双主备份](https://blog.51cto.com/huangzhijun/1725606)
- [keepalived+redis 高可用redis主从解决方案](https://developer.aliyun.com/article/524588)

## 复制
- [Redis主从同步与故障切换，有哪些坑？](https://new.qq.com/omn/20201125/20201125A0GFNT00.html)
- [+redis 主从同步-slave端](https://www.jianshu.com/p/e10d21ecdd0b)
- [boazjohn/redis-keepalived](https://github.com/boazjohn/redis-keepalived)