# cut

## 描述

显示行中的指定部分或删除文件中指定部分

## 选项

- -b : 仅显示行中指定字节范围的内容
- -c : 仅显示行中指定字符范围的内容
- -d : 指定定界符
- -f : 仅显示行中指定字段范围的内容
- -n : 取消分割多字节字符, 与`-b`联用
- --complement：补足被选择的字节、字符或字段

cut选取的范围:
- N : N
- N- : N~结尾
- -M : 1~M
- N-M : N到M

缺点: `-d`只擅长处理`以一个字符间隔`的文本内容, 复杂情况可用awk代替.

## 例
```sh
$cut -f2,4 filename # 截取文件的第2列和第4列
$cut -f3 --complement filename # 取文件除第3列外的所有列
$cat -f2 -d";" filename
$cut -c1-5 file # 打印第一到5个字符
$cut -c-2 file  # 打印前2个字符
$cut -b 3-5,10 file
$cut -b -3,3- file # 显示整行, 不会重复输出第3个字符
```

# column
## example
```bash
column -t -s, file.csv
```