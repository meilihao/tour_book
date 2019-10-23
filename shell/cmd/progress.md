# progress
用于显示任何核心组件命令（如：cp、mv、dd、tar、gzip、gunzip、cat、grep、fgrep、egrep、cut、sort、xz、exiting）的进度. 它使用文件描述信息来确定一个命令的进度.

## 选项
- -q (--quiet) : 隐藏所有消息
- -d (--debug) : 显示所有警告/错误消息
- -w (--wait) : 估计I / O吞吐量和估计剩余时间（显示较慢）             
- -W (--wait-delay secs) : I / O估计等待'秒'秒（暗示-w）
- -m (--monitor) : 在受监视的进程仍在运行时循环       
- -M (--monitor-continuously) : 监视永不停止（类似于观看进度）
- -c (--command cmd) : 仅监视此命令名称（例如：firefox）, 可以在命令行上多次使用此选项               
- -a (--additional-command cmd) : 将此命令添加到默认列表, 可以在命令行上多次使用此选项  
- -p (--pid id) : 仅监视此数字进程ID（例如：`pidof firefox`）, 可以在命令行上多次使用此选项       
- -o (--open-mode {r|w}) : 仅报告为进程读取或写入而打开的文件. 当您只想监视进程的输出文件（或输入文件）时，此选项很有用.                 
- -i (--ignore-file file) : 不报告“文件”的过程. 如果文件尚不存在，则必须提供完整且干净的绝对路径, 可以在命令行上多次使用此选项.

## example
```
$ sha256sum SHA256.txt & progress -mp $! # 使用$！启动并监控命令并重定向输出
```