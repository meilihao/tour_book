# bareos
Bareos 由 bacula fork而來.

[Bareos组成](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-components-or-services):
- bconsole : 全功能cli, 与Director进行通信
- webui : 只用于备份和恢复功能, 同时支持基于Web的bconsole界面
- Director : Bareos中控, 它计划并监督所有备份, 还原, 验证和存档操作.
- Storage Daemon : 在Bareos上作为备份数据的存储空间, 允许多个

    应Bareos Director的请求负责从Bareos File后台驻留程序接受数据, 并将文件属性和数据存储到物理备份介质或卷中
- File Daemon : **在客户机**, 管理本地文件的备份和恢复.

    它会应Director的请求它会找到要备份的文件, 并将指定数据发送到Bareos Storage Daemon
- Catalog : 目录服务由负责维护所有备份文件的文件索引和卷数据库.

    目录服务允许系统管理员或用户快速查找和还原任何所需的文件.

> Bareos推荐使用postgres, mysql/mariadb已废弃.

> 要成功执行保存或还原, 必须配置并运行以下四个守护程序：Director daemon, File daemon, Storage daemon 以及 Catalog service(即DB).

> [Bareos所有相关package](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-packages)

> [bareos网络连接概览](https://docs.bareos.org/TasksAndConcepts/NetworkSetup.html#network-connections-overview)

## 编译
env: Ubuntu 20.04

```bash
apt install libreadline-dev libpq-dev chrpath
# mkdir build && cd build
# cmake -Dpostgresql=yes -Dtraymonitor=no -Dmysql=no -Dsqlite3=no .. # make install时用, 而非deb打包时, cmake参数参考`debian/rules`
dpkg-checkbuilddeps
# generate changelog from [here](https://github.com/bareos/bareos/blob/15f82cd288f295f4ae13c3f27775eb2df46f2c98/.travis.yml)
NOW=$(LANG=C date -R -u)
BAREOS_VERSION=$(cmake -P get_version.cmake | sed -e 's/-- //')
printf "bareos (%s) unstable; urgency=low\n\n  * See https://docs.bareos.org/release-notes/\n\n -- nobody <nobody@example.com>  %s\n\n" "${BAREOS_VERSION}" "${NOW}" | tee debian/changelog
vim debian/rules # ~~根据上面的cmake参数定制deb打包编译bareos时需要的参数~~. 不能修改参数, 只能装全依赖, 因为deb打包时dh_install并没有根据参数(比如`-Dmysql=no`, `-Dtraymonitor=no`等)忽略相关依赖文件.
fakeroot debian/rules binary
```

仅cmake编译(非fakeroot打包编译)的缺陷:
1. arm没有vmware插件, 因为依赖的vmware不提供arm so
1. xxx.service 没有User/Group, 可使用root
1. 数据库配置在`/usr/local/etc/bareos-dir.d/catalog/Mycatalog.conf`, 且默认使用sqlite3, 需改用postgres
1. `bareos-dir -t -f -d 500 -v`发现database bareos不存在. 需手动配置db见[这里](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms)

## 概念
- volume : Bareos将在其上写入备份数据的单个物理磁带（或可能是单个文件）
- pool : 定义接收备份数据的多个volume（磁带或文件）组成的逻辑组


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

sudo -u postgres psql # 进入psql
> alter user postgres with password 'postgres'; # 为postgres创建密码
psql -h localhost -p 5432 -U postgres -W # 测试密码登录

systemctl restart postgres

apt install bareos bareos-database-postgresql # 输入db密码. bareos-database-postgresql会利用dbconfig-common mechanism, 在apt install过程中配置db, db配置在`/etc/dbconfig-common/bareos-database-common.conf`. 可用`dpkg-reconfigure bareos-database-common`重新配置, 手动配置db见[这里](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms)

systemctl restart bareos-dir # db config in /etc/bareos/bareos-dir.d/catalog/MyCatalog.conf
systemctl restart bareos-sd
systemctl restart bareos-fd

bareos-dir -t -f -d 500 -v # 测试bareos-dir是否正常, 包括与pg的连接
bareos-fd -t -f -d 500 -v
bareos-dbcheck -B # 作用同上, 显示db的连接信息

apt install bareos-webui # 基于php+apache2
systemctl restart bareos-dir

# -- 配置webui, 也可使用[bconsole configure子命令](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareosWebui.html#create-a-restricted-consoles)
cp /etc/bareos/bareos-dir.d/console/admin.conf.example vim /etc/bareos/bareos-dir.d/console/admin.conf && chown bareos:bareos /etc/bareos/bareos-dir.d/console/admin.conf
vim /etc/bareos/bareos-dir.d/console/admin.conf # 设置bareos-dir admin用于bareos-webui
systemctl restart bareos-dir # 不能省略, 否则可能webui无法登入(账号正确)
systemctl restart apache2 # 访问http://HOSTNAME/bareos-webui即可使用webui, webui也可使用[nginx](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareosWebui.html#nginx), 但访问地址要变为`http://bareos:9100/`
```

### bareos-fd部署
1. 需备份的机器(client端, 使用9102端口, 等待来自bareos-dir的连接)安装客户端软件bareos-filedaemon

    - `apt install bareos-filedaemon`

    > [Windows client下载地址](http://download.bareos.org/bareos/release/20/windows/), `netstat -ano|findstr "9102"`
1. bareos director配置bareos-dir
```bash
$ bconsole
* configure add client name=client2-fd address=192.168.0.2 password=secret # 注册client, 会创建`/etc/bareos/bareos-dir.d/client/client2-fd.conf`和`/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf`(bareos-fd访问bareos-dir的授权, **如果其中不包含Address-<dir_ip>时请添加**)
reload # 不能丢
exit
```
1. 配置clients

    需备份的机器(client端, 使用9102端口, 等待来自bareos-dir的连接)安装客户端软件

    - linux

        ```bash
        # apt install bareos-filedaemon
        # scp dareos-server:/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf /etc/bareos/bareos-fd.d/director # director对client的授权
        # vim /etc/bareos/bareos-fd.d/director/director-mon.conf # 用/etc/bareos/bareos-dir.d/console/bareos-mon.conf文件中的Password替换该文件中的Password
        # systemctl restart bareos-fd
        ```

    - windows

        需要设置的参数:
        - Client Name: bconsole注册clients时的名称, 最好是clients os的hostname
        - Director Name: 不修改
        - Password: dareos-server:/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf中的Password
        - Network Address: 注册client时本机的ip
        - Client Monitor Password: 用/etc/bareos/bareos-dir.d/console/bareos-mon.conf文件中的Password
1. 测试client by bconsole

  ```bash
  status client=client2-fd
  ```

## bconsole cmd
```bash
* reload # 重载配置
* status client # 测试client connection
* status storage
* show client
* show fileset[=xxx]
* list clients
* list pools
* list volumes
* list jobs
* llist jobs # 更详细的`list jobs`
* .jobs # 更精简的`list jobs`, 只有job name
* configure add console name=admin password=pwd111111 profile=webui-admin # 注册bconsole
* configure add client name=client2-fd address=192.168.0.2 password=secret # 注册client, 需要重启bareos-dir
* setdebug client=bareos-fd level=200 # [测试client](https://docs.bareos.org/TasksAndConcepts/TheWindowsVersionOfBareos.html#enable-debuggging)
* configure add job name=client2-job client=client2-fd jobdefs=DefaultJob # 添加job
* run [job=Client1 yes]# 手动开始job, 未指定job时需要选择job
* status [Director]
* restore # 选择文件的命令在[restore-command](https://docs.bareos.org/TasksAndConcepts/TheRestoreCommand.html#restore-command), 被选中的文件名前带`*`
```

## webui
### 还原
- 客户端：从下拉菜单中选择备份所属的客户端
- 备份作业：从下拉菜单中选择需要的备份作业
- 合并所有客户端文件集：自动把该客户端该作业和该作业以前的所有备份（**含不同作业**）集合在一起供恢复文件使用; 如选"否", 只从选择的备份中恢复文件
- 合并所有相关作业：如选"是", 自动把该客户端该作业和该作业以前的所有**同一作业**的备份集合在一起供恢复文件使用; 如选"否", 只从选择的备份中恢复文件
- 还原到客户端：从下拉菜单中选择恢复文件的目标客户端
- 还原作业：从下拉菜单中选择预定义的还原作业
- 替换客户端上的文件：选择同名文件的覆盖规则. 可选规则为：总是、从不、比现有文件旧和比现有文件新
- 要恢复到客户端的位置：指定恢复文件的目标路径
- 文件选择：点击文件/路径前的`□`来选择是否要恢复此文件/路径; 如选择路径, 在该路径下的所有文件都会被恢复

## 流程图/时序

[连接流程](https://docs.bareos.org/TasksAndConcepts/NetworkSetup.html#network-connections-overview)
如果dir-fd之间的连接有正在执行的job, 那么fd不会复用该连接而是在需要时发起新连接, 用于发送其他命令, 比如`cancel`.

[job的网络时序图](https://docs.bareos.org/DeveloperGuide/netprotocol.html#network-sequence-diagrams)

[使用sd device的流程图](https://docs.bareos.org/DeveloperGuide/reservation.html#usedevicecmd)

[job执行的流程图](https://docs.bareos.org/DeveloperGuide/jobexec.html)

## api
bareos console支持非交互式的[点命令](https://docs.bareos.org/DeveloperGuide/api.html#dot-commands), 同时支持json输出(执行`.api json`即可).

### python-bareos
[python-bareos](https://github.com/bareos/bareos/tree/master/python-bareos/)是bareos官方的python sdk, 用于与bareos-dir通信.

### Bareos REST API (based on python-bareos)
参考:
- [README](https://github.com/bareos/bareos/tree/master/rest-api#readme)

```bash
wget https://github.com/bareos/bareos/archive/refs/tags/Release/20.0.1.tar.gz
tar -xf 20.0.1.tar.gz && cd bareos/rest-api
pip3 install -r requirements.txt
vim api.ini # 配置Director并设置secret_key
uvicorn [--host 0.0.0.0 --port 8000] bareos-restapi:app --reload
```

Serve the Swagger UI to explore the REST API: http://127.0.0.1:8000/docs
Alternatively you can use the redoc format: http://127.0.0.1:8000/redoc

原理: 将账户名和密码创建[bareos.bsock.DirectorConsoleJson](https://pypi.org/project/python-bareos/), 再将DirectorConsoleJson和用户名关联, 返回包含该用户名的JWT, 调用restful api with JWT即使用`DirectorConsoleJson.call(cmd)`执行拼接好的cmd.

> 页面有cdn资源依赖. 该功能由fastapi提供, [离线资源加载看这里](https://fastapi.tiangolo.com/advanced/extending-openapi/#self-hosting-javascript-and-css-for-docs), 在自身项目上引入fastapi资源来解决. 注意不能忘记这两属性`FastAPI(docs_url=None, redoc_url=None)`, 否则应用还是使用fastapi默认的渲染函数.

> 只需设置`http://127.0.0.1:8000/docs`页面的"Authorize"按钮里的username和password即可使用openapi的`try it out`

## plugin
> [官方plugins](https://github.com/bareos/bareos/tree/master/contrib)

bareos原生支持dir, storage, filedaemon的插件扩展. 使用插件前必须在配置中启用它们, **修改后需要重启服务**, 当前支持python 2/3. **bareos 20开始推荐使用python3, 虽然官方20.0.1目前plugins都是python2的**.

> 前提: `apt install bareos-{director,storage,filedaemon}-python3-plugin`或`apt install bareos-{director,storage,filedaemon}-python2-plugin`, 都装时先安装python3的.

> [Porting existing Python plugins和Switching to Python 3](https://docs.bareos.org/TasksAndConcepts/Plugins.html)

> [bpluginfo](https://docs.bareos.org/Appendix/BareosPrograms.html#bpluginfo)可用于查看plugin相关信息, 比如`bpluginfo -v /usr/lib/bareos/plugins/python3-fd.so`

> 插件依赖的python package在`core/src/plugins/{dir,file,store}d/python/pyfiles`下, 会由`bareos-{directoor,storage,filedaemon}-python-plugins-common`安装在`/usr/lib/bareos/plugins`下

因为最常用的是fd-plugins, 这里重点介绍. 其他两种请参考[bareos docs](https://docs.bareos.org/TasksAndConcepts/Plugins.html)

### fd-plugins
以官方MySQL Plugin举例:
1. 配置

    - client的`bareos-fd.d/client/myself.conf`的`Plugin Directory`
    - director中的`bareos-dir.d/fileset/mysql.conf`: `Include.Plugin = "python:module_path=/usr/lib64/bareos/plugins:module_name=bareos-fd-mysql"`

        插件参数拼接在Plugin中以`:`分隔即可, 比如`Plugin = "python:module_path=/usr/lib64/bareos/plugins:module_name=bareos-fd-mysql:mysqlhost=dbhost:mysqluser=bareos:mysqlpassword=bareos"`

    > bareos-fd-mysql插件中的[_mysqlbackups_](https://docs.bareos.org/Appendix/Howtos.html#backup-of-mysql-databases-using-the-python-mysql-plugin)是虚拟目录, 说明fd plugins可将io流(mysqldump的输出)发送到storage中.

其他官方插件:
- [`bareos-fd-local-fileset`](https://github.com/aussendorf/bareos-fd-python-plugins/wiki): 备份时动态将filename=/etc/bareos/extra-files中的文件列表加入fileset

fd-plugins其实就是操作fileset, fliter或添加需要备份的文件列表.

## FAQ
### bconsole配置
`/etc/bareos/bconsole.conf`

### bareos-sd配置
> 修改bareos-sd的配置后, 必须重启bareos-sd. 在重启bareos-sd前, 请首先使用`bareos-sd -t -v`检查bareos-sd配置文件, 如它没有任何输出, 说明配置文件没有任何语法问题.

`/etc/bareos/bareos-sd.d`:
- device : [数据存储位置](https://docs.bareos.org/Configuration/StorageDaemon.html#device-resource)

    ```conf
    # HDD 存储设备
    Device {
      Name = FileStorage                  # 设备名称
      Media Type = File                   # 类型, bareos是基于文件的备份/恢复系统, 类型永远是文件
      Archive Device = /bareos/hdd        # Ubuntu下的备份文件目录（或mount point）
      LabelMedia = yes;                   # lets Bareos label unlabeled media
      Random Access = yes;                # 可随机读写
      AutomaticMount = yes;               # 自动加载
      RemovableMedia = no;                # 媒体介质不可移除
      AlwaysOpen = yes;                   # 建议总是打开, FIFO存储设备除外
      Description = "File device. A connecting Director must have the same Name and MediaType"
    }

    # 磁带存储设备
    Device {
      Name = TapeStorage                  # 设备名称
      Media Type = File                   # 类型, bareos是基于文件的备份/恢复系统, 类型永远是文件
      Archive Device = /bareos/tape       # Ubuntu下的mount point
      LabelMedia = yes;                   # lets Bareos label unlabeled media
      Random Access = no;                 # 不能随机读写
      AutomaticMount = no;                # 不自动加载
      RemovableMedia = yes;               # 媒体介质可移除
      AlwaysOpen = yes;                   # 按需打开
    }
    ```
- director

    - bareos-dir.conf : 管理storage对director的授权
    - bareos-mon.conf : 管理storage对bareos traymonitor的授权
- message : storage message管理
    
    - Standard.conf : bareos-sd日志处理
- storage

    - bareos-sd.conf : bareos-sd配置

### bareos-dir配置
> 修改bareos-dir的配置后(比如添加fileset), 必须重启Director. 在重启Director前, 请首先使用`bareos-dir -t -v`检查bareos-dir配置文件. 如命令没有任何输出, 说明配置文件没有任何语法问题.

> 创建文件时注意owner需要是bareos, 否则`systemctl restart bareos-dir`会因为权限导致执行失败.

`/etc/bareos/bareos-dir.d`:
- catalog : 备份/还原索引信息来源

    - MyCatalog.conf : db配置
- client : clients信息

    - xxx.conf : client注册信息
- console

    - admin.conf : web ui访问的授权
    - bareos-mon.conf : monitor访问bareos-dir的授权
- director

    - bareos-dir.conf : bareos-dir配置
- fileset : 备份文件组(定义如何备份一组文件)配置

    - win.conf

        ```conf
        # all office files in users (c:/ and d:/)
        # for win 7     = D
        # for win 10    = C 


        FileSet {
          Name = "Win7_office"
          
          # volume shadow copy service
          Enable VSS = yes
          Include {
          
          # location
            File = "D:/Users"
            File = "D:/My Documents"
          
          Options {
            # config
            Signature = MD5
            compression = LZ4
            IgnoreCase = yes
            noatime = yes
            
            # Word
            WildFile = "*.doc"
            WildFile = "*.dot"
            WildFile = "*.docx"
            WildFile = "*.docm"

            # Excel
            WildFile = "*.xls"
            WildFile = "*.xlt"
            WildFile = "*.xlsx"
            WildFile = "*.xlsm"
            WildFile = "*.xltx"
            WildFile = "*.xltm"

            # Powerpoint
            WildFile = "*.ppt"
            WildFile = "*.pot"
            WildFile = "*.pps"
            WildFile = "*.pptx"
            WildFile = "*.pptm"
            WildFile = "*.ppsx"
            WildFile = "*.ppsm"
            WildFile = "*.sldx"

            # access
            WildFile = "*.accdb"
            WildFile = "*.mdb"
            WildFile = "*.accde"
            WildFile = "*.accdt"
            WildFile = "*.accdr"

            # publisher
            WildFile = "*.pub"

            # open office
            WildFile = "*.odt"
            WildFile = "*.ods"
            WildFile = "*.odp"

            # pdf
            WildFile = "*.pdf"
            
            # flat text / code
            WildFile = "*.xml"
            WildFile = "*.log"
            WildFile = "*.rtf"
            WildFile = "*.tex"
            WildFile = "*.sql"
            WildFile = "*.txt"
            WildFile = "*.tsv"
            WildFile = "*.csv"
            WildFile = "*.php"
            WildFile = "*.sh"
            WildFile = "*.py"
            WildFile = "*.r"
            WildFile = "*.rProj"
            WildFile = "*.js"
            WildFile = "*.html"
            WildFile = "*.css"
            WildFile = "*.htm"
          } 

          # exclude everything else
            Options {
            
            # all files not in include
            RegExFile = ".*"
            
            # default user profiles
            WildDir = "[C-D]:/Users/All Users/*"
            WildDir = "[C-D]:/Users/Default/*"
            
            # explicit don't backup
            WildDir = "[C-D]:/Users/*/AppData"
            WildDir = "[C-D]:/Users/*/Music"
            WildDir = "[C-D]:/Users/*/Videos"
            WildDir = "[C-D]:/Users/*/Searches"
            WildDir = "[C-D]:/Users/*/Saved Games"
            WildDir = "[C-D]:/Users/*/Favorites"
            WildDir = "[C-D]:/Users/*/Links"
          
            # application specific
            WildDir = "[C-D]:/Users/*/MicrosoftEdgeBackups"
            WildDir = "[C-D]:/Users/*/Documents/R"
            WildDir = "*iCloudDrive*"
            WildDir = "*.svn/*"
            WildDir = "*.git/*"
            WildDir = "*.metadata/*"
            WildDir = "*cache*"
            WildDir = "*temp*"
            WildDir = "*OneDrive*"
            WildDir = "*RECYCLE.BIN*"
            WildDir = "[C-D]:/System Volume Information"
            Exclude = yes
          }
           
          }
        }

        FileSet {
          Name = "Win10_office"
          
          # volume shadow copy service
          Enable VSS = yes
          Include {
          
          # location
            File = "C:/Users"
          
          Options {
            # config
            Signature = MD5
            compression = LZ4
            IgnoreCase = yes
            noatime = yes
            
            # Word
            WildFile = "*.doc"
            WildFile = "*.dot"
            WildFile = "*.docx"
            WildFile = "*.docm"

            # Excel
            WildFile = "*.xls"
            WildFile = "*.xlt"
            WildFile = "*.xlsx"
            WildFile = "*.xlsm"
            WildFile = "*.xltx"
            WildFile = "*.xltm"

            # Powerpoint
            WildFile = "*.ppt"
            WildFile = "*.pot"
            WildFile = "*.pps"
            WildFile = "*.pptx"
            WildFile = "*.pptm"
            WildFile = "*.ppsx"
            WildFile = "*.ppsm"
            WildFile = "*.sldx"

            # access
            WildFile = "*.accdb"
            WildFile = "*.mdb"
            WildFile = "*.accde"
            WildFile = "*.accdt"
            WildFile = "*.accdr"

            # publisher
            WildFile = "*.pub"

            # open office
            WildFile = "*.odt"
            WildFile = "*.ods"
            WildFile = "*.odp"

            # pdf
            WildFile = "*.pdf"
            
            # flat text / code
            WildFile = "*.xml"
            WildFile = "*.log"
            WildFile = "*.rtf"
            WildFile = "*.tex"
            WildFile = "*.sql"
            WildFile = "*.txt"
            WildFile = "*.tsv"
            WildFile = "*.csv"
            WildFile = "*.php"
            WildFile = "*.sh"
            WildFile = "*.py"
            WildFile = "*.r"
            WildFile = "*.rProj"
            WildFile = "*.js"
            WildFile = "*.html"
            WildFile = "*.css"
            WildFile = "*.htm"
          } 

          # exclude everything else
            Options {
            
            # all files not in include
            RegExFile = ".*"
            
            # default user profiles
            WildDir = "[C-D]:/Users/All Users/*"
            WildDir = "[C-D]:/Users/Default/*"
            
            # explicit don't backup
            WildDir = "[C-D]:/Users/*/AppData"
            WildDir = "[C-D]:/Users/*/Music"
            WildDir = "[C-D]:/Users/*/Videos"
            WildDir = "[C-D]:/Users/*/Searches"
            WildDir = "[C-D]:/Users/*/Saved Games"
            WildDir = "[C-D]:/Users/*/Favorites"
            WildDir = "[C-D]:/Users/*/Links"
          
            # application specific
            WildDir = "[C-D]:/Users/*/MicrosoftEdgeBackups"
            WildDir = "[C-D]:/Users/*/Documents/R"
            WildDir = "*iCloudDrive*"
            WildDir = "*.svn/*"
            WildDir = "*.git/*"
            WildDir = "*.metadata/*"
            WildDir = "*cache*"
            WildDir = "*temp*"
            WildDir = "*OneDrive*"
            WildDir = "*RECYCLE.BIN*"
            WildDir = "[C-D]:/System Volume Information"
            Exclude = yes
          }
           
          }
        }
        ```
- jobdefs : 备份任务定义, 可被多个作业重复调用, 类似于job template

    ```conf
    JobDefs {
      Name = "TestJob"                                          # 测试任务
      Type = Backup                                             # 类型：备份（Backup）
      Level = Incremental                                       # 方式：递进（Incremental）
      Client = bareos-fd                                        # 被备份客户端：bareos-fd （在Client中定义）
      FileSet = "TestSet"                                       # 备份文件组：TesetSet （在FileSet中定义）
      Schedule = "WeeklyCycle"                                  # 备份周期：WeeklyCy（在schedule中定义）
      Storage = File                                            # 备份媒体： File（在Storage中定义）
      Messages = Standard                                       # 消息方式：Standard（在Message中定义）
      Pool = Incremental                                        # 存储池：Incremental（在pool中定义） 
      Priority = 10                                             # 优先级：10
      Write Bootstrap = "/var/lib/bareos/%c.bsr"                # 
      Full Backup Pool = Full                  # Full备份, 使用 "Full" 池（在storage中定义）
      Differential Backup Pool = Differential  # Differential备份, 使用 "Differential" 池（在storage中定义）
      Incremental Backup Pool = Incremental    # Incremental备份, 使用 "Incremental" 池（在storage中定义）
    }
    ```
- job : 任务配置

    任务类型分: Backup(备份)/Restore(还原), 默认存在的backup-bareos-fd.conf和BackupCatalog.conf是备份job, RestoreFiles.conf是还原job.

    ```conf
    Job {
      Name = "backup-test-on-bareos-fd"              # 任务名
      JobDefs = "TestJob"                            # 使用已定义的备份任务TestJob （在jobdefs中定义）
      Client = "bareos-fd"                           # 客户端名称： bareos-fd（在client中定义）
    }
    ```
- storage : 备份保存位置的配置

    ```conf
    Storage {
      Name = File
      Address = bareos                # director-sd名字, 使用FQDN (不要使用 "localhost" ).
      Password = "JgwtSYloo93DlXnt/cjUfPJIAD9zocr920FEXEV0Pn+S"
      Device = FileStorage            # 在bareos-sd中定义
      Media Type = File
    }
    ```

    > Device, Media Type项必须与bareos-sd定义的一致
- pool : pool配置

    - full : 完整备份

        ```conf
        Pool {
          Name = Full
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 365 days         # Volume有效时间
          Maximum Volume Bytes = 50G          # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Full-"              # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - incremental : 增量备份, 备份所有状态变化的文件. 前提是有full备份, 否则转为full备份.

        ```conf
        Pool {
          Name = Incremental
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 30 days          # Volume有效时间
          Maximum Volume Bytes = 1G           # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Incremental-"       # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - differential : 差异备份, 备份所有modified标志变化的文件. 前提是有full备份, 否则转为full备份.

        ```conf
        Pool {
          Name = Differential
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 90 days          # Volume有效时间
          Maximum Volume Bytes = 10G          # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Differential-"      # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - scratch: 当系统找不到需要的volume时, 自动使用该pool. 该pool名称不可修改, 其他pool名称没有重命名限制.
- schedule: 计划配置

    ```conf
    Schedule {
      Name = "WeeklyCycle"
      Run = Full 1st sat at 21:00                   # 每月第一个周六/晚九点, 完整备份
      Run = Differential 2nd-5th sat at 21:00       # 其余周六/晚九点, 差异备份
      Run = Incremental mon-fri at 21:00            # 周一至周五, 递增备份
    }
    ```
- message : 提示信息(job完成后如何发送提示信息)的配置

    ```conf
    Messages {
      Name = Standard
      Description = "Reasonable message delivery -- send most everything to email address and to the console."
      # operatorcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: Intervention needed for %j\" %r"
      # mailcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: %t %e of %c %l\" %r"
      operator = root@localhost = mount                                 # 执行operatorcommand命令, 用户：root@localhost, 操作：mount
      mail = root@localhost = all, !skipped, !saved, !audit             # 执行mailcommand, 用户：root@localhost, 操作：所有（除skipped, saved和audit）
      console = all, !skipped, !saved, !audit                           # 所有操作, 除skipped, saved和audit
      append = "/var/log/bareos/bareos.log" = all, !skipped, !saved, !audit  # 所有操作, 除skipped, saved和audit
      catalog = all, !skipped, !saved, !audit                           # 所有操作, 除skipped, saved和audit
       # 可用参数
      # %% = %
      # %c = Client’s name
      # %d = Director’s name
      # %e = Job Exit code (OK, Error, ...)
      # %h = Client address
      # %i = Job Id
      # %j = Unique Job name
      # %l = Job level
      # %n = Job name
      # %r = Recipients
      # %s = Since time
      # %t = Job type (e.g. Backup, ...)
      # %v = Read Volume name (Only on director side)
      # %V = Write Volume name (Only on director side)
      # console：定义发送到console的信息
      # append：定义发送到日志文件的信息
      # catalog：定义发送到数据库的信息
    }
    ```
- profile : 定义一组访问控制用于针对不同控制台或角色

### fileset
- `One FS=no` : no, 不检查是否在同一个fs上; yes, 检查是否在同一个fs上
- `FS Type=ext4` : 支持备份的fs类型
- `File=/` : 备份开始位置
- `Exclude {}` : 排除位置
- `WildDir` : 指定文件
- `Exclude = yes`: 排除`WildDir`指定的文件

### backup参数
```conf
Run Backup job
JobName:  backup-test-on-bareos-fd
Level:    Full
Client:   lswin7-1-fd
Format:   Native
FileSet:  TestSet
Pool:     Full (From Job FullPool override)
Storage:  File (From Job resource)
When:     2018-10-05 10:39:59
Priority: 10
OK to run? (yes/mod/no):
```
### restore参数
```conf
Run Restore job
JobName:         RestoreFiles
Bootstrap:       /var/lib/bareos/client1.restore.3.bsr
Where:           /tmp/bareos-restores
Replace:         Always
FileSet:         Full Set
Backup Client:   client1
Restore Client:  client1
Format:          Native
Storage:         File
When:            2013-06-28 13:30:08
Catalog:         MyCatalog
Priority:        10
Plugin Options:  *None*
OK to run? (yes/mod/no):
```

### bconsole命令行调用形式
bconsole是交互式命令, 无法直接后接子命令的形式试用, 因此使用:
```bash
bconsole -c ./bconsole.conf <<END_OF_DATA
show pool
quit
END_OF_DATA
```

[组合使用(备份+还原)](https://docs.bareos.org/TasksAndConcepts/BareosConsole.html#running-the-console-from-a-shell-script):
```bash
bconsole <<END_OF_DATA
@output /dev/null
messages
@output /tmp/log1.out
label volume=TestVolume001
run job=Client1 yes
wait
messages
@#
@# now do a restore
@#
@output /tmp/log2.out
restore current all
yes
wait
messages
@output
quit
END_OF_DATA
```

### job执行过程中报`BnetHost2IpAddrs() for host "ubuntu-18" failed: ERR=`
ubuntu-18是storage daemon的参数在`/etc/bareos/bareos-dir.d/storage/File.conf`的`Address`.

file daemon备份时, 从dareos-dir获取storage参数, 因为网络中没有dns, 因此无法获取到storage的ip.

解决方法: 将Address的参数换成ip即可.

> 错误来源: `/var/log/bareos/bareos.log`或 webui中job的log

### job备份windows文件时报`no drive letters found for generating vss snapshots`
fileset中备份文件路径错误.

错误路径: `File=/c/dsDefault.log`, 正确路径: `File="C:/dsDefault.log"`

### job备份Windows 10文件报`error:14094417:SSL routines:ssl3_read_bytes:sslv3 alert illegal parameter`, `TLS negotiation failed(while probing client protocol)`和`Network error during CRAM MD5 With 192.168.0.197`

此时Windows 10 log报"SSL routines:tls_psk_do_binder:binder does not verify", `TLS negotiation failed`.

解决方法: 卸载并重新安装bareos windows client, 安装时填入正确的参数即可.

> 出问题时安装是使用默认参数(即错误参数), 安装完成后修正`C:\Program Files\Bareos\defaultconfigs\bareos-fd.d\director`下的`*.conf`并重启`Bareos File Backup Service`进行配置的.

### 修改Director邮件发送命令
参考:
- [备份/恢复系统BAREOS的安装、设置和使用（四）](https://blog.csdn.net/laotou1963/article/details/82939355)

在Director默认使用bsmtp发送邮件, 由于bsmtp的局限性，无法使用一般外部商业SMTP服务，我们必须对此进行修改。在示例中，我们对/etc/bareos/bareos-dir.d/message/Standard.conf做修改，您可以参照示例，对其他的邮件发送配置做对应的修改。

配置文件中的默认邮件命令为：
`mailcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: %t %e of %c %l\" %r"`

改为: `mailcommand = "/usr/local/bin/sendmail -c %c -d %d -e %e -h %h -i %i -j %j -n %n -r %r -t %t -s \"%s\"  -l %l -v \"%v\" -V \"%V\%"`

`/user/local/bin/sendmail`是自定义的发送邮件脚本程序. 以`%`开头的是在Bareos中可用的参数，可把所有可用参数全部传递到脚本程序.

> ps: `%s、%v和%V`用`" "`包起来的原因是，这些参数有可能为空，如不把它们包起来，当它们为空时，会造成参数处理问题.

```bash
#!/usr/bin/env bash
# available mailcommand parameters
# %% = %
# %c = Client’s name
# %d = Director’s name
# %e = Job Exit code (OK, Error, ...)
# %h = Client address
# %i = Job Id
# %j = Unique Job name
# %l = Job level
# %n = Job name
# %r = Recipients
# %s = Since time
# %t = Job type (e.g. Backup, ...)
# %v = Read Volume name (Only on director side)
# %V = Write Volume name (Only on director side)

bareos_admin="admin@lswin.cn"
mail_receiver="s.zhang@lswin.cn"

# get input opts
while getopts ":c:d:e:h:i:j:l:n:r:s:t:v:V:" o; do
  case "${o}" in
    c)
       client_name=${OPTARG}
       ;;
    d)
       director_name=${OPTARG}
       ;;
    e)
       job_exit_code=${OPTARG}
       ;;
    h)
       client_address=${OPTARG}
       ;;
    i)
       job_id=${OPTARG}
       ;;
    j)
       unique_job_name=${OPTARG}
       ;;
    l)
       job_level=${OPTARG}
       ;;
    n)
       job_name=${OPTARG}
       ;;
    r)
       recipients=${OPTARG}
       ;;
    s)
       since_time=${OPTARG}
       ;;
    t)
       job_type=${OPTARG}
       ;;
    v)
       read_volume_name=${OPTARG}
       ;;
    V)
       write_volume_name=${OPTARG}
       ;;
    *)
       ;;
    esac
