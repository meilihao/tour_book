# exec

## 描述

exec命令用于调用并执行指令的命令,或创建自定义的文件描述符

exec并不启动新的shell，而是用要被执行命令替换当前的shell进程，并且将老进程的环境清理掉,而且exec命令后的其它命令将不再执行.

因此，如果你在一个shell里面，执行exec ls那么，当列出了当前目录后，这个shell就自己退出了，因为这个shell进程已被替换为仅仅执行ls命令的一个进程，执行结束自然也就退出了。为了避免这个影响,一般将exec命令放到一个shell脚本里面，用主shell调用这个脚本，调用点处可以用bash a.sh，（a.sh就是存放该命令的脚本），这样会为a.sh建立一个sub shell去执行，当执行到exec后，该子脚本进程就被替换成了相应的exec的命令.

## 例子

 # exec 3<input.txt # 使用文件描述符3打开并读取文件
 # cat<&3 # 在命令中使用文件描述符3
 # exec 4>a.txt # 同理也可用`>>`代替`>`
 # echo newline >&4
 # cat a.txt
 newline
