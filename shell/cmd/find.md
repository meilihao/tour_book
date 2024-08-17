# find
参考:
- [Linux find运行机制详解](https://www.cnblogs.com/f-ck-need-u/p/6995529.html)

## 描述

查找目录或者文件

## 格式

    find [-H] [-L] [-P] [-D debugopts] [-Olevel] [starting-point...] [expression]

`[-H] [-L] [-P]`如何处理符号链接; `[starting-point...]`要查找的路径; `[expression]` = 参数(`options`) + 限定条件(`tests`) + 执行的动作(`actions`)
`-D`调试find逻辑时很有用

## 选项

options:
- -depth : 从指定目录的最深层开始往上查找, 即搜索到目录时，先处理目录中的文件(子目录)，再处理目录本身. 对于"-delete"这个action，它隐含"-depth"选项
- -maxdepth level: 限制查找的最大深度
- -mindepth level: 限制开始查找的最小深度
- -follow : 排除符号链接
- -empty : 查找内容为空的文件
- -regextype type : 修改正则表达式的模式, 默认是emacs, 还支持posix-awk,posix-basic,posix-egrep, posix-extended

> `-maxdepth`和`-mindepth`应该作为find的第3个参数出现,在之后出现可能会影响查找效率.

tests:
- -depth : 先处理当前目录下的子`文件/目录`, 然后再是自身
- -type : 文件类型
- -uid <N> : 指定uid
- -gid <N> : 指定gid
- -user <'string'>: 指定用户名
- -group <'string'>: 指定组名
- -empty : 指定大小为0
- -size : 文件大小,`+ 大于  -小于   无符号，恰好等于`, 比如`+50K`
- -regex : 使用正则表达式来匹配**文件路径**(给出的expr必须要匹配完整的文件路径),因此`find . -regex "a"`与`find . |grep a`的结果可能会不一致
- -iregex : 和`-regex`类似,但匹配时忽略大小写
- -path pattern: 按照指定**的路径名(并不是单纯的字符串匹配)**来匹配**文件路径**,支持通配符("*"、"?"和"[]", 但"/"或"."是特殊字符会跳过);不使用通配符时表示匹配具体路径. **`-path`可配合`-prune`用于排除指定目录**
- -perm ： 查找符合指定的权限数值的文件或目录
- -name : 按照指定的字符串来匹配**文件名**,仅支持通配符`*,?,[]`
- -lname : 按照指定的字符串来匹配**符号链接的文件名**,仅支持通配符`*,?,[]`
- -iname : 和`-name`类似,但忽略大小写
- -newer : 查找其更改时间较指定文件或目录的更改时间更接近现在的文件或目录. 比如`-newer f1 !f2`匹配比f1新但比f2旧的文件
- [amc]time [-n|n|+n] : 按照指定的文件时间戳来查找. `-n`: 时间距现在n天内; `+n`:时间距现在n天前, `n`:距现在n天
- [amc]min [-n|n|+n] : 与`[amc]time [-n|n|+n]`相同, 但时间单位是分钟
- -nogroup : 查找无有效所属组的文件，即该文件所属的组在/etc/groups中不存在
- -nouser : 查找无有效属主的文件，即该文件的属主在/etc/passwd中不存在

actions:
- -delete : 删除匹配的文件,**`-delete`位置一定在最后**
- `-exec cmd {}\;`： 假设find指令的返回值为True，就执行该指令.**`-exec`后只能接受单个命令,多个命令时使用shell脚本(`-exec ./cmd.sh {} \;`)**,`{}`是占位符,`;`表示结束. 为了效率推荐使用xargs, 特别是文件多的时候, 因为`-exec`是逐个传递查询结果
- -ok : 与`-exec`作用相同, 但执行的每个命令都需用户确认
- -prune : 忽略(即跳过)**指定目录**. 因为要忽略的是目录，所以一般都只和`-path`配合. 没有给定`-depth`时，总是返回true，如果给定`-depth`，则直接返回false，所以`-delete`(隐含了`-depth`)是不能和`-prune`一起使用的.
- -print : 默认action, 将查找结果(即包含路径的文件名)导入标准输出,并用`\n`换行
- -print0 : 和`-print`类似,但用`\0`换行,在文件名包含换行符时可用

> 如果find**评估完所有表达式**后发现没有action(`-prune`这个action除外)，则在**最末尾加上-print作为默认的action**. 还需注意，如果只有`-prune`这个action，它还是会补上`-print`. 比如假设`./a`是非空目录, 那么`find . -path "./a" -prune` = `find . -path "./a" -prune -print`, 因为`-path "./a"`和`-prune `均为true, 因此即使有`-prune`还是会输出`./a`
> **`-prune`的一个弱点是不适合通过通配符来忽略目录**，因为通配符出来的很可能导致非预期结果,因此想要忽略多个目录，最好的方法是多次使用`-path`

operators, find支持的逻辑运算符:
- ! : not, 取反
- -a : and, 取交集
- -o : or, 取并集

> and的优先级高于or, 即`expr1 -o expr2 -a expr3`等价于`expr1 -o (expr2 -a expr3)`

文件类型:

- b ：块设备
- c ：字符设备
- d ：目录
- f ：普通文件
- l ：符号链接
- p ：管道文件
- s ：Socket

文件大小

- b : 块(512字节)
- c : 字节
- w : 字（2字节）
- k : 千字节
- M : 兆字节
- G : 吉字节

## 例

    # find / \(-path /var/log -o -path /usr/bin\) -prune -o -name "main.c" -pirnt # 在/下查找不包括/var/log和/usr/bin的所有普通文件
    # find / -path "/usr/bin" -prune -o -name "main.c" -mtime +2 -print # 在/但不包括/usr/bin下查找两天前的名为main.c的普通文件
    # find . -atime -2 # 查找两天内访问过的问题
    # find . -perm 755 # 查找权限是755的文件
    # find . -size +1000c # 查找大小超过1000字节的文件
    # find . path "/data/a" -prune -o -print # 查找时排除目录"/data/a"
    # find path # 列出(递归)当前目录及子目录下所有的文件和文件夹
    # find . -print # 输出效果和上面的命令相同
    # find /root -name "*.txt" -print
    # find . \( -name "*.txt" -o -name "*.py" \) # 多个条件 or,使用转义避免内容被bash当做子命令来解析
    # find /home/users -path "*slynux*" -print
    # find . -regex ".*\(\.py\|\.sh\)$"
    # find . ! -name "*.txt" -print # 匹配所有不以`.txt`结尾的文件
    # find . -maxdepth 1 -type f -print # 仅查找一层,即只列出当前目录下的文件
    # find . -mindepth 2 -type f -print # 从距离当前目录至少两层的子目录开始查找
    # find . -type f -atime -7 -print # 最近七天内被访问的文件
    # find . -type f -atime 7 -print  # 恰好在七天前被访问的文件
    # find . -type f -atime +7 -print # 超过七天前被访问的文件
    # find . -type f -newer file.txt -print # 比参考文件更新(即修改时间更接近当前时间)的文件
    # find . -type f -size +2k # 大于2K
    # find . -type f -perm 644 -print
    # find . -type f -user slynux -print
    # find . -type f -user slynux -exec chown slynux {} \; # `{}`表示查询结果的占位符,其两边有空格;`\;`表`-exec`指定命令的结束
    # find . -type f -name "*.txt" -exec cat {} \;> all.txt # 这里没用`>>`是因为find命令的全部输出都是单数据流(stdin)
    # find . -type f -name "*.txt" -exec printf "File: %s\n" {} \; # `-exec`和`-printf`结合来格式化输出
    # find . -name ".git" -prune -o -name "*.txt" -print # 找出不在'.git'文件夹内的所有txt文件
    # find . -type f -name "*.swp" -delete <=> find . type f -name "*.swp" | xargs rm
    # find . -type f|xargs file|grep "CRLF" # 查找所有使用`CRLF`换行的文件
    # find . -newer f1 ! -newer f2 # 查找更改时间比f1新但比f2旧的问题
    # find -D rates . -path "./a" -prune -o -name a # 使用`-D`调试find逻辑
    # find . -name "__pycache__" -type d -print -exec rm -rf {} \;
