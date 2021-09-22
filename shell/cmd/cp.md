# cp

## 描述

用来将一个或多个源文件或者目录复制到指定的目的文件或目录

## 格式

    cp [OPTION]... SOURCE... DIRECTORY

## 选项

- -a : 等同`-pdr`
- -d : 如果复制的源文件是符号链接, 那么仅复制符号链接本身
- -i : 覆盖已有文件前需用户确认
- -p : 复制时保持源文件的所有者, 权限及时间戳
- -r : 递归处理

## 例
```sh
$cp -r source_dir  dest_dir # 复制目录
$ sudo cp -rpv /nvim-linux64/*  / # 复制安装nvim, 有问题:`无法以目录'./nvim-linux64/bin' 来覆盖非目录'/bin'`, 可先处理`./nvim-linux64/bin/nvim`, 再用cp复制剩余内容.
$ sudo cp -r usr/* /usr # 同上, 实现绿色安装软件
$ cp –r test/ newtest  # 将当前目录 test/ 下的所有文件复制到新目录 newtest 
```
