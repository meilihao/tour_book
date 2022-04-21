# 安装
参考:
- [rust离线安装](https://forge.rust-lang.org/infra/other-installation-methods.html)

通过官方的安装脚本安装.

[安装过程](https://www.rust-lang.org/zh-CN/tools/install):
```bash
$ curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > rust.sh
$ export RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static # 用于更新 toolchain
$ export RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup # 用于更新 rustup
$ bash rust.sh
$ echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc # from $HOME/.cargo/env
$ rustc --version
```

> [rust 使用国内镜像，快速安装方法](https://www.cnblogs.com/hustcpp/p/12341098.html)
> [rust ustc mirror](https://lug.ustc.edu.cn/wiki/mirrors/help/rust-static)
> [rust tsinghua mirror](https://mirrors.tuna.tsinghua.edu.cn/help/rustup/)

## rustup

参考:
- [Rust 版本管理工具: rustup](https://github.com/rustcc/RustPrimer/blob/master/install/rustup.md)

## 安装rls
`rustup component add rls rust-analysis rust-src`

修改rustup默认安装位置, 配置RUSTUP_HOME(默认`~/.rustup`, 保存工具链和配置文件)和CARGO_HOME(默认`~/.cargo`,保存cargo的cache)即可.

## 常用命令
- `rustc -V` : rust version
- 使用 rustc -h : rustc 的基本用法
- 使用 cargo -h : cargo 的基本用法
- 使用 rustc -C help : rustc 的一些跟代码生成相关的选项
- 使用 rustc -W help : rustc 的一些跟代码警告相关的选项
- 使用 rustc -Z help : rus tc 的一些跟编译器内部实现相关的选项
- 使用 rustc -help -V : rustc 的更详细的选项说明

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

### vscode插件
- rust-analyzer: 它会实时编译和分析 Rust 代码，提示代码中的错误，并对类型进行标注; 也可以使用官方的 rust 插件取代
- rust syntax：为代码提供语法高亮
- crates：帮助分析当前项目的依赖是否是最新的版本
- better toml：Rust 使用 toml 做项目的配置管理. better toml 可以帮你语法高亮, 并展示 toml 文件中的错误
- rust test lens：可以帮快速运行某个 Rust 测试
- Tabnine：基于 AI 的自动补全，可以帮助开发者更快地撰写代码