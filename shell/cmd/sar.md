# sar

## 描述

sar（System Activity Reporter系统活动情况报告）是目前 Linux 上最为全面的系统性能分析工具之一，可以从多方面对系统的活动进行报告，包括：文件的读写情况、系统调用的使用情况、磁盘I/O、CPU效率、内存使用状况、进程活动及IPC有关的活动等.

## install/安装

```shell
# ubuntu
sudo apt-get install sysstat
```

## 选项

- -o file : 将命令结果以二进制格式存放在名为file的文件中
- -q : 查看运行队列中的进程数、系统上的进程大小、平均负载等
- -r : 查看内存使用状况
- -u : 显示CPU利用率
- -W : 显示交换分区状态

## 说明

输出项说明：

CPU：all 表示统计信息为所有 CPU 的平均值。

%user：显示在用户级别(application)运行使用 CPU 总时间的百分比。

%nice：显示在用户级别，用于nice操作，所占用 CPU 总时间的百分比。

%system：在核心级别(kernel)运行所使用 CPU 总时间的百分比。

%iowait：显示用于等待I/O操作占用 CPU 总时间的百分比。

%steal：管理程序(hypervisor)为另一个虚拟进程提供服务而等待虚拟 CPU 的百分比。

%idle：显示 CPU 空闲时间占用 CPU 总时间的百分比。

1. 若 %iowait 的值过高，表示硬盘存在I/O瓶颈

2. %idle值高，表示CPU较空闲,若 %idle 的值高但系统响应慢时，有可能是 CPU 等待分配内存，此时应加大内存容量

3. 若 %idle 的值持续低于10，则系统的 CPU 处理能力相对较低，表明系统中最需要解决的资源是 CPU

## 例

```shell
$ sar 16 4 # sar 采集间隔 采集次数
Linux 3.16.0-34-generic (localhost) 	2015年04月20日 	_x86_64_	(8 CPU)

13时34分08秒     CPU     %user     %nice   %system   %iowait    %steal     %idle
13时34分24秒     all      2.03      0.00      0.62      0.04      0.00     97.31
13时34分29秒     all      2.23      0.00      0.74      0.00      0.00     97.03
13时34分35秒     all      2.10      0.00      0.73      0.04      0.00     97.13
平均时间:     all      2.08      0.00      0.67      0.03      0.00     97.22
```

## FAQ

1 . `无法打开 /var/log/sysstat/saXX: 没有那个文件或目录 Please check if data collecting is enabled in /etc/default/sysstat`

>方法1(**ubuntu推荐**):修改`/etc/default/sysstat`文件， 将 ENABLED 设置为 true,再重启系统.
>方法2:执行` sudo sar -o XX`创建文件
>方法3(fedora22):执行`sudo systemctl enable sysstat && sudo systemctl start sysstat`创建文件
