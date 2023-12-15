# kill

## 描述

杀死进程

## 选项

- -l <信号编号> : 若不加<信息编号>选项，则-l参数会列出全部的信号名称
- -s <信号名称或编号> : 发送指定信号
- -<sigal> : 和`-s`相同, 默认15表示正常结束, 9,强制结束; 2, 结束进程但不强制, ctrl+c即是触发了`kill -2`.

## 例

```sh
$ kill PID
$ kill %job # 杀死job工作 (job为job number)
```

# killall命令
killall （kill processes by name）用于杀死进程，与 kill 不同的是killall 会杀死指定名字的所有进程。kill 命令杀死指定进程 PID，需要配合 ps 使用，而 killall 直接对进程对名字进行操作，更加方便。

```bash
killall -9 mysql         //结束所有的 mysql 进程
```

# pkill命令
pkill 命令和 killall 命令的用法相同，都是**通过进程名杀死一类进程**，除此之外，pkill 还有一个更重要的功能，即按照终端号来踢出用户登录

```bash
pkill mysql         //结束 mysql 进程
pkill -u mark,danny //结束mark,danny用户的所有进程
w  //#使用w命令查询本机已经登录的用户
pkill -9 -t pts/1  //#强制杀死从pts/1虚拟终端登陆的进程
```