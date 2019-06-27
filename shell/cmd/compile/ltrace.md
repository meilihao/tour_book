# ltrace
library trace(ltrace,跟踪库函数调用)与 strace(跟踪系统调用)非常类似, ltrace 会解析共享库(可执行文件的动态段),即一个程序的链接信息,并打印出共享库和静态库的实际符号和函数, 还可以使用`-S` 选项查看系统调用

## 类似
[function trace(ftrace,函数追踪)](https://github.com/elfmaster/ftrace), 它的功能与 ltrace 类似,但还可以显示出二进制文件本身的函数调用

## 实践
- [谁偷走了你的服务器性能——Strace & Ltrace篇](http://rdc.hundsun.com/portal/article/597.html)