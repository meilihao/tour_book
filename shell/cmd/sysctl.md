# sysctl

## 描述

设置内核参数.

## 格式

    sysctl [options] [variable[=value]] [...]

## 选项

- -a : 打印全部参数
- -w : 设置参数, 设置后即已生效
- -p[file] : 参数使能 

## 例

    # sysctl -w net.core.rmem_max=67108864
