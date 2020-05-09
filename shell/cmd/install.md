# install
作用是安装或升级软件或备份数据，它的使用权限是所有用户. install命令和cp命令类似，都可以将文件/目录拷贝到指定的地点. 但是，**install允许控制目标文件的属性, install通常用于程序的makefile**，使用它来将程序拷贝到目标（安装）目录.

## 选项
- -d，--directory：所有参数都作为目录处理(支持多个用空格分隔)，而且会创建指定目录的所有目录, 类似于`mkdir -p`
- -m，--mode=模式：自行设定权限模式 (像chmod 0700)，而不是rwxr-xr-x

## example
```bash
$ install a/e c # 类似 cp a/e c    # 注意c必须是文件
$ install -D x a/b/c # 类似 mkdir -p a/b && cp x a/b/c
$ install a/* d # 复制多个SOURCE文件到目的目录
```