# 远程仓库

#### 基础

```shell
# 列出已经存在的远程分支
$ git remote
origin
# 列出远程分支的详细信息(显示对应的克隆地址)
$ git remote -v
# 查看某个远程仓库的详细信息
$ git remote show [remote-name]
# 添加一个新的远程仓库
git remote add [shortname] [url]
# 重命名某个远程仓库在本地的简称
git remote rename [old-name] [new-name]
# 删除某个远端仓库
git remote rm [remote-name]
# 修改远程仓库的url
git remote set-url [remote-name] url
```