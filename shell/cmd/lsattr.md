# lsattr

## 描述

查看文件的扩展属性. 与chmod相比, chmod只能改变文件rwx权限, 更底层的属性是由lsattr来修改的.

## 格式

    lsattr [ options ] files...

## 选项

- -R : 递归查看
- -a : 显示所有文件包括隐藏文件的扩展属性
- -d : 显示目录的扩展属性