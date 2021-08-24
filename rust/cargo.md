# cargo
Cargo是Rust中的包管理工具，第三方包叫做crate.

Cargo一共做了四件事：
- 使用两个元数据（metadata）文件来记录各种项目信息
- 获取并构建项目的依赖关系
- 使用正确的参数调用rustc或其他构建工具来构建项目
- 为Rust生态系统开发建议了统一标准的工作流

Cargo文件：
- Cargo.lock：只记录依赖包的详细信息，不需要开发者维护，而是由Cargo自动维护

    Cargo 用来确保任何人在任何时候重新构建代码，都会产生相同的结果.

    它在第一次运行 cargo build 时创建. 当第一次构建项目时，Cargo 计算出所有符合要求的依赖版本并写入 Cargo.lock 文件. 当将来构建项目时，Cargo 会发现 Cargo.lock 已存在并使用其中指定的版本，而不是再次计算所有的版本. 这使得我们拥有了一个自动化的可重现的构建.

    Cargo 提供了另一个命令，update，它会忽略 Cargo.lock 文件，并计算出所有符合 Cargo.toml 声明的最新版本。如果成功了，Cargo 会把这些版本写入 Cargo.lock 文件.
- Cargo.toml：描述项目所需要的各种信息，包括第三方包的依赖

cargo编译默认为Debug模式, 在该模式下编译器不会对代码进行任何优化; 也可以使用--release参数来使用发布模式, 此时编译器会对代码进行优化，使得编译时间变慢，但是代码运行速度会变快.

官方编译器rustc，负责将rust源码编译为可执行的文件或其他文件`.a, .so, .lib`等.

Rust还提供了包管理器Cargo来管理整个工作流程:
```bash
cargo new first_pro_create # 创建一个编写可执行文件的项目 = `cargo new --bin first_pro_create`
cargo new --lib first_lib_create # 创建用于编写库的项目
cargo doc
cargo doc --open # 构建所有本地依赖提供的文档，并在浏览器中打开
cargo check # 快速检查当前代码是否可以通过编译，而不需要构建程序, 比cargo build快. 因此大部分Rust用户在编写程序的过程中都会周期性地调用cargo check以保证自己的程序可以通过编译，只有真正需要生成可执行程序时才会调用.
cargo build
外的时间去真正生成可执行程序：
cargo test
cargo test -- --test-threads=1
cargo build # 编译项目. 用于开发，它允许快速地反复执行构建操作
cargo build --release # 这种模式会以更长的编译时间为代价来优化代码, 从而使代码拥 有更好的运行时性能
cargo run # 运行项目
cargo install --path
cargo uninstall first_pro_create
```

## 代码组织
Rust 有许多功能可用于组织和管理代码, 包括哪些内容可以被公开, 哪些内容作为私有部分, 以及程序每个作用域中的名字. 这些功能有时被称为`模块系统(the module system)`包括:
- 包（Packages） : Cargo 的一个功能，它允许构建、测试和分享 crate

    一个包会包含有一个 Cargo.toml 文件，阐述如何去构建这些 crate.

    包中所包含的内容由几条规则来确立:
    1. 一个包中至多只能包含一个库 crate(library crate)
    1. 包中可以包含任意多个二进制 crate(binary crate)
    1. 包中至少包含一个 crate，无论是库的还是二进制的

    执行`cargo new my-project`, 查看 Cargo.toml 的内容，会发现并没有提到 src/main.rs, 因为 Cargo 遵循的一个约定: src/main.rs 就是一个与包同名(这里就是my-project)的二进制 crate 的 crate 根. 同样的, Cargo 知道如果包目录中包含 src/lib.rs, 则 src/lib.rs 就是一个与包同名的 lib crate的 crate 根. crate 根文件将由 Cargo 传递给 rustc 来实际构建库或者二进制项目.

    如果一个包同时含有 src/main.rs 和 src/lib.rs，则它有两个 crate：一个库和一个二进制项，且名字都与包相同. 通过将文件放在 src/bin 目录下，一个包可以拥有多个二进制 crate：每个 src/bin 下的文件都会被编译成一个独立的二进制 crate.

- Crates : 一个模块的树形结构，它形成了库或二进制项目

    crate root 是一个源文件, Rust 编译器以它为起始点, 同时也是 crate 的根模块.
- 模块（Modules）和 use : 允许控制作用域和路径的私有性
- 路径（path）: 一个命名例如结构体、函数或模块等项的方式

对于一个由一系列相互关联的包组合而成的超大型项目, Cargo 提供了`工作空间`这一解决方案.

> crate是rust最小的编译单元, package是若干crate的集合, 它们都可被称为包. 只在两者同时出现且需要区别对
待时，将crate译为单元包，将package译为包.

Rust的代码从逻辑上是分 crate 和 mod 管理的. crate 可以理解为"项目". 每个 crate 是一个完整的编译单元，它可以生成为一个 lib 或者可执行文件. 而在 crate 内部则由 mod 管理, mod 大家可以理解为 namespace. 可以使用 use 语句把其他模块中的内容引入到当前模块中来.

Rust 有一个极简标准库， 叫作 std ，除了极少数嵌入式系统下无法使用标准库之外，绝大部分情况下，我们都需要用到标准库里面的东西. 为了给大家减少麻烦， Rust 编译器对标准库有特殊处理. 默认情况下，用户不需要手动添加对标准库的依赖 ，编译器会自动引人对标准库的依赖. 除此之外 ，标准库中的某些 type, trait、 function, macro 等实在是太常用了, 每次都写 use 语句确实非常无聊，因此标准库提供了一个[`std::prelude`'](https://github.com/rust-lang/rust/blob/master/library/std/src/prelude/mod.rs)模块(doc在[这](https://doc.rust-lang.org/std/prelude/index.html))，在这个模块中导出了一些最常见的类型, trait 等东西, 编译器会为用户写的每个 crate 自动插入一句话：`use std: :prelude::*;`, 这样，标准库里面的这些最重要的类型, trait 等名字就可以直接使用，而无须每次都写全称或者 use 语句. 如果需要的类型不在 prelude 中，你必须使用 use 语句显式地将其引入作用域.

> [std::prelude](https://doc.rust-lang.org/std/prelude/)目前的[mod.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/mod.rs) 中, 直接导出了 v1 模块中的内容, 而 [v1.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/v1.rs) 中则是编译器为我们自动导人的相关trait和类型.

### 使用第三方包
Rust可以在Cargo.toml中的[dependencies]下添加想依赖的包来使用第三方包. 然后在`src/main.rs`或`src/lib.rs`文件中, 使用`extern crate`命令声明引入该包即可使用.

值得注意的是, 使用extern crate声明包的名称是linked_list，用的是下划线`_`, 而在Cargo.toml中用的是连字符`-`, 其实Cargo默认会把连字符转换成下划线.

Rust也不建议以`-rs`或`_rs`为后缀来命名包名, 而且会强制性的将此后缀去掉.