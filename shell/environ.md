# TIME_STYLE

自定义显示时间格式,例如控制ls命令输出时间的格式

export TIME_STYLE='+%Y-%m-%d %H:%M:%S'

# LD_LIBRARY_PATH

Linux环境变量名，该环境变量主要用于指定查找共享库（动态链接库）时除了默认路径之外的其他路径(该路径在默认路径之前查找).
移植程序时的经常碰到需要使用一些特定的动态库，而这些编译好的动态库放在我们自己建立的目录里，这时可以将这些目录设置到LD_LIBRARY_PATH中.

```
# NEW_DIRS是新的动态库路径
export LD_LIBRARY_PATH=NEW_DIRS: $LD_LIBRARY_PATH
```

# UID

用户Id,root的UID是0.

# PS1

sh的提示字符串.可扩展的参数:`\u`:用户名;`\h`:主机名;`\w`:当前工作目录.ubuntu下直接修改无效,需修改.bash_profile等配置文件来生效.

# IFS(内部字段分隔符，Internal Field Separator)

存储定界符,FS默认为空白字符（换行符，制表符，空格）
