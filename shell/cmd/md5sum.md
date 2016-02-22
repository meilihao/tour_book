# md5sum

## 描述

计算md5值,类似的有sha1sum.

## 选项

- -c：从指定文件中读取MD5校验和，并进行校验.

## 例

    # md5sum filename > checksum.md5
    # md5sum file1 file2 # 输出每个文件的校验和
    [checksum1] file1
    [checksum2] file2
    # md5sum -c checksum.md5

### md5deep : 对目录进行校验
和md5deep类似的有sha1deep.

    $ md5deep -rl dirname > dirname.md5
    # -r : 使用递归方式
    # -l : 使用相对路径(默认情况下,md5deep输出文件的绝对路径).
    $ md5sum -c dirname.md5
