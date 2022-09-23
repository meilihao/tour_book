# ps

## 描述

查看系统的进程的瞬间状态

推荐监控间隔: 5m.

## 选项

- -a : 显示所有终端机下执行的程序,除没有关联到终端的进程和session leaders.
- a : 显示所有用户的程序，包括每个程序的完整路径
- -A : 显示所有程序
- c : 只显示进程的名称, 不显示它的完整路径
- -C : 筛选指定进程名的进程
- e : 列出程序时，显示每个程序所使用的环境变量
- -d : 选择除session leader以外的所有进程
- -e : 同`-A`
- -f : 显示UID,PPIP,C与STIME栏位
- -l : 采用详细的格式来显示程序状况
- -L : 显示线程(LWP, light-weight-process)
- -o : 显示指定列(包含pid列), 比如`pid,user,cgroup,args`

    - lwp ： 线程id
    - psr ： 当前分配给进程运行的处理器编号
    - ruser ： 运行进程的用户
    - args : 运行的命令及其参数
    - psr: cpu id
    - cgroup : 所属cgroup, 信息来自`/proc/<pid>/cgroups`
- -O : 显示指定列(不包含pid列), 类同`-o`
- -T : 显示线程
- u : 以用户为主的格式来显示程序状况
- -U <用户名/Uid> : 列出属于该用户的程序的状况(不含参数)
- -u : 同`-U`
- U <用户名/Uid> : 列出属于该用户的程序的状况(含参数)
- x : 显示所有系统程序，包括那些没有终端的程序
- -x : 同`x`
- r : 只选择正在运行的进程
- X : 显示寄存器信息
- -l : 长格式
- v : 显示虚拟内存格式

## 例

    # ps lax // 比`ps aux`实用
    # ps aux
    # ps -elf
    # ps auxf # 输出进程树
    # ps -A -O pgid # 基本属性+`-O`指定的属性
    # ps -C smbd --no-header
    # ps xawf -eo pid,user,cgroup,args # 查看cgroup
    # ps -eLo psr | grep -e "^[[:blank:]]*3$" | wc -l # 显示线程使用的cpu id并筛选出id=3的线程并汇总个数
    # ps -eLo ruser,pid,ppid,lwp,psr,args | awk '{if($5==3) print $0}' # 显示运行在cpu id=3上的线程
    # ps aux | awk 'NR>0 {$6=int($5/1024)"M";}{ print;}' | column -t # 查看进程rss内存(MB)
    # ps aux | awk '{mem += $6} END {print mem/1024/1024}'
    # ps -f --ppid 2 -p 2 # 查找2的子进程
    # ps -ef | grep "\[.*\]" # 查找内核进程
    # ps -o pid,psr,comm -p <pid> # 进程/线程目前分配到的 （在“PSR”列）CPU ID
    # ps -T -p <pid> 显示指定进程的线程
    # ps -p <pid> # 获取指定进程信息, 可检查进程是否alive

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
- STAT : 进程状态
- COMMAND : 进程执行时的命令, 内核态会用`[]`包裹, 而用户态不会.

### 进程状态

linux上进程有5种状态:

- 运行(正在运行或在运行队列中等待)
- 中断(休眠中, 受阻, 在等待某个条件的形成或接受到信号)
- 不可中断(休眠中,收到信号不唤醒和不可运行, 进程必须等待直到有中断发生)
- 僵死(进程已终止, 但进程描述符存在, 直到父进程调用wait4()系统调用后释放)
- 停止(进程收到SIGSTOP, SIGSTP, SIGTIN, SIGTOU信号后停止运行运行)

ps工具标识进程的5种状态码:

- D 不可中断睡眠 uninterruptible sleep (usually IO),进程繁忙或挂起，不响应信号，此时即便用 kill 命令也不能将其中断, 通常是等待io
- R 正在运行或在运行队列中等待(on run queue)
- I TASK_INTERRUPTIBLE：进程处于睡眠状态，正在等待某些事件发生, 进程可以被信号中断. 接收到信号或被显式的唤醒呼叫唤醒之后，进程将转变为 TASK_RUNNING 状态
- S 可中断的sleeping, 当某个条件形成后或者接收到信号时，则脱离该状态. 例如，终端进程和 Bash 通常处于此状态，等待键入某些内容
- T 由任务控制信号停止, 即进程收到停止信号后停止运行
- t 在跟踪期间由调试器停止
- X 死亡（不应该看到)
- Z 僵死进程, 已经终止但是其父进程还没回收
- W 进程被交换出去 (自2.6.xx内核以后无效）

附加标记:
- < 高优先级
- N 低优先级
- L	将页面锁定到内存中（用于实时和自定义 IO)
- s 是会话领导。Linux 中的相关进程被视为一个单元，并具有共享会话 ID（SID）。如果进程 ID（PID）= 会话 ID（SID），则此进程将是会话领导。
- l	是多线程的（使用 CLONE_THREAD，例如 NPTL pthreads）
- +	位于前台进程组,这样的处理器允许输入和输出到tty

# pidof
根据进程名查找pid, 比pgrep准确.