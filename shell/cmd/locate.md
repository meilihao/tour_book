# locate

## 描述

查找文件

> 用来查找符合条件的文件或目录，条件可以使用正则表达式。locate的查找速度要比find快得多，其方法是先建立一个包括系统内所有的文件、路径的数据库。在数据库中查找，而不需要实际进入文件系统查找文件,但locate在查找时，无法感知新增和已删除的文件.数据库可用updatedb命令更新,通常可由cron来定期更新.

whereis和locate检索的是同一数据库里的内容,只是参数有区别而已.

## 核心文件
- /etc/cron.daily/mlocate
- /etc/updatedb.conf
- /usr/bin/mlocate # locate
- /usr/bin/updatedb.mlocate # updatedb
- /var/lib/mlocate

> locate由alternatives管理.