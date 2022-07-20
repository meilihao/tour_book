# zstd

安装: `apt/yum/dnf install zstd`

## example
```bash
man zstd
zstd -z my_file.txt # 压缩
zstd my_file.txt
zstd -d my_file.txt.zst # 解压
unzstd my_file.txt.zst # `unzstd`=`zstd -d`
zstd --rm my_file.txt # 压缩并删除源文件
zstd --rm -d my_file.txt.zst # 解压缩并删除压缩文件
# 使用 tar与zstd压缩目录
tar -I zstd -cvf my_folder.tar.zst my_folder/

# 使用 tar 和 zstd 压缩多个文件
tar -I zstd -cvf my_files.tar.zst file1.txt file2.txt
# zstd -l my_file.txt.zst # 查看压缩文件的内容
# zstd -t my_file.txt.zst # 检查压缩文件是否正常. `-t`的效果等价于`--decompress 1>/dev/null`
# zstd -dcf my_file.txt.zst <=> zstdcat my_file.txt.zst # 解压缩并其输出内容
# zstd -9 my_file.txt # 只能压缩级别
# 指定并发压缩的线程数
zstd -T4 my_file.txt
# 指定并发解压的线程数
zstd -T4 -d my_file.txt.zst
zstdmt my_file.txt # zstdmt  <=> `zstd -T0`, `T0`表示线程数等于cpu核数
zstd -b5 # 测试指定压缩级别的效率
zstd --priority=rt my_file.txt # 将压缩进程设定为rt(实时)以加速压缩
zstd -v my_file.txt.zst # 输出detail
```

> go example: [github.com/klauspost/compress/zstd](https://github.com/TeaOSLab/EdgeNode/blob/d5f6acf6903ca42d1dc75004bb95b6462cb41033/internal/compressions/writer_zstd.go)