# gdb

注：如果如果程序编译时开启了优化选项，那么在用GDB调试被优化过的程序时，可能会发生某些变量不能访问，或是取值错误码的情况. 对付这种情况时，需要在编译程序时关闭编译优化. GCC中可以使用`-gstabs`选项来解决这个问题.
## 选项
- -cd : 设置工作目录
- -d : 添加文件查找路径
- -q : 安静模式, 启动时不打印介绍信息和版本信息
- -s : 设置读取的符号表文件
- -tui : 有窗口

## 命令模式
```gdb
$ gdb
(gdb) target remote :1234 # 连接gdb server
(gdb) set [arch/architecture] i386:x86-64:intel # 可直接输入`set arch/architecture`查询可得
(gdb) info reg # 获取寄存器信息
```

## gdb调试工具常用命令
编译程序时需要加上`-g`(为了调试时显示源码)，之后才能用gdb进行调试：`gcc -g main.c -o main`, 或者`GDB yourpram`, 再或者
```bash
$ gdb
(gdb) file yourpram
(gdb) run/r:开始程序的执行; 或`run parameter`将参数传递给该程序
```

gdb中命令：
```
(gdb) 回车键：重复上一命令
(gdb) help：查看命令帮助，简写h

    help <命令> : 查询gdb的具体命令信息
(gdb) show language : 查看当前的语言环境, 默认是c
(gdb) file : 装入需要调试的程序或加载符号表
(gdb) add-symbol-file android_test/system/bin/linker –s .text 0x6fee6180 # 将代码段(.data)段的调试信息加载到0x6fee6180上
(gdb) where : 查看程序出错的地方
(gdb) run：重新开始运行文件（run-text：加载文本文件，run-bin：加载二进制文件）,简写r. **如果在程序运行结束后再执行r, 会重新开始调试该程序**

     run argv[1] argv[2]：调试时命令行传参
(gdb) start：单步执行，运行程序，停在第一执行语句, 简写st
(gdb) layout src : 弹出窗口用于查看代码
(gdb) layout asm : 弹出窗口用于查看汇编. `CTRL+x+a`: 退出
(gdb) list：查看原代码,简写l

    list : 默认显示当前行和之后的10行，再执行又下滚10行
    list n : 从第n行开始查看代码
    list n1,n2 : 显示n1行和n2行之间的代码
    list ,n2 : 显示当前行和n2行之间的代码
    list <function>：查看具体函数
    list + : 向后显示源代码
    list - : 向前显示源代码
    set listsize <count> : 设置一次显示源代码的行数
    show listsize : 查看listsize的值
(gdb) set：设置变量的值

    set environment varname [= value]  : 设置传递给调试程序的环境变量varname的值为value，不指定value时，值默认为NULL
    unset environment varname : 取消环境变量
    set args 5.in.txt : 设置被调试程序的参数. `gdb ./a.out 5.in.txt`, **5.in.txt是gdb参数而不是a.out**. 或直接使用`gdb --args ./a.out 5.in.txt`
    set 变量名=表达式  : 可以修改变量的值
    set $cs = 0xf000 : 设置cs寄存器
    set arch i386:x86-64 : 设置arch, 不正确的arch会导致反编译出错. 因为i8086是16bit指令, i386是32bit, i386:x86-64是64bit.
    set print address on : 打开地址输出，当程序显示函数信息时，GDB会显出函数的参数地址
    set print array on  : 设置数组显示方式，打开后当显示数组时，每个元素占一行，否则每个元素则以逗号分隔
    set print elements : 设置数组的最大展示长度，当到达这个长度时，GDB就不再往下显示了. 如果设置为0，则表示不限制
    set print null-stop : 如果打开了这个选项，那么当显示字符串时，遇到结束符则停止显示, 默认为off
    set print pretty on : 如果打开, 那么当GDB显示结构体时会比较漂亮
    set print union : 设置显示结构体时，是否显式其内的联合体数据
    set print object : 在C++中，如果一个对象指针指向其派生类，如果打开这个选项，GDB会自动按照虚方法调用的规则显示输出，如果关闭这个选项的话，GDB就不管虚函数表了
    set language <lan> : 设在程序语言类型
    set disassemble-next-line on: si时默认显示下一条指令
(gdb) next：单步调试（c语言级别, 逐过程，函数直接执行）,简写n

    nexti/ni : 单步跳过下一条指令(不进入函数) 
(gdb) step：单步调试（c语言级别, 逐语句：跳入自定义函数内部执行）,简写s

    si [N]: 执行N个指令, 默认是1个
(gdb) backtrace：查看函数的调用的栈帧和层级关系,简写bt

    bt <-n> : 只打印栈底下n层信息
(gdb) frame：切换函数的栈帧,简写f
    
    f n : 查看某一层的栈
    up n : 向栈的上面移动
    down n : 下栈的下面移动
(gdb) show : 显示

    - show environment : 显示环境变量
    - show paths : 显示$PATH
    - show args：查看设置好的参数
(gdb) info：查看函数内部局部变量的数值,简写i

    info f : 打印详细的栈信息
    info frame : 查看当前函数的程序语言
    info source : 查看当前文件的程序语言
    info files : 显示被调试文件的详细信息
    info func : 显示被调试程序的所有函数名称
    info prog : 显示被调试程序的执行状态
    info args : 打印当前函数的参数名和value
    info local : 打印当前函数的局部变量及其值
    info var : 打印被调试程序的所有全局和静态变量名称
    info catch : 打印当前函数中的异常信息调用
    info breakpoints/break ：查看当前设置的所有断点
    info registers/reg [<regname>...]: 查看寄存器(包括浮点寄存器)
    info all-registers : 查看寄存器(除了浮点寄存器)
    info $<register_name> : 查看指定寄存器
    info registers eflags : 查看eflags的结果
    info display : 查看display设置的自动显示的信息
    info line <filename>:<function> : 显示指定函数在文件中的行数, 内存的起始地址及终止地址
    info proc m : 查看内存映射
(gdb) examine : 查看内存地址中的值, 简写x

    格式: `x/<n><format><size> ADDRESS`:
    参数: n、f、u是可选的参数:
    - n 是一个正整数，表示读取n个单位的内存
    - format 表示显示的格式，print的格式. 如果地址所指的是字符串，那么格式可以是s; 如果是指令地址，那么格式可以是i; x(hex); d(deciaml), c(char)等
    - size 表示每个单位的字节数，如果不指定的话，GDB默认是4个bytes. u参数可以用下面的字符来代替，b表示单字节，h表示双字节，w表示四字节，g表示八字节. 当指定了字节长度后，GDB会从指内存定的内存地址开始，读写指定字节，并把其当作一个值取出来.

    example:
    - `x /3uh ${addr}` ：比如从内存地址0x54320读取内容, 否则从上次读取的结尾开始，h表示以双字节为一个单位，3表示三个单位，u表示按十六进制显示
    - x/5i $cs*16+$pc <=> display /5i $cs*0x10+$pc  # 基于下一条指令 disassemble 5条汇编指令(包括下一条指令自身)
    - `x/2xb 0x7c00+510`: mbr标志
(gdb) finish：结束当前函数，返回到函数调用点
(gdb) continue：继续运行,简写c

    continue [N] : 使程序在运行过程中忽略该断点num次，就是说在num+1次执行到该断点时才暂停程序的运行，在循环中作用较大
(gdb) print：打印值及地址,简写p

    print 操作符:
    - @ : 是一个和数组有关的操作符
    - `::` : 指定一个在文件或是一个函数中的变量
    - {} : 表示一个指向内存地址的类型为type的一个对象

    print x=8 : 修改变量值
    print 开始表达式@连续打印空间的大小 : 还可以打印出内存的某个部分开始的连续值
    p 'file'::variable : 查看指定文件的变量
    p 'func_name'::variable : 查看指定函数的变量
    p array@len : 打印动态数组, 静态数组可直接print. array:数组的首地址，len:数据的长度. 比如`p *0x6001a8@500`, 0x6001a8为数组地址
    p $<regname> : 查看reg
    p /<format> <name> : 打印格式, 比如`p /x $rax`

        x 按十六进制格式显示变量
        d 按十进制格式显示变量
        u 按十六进制格式显示无符号整型
        o 按八进制格式显示变量
        t 按二进制格式显示变量
        a address
        c 按字符格式显示变量
        f 按浮点数格式显示变量
        s string
        z 填充前导零
        i 汇编指令(机器码)
    print $cs * 0x10 + $pc # 打印real-mode下的下一条指令的地址
(gdb) kill : 终止正在调试的程序,简写k
(gdb) quit：退出gdb,简写q
(gdb) break：设置断点,简写b

    break <function>, c++可用`class::function`指定
    break <linenum>
    break +offset
    break -offset
    b filename:linenum
    b filename:function
    break : 在下一条指令处停在
    break ... if <condition> : 比如`break if i=100`
    b *addr : 在该地址设置断点
    break *_start : 在_start设置断点, 该用法常见于调试as生成的汇编程序
(gdb) tbreak : 设置临时断点, 它会在调用一次后自动删除. 设置规则参考break
(gdb) delete breakpoints num：删除第num个断点,简写d
(gdb) display：追踪查看具体变量值, 简写disp

    格式: display[/i|s] [expression | addr], i表示机器指令码, 即汇编; s,字符串
    **display命令会被gdb记忆，如果打印一个值，后续遇到该值均会被打印出来.**
    
    disable和enalbe不删除自动显示的设置，而只是让其失效和恢复:
    - disable display
    - enable display
(gdb) undisplay：取消先前的display设置，编号从1开始递增
(gdb) watch：被设置观察点的变量发生修改时，打印显示
(gdb) i watch：显示观察点
(gdb) enable breakpoints：启用断点
(gdb) disable breakpoints：禁用断点
(gdb) delete 断点编号 :  直接删除该断点
(gdb) clear 断点所在行号 : 直接删除所在行断点

    delete和clear不同之处是delete跟的是断点号，clear跟的是行号.
(gdb) disassemble : 反汇编, 简写disas

    disassemble main : 反汇编main函数
    disassemble /m addme : 带源码的反汇编addme函数
    set disassembly-flavor intel/att : 切换汇编方言 : GDB默认汇编方言是AT&T格式
    disassemble start,end/start,+length : 反汇编地址范围. length为连续解析的字节长度
    disassemble /r : 可以用16进制形式显示程序的原始机器码

    example:
    - `disas 0x7c00,+10`, 反汇编结果是32bit寄存器格式, 可当前实际是在实模式下, 自己换算成16bit寄存器位宽即可.
(gdb) set follow-fork-mode child#Makefile项目管理：选择跟踪父子进程（fork()）
      ctrl+c：退出输入
(gdb) signal <singal>: 产生信号量
(gdb) jump <line>/<addr> : 修改程序的执行顺序, **不推荐**, 因为jump不会修改当前栈内容, 函数运行返回进行弹栈时必定出错, 且错误未可知.
(gdb) call <expr>: 强制调用函数
(gdb) return : 强制函数返回
(gdb) shell : 使用户不离开gdb就能使用shell
(gdb) make : 使用户不离开gdb就能使用make
 ```

## tools
- [gdb-dashboard](https://github.com/cyrus-and/gdb-dashboard) 

## FAQ
### gdb命令连写
```bash
$ gdb -ex 'target remote :1234' \
-ex 'break *0x7c00' \
-ex 'continue' \
-ex 'x/3i $pc'
$ gdb -ix gdb_init_real_mode.txt \ # 使用gdb script + 连写
-ex 'target remote localhost:8000' \
-ex 'break *0x7c00' \
-ex 'continue'
$ vim gdbinit
target remote 127.0.0.1:1234
file bin/kernel
$ gdb -x gdbinit # 仅使用gdb script
```
### gdb传递空环境变量
`env - gdb /home/chen/hello` : 此时 gdb 环境下所包含的环境变量仅为其新增加的 LINES 和 COLUMNS.

### 开启core dump
ref:
- [Core文件去哪里了](https://www.jianshu.com/p/7317910210a4)

生成core文件：先用`$ ulimit -c ${0|1024|unlimited}`(0表示关闭)开启core，当程序出错会自动生成core文件, 调试core时用`gdb a.out core`

core dump位置:
- ubuntu 20.04:/var/lib/apport/coredump