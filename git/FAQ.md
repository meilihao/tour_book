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
该仓库使用`git clone --depth=1`**浅克隆**而来, 其有限制: 不能将它推送到一个新的存储库.

解决方法: 先使用`git fetch --unshallow <origin_repo>`补全再推送即可.

### submodule 和 subtree的区别
[submodule is link; subtree is copy](https://gb.yekai.net/concepts/subtree-vs-submodule)

### git中Untracked files如何清除
`git clean [-x -n] -f -d`

选项:
- -f : 删除当前目录下untracked的文件, 不包括文件夹和.gitignore中指定的文件和文件夹
- -d : 删除当前目录下untracked的文件和文件夹, 不包括.gitignore中指定的文件和文件夹
- -x : 删除当前目录下所有的untracked的文件和文件夹
- -n : 仅显示会被删除的文件和文件夹, 但并不实际删除

### 查看远程分支的tag信息
`git ls-remote -t  https://review.coreboot.org/coreboot.git`

### git clone --depth导致无法checkout origin的其他分支
origin的其他分支是基于某个历史commit开始的, 因为使用了`--depth`导致该commit不存在, 因此`git checkout -b dev origin/${branch}`会失败.

解决方法:
1. 补全commit
1. 重新完整地clone

### git diff filter
```bash
git diff b446759c1be2fa8a2c4532ac2bffb4b0449994bc HEAD -- ./xxx/* # 按照path过滤
```

### git查看commit提交记录详情
git log : 查看所有的commit提交记录
git show [commit_id] : 查看最新或指定的commit
git show commitId fileName # 查看某次commit中具体某个文件的修改

### [迁移git commit](https://git-scm.com/book/zh/v2/Git-%E5%B7%A5%E5%85%B7-%E6%89%93%E5%8C%85)

### git拉取远程分支
```bash
# git fetch [origin <remote-branch>] # 刷新本地的远程分支列表, 避免`git branch -a`找不到远程分支
## 可选择是否合并到local current branch
## - 需要合并到local current branch, 假设local current branch是master
# git pull <远程repo名> <远程分支名>:<本地分支名>
# git pull origin xxx = git pull origin xxx:master # git pull **会合并到本地分支, 且本地分支必须存在**
## - 不需要合并到local current branch
# git checkout -b xxx orgin/xxx
```

### Git cherry-pick `from another repository`/`cross repository`
参考:
- [使用 Git 時如何做出跨 repo 的 cherry-pick](https://blog.m157q.tw/posts/2017/12/30/git-cross-repo-cherry-pick/)

> git cherry-pick会改变commit id.


共计3种方法:
1. git remote add + git fetch + git cherry-pick(**推荐**)
```bash
# Cloning our fork
$ git clone git clone git@github.com:ifad/rest-client.git

# Adding (as "endel") the repo from we want to cherry-pick
$ git remote add endel git://github.com/endel/rest-client.git

# Fetch their branches
$ git fetch endel

# List their commits
$ git log endel/master

# Cherry-pick the commit we need
$ git cherry-pick 97fedac

# Pushing to our master
$ git push origin master
```

1. git format-patch + git am(**推荐**)
```bash
# git format-patch -k --stdout ${A_COMMIT_HASH}..${B_COMMIT_HASH} > xxx.patch # `-k`, 不加时patch中的git log's subject会被冠以"[PATCH] "前缀, 不要省略commit range中间的"..", rang = (A_COMMIT_HASH, B_COMMIT_HASH]
# git am -k -3 [--signoff] < xxx.patch # "-3",使用 three-way merge; "-k", 去除导入patch中git log's subject中的"[PATCH] "前缀
```

1. git diff + git apply
```bash
# git diff ${A_COMMIT_HASH} ${B_COMMIT_HASH} > xxx.patch # 会丢失git log的subject
# git apply xxx.patch
```

### git pull origin master 如何不自动生成 merge commit
```bash
# [git pull --rebase的正确使用](https://juejin.im/post/5d3685146fb9a07ed064f11b)和[聊下git pull --rebase](https://www.cnblogs.com/wangiqngpei557/p/6056624.html)
$ git stash
$ git pull --rebase
$ git push
$ git stash pop
```

### 如何在 git 里查找哪一个 commit 删除了代码的一行
`git log -S <string> path/to/file`: 使用git log 的 -S 选项来指定 commit 里包括的字符串.

### 查看某行代码由谁写的，在哪个commit中提交的
`git blame file_name`

### 查看哪个Git分支正在跟踪哪个远程/上游分支
```bash
# git branch -vv   # doubly verbose
# git branch --set-upstream-to=<remote>/<branch> <local-branch> # git pull/push upstream
# git push --set-upstream <remote> <remote-branch> # 为当前分支设在upstream
# git pull --set-upstream <remote> <remote-branch> # 为当前分支设在upstream, 作用同上
```

### git clone断点续传
没有, [网上找的解决方法](https://gist.github.com/arliang/0019de079fbe77f946c13a010e7f97c6)无用.

### git pull拉取远程分支合并到本地分支
git pull <远程主机名> <远程分支名>:<本地分支名> , 比如`git pull origin master:wy`

### git checkout远程分支到本地
`git checkout -b 本地分支名x origin/远程分支名x`

### [配置git diff的输出颜色](http://ericnode.info/post/colorize_git_diff/), 以区别terminal的背景色
可执行`git help config`，在里面搜找`color.diff`相关配置的说明.

比如:
```conf
[color "diff"]
    meta = white reverse
    frag = cyan reverse
    old = red reverse
    new = green reverse
```

其中修改了四个slot的颜色属性:
- meta: meta information，分割了不同的文件。设置为white reverse
- frag: hunk header, 文件中一个修改的头，分割了同一个文件内的不同修改。设置为cyan reverse
- old: 被删除的代码，设置为red bold
- new: 新增的代码，设置为green bold

其实`color.diff.<slot>`支持的slot有:
- plain
- meta
- frag
- func
- old
- new
- commit
- whitespace

### error: 推送一些引用到 'git@gitee.com:chenhao/hello_minio.git' 失败
```log
# git push -u read  master
...
To gitee.com:chenhao/hello_minio.git
 ! [remote rejected] master -> master (shallow update not allowed)
error: 推送一些引用到 'git@gitee.com:chenhao/hello_minio.git' 失败
```

在 clone 原仓库时用了`git clone --depth 1`，导致本地为`shallow repo`, 解决方法也很简单:
1. 补全repo再push

    ```
    git fetch --unshallow origin
    ```
1. 删除`.git`, 再重建repo并push

### 显示git调用curl过程中的详细信息
`GIT_CURL_VERBOSE=1 git ls-remote https://github.com/`

`echo | openssl s_client -connect github.com:443`

### 查看git使用openssl还是guntls
`/usr/lib/git-core/ldd git-http-fetch | grep libcurl`

### go get报"Error -50 setting GnuTLS cipher list starting with +VERS-TLS1.3:+SRP:
git依赖的guntls不支持tls 1.3, 让go get使用git ssh即可: `git config --global url.git@github.com:.insteadOf https://github.com/`

### git remote test
`ssh -T git@github.com`

### git submodule update --init --recursive
`git submodule init` + `git submodule update`

### `git submodule add git@gitee.com:chenhao/hello_zstack.git __read_source`报`fatal: You are on a branch yet to be born`
hello_zstack.git是全新repo, 必须有git log(即有内容)才行.

### git只克隆仓库某个目录
```bash
# mkdir devops
# cd devops/
# git init                  #初始化空库

## step 2 :  拉取remote的all objects信息
# git remote add -f origin git@github.com:gopherchina/conference.git   #拉取remote的all objects信息


## step 3 :  #3.1 开启sparse clone, #3.2 设置需要pull的目录 devlops
# git config core.sparsecheckout true   #开启sparse clone
# echo "devops" >> .git/info/sparse-checkout   #设置需要pull的目录，*表示所有，!表示匹配相反的
# more .git/info/sparse-checkout

## step 4 :  # 将origin 端，由第三步（文件 .git/info/sparse-checkout）设置的 目录下的文件 pull 到本地
# git pull --depth 1 origin master
```

### git log filter
```
# git log --author="meilihao" --reverse --pretty="%H" |xargs -I {} bash -c "git show {} --stat --pretty='%b'|grep src"
# git log --since="2019-02-20" --until='2012-12-10'
# git log --after="2019-02-20"
```

### [git筛选commit制作patch](https://stackoverflow.com/questions/14591888/apply-patches-created-with-git-log-p?rq=1)
`git log --reverse --author="John\|Mary" --since "Fri Jul 23 09:38:07 2021 +0800" --no-merges --pretty="%H" -- <path>`:
- `--reverse` : git log本身输出是倒叙的
- `-since` : 指定大于等于某个时刻
- `--no-merges` : 不包含merge commit
- `-- <path>` : 仅包含指定的路径
- --pretty=<format> : 格式化commit信息: `%H`是提交对象（commit）的完整哈希字串


因此完整命令是即`git format-patch $(git log --reverse --author="John Doe" --pretty="%H" -- the/needed/path | awk '{print -NR, " ", $0}') -o patches`:
- awk用于生成patch的序号
- `-o patches` : 所有生成的patch保存在patches文件夹中

实际操作发现git log生成的结果有多条时, `git format-patch $(git log ...)`生成的patch不正确, 因此改为script形式:
```bash
#!/usr/bin/env bash
COUNTER=1;
for c in `git log --reverse --author="John Doe" --pretty="%H" -- the/needed/path`;do
    git format-patch --start-number=$COUNTER "$c^1".."$c" -o patches # 使用`"$c^1"`因为不包含起点
    let COUNTER=COUNTER+1
done
```

> 上述script使用`git log xxx | awk '{print -NR, " ", $0}'` + `git format-patch $c -o patches`时会因为`$c`的内的换行导致执行结果与实际不符.

### 查看commit的文件变动
```bash
git show <commit> --stat # 使用`--stat`时可能因输出内容过宽导致出现`...`替换从而导致后接的grep不匹配
git show <commit> --name-only # 对merge无效
git show <commit> --name-status # 对merge无效
git show 6db185d0118853b9382f6550b3741f13557450a1 --stat --pretty="%b"
```

### 迁移repo
```bash
git remote add new_origin xxx
git push -u new_origin --all
git push -u new_origin --tags
```

### [`git push`报`shallow update not allowed`](https://stackoverflow.com/questions/28983842/remote-rejected-shallow-update-not-allowed-after-changing-git-remote-url)
解决方法: `git fetch --unshallow origin && git push`

### commitizen-go
ref:
- [规范你的代码 -Commitizen](https://juejin.cn/post/7024103006752735269)
- [git commit规范及自动检查工具安装小记](https://juejin.cn/post/6844904033635794958)
- [约定式提交](https://www.conventionalcommits.org/zh-hans/v1.0.0/)

```bash
$ go install github.com/lintingzhen/commitizen-go@master
$ whereis commitizen-go
/home/chen/git/go/bin/commitizen-go
$ sudo /home/chen/git/go/bin/commitizen-go install # 需要在git repo里执行, 此后即可在git repo里使用`git cz`
$ git cz
```

node版本:
```bash
npm install -g commitizen cz-conventional-changelog
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc
cz
```

### git pull报`无法找到远程引用 refs/heads/xxx`
原因: 以前重命名过分支

解决方法: `vim .git/config`, 将其中错误的分支名替换为正确的分支名.

### 访问github报`kex_exchange_identification: Connection closed by remote host`
原因: 被封禁了 Github 端口 22 的连接

```bash
# vim ~/.ssh/config
Host github.com
    HostName ssh.github.com
    User git
    Port 443
# ssh -T git@github.com -v
```

### 比较分支文件差异
```bash
git diff branch1 branch2 [--stat] # 显示差异. `--stat`: 仅显示差异的文件列表
git diff branch1 branch2 <file_path> # 比较指定文件差异
```

### 将修改文件发送到remote+删除旧文件
`git diff --name-only [--relative=xxx] [--cached] [--diff-filter=AM] | xargs -I '{}' scp -i ~/.ssh/xxx '{}' root@192.168.16.100:/opt/xxx/{}`
`git diff --name-only [--relative=xxx] [--cached] [--diff-filter=AM] | xargs -I {} echo {}|sed 's/\.py/\.so' | xargs -I {} ssh -i ~/.ssh/xxx root@192.168.16.100 "rm /opt/{}"`

选项:
- relative : 调整相对路径

> diff-filter from `git diff --name-status`

### clone带submodule的repo
`git clone --recursive https://github.com/cloudflare/quiche`

### 设置默认remote
`git branch --set-upstream-to origin/master`

### git push无响应且`ssh -T git@github.com -vvv`卡在`expecting SSH2_MSG_KEX_ECDH_REPLY`
ref:
- [SSH使用问题以及解决方案(expecting SSH2_MSG_KEX_ECDH_REPLY) ](https://github.com/johnnian/Blog/issues/44)

git push之前正常, 2024.3月末出现该问题.

ifconfig tun0 mtu 1400

### gitcode git push报`Deny by project hooks setting 'default': size of the file...`
ref:
- [gitcode 上传文件报错文件太大has exceeded the upper limited size](https://blog.csdn.net/downanddusk/article/details/138187765)

    验证修改无效
- [处理 github 不允许上传超过 100MB 文件的问题](https://www.liuxiao.org/2017/02/git-%E5%A4%84%E7%90%86-github-%E4%B8%8D%E5%85%81%E8%AE%B8%E4%B8%8A%E4%BC%A0%E8%B6%85%E8%BF%87-100mb-%E6%96%87%E4%BB%B6%E7%9A%84%E9%97%AE%E9%A2%98/)

gitcode限制10M, 解决方法:
```bash
git rm --cached [-r] path_of_a_giant_file
git commit --amend
git push
```

### tag version
version.xxx时全局变量, 它将在编译时被设置

```bash
# cd <project>
# cat pkg/lib/version/version.go
package version

import (
	"fmt"
	"runtime"
)

var (
	gitTag    string
	gitBranch string
	gitHash   string
	buildTime string
)

type VersionInfo struct {
	GitTag    string
	GitBranch string
	GitHash   string
	BuildTime string
	GoVersion string
	Compiler  string
	Platform  string
}

func GetVersion() *VersionInfo {
	return &VersionInfo{
		GitTag:    gitTag,
		GitBranch: gitBranch,
		GitHash:   gitHash,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s", v.GitBranch, v.GitTag, v.GitHash, v.BuildTime, v.GoVersion, v.Platform, v.Compiler)
}
```

```bash
GitTag=$(git describe --tags --dirty --always)
GitBranch=$(git rev-parse --abbrev-ref HEAD)
GitHash=$(git rev-parse HEAD)
# BuildTS=$(date -u --rfc-3339=seconds)
BuildTS=$(date -u +'%Y-%m-%dT%H:%M:%S')

LDFLAGS="-X <project>/pkg/lib/version.gitTag=${GitTag}
         -X <project>/pkg/lib/version.gitBranch=${GitBranch}
         -X <project>/pkg/lib/version.gitHash=${GitHash}
         -X '<project>/pkg/lib/version.buildTime=${BuildTS}'"

go build -ldflags "$LDFLAGS"
```

### git pull报`cannot lock ref`
```
# git pull
error: cannot lock ref 'refs/remotes/origin/feature': 'refs/remotes/origin/feature/xdd-manager-backend' exists; cannot create 'refs/remotes/origin/feature'
From codeup.aliyun.com:xxx/server
 ! [new branch]        feature    -> origin/feature  (unable to update local ref)
```

手动删除`refs/remotes/origin/feature/xdd-manager-backend`没效果

解决: `git remote prune origin`