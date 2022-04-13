
## git config init
```
# ~/.gitconfig
[core]
	editor = vim
	quotepath = false
	filemode = false
[user]
	name = meilihao
	email = 563278383@qq.com
[alias]
	br = branch
	ci = commit
	cm = commit -m
	co = checkout
	df = diff
	d = diff
	dv = difftool -t vimdiff -y
	lg = log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
	ll = log --oneline
	st = status
	st2 = status -sb
	last = log -1 HEAD --stat
	rv = remote -v
	gl = config --global -l
	se = !git rev-list --all | xargs git grep -F
[color "diff"]
	meta = white reverse
	frag = cyan reverse
	old = red reverse
	new = green reverse
[url "git@github.com:"]
	insteadOf = git://github.com/
[url "https://github.com/"]
	insteadOf = git://github.com/
```

### prompt/提示
1. 下载git prompt
下载`github.com/git/git/tree/master/contrib/completion`的`git-completion.bash`和`git-prompt.sh`到`~/.config`
 
1. 把下面代码加入~/.bashrc中
```
# git
[[ $- == *i* ]] && . /home/chen/.config/git-prompt.sh
[[ $- == *i* ]] && . /home/chen/.config/git-completion.bash
PS1='\[\033[1;32m\]\u \[\033[1;34m\]\W\[\033[1;33m\]$(__git_ps1 " (%s)")\[\033[0m\] \$ '
```
> [PS1中使用的函数](https://gist.github.com/richarddong/1981392)

1. 使修改生效

### 默认参数

- `git config --global core.editor vim`,修改git默认编辑器