# extract(提取)
## bin
```bash
$ sudo apt install binwalk
$ binwalk -e xxx.bin # 已用ZStack-Cloud-installer-4.3.8.bin验证
```

`vim xxx.bin`可见其开头是"#!/bin/bash", 因此它就是一个bash script, 核心是从bin文件的指定行(by`tail`命令)开始提取内容, 再对提取的内容进行操作, 比如直接执行/解压再执行等等.