# ldap
目录服务是一个特殊的数据库, 用来保存描述性的、基于属性的详细信息，支持过滤功能.

目录是一个为查询、浏览和搜索而优化的数据库，它成树状结构组织数据，类似文件目录一样.

目录数据库和关系数据库不同，它有优异的读性能，但写性能差，并且没有事务处理、回滚等复杂功能，不适于存储修改频繁的数据. 所以目录天生是用来查询的.

LDAP目录服务是由目录数据库和一套访问协议组成的系统.

LDAP(Light Directory Access Portocol)是基于X.500标准的轻量级目录访问协议, 目前的版本是`v3.0`. LDAP通过IP网络用于管理和访问分布式目录服务. LDAP的主要目的是在分层结构中提供一组记录.

OpenLDAP是LDAP的自由和开源的实现, slapd是其ldap daemon. ApacheLdapStudio是目前最靠谱的ldap客户端工具.

新版的OpenLDAP里已没有`slapd.conf`！取而代之的是`/etc/slapd.d`的文件夹. **配置OpenLDAP最正确的姿势是通过ldapmodify命令执行一系列自己写好的ldif文件，而不要修改任何OpenLDAP装好的配置文件**.

> /etc/ldap/slapd.d 是2.4.x版本新采用的配置文件目录.

## 概念
目录树(DIrecotry Information Treesor简称为DIT)：在一个目录服务系统中，整个目录信息集可以表示为一个目录信息**树**，**树中的每个节点是一个条目**. 目录服务的数据类型主要是字符型.

>　举例: 树的根结点是一个组织的域名（dlw.com），其下分为3个部分，分别是managers、people和group，可将这3个组看作组织中的3个部门：如managers用来管理所有管理人员，people用来管理登录系统的用户，group用来管理系统中的用户组.

条目(Entry)：**每个条目就是一条记录**, 每个条目由属性构成. 每个条目有自己在目录树的唯一名称标识(DN,DistinguishedName).

对象类：与某个实体类型对应的一组属性，**对象类是可以继承的**，这样父类的必须属性也会被继承下来

属性(Attribute-value,简称AV)：描述条目的某个方面的信息，**一个属性由一个属性类型和一个或多个属性值组成**，属性有必须属性和非必须属性.

    LDIF（LDAP Data Interchange Format，数据交换格式）是LDAP数据库信息的一种文本格式，用于数据的导入导出，每行都是"属性: 值"对. LDIF用文本格式表示目录数据库的信息，以方便用户创建、阅读和修改.

    常见属性:
    - dc : Domain Component

    域名的部分，其格式是将完整的域名分成几部分，如域名为example.com变成dc=example,dc=com(一条记录的所属位置)
    - uid : User Id

        用户ID，如`tom`(一条记录的ID)
    - c : Country

        国家，如"CN"或"US"等
    - o : Organization

        组织名(公司)，如"Example, Inc."
    - ou : Organization Unit

        组织单位(部门)，类似于Linux文件系统中的子目录，它是一个容器对象，组织单位可以包含其他各种对象（包括其他组织单元）, 如`market`, 是一条记录的所属组织
    - cn : Common Name

        公共名称(公共名称)，如"Thomas Johansson"(一条记录的名称)
    - sn : Surname

        姓，如"Johansson"
    - dn  Distinguished Name

        惟一辨别名，类似于Linux文件系统中的绝对路径，每个对象都有一个惟一的名称，如"uid= tom,ou=market,dc=example,dc=com"，在一个目录树中DN总是**惟一**的

        Base DN：LDAP目录树的最顶部就是根，也就是所谓的"Base DN"，如"dc=mydomain,dc=org".
    - rdn : Relative dn

        相对辨别名,一般指dn逗号最左边的部分，类似于文件系统中的相对路径，它是与目录树结构无关的部分，如"uid=tom"或"cn= Thomas Johansson"
    - givenName：指一个人的名字，不能用来指姓
    - mail：电子信箱地址
    - telephoneNumber：电话号码
    - l：指一个地名，如一个城市或者其他地理区域的名字

    下面列出部分常用objectClass要求必设的属性:
    - account：userid。
    - organization：o。
    - person：cn和sn。
    - organizationalPerson：与person相同。
    - organizationalRole：cn。
    - organizationUnit：ou。
    - posixGroup：cn、gidNumber。
    - posixAccount：cn、gidNumber、homeDirectory、uid、uidNumber

