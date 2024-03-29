# zip/unzip
生成或解压后缀为".zip"的压缩文件

## zip
zip默认输出包含两个不同的词，即 deflate 和 store, zip 默认使用的压缩方法是 deflate.

zip 的加密算法是使用流式加密的 PKZIP。而它很容易被破解。如果你想保护你的文件，请使用 7-Zip 或其他高级工具.
### 选项
- -r : 递归压缩
- -d : 从压缩文件内删除指定的文件
- -e : 对压缩文件进行密码保护
- -i : "文件列表" : 只压缩文件列表中的文件
- -j : 不保留所有需要压缩文件的绝对路径
- -q : quiet
- -x : "文件列表" : 压缩时排除文件列表中的指定文件, 需要放在其他参数的最后面, 文件用法应是`-x "xxx.txt"`, 目录则是`-x "xxx/*"`. 使用相对路径是要注意, 它相对的是当前工作路径而不是被打包目录的路径.
- -u : 更新文件到压缩文件中
- -m : 将文件加入压缩文件压缩后, 删除原文件. 即把文件移入压缩文件
- -F : 尝试修复损坏的压缩文件
- -s : 分割文件
- -T <size>: 检查压缩文件内的每个文件是否正确, 比如`1g`
- <压缩级别> : 1~9, 9最高, 默认是6
- -y : 保留符号链接

### example
```
# zip -9r /opt/etc.zip /etc # 将/etc进行压缩保存入/opt/etc.zip.
# zip -r /opt/var.zip /var -x "*.log" # 将/var中除"*.log"外的所有文件压入/opt/var.zip
# zip /opt/etc.zip -d etc/passwd # 将etc.zip中的etc/passwd删除
# zip -u /opt/etc.zip /etc/inittab # 将/etc/inittab更新入etc.zip
# zip -u /opt/etc.zip /etc/* # 将/etc下所有内容更新入etc.zip
# zip -ru /opt/etc.zip /etc # 将/etc下所有内容更新入etc.zip
# zip -q -r -P 123456 zipfile.zip /etc
# cd /home/chen/1/2021 && zip -r xx.zip img 01.pdf # 相对路径压缩
```

> 含空格的path可用`''`处理

## unzip
### 选项
- -x "文件列表" : 解压文件时排除文件列表中的文件
- -t : 测试压缩文件是否损坏, 但不解压
- -v : 查看压缩文件的详细信息: 文件, 文件大小, 压缩比等, 但不解压
- -n : 解压时不覆盖已存在的文件
- -o : 解压时覆盖已存在的文件, 并且不需要用户确认
- -d <目录> : 把压缩文件解压到指定目录下
- -r 递归, 将指定目录下的所有文件和子目录一起处理
- -S : 包含系统和隐藏文件
- -y : 直接保存符合连接, 而不是该连接所指向的文件
- -P : 指定密码
- -q : 不显示命令的执行过程

### example
```
# unzip -o /opt/etc.zip -x etc/inittab -d /etc # 将etc.zip解压到/etc下, 但etc/inittab除外, 且解压时同名覆盖
$ zip -9r lfs/lfs_root/iso/lfs-fsroot.zip lfs -x="lfs/lfs_root/*"
$ unzip lfs-fsroot.zip
# unzip tilix.zip -d / # 绿色安装tilix
# zip -9r - main.go |cat - > a.zip # 默认标准输出是打包进度
```

## FAQ
### windows zip解压乱码
Windows下生成的zip文件中的编码是GBK/GB2312等, 而Linux下的默认编码是UTF8.

```bash
$ export UNZIP="-O GBK"
$ export ZIPINFO="-O GBK"
$ unzip xxx.zip
```

推荐: 将env放入`~/.bashrc`