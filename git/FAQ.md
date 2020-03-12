### 中文文件名是乱码(其实是UNICODE编码)

```shell
# 此时不会对0x80以上的字符进行quote,中文显示正常
git config --global core.quotepath false
```
### git diff

`git diff`：是查看working tree与index file的差别
`git diff --cached`：是查看index file与commit的差别
`git diff HEAD`：是查看working tree和commit的差别

`git diff commit_id1 commit_id2` : 比较两个历史版本之间的差异
`git diff dev` : 当前目录和"dev"分支间的差异
`git diff dev...master` : dev分支和master分支的差异,`git diff ...dev`表示当前分支和dev分支的差异

`-- file_name`参数可指定文件名或目录,比如`git diff HEAD -- readme.txt`,`git diff HEAD -- ./lib`.

#### git diff 出现"^M"

Windows用CR LF来定义换行，Linux用LF,Mac用LF(Mac OS 9 及之前是CR，之后换成 UNIX 的 LF). CR全称是Carriage Return ,或者表示为\r, 意思是回车; LF全称是Line Feed，它才是真正意义上的换行表示符.

    git config --global core.whitespace cr-at-eol //让git diff忽略换行符的差异

### 添加/修改配置

    git config --global core.autocrlf input

### 删除配置

    git config --local --unset core.autocrlf

### git crlf换行符

core.autocrlf:

- true : 提交时自动地把行结束符CRLF转换成LF，而在签出代码时把LF转换成CRLF.适用于windows.
- input : 提交时把CRLF转换成LF，签出时不转换.适用于Linux,Mac等类UNIX系统.
- false : 提交和签出时均不转换.

### `.gitignore`无效，不能过滤某些文件



1. `.gitignore` 文件只能作用于 Untracked Files，也就是那些从来没有被 Git 记录过的文件（即自添加以后，从未 add 及 commit 过的文件).

 之所以规则不生效，是因为那些文件曾经被 Git 记录过，因此 .gitignore 对它们完全无效.

 解决方法：

    1. 用git rm从 Git 的数据库中删除对于该文件的追踪；
    2. 把对应的规则写入 .gitignore，让忽略真正生效；
    3. 提交＋推送。

 最后有一点需要注意的，`git rm --cached filename`或`git rm -rf --cached foldername`删除的是追踪状态，而不是物理文件.

2. 注释导致

       ×.go.2 # 忽略测试文件

 注释导致该规则无效，将注释去掉或放在规则的上一行即可.

### 撤销

#### 修改最后一次提交

提交信息写错:

    $ git commit --amend

提交时忘了暂存某些修改,下面的三条命令最终只是产生一个提交，第二个提交命令会修正第一个的提交内容:

    $ git commit -m 'initial commit'
    $ git add forgotten_file
    $ git commit --amend

撤销上一次的提交,提交内容回到暂存区:

    git reset --soft HEAD^

git reset --mixed id,是将git的HEAD变了（也就是提交记录变了），且文件变动回到工作目录.
git reset --soft id,实际上，是git reset –mixed id 后，又做了一次git add
git reset --hard id,是将git的HEAD变了,抛弃文件变动.

HEAD指向的版本就是当前版本,Git允许我们使用命令`git reset --hard commit_id`在版本的历史之间穿梭.
`git reset`可修改当前HEAD指针.

#### 后悔药

- 要用git reflog查看命令历史，以便确定要回到未来的哪个版本,再用`git reset --hard commit_id`返回指定版本.

参考: [`git寻根——^和~的区别`](http://mux.alimama.com/posts/799)

- 在reset之前的提交可以参照ORIG_HEAD,Reset错误时，在ORIG_HEAD上reset 就可以还原到reset前的状态.
```
$ git reset --hard ORIG_HEAD
```

### 历史

#### 更改commit信息

修改历史commit信息:

更新前和remote同步一下.

```
git filter-branch --env-filter '
if test "$GIT_AUTHOR_EMAIL" = "OldEmail"
then
    GIT_AUTHOR_NAME="NewName"
    GIT_AUTHOR_EMAIL="NewEmail"
    GIT_COMMITTER_NAME="NewName"
    GIT_COMMITTER_EMAIL="NewEmail"
fi
export GIT_AUTHOR_NAME
export GIT_AUTHOR_EMAIL
export GIT_COMMITTER_NAME
export GIT_COMMITTER_EMAIL
'
```
如果git报错`Cannot rewrite branches: You have unstaged changes.` 只需要 git stash再运行上面代码.
此时查看`git log`，确认名字和邮箱改好以后，`git push origin master --force`，大功告成！

参考:[Git-工具-重写历史](http://git-scm.com/book/zh/v1/Git-%E5%B7%A5%E5%85%B7-%E9%87%8D%E5%86%99%E5%8E%86%E5%8F%B2)

## 找回丢失的commit/stash
1. 使用`git fsck --no-reflog`或`git fsck --unreachable`找到丢失的commit/stash
1. 使用`git archive`导出数据:
```
# 不推荐使用git show fc30e953fd6cd052fcefda9129de4b1df8830be3,因为有些内容会没展示全.
$ git archive fc30e953fd6cd052fcefda9129de4b1df8830be3 -o a.zip
```

### alias/别名
```sh
$ git config --global alias.br branch
$ git config --global alias.ci commit
$ git config --global alias.co checkout
$ git config --global alias.df diff
$ git config --global alias.lg "log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"
$ git config --global alias.st status
```

### 第一次上传代码时碰到"更新被拒绝，因为您当前分支的最新提交落后于其对应的远程分支..."
这是因为git拒绝合并无关的历史纪录，解决方法是在git pull时加上–allow-unrelated-histories：
```
git pull origin master --allow-unrelated-histories
```

### 如何将一个裸GIT仓库(本地)转换成一个正常的

```bash
git clone -l <path_to_bare_repo> <new_normal_repo>
```

### the remote end hung up unexpectedly
网络比较缓慢(质量差)情况下，就断开了

```
git config --global http.lowSpeedLimit 1000
git config --global http.lowSpeedTime 300
```

以上配置后，只有持续300秒以上的时间内传输速率都低于 1KB/s 的话才会 timeout

### shallow update not allowed
该仓库使用`git clone –depth=1`**浅克隆**而来, 其有限制: 不能将它推送到一个新的存储库.

解决方法: 先使用`git fetch --unshallow old_repo`补全再推送即可.

### submodule 和 subtree的区别
[submodule is link; subtree is copy](https://gb.yekai.net/concepts/subtree-vs-submodule)

### git中Untracked files如何清除
`git clean -f -d`