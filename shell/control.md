# if

## 描述

if 语句通过关系运算符判断表达式的真假来决定执行哪个分支,Shell 有三种 if ... else 语句：

- [if](#if)
    - [描述](#%E6%8F%8F%E8%BF%B0)
    - [例](#%E4%BE%8B)
- [case](#case)
    - [描述](#%E6%8F%8F%E8%BF%B0)
- [for](#for)
- [while](#while)
- [until](#until)

> 注意:[]中判断条件两侧的空格不可省略

## if
```shell
#!/bin/bash

function show_usage {
   echo "Usage: $0 source_dir dest_dir"
   exit $1  #会使用show_usage指定的出错码退出
}

if [ $# -ne 2 ];then # `[`实际上是一个命令，必须将其与剩余的字符串用空格隔开
   show_usage
else
   if [ -d $1 ];then
   source_dir=$1
   else
      echo "Invalid source directory"
      show_usage 1 #给函数一个确定的出错码值1
fi
if [ -d $2 ];then
   dest_dir=$2
else
   echo "Invalid destination directory"
   show_usage 2
fi
fi

printf "Source directory is ${source_dir}\n"
printf "Destinatin directory is ${dest_dir}\n"
```
```shell
# 用逻辑运算符进行简化, 短路运算更简洁
[ condition ] && action; # 如果condition为真,执行action
[ condition ] || action; # 如果condition为假,执行action
```
```shell
# -o = or , -a = and , 但推荐只用 || 或者 &&
if [[ $a > $b ]] || [[ $a < $c ]] <=>  if [ $a -gt $b -o $a -lt $c ]
```

> `[[]]`写法是保护性写法, 保证内层的`[]`失败(比如其变量没有被赋值)相当于是表达式为假.

# case

## 描述

允许通过判断来选择代码块中多条路径中的一条,类似c中的switch.

格式:
```sh
case $var in
模式1)
   命令1
   ...
   ;;
模式2)
   命令2
   ...
   ;;
*)
   其他
......
esac
```
# for

循环

格式1:
```shell
for var in 对象列表 #任何以环境变量IFS分隔的对象列表,包括变量内容,都可作为for...in的目标体
do
  $var
done
```
例:
```shell
#打印当前目录下所有匹配的文件列表
for FILE in *.sh
do
    echo $FILE
done
```

格式2:
```shell
for((i=1;i<=$var;i++));
do
   $i
done;
```
> {1..5},{a..h}均可轻松地生成不同的序列.

格式3:
```shell
# 和C语言的for循环类似,`{`后必须紧跟一个空格,换行或tab
for((i=1;i<=10;i++));{ echo $(expr $i \* 4); }
```

# while

循环

格式:
```shell
while condition
do
   ...
done
```
```shell
#!/bin/bash
exec 0<$1 #将内容输入到stdin
counter=1
while read line;do
   echo "$((counter++)):$line" # "(())"结构语句是对shell中算数及赋值运算的扩展,所有表达式可以像c语言一样
done
```
# until

一直循环直到给定的条件为真.
```shell
until condition
do
    commands;
done
```
