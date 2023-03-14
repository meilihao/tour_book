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

## nsis
ref:
- [Can I decompile an existing installer?](https://nsis.sourceforge.io/Can_I_decompile_an_existing_installer)

```bash
# zypper install p7zip-full
# 7z e winbareos-21.0.0-release-64-bit.exe # 解压出来没有目录层级
```

或使用[Universal Extractor 2](https://github.com/Bioruebe/UniExtract2)或[7-Zip](https://sourceforge.net/projects/sevenzip/files/7-Zip/15.05/)

`7-Zip 15.05`可提取安装文件的nsi配置.

## FAQ
### 制作bin
linux 下制作二进制`.bin`的文件的方法是使用cat 命令将执行脚本和打包文件同时放到一个文件里, 在给它可执行权限. 这样安装的时候只要使用一个包, 直接执行该包即可安装完毕, 简单方便.

准备demo:
```bash
# tree demo
demo
├── 新建文件.txt
└── t.sh
# cat t.sh
#!/bin/bash
demo_pkg="`ls demo*.rpm 2>/dev/null |grep -v demo_ui 2>/dev/null |head -1`" # 变量名不能出现`-`
if [ -f ${demo_pkg} ]; then
	echo "found "${demo_pkg}
fi
# tar -cvf demo.tar demo
```

#### [zstack](https://github.com/zstackio/zstack-utility/blob/master/zstackbuild)
```bash
$ cat #!/bin/bash
cat >$1 <<EOF
#!/bin/bash
#set -x
line=\`wc -l \$0|awk '{print \$1}'\`
line=\`expr \$line - 12\` # 12是生成的setup.sh的行数, 需去掉第一行和最后一行的换行
tmpdir=\`mktemp\`
/bin/rm -f \$tmpdir
mkdir -p \$tmpdir
tail -n \$line \$0 |tar -x -C \$tmpdir --strip-components 1
cd \$tmpdir
bash ./t.sh \$*
ret=\$?
#rm -rf \$tmpdir
exit \$ret
EOF
# cat build_installation_bin.sh
#!/bin/bash
cat $1 $2 > $3
chmod a+x $3
# ./gen_setup.sh setup.sh
# ./build_installation_bin.sh setup.sh demo.tar upgrade.bin
# ./upgrade.bin
# strings upgrade.bin # 可用strings验证
# tail -n 3 ./upgrade.bin # 可验证获取的demo.tar是否正确
```

#### 网上
ref:
- [linux下制作二进制bin 文件制做方法](https://blog.csdn.net/xiaotengyi2012/article/details/8493929)

```bash
# cat setup2.sh
#!/bin/bash
#set -x
tmpdir=`mktemp --suffix=xxx`
/bin/rm -f $tmpdir
mkdir -p $tmpdir
sed -n -e '1,/^exit 0$/!p' $0 |tar -x -C $tmpdir --strip-components 1
cd $tmpdir
bash ./t.sh $*
rm -rf $tmpdir
exit 0

# ./build_installation_bin.sh setup2.sh demo.tar upgrade.bin
# ./upgrade.bin
```

**setup2.sh必须存在最后一行空行, 便于sed提取, 否则sed会匹配不到需要截断的行, 因此没有最后一行空行时, sed会匹配到整个upgrade.bin**

> 也遇到过setup2.sh最后一行空行不能有换行, 当时是直接cat了rpm包(没用tar封装), 否则sed提取结果开头多了一个换行.