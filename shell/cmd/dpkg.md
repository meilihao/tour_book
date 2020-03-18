# dpkg

## example
```
$ sudo dpkg -i --force-bad-verify  acl_2.2.52-3_amd64.deb # 跳过签名验证
$ dpkg -S file # 这个文件属于哪个已安装软件包
$ dpkg -L package # 列出软件包中的所有文件
$ dpkg -s package # 列出软件包中的描述信息
$ echo "PACKAGE hold" | sudo dpkg --set-selections  ##锁定软件包
$ dpkg --get-selections | grep hold  ##显示锁定的软件包列表
$ echo "PACKAGE install" | sudo dpkg --set-selections  ##解除对软件包的锁定
```

按文件搜索package也可直接使用[debian package服务](https://www.debian.org/distrib/packages)