done

# 建立邮件 Subject
ubject="BAREOS任务执行"
if [[ "$job_exit_code" == "OK" ]]
then
  Subject=$Subject"完成通知"
else
  Subject=$Subject"失败通知！"
fi

# 建立邮件内容
Content="\"任务 "$job_name" 执行简况:\n 任务ID："$job_id"\n 任务名字："$unique_job_name"\n 任务类型："$job_type
if [[ ! -z "$job_level" && "$job_type" == "Backup" ]]; then Content=$Content"\n 备份级别："$job_level; fi
Content=$Content"\n 完成情况："$job_exit_code"\n 主控端名字："$director_name"\n 客户端名字："$client_name"\n 客户端地址："$client_address
if [[ ! -z "$read_volume_name" && "$job_type" == "RestoreFiles" ]]; then Content=$Content"\n 读取卷名字："$read_volume_name; fi
if [[ ! -z "$write_volume_name" && "$job_type" == "Backup" ]]; then Content=$Content"\n 写入卷名字："$write_volume_name; fi
Content=$Content"\""

# 建立邮件发送命令
cmd="echo -e $Content | /usr/bin/mail -s \"Subject: $Subject\" -r $bareos_admin $mail_receiver"

# 执行邮件发送命令
eval $cmd

exit 0
```

email example:
```conf
Subject: BAREOS任务执行完成通知

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:26  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 backup-bareos-fd 执行简况:
 任务ID：52
 任务名字：backup-bareos-fd.2018-10-18_10.26.39_12
 任务类型：Backup
 备份级别：Full
 完成情况：OK
 主控端名字：bareos-dir
 客户端名字：bareos-fd
 客户端地址：localhost
 写入卷名字：Full-0001

