# tee

## 描述

用来将标准输入的内容输出到标准输出,并将同样的内容重定向到其他文件. 即多重定向

存在缓存机制，每1024个字节将输出一次. 若从管道接收输入数据，应该是缓冲区满，才将数据转存到指定的文件中.

## 格式

  tee [OPTION]... [FILE]...

## 选项

- -a：追加到文件,tee默认为覆盖.

## 例

    # tee test.txt
    #  cat a* | tee -a out.txt | cat –n
    # ./a.sh  2>&1 | tee  l.log # 测试script常用, 但不适合用于内含dialog命令的script.
