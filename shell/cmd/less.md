# less

## 描述

文本浏览器, 分页显示文件内容

## 格式

    less [OPTION]... [FILE]...

## 选项
- -i : 搜索时忽略大小写
- -m : 显示进度百分比
- -N : 显示行号
- -r : 显示原始控制字符, 比如`ANSI escape code`即terminal color
- -s : 将连续的空行压缩为一行

## 子命令
- j : 向上移动
- k : 向下移动
- q : 退出less
- / : 搜索
- & : 只显示文件中包含某些内容的行
- n : 下一个, 常与`/`结合使用
- b : 上一个, 常与`/`结合使用

## 扩展
### more命令
less是more的增强版, 不推荐使用more