# grep/egrep
> egrep是`grep -E "$@"`的脚本链接.

## 描述
参考:
- [Linux 命令行下搜索工具大盘点，效率提高不止一倍！](https://www.tuicool.com/articles/bMnmymY)

文本检索

## 语法格式

```
grep [OPTIONS] 关键词 [FILE...]
```

## 选项

- -a, --text equivalent to --binary-files=text : 让二进制文件等价于文本
- -A : 除了列出符合条件的行外, 同时列出每个符合条件行的**后NUM行**
- -b : 将可执行文件（binary）当作文本文件（text）来搜索
- -B : 与`-A`类似, 除了列出符合条件的行外, 同时列出每个符合条件行的**前NUM行**
- -c : 统计文件中包含文本的次数
- -e : 使用模式匹配
- -E : 使用extended regular expression(ERE,扩展的正则表达式)模式进行匹配,默认是使用基本正则表达式(BRE)
- -f : 从指定文件获取匹配模式, 然后按该模式进行匹配
- -i : 忽略大小写
- -l : 只打印文件名
- -n : 打印匹配的行号
- -o : 只输出匹配的文本行
- -r : 递归的检索目录下的所有文件(包括子目录)
- -R : 类似`-r`, 但支持符号链接
- -s : 不显示错误信息
- -v : 只输出没有匹配的文本行
- --exclude-dir : 排除目录

## 例
```
$ grep -c “text” filename
$ grep -r "Rows" # 检索当前目录下包含字符串"Rows"的文件
$ grep -e "class" -e "vitural" file # 匹配多个模式
$ cat LOG.* | tr a-z A-Z | grep "FROM " | grep "WHERE" > b # 查找日志中的所有带where条件的sql
$ grep -r $'\r' * # 查找`^M`字符.($：锚定行尾，此字符前面的任意内容必须出现在行尾)
$ lsmod | grep -E "drbd|xxx" # grep 或
$ cat /var/log/syslog |grep zfs
Binary file (standard input) matches
$ grep -a zfs /var/log/syslog # 可解决`Binary file (standard input) matches`的问题
$ grep -r --include="*.lua"  "ToSearchString"  Path # 按扩展名搜索
$ cat id_rsa.pub |grep -c - authorized_keys # 判断是否指定文件. **不推荐**: 某些情况下发现匹配出错(key不存在于authorized_keys, 但输出是1)
$ cat id_rsa.pub |grep -c -f - authorized_keys # 解决上面的异常情况, **推荐**
```
