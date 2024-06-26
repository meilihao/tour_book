# SHELL
当前使用的是哪种shell

# HOME
用户的主目录

# TIME_STYLE

自定义显示时间格式,例如控制ls命令输出时间的格式

export TIME_STYLE='+%Y-%m-%d %H:%M:%S'

# TMOUT
设置`<N>`秒内用户无操作就自动断开终端

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

sh的提示字符串.可扩展的参数:`\u`:用户名;`\h`:主机名;`\w`:当前工作目录.

> 可通过全局配置文件`/etc/bash.bashrc`或`/etc/profile`进行调整
> ubuntu下直接修改无效,需修改.bash_profile等配置文件来生效

# IFS(内部字段分隔符，Internal Field Separator)

存储定界符,FS默认为空白字符（换行符，制表符，空格)

```sh
#!/bin/bash
# shell字段分隔符IFS，以逗号（，）为分隔符获取字符串内容
function test_for
{
    ifs_old=$IFS
    IFS=$','
    for i in $(echo "${1}")
    do
    echo "${i}"
    done
}
    
test_for "test1,test2"

IFS=';' read -ra ADDR <<< "$IN"
for i in "${ADDR[@]}"; do
    # 处理 "$i"
done

# 类似的单行命令
while IFS=';' read -ra ADDR; do
    for i in "${ADDR[@]}"; do
        # 处理 "$i"
    done
done <<< "$IN"
```

# PWD

存储当前的工作目录

# LANG
当前系统使用的字符集, 在`/etc/default/locale`里.
可通过`locale -a`查看可用的字符集.

# 其他环境变量
- SHELL : 当前使用的shell
- HISTSIZE : 输出的历史命令记录条数
- HISTFILESIZE : 保存的历史命令记录条数
- LANG : 系统语言、语系名称
- RANDOM 生成一个随机数字
- PATH 定义解释器搜索用户执行命令的路径
- EDITOR 用户默认的文本编辑器
- TERM 当前session终端的名称, 更加名称可知它的特性
- LINES : 当前终端可显示的行数
- COLUMNS : 当前终端可显示的列数

# XDG_SESSION_TYPE
当前使用的显示服务器类型

> loginctl + `loginctl show-session [SESSION_ID] | egrep 'Display|Type'`(SESSION_ID from `loginctl`的输出), 看结果包含x11还是wayland

# XDG_CURRENT_DESKTOP
当前使用的桌面类型, 比如UKUI, GNOME

# 配置文件
`bashrc`与`profile`都用于保存用户的环境信息，`bashrc`用于非交互式`no-login shell`,比如x-window下启动的终端，而profile用于交互式`login shell`.

`/etc/profile，/etc/bashrc`是系统全局环境变量设定;`~/.profile，~/.bashrc`用户目录下的私有环境变量设定.

> 注意: 切换sh,比如 fish -> bash, 使用`no-login`操作,仅执行`~/.bashrc`.

## profile和bashrc的区别
### profile
profile 是某个用户唯一的用来设置环境变量的地方, 因为用户可以有多个 shell 比如 bash, sh, zsh 之类的, 但像环境变量这种其实只需要在统一的一个地方初始化就可以了, 而这就是 profile.
### bashrc

bashrc是专门用来给 bash 做初始化的比如用来初始化 bash 的设置, bash 的代码补全, bash 的别名, bash 的颜色. 以此类推也就还会有 shrc, zshrc 这样的文件存在了, 只是 bash 太常用了而已.

## 流程

参考:[理解 Linux/Unix 登录脚本](https://www.sdk.cn/news/5585)

### login shell
"login shell"代表用户登入, 比如使用`su - xxx`命令, 或者用 ssh 连接到某一个服务器上, 都会使用该用户默认 shell 启动 login shell 模式.

1. login
2. `/etc/profile`,根据其内容读取额外的文档，如/etc/profile.d和/etc/inputrc等
3. 个人配置文件(`~/.bash_profile`,`/.bash_login`和`~/.profile`,主要是获取与用户有关的环境、别名和函数等)
    - 在列出的顺序中第一个被找到的文件会被作为配置文件，其余的都会被忽略

    > ubuntu使用`~/.profile`加载`~/.bashrc`
4. 如果`~/.bashrc`(其会调用`~/.bash_alias`)存在的话，`~/.bash_profile`还会调用它
5. /etc/bashrc
6. `~/.bash_logout`

ps :

其他的shell，例如Dash，支持相似的东西，但是只会查找~/.profile文件。这允许用户为Bash特定的应用场景配置单独的.bash_profile文件，如果在某些时候需要切换到Dash或其他shell作为登录shell（例如通过chsh -s dash命令）。可以保留~/.profile作为这些shell的配置文件。

需要牢记的一点是，默认的Debian框架目录（/etc/skel，用于存放要复制到新用户账户主目录的文件和目录）包含.profile文件，但不包含.bash_profile和.bash_login文件。此外Debian使用Bash作为默认的shell，因此，许多Debian用户习惯于将他们的Bash 登录shell设置放在.profile文件中.

> 不同的发行版可能有点不同,具体可以看其相关脚本的内容.

### non-login shell
no-login shell是在终端下直接输入 bash 或者`bash -c "CMD"`来启动的 shell

1. `~/.bashrc`
2. /etc/bashrc(或/etc/bash.bashrc)
3. /etc/profile.d/*.sh

## ~/.profile与~/.bashrc的区别:

这两者都具有个性化定制功能,但

- `~/.profile`可以设定本用户专有的路径，环境变量等，它只能登入的时候执行一次
- `~/.bashrc`也是某用户专有设定文档，可以设定路径，命令别名，每次shell script的执行都会使用它一次.

## X11

见`/etc/X11/Xsession`

### /etc/profile、/etc/bash.bahsrc、~/.profile、~/.bashrc的用途
> /etc/bashrc(ubuntu/debian没有这个文件，对应地，其有/etc/bash.bashrc文件)

打开一个新的shell（包括打开一个新终端和在终端上输入bash），都会重新读取/etc/bash.bashrc 和 ~/.bashrc文件里面的内容.


/etc/profile、/etc/bash.bashrc文件是针对所有用户来说的，每个用户登录时都会执行，其中/etc/profile只执行一次，而/etc/bash.bashrc在每次Shell登录时都会执行.
~/.profile、~/.bashrc文件是针对单个用户来说的，每个用户目录下都会有这两个文件，其中~/.profile在Login Shell登录时执行，~/.bashrc在Non-login Shell登录时执行.

![](/misc/img/shell/20170405223231326.png)

## FAQ
### web terminal bash 输入后退键变空格
解决方法: 缺`ncurses-base`或使用`export TERM=linux`/`export TERM=xterm-256color`
