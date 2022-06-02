# sudo
可把特定命令的执行权限赋予给指定用户

sudo 命令具有如下功能：
- 限制用户执行指定的命令
- 记录用户执行的每一条命令
- 配置文件（/etc/sudoers）提供集中的用户管理、权限与主机等参数
- 验证密码的后 5 分钟内（默认值）无须再让用户再次验证密码

可使用 visudo 命令配置 sudo 命令的配置文件. 比如追加`chen ALL=(ALL) NOPASSWD: ALL`可让用户chen使用sudo时无需密码且没有使用命令的限制.

## 选项
- -l : 列出当前用户可执行的命令
- -u : 用户名或 UID 值 以指定的用户身份执行命令
- -k : 清空密码的有效时间，下次执行 sudo 时需要再次进行密码验证
- -b : 在后台执行指定的命令
- -p : 更改询问密码的提示语

## examples
```bash
$ sudo KKZONE=cn env # sudo 传入env
```

## FAQ
### 添加sudo
```sh
# visudo # 添加`%sudo	ALL=(ALL:ALL) ALL`, 即属于sudo用户组的用户均可使用sudo命令
```

`tidb ALL=(ALL) NOPASSWD:ALL`表示使用sudo命令时无需密码.