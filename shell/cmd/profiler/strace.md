# [strace](https://linuxtools-rst.readthedocs.io/zh_CN/latest/tool/strace.html)

跟踪进程执行时的系统调用和所接收的信号.

> strace的负载很高.

类似的不同作用的trace还有:
- ltrace: 库文件调用跟踪器

  strace也能查看系统调用, 但它查看的并不是最终的系统调用, 而是系统调用的封装函数.
- ptrace: 进程跟踪器
- ftrace: 包含一系列跟踪器，用于不同的场合，比如跟踪内核函数调用(function tracer)、跟踪上下文切换(sched_switch tracer)、查看中断被关闭的时长(irqsoff tracer)、跟踪内核中的延迟以及性能问题等. 它是内建于Linux的内核跟踪工具，依赖于内核配置宏(Kernel Hacking->Tracers)和debugfs.

## 选项
- -c : 统计每一系统调用的所执行的时间,次数和出错的次数等, 并在退出时汇总信息
- -d : 在stderr上输出strace自身的调试信息
- -f : 跟踪当前进程fork所产生的子进程
- -ff : 如果提供-o filename,则所有进程的跟踪结果输出到相应的filename.pid中,pid是各进程的进程号. 与`-c`冲突.
- -i : 输出系统调用的入口指针
- -q : 禁止显示有关attaching(附加), detaching(分离)等消息
- -r : 打印每个syscall的相对间
- -t : 在输出中的每一行前加上时间信息
- -tt : 在输出中的每一行前加上时间信息,微秒级
- -T : 显示每一调用所耗的时间.
- -v : 详细模式. 比如未缩写的argv,stat,termios等参数
- -x : 以十六进制形式输出非ASCII字符串
- -xx : 所有字符串以十六进制形式输出
- -p pid : 指定进程
- -o filename : 指定信息输出位置
- -s STRLEN : 设置打印的字符串最大长度
- -E var : 运行时移除环境变量var
- -E var=val : 运行时设置环境变量var
- -e : 使用过滤表达式, 格式`option=[!]all or option=[!]val1[,val2]...`, 具体见man.
    - -e trace=set
      只跟踪指定的系统调用.例如:-e trace=open,close,rean,write表示只跟踪这四个系统调用, 默认的为set=all
    - -e trace=file
      只跟踪有关文件操作的系统调用
    - -e trace=process
      只跟踪有关进程控制的系统调用
    - -e trace=network
      跟踪与网络有关的所有系统调用
    - -e trace=signal
      跟踪所有信号有关的系统调用
    - -e trace=ipc
      跟踪所有ipc有关的系统调用
    - -e trace=desc
      跟踪所有描述符有关的系统调用
    - -e abbrev=set
      设定strace输出的系统调用的结果集. `-v`等同于`abbrev=none`.默认为abbrev=all
    - -e verbose=set
      对特定系统调用集合解除参照结构.
    - -e raw=set
      对特定系统调用集合显示原始的未解码的参数, 并以十六进制显示该参数
    - -e signal=set
      指定跟踪的系统信号.默认为all.如 signal=!SIGIO(或者signal=!io),表示不跟踪SIGIO信号
    - -e read=set
      输出从指定文件中读出的数据.例如:`-e read=3,5`
    - -e write=set
      输出写入到指定文件中的数据

## example
```
# strace command arg ...
# strace -e trace=openat pkg-config --cflags systemd # 仅输出syscall=openat的记录
```