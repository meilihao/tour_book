## shell执行位置

    $ bash hello.sh //等于"./hello.sh",在子shell中执行
    $ source hello.sh //在当前shell执行,用于设置当前环境变量或针对当前shell做定制时使用
 
> 可在hello.sh中加echo "$SHLVL" 来说明

    echo "$SHLVL" #找出子shell的层级或临时shell的嵌套层级,每个bash实例启动后，变量$SHLVL的值都会加一

## 当前shell及版本

    $ echo $SHELL #查看当前是什么shell
    /bin/bash
    $ bash --version #再用"shell_name --version"查看具体版本,即"$SHELL --version"
 
参考:http://linux.cn/article-5151-1.html

## 判断内建命令

    type -a command

## 显示内建命令的help信息

    help command