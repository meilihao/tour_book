## shell首行的"#!"
1. `#! /bin/env bash`,因为解释器在linux中可能被安装到不同的目录，env可以在系统的PATH环境变量中所定义的目录里查找解析器路径,**推荐**.
2. `#! /bin/bash`,会固定解析器的路径.

## shell执行位置
1. `$ ./hello.sh` 在子shell中执行(即脚本设置的变量不会影响当前shell),需`执行权限`
2. `$ bash hello.sh` 在子shell中执行,可无`执行权限`,将脚本作为sh的命令行参数来运行会忽略首行的"#! ..."
3. `$ source hello.sh`
在当前shell执行,用于设置当前环境变量或针对当前shell做定制时使用,可无`执行权限`,该命令通常用命令`.`来替代.

> 可在hello.sh中加echo "$SHLVL" 来说明

    echo "$SHLVL" #找出子shell的层级或临时shell的嵌套层级,每个bash实例启动后，变量$SHLVL的值都会加一

## 命令序列(即多个命令在一行)
通过`;`来分隔.

## 当前shell及版本

    $ echo $SHELL #查看当前是什么shell
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

- `>` : 重定向文本,但会先清空目标文件
- `>>` : 追加文件到目标文件

```shell
cmd > &m // 将标准输出定向到文件描述符m中.
```

## 读取
- `<` : 从标准输入中读入
- `<<` : 从标准输入中读入, 直到遇到分隔符

```
cmd <&- // 关闭标准输入
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

## 脚本调试

使用选项`-x`,启动追踪调试shell脚本.

    # bash -x script.sh

> 可使用_DEBUG环境变量来自定义显示调试信息.

### 选项

1. 在脚本中设置开关
```shell
set -x # 在执行时显示参数和命令
set +x # 关闭调试
set -v # 当命令进行读取时显示输入
set +v # 禁止打印输入
```

2. 直接修改脚本
```shell
#!/bin/bash -xv
```

## 读取命令序列的输出

1. 子shell(subshell),**推荐**
        cmd_output=$(COMMANDS)

2. 反引用(其实该方法也是用子shell来运行命令)
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

## Here Document
Here Document 是在Linux Shell 中的一种特殊的重定向方式，它的基本的形式如下:
```
cmd << delimiter
  Here Document Content
delimiter
```
它的作用就是将两个 delimiter 之间的内容(Here Document Content 部分) 传递给cmd 作为输入参数.

注意:
- EOF 只是一个标识而已，可以替换成任意的合法字符
- 作为结尾的delimiter一定要顶格写，前面不能有任何字符
- 作为结尾的delimiter后面也不能有任何的字符（包括空格）
- 作为起始的delimiter前后的空格会被省略掉

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

## bash命令中的反引号
将反引号内容作为子命令并优先执行, 比如"mv `find . -name ".txt"` dir"