# extract(提取)
## bin
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
