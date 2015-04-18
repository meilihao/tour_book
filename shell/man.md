# man

## 描述

显示man手册

## 参数

- -k: 相当于apropos命令,显示所有包含匹配项的man页面的简短描述

## man手册各章节

1. 用户在shell环境中可以操作的命令或可执行文件
2. 系统内核可调用的函数与工具等
3. 一些常用的函数（function）与函数库（library），大部分为C的函数库（libc）
4. 设备文件的说明，通常是在/dev下的文件
5. 配置文件或者是某些文件的格式
6. 游戏（games）
7. 惯例与协议等，例如Linux文件系统、网络协议、ASCII code等说明
8. 系统管理员可用的管理命令
9. 跟kernel有关的文件

>1，5，7这三个号码常用

>重构man数据库的方法:RedHat:makewhatis命令;Ubuntu,SUSE:mandb命令.