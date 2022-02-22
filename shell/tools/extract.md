# extract(提取)
## bin
ref:
- [binwalk使用整理](https://blog.csdn.net/weixin_44932880/article/details/112478699)

```bash
$ sudo apt install binwalk
$ binwalk -e xxx.bin # 已用ZStack-Cloud-installer-4.3.8.bin验证
$ binwalk xxx.bin # 查看bin文件的布局
$ dd if=zstack-installer.bin of=i.sh bs=1 skip=0 count=312 # 提前bin开头的sh script部分
$ --- 手动安装(比apt install少很多依赖) for centos 7.6
$ wget https://github.com/ReFirmLabs/binwalk/archive/master.zip # binwalk由python3编写
$ unzip master.zip
$ cd binwalk-master && sudo python3 setup.py uninstall && sudo python3 setup.py install
```

`vim xxx.bin`可见其开头是"#!/bin/bash", 因此它就是一个bash script, 核心是从bin文件的指定行(by`tail`命令)开始提取内容, 再对提取的内容进行操作, 比如直接执行/解压再执行等等.

binwalk选项:
- -D, --dd=<type[:ext[:cmd]]>

	- type是签名描述中包含的小写字符串（支持正则表达式）
	- ext是保存数据磁盘时使用的文件扩展名（默认为none）
	- cmd是在将数据保存到磁盘后执行的可选命令


	`binwalk -D 'zip archive:zip:unzip %e' -D 'png image:png' firmware.bin`:
	1. 该选项将提取包含字符串“zip archive”,文件扩展名为“zip”的文件，然后执行“unzip”命令. 请注意使用’％e’占位符。执行unzip命令时，此占位符将替换为解压缩文件的相对路径
	1. 此外，PNG图像按原样提取，带有’png’文件扩展名。
