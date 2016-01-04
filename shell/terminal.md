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
