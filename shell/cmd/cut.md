# cut

## 描述

显示行中的指定部分或删除文件中指定字段

## 选项

- -b : 仅显示行中指定字节范围的内容
- -c : 仅显示行中指定字符范围的内容
- -d : 指定定界符
- -f : 仅显示行中指定字段范围的内容
- --complement：补足被选择的字节、字符或字段

cut选取的范围:
- N- N~结尾
- -M 1~M
- N-M N到M

## 例
```sh
$cut -f2,4 filename # 截取文件的第2列和第4列
$cut -f3 --complement filename # 取文件除第3列外的所有列
$cat -f2 -d";" filename
$cut -c1-5 file # 打印第一到5个字符
$cut -c-2 file  # 打印前2个字符
```
