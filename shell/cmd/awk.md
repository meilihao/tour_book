# awk

## 描述

[awk](http://man.linuxde.net/awk)是一种编程语言，用于在linux/unix下对文本和数据进行处理.

## 选项

- -f 指定包含awk文本处理程序的脚本文件
- -F 设置定界符（默认为空格）

## awk脚本基本结构
```
awk 'BEGIN{ print "start" } pattern{ commands } END{ print "end" }' file
```

## awk的工作原理
```
awk 'BEGIN{ commands } pattern{ commands } END{ commands }'
```
执行步骤:
1. 执行BEGIN{ commands }语句块中的语句
2. 从文件或标准输入(stdin)读取一行，然后执行pattern{ commands }语句块，它逐行扫描文件，从第一行到最后一行重复这个过程，直到文件全部被读取完毕
3. 当读至输入流末尾时，执行END{ commands }语句块

- `BEGIN语句块`在awk开始从输入流中读取行之前被执行，这是一个可选的语句块，比如变量初始化、打印输出表格的表头等语句通常可以写在BEGIN语句块中
- `END语句块`在awk从输入流中读取完所有的行之后即被执行，比如打印所有行的分析结果这类信息汇总都是在END语句块中完成，它也是一个可选语句块
- `pattern语句块`中的通用命令是最重要的部分，它也是可选的。如果没有提供pattern语句块，则默认执行``{ print }`，即打印每一个读取到的行，awk读取的每一行都会执行该语句块.

## awk的特殊变量

- NR : 表示记录数量，在执行过程中对应当前行号
- NF : 表示字段数量，在执行过程总对应当前行的字段数
- $0 : 这个变量包含执行过程中当前行的文本内容
- $1 : 第一个字段的文本内容
- $2 : 第二个字段的文本内容

## 例
```
$ ls -lrt | awk '{print $1}' # 打印第一列
$ awk ' END {print NR}' file # 统计文件的行数
$ awk '{print $(NF-1)}' file # 打印倒数第二列
$ awk '{sum += $1} END {print sum}' file # 汇总第一列
```
