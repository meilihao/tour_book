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
其通过`;`来分隔.

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
- `>` : 追加文件到目标文件

## Bash 中的特殊字符大全

[Bash 中的特殊字符大全](https://linux.cn/article-5657-1.html)

## shell中的叹号
shell中`！`叫做[事件提示符](https://linux.cn/article-5658-1.html)，英文是：Event Designators,可以方便的引用历史命令，也就是history中记录的命令.可用`set +H`来关闭,`set -H`来打开.

- `!`: 当后面跟随的字符串不是“空格、换行、回车、制表符、=”时，做历史命令的替换,否则只是单纯 的叹号字符.
- `!$`:是上一条命令中最后一个parameter
- `!n`: 会引用history中的第n个命令，比如输入`！100`，就是执行history列表中的第100条命令
- `!-n`: 获取history中倒数第N个命令并执行，比如输入`!-1`,就会执行上一条命令,`!!`是`!-1`的一个alias
- `!string`:引用最近的以 string 开始的命令
- `!?string[?]`:引用最近包含这个字符串的命令
