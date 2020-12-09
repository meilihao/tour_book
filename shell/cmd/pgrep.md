# pgrep

## 描述

获取某个进程名对应进程id

## 选项

- -l : 显示进程名称
- -o : 显示进程起始的id
- -n : 显示进程终止的id
- -f : 匹配command中的关键字
- -G : 匹配指定用户组的pid
- -x : 精确匹配

## 例

    # pgrep gedit # 所得结果的列表是以null字符(`\0`)分隔,可用tr命令来优化显示
    # cat /proc/14502/environ |tr '\0' '\n'
    # pgrep -f sshd
    # pgrep -G www_data
    # pgrep -x "cdp" 2>/dev/null # 精确查找名为cdp的进程

