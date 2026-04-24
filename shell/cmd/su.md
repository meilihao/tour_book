# su

## 选项
- - : 加载相应用户的环境变量
- -c : 改变身份运行一个指令后就结束
- -m : 改变用户, 但不加载环境变量

## example
```
$ sudo su - freeswitch  -s /bin/bash -c "/home/chen/node_server" # 以freeswitch运行node_server(su xxx不能切换用户时使用), 要使用绝对路径.
```

## FAQ
## su root报错`无法执行 fish: 没有那个文件或目录`
```sh
$ cat /etc/passwd|grep root
root:x:0:0:root:/root:fish # root的shell的路径不对
```