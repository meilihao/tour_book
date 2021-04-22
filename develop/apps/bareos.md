# bareos
Bareos 由 bacula fork而來.

[Bareos组成](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-components-or-services):
- bconsole : 全功能cli, 与Director进行通信
- webui : 只用于备份和恢复功能
- Director : Bareos中控, 它计划并监督所有备份, 还原, 验证和存档操作.
- Storage Daemon : 在Bareos上作为备份数据的存储空间, 允许多个

    应Bareos Director的请求负责从Bareos File后台驻留程序接受数据，并将文件属性和数据存储到物理备份介质或卷中
- File Daemon : **在客户机**, 管理本地文件的备份和恢复.

    它会应Director的请求它会找到要备份的文件，并将指定数据发送到Bareos Storage Daemon
- Catalog : 目录服务由负责维护所有备份文件的文件索引和卷数据库.

    目录服务允许系统管理员或用户快速查找和还原任何所需的文件.

> Bareos推荐使用postgres, mysql/mariadb已废弃.

> 要成功执行保存或还原，必须配置并运行以下四个守护程序：Director daemon, File daemon, Storage daemon 以及 Catalog service(即DB).

> [Bareos所有相关package](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-packages)


## [部署](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#install-the-bareos-software-packages)
```bash
# -- pg
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
# -- bareos
wget -q http://download.bareos.org/bareos/release/20/xUbuntu_20.04/Release.key -O- | apt-key add -
wget -O /etc/apt/sources.list.d/bareos.list http://download.bareos.org/bareos/release/20/xUbuntu_20.04/bareos.list
# -- install
sudo apt install postgresql-12 postgresql-client-12 pgadmin4

vim ${pg}/pg_hba.conf
local bareos bareos md5 # bareos默认使用本地pg, 因此添加该匹配规则

systemctl restart postgres

apt install bareos bareos-database-postgresql # 输入db密码. bareos-database-postgresql会利用dbconfig-common mechanism, 在apt install过程中配置db, db配置在`/etc/dbconfig-common/bareos-database-common.conf`. 可用`dpkg-reconfigure bareos-database-common`重新配置, 手动配置db见[这里](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms)

systemctl restart bareos-dir
systemctl restart bareos-sd
systemctl restart bareos-fd

bareos-dir -t -f -d 500 -v # 测试bareos-dir是否正常, 包括与pg的连接
bareos-dbcheck -B # 作用同上, 显示db的连接信息

apt install bareos-webui # 基于php+apache2
systemctl restart bareos-dir

# -- 配置webui, 也可使用[bconsole configure子命令](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareosWebui.html#create-a-restricted-consoles)
cp /etc/bareos/bareos-dir.d/console/admin.conf.example vim /etc/bareos/bareos-dir.d/console/admin.conf && chown bareos:bareos /etc/bareos/bareos-dir.d/console/admin.conf
vim /etc/bareos/bareos-dir.d/console/admin.conf # 设置bareos-dir admin用于bareos-webui
systemctl restart bareos-dir # 不能省略, 否则可能webui无法登入(账号正确)
systemctl restart apache2 # 访问http://HOSTNAME/bareos-webui即可使用webui
```