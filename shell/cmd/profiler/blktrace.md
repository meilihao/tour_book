# blktrace
参考:
- [IO神器blktrace使用介绍](https://developer.aliyun.com/article/698568)
- [利用blktrace分析IO性能](http://linuxperf.com/?p=161)

## example
```bash
# blktrace -d /dev/sda -o - |blkparse -i -
```

## FAQ
### blkparse输出
```bash
# blktrace -d /dev/sda -o - |blkparse -i -
  8,0    3        1     0.000000000   594  A FWFSM 43781957 + 17 <- (8,4) 18782021
  8,0    3        2     0.000001967   594  Q FWFSM 43781957 + 17 [(null)]
  8,0    3        3     0.000006490   594  G FWFSM 43781957 + 17 [(null)]
  8,0    3        4     0.000031145   195  D  FN [kworker/3:1H]
  8,0    2        7     0.001628764     0  C  FN 0 [0]
  8,0    2        8     0.001671798   323  D WSM 43781957 + 17 [kworker/2:1H]
  8,0    2        9     0.001843458     0  C WSM 43781957 + 17 [0]
  8,0    2       10     0.001863050   323  D  FN [kworker/2:1H]
  8,0    2       11     0.002735974     0  C  FN 0 [0]
  8,0    2       12     0.002742812     0  C WSM 43781957 [0]
  8,0    3        5     1.477349073   233  A  WM 26379008 + 32 <- (8,4) 1379072
  8,0    3        6     1.477351508   233  Q  WM 26379008 + 32 [xfsaild/sda4]
  8,0    3        7     1.477362690   233  G  WM 26379008 + 32 [xfsaild/sda4]
  8,0    3        8     1.477364537   233  P   N [xfsaild/sda4]
  8,0    3        9     1.477375143   233  A  WM 26713248 + 32 <- (8,4) 1713312
  8,0    3       10     1.477376067   233  Q  WM 26713248 + 32 [xfsaild/sda4]
  8,0    3       11     1.477379594   233  G  WM 26713248 + 32 [xfsaild/sda4]
  8,0    3       12     1.477389324   233  A  WM 26713792 + 8 <- (8,4) 1713856
  8,0    3       13     1.477390200   233  Q  WM 26713792 + 8 [xfsaild/sda4]
  8,0    3       14     1.477393235   233  G  WM 26713792 + 8 [xfsaild/sda4]
...
```


1. 第一个字段：8,0 这个字段是设备号 major device ID和minor device ID。
1. 第二个字段：3 表示CPU
1. 第三个字段：11 序列号
1. 第四个字段：0.009507758 Time Stamp是时间偏移
1. 第五个字段：PID 本次IO对应的进程ID
1. 第六个字段：Event，这个字段非常重要，反映了IO进行到了那一步
1. 第七个字段：R表示 Read， W是Write，D表示block，B表示Barrier Operation
1. 第八个字段：223490+56，表示的是起始block number 和 number of blocks，即我们常说的Offset 和 Size
1. 第九个字段： 进程名

其中第六个字段非常有用：每一个字母都代表了IO请求所经历的某个阶段:
Q – 即将生成IO请求
|
G – IO请求生成
|
I – IO请求进入IO Scheduler队列
|
D – IO请求进入driver
|
C – IO请求执行完毕