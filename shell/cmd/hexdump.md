# hexdump
以指定格式查看`二进制`文件的内容

## 选项:
- -C : 同时输出规范的十六进制和ASCII码
- -n <length> : 只格式化输入文件的前length个字节

## example
```
$ hexdump -C boot.bin
```