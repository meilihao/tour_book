# iostat

## 描述

iostat 被用来报告**CPU**的统计和**设备与分区的输出/输出**的统计

> iostat 工具是 sysstat 包的一部分

## 选项

- -d : 查看所有设备的 I/O 统计
- -p [设备名] : 查看所有/具体设备和它的分区的 I/O 统计
- -x : 显示所有设备的详细的 I/O 统计信息
- -m : 以 MB 为单位而不是 KB 查看所有设备的统计. 默认以 KB 显示输出
- -N : 查看 LVM 磁盘 I/O 统计报告

## 例

    # iostat
    # iostat 5 2 # 打算以 5 秒捕获的间隔捕获两个报告, iostat [Interval] [Number Of Reports], 使用特定的间隔输出
