# make
最常用的构建工具.

make和makefile是项目编译的管理方式.

## 选项
- -B : make 命令不会编译那些自从上次编译之后就没有更改的文件, 但此参数会忽略该设定, 全部重新编译
- -C : 将当前工作目录转移到指定的位置, 再还行该目录下的Makefile, 最后返回原目录
- -d : 打印详细信息
- -e : 环境变量覆盖 makefile 中的变量
- -f : 指定Makefile文件名称
- -k : 执行命令出错时, 放弃当前目标, 继续其他目标
- -i : 忽略错误
- -I : 在 <目录> 中搜索被包含的 makefile
- M= : 当用户需要以某个内核为基础编译一个外部模块的话，需要在make modules 命令中加入`M=dir`, 程序会自动到指定的dir目录中查找模块源码，将其编译，生成ko文件
- -n : 只打印命令过程，**不实际执行**
- -p : 打印 make 的内部数据库, 即显示Makefile中所有的变量和隐式规则
- -r : 禁用内置隐含规则
- -s : 执行但不输出命令过程, 常用于检查Makefile的正确性
- -S : 关闭`-k`, 即执行命令出错就退出
- -t : touch 目标（更新修改时间）
- -V : 显示make的version

## makefile
make主要功能就是通过makeflie来实现的. 它定义了各种源文件间的依赖关系, 阐明了源文件如何进行编译.

linux下, 通常用Makefile代替makefile, 通过`configure`来生成. 在命令行执行make时, make默认会在当前目录查找Makefile或makefile. 如果使用其他文件作为Makefile则需要用`-f <makefile>`参数明确指明, make默认会执行Makefile中的第一个target, 且make或递归查找依赖, 如果被依赖的文件不存在, make就会退出.

对于Makefile, 最主要的就是目标, 条件和命令三大要素.

> make会对比target和它所有依赖的时间戳, 当发现target比它的任一依赖的时间戳小时会重新构建target, 该过程会递归进行, 因为该target的依赖可能是其他target.

### 规则
**shell 命令必须放到target里或使用`$(shell ...)`的形式**

Makefile由一组规则（Rule）组成，每条规则的格式是:
```makefile
targets ... : prerequisites ...
    command
    ...
    ...
```

- targets : 目标文件名, 多个文件以空格分隔, 可使用通配符(`*`,`?`,`[...]`), 至少要有一个target.
- prerequisites : 目标所依赖的文件或target, 任意个数包括0. 0个时只有工作路径下target所代表的文件不存在时才执行command, 因为此时没法对比时间戳. 
- command : 如果它不与`target : prerequisites`写在一行则**必须以tab键开头**, 否则和prerequisites在同一行, 用`;`分隔. 如果命令过长, 可用`\`作为换行连接.

target也就是一个目标文件，可以是object file，也可以是执行文件, 还可以是一个标签（label）. prerequisites就是，要生成那个target所需要的文件或是目标. command也就是make需要执行的命令（任意的shell命令）. 这是一个文件的依赖关系，也就是说，target这一个或多个的目标文件依赖于prerequisites中的文件，其生成规则定义在 command中. 如果prerequisites中有一个以上的文件比target文件要新(修改日期)或target不存在，那么command所定义的命令就会被执行. 这就是makefile的规则, 也就是makefile中最核心的内容.

常用target名称:
- all : 表示编译所有内容, 是不给定target时make执行的默认target
- clean : 删除编译生成的二进制文件
- distclean : 不仅删除编译生成的二进制文件，也删除其它生成的文件，例如配置文件和格式转换后的文档，执行make distclean之后应该清除所有这些文件，只留下源文件
- install : 需安装的内容,即执行编译后的安装工作，把可执行文件、配置文件、文档等分别拷到不同的安装目录

make用`.PHONY`显式指明伪目标. 伪目标不是文件, 仅是一个标签, 因此Makefile不会生成伪目标对应的文件. 伪目标的特性是总能被执行.
make使用"include"来包含其他文件, 用空格分隔.

Makefile组成:
- 显式规则

	它说明了如何生成一个或多个目标, 由Makefile的书写者**显示指出**要生成的文件, 文件的依赖及生成的命令.
- 隐式规则

	因为make有自动推导功能, 会选择一套默认的方法进行make, 它让开发者可以比较简略地书写Makefile. **不推荐使用**.

	如果一个目标在Makefile中的所有规则都没有命令列表，make会尝试在内建的隐含规则（Implicit Rule）数据库中查找适用的规则. make的隐含规则数据库可以用make -p命令打印.

	模式规则: `%`定义对文件名的匹配, 表示任意长度的非空字符串. `%.o : %.c ; <cmd ...>`表明了从所有`.c`文件生成相应的`.o`文件的规则
- 变量定义

	在Makefile中定义一系列需要的变量, 类似C语言中的宏, 当Makefile执行时, 其中的变量会被扩展到相应的引用位置上.

	使用时使用`$()`或`${}`来引入变量.

	推荐使用`:=`定义变量, 保证引用的变量在前面已定义过.

	`?=`表示未定义时赋值, 否则跳过.

	`+=`表示给变量追加值, 会延用定义该变量时的`:=`或`=`

	目标变量即Makefile的局部变量, 仅在target中有效, 离开作用域后会恢复旧值, 语法为：
	```
	<target...>:<variable-assignment>
	<target...>:override <variable-assignment>
	```
	<variable-assignment>可以是前面讲过的各种赋值表达式，如`=`、`:=`、`+=`或是`?=`

	常用的自动变量(随上下文的不同而改变)有：
    - $@，表示规则(Rule)中的目标
    - $<，表示规则中的第一个条件
    - $?，表示规则中所有时间戳比目标新的条件，组成一个列表，以空格分隔
    - $^，表示规则中的所有条件，组成一个列表(已去重)，以空格分隔
	- $+, 与$^类似, 但没有排除重复条件
	- $*, 目标的主文件名, 不包括扩展名

	```makefile
	main: main.o stack.o maze.o
	gcc $^ -o $@
	# --- 等价于
	main: main.o stack.o maze.o
	gcc main.o stack.o maze.o -o main
	```
- 文件指示

	1. 一个Makefile引用另一个Makefile
	1. 根据某些情况指定Makefile中的有效部分
	1. 定义一个多行的命令
- 注释

	仅有行注释, 用`#`开头


