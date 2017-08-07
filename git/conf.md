
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