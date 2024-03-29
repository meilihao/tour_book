# bash
使用`#`表示注释.

## 通配符
bash支持通配符有:
- `*` : 一个或多个字符
- `?` : 一个字符
- `[]`: 中括号内的任意字符, 支持`-`

## 引号
shell支持单引号, 双引号和反引号.

```bash
$ echo '$PATH'
$PATH
```

### 单引号
单引号会忽略所有的特殊字符, 因此用单引号包裹的字符仅作为普通字符来解析.

### 双引号
双引号会忽略大部分的特殊字符, 但`$`, `\`, ```除外.

### 反引号
反引号包裹的内容会被当作shell命令来解析.

## shell首行的"#!"
1. `#!/bin/env bash`,因为解释器在linux中可能被安装到不同的目录，env可以在系统的PATH环境变量中所定义的目录里查找解析器路径,**推荐**.
2. `#!/bin/bash`,会固定解析器的路径.

## shell执行位置
1. `$ ./hello.sh` 在子shell中执行(即脚本设置的变量不会影响当前shell),需`执行权限`
2. `$ bash hello.sh` 在子shell中执行,可无`执行权限`,将脚本作为sh的命令行参数来运行会忽略首行的"#! ..."
3. `$ source hello.sh`
在当前shell执行,用于设置当前环境变量或针对当前shell做定制时使用,可无`执行权限`,该命令通常用命令`.`来替代.

> 可在hello.sh中加echo "$SHLVL" 来说明

    echo "$SHLVL" #找出子shell的层级或临时shell的嵌套层级,每个bash实例启动后，变量$SHLVL的值都会加一

## 命令序列(即多个命令在一行)
通过`;`来分隔.

```sh
$ cmd1 ; cmd2
# 这等同于：
$ cmd1
$ cmd2
```

## 当前shell及版本

    $ echo $SHELL #查看当前是什么shell # $是变量符号
    /bin/bash
    $ echo $0
    bash
    $ bash --version #再用"shell_name --version"查看具体版本,即"$SHELL --version"

参考:http://linux.cn/article-5151-1.html

## 判断内建命令

    type -a command

## 显示内建命令的help信息

    help command

## 标准输入、输出和错误

	0 – stdin (standard input)，1 – stdout (standard output)，2 – stderr (standard error)

## 重定向
输入重定向中用到的符号及其作用:
- `命令 < 文件` : 将文件作为命令的标准输入
- `命令 << 分界符` : 从标准输入中读入，直到遇见分界符才停止
- `命令 < 文件 1 > 文件 2` : 将文件 1 作为命令的标准输入并将标准输出到文件 2

输出重定向中用到的符号及其作用:
- `cmd > file` : 将标准输出重定向到一个文件中（清空原有文件的数据）, 等价于`cmd 1> file`
- `cmd 2> file` : 将错误输出重定向到一个文件中（清空原有文件的数据）
- `cmd >> file` : 将标准输出重定向到一个文件中（追加到原有内容的后面）
- `cmd 2>> file` : 将错误输出重定向到一个文件中（追加到原有内容的后面）
- `cmd 1>file1 2>file2`
- `cmd >file 2>&1` : 将标准输出与错误输出共同写入到文件中（清空原有文件的数据）或`cmd 2>&1 file`<=>`cmd &> file`
- `cmd >> file 2>&1`: 将标准输出与错误输出共同写入到文件中（追加到原有内容的后面）= `cmd &>> file`
- `cmd >&m` : 把stdout重定向到文件描述符m
- `cmd <&m` : 把文件描述符m作为输入

注意: `python 1.py' > /var/log/1.err 2>&1`正常输出到文件; `python &> /var/log/1.err `不正常, 它输出到terminal

重定向快速记忆:
- `<` : 从标准输入中读入
- `<<` : 从标准输入中读入, 直到遇到分隔符
- `>` : 重定向文本,但会先清空目标文件
- `>>` : 追加文件到目标文件

```shell
cmd > &m // 将标准输出定向到文件描述符m中.
cmd <&- // 关闭标准输入

# mail -s "Readme" root@example.com << over
> hello chen
> over
```

