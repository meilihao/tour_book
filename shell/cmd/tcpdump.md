# tcpdump

## 描述

根据使用者的定义对网络上的数据包进行截获的包分析工具,常与wireshark组合使用(tcpdump捕获数据,再用wireshark分析)

## 选项
- -i : 指定监听的网络接口
- -w : 将捕获的数据写入文件

## 例子
    # tcpdump -i lo host 127.0.0.1 and port 8210  -w out.cap # 捕获发送到127.0.0.1:8210的数据
