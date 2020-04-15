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
```

# newgrp
切换有效用户组.

# groupdel
删除用户组

## example
```bash
# groupdel abc
```