模式（Schema）
模式是对象类（ObjectClass）、属性类型（AttributeType）、属性语法（Syntax）和匹配规则（MatchingRules）的集合


LDAP操作(4类10中, 除了扩展操作, 其他9中是ldapd标准操作):
1. 查询类

    搜索, 比较
1. 更新类
    
    添加/删除/修改条目, 修改条目名
1. 认证类

    绑定, 解绑定
1. 其他

    放弃和扩展操作

ldap的安全模型:
- 身份认证

    - 匿名认证：即不对用户进行认证，该方法仅对完全公开的方式适用
    - 基本认证：通过用户名和密码进行身份识别, 又分为简单密码和MD5密码认证
    - SASL(Simple Authentication and Secure Layer)认证：即LDAP提供的在SSL和TLS安全通道基础上进行的身份认证，包括数字证书的认证

        SASL有几大工业实现标准：Kerveros V5、DIGEST-MD5、EXTERNAL、PLAIN、LOGIN.

        EXTERNAL一般用于初始化添加schema时使用.

        Kerveros V5是里面最复杂的一种，使用GSSAPI机制，必须配置完整的Kerberos V5安全系统，密码不再存放在目录服务器中，每一个dn与Kerberos数据库的主体对应。DIGEST-MD5稍微简单一点，密码通过saslpasswd2生成放在sasldb数据库中，或者将明文hash存到LDAP dn的userPassword中，每一个authid映射成目录服务器的dn，常和SSL配合使用.
- 安全通道
- 访问控制

统一身份认证主要是改变原有的认证策略，使需要认证的软件都通过LDAP进行认证，在统一身份认证之后，用户的所有信息都存储在AD Server中. 终端用户在需要使用公司内部服务的时候，都需要通过AD服务器的认证.

ldap的操作流程是:
1. 连接到LDAP服务器
1. 绑定到LDAP服务器
1. 在LDAP服务器上执行所需的任何操作
1. 释放LDAP服务器的连接

## 部署openldap
env: Ubuntu 20.04

```bash
apt install slapd ldap-utils
yum install -y openldap openldap-clients openldap-servers

systemctl start slapd
slapd -VV
ss -tnlp | grep 389
slapd -d ? # 支持哪些模块的log
man slapd-config # slapd配置
```

安装过程中需要输入admin entry的密码.

`slapcat`可获取openldap DIT的配置.
`slappasswd`生成根节点管理员密码, 与admin entry的密码相同.
`sudo dpkg-reconfigure slapd` 配置openldap.
`ldapwhoami -H ldap:// -x`测试openldap, 正常时返回`anonymous`.
OpenLDAP主配置文件slapd.conf,该配置文件一般保存在安装目录下的etc/openldap/目录下

## example
ldapadd
- -x : 进行简单认证
- -D : 用来绑定服务器的DN
- -h : 目录服务的地址
- -w : 绑定DN的密码
- -f : 使用ldif文件进行条目添加的文件

ldapmodify
- -a 添加新的条目.缺省的是修改存在的条目.
- -C 自动追踪引用.
- -c 出错后继续执行程序并不中止.缺省情况下出错的立即停止.
- -D binddn 指定搜索的用户名(一般为一dn 值).
- -e 设置客户端证书文件,例: -e cert/client.crt
- -E 设置客户端证书私钥文件,例: -E cert/client.key
- -f file 从文件内读取条目的修改信息而不是从标准输入读取.
- -H ldapuri 指定连接到服务器uri。常见格式为ldap://hostname:port
- -h ldaphost 指定要连接的主机的名称/ip 地址.它和-p 一起使用.
- -p ldapport 指定要连接目录服务器的端口号.它和-h 一起使用.
- -M[M] 打开manage DSA IT 控制. -MM 把该控制设置为重要的.
- -n 用于调试到服务器的通讯.但并不实际执行搜索.服务器关闭时,返回错误；服务器打开时,常和-v 参数一起测试到服务器是否是一条通路.
- -Q : 安静执行
- -v 运行在详细模块.在标准输出中打出一些比较详细的信息.比如:连接到服务器的ip 地址和端口号等.
- -V 启用证书认证功能,目录服务器使用客户端证书进行身份验证,必须与-ZZ 强制启用TLS 方式配合使用,并且匿名绑定到目录服务器.
- -W 指定了该参数,系统将弹出一提示入用户的密码.它和-w 参数相对使用.
- -w bindpasswd 直接指定用户的密码. 它和-W 参数相对使用.
- -x 使用简单认证.
- -Y SASL机制
- -Z[Z] 使用StartTLS 扩展操作.如果使用-ZZ,命令强制使用StartTLS 握手成功.

