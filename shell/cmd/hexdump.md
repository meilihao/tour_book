# hexdump
以指定格式查看`二进制`文件的内容

> `*`表示该区间输出同上一行

## 选项:
- -C : 同时输出规范的十六进制和ASCII码
- -n <length> : 只格式化输入文件的前length个字节
- -s <length> : 跳过length个字节

## example
```
$ hexdump -C boot.bin
$ hexdump -C /dev/sda1 -n 1024 -s 1024 # 查看ext4 block group 0 的super block信息
```