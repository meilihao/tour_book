# example
1. Check if run in rescue mode
```bash
if [ "`grep 'rescue' /proc/cmdline 2> /dev/null || true`" != "" ] # 舍弃stderr, [`|| true`](https://unix.stackexchange.com/questions/325705/why-is-pattern-command-true-useful/325727)保证命令的exit code是0, 因此不会中断脚本
```

1. trap
参考:
- [为shell布置陷阱：trap捕捉信号方法论](https://www.cnblogs.com/f-ck-need-u/p/7454174.html)

```bash
trap on_exit EXIT # Assign exit handler. bash 提供的一个叫做 EXIT 的伪信号(exit为0)，trap 它时当脚本因为任何原因退出时，相应的命令或函数就会执行. `on_exit`函数用`-`替换时表示不执行exit handler
```