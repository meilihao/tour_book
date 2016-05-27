# pgrep

## 描述

获取某个进程名对应进程id

## 选项

- -l : 显示进程名称

## 例

    # pgrep gedit # 所得结果的列表是以null字符(`\0`)分隔,可用tr命令来优化显示
    # cat /proc/14502/environ |tr '\0' '\n'
