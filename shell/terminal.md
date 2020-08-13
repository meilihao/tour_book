# 终端处理工具

tput和stty

## tput

    # tput cols # 终端的行数
    # tput lines # 终端的列数
    # tput longname # 终端的名称
    # tput cup 9 10 # 移动光标到（第9行，第10列）
    # tpu serf no # 设置文本前景色（no：0~7）
    # tpu setb no # 设置终端背景色
    # tput bold # 设置文本为粗体

## stty

    # stty -echo # 禁止将输出内容发送给终端，常用于输入密码但不显示内容.
    # stty echo # 允许将输出内容发送给终端

## 快捷键
- `ctrl + a` : 光标回到命令行首
- `ctrl + e` : 光标回到命令行尾

## 颜色
参考:
- [Bash tips: Colors and formatting (ANSI/VT100 Control sequences)](https://misc.flogisoft.com/bash/tip_colors_and_formatting)

```bash
$ echo -e "\e[1;31m This is red text \e[0m" # 8/16 Colors
```

`\e[1;<color>m`开始设置文本颜色, `\e[0m` 重置颜色.

在bash中，escape字符可以是以下任一项开头：
- `\e`
- `\033`(八进制)
- `\x1B`(十六进制)

颜色区分:
- 8-bit color 又名 256 color
- 24-bit color 又名 true color，一共有 16,777,216 colors
- 32-bit color 基于 24-bit color 而生，增加了 8-bit 透明通道. 现代虚拟终端 deepin-terminal 就支持, 以致它的背景可以变成半透明.

`TERM`: 环境变量TERM就是用来告知程序当前虚拟终端支持的颜色.

## 终端信息采集
```bash
# tput cols/lines # 获取终端的行数和列数
# tput longname # 答应你 # 查看终端名称
# tput setb <n> # 设置终端背景色(n是0~7)
# tput setf <n> # 设置终端前景色(n是0~7)
```

`stty -echo`可禁止将输出发送到终端, 在输入密码场景可禁止显示输入的密码内容；`stty echo`则允许发送输出到终端, 适用于恢复输出的场景, 比如结束密码输出后.