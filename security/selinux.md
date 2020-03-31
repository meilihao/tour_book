# SELinux（Security-Enhanced Linux）
Linux 上的一个强制访问控制（MAC，Mandatory Access Control）的安全子系统.

RHEL 7 系统使用SELinux 技术的目的是为了让各个服务进程都受到约束，使其仅获取到本应获取的资源.

SELinux 功能：
- 对服务程序的功能进行限制（SELinux 域限制可以确保服务程序做不了出格的事情）
- 对文件资源的访问限制（SELinux 安全上下文确保文件资源只能被其所属的服务程序进行访问）

    在文件上设置的 SELinux 安全上下文(可用`ls -Zd ${path}`查看)是由用户段、角色段以及类型段等多个信息项共同组成的. 其中，用户段 system_u 代表系统进程的身份，角色段 object_r 代表文件目录的角色，类型段 httpd_sys_content_t 代表网站服务的系统文件.

SELinux 服务有三种配置模式，具体如下:
- enforcing：强制启用安全策略模式，将拦截服务的不合法请求
- permissive：遇到服务越权访问时，只发出警告而不强制拦截
- disabled：对于越权的行为不警告也不拦截

配置文件在`/etc/selinux/config`.

## example
```bash
#  getenforce  # 获得当前 SELinux服务的运行模式
#  setenforce 0 # 修改 SELinux 当前的运行模式（0 为禁用，1 为启用）
#  getsebool -a | grep http # 查询并过滤出所有与 HTTP 协议相关的安全策略
# setsebool -P httpd_enable_homedirs=on # `-P `让修改后的 SELinux 策略规则永久生效且立即生效 
```

## semanage
semanage 命令用于管理 SELinux 的策略，格式为`semanage [选项] [文件]`.

选项:
- -l : 参数用于查询
-  -a : 参数用于添加
-  -m : 参数用于修改
-  -d : 参数用于删除

### example
```bash
# ### 以向新的网站数据目录中新添加一条 SELinux 安全上下文，让这个目录以及里面的所有文件能够被 httpd 服务程序所访问到
# semanage fcontext -a -t httpd_sys_content_t /home/wwwroot
# semanage fcontext -a -t httpd_sys_content_t /home/wwwroot/*
# restorecon -Rv /home/wwwroot/ # 让设置好的 SELinux 安全上下文立即生效
```