# bashrc
env: xubuntu 20.04/deepin v20 

1. 先安装[oh-my-bash](https://github.com/ohmybash/oh-my-bash)

1. `.bashrc`
```
...
# --- alias
alias ls="ls --color"
alias vim="nvim"
alias grep="grep --color=auto" # [设置grep高亮显示匹配项和基本用法](https://www.cnblogs.com/lazyfang/p/7645627.html)

alias python="/usr/bin/python3.8"

# --- fcitx
export GTK_IM_MODULE="fcitx"
export QT_IM_MODULE="fcitx"
export XMODIFIERS="@im=fcitx"

# --- go
export GOROOT=/usr/local/go
export GOPATH=/home/chen/git/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# --- rust
export PATH=$HOME/.cargo/bin:$PATH
export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static # 用于更新 toolchain
export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup # 用于更新 rustup
```