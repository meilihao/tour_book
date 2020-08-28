# tar

[参考](https://linux.cn/article-7802-1.html)

## 描述

归档工具

## 选项

- -C : 执行归档动作前变更工作目录到指定路径
- -c : 建立档案
- -f : 指定要操作的归档文件名
- -h : tar默认保留软连接, 使用该参数后打包被指向的文件(文件名与软链名相同)
- -p : 保留原文件的访问权限
- -P : 使用绝对路径来压缩
- -x : 还原归档文件
- -t : 查看内容即文件列表
- -r : 向归档文件末尾追加文件
- -u : 更新归档中的文件
- --exclude : 排除目录或文件, 参数必须紧跟在tar后, 如果放在被压缩路径后可能会导致循环压缩
- --strip-components N : 解压时跳过目录的前N层前缀

上面五个是独立的命令，压缩解压都要用到其中一个，可以和别的命令连用但只能用其中一个.

- -z : gzip属性的(即tar.gz)
- -j : bz2属性的(即tar.bz2)
- -J : xz属性
- -Z : compress属性的(tar.Z)
- -v : 显示所有过程
- -O : 将文件解开到标准输出
- -f : 归档位置, **`-f`*后不能有其他选项
- -w : 在还原归档时, 把所有文件的修改时间设置现在时间
- -p : 归档时, 保持文件属性不变
- -N "yyyy/mm/dd" : 在指定日期之后的文件才会打包到归档中

上面的参数是根据需要在压缩或解压档案时可选的.

总结

1. *.tar 用 tar –xvf 解压
2. *.gz 用 gzip -d或者gunzip 解压
3. *.tar.gz和*.tgz 用 tar –xzf 解压
4. *.bz2 用 bzip2 -d或者用bunzip2 解压
5. *.tar.bz2用tar –xjf 解压
6. *.Z 用 uncompress 解压
7. *.tar.Z 用tar –xZf 解压
8. *.rar 用 unrar e解压
9. *.zip 用 unzip 解压

## 例

### 打包/压缩

```
tar -zcvf - /etc |tar -zxvf - # 第一个"-"表示输出到stdout, 第二个"-"是将管道传入的信息作为解压的数据来源
tar -N "2008/7/21" -zcvf log.tar.gz /var/log # 压缩/var/log中2008/7/71以后的文件
tar -ztvf /opt/etc.tar.gz # 查看内容
tar –cvf jpg.tar *.jpg //将目录里所有jpg文件打包成tar.jpg
tar –czf jpg.tar.gz *.jpg   //将目录里所有jpg文件打包成jpg.tar后，并且将其用gzip压缩，生成一个gzip压缩过的包，命名为jpg.tar.gz
tar –cjf jpg.tar.bz2 *.jpg //将目录里所有jpg文件打包成jpg.tar后，并且将其用bzip2压缩，生成一个bzip2压缩过的包，命名为jpg.tar.bz2
tar –cZf jpg.tar.Z *.jpg   //将目录里所有jpg文件打包成jpg.tar后，并且将其用compress压缩，生成一个umcompress压缩过的包，命名为jpg.tar.Z
tar --exclude=tomcat/logs --exclude=tomcat/libs --exclude=tomcat/xiaoshan.txt -zcvf tomcat.tar.gz tomcat # 排除logs和libs两个目录及文件xiaoshan.txt
rar a jpg.rar *.jpg //rar格式的压缩，需要先下载rar for linux
zip jpg.zip *.jpg //zip格式的压缩，需要先下载zip for linux
tar -rvf data.tar /etc/fstab //在压缩过的 tar 文件中无法进行追加文件操作
tar -zcvpf optbackup-$(date +%Y-%m-%d).tgz /opt/ //使用 tar 命令进行定时备份
split -b <Size-in-MB> <tar-file-name>.<extension> “prefix-name”//分割体积庞大的 tar 文件为多份小文件
tar --exclude=${LFSRoot} -cJpf ${LFSRoot}/iso/lfs-temp-tools-10.0-systemd-rc1.tar.xz . # tar.xz
```

### 解压

```
tar –xvf file.tar //解压 tar包
tar -xzvf file.tar.gz //解压tar.gz
tar -xjvf file.tar.bz2   //解压 tar.bz2
tar –xZvf file.tar.Z   //解压tar.Z
unrar e file.rar //解压rar
unzip file.zip //解压zip
tar -xvf myarchive.tar -C /tmp/ //释放 tar 文件到指定目录
tar -xvf myarchive.tar root/anaconda-ks.cfg -C /tmp/
root/anaconda-ks.cfg //释放 tar 文件中的指定文件或目录
tar -zcpvf myarchive.tar.gz /etc/ /opt/
tar --exclude=*.html -zcpvf myarchive.tgz /etc/ /opt/ //排除指定文件或类型后创建 tar 文件
tar -tvf myarchive.tar.gz  | more //列出 .tar.gz 文件中的内容
tar -tvf /lfs/sources/mpfr-*.tar.xz // 列出.tar.xz内容
tar czf xx.tgz -C /xxx/xxx A //使用-C指定相对路径
tar -xf binutils-2.35.tar.xz -C a --strip-components 1 # 解压时生成的路径不包括父目录`binutils-2.35`
```

## FAQ
### tar: 由于前次错误，将以上次的错误状态退出
解压时使用参数`-C`指定解压目的即可. 但有时还是不行, 换用zip压缩解压缩即可.

```bash
$ zip -9r lfs/lfs_root/iso/lfs-fsroot.zip lfs -x="lfs/lfs_root/*"
$ unzip lfs-fsroot.zip
```