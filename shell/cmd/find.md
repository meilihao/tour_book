# find

## 描述

查找目录或者文件

## 选项

- -delete : 删除匹配的文件,**`-delete`位置一定在最后**
- -exec ： 假设find指令的返回值为True，就执行该指令.**`-exec`后只能接受单个命令,多个命令时使用shell脚本(`-exec ./cmd.sh {} \;`)**
- -iname : 和`-name`类似,但忽略大小写
- -iregex : 和`-regex`类似,但匹配时忽略大小写
- -maxdepth : 限制查找的最大深度
- -mindepth : 限制查找的最小深度
- -name : 按照指定的字符串来匹配**文件名**,支持通配符
- -newer : 查找其更改时间较指定文件或目录的更改时间更接近现在的文件或目录
- -path : 按照指定的字符串来匹配**文件路径**,支持通配符
- -perm ： 查找符合指定的权限数值的文件或目录
- -print : 将查找结果(即包含路径的文件名)导入标准输出,并用`\n`换行
- -print0 : 和`-print`类似,但用`\0`换行,在文件名包含换行符时可用
- -prune : 忽略(即跳过)某个目录
- -regex : 使用正则表达式来匹配**文件路径**
- -size : 文件大小,`+ 大于  -小于   无符号，恰好等于`
- -type : 文件类型
- -user : 指定用户(用户名/UID)
- ! : 表示否定参数

> `-maxdepth`和`-mindepth`应该作为find的第3个参数出现,在之后出现可能会影响查找效率.

文件类型:

- b ：块设备
- c ：字符设备
- d ：目录
- f ：普通文件
- l ：符号链接
- p ：管道文件
- s ：Socket

时间戳

--- 计量单位: 天
- -atime : 访问时间,最近一次访问文件的时间
- -mtime : 修改时间,最后一次被修改的时间
- -ctime : 变化时间,文件元数据(metadata,权限或所有权等)最近一次改变的时间

--- 计量单位: 分钟
- -amin : 访问时间
- -mmin : 修改时间
- -cmin : 变化时间

文件大小

- b : 块(512字节)
- c : 字节
- w : 字（2字节）
- k : 千字节
- M : 兆字节
- G : 吉字节

## 例

    # find path # 列出(递归)当前目录及子目录下所有的文件和文件夹.path为路径
    # find . -print # 输出效果和上面的命令相同
    # find /root -name "*.txt" -print
    # find . \( -name "*.txt" -o -name "*.py" \) # 多个条件 or
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
