### `error: insufficient permission for adding an object to repository database .git/objects`

原因:

某次用root账号进行了commit,导致`.git/objects`里的几个文件夹的拥有者和群组变成了root:root,当前用户操作时权限不足而报错.

解决:

使用`chown`命令改回原有的文件夹拥有者和群组即可.
