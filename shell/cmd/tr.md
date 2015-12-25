# tr

## 描述

可以对来自标准输入的字符进行替换，删除以及压缩(translate, 可以将一组字符变成另一组字符)
tr只能通过stdin，无法通过其他命令行进行接收参数.

## 格式

    tr [options] source-char-set replace-char-set

字符类

    alnum 字母和数字
    alpha 字母
    cntrl 控制字符
    digit 数字
    graph 图形字符
    lower 小写字母
    print 可打印字符
    punct 标点符号
    space 空白字符
    upper 大写字母
    xdigit 十六进制字符    

## 选项

- c : 取source-char-set补集，通常与-d/-s配合
- d : 删除字source-char-set中的所列的字符
- s : 浓缩重复字符，连续多个变成一个

## 例

    # cat /proc/12501/environ | tr '\0' '\n' # 字符替换
    # echo  "HELLO" | tr 'A-Z' 'a-z' # 大小写替换
    # cat text | tr '\t' '' # 删除字符
    # echo "hello 123 world 456"| tr -d '0-9' # 取字符集的补集
    hello  world
    # echo "hello 1 char 2" | tr -d -c '0-9'  #删除非0-9
    12
    # echo "GNU is    not UNix" | tr -s ' ' # 压缩字符(连续的重复字符)
    GNU is not UNix
    # echo "GNU is not UNix" | tr '[:lower:]' '[:upper:]'
    GNU IS NOT UNIX
