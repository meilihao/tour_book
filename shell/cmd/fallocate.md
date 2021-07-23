# fallocate
参考:
- [Linux文件空洞 与稀疏文件](https://lrita.github.io/images/posts/filesystem/Linux_File_Hole_And_Sparse_Files.pdf)

创建稀疏文件

## 选项

- -l 指定文件大小, 支持KB, MB, and so on for GB, TB, PB, EB, ZB, and YB.
- -o 设置偏移量

## example

    # fallocate -l 1G # 创建1G大小的稀疏文件