# rm (remove)

## 描述

删除文件或目录

## 语法格式

```
rm [OPTION]... FILE...
```

## 选项

- -d 删除可能仍有数据的目录
- -f 强制删除(会忽略不存在的文件且处理过程不提示任何信息)
- -i 交互模式(即删除确认)
- -r 递归删除(同时删除该目录下的所有子目录)
- -v 显示详细的处理信息

## example
```bash
$ rm -rf ./* # 加上" ./"来配合通配符, 避免错误执行`rm -rf *`
## Delete all file except file1 ##
rm  !(file1)
 
## Delete all file except file1 and file2 ##
rm  !(file1|file2) 
 
## Delete all file except all zip files ##
rm  !(*.zip)
 
## Delete all file except all zip and iso files ##
rm  !(*.zip|*.iso)
 
## You set full path too ##
rm /Users/vivek/!(*.zip|*.iso|*.mp3)
 
## Pass options ##
rm [options]  !(*.zip|*.iso)
rm -v  !(*.zip|*.iso)
rm -f  !(*.zip|*.iso)
rm -v -i  !(*.php)
# ls
'--exclude=lfs_root'   boot
# rm '--exclude=lfs_root'
rm: unrecognized option '--exclude=lfs_root'
Try 'rm ./'--exclude=lfs_root'' to remove the file '--exclude=lfs_root'.
Try 'rm --help' for more information.
# rm ./'--exclude=lfs_root' # ok
```