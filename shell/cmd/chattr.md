# chattr

## 描述

修改文件的扩展属性. 与chmod相比, chmod只能改变文件rwx权限, 更底层的属性是由chattr来修改的.

## 格式

    chattr [ options ] [ mode ] files...

## 选项

- -R : 递归修改
- -V : 显示命令的执行过程

mode:
- + : 增加参数
- - : 移除参数
- = : 更新为指定参数
- A : 不修改该文件的atime
- a : 只能向文件追加内容, 而不能删除. 多用于保障日志文件的安全, 比如`.bash_history`
- i : 设定文件不能被删除, 改名, 写入或新增内容, 比如`/etc/passwd`

## 例
```bash
# lsattr ad.js # 查看文件的扩展属性
--------------e---- ad.js
# chattr +a ad.js # 追加属性`a`
# lsattr ad.js
-----a--------e---- ad.js
# rm -f ad.js 
rm: 无法删除'ad.js': 不允许的操作
# echo 111 >> ad.js # 能追加
# echo 222 > ad.js # 不能清空
bash: ad.js: 不允许的操作
# chattr +i ad.js
# rm ad.js # 即使root也不能删除
rm: 无法删除'ad.js': 不允许的操作
```
