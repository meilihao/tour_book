# rsync
参考:
- [rsync命令](https://man.linuxde.net/rsync)

rsync命令是一个远程数据同步工具，可通过LAN/WAN快速同步多台主机间的文件. rsync使用所谓的`rsync算法`来使本地和远程两个主机之间的文件达到同步，这个算法只传送两个文件的不同部分，而不是每次都整份传送，因此速度相当快.

## 选项
- -a, --archive : 归档模式，表示以递归方式传输文件，并保持所有文件属性，等于-rlptgoD
- -c, --checksum : 打开校验开关，强制对文件传输进行校验
- -e "ssh -i $HOME/.ssh/somekey" : `-e`表示指定使用rsh、ssh方式进行数据同步, 这里同时指定了使用的ssh key
- --exclude : 排除路径, 相对于源地址
- -n : dry-run
- --progress : 显示备份过程
- -P 等同于 --partial : 断点续传
- --update : 跳过目标文件比源文件新的文件
- -v : 显示执行细节
- -z, --compress : 对备份的文件在传输时进行压缩处理

## 格式
```
rsync [OPTION]... SRC DEST
rsync [OPTION]... SRC [USER@]host:DEST
rsync [OPTION]... [USER@]HOST:SRC DEST
rsync [OPTION]... [USER@]HOST::SRC DEST
rsync [OPTION]... SRC [USER@]HOST::DEST
rsync [OPTION]... rsync://[USER@]HOST[:PORT]/SRC [DEST]
```

 对应于以上六种命令格式，rsync有六种不同的工作模式：
- 拷贝本地文件. 当SRC和DES路径信息都不包含有单个冒号":"分隔符时就启动这种工作模式。如：rsync -a /data /backup
- 使用一个远程shell程序(如rsh、ssh)来实现将本地机器的内容拷贝到远程机器。当DST路径地址包含单个冒号":"分隔符时启动该模式。如：rsync -avz *.c foo:src
- 使用一个远程shell程序(如rsh、ssh)来实现将远程机器的内容拷贝到本地机器。当SRC地址路径包含单个冒号":"分隔符时启动该模式。如：rsync -avz foo:src/bar /data
- 从远程rsync服务器中拷贝文件到本地机。当SRC路径信息包含"::"分隔符时启动该模式。如：rsync -av root@192.168.78.192::www /databack
- 从本地机器拷贝文件到远程rsync服务器中。当DST路径信息包含"::"分隔符时启动该模式。如：rsync -av /databack root@192.168.78.192::www
- 列远程机的文件列表。这类似于rsync传输，不过只要在命令中省略掉本地机信息即可。如：rsync -v rsync://192.168.78.192/www

## example
```bash
$ rsync -avc --dry-run --update ./* root@192.168.0.137:/opt/test # 仅计算同步
$ rsync -avc --update --exclude="adapter" ./* root@192.168.0.137:/opt/test # 会排除./adapter
$ rsync -P --rsh=ssh aliyun:~/git/lfs.img.zstd . # 断点续传
$ rsync -ah --progress source destination # 拷贝带进度
$ sspass -p "<password>" rsync -ah --progress source destination # 传入ssh password
$ rsync -av -e ssh --exclude='*.new' ~/virt/ root@centos7:/tmp
```