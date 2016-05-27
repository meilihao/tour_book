# TIME_STYLE

自定义显示时间格式,例如控制ls命令输出时间的格式

export TIME_STYLE='+%Y-%m-%d %H:%M:%S'

# LD_LIBRARY_PATH

Linux环境变量名，该环境变量主要用于指定查找共享库（动态链接库）时除了默认路径之外的其他路径(该路径在默认路径之前查找).
移植程序时的经常碰到需要使用一些特定的动态库，而这些编译好的动态库放在我们自己建立的目录里，这时可以将这些目录设置到LD_LIBRARY_PATH中.

```
# NEW_DIRS是新的动态库路径
export LD_LIBRARY_PATH=NEW_DIRS: $LD_LIBRARY_PATH
```

# UID

用户Id,root的UID是0.

# PS1

sh的提示字符串.可扩展的参数:`\u`:用户名;`\h`:主机名;`\w`:当前工作目录.ubuntu下直接修改无效,需修改.bash_profile等配置文件来生效.

# IFS(内部字段分隔符，Internal Field Separator)

存储定界符,FS默认为空白字符（换行符，制表符，空格

# PWD

存储当前的工作目录

# 配置文件

`bashrc`与`profile`都用于保存用户的环境信息，`bashrc`用于交互式`non-login shell`,比如x-window下启动的终端，而profile用于交互式`login shell`.

`/etc/profile，/etc/bashrc`是系统全局环境变量设定;`~/.profile，~/.bashrc`用户目录下的私有环境变量设定.

## 流程

### login shell

1. login
2. `/etc/profile`,根据其内容读取额外的文档，如/etc/profile.d和/etc/inputrc等
3. 个人配置文件(`~/.bash_profile`,`/.bash_login`和`~/.profile`,主要是获取与用户有关的环境、别名和函数等)
    - 如果`~/.bash_profile`存在，那么bash就不会理睬其他两个文件.
    - 如果`~/.bash_profile`不存在，bash才会读取`~/.bash_login`.
    - 而前两个文件都不存在的话，bash才会读取`~/.profile`文件.
4. 如果`~/.bashrc`存在的话，`~/.bash_profile`还会调用它
5. bash启动

### non-login shell

1. `~/.bashrc`
2. /etc/bashrc(或/etc/bash.bashrc)
3. bash启动

## ~/.profile与~/.bashrc的区别:

这两者都具有个性化定制功能,但

- `~/.profile`可以设定本用户专有的路径，环境变量等，它只能登入的时候执行一次
- `~/.bashrc`也是某用户专有设定文档，可以设定路径，命令别名，每次shell script的执行都会使用它一次.
