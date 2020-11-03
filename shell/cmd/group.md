# group

# groupadd
创建用户组

## 格式
`groupadd [选项] 群组名`

## 选项
- -g : 指定gid, 必须唯一, 不可与其他gid重复.

## FAQ
### 将用户加入组
```
# gpasswd -a  ${USER} docker
# gpasswd -d userName groupName # 移出组, 需要重开terminal
```

# newgrp
切换有效用户组.

# groupmod
修改组信息

# groupdel
删除用户组

## example
```bash
# groupdel abc
```

# groupmems
查看组内成员

# groups
查看用户所在的组

# gpasswd
维护组

## 选项
- -M, --members USER,...   :     设置组 GROUP 的成员列表, **会覆盖原有members**

## example
```bash
$ gpasswd -M hu,hua nogroup # 将ha, hua加入nogroup
```
