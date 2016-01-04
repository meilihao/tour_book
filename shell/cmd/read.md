# read

## 描述

read命令接收标准输入（键盘）或其他文件描述符的输入,从中读取一行.

## 例

    # read -n number_of_chars variable_name # 读取n个字符存入变量
    # read -s var # 以不在终端显示内容的方式来读取输入信息
    # read -p "Enter input:" var # 显示提示信息
    # read -t  timeout-seconds var # 读取指定时间内的输入
    # read -d ":" var # 设置界定符
