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
> [字节跳动新的 Rust 镜像源](https://rsproxy.cn/)
> [rust ustc mirror](https://lug.ustc.edu.cn/wiki/mirrors/help/rust-static)
> [rust tsinghua mirror](https://mirrors.tuna.tsinghua.edu.cn/help/rustup/)

## rustup

参考:
- [Rust 版本管理工具: rustup](https://github.com/rustcc/RustPrimer/blob/master/install/rustup.md)

rust有stable和nightly两个版本, 默认版本可通过`rustup show`查看.

修改rustup默认安装位置, 配置RUSTUP_HOME(默认`~/.rustup`, 保存工具链和配置文件)和CARGO_HOME(默认`~/.cargo`,保存cargo的cache)即可. 卸载rust用`rustup self uninstall`.

## nightly
```bash
# rustup install nightly
# rustup default nightly
# rustup toolchain remove stable
```

## 安装rls
`rustup component add rls rust-analysis rust-src`

## toolchians(工具链)
ref:
- [Clippy：Rust 官方代码质量增强工具完全指南](https://mp.weixin.qq.com/s/y_YSZKp9sNd2JoKQ-xlTEQ)

- rustc : 编译器
- cargo : 项目管理工具
- rustup : 管理工具

	- `rustup toolchain list` : 查看支持的toolchian
	- `rustup default <stable|nightly>` : 切换rust版本
	- `rustup show`: 查看默认设置
	- `rustup component add clippy` # 用`cargo clipy`格式化rust代码, 类似eslint

		```bash
		# cargo clippy # 项目目录下运行
		# cargo clippy --tests # 针对特定目标（如测试代码）
		```
	- `rustup update nightly` && ` rustup override set nightly` : 为项目指定使用的rust版本 
- rustdoc : 文档工具
- rustfmt : 格式化工具
- rust-gdb : 调试工具
- rust-lldb: 调试工具

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

## rls
```bash
rustup component add rls --toolchain 
rustup component add rust-analysis --toolchain 
rustup component add rust-src --toolchain
```