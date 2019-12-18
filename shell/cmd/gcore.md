# gcore

## 描述

dump linux 进程.
前提: 使用调试参数编译程序（例如： gcc中使用"-g"选项），编译后不要去除文件的调试符号信息.

> from:　`apt install gdb`

## 例

    # gcore ${pid}
    # gdb core.`${pid}`
