
### prompt/提示

```
# 下载git prompt
github.com/lvv/git-prompt/git-prompt.sh
 
# 把下面代码加入~/.bashrc中
[[ $- == *i* ]] && . ~/github/git-prompt/git-prompt.sh
 
# 使修改生效
$ source ~/.bashrc
```

### 默认参数

- `git config --global core.editor vim`,修改git默认编辑器