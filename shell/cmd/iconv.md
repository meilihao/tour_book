# iconv

## 描述

转换文件的编码格式

## 格式

    iconv -f encoding [-t encoding] [inputfile]... 

## 选项
- -f encoding : 把字符从encoding编码转出
- -t encoding : 把字符转换到encoding编码
- -l : 列出系统支持的编码字符集合
- -o file : 指定输出文件
- -c : 忽略输出的非法字符
- --verbose : 显示进度信息