### Here Document
Here Document 是在Linux Shell 中的一种特殊的重定向方式，它的基本的形式如下:
```
cmd << delimiter // 此处如果使用了`-delimiter`则结尾的delimiter可不用顶格, 但要用tab缩进
  Here Document Content
delimiter
```
它的作用就是将两个 delimiter 之间的内容(Here Document Content 部分) 作为标准输入的内容传递给cmd.

注意:
- EOF 只是一个标识而已，可以替换成任意的合法字符
- **作为结尾的delimiter一定要顶格写，前面不能有任何字符**
- **作为结尾的delimiter后面也不能有任何的字符（包括空格）**

    否则报: "warning: here-document at line N delimited by end-of-file (wanted `EOF')"
- 作为起始的delimiter前后的空格会被省略掉

```bash
# cat << __EOF__ | echo "abc"
__EOF__
abc
```

```bash
cat << EOF > abcd.txt
hello world
EOF
```

> cat+delimiter时内容中的tab会被解析成bash中的按tab键操作, 需要用空格替代tab.

delimiter包含变量:
```bash
$ LFSVersion="10.0-systemd"
$ cat << EOF
$LFSVersion
DISTRIB_ID="Linux From Scratch"
DISTRIB_RELEASE="${LFSVersion}"
DISTRIB_CODENAME="BigBang"
DISTRIB_DESCRIPTION="Linux From Scratch for fun"
EOF
10.0-systemd
DISTRIB_ID="Linux From Scratch"
DISTRIB_RELEASE="10.0-systemd"
DISTRIB_CODENAME="BigBang"
DISTRIB_DESCRIPTION="Linux From Scratch for fun"
$ cat << "EOF" # delimiter使用""包裹表示直接按文本格式输出 
$LFSVersion
DISTRIB_ID="Linux From Scratch"
DISTRIB_RELEASE="${LFSVersion}"
DISTRIB_CODENAME="BigBang"
DISTRIB_DESCRIPTION="Linux From Scratch for fun"
EOF
$LFSVersion
DISTRIB_ID="Linux From Scratch"
DISTRIB_RELEASE="${LFSVersion}"
DISTRIB_CODENAME="BigBang"
DISTRIB_DESCRIPTION="Linux From Scratch for fun"
```

## 管道
```sh
# echo "linuxprobe" | passwd --stdin root
```

## Bash 中的特殊字符大全

