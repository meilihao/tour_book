# tar

[参考](https://linux.cn/article-7802-1.html)

## 描述

解压缩

## 选项

- -C : 执行归档动作前变更工作目录到指定路径
- -c : 建立压缩档案
- -f : 指定要操作的归档文件名
- -p : 保留原文件的访问权限
- -x : 解压
- -t : 查看内容
- -r : 向归档文件末尾追加文件
- -u : 更新原压缩包中的文件

上面五个是独立的命令，压缩解压都要用到其中一个，可以和别的命令连用但只能用其中一个.

- -z : gzip属性的(即tar.gz)
- -j : bz2属性的(即tar.bz2)
- -J : xz属性
- -Z : compress属性的(tar.Z)
- -v : 显示所有过程
- -O : 将文件解开到标准输出

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

### 压缩

```
tar –cvf jpg.tar *.jpg //将目录里所有jpg文件打包成tar.jpg
tar –czf jpg.tar.gz *.jpg   //将目录里所有jpg文件打包成jpg.tar后，并且将其用gzip压缩，生成一个gzip压缩过的包，命名为jpg.tar.gz
tar –cjf jpg.tar.bz2 *.jpg //将目录里所有jpg文件打包成jpg.tar后，并且将其用bzip2压缩，生成一个bzip2压缩过的包，命名为jpg.tar.bz2
tar –cZf jpg.tar.Z *.jpg   //将目录里所有jpg文件打包成jpg.tar后，并且将其用compress压缩，生成一个umcompress压缩过的包，命名为jpg.tar.Z
rar a jpg.rar *.jpg //rar格式的压缩，需要先下载rar for linux
zip jpg.zip *.jpg //zip格式的压缩，需要先下载zip for linux
tar -rvf data.tar /etc/fstab //在压缩过的 tar 文件中无法进行追加文件操作
tar -zcvf optbackup-$(date +%Y-%m-%d).tgz /opt/ //使用 tar 命令进行定时备份
split -b <Size-in-MB> <tar-file-name>.<extension> “prefix-name”//分割体积庞大的 tar 文件为多份小文件
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
tar -zcpvf myarchive.tgz /etc/ /opt/ --exclude=*.html //排除指定文件或类型后创建 tar 文件
tar -tvf myarchive.tar.gz  | more //列出 .tar.gz 文件中的内容
tar czf xx.tgz -C /xxx/xxx A //使用-C指定相对路径
```
