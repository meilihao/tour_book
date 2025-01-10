# lscpu

 lscpu 显示CPU 的架构信息. lscpu 从 sysfs 和 proc/cpuinfo 中收集信息.

## 选项
- -a : 打印在线和离线CPU(默认为-e)
- -e : 打印出一个扩展的可读格式
- -p : 打印出可解析的格式

## 字段
           CPU  逻辑CPU编号
          CORE  逻辑核心号码
        SOCKET  逻辑套接字号
          NODE  逻辑NUMA节点号
          BOOK  逻辑书号
         CACHE  显示了如何在CPU之间共享高速缓存
  POLARIZATION  虚拟硬件上的CPU调度模式
       ADDRESS  CPU的物理地址
    CONFIGURED  显示管理程序是否分配了CPU
        ONLINE  显示Linux是否正在使用CPU
Virtualization  虚拟化支持情况
