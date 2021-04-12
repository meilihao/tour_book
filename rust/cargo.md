# cargo
Cargo是Rust中的包管理工具，第三方包叫做crate.

Cargo一共做了四件事：
- 使用两个元数据（metadata）文件来记录各种项目信息
- 获取并构建项目的依赖关系
- 使用正确的参数调用rustc或其他构建工具来构建项目
- 为Rust生态系统开发建议了统一标准的工作流

Cargo文件：
- Cargo.lock：只记录依赖包的详细信息，不需要开发者维护，而是由Cargo自动维护
- Cargo.toml：描述项目所需要的各种信息，包括第三方包的依赖

cargo编译默认为Debug模式, 在该模式下编译器不会对代码进行任何优化; 也可以使用--release参数来使用发布模式, 此时编译器会对代码进行优化，使得编译时间变慢，但是代码运行速度会变快.

官方编译器rustc，负责将rust源码编译为可执行的文件或其他文件`.a, .so, .lib`等.

Rust还提供了包管理器Cargo来管理整个工作流程:
```bash
cargo new first_pro_create # 创建一个编写可执行文件的项目 = `cargo new --bin first_pro_create`
cargo new --lib first_lib_create # 创建用于编写库的项目
cargo doc
cargo doc --open
cargo test
cargo test -- --test-threads=1
cargo build # 编译项目
cargo build --release
cargo run # 运行项目
cargo install --path
cargo uninstall first_pro_create
```

## 代码组织
Rust的代码从逻辑上是分 crate 和 mod 管理的. crate 可以理解为"项目". 每个 crate 是一个完整的编译单元，它可以生成为一个 lib 或者可执行文件. 而在 crate 内部则由 mod 管理, mod 大家可以理解为 namespace. 可以使用 use 语句把其他模块中的内容引入到当前模块中来.

Rust 有一个极简标准库， 叫作 std ，除了极少数嵌入式系统下无法使用标准库之外，绝大部分情况下，我们都需要用到标准库里面的东西. 为了给大家减少麻烦， Rust 编译器对标准库有特殊处理. 默认情况下，用户不需要手动添加对标准库的依赖 ，编译器会自动引人对标准库的依赖. 除此之外 ，标准库中的某些 type, trait、 function, macro 等实在是太常用了, 每次都写 use 语句确实非常无聊，因此标准库提供了一个[`std::prelude`'](https://github.com/rust-lang/rust/blob/master/library/std/src/prelude/mod.rs)模块，在这个模块中导出了一些最常见的类型, trait 等东西, 编译器会为用户写的每个 crate 自动插入一句话：`use std: :prelude::*;`, 这样，标准库里面的这些最重要的类型, trait 等名字就可以直接使用，而无须每次都写全称或者 use 语句.

> [std::prelude](https://doc.rust-lang.org/std/prelude/)目前的[mod.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/mod.rs) 中, 直接导出了 v1 模块中的内容, 而 [v1.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/v1.rs) 中则是编译器为我们自动导人的相关trait和类型.

### 使用第三方包
Rust可以在Cargo.toml中的[dependencies]下添加想依赖的包来使用第三方包. 然后在`src/main.rs`或`src/lib.rs`文件中, 使用`extern crate`命令声明引入该包即可使用.

值得注意的是, 使用extern crate声明包的名称是linked_list，用的是下划线`_`, 而在Cargo.toml中用的是连字符`-`, 其实Cargo默认会把连字符转换成下划线.

Rust也不建议以`-rs`或`_rs`为后缀来命名包名, 而且会强制性的将此后缀去掉.