[Bash 中的特殊字符大全](https://linux.cn/article-5657-1.html)

- EOF : "end of file"，表示文本结束符
- << : Here-document，用来将后继的内容重定向到左侧命令的stdin中.<<紧跟一个标识符，从下一行开始是想要引用的文字，然后再在单独的一行用相同的标识符关闭
- 冒号 : 空命令，这个命令什么都不做，但是有返回值，返回值为0（即：true）

## shell中的叹号
shell中`！`叫做[事件提示符](https://linux.cn/article-5658-1.html)，英文是：Event Designators,可以方便的引用历史命令，也就是history中记录的命令.可用`set +H`来关闭,`set -H`来打开.

- `!`: 当后面跟随的字符串不是“空格、换行、回车、制表符、=”时，做历史命令的替换,否则只是单纯 的叹号字符.
- `!$`:是上一条命令中最后一个parameter
- `!n`: 会引用history中的第n个命令，比如输入`！100`，就是执行history列表中的第100条命令
- `!-n`: 获取history中倒数第N个命令并执行，比如输入`!-1`,就会执行上一条命令,`!!`是`!-1`的一个alias
- `!string`:引用最近的以 string 开始的命令
- `!?string[?]`:引用最近包含这个字符串的命令

## 读取命令序列的输出
> 在bash shell中, $()与``(反引号)都是用来做命令替换(command substitution)的.

1. 子shell(subshell),**推荐**
        cmd_output=$(COMMANDS)

        ```bash
        $  (find . | grep -v ".so" |xargs ls -al --time-style=+"" ) > my.log
        ```

2. 反引用(其实该方法也是用子shell来运行命令), 唯一优点: 跟其他unix shell的兼容性高.
        cmd_output=`COMMANDS`

>使用$()显然比``优越，这是因为：
1. 前者更易读，不会产生歧义;而反引号`常常被初学者当成单引号'；
1. 前者嵌套时更简单，直接使用就行;而后者嵌套时内部的反引号必须用\转义；
1. 它们对反斜杠\的处理不一样，在$()中可以减少转义的麻烦。而这一点与第二点是前因后果的关系。正是因为$()嵌套时不需转义，所以\在$()中就不需要作为一个特殊字符了。而``中的\必须是特殊字符，否则就无法嵌套使用了.

```shell
out=$(cat text.txt)
echo $out  #丢失所有换行符

out="$(cat text.txt)"
echo $out  #保留换行符
```

一般情况下，$var与${var}并没有啥不一样,  但是用${}会比较精准的界定变量名称的范围， 比方说:
```sh
$ A=B
$ echo $AB

$ A=B
$ echo ${A}B
$ BB
```

当然`${}`还有其他更多功能.

## 子shell

子shell本身是独立进程, 不会对当前shell有任何影响,可用`()`来定义一个子shell.

xargs只能以有限的几种方式来提供参数,而且它也不能为多组命令提供参数.要执行一些包含来自标准输入的多个参数的命令时可用子shell来处理.
```shell
cat file | ( while read arg; do cat $arg; done ) <==> cat file | xargs -I {} cat {}
```
在上述while循环中,可将`cat $arg`替换成任意数量的命令,这样我们就可以对同一个参数执行多项命令,同时也可以不借助管道,将参数传递给其他命令.

## 特殊符号
```sh
$ : > a.txt # 清空文件
```

## bash快捷键
```
Ctl-U   删除光标到行首的所有字符,在某些设置下,删除全行
Ctl-W   删除当前光标到前边的最近一个空格之间的字符
Ctl-H   backspace,删除光标前边的字符
```

## 追踪
使用`-x`选项, 比如`#!/bin/sh -x`或`sh -x xxx.sh`,可追踪shell每个命令的执行结果.

## assert函数

## {}的特殊用法
```bash
$ echo a{a,b}
aa ab
$ echo a{,b}
a ab
$ echo a{a..c}
aa ab ac
```

## 屏蔽别名
1. 使用绝对路径
1. 使用`\`屏蔽系统别名, 比如`\mv`

查找别名:
```bash
$ alias ll
alias ll='ls -l'
```

## bash
- 内置命令`fc`能捕获上次命令并送入编辑器.
- 使用`type 命令名称`可判断输入的命令是内部命令还是外部命令
- `help`命令可查看bash的help doc.

## =
`var = value`不同于`var=value`, 两边没有空格的等号是赋值操作符，加上空格的等号表示的
是等量关系测试.

## 安全

- 执行当前可执行程序加`./`的原因：

	主要是安全原因，因为在linux中执行程序时，会先搜索当前目录然后是系统目录，所以如果当前目录中有与系统可执行程序重名的程序，比如cp，她就会优先执行当前目录中的cp，但是如果当前目录的cp是木马，就会威胁到系统安全，所以这是Linux的一种安全策略，所以默认并没有把当前目录加到环境变量PATH中去.

## `. FILE`
`. ./buildenv.mk`(`.`+空格+路径)作用: 点是内置命令`source`的快捷方式，该命令是在当前bash进程中运行脚本, 有点类似编程语言中的import, 可以直接理解为文件包含.

## logger/日志
logger是一个shell命令接口，可以通过该接口使用Syslog的系统日志模块，还可以从命令行直接向系统日志文件写入一行信息

## set命令
```conf
set -e # 若指令传回值不等于0, 则立即退出shell
```

## 调试
参考:
- [Shell脚本调试的几种方式](https://blog.csdn.net/Jerry_1126/article/details/52096886)

### 脚本调试

设置`export PS4='+${BASH_SOURCE}:${LINENO}:${FUNCNAME[0]}: '`, 允许调试shell脚本时，在跟踪里输出文件/行号, 此时可能导致某些脚步报错.

使用选项`-x`,启动追踪调试shell脚本.

    # bash -x script.sh

> 可使用_DEBUG环境变量来自定义显示调试信息.

### 调试选项

1. 在脚本中设置开关
```shell
set -x # 在执行时显示参数和命令, 必须设置在最前面, 否则不能生效. 之前有写在`echo "xxx"`下失效的情况.
set +x # 关闭调试
set -v # 当命令进行读取时显示输入
set +v # 禁止打印输入
```

2. 直接修改脚本
```shell
#!/bin/bash -xv
```

### 具体调试

1 . `-c` :  使Shell解析器从字符串而非文件中读取并执行命令
场合: 当需要调试一小段脚本的执行结果时，非常方便
示例: `# bash -c 'x=1;y=2;let z=x+y;echo "z=$z"'`

1. `-x` : 提供跟踪执行信息,将执行脚本的过程中把实际执行的每个命令显示出来，行首显示+, +后面显示经过替换之后的命令行内容，有助于分析实际执行的是什么命令.
场合: 是调试Shell脚本的强有力工具，**是Shell脚本首选的调试手段**
示例: 
    1. 在命令行提供参数：`$ sh -x script.sh`
    2. 脚本开头提供参数：`#!/bin/sh -x`
    3. 在脚本中用set命令启用or禁用参数：`其中set -x表启用，set +x表禁用`

1. `-v` : 区别于-x参数,该选项打印命令行的原始内容，-x参数打印出经过替换后命令行的内容
场合: 仅想显示命令行的原始内容
示例: `# bash -v script.sh`

1. 使用调试工具bashdb, 脚本中的gdb
常用参数：

    1. 列出代码和查询代码类：
    - `l` : 列出当前行以下的10行
    - `-` :  列出正在执行的代码行的前面10行
    - `.` :  回到正在执行的代码行
    - `w` :  列出正在执行的代码行前后的代码
    - `/pat/` : 向后搜索pat
    - `？pat？` : 向前搜索pat

    1. Debug控制类：
    - `h` :  帮助
    - `help  命令` : 得到命令的具体信息
    - `q` : 退出bashdb
    - `x`  :   算数表达式 计算算数表达式的值，并显示出来
    - `!!` :  空格Shell命令 参数 执行shell命令

    1. 控制脚本执行类(使用bashdb进行debug的常用命令)：
    - `n` :  执行下一条语句，遇到函数，不进入函数里面执行，将函数当作黑盒
    - `s + n` : 单步执行n次，遇到函数进入函数里面
    - `b   行号n` : 在行号n处设置断点
    - `del 行号n` :  撤销行号n处的断点
    - `c 行号n` :  一直执行到行号n处
    - `R` :   重新启动当前调试脚本
    - `Finish` :  执行到程序最后
    - `cond n expr` : 条件断点

1. caller N : 在函数中使用caller可在stdout中打印该函数的相关信息. N表示打印哪层的frame
1. trap : 指定在收到信号后采取的动作, 格式是`trap <command> <signals>`

        > `trap -l` : 打印所有signal

## FAQ
### 执行bash脚本后机器变得很卡, 此时cpu和内存都不高
该脚本来自git repo(在liunx上执行正常), 先git pull到windows上, 再上传到oracle linux 7.9, 此时该脚本的换行是windows换行.

```bash
# ./build.sh # 执行报如下错误, 且执行后os变得很卡
: No such file or director
# bash build.sh
build.sh: line 2: $'\r': command not found
# echo $SHELL
/bin/bash
# bash --version
... 4.2.46(2)...
```

### 程序调用命令行执行报`shell-init: error retrieving current directory: getcwd: cannot access parent directories: No such file or directory`
getcwd 命令无法定位到当前工作目录, 即工作目录被删除了, 或者曾经被删除过.

### 切换工作目录到当前脚本所在目录
```bash
#!/bin/bash
cd "$(dirname "$0")"
```