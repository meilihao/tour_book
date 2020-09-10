# cat

## 描述

合并/查看文件内容

## 格式

    cata [option] [file]

## 选项
- -A : 等同`-vET`
- -b : 与`-n`类似, 但会忽略空白行的行号
- -E : 在每行行尾显示结束符`$`
- -n : 显示行号
- -s : 将两行及以上的空白行压缩为一行
- -T : 用`^I`代替隐藏的`TAB`并显示出来

## 例

	# cat /dev/null > filename # 清空文件
    # cat test.txt
    # cat test1.txt test2.txt
    # OUTPUT_FROM_SOME_CMDS | cat # 从管道中读取
    # echo "test" | cat - file1 # 拼接标准输入和文件的内容,这里的`-`表示stdin的重定向的源或目的,简单理解即是stdin.
    # cat >> /root/a.txt <<EOF # 在a.txt文件后面加上三行代码
    123456789
    bbbbbbbb
    EOF # 最后一个标识符(Here-document,这里是EOF)一定要顶格写

## 扩展
### tac命令
tac是cat的反向书写, 命令的功能是反向显示文件内容(从最后一行->第一行)

### 清空文件
1. 将文件清空，文件大小为0k, 需配合`set +o noclobber`使用
```bash
$ : > filename # : 符号，它是 shell 的一个内置命令，等同于 true 命令，它可被用来作为一个 no-op（即不进行任何操作）**推荐**
$ true > my_file
$ > filename # 通过 shell 重定向 null （不存在的事物）到该文件, **推荐**
$ cat /dev/null > filename
$ cp /dev/null access.log
$ dd if=/dev/null of=access.log
$ echo -n "" > access.log # 要将 null 做为输出输入到文件中，你应该使用 -n 选项，这个选项将告诉 echo 不再像上面的那个命令那样输出结尾的那个新行
$ truncate -s 0 access.log
```

1. 清空文件，但会有一个换行，文件大小为4k(block大小)
```bash
$ echo "" > filename # 空字符串并不等同于 null
$ echo > filename
```