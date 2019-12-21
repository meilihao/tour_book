# seq
产生从某个数到另外一个数之间的所有整数

## example
```sh
$ seq 40 |xargs -i dd if=/dev/zero 0f={}.mp4 bs=1001001024 count=1 # 批量随机生成 40个1G的大mp4文件
$ seq 18 |xargs -i cp /root/deepin-15.11-amd64.iso ./deepin-15.11-amd64_{}.iso # 造数据, dd+/dev/zero生成的数据与实际所需大小不一致
```