ldapsearch
 - -x : 进行简单认证
 - -D : 用来绑定服务器的DN
 - -w : 绑定DN的密码
 - -b : 指定要查询的根节点
 - -H : 制定要查询的服务器
 - -L : 输出ldif格式, `-LLL`输出内容不包含注释

ldappasswd(设置使用者密码)
- -x :　进行简单认证
- -D :　用来绑定服务器的DN
- -w :　绑定DN的密码
- -S :　提示的输入密码
- -s :　pass 把密码设置为pass
- -a :　pass 设置old passwd为pass
- -A :　提示的设置old passwd
- -H :　是指要绑定的服务器
- -I :　使用sasl会话方式

LDAP 数据交换格式文件，它以文本形式存储，用于在服务器之间交换数据, 可以跟关系型数据库的 SQL 文件做类比. LDIF 文件的格式一般如下：
```conf
dn: <识别名>
<属性 1>: <值 1>
<属性 2>: <值 2>
...
```

```bash
# cat baseldapdomain.ldif
dn: dc=example,dc=com
objectClass: top
objectClass: dcObject
objectclass: organization
o: example com
dc: example

dn: cn=Manager,dc=example,dc=com
objectClass: organizationalRole
cn: Manager
description: Directory Manager

dn: ou=People,dc=example,dc=com
objectClass: organizationalUnit
ou: People

dn: ou=Group,dc=example,dc=com
objectClass: organizationalUnit
ou: Group
# ldapadd -Y EXTERNAL -x -D cn=Manager,dc=example,dc=com -W -f baseldapdomain.ldif
# cat add_entry.ldif
dn: ou=Marketing, cn=root,dc=kevin,dc=com
changetype: add
objectclass: top
objectclass: organizationalUnit
ou: Marketing
  
dn: cn=Pete Minsky,ou=Marketing,cn=root,dc=kevin,dc=com
changetype: add
objectclass: person
objectclass: organizationalPerson
objectclass: inetOrgPerson
cn: Pete Minsky
sn: Pete
ou: Marketing
description: sb, sx
description: sx
uid: pminsky
# ldapadd -x -D "cn=root,dc=kevin,dc=com" -w secret -f /root/add_entry.ldif # 将 add_entry.ldif 中的数据导入 ldap, 创建一个Marketing部门，添加一个dn记录. `cn=root,dc=kevin,dc=com`是ldap administrator. `-x`表示使用`simple authentication`(即明文用户密码)
# ldapsearch -x -b 'cn=root,dc=kevin,dc=com' '(objectClass=*)' # 验证插入
# cat delete_entry.ldif
dn: cn=Susan Jacobs,ou=Marketing,cn=root,dc=kevin,dc=com
changetype: delete
# ldapdelete -x -D "cn=root,dc=kevin,dc=com" -w secret -f /root/delete_entry.ldif
# ldapdelete -x -D "cn=Manager,dc=test,dc=com" -w secret "uid=test1,ou=People,dc=test,dc=com" 删除：ldapdelete
# cat modify_entry.ldif
dn: cn=Pete Minsky,ou=Marketing,cn=root,dc=it,dc=com
changetype: modify
add: mail
mail: pminsky@example.com
-
replace: sn
sn: Minsky
-
delete: description
description: sx
# ldapmodify -x -D "cn=root,dc=it,dc=com" -W -f modify_entry.ldif # 添加mail属性，修改sn的值，删除一个description属性
# ldappasswd -x -D cn=admin,dc=wecash,dc=net -w weopenldap -H ldapi:/// "cn=wedba,ou=Groups,dc=wecash,dc=net" -S # ldappasswd 命令用于修改密码
# ldapsearch -H ldapi://192.168.0.245:389 -b dc=xx,dc=cn -LLL 
# ldapsearch -Q -LLL -Y EXTERNAL -H ldapi:/// -b cn=config dn # 查看初始化信息
# ldapsearch -x -D "cn=root,dc=kevin,dc=com" -w secret -b "dc=kevin,dc=com" # 使用简单认证，用 "cn=root,dc=kevin,dc=com" 进行绑定，要查询的根是 "dc=kevin,dc=com"。这样会把绑定的用户能访问"dc=kevin,dc=com"下的所有数据显示出来
# ldapsearch -x -W -D "cn=administrator,cn=users,dc=osdn,dc=cn" -b "cn=administrator,cn=users,dc=osdn,dc=cn" -h troy.osdn.zzti.edu.cn
# ldapsearch -b "dc=canon-is,dc=jp" -H ldaps://192.168.10.192:389
# ldappasswd -x -D 'cm=root,dc=it,dc=com' -w secret 'uid=zyx,dc=it,dc=com' -S
# slapcat -l export.ldif # 数据导出
cat change_rootdn.ldif # `ls /etc/ldap/slapd.d/cn=config/`, 找`olcDatabase={N}<M>db.ldif`, 我这里是`olcDatabase={1}mdb.ldif`
dn: olcDatabase={1}mdb,cn=config
changetype: modify
replace: olcRootDN
olcRootDN: cn=admin,dc=qiban,dc=com
-
replace: olcSuffix
olcSuffix: dc=qiban,dc=com
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f change_rootdn.ldif
# cat add_content.ldif
dn: ou=People,dc=example,dc=com
objectClass: organizationalUnit
ou: People

dn: ou=Groups,dc=example,dc=com
objectClass: organizationalUnit
ou: Groups

dn: cn=miners,ou=Groups,dc=example,dc=com
objectClass: posixGroup
cn: miners
gidNumber: 5000

dn: uid=john,ou=People,dc=example,dc=com
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
uid: john
sn: Doe
givenName: John
cn: John Doe
displayName: John Doe
uidNumber: 10000
gidNumber: 5000
userPassword: johnldap
gecos: John Doe
loginShell: /bin/bash
homeDirectory: /home/john
shadowLastChange: 0
shadowMax: 0
shadowWarning: 0
# ldapadd -x -W -D "cn=admin,dc=wecash,dc=net" -f add_content.ldif # posixGroup, posixAccount会建出Linux账号
# cat logging.ldif
dn: cn=config
changetype: modify
replace: olcLogLevel
olcLogLevel: stats
# ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f logging.ldif #  Logging设置
# cat << EOF > chrootpw.ldif 
dn: olcDatabase={0}config,cn=config
changetype: modify
add: olcRootPW
olcRootPW: {SSHA}XsxctHt+Ae3Saq2Kcead4UdZ0kOTZRn8
EOF
# ldapadd -Y EXTERNAL -H ldapi:/// -f chrootpw.ldif # 设置LDAP管理员密码
# slapcat -v -l openldap.ldif # slapcat 命令用于将数据条目转换为 OpenLDAP 的 LDIF 文件，可用于 OpenLDAP 条目的备份以及结合 slapdadd 指定用于恢复条目.
# slaptest -f slapd.conf -F slapd.d # 可转换配置
# slaptest [-d 3] -F /etc/ldap/slapd.d # 检测配置文件的可用性，可设置输出级别
# slappasswd -h {SHA} # 输入密码, slappasswd会输出密码编码后的结果. 也可将SHA替换为SSHA(更安全).
{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g=
# cat change_pwd.ldif
dn: olcDatabase={0}config,cn=config
changetype: modify
add: olcRootPW
olcRootPW: {SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g=
# ldapmodify -Y EXTERNAL -H ldapi:/// -f change_pwd.ldif # olcRootPW已存在时用`replace`替换`add`即可
# ldapsearch -H ldapi:/// -Y EXTERNAL -b "cn=config" -LLL -Q "olcDatabase=*" dn # 查找所有db
olcDatabase={-1}frontend,cn=config # 此条目用于定义特定“frontend”数据库的功能。这是一套pseudo-database，用于定义适用于全部其它数据库（除非特别指定）的全局设置。
olcDatabase={0}config,cn=config # 此条目用于定义我们目前正在使用的cn=config数据库设置。多数情况下，其主要负责访问控制设置与复制配置等。
olcDatabase={1}mdb,cn=config # 此条目定义设置特定类型的数据库（本示例中为mdb）. 其一般负责定义访问控制、数据存储细节、缓存与缓冲、DIT的root条目以及管理细节
# ldapsearch -H ldapi:// -Y EXTERNAL -b "cn=schema,cn=config" -s base -LLL -Q # 查看cn=schema,cn=config条目下内置schema
# ldapsearch -H ldapi:// -Y EXTERNAL -b "cn=schema,cn=config" -s one -Q -LLL dn # 查看系统中已载入的其它schema. 大括号加数字代表该schema被系统读取时的顺序。在添加schema时, 数字一般由系统自动添加
```

