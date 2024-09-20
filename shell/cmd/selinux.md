# selinux
SELinux（ Security-Enhanced Linux）是美国国家安全局在 Linux 开源社区的帮助下开发的一个强制访问控制（ MAC， Mandatory Access Control）的安全子系统. Linux 系统使用 SELinux
技术的目的是为了让各个服务进程都受到约束，使其仅获取到本应获取的资源.

SELinux 服务有 3 种配置模式:
- enforcing：强制启用安全策略模式，将拦截服务的不合法请求
- permissive：遇到服务越权访问时，只发出警告而不强制拦截
- disabled：对于越权的行为不警告也不拦截

## example
```bash
# getenforce # 获得当前 SELinux服务的运行模式
# setenforce 0 # 修改SELinux 当前的运行模式（ 0 为禁用， 1 为启用）
# vim /etc/selinux/config # 禁用selinux, 需要reboot
...
SELINUX=disabled
...
# ls -Zd /home/wwwroot
drwxrwxrwx. root root unconfined_u:object_r:home_root_t:s0 /home/wwwroot
```

`unconfined_u:object_r:home_root_t:s0`说明: 在文件上设置的 SELinux 安全上下文是由用户段、角色段以及类型段等多个信息项共同
组成的. 其中，用户段 system_u 代表系统进程的身份，角色段 object_r 代表文件目录的角色, 类型段 httpd_sys_content_t 代表网站服务的系统文件.

# semanage
用于管理 SELinux 的策略，全称为`SELinux manage`.

## 选项
- -l : 查询
- -a : 添加
- -m : 修改
- -d : 删除

## example
```bash
# --- 让这个目录以及里面的所有文件能够被 httpd 服务程序访问到
# semanage fcontext -a -t httpd_sys_content_t /home/wwwroot
# semanage fcontext -a -t httpd_sys_content_t /home/wwwroot/*
# semanage port -a -t http_port_t -p tcp 6111 # 允许端口
# restorecon -Rv /home/wwwroot/ # 将设置好的 SELinux 安全上下文立即生效. -Rv 表示对指定的目录进行递归操作，以及显示 SELinux 安全上下文的修改过程
# getsebool -a | grep http # 查询并过滤出所有与 HTTP 协议相关的安全策略
# setsebool -P httpd_enable_homedirs=on # `-P`让修改后的 SELinux 策略规则永久生效且立即生效
```

## FAQ