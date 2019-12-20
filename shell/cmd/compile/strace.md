# strace
system call trace(strace,系统调用追踪)是基于 ptrace(2)系统调用的一款工具. strace 通过在一个循环中使用 PTRACE_SYSCALL 请求来显示运行中程序的系统调用(也称为 syscalls)活动相关的信息以及程序执行中捕捉到的信号量.

## 示例
```
# strace -T -tt -e trace=all -p 26143
# strace /bin/ls –o ls.out // 跟踪一个程序
# strace –p <pid> -o daemon.out // 附加到一个现存的进程上
# cat /proc/${pid}/stack # 查看进程调用栈
# gdb # gdb调试
# attach ${pid}
```

strace的输出将会显示每个系统调用的文件描述编号(系统调用会将文件描述符作为参数), 比如`SYS_read(3, buf, sizeof(buf));`, 也可通过`strace –e [read/write]=3 /bin/ls`查看读入/写入到文件描述符 3 中的所有数据. 它还能解析参数和给出由内核返回的结果代码.