# apache2
## httpd.conf配置
参考:
- [Apache2.4之httpd.conf配置详解](https://blog.csdn.net/a88073327/article/details/80921808)

httpd.conf:
- Alias ：设置路径别名

    语法：Aliase    /alias/        "/path/to/somewhere/"
    这意味着访问http://Server_IP/alias时，其页面文件来自于/path/to/somewhere/这个位置

    例如：Aliase    /images/    "/www/htdocs/imgs/"
    注释：访问：http://192.168.180.100/images/1.gif，就相当于去访问192.168.180.100这台主机  的/www/htdocs/imgs/1.gif

- Direcrory : 设置接入访问权限

  ```conf
    <Direcrory "/path/to/somewhere">
        AddType application/x-httpd-php .htm .html
        DirectoryIndex index.php /index.php # 指定默认页
        Options: 开启下面哪些的哪些指令, 多个时用空格分隔
            AllowOverride None : 忽略所有.htaccess文件
            Indexes：缺少指定的默认页面时，允许将目录中的所有文件以列表形式返回给用户；
            FollowSymLinks：是否将符号连接所指向的文件打开；
            None：所有选项都不启用
            All：所有选项都启用
            ExecCGI：允许使用mod_cgi模块执行CGI脚本
            Includes：允许使用mod_include模块实现SSI(服务器端包含)
            MultiViews：允许使用mod_negotiation(协商),实现内容协商
            SymLinksifOwnerMatch：在链接文件属主数组与原始文件的属主属组相同时，允许跟随符号链接所指向的原始文件；
    </Direcrory>
   ```
   
- LoadModule :  Apache加载动态文件，如果要与PHP结合，就需要加载PHP的.so文件
