## systemd 文件类型及存放路径

systemd配置文件被称为unit单元，根据类型不同，以不同的扩展名结尾:

- `.service` 系统服务；
- `.target` 一组系统服务；
- `.automount` 自动挂载点；
- `.device` 能被内核识别的设备；
- `.mount` 挂载点；
- `.path` 文件系统的文件或者目录；
- `.scope` 外部创建的进程；
- `.slice` 一组分层次管理的系统进程；
- `.snapshot` 系统服务状态管理；
- `.socket` 进程间通讯套接字；
- `.swap` 定义swap文件或者设备；
- `.timer`	定义定时器。

systemd单元文件放置位置:

- `/usr/lib/systemd/system/` systemd默认单元文件安装目录
- `/run/systemd/system systemd` systemd单元运行时创建，这个目录优先于安装目录
- `/etc/systemd/system` 系统管理员创建和管理的单元目录，优先级最高