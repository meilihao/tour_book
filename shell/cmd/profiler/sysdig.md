# sysdig
ref:
- [Sysdig速查手册 Sysdig常用命令参数参考](https://xunyidian.com/t/SysdigManual)

sysdig是一个能够让系统管理员和开发人员以前所未有方式洞察其系统行为的监控工具.

它可以认为是"strace+tcpdump+lsof+lua的集合"，sysdig的最棒特性之一在于，它不仅能分析Linux系统的"现场"状态，也能将该状态保存为转储文件以供离线检查

原理: Sysdig 通过在内核的 driver 模块注册系统调用的 hook，这样当有系统调用发生和完成的时候，它会把系统调用信息拷贝到特定的 buffer，然后用户模块的组件对数据信息处理（解压、解析、过滤等），并最终通过 Sysdig 命令行和用户进行交互.

## example
```bash
# sysdig -c topprocs_cpu # 查看top cpu占用, 可显示已隐藏的进程(查找隐藏进程的其他工具还有`unhide proc`)
# sysdig -c topprocs_net # 查看占用网络带宽最多的进程
# sysdig -c topprocs_file # 查看使用硬盘带宽最多的进程
# sysdig -c topfiles_bytes # 根据读取+写入字节查看top文件
# sysdig -c topfiles_bytes proc.name=httpd #打印apache一直在阅读或写入的top文件
# sysdig -c fdcount_by proc.name "fd.type=file" # 列出使用大量文件描述符的进程
# sysdig -c fdbytes_by fd.directory # 根据R + W磁盘活动查看top目录
# sysdig -c fdbytes_by fd.filename "fd.directory=/tmp/" # 查看/tmp目录中有关R + W磁盘活动的顶级文件
# sysdig -A -c echo_fds "fd.filename=passwd" # 查看passwd'的所有文件的I/O活动
# sysdig -s2000 -X -c echo_fds fd.cip=192.168.0.1  # 显示主机192.168.0.1的网络传输数据, 以16进制显示
# sysdig -s2000 -A -c echo_fds fd.cip=192.168.0.1 # 同上, 但以ASCII显示
# sysdig -c fdcount_by fd.sport "evt.type=accept" # 查看连接最多的服务器端口
# sysdig -c fdcount_by fd.cip "evt.type=accept"  # 查看客户端连接最多的ip
# sysdig -p"%proc.name %fd.name" "evt.type=accept and proc.name!=httpd" # 列出所有不是访问apache服务的访问连接
# csysdig -vcontainers # 查看机器上运行的容器列表及其资源使用情况
# csysdig -pc # 查看容器上下文的进程列表
# sysdig -pc -c topprocs_cpu container.name=wordpress1 # 查看运行在wordpress1容器里CPU的使用率
# sysdig -pc -c topprocs_net container.name=wordpress1 # 运行在wordpress1容器里网络带宽的使用率
```