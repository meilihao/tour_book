# apache2 推荐使用nginx.

也叫httpd. 它支持VirtualHost(基于用户请求的不同 IP 地址、主机域名或端口号).

Apache提供了方便的工具用于切换站点，就是a2ensite和a2dissite，它们都在apache2-common包里.  a2ensite可以激活apache文件夹下sites-available里包含配置文件的站点，a2dissite的作用正好相反.

> 启动命令: `systemctl start httpd`/`service apache2 start`

> 虚拟主机以其指定的DocumentRoot为主.

Linux 系统中的配置文件:
配置文件的名称 存放位置
服务目录 /etc/httpd 
主配置文件 /etc/httpd/conf/httpd.conf 
网站数据目录 /var/www/html 
访问日志 /var/log/httpd/access_log 
错误日志 /var/log/httpd/error_log

## httpd.conf配置
参考:
- [Apache2.4之httpd.conf配置详解](https://blog.csdn.net/a88073327/article/details/80921808)

httpd.conf:
- ServerRoot 服务目录
- ServerAdmin 管理员邮箱
- User 运行服务的用户
- Group 运行服务的用户组
- ServerName 网站服务器的域名
- DocumentRoot 网站数据目录. 定义网站数据的保存路径，默认是`/var/www/html`. 注意修改DocumentRoot后, selinux可能会拦截请求
- Listen 监听的 IP 地址与端口号
- DirectoryIndex 默认的索引页页面
- ErrorLog 错误日志文件
- CustomLog 访问日志文件
- Timeout 网页超时时间，默认为 300 秒
- Alias ：设置路径别名

    语法：Aliase    /alias/        "/path/to/somewhere/"
    这意味着访问http://Server_IP/alias时，其页面文件来自于/path/to/somewhere/这个位置

    例如：Aliase    /images/    "/www/htdocs/imgs/"
    注释：访问：http://192.168.180.100/images/1.gif，就相当于去访问192.168.180.100这台主机  的/www/htdocs/imgs/1.gif

- Direcrory : 设置接入访问权限, 即网站数据目录的权限

  ```conf
    <Direcrory "/path/to/somewhere">
        AddType application/x-httpd-php .htm .html
        DirectoryIndex index.php /index.php # 指定默认页
        Order Allow, Deny # 来定义 Allow或 Deny 指令起作用的顺序，其匹配原则是按照顺序进行匹配，若匹配成功则不执行后面的相应指令
        Options: 开启下面哪些的哪些指令, 多个时用空格分隔
            AllowOverride None : 是否读取目录中的.htaccess文件, 以改变原来设置的权限
              - ALL : 允许
              - None : 不允许
            Indexes：缺少指定的默认页面时，允许将目录中的所有文件以列表形式返回给用户；
            FollowSymLinks：在Direcrory目录下的文件是否允许符号链接到其他目录
            None：所有选项都不启用
            All：所有选项都启用
            ExecCGI：允许使用mod_cgi模块执行CGI脚本
            Includes：允许使用mod_include模块实现SSI(服务器端包含)
            MultiViews：允许使用mod_negotiation(协商),实现内容协商
            SymLinksifOwnerMatch：在链接文件属主数组与原始文件的属主属组相同时，允许跟随符号链接所指向的原始文件；
    </Direcrory>
   ```

  ```conf
  <Directory "/var/www/html/server"> 
  Order allow,deny 
  Allow from 192.168.10.20 
  Order allow,deny 
  Allow from env=ie 
  </Directory> 
   ```
   
- LoadModule :  Apache加载动态文件，如果要与PHP结合，就需要加载PHP的.so文件

## example
```bash
#  htpasswd -c /etc/httpd/passwd linux # 使用 htpasswd 命令生成密码数据库. -c 参数表示第一次生成；后面再分别添加密码数据库的存放文件，以及验证要用到的用户名称（该用户不必是系统中已有的本地账户）
#  vim /etc/httpd/conf.d/userdir.conf
 <Directory "/home/*/public_html"> 
 AllowOverride all 
#刚刚生成出来的密码验证文件保存路径
authuserfile "/etc/httpd/passwd" 
#当用户尝试访问个人用户网站时的提示信息
 authname "My privately website" 
 authtype basic 
#用户进行账户密码登录时需要验证的用户名称
 require user linuxprobe 
 </Directory> 
```

## FAQ
### 多站点
1. 基于端口
```bash
# cat /etc/httpd/conf/httpd.conf # apache2是`/etc/apache2/ports.conf` 
Listen  81

# cat /etc/httpd/conf.d/xx.conf # apache2是`/etc/apache2/sites-available/default`
NameVirtualHost *:81
<VirtualHost *:81>
  ServerName localhost
  include /etc/httpd/conf-available/xxx.conf # 或不用include而是在/etc/httpd/conf-available/xxx.conf添加Alias /xxx /usr/share/bareos-webui/public来启用
</VirtualHost>

# cat /etc/httpd/conf-available/xxx.conf
<Directory "/usr/share/bareos-webui/public"> 
...
</Directory> 
```

### log
```bash
$ grep -r CustomLog /etc/apache2
$ grep -r ErrorLog /etc/apache2
```

### neokylin安装php报缺httpd-mmn, 但Google httpd-mmn没找到具体package
参考: [Run-Time Dependencies](https://fedoraproject.org/wiki/PackagingDrafts/ApacheHTTPModules#Run-Time_Dependencies)

使用`yumdownloader httpd-mmn`发现就是httpd, 再用`dnf install httpd-2.4.34-18...aarch64.rpm`, 最后再用`dnf install php`即可(也可先用yumdownloader下载再安装).

> [Module Magic Number (MMN)](http://httpd.apache.org/docs/2.0/glossary.html), 其实httpd-mmn就是httpd, 用于确保mod_xxx仅与提供相同二进制模块接口的 httpd 包一起使用

### http2 proxy访问后端https服务报`AH01097 pass request body failed to`
```bash
SSLProxyEngine on # 启用proxy
SSLProxyVerify none
SSLProxyCheckPeerCN off
SSLProxyCheckPeerName off # 当前及上面2项是禁用ssl的部分检查
SSLProxyCheckPeerExpire off # 比较server time和cert created_at
```