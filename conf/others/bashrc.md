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
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOROOT=/usr/local/go
export GOPATH=/home/chen/git/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# --- rust
export RUSTUP_HOME=/opt/rust/rustup # default: ~/.rustup
export CARGO_HOME=/opt/rust/cargo   # default: ~/.cargo
source $CARGO_HOME/env
# https://lug.ustc.edu.cn/wiki/mirrors/help/rust-static/
export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static # 用于更新 toolchain
export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup # 用于更新 rustup

# --- git
export LESSCHARSET=utf-8 # git diff中文乱码

# --- llvm
alias clang="clang-13"
alias opt="opt-13"
alias llvm-dis="llvm-dis-13"
alias llvm-as="llvm-as-13"
alias llvm-link="llvm-link-13"
alias llvm-mc="llvm-mc-13"
alias lli="lli-13"
alias llc="llc-13"
alias ncdu="ncdu --exclude-kernfs "

# --- liteide
export LD_LIBRARY_PATH="/opt/liteide/lib:$LD_LIBRARY_PATH"
```