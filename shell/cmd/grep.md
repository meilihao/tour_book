# grep

## 描述

文本检索

## 语法格式

```
grep [OPTIONS] 关键词 [FILE...]
```

## 选项

- -c : 统计文件中包含文本的次数
- -e : 使用模式匹配
- -E : 使用extended regular expression模式进行匹配
- -i : 忽略大小写
- -l : 只打印文件名
- -n : 打印匹配的行号
- -o : 只输出匹配的文本行
- -r : 递归的检索目录下的所有文件(包括子目录)
- -v : 只输出没有匹配的文本行

## 例
```
$ grep -c “text” filename
$ grep -r "Rows" # 检索当前目录下包含字符串"Rows"的文件
$ grep -e "class" -e "vitural" file # 匹配多个模式
$ cat LOG.* | tr a-z A-Z | grep "FROM " | grep "WHERE" > b # 查找日志中的所有带where条件的sql
```