每当target中的command执行完毕后, make会检测它们的返回码, 如果成功则继续执行, 否则中止该target. 在命令前加`-`则会忽略检查返回码, 认为都成功. 命令前加`@`表示不将要执行的命令输出到屏幕.命令前加`+`表示只显示命令但不执行, 常用于递归式的makefile. 如果要让上一条命令的结果应用在下一条上, 则它们需要写在一行并用`;`分隔.

### 宏
```
define MACRO_NAME
	...
endef
```

带有参数的宏就是函数, 语法是`$(call MACRO_NAME [, arg1, ... , argn])`, 宏中引用参数用`$1, ... $n`表示

### 函数调用
`$(<function> <args>)`或`${<function> <args>}`

### 条件指令
ifdef, ifndef, ifeq, ifneq, else, 结尾统一用endif.

## FAQ
### 了解make时执行了哪些命令
```bash
$ make [test] "V=1"
```

### Makefile定义的操作
`make help`

### make if判断明明正确却没有日志输出
```makefile
pkg:
    @echo "----"
	@if [ -d "$(TOPDIR)/www_replace" ]; then \ # 条件判断明明成功却没输出make执行的命令
		echo "replace www"; \
		rm -rf  $(PKG_DIR)$(PRODUCT_ROOT)/www; \
		cp -r $(TOPDIR)/www_replace $(PKG_DIR)$(PRODUCT_ROOT)/www; \
	fi
```

makefile执行时, `if`即使为true, 里面的命令执行日志也不会输出, 因此建议在if中手动添加`echo`.

### make xxx Is a directory. Stop
Makefile要求每行结尾，一定要确认没有空格，直接是换行.

原因:
```makefile
TOPDIR = $(realpath .) # in docker : `/app/xxx`
```

解决:
```makefile
# in docker : `TOPDIR=/app/xxx
TOPDIR = $(realpath .)
```

### *** 遗漏分隔符
makefile中的command必须以tab键开头.

### Makefile调用栈
类似函数栈

### 如何重新执行某项 make check
根据make check日志, 重新进入相应的目录, 直接执行`make check`即可.

### missing separator.  Stop.
在makefile中，命令行要以tab键开头. vscode因为配置原因将tab转成了4个空格.

### recipe commences before first target. stop
在target前使用了`ifeq ... else ifeq...`, 将这些代码移入一个target即可.

```makefile
SCRIPTS = abc

# 不使用`uname -p`的原因: 它可能返回`unknown`
ARCH = $(shell arch)

deps-arch:
ifeq ($(ARCH), x86_64)
	cp -f app-amd64 app
	$(eval SCRIPTS += app)
else ifeq ($(ARCH), aarch64)
	cp -f app-arm64 app
	$(eval SCRIPTS += app)
else
endif
```

扩展:
```makefile
SCRIPTS = abc

ARCH = $(shell arch)

ifeq ($(ARCH), x86_64)
                tmp = $(shell cp -f a.key a.key1) # cp命令不执行, 原因是没有使用临时变量tmp导致该语句被忽略(使用`tmp := ...`声明的变量即使不使用也不会被忽略). 将下面的`111`替换为`$(tmp)`, cp命令就执行了.
                SCRIPTS += 111
else ifeq ($(ARCH), aarch64)
                SCRIPTS += "456"
else
                SCRIPTS += 000
endif

all:
                echo "-----"
                echo "$(SCRIPTS)"
```

### 调用shell
```makfile
buildtime = $(shell date +"%Y%m%d%H%M")
```

### 调用python
```makfile
setversion:
    python3 xxx.py
```

### make 传递参数
```makefile
version = $(version)
```

执行`make version=1.0.0`即可

### if...elseif...else...
```makefile
TARGET_CPU = unknown
ifeq ($(TARGET_CPU),x86)
  TARGET_CPU_IS_X86 := 1
else ifeq ($(TARGET_CPU),x86_64)
  TARGET_CPU_IS_X86 := 1
else
  TARGET_CPU_IS_X86 := 0
endif
```

### [使用环境变量](https://stackoverflow.com/questions/4728810/how-to-ensure-makefile-variable-is-set-as-a-prerequisite)
```makefile
check-env:
    if test "$(ENV)" = "" ; then \
        echo "ENV not set"; \ # 分号不能丢
        exit 1; \
    fi
```

## cmake
参考:
- [抛弃 Autotools 向 CMake 迈进吧](https://linux.cn/article-13419-1.html)