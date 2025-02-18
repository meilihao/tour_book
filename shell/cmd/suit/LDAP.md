# LDAP
ref:
- [Linux安装ldap服务端](https://blog.csdn.net/weixin_54975325/article/details/138959072)
- [Ubuntu 24.04 : OpenLDAP](https://www.server-world.info/en/note?os=Ubuntu_24.04&p=openldap&f=1)

LDAP(Lightweight Directory Access Protocol)是轻量目录访问协议, 是用于管理分层数据的协议. 它提供了一种存储，组织和管理组织数据的方法，例如员工帐户和计算机. 它促进了集中式身份验证和授权管理.

LDAP的目录服务其实也是一种数据库系统（Berkeley DB），只是这种数据库是一种树形结构（B Tree），适合读不适合频繁写，不支持事务不能回滚.

## LDAP的树形结构
从叶子到根的这条"路径"是一条数据称为条目（Entry），这条数据的全局唯一标识叫做dn.
- dc（Domain Component）是域名的一部分，把完整的域名拆开
- ou（Organization Unit）是组织单元
- cn（Common Name）一般是用户的名字
- uid（User Id）一般是用户登录id

简单理解就是：dn是一条记录的详细位置； dc是一个区域（相当于哪颗树）；ou是一个组织（相当于哪一个分支）；cn/uid（分支上的哪个苹果）.

## 安装server
```bash
# yum install -y openldap openldap-clients openldap-servers migrationtools
# apt -y install slapd ldap-utils # 安装时需要输入Administrator密码
# dpkg-reconfigure slapd # 配置slapd
- DNS domain name: example.com
- Organization name: test
- Administrator password: password
```

安装slapd后将创建一个最小的工作配置, 其中包含一个顶级条目和一个管理员dn.

配置:
```bash
# slapcat #查看配置
# slaptest -u # 检测配置是否正确, 出现succeeded表示验证成功
# systemctl status slapd # 查看ldap server状态
# ss -antup | egrep "389|636" # 检测服务及端口
```

slapd被设计为在服务本身内进行配置，为此目的专用一个单独的DIT(Directory Information Tree). 这允许动态配置slapd，而无需重新启动服务或编辑配置文件. 此配置数据库由一组基于文本的LDIF文件组成，这些文件位于/etc/ldap/slapd.d下，但决不能直接编辑这些文件. 配置方法有:
1. slapd-config

    by ldapsearch
2. Real Time Configuration (RTC)
3. cn=config

安装slapd后, 将获得两个数据库或后缀：一个用于用户的数据，基于用户的主机域（dc=example,dc=com），另一个用于配置，其根位于cn=config. 要更改每个数据，需要不同的凭据和访问方法
1. `dc=example,dc=com`: 此后缀的管理用户为cn=admin,dc=example,dc=com，其密码为安装slapd软件包时选择的密码
2. `cn=config`: slapd本身的配置存储在此后缀下. 可以通过特殊dn `gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth`进行修改. 当使用 SASL EXTERNAL 身份验证通过 /run/slapd/ldapi Unix 套接字的 `ldapi:///` 传输时，目录会这样看到本地系统的 root 用户 (uid=0/gid=0). 这实际上意味着只有本地 root 用户才能更新 cn=config 数据库

ldapsearch:
```bash
# ldapsearch -Q -LLL -Y EXTERNAL -H ldapi:/// -b cn=config dn
dn: cn=config

dn: cn=module{0},cn=config

dn: cn=schema,cn=config

dn: cn={0}core,cn=schema,cn=config

dn: cn={1}cosine,cn=schema,cn=config

dn: cn={2}nis,cn=schema,cn=config

dn: cn={3}inetorgperson,cn=schema,cn=config

dn: olcDatabase={-1}frontend,cn=config

dn: olcDatabase={0}config,cn=config

dn: olcDatabase={1}mdb,cn=config
# ldapsearch -x -LLL -H ldap:/// -b dc=example,dc=com dn
dn: dc=example,dc=com
```

理解:
- cn=config: Global settings # cn=config：全局设置
- cn=module{0},cn=config: # 动态加载模块
- cn=schema,cn=config: # 包含硬编码的系统级schema
- cn={0}core,cn=schema,cn=config: 硬编码的核心schema
- cn={1}cosine,cn=schema,cn=config: The Cosine schema
- cn={2}nis,cn=schema,cn=config: 网络信息服务 (NIS) schema
- cn={3}inetorgperson,cn=schema,cn=config: InetOrgPerson schema
- olcDatabase={-1}frontend,cn=config: 前端数据库，其他数据库的默认设置
- olcDatabase={0}config,cn=config: slapd配置数据库（cn=config）
- olcDatabase={1}mdb,cn=config: 用户的数据库实例, 这里即`dc=example,dc=com`

- dc=example,dc=com: DIT的基础

ldapsearch两种不同的身份验证机制:
- `-x`: 简单绑定, 本质上是一个纯文本身份验证. 由于没有提供绑定DN（通过-D），这就变成了一个匿名绑定. 如果没有-x，默认情况下将使用简单身份验证安全层（SASL）绑定
- `-Y EXTERNAL`: 使用SASL绑定（没有提供-x），并进一步指定EXTERNAL类型. 与`-H ldapi:///`一起，它使用本地UNIX套接字连接

验证身份验证工具:
```bash
# ldapwhoami -x
# ldapwhoami -x -D cn=admin,dc=example,dc=com -W
# ldapwhoami -Y EXTERNAL -H ldapi:/// -Q
dn:gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth
```

当使用简单绑定（-x）并指定绑定DN（-D作为身份验证DN）时，服务器将在条目中查找userPassword属性，并使用该属性来验证凭据. 但使用数据库根DN条目，即实际的管理员时，这是一种特殊情况，其密码在安装软件包时在配置中设置.

当通过ldapi:///传输使用SASL EXTERNAL时，绑定DN变成连接用户的uid和gid的组合，后跟后缀cn=peercred,cn=external,cn=auth. 服务器会知道这一点，并通过SASL机制授予本地root用户对cn=config的完全写访问权限

演示要求:
1. A node called People, to store users:名为People的节点，用于存储用户

    A user called john:一个叫john的用户

1. A node called Groups, to store groups:一个名为Groups的节点，用于存储组

    1. A group called miners:一个叫矿工的组织

```bash
# vim add_content.ldif # 目录中的uid和gid值不要与本地值冲突
dn: ou=People,dc=example,dc=com
objectClass: organizationalUnit
ou: People

dn: ou=Groups,dc=example,dc=com
objectClass: organizationalUnit
ou: Groups

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
userPassword: {CRYPT}x
gecos: John Doe
loginShell: /bin/bash
homeDirectory: /home/john

dn: cn=miners,ou=Groups,dc=example,dc=com
objectClass: posixGroup
cn: miners
gidNumber: 5000
memberUid: john
# ldapadd -x -D cn=admin,dc=example,dc=com -W -f add_content.ldif
# ldapsearch -x -LLL -b dc=example,dc=com '(uid=john)' cn gidNumber # 检查信息是否已正确添加, 比如搜索`john`条目的cn和gidnumber属性
# ldappasswd -x -D cn=admin,dc=example,dc=com -W -S uid=john,ou=people,dc=example,dc=com # 重置john的密码
```

`john`条目的userPassword字段设置为加密值{CRYPT}x, 这本质上是一个无效的密码，因为任何哈希运算都不会产生x. 这是添加没有默认密码的用户条目时的常见模式。要将密码更改为有效的密码，可以使用ldappasswd设置

更改RootDN密码:
```bash
# --- 方法1:
# slappasswd : 生成密码的哈希值
...
{SSHA}VKrYMxlSKhONGRpC6rnASKNmXG2xHXFo
# vim changerootpw.ldif
dn: olcDatabase={1}mdb,cn=config
changetype: modify
replace: olcRootPW
olcRootPW: {SSHA}VKrYMxlSKhONGRpC6rnASKNmXG2xHXFo
# ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f changerootpw.ldif
# --- 方法2:
# ldappasswd -x -D cn=admin,dc=example,dc=com -W -S
```

turenas配置ldap:
1. 打开webui->Credentials->Directory Services
1. 配置, 并查看LDAP Status=HEALTHY即成功

    - hostname: ip/域名
    - base dn: dc=example,dc=com
    - Bind DN: cn=admin,dc=example,dc=com
    - Bind Password: password
    - Enable: true

## FAQ
### dpkg-reconfigure slapd报`... Make sure that the DNS domain name is syntactically valid, ...`
键盘问题导致密码输错