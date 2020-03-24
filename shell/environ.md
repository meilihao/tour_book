# SHELL
当前使用的是哪种shell

# HOME
用户的主目录

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
- HISTSIZE : 输出的历史命令记录条数
- HISTFILESIZE : 保存的历史命令记录条数
- LANG : 系统语言、语系名称
- RANDOM 生成一个随机数字
- PATH 定义解释器搜索用户执行命令的路径
- EDITOR 用户默认的文本编辑器

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

1. login
2. `/etc/profile`,根据其内容读取额外的文档，如/etc/profile.d和/etc/inputrc等
3. 个人配置文件(`~/.bash_profile`,`/.bash_login`和`~/.profile`,主要是获取与用户有关的环境、别名和函数等)
    - 在列出的顺序中第一个被找到的文件会被作为配置文件，其余的都会被忽略
4. 如果`~/.bashrc`(其会调用`~/.bash_alias`)存在的话，`~/.bash_profile`还会调用它
5. bash启动
6. `~/.bash_logout`

ps :

其他的shell，例如Dash，支持相似的东西，但是只会查找~/.profile文件。这允许用户为Bash特定的应用场景配置单独的.bash_profile文件，如果在某些时候需要切换到Dash或其他shell作为登录shell（例如通过chsh -s dash命令）。可以保留~/.profile作为这些shell的配置文件。

需要牢记的一点是，默认的Debian框架目录（/etc/skel，用于存放要复制到新用户账户主目录的文件和目录）包含.profile文件，但不包含.bash_profile和.bash_login文件。此外Debian使用Bash作为默认的shell，因此，许多Debian用户习惯于将他们的Bash 登录shell设置放在.profile文件中.

> 不同的发行版可能有点不同,具体可以看其相关脚本的内容.

### non-login shell

1. `~/.bashrc`
2. /etc/bashrc(或/etc/bash.bashrc)
3. bash启动

## ~/.profile与~/.bashrc的区别:

这两者都具有个性化定制功能,但

- `~/.profile`可以设定本用户专有的路径，环境变量等，它只能登入的时候执行一次
- `~/.bashrc`也是某用户专有设定文档，可以设定路径，命令别名，每次shell script的执行都会使用它一次.

## X11

见`/etc/X11/Xsession`
