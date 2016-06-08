# ps

## 描述

查看系统的进程状态

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

    # ps aux
    # ps -ef

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

### 进程状态

linux上进程有5种状态:

- 运行(正在运行或在运行队列中等待)
- 中断(休眠中, 受阻, 在等待某个条件的形成或接受到信号)
- 不可中断(休眠中,收到信号不唤醒和不可运行, 进程必须等待直到有中断发生)
- 僵死(进程已终止, 但进程描述符存在, 直到父进程调用wait4()系统调用后释放)
- 停止(进程收到SIGSTOP, SIGSTP, SIGTIN, SIGTOU信号后停止运行运行)

ps工具标识进程的5种状态码:

- D 不可中断 uninterruptible sleep (usually IO)
- R 运行 runnable (on run queue)
- S 中断 sleeping
- T 停止 traced or stopped
- Z 僵死 a defunct (”zombie”) process
