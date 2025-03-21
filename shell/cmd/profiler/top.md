# top
参考:
- [如何使用 top 命令](https://diabloneo.github.io/2019/08/29/How-to-use-top-command/)
- [Linux top 命令里的内存相关字段（VIRT, RES, SHR, CODE, DATA）](https://liam.page/2020/07/17/memory-stat-in-TOP/)
- [计算Linux系统的CPU利用率](https://bravey.github.io/2019-03-31-%E8%AE%A1%E7%AE%97Linux%E7%B3%BB%E7%BB%9F%E7%9A%84CPU%E5%88%A9%E7%94%A8%E7%8E%87)

> 默认刷新频率: 3s.

> 默认排序列是`%CPU`

> 将各个进程中的RSS值相加后,一般都会超出整个系统的内存消耗,这是因为RSS中包含了各个进程之间的共享内存

top 的界面其实有两个模式：full-screen mode(全屏模式) 和 alternate-display mode(多窗口模式)，默认是全屏模式.

full-screen mode分为两个部分：summary area(上面的汇总信息) 和 tasks area(下面的进程信息).

alternate-display mode有一个 summary area 和 4 个 window. 4 个窗口从上到下编号，对应 1 到 4. summary area 默认为 window 1 的 summary area，可按 a （从上到下） 或者 w （从下到上）切换为不同 window 的 summary area，这个可以从 summary area 左上角的提示看出来. 这个提示区会展示 summary area 当前对应的 window 的编号和名字，比如 1:Def.

top 默认定义了 4 个 window，名字和编号分别是：
1. Def
2. Job
3. Mem
4. Usr

top 其实有两个概念是混合使用的，就是 window 和 Field Group，可以认为它们就是等价的，这里使用 window 这个术语. top 默认定义了 4 个 window，每个 window 配置要显示的列，就称为这个 window 的 Field Group/ 在全屏模式下，可以通过快捷键 g 切换当前要显示的 Field Group 是哪个，使用序号 1 到 4 到进行选择. 修改当前展示的 Field Group 不仅是修改了要显示的内容，样式也会使用 window 的配置，包括过滤器和颜色配置等.

在多窗口模式下，快捷键 g 就对应到切换 window，效果和 a w 是一样的. 在多窗口模式下，快捷键 G 用来修改 window 的名字.

### 界面的滚动和定位

在任何一个 window 下，或者全屏模式下，top 都支持如下快捷键来进行界面的滚动：

    Up or alt+k
    Down or alt+j
    Left
    Right
    PgUp
    PgDown
    Home
    End

top 支持显示一个定位信息，可以看到自己现在处于第几行，第几个字段，通过快捷键 C 切换. y 表示当前的第一行是第几个 task 以及一共有几个 task；x 表示当前的第一列是第几列，以及一共有多少列.

## 描述

实时动态地查看系统的整体运行情况，是一个综合了多方信息监测系统性能和运行信息的实用工具.

第一行(uptime的输出):
1. 当前时间
1. 已运行的时间
1. 目前有多少终端用户登录
1. 过去 1，5 和 15 分钟内的 CPU 负载(TASK_RUNNING + TASK_UNINTERRUPTIBLE进程总和的平均值). 这不是规范化的，所以负载均值为 1 意味着单个 CPU 的满负载，但是在 4 个 CPU 的系统上，这意味着它有 75% 空闲时间.

> 负载是就绪状态等待CPU调用的进程数量统计

第二行(Tasks):
1. 进程总数
1. 当前正在执行的进程数
1. 当前正在睡眠的进程数
1. 被停止的进程数(例如使用CTRL + Z)
1. 僵死的进程数

第三行(CPU):
1. 用户进程的cpu百分比
1. 系统(内核)进程的cpu百分比
1. 用户进程空间内改变过优先级的进程的cpu百分比
1. 空闲的cpu百分比
1. 等待 I/O 的cpu百分比
1. 处理硬件中断花费的时间
1. 处理软件中断花费的时间
1. 由管理程序hypervisor从这个虚拟机`偷走`的时间，用于其他任务(例如启动另一个虚拟机)

第四行(Mem):
1. 物理内存总量
1. 完全空闲的物理内存
1. 已使用的物理内存
1. 用作buff和cache的内存

第五行(Swap):
1. swap总量
1. 空闲的swap
1. 已使用的swap
1. 系统实际可用的空闲内存, 来自`/proc/meminfo`的MemAvailable, 是系统的估算值，表示可用于启动新程序的物理内存大小（不包括 swap 空间）

第六行(进程信息, 可通常快捷键f查看所有列):
- PID : 进程id/标识符
- PPID : 父进程id
- USER : 任务所有者的有效用户名
- RUSER : real user name
- UID : 进程所有者的uid
- GROUP : 进程所有者的组名
- TTY : 启用进程的终端名. 不是从终端启动的进程则显示为`?`
- PR : 任务的调度优先级, 是内核实际使用的任务的优先级，范围是 0 到 39，映射到内核的值是 100 到 139；也可以是 rt ，表示实时任务. 内核表示一个 task 的优先级的范围是 0 到 139. 其中，0 到 99 是实时进程的优先级，100 到 139 是非实时进程的优先级. PR 的默认值是 20，对应到内核是 120，和 NI 的关系是: PR = 20 + NI
- NI : 任务的nice值. 范围是 -20 到 +19，默认值是 0，值越低优先级越高
- P : 最后使用的cpu. 仅在多cpu下有意义
- VIRT : 任务使用的虚拟内存总量(=SWAP+RES). 它包括所有代码，数据和共享库, 交换出的分页以及已被映射但未被使用的页面.
- SWAP : 进程使用的虚拟内存中被换出的大小, 单位kb.
- RES : 常驻内存(进程使用的, 未被换出的物理内存)大小(=CODE+DATA), 即任务已使用的物理内存
- CODE : 代码段占用的物理内存大小, 单位kb
- DATA : 可执行代码以外的部分(数据段+栈)占用的物理内存大小, 单位kb
- SHR : 任务使用的共享内存量. 它只是反映可能与其他进程共享的内存.
- S : 任务的状态可以是以下之一：D=不可中断的睡眠，R=运行，S=睡眠，T=跟踪或停止，Z=僵尸
- %CPU : 上次更新到现在的CPU时间占用百分比
- %MEM : 任务当前使用的可用物理内存的百分比
- TIME :	进程使用的CPU时间总计, 单位秒
- TIME+ : 任务使用CPU的总时间, 单位是1/100s 
- COMMAND命令 : 命令行或程序名称
- WCHAN : 若该进程在睡眠, 则显示睡眠中的系统函数名
- Flags : 任务标志, 参考 sched.h
- nFLT : 页面错误次数
- nDRT : 最后一次写入到现在, 被修改过的页面(脏页)数
- nMaj: Major Page Faults. 该 task 遇到的 Major Page Faults 的数量. Major Page Fault 是指需要访问通过访问 swap 分区或者硬盘（mmap 一个文件，但是还没把内容读取到内存时）来处理的缺页异常.
- nMin: Minor Page Faults. 该 task 遇到的 Minor Page Faults 的数量. Minor Page Fault 是指不需要通过访问 swap 分区或者硬盘来处理的缺页异常，简单的说，没用到磁盘 I/O.
- vMj: Major Faults delta. 上次 top 刷新数据依赖的 Major Page Faults 增加数量
- vMn: Minor Faults delta. 上次 top 刷新数据依赖的 Minor Page Faults 增加数量
- P : cpu id
- NU : numa id

## 选项
- -b : batch模式, 重定向不乱码. 分屏显示输出信息, 即每次刷新算一屏
- -c : 显示进程的整个命令路径参数
- -d N : 整个界面更新的秒数，默认是5秒
- -i : 不显示闲置或僵死的进程
- -n : top输出的更新次数然后退出, 比如`top -n 1 -b > top-output.txt`
- -p PID : 监控指定的进程
- -s : 使top在安全模式下运行(即禁用交互式指令被取消)
- -u <user>: 指定用户

## 快捷键(交互式命令)
- h/? : help
- A : 切换全屏或多窗口模式
- 1 : 切换显示(总体/每个逻辑)cpu状态
- 2 : 按NUMA node显示
- 3 : 显示指定NUMA node下的所有cpu
- b : 切换粗体/反色显示, 仅对x/y有效
- c : **切换显示完整命令行和命令名称**
- d/s : 改变top的刷新间隔. 默认是秒, 小数会换成毫秒, 0是不断刷新. 间隔过短会导致来不及看刷新及系统负载增大
- e : 底下的进程信息切换，每次切换转换率为1000，切换的单位也是 k,m,g,t,p
- E : 顶部的内存信息切换，每次切换转换率为1000，只是没有单位，切换的单位为 k,m,g,t,p
- f/F : 从当前显示列表中添加或删除项. 按f后进入列管理界面(顶部**第一行会显示当前排序列**, 其他是help), 按上下键选择列, 空格选择是否显示; 选择一项, 按方向键右，则可以将标注栏扩大，再按方向键上和下，可以移动标注对象, 按方向键左则取消移动操作; 选择一项, 按s设置排序列.
- H : 显示线程, 默认是进程
- k : 终止一个进程, 系统会提示输入一个指定的pid来kill
- l : 显示/不显示平均负载和启动时间
- L : 搜索

    搜索不考虑字段，会在 tasks area 显示的所有内容中找到指定的行. & 表示跳转到下一个匹配项. 搜索关键字会被高亮，也可使用 PgUp 和 PgDown 等来调整展示的内容
- o/O : 按**指定列**过滤进程信息

    过滤条件的格式: `[!]FLD?VAL`
    - ! 表示非，不是必选字段，其他都是必填字段
    - FLD 表示字段名字，即按f后显示的列, 是大小写敏感的，要完全一致
    - ? 代表一个要求的操作符(=, <, >):
        - < 或者 > 表示字符串比较，或者数值比较. 注意比较时的数字的单位问题，100.0m 会比 1.0g 大，所以需要先进行单位切换（例如内存的单位切换，e）
        - = 表示部分匹配

    可以使用多个条件进行过滤，每个按一次 o/O 输入一个条件，所有条件会进行与操作.可以按 Ctrl-o 显示所有条件.

    举些例子：

        过滤 nice 值小于 0 的任务，即以减号开头的 nice 值：NI=-
        过滤所有 postgres 进程：COMMAND=postgres

- P : 根据CPU使用百分比大小进行排序
- q : quit
- m : 切换显示(是否显示)内存信息
- M : 根据驻留内存大小进行排序
- N : 按pid排序
- i : 使top不显示任何闲置或者僵死进程
- I : 切换 %CPU 的模式，Irix/Solaris 两个模式。Irix mode 的计算方式是跑满一个 CPU 为 100%，%CPU 可能会超过 100%。Solaris mode 则是会把总体利用率除以 CPU 核数，保证不会超过 100%
- r : 重新设置一个进程的优先级,系统会要求输入pid和要设置的优先级
- R : 切换排序方向
- t : 显示/不显示进程和cpu状态
- S : 切换到累计模式. Cumulative time，off 展示两次刷新时间的即时值，而不是从进程启动到现在的累加值
- T : 按时间/累计时间排序
- V : 显示进程父子关系. 这个模式下无法按照字段排序.
- W : 将当前top设置`~/.toprc`
- x : 将排序列高亮. `shift + >`或`shift + <`可以向右或左改变排序列
- u/U : 仅显示指定用户, 系统要求输入用户名称
- y : 将运行状态是R的task高亮
- Y : top 有一个模式Inspect Mode, 需提前定义inspect entries(在`~/.toprc`中top设置的下方)，否则无效. 它可以允许执行指定的命令来查看**一个**任务的更多信息, 需指定pid. 参考[这里](https://diabloneo.github.io/2019/08/29/How-to-use-top-command/).
- z : 改名进程信息的颜色
- Z : 进入颜色定制界面, 便于清晰展示内容


## example
```bash
# top -n 1 # 仅执行一次, 且分屏
# top -bn 1 # 完成显示(不分屏)top信息1次
# top -d 2 -n 1 -b # 输出重定向
# top -p <pid> # 仅显示指定进程, 适合排查numa下cpu问题, 再通过F选择P和NU列即可
# top -H -p <pid> # 显示指定进程的线程
```

## FAQ
### 进程字段排序
先按键盘`x`, 可打开排序列的高亮效果.

- 按CPU占用率排序: `shift + p`
- 按内存占用率排序: `shift + m`
- 按CPU占用时间排序: `shift + t`
- 按PID排序: `shift + n`

### 显示线程
`shift + h`

### 等式
不一定有等式 CODE + DATA = RES 成立，但一定成立等式 ANON(在堆上分配的内存) = RES - SHR 及不等式 ANON <= DATA (vm_physic) <= DATA. 如果观察到程序稳定运行时 RES - SHR 不断增长，则可能预示着程序存在内存泄漏现象

### 保存top当前配置
按下 shift + w 键将当前设置保存为默认配置(` ~/.toprc`), 下次启动时，top将按当前配置显示