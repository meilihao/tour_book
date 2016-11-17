### `error: insufficient permission for adding an object to repository database .git/objects`

原因:

某次用root账号进行了commit,导致`.git/objects`里的几个文件夹的拥有者和群组变成了root:root,当前用户操作时权限不足而报错.

解决:

使用`chown`命令改回原有的文件夹拥有者和群组即可.

### `error: src refspec master does not match any`

原因: 本地版本库为空, 空目录不能提交 (只进行了init, 没有add和commit)

### `remote: fatal: Unable to create temporary file '/xxx.git/./objects/pack/tmp_pack_XXXXXX': Permission denied`

git remote url中的账户无权写入远程库目录，方法就是修改远程库目录的所属用户和所属用户组或在remote url中使用有权限的账户,通常时是`git`

### `Permission denied (publickey,gssapi-keyex,gssapi-with-mic).`

今天迁移gogs,`/home/git/.ssh/authorized_keys`不存在,在gogs上重新添加公钥即可.
