# man

## 描述

显示man手册

## 类似
- `xxx --help` : 具体命令的help内容
- `help xxx` : 获取shell的内置命令信息

## 参数

- -k: 相当于apropos命令,根据命令中部分关键字来查询命令,会显示所有包含匹配项的man页面的简短描述.

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

man 命令中常用按键以及用途:
空格键 向下翻一页
PaGe down 向下翻一页
PaGe up 向上翻一页
home 直接前往首页
end 直接前往尾页
/ 从上至下搜索某个关键词，如“/linux”
? 从下至上搜索某个关键词，如“?linux”
n 定位到下一个搜索到的关键词
N 定位到上一个搜索到的关键词
q 退出帮助文档

man 命令帮助信息的结构以及意义:
结构名称 代表意义
NAME 命令的名称
SYNOPSIS 参数的大致使用方法
DESCRIPTION 介绍说明
EXAMPLES 演示（附带简单说明）
OVERVIEW 概述
DEFAULTS 默认的功能
OPTIONS 具体的可用选项（带介绍）
ENVIRONMENT 环境变量
FILES 用到的文件
SEE ALSO 相关的资料
HISTORY 维护历史与联系方式

## 例
```sh
$man 3 printf
$man man # 查看 man手册的 章节 号及其包含的手册页类型
```
