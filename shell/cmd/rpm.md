# rpm
已被dnf取代.

## 常用的 RPM 软件包命令
安装软件的命令格式 rpm -ivh filename.rpm 
升级软件的命令格式 rpm -Uvh filename.rpm 
卸载软件的命令格式 rpm -e filename.rpm 
查询软件描述信息的命令格式 rpm -qpi filename.rpm 
列出软件文件信息的命令格式 rpm -qpl filename.rpm 
查询文件属于哪个 RPM 的命令格式 rpm -qf filename

## 常见的 Yum 命令
yum repolist all 列出所有仓库
yum list all 列出仓库中所有软件包
yum info 软件包名称 查看软件包信息
yum install 软件包名称 安装软件包
yum reinstall 软件包名称 重新安装软件包
yum update 软件包名称 升级软件包
yum remove 软件包 移除软件包
yum clean all 清除所有仓库缓存
yum check-update 检查可更新的软件包
yum grouplist 查看系统中已经安装的软件包组
yum groupinstall 软件包组 安装指定的软件包组
yum groupremove 软件包组 移除指定的软件包组
yum groupinfo 软件包组 查询指定的软件包组信息
