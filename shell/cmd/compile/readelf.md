# readelf
一个非常有用的解析 ELF 二进制文件的工具. 我们可使用 readelf 命令收集
符号、段、节、重定向入口、数据动态链接等相关信息.

## 选项
- -a : 显示全部信息, 等价于 -h -l -S -s -r -d -V -A -I
- -h : 显示elf文件开始的文件头信息
- -l : 显示程序头（及其段头）信息(如果有的话)
- -S : 显示节头信息(如果有的话)
- -g : 显示节组信息(如果有的话)
- -t : 显示节的详细信息, 与`-S`联用
- -s : 显示符号表段中的项（如果有的话）
- -e : 显示全部头信息，等价于: -h -l -S
- -n : 显示note段（内核注释）的信息
- -r : 显示可重定位段的信息
- -u : 显示unwind段信息. 当前只支持IA64 ELF的unwind段信息
- -d : 显示动态段的信息
- -V : 显示版本段的信息
- -A : 显示CPU构架信息
- -D : 使用动态段中的符号表显示符号，而不是使用符号段
- -x <number or name> : 以16进制方式显示指定段内内容. number指定段表中段的索引,或字符串指定文件中的段名
- -w[liaprmfFsoR] 显示调试段中指定的内容
- -I 显示符号的时候，显示bucket list长度的柱状图
- -v : 显示readelf的版本信息
- -W : 宽行输出

## 示例
```
readelf -h <object> # 文件开始部分的ELF header信息
readelf -S <object> # 查询节头表
readelf -l <object> # 查询程序头表
readelf -s <object> # 查询符号表
readelf -e <object> # 显示文件中的所有header, `-e` = `-h -l -S`
readelf -r <object> # 查询重定位入口
readelf -d <object> # 查询动态段
readelf a.out -x .rodata # 以16进制输出dump的内容
```