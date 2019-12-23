# tee

## 描述

用来将标准输入的内容输出到标准输出,并将同样的内容重定向到其他文件. 即多重定向

## 格式

  tee [OPTION]... [FILE]...

## 选项

- -a：追加到文件,tee默认为覆盖.

## 例

    # tee test.txt
    #  cat a* | tee -a out.txt | cat –n
