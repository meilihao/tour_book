# ps

## 描述

查看系统的进程状态

推荐监控间隔: 5m.

## 选项

- -a : 显示所有终端机下执行的程序,除没有关联到终端的进程和session leaders.
- a : 显示现行终端机下的所有程序，包括其他用户的程序
- -A : 显示所有程序
- e : 列出程序时，显示每个程序所使用的环境变量
- -e : 同`-A`
- -f : 显示UID,PPIP,C与STIME栏位
- -l : 采用详细的格式来显示程序状况
- u : 以用户为主的格式来显示程序状况
- -U <用户名/Uid> : 列出属于该用户的程序的状况(不含参数)
- -u : 同`-U`
- U <用户名/Uid> : 列出属于该用户的程序的状况(含参数)
- x : 显示所有程序，不以终端机来区分
- -x : 同`x`

## 例

    # ps lax // 比`ps aux`实用
    # ps aux
    # ps -elf

## 扩展

- [session leaders](http://www.win.tue.nl/~aeb/linux/lk/lk-10.html#ss10.3)

Process ID (PID)

This is an arbitrary number identifying the process. Every process has a unique ID, but after the process exits and the parent process has retrieved the exit status, the process ID is freed to be reused by a new process.
Parent Process ID (PPID)

This is just the PID of the process that started the process in question.
Process Group ID (PGID)

This is just the PID of the process group leader. If PID == PGID, then this process is a process group leader.
Session ID (SID)

This is just the PID of the session leader. **If PID == SID, then this process is a session leader**.

### 字段
- %CPU : 正在使用的CPU百分比
- %MEM : 正在使用的实际内存百分比
- VSZ : 进程分配的虚拟内存, 包括进程可以访问的所有内存，包括进入交换分区的内容，以及共享库占用的内存
- RSS : 常驻内存集（Resident Set Size），表示该进程分配的**内存大小**. 因此不包括进入交换分区的内存, 但包括共享库占用的内存（只要共享库在内存中）,包括所有分配的栈内存和堆内存.
- TTY : 控制终端的id
- Time : 进程消耗的cpu时间
- WCHAN : 进程正在等待的资源
- SZ : 进程在主存的大小(按页面数计算)
- NI : nice值
- PRI : 调度优先级(内核的内部表示,与nice不同)

### 进程状态

linux上进程有5种状态:

- 运行(正在运行或在运行队列中等待)
- 中断(休眠中, 受阻, 在等待某个条件的形成或接受到信号)
- 不可中断(休眠中,收到信号不唤醒和不可运行, 进程必须等待直到有中断发生)
- 僵死(进程已终止, 但进程描述符存在, 直到父进程调用wait4()系统调用后释放)
- 停止(进程收到SIGSTOP, SIGSTP, SIGTIN, SIGTOU信号后停止运行运行)

ps工具标识进程的5种状态码:

- D 不可中断睡眠 uninterruptible sleep (usually IO),进程繁忙或挂起，不响应信号，通常是等待io
- R 正在运行或可运行 (on run queue)
- I TASK_INTERRUPTIBLE：进程处于睡眠状态，正在等待某些事件发生, 进程可以被信号中断. 接收到信号或被显式的唤醒呼叫唤醒之后，进程将转变为 TASK_RUNNING 状态
- S 可中断的sleeping,例如，终端进程和 Bash 通常处于此状态，等待你键入某些内容
- T 由任务控制信号停止
- t 在跟踪期间由调试器停止
- X	死亡（不应该看到)
- Z 僵死 a defunct (”zombie”) process,这种情况发生在错误终止的进程上
- W 进程被交换出去 (自2.6.xx内核以来无效）

附加标记:
- < 高优先级
- N 低优先级
- L	将页面锁定到内存中（用于实时和自定义 IO)
- s 是会话领导。Linux 中的相关进程被视为一个单元，并具有共享会话 ID（SID）。如果进程 ID（PID）= 会话 ID（SID），则此进程将是会话领导。
- l	是多线程的（使用 CLONE_THREAD，例如 NPTL pthreads）
- +	位于前台进程组,这样的处理器允许输入和输出到tty
