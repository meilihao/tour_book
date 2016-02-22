# cat

## 描述

合并文件内容并输出到stdout.

## 选项

- -n : 显示行号
- -s : 将两行及以上的空白行压缩为一行
- -T : 用`^I`代替隐藏的`TAB`并显示出来

## 例

    # cat test.txt
    # cat test1.txt test2.txt
    # OUTPUT_FROM_SOME_CMDS | cat # 从管道中读取
    # echo "test" | cat - file1 # 拼接标准输入和文件的内容,这里的`-`表示stdin的重定向的源或目的,简单理解即是stdin.
    # cat >> /root/a.txt <<EOF # 在a.txt文件后面加上三行代码
    123456789
    bbbbbbbb
    EOF # 最后一个标识符(Here-document,这里是EOF)一定要顶格写
