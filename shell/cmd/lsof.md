# lsof

## 描述

查看当前系统文件的工具,但要记住：linux 下 “一切皆文件”.

lsof打开的文件有：

- 普通文件
- 目录
- 网络文件系统的文件
- 字符或设备文件
- (函数)共享库
- 管道，命名管道
- 符号链接
- 网络文件（例如：NFS file、网络socket，unix域名socket）
- 还有其它类型的文件，等等

## 选项

- -a : 列出打开文件的进程
- -c 进程名 : 列出指定进程所打开的文件
- -d <文件号> : 列出占用该文件号的进程
- +d 目录 : 列出目录下被打开的文件
- +D 目录 : 递归列出目录下被打开的文件
- -i条件 : 列出符合条件的进程（4、6、协议、:端口、 server,进程的服务名,比如nfs, ssh、 @ip ）
- -p 进程号 : 列出指定进程号所打开的文件
- -u 用户名/uid : 查看用户username的进程所打开的文件
- -w : 不打印警告信息

###  lsof输出各列信息的意义

- COMMAND：进程的名称
- PID：进程标识符
- PPID：父进程标识符（需要指定-R参数）
- USER：进程所有者
- FD：文件描述符，应用程序通过文件描述符识别该文件
```
（1）cwd：表示current work dirctory，即：应用程序的当前工作目录，这是该应用程序启动的目录，除非它本身对这个目录进行更改
（2）txt ：该类型的文件是程序代码，如应用程序二进制文件本身或共享库，如上列表中显示的 /sbin/init 程序
（3）lnn：library references (AIX);
（4）er：FD information error (see NAME column);
（5）jld：jail directory (FreeBSD);
（6）ltx：shared library text (code and data);
（7）mxx ：hex memory-mapped type number xx.
（8）m86：DOS Merge mapped file;
（9）mem：memory-mapped file;
（10）mmap：memory-mapped device;
（11）pd：parent directory;
（12）rtd：root directory;
（13）tr：kernel trace file (OpenBSD);
（14）v86  VP/ix mapped file;
（15）0：表示标准输出
（16）1：表示标准输入
（17）2：表示标准错误
 (18) n : 程序的文件描述符，这是打开该文件时返回的一个非负整数,因为`0,1,2`由POSIX定义,每个进程都有的，所以大多数应用程序所打开的文件的FD都是从3开始 
一般在标准输出、标准错误、标准输入后还跟着文件状态模式：r、w、u等
（1）u：表示该文件被打开并处于读取/写入模式
（2）r：表示该文件被打开并处于只读模式
（3）w：表示该文件被打开并处于
（4）空格：表示该文件的状态模式为unknow，且没有锁定
（5）-：表示该文件的状态模式为unknow，且被锁定
同时在文件状态模式后面，还跟着相关的锁
（1）N：for a Solaris NFS lock of unknown type;
（2）r：for read lock on part of the file;
（3）R：for a read lock on the entire file;
（4）w：for a write lock on part of the file;（文件的部分写锁）
（5）W：for a write lock on the entire file;（整个文件的写锁）
（6）u：for a read and write lock of any length;
（7）U：for a lock of unknown type;
（8）x：for an SCO OpenServer Xenix lock on part      of the file;
（9）X：for an SCO OpenServer Xenix lock on the      entire file;
（10）space：if there is no lock.
```
- TYPE：文件类型
```
（1）DIR：表示目录
（2）CHR：表示字符类型
（3）BLK：块设备类型
（4）UNIX： UNIX 域套接字
（5）FIFO：先进先出 (FIFO) 队列
（6）IPv4：网际协议 (IP) 套接字
```
- DEVICE：指定磁盘的名称
- SIZE：文件的大小
- NODE：索引节点（文件在磁盘上的标识）
- NAME：打开文件的确切名称

## 例

    # lsof -p ${pid} 2>&1 # 获取进程打开的文件列表
    # lsof -i@192.168.1.6
    # lsof -i udp@127.0.0.1:53 # 使用本机localhost+udp+53端口的进程
    # lsof -i tcp:25 # 使用tcp+25端口的进程
    # lsof -i :3306 # 查看端口占用的进程状态
    # lsof -u username # 查看某个用户打开的文件
    # lsof /bin/bash # 查看使用某个文件的进程
    # lsof -c mysql # 查看某个进程打开的文件
    # losf -u username -c mysql # 查看某个用户以及某个进程所打开的文件
    # losf -g gid : 指定gid打开的文件
    # losf -i # 查看所有的网络连接
    # losf -i tcp # 查看所有tcp网络连接信息
    # losf -i tcp:3306 # 查看3306端口上tcp网络连接信息
    # lsof -a -u test -i # 列出某个用户的所有使用的网络端口
    # lsof -d file_description(like 2) # 根据文件描述列出对应的文件信息
    # lsof -i 4 -a -p 1234 # 列出被进程号为1234的进程所打开的所有IPV4 network files
    # lsof -i @nf5260i5-td:20,21,80 -r 3 # 列出目前连接主机nf5260i5-td上端口为：20，21，80相关的所有文件信息，且每隔3秒重复执行
