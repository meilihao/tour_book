# pwd (print working directory)

## 描述

显示当前所在目录

## 选项

- -P 显示物理路径. 比如当前目录是软链接的话会显示源文件的路径

## 例
```bash
chen go_path $ pwd
/home/chen/go_path
chen go_path $ ls -l /home/chen/go_path
lrwxrwxrwx 1 chen chen 22 1月  14  2017 /home/chen/go_path -> /home/chen/git/go/src/
chen go_path $ pwd -P
/home/chen/git/go/src
```