# 安装

通过官方的安装脚本安装.

[安装过程]():
```bash
$ curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > rust.sh
$ export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static # 用于更新 toolchain
$ export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup # 用于更新 rustup
$ bash rust.sh
$ echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc # from $HOME/.cargo/env
```

> [rust 使用国内镜像，快速安装方法](https://www.cnblogs.com/hustcpp/p/12341098.html)
> [rust ustc mirror](https://lug.ustc.edu.cn/wiki/mirrors/help/rust-static)
> [rust tsinghua mirror](https://mirrors.tuna.tsinghua.edu.cn/help/rustup/)

## rustup

参考:
- [Rust 版本管理工具: rustup](https://github.com/rustcc/RustPrimer/blob/master/install/rustup.md)

## 安装rls
`rustup component add rls rust-analysis rust-src`

## 常用命令
- 使用 rustc -h : rustc 的基本用法
- 使用 cargo -h : cargo 的基本用法
- 使用 rustc -C help : rustc 的一些跟代码生成相关的选项
- 使用 rustc -W help : rustc 的一些跟代码警告相关的选项
- 使用 rustc -Z help : rus tc 的一些跟编译器内部实现相关的选项
- 使用 r ustc -help -V : rustc 的更详细的选项说明

### 获取源码
```
rustup component add rust-src
```

### 获取rust crates mirror
[Rust Crates 镜像使用帮助](https://lug.ustc.edu.cn/wiki/mirrors/help/rust-crates)

## tool
参考:
- [将 Vim 设置为 Rust IDE](https://linux.cn/article-12530-1.html)

### sublime3插件
```
Rust Enhanced // 需先通过package control的disable package禁用st3自带的Rust插件
RustAutoComplete
```