# ----
Subject: BAREOS任务执行失败通知！

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 backup-test-on-bareos-fd 执行简况:
 任务ID：53
 任务名字：backup-test-on-bareos-fd.2018-10-18_10.42.13_17
 任务类型：Backup
 备份级别：Full
 完成情况：Error
 主控端名字：bareos-dir
 客户端名字：lscms-fd
 客户端地址：lscms.lswin.cn

 # ---
Subject: BAREOS任务执行完成通知

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang <s.zhang@lswin.cn>
任务 RestoreFiles 执行简况:
 任务ID：54
 任务名字：RestoreFiles.2018-10-18_10.43.18_37
 任务类型：Restore
 完成情况：OK
 主控端名字：bareos-dir
 客户端名字：bareos-fd
 客户端地址：localhos

# ---
Subject: BAREOS任务执行失败通知！

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 RestoreFiles 执行简况:
 任务ID：55
 任务名字：RestoreFiles.2018-10-18_10.44.20_01
 任务类型：Restore
 完成情况：Error
 主控端名字：bareos-dir
 客户端名字：lswin7-1-fd
 客户端地址：lswin7-1.lswin.cn
```
### BVFS
BVFS（Bareos虚拟文件系统）提供了一个API来浏览目录中的备份文件并选择文件进行恢复.

### bareos webui如何获取data
以job列表页`localhost:9100/job/`举例, 找到其ajax req(`localhost:9100/job/getData/?data=jobs&period=7&sort=jobid&order=desc`)

在bareos webui root(`/usr/share/bareos-webui/module/Job`)下执行`grep -r getData`, 在`src/Job/Controller/JobController.php`中找到`getDataAction()`, 再在其中找到关键函数`getJobs`.

执行`grep -r getJobs`, 在`src/Job/Model/JobModel.php`中找到它, 看其实现基本可推断是基于bsock, 通过`$bsock->send_command()`逆推, 在`src/Job/Controller/JobController.php`中找到`$this->bsock=$this->getServiceLocator()->get('director')`.

在`/usr/share/bareos-webui`执行`grep -r "send_command" |grep -v "bsock"`, 在`vender/Bareos/library/Bareos/BSock/BareosBSock.php`找到其实现(需考虑send_command有参数列表). 在找到它的上层函数send(), 发现它是操作`fwrite($this->socket,...)`, 找到socket定义: [`stream_socket_client()`](https://php.golaravel.com/function.stream-socket-client.html).

### log
使用`-d 500`参数, 可打印详细日志

bareos-dird log在`/var/log/bareos/bareos.log`
bareos-fd log在systemd.

### 使用官方plugin [bareos-fd-mysql](https://docs.bareos.org/Appendix/Howtos.html#backup-mysql-python)执行job时报`... PluginSave: Command plugin "<python plugin>" required, but is not loaded`
fd `/etc/bareos/bareos-fd.d/client/myself.conf`配置:
```
Client {
  ...

  # remove comment from "Plugin Directory" to load plugins from specified directory.
  # if "Plugin Names" is defined, only the specified plugins will be loaded,
  # otherwise all filedaemon plugins (*-fd.so) from the "Plugin Directory".
  #
  Plugin Directory = "/usr/lib/bareos/plugins"
  Plugin Names = "python"

  ...
}
```

使用`-d 500`参数, 打印详细日志可见, fd log提示`field/fd_plugins.cc:1750-0 No plugin loaded`.

结合myself.conf和日志调试发现, 只要启用了`Plugin Names`即使其value为空, 均会按`Plugin Names`指定的名称去load plugin. 将`Plugin Names`注释默认加载全部插件即可.

### 使用自编译bareos 20.0.1 arm版本, linux备份还原正常, 官方对应版本的windows client无法备份
dir, sd, fd均无报错.