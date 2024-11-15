# iperf3
## 选项
1. 通用参数:
-p        端口号
-f        指定带宽输出格式： Kbits、Mbits、Gbits、Tbits
-i        监控报告时间间隔，单位秒(s)
-J        Json格式输出结果
--logfile 将结果输出到指定文件中

2. 服务端参数:
-s  以服务器模式运行
-D  后台运行服务器模式
 
3. 客户端参数:
-c  以客户端模式运行，连接到服务端
-F	传输或接收特定的文件
-t  传输时间，默认10秒
-n  传输内容大小，不能与-t同时使用
-b  目标比特率(0表示无限)(UDP默认1Mbit/sec，TCP不受限制)
-l  要读取或写入的缓冲区长度(TCP默认128 KB，UDP默认1460)
-O  忽略前几秒
-P	客户端到服务器的连接数，默认值为1
-R  反向模式运行，即服务端发送，客户端接收
-u  使用UDP协议，默认使用TCP协议
--get-server-output      输出服务端的结果
-w	设置套接字缓冲区为指定大小，对于TCP方式，此设置为TCP窗口大小，对于UDP方式，此设置为接受UDP数据包的缓冲区大小，限制可以接受数据包的最大值.

## example
```bash
# --- server
iperf3 -s -p 5201
# --- client
iperf3 -c x.x.x.x -p 5201 -t 5 -P 1 -R  # -c, 指定测速服务器IPx.x.x.x，-p, 指定端口为5201，-t测速时间5s，-P指定发送连接数1，-R表示下载测速  
```