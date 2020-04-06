# gzip/gunzip
生成或解压后缀为".gz"的压缩文件

> gunzip是`gzip -d`的脚本

## gzip选项
- -d : 解压
- -r : 递归压缩
- -t : 检查压缩文件的完整性
- -v : 对每个压缩或解压缩的文档, 显示相应的文件名和压缩比
- -l : 显示压缩文件的压缩信息
- -<num> : 指定压缩比. 1~9, 默认是6

## example
```
# gzip -9v /opt/etc.zip # 对etc.zip进行gzip压缩, 生成etc.zip.gz
# gzip -l /opt/etc.zip.gz # 显示etc.zip.gz的压缩信息
# gzip -d /opt/etc.zip.gz # 对etc.zip.gz解压
# gunzip /opt/etc.zip.gz # 对etc.zip.gz解压
```