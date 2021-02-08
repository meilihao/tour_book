# addr2line
转换地址到文件名和行号.

通过分析调试信息中的Line Number Table(程序源程序的行号和编译后的机器代码之间的对应关系)就能把源码中的出错位置找出来.

DWARF格式的Line  Number Table是一种高度压缩的数据，存储的是表格前后两行的差值，在解析调试信息时，需要按照规则在内存里重建Line Number  Table才能使用.

Line Number Table存储在可执行程序的.debug_line域，使用命令`$ readelf -w ${program}`查看,比如
```sh
$ readelf -w test
...
Special opcode 146: advance Address by 10 to 0x4004fe and Line by 1 to 5
Special opcode 160: advance Address by 11 to 0x400509 and Line by 1 to 6
...
```

说明机器二进制编码的0x4004fe位置开始，对应于源码中的第5行，0x400509开始就对应与源码的第6行了，所以400506这个地址对应的是源码第5行位置.

## example
```bash
objdump -d a.out | grep -A 2 -E 'main>:|function1>:|function2>:'
000000000040051d :
40051d: 55 push %rbp
40051e: 48 89 e5 mov %rsp,%rbp
addr2line 0x40051d -e self-monitor # 0x40051d 为出错时的`pc addr`.
```