## FAQ
### 复制db
把sldap关闭, 然后通过以下三步操作静态同步主从服务器上的数据：
- 把主服务器上/var/lib/ldap目录下的所有数据库文件全部拷贝到从服务器的同目录中，覆盖原有文件
- 把主服务器上的/etc/ldap/schema目录下的所有schema文件拷贝到从服务器的同目录中，覆盖原有文件
- 把主服务器上/etc/ldap/slapd.d目录拷贝到从服务器的同目录中，覆盖原有文件

### 启用tls
[Ubuntu下LDAP 部署文档](https://blog.csdn.net/qq_40907977/article/details/108871300)

```bash
# cat certinfo.ldif
# create new
dn: cn=config
changetype: modify
add: olcTLSCACertificateFile
olcTLSCACertificateFile: /etc/ldap/certs/cacert.pem
-
replace: olcTLSCertificateFile
olcTLSCertificateFile: /etc/ldap/certs/tldap.wecash.net.pem
-
replace: olcTLSCertificateKeyFile
olcTLSCertificateKeyFile: /etc/ldap/certs/tldap.wecash.net-key.pem
# ldapmodify -Y EXTERNAL -H ldapi:/// -f certinfo.ldif # 使用ldapmodify命令通过slapd-config数据库告诉slapd启用TLS
# vim /etc/default/slapd # 需要在/etc/default/slapd中添加`ldap:///`才能使用加密
SLAPD_SERVICES="ldap:/// ldapi:/// ldaps:///"
# cat slapd.ldif
# log
dn: cn=config
changetype: modify
replace: olcLogLevel
olcLogLevel: stats
-
add: olcIdleTimeout
olcIdleTimeout: 30
-
add: olcReferral
olcReferral: ldaps://tldap.wecash.net
-
add: olcLogFile
olcLogFile: /var/log/sladp.log

# ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f slapd.ldif # 修改请求域名. olcLogLevel来自`sladpd- d ?`
# ldapsearch -Y external -H ldapi:/// -b cn=config "(objectClass=olcGlobal)"  olcReferral
# systemctl restart slapd.service
```

### 备份和恢复
[Ubuntu下LDAP 部署文档](https://blog.csdn.net/qq_40907977/article/details/108871300)

```bash
wget --no-check-certificate https://raw.githubusercontent.com/alexanderjackson/ldap-backup-and-restore/master/ldap-backup -O /usr/local/sbin/ldap-backup
wget --no-check-certificate https://raw.githubusercontent.com/alexanderjackson/ldap-backup-and-restore/master/ldap-restore -O /usr/local/sbin/ldap-restore
```

### ldap_bind: Invalid credentials
`ldappasswd -H ldap://172.16.0.21 -x -D "cn=admin,ou=People,dc=expmale,dc=com" -W -S "uid=zhang3,ou=People,dc=expmale,dc=com"`报该错.

管理员DN不正确/管理员密码错误

### 连接ldap server的必选配置
1. server ip
1. server port
1. base dn
1. admin dn

### ldap log
直接配置`olcLogFile: /var/log/sladp.log`没有效果. 因此间接方式(重定向到文件):
1. 配置`olcLogLevel: stats`或修改/etc/default/slapd

    ```bash
    $ grep SLAPD_OPTIONS /etc/default/slapd
    SLAPD_OPTIONS="-s 256"
    ```
1. 通过systemd unit查看(**推荐**)或使用rsyslog

    ```bash
    # vim etc/rsyslog.d/10-slapd.conf
    $template slapdtmpl,"[%$DAY%-%$MONTH%-%$YEAR% %timegenerated:12:19:date-rfc3339%] %app-name% %syslogseverity-text% %msg%\n"
    local4.*    /var/log/slapd.log;slapdtmpl
    # systemctl restart rsyslog
    # ldapsearch -Y external -H ldapi:/// -b dc=ldaptuto,dc=net
    # cat /var/log/slapd.log
    ```