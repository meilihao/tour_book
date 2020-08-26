# zip/unzip
生成或解压后缀为".zip"的压缩文件

## zip

### 选项
- -r : 递归压缩
- -d : 从压缩文件内删除指定的文件
- -i : "文件列表" : 只压缩文件列表中的文件
- -x : "文件列表" : 压缩时排除文件列表中的指定文件, 文件用法应是`-x "xxx.txt"`, 目录则是`-x "xxx/*"` 
- -u : 更新文件到压缩文件中
- -m : 将文件加入压缩文件压缩后, 删除原文件. 即把文件移入压缩文件
- -F : 尝试修复损坏的压缩文件
- -T : 检查压缩文件内的每个文件是否正确
- <压缩级别> : 1~9

### example
```
# zip -9r /opt/etc.zip /etc # 将/etc进行压缩保存入/opt/etc.zip.
# zip -r /opt/var.zip -/var -x "*.log" # 将/var中除"*.log"外的所有文件压入/opt/var.zip
# zip /opt/etc.zip -d etc/passwd # 将etc.zip中的etc/passwd删除
# zip -u /opt/etc.zip /etc/inittab # 将/etc/inittab更新入etc.zip
# zip -u /opt/etc.zip /etc/* # 将/etc下所有内容更新入etc.zip
# zip -ru /opt/etc.zip /etc # 将/etc下所有内容更新入etc.zip
# zip -q -r -P 123456 zipfile.zip /etc
```

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
```