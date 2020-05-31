# 安装

通过官方的安装脚本安装.

[安装过程]():
```bash
$ curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > rust.sh
$ export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static
$ export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup
$ bash rust.sh
$ echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc # from $HOME/.cargo/env
```

> [rust 使用国内镜像，快速安装方法](https://www.cnblogs.com/hustcpp/p/12341098.html)

## rustup

参考:
- [Rust 版本管理工具: rustup](https://github.com/rustcc/RustPrimer/blob/master/install/rustup.md)

### 获取源码
```
rustup component add rust-src
```

### 获取rust crates mirror
[Rust Crates 镜像使用帮助](https://lug.ustc.edu.cn/wiki/mirrors/help/rust-crates)

## tool
### sublime3插件
```
Rust Enhanced // 需先通过package control的disable package禁用st3自带的Rust插件
RustAutoComplete
```
