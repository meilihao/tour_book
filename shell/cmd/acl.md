# acl
配置fs 的acl.

ACL是由一系列的Access Entry所组成的，每一条Access Entry定义了特定的类别可以对文件拥有的操作权限.
Access Entry有三个组成部分：
- Entry tag type
- qualifier (optional) : 定义了特定用户和拥护组对于文件的权限. 只有user和group才有qualifier，其他的都为空.
- permission

Entry tag type它有以下几个类型：
- ACL_USER_OBJ： 	相当于Linux里file_owner的permission
- ACL_USER： 	    定义了额外的用户可以对此文件拥有的permission
- ACL_GROUP_OBJ： 	相当于Linux里group的permission
- ACL_GROUP： 	    定义了额外的组可以对此文件拥有的permission
- ACL_MASK： 	    定义了ACL_USER,ACL_GROUP_OBJ和ACL_GROUP的最大权限
- ACL_OTHER： 	    相当于Linux里other的permission

> mask即ACL_MASK，　｀#effective:...｀即当前mask限制后的权限
> 当权限位只包含"-"时，可用`-`代替`---`
> 使用`-m`时， 指定的user，group必须存在，否则报错

## example
```
# getfacl test # 查看acl
# getfacl --omit-header ./test.sh
# setfacl -b . # 移除所有acl
# setfacl -k . # 移除所有default acl
# setfacl -m m::--- . # 修改mask
# setfacl [-R] -m u:zhangy:rw- test    #  添加/修改一个用户权限, `-R`:递归修改
# setfacl -m u::r-- a # 没有指定用户时即修改文件所有者的权限
# setfacl -m g:zhangying:r-w test      # 添加/修改一个组权限
# setfacl -x u:tank test    # 清除tank用户在test文件acl规则
# setfacl -m d:u:user1:rwx /test <=> setfacl -d -m u:user1:rwx /test # Default ACL是指对于一个目录进行Default ACL设置，并且在此目录下建立的文件都将继承此目录的ACL
# setfacl --set u::rw,u:testu1:rw,g::r,o::- file1 # --set选项会把原有的ACL项都删除，用新的替代(此时会设置mask)，需要注意的是**一定要包含UGO的设置**，不能象-m一样只是添加ACL就可以了
```