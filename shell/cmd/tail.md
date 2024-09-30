# tail
ref:
- [src](https://github.com/coreutils/coreutils/blob/master/src/tail.c)

    有多种实现:
    1. 一种是定时去检测文件化，然后输出到控制台 by tail_forever 轮询式跟踪

        检测主要依赖于函数: fstat (fd, &stats) , 通过比较 mtime 进行文件是否变化判定

    1. 第二种则高级些，利用操作系统提供的文件通知功能，进行实时内容输出 by tail_forever_inotify 更新时通知

        它会先进行 inotify_init(); 然后再进行 tail_forever_inotify. 即基于通知的方式实现文件跟踪，首先是注册事件，然后是轮询事件，然后是事件处理.

## 描述

显示文件中的尾部内容, 默认显示10行

## 格式

    tail [OPTION]... [FILE]...

## 选项

- -c<数字> ：显示指定的字节数
- -n<数字> ：显示指定的行数
- -f : 显示文件最新追加的内容,即实时追踪文件的变化.
- -F : = `-f --retry`
- -n : 输出文件的尾部N（N位数字）行内容
- --retry : 不断尝试打开文件直至成功
- --pid=<pid> : 与-f参数连用, 在进程结束后自动退出tail命令
- -s <秒数N> : 监控文件变化的间隔时间
- -v ： 总是显示文件名的头信息
- -q ： 不显示文件名的头信息

## 例
```shell
tail -n 5 m.txt //显示最后5行
tail -n -5 m.txt //显示最后5行
tail -n +5 m.txt //从第5行开始显示
ps -ef
sudo tail  /proc/<pid>/fd/1 # 查看运行进程的输出
strace -p {pid} -e write # 同上
```