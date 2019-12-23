# dpkg

## example
```
$ sudo dpkg -i --force-bad-verify  acl_2.2.52-3_amd64.deb # 跳过签名验证
$ dpkg -S file # 这个文件属于哪个已安装软件包
$ dpkg -L package # 列出软件包中的所有文件
```