# supervisor

## FAQ
### 解决 Python多进程
```conf
[program:theprogramname]
stopasgroup=false             ; 用于supervisord管理的子进程本身还有的子进程的情况, 比如解决 Python多进程问题. 如果不设置那么如果仅仅干掉supervisord的子进程的话, 子进程的子进程有可能会变成孤儿进程. 设置这个选项，会把整个该子进程的整个进程组都干掉. 该选项发送的是stop信号, 配合killasgroup后才发送kill信号.
killasgroup=false             ; SIGKILL the UNIX process group (def false). 这个和上面的stopasgroup类似, 不过发送的是kill信号
```