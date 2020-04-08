# cpio
通过重定向的方式将文件进行打包/备份/还原/恢复的工具, 也可解压以".cpio"或".tar"结尾的文件.

## 格式
```
cpio [option] > 文件或设备名
cpio [option] < 文件或设备名
```
> cpio在打包/备份时使用绝对/相对路径, 那么在还原时就使用相应的路径.

> cpio无法直接读取文件, 它需要每个文件或目录的完整路径才能读取, 而find的输出正好符合这点, 因此它一般与find联用.

##  选项
- -o : 将文件复制/打包成文件或者将文件输出到设备上
- -i : 将打包文件解压或将设备上的备份还原到系统中
- -t : 查看cpio打包的文件内容或者输出到设备上的文件内容
- -v : 显示打包过程中的文件名称 
- -d : 在cpio还原文件的过程中, 自动创建相应的目录 
- -c : 一种较新的存储方式
- -B : 让默认块可以增大到5120B, 默认是512B. 这样可加快存取速度
- -u : 无条件覆盖所有文件

## example
```
# find /etc -type f | cpio -ocvB > /opt/etc.cpio # 将/etc下的普通文件备份到etc.cpio
# find / -print | cpio -ocvB > /dev/st0 # 将/备份到scsi 磁带机上
# cpio -icdvt < /dev/st0 #  查看磁带机上的备份文件
# cpio -icdvt < /dev/st0  >/tmp/st_content #  查看磁带机上的备份文件, 因为内容太多, 无法全部显示, 因此将信息存入st_content
# cpio -icduv < /opt/etc.cpio # 将备份还原到相应的位置, 同名文件则覆盖
```