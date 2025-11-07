# cargo
ref:
- [Rust Cargo 官书（非官方翻译)](https://llever.com/cargo-book-zh/index.zh.html)

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
cargo new first_pro_create [--vcs none] # 创建一个编写可执行文件的项目 = `cargo new --bin first_pro_create`. vcs默认是git
cargo new --lib first_lib_create # 创建用于编写库的项目
cargo doc
cargo doc --open --no-deps # 构建文档(在`target/doc`)，并在浏览器中打开. `--no-deps`即忽略生成依赖项的文档
cargo check # 快速检查当前代码是否可以通过编译，而不需要构建程序, 比cargo build快. 因此大部分Rust用户在编写程序的过程中都会周期性地调用cargo check以保证自己的程序可以通过编译，只有真正需要生成可执行程序时才会调用.
cargo build
外的时间去真正生成可执行程序：
cargo test
cargo test -- --test-threads=1
cargo check # 它基本上跳过了编译器的代码生成部分，只通过前端阶段运行代码, 即编译器的解析和语义分析
cargo build # 编译项目. 用于开发，它允许快速地反复执行构建操作. 默认情况下，cargo build 编译出来的二进制，在项目根目录的 target/debug 下.
cargo build --release # 这种模式会以更长的编译时间为代价来优化代码, 从而使代码拥有更好的运行时性能, 但不会显示 panic backtrace 的具体行号; 不会检查整数溢出; 会跳过 debug_assert!() 断言
cargo run [--release] # 运行项目
cargo run --examples <file_name> #  examples/目录下的代码
cargo run -- assets/ferris.png # `--`表示将之后的内容传递给当前项目的可执行程序
cargo install [--path] cargo-watch # 安装cargo自定义命令cargo-watch, 安装后在cargo项目下可用`cargo watch`命令. 其他命令cargo-edit(自动添加依赖到Cargo.toml), cargo-deb(构建deb软件包), cargo-outdated(显示过时依赖), racer(代码补全)
cargo uninstall first_pro_create
cargo fmt: 类似gofmt, 格式化代码
cargo clippy: 类似eslint, 检查代码规范
cargo tree: 查看第三方库的版本和依赖关系
cargo bench: 运行benchmark(基准测试,性能测试)
cargo udeps(第三方): 检查项目中未使用的依赖
cargo update [-p <crate>] : 更新全部或某个依赖. 默认只更新小版本, 比如`0.3.14->0.3.20`, 大版本升级需修改Cargo.toml
cargo package : 创建一个文件(`target/package/xxx-<version>.crate`）, 其中包含库的所有源文件，以及 Cargo.toml
cargo package --list : 查看其中包含什么文件
cargo login <apikey> : 登录crates.io
cargo publish : 将包发布到 crates.io
```

## 代码组织
Rust 有许多功能可用于组织和管理代码, 包括哪些内容可以被公开, 哪些内容作为私有部分, 以及程序每个作用域中的名字. 这些功能有时被称为`模块系统(the module system)`包括:
- workspace：项目复杂时，管理多个package

    工作区只是一个包含 Cargo.toml 文件的目录. 整个workspace共享一个Cargo.lock，也只有一个target目录

    使用 Cargo workspace可以节省编译时间和磁盘空间.

- 包（Packages） : **Cargo 的一个功能概念**, 用于管理crate, 比如允许构建、测试和分享 crate

    **一个包会包含有一个 Cargo.toml 文件**，描述了包的基本信息以及依赖项以及如何去构建这些 crate.

    > Cargo.toml使用semver版本语义

    包中所包含的内容由几条规则来确立:
    1. 一个包中至多只能包含一个库 crate(library crate)
    1. 包中可以包含任意多个二进制 crate(binary crate)
    1. 包中至少包含一个 crate，无论是库的还是二进制的

    执行`cargo new my-project`, 查看 Cargo.toml 的内容，会发现并没有提到 src/main.rs, 因为 Cargo 遵循的一个约定: src/main.rs 就是一个与包同名(这里就是my-project)的二进制 crate 的 crate 根. 同样的, Cargo 知道如果包目录中包含 src/lib.rs, 则 src/lib.rs 就是一个与包同名的 lib crate的 crate 根. crate 根文件将由 Cargo 传递给 rustc 来实际构建库或者二进制项目.

    如果一个包同时含有 src/main.rs 和 src/lib.rs，则它有两个 crate：一个库和一个二进制项，且名字都与包相同. 通过将文件放在 src/bin 目录下，一个包可以拥有多个二进制 crate：每个 src/bin 下的文件都会被编译成一个独立的二进制 crate.

- crate : mod的集合, 一个模块的树形结构，它形成了库或二进制项目, 存在于`包`中.

    crate root 是一个源文件(main.rs/lib.rs), Rust 编译器以它为起始点, 同时也是 crate 的根模块, 其名与package名相同.

    > 只是如果有多个bin类型的crate，一个src/main.rs就不行了，就得放到 src/bin 下面，每个crate一个文件，换句话说，每个文件都是一个不同的crate

    > `extern crate iron`: 将 Cargo.toml 文件中指定的 iron 引入程序中

    > `#[macro_use] extern crate mime;` : 将 Cargo.toml 文件中指定的 mine 引入程序中, 并会使用mine导出的宏

- 模块（Modules即mod）和 use : 允许控制作用域和路径的私有性

    模块既是 Rust 的命名空间，也是函数、类型、常量等构成 Rust 程序或库的容器.

    > Java 组织功能模块的主要单位是类; JavaScript 组织模块的主要方式是 function.

    > mod可以是一个文件或目录

    模块可以将一个 crate 中的代码进行分组, 以提高可读性与重用性. 模块还可以控制项的 私有性，即项是可以被外部代码使用的（public），还是作为一个内部实现的内容，不能被外部代码使用（private）.

    Rust 中**默认所有项（函数、方法、结构体、枚举、模块和常量）**都是私有的.

    父模块中的项不能使用子模块中的私有项，但是子模块中的项可以使用它们父模块中的项. 这是因为子模块封装并隐藏了它们的实现详情，但是子模块可以看到他们定义的上下文. 使用同级的项可忽略私有性.

    ```rust
    mod front_of_house {
        pub mod hosting {
            pub fn add_to_waitlist() {}
        }
    }

    pub fn eat_at_restaurant() {
        // 绝对路径
        crate::front_of_house::hosting::add_to_waitlist();

        // 相对路径
        front_of_house::hosting::add_to_waitlist(); // 使用同级的front_of_house
    }
    ```

    模块上的 pub 关键字只允许其父模块引用它.

    使用 super 开头可以用来构建从父模块开始的相对路径.

    如果在一个结构体定义的前面使用了 pub, 这个结构体会变成公有的，但是这个结构体的字段仍然是私有的. 如果将枚举设为公有，则它的所有成员都将变为公有.
    
    在作用域中增加 use 和路径类似于在文件系统中创建软连接（符号连接，symbolic link）.
    使用 use 将函数的父模块引入作用域，就必须在调用函数时指定父模块，这样可以清晰地表明函数不是在本地定义的，同时使完整路径的重复度最小化. 另一方面，使用 use 引入结构体、枚举和其他项时，习惯是指定它们的完整路径. 这种习惯用法背后没有什么硬性要求：它只是一种惯例，人们已经习惯了以这种方式阅读和编写 Rust 代码.

    使用 use 将两个同名类型引入同一作用域时可在这个类型的路径后面使用 as 指定一个新的本地名称或者别名.

    当使用 use 关键字将名称导入作用域时，在新作用域中可用的名称是私有的. 如果为了让调用你编写的代码的代码能够像在自己的作用域内引用这些类型，可以结合 pub 和 use. 这个技术被称为 “重导出（re-exporting）”，因为这样做将项引入作用域并同时使其可供其他代码引入自己的作用域.

    在 Cargo.toml 中加入 xxx 依赖告诉了 Cargo 要从 [crates.io](https://crates.io/) 下载 xxx 和其依赖，并使其可在项目代码中使用.

    注意标准库（std）对于你的包来说也是外部 crate。因为标准库随 Rust 语言一同分发，无需修改 Cargo.toml 来引入 std，不过需要通过 use 将标准库中定义的项引入项目包的作用域中来引用它们, 比如`use std::collections::HashMap;`, 这是一个以标准库 crate 名 std 开头的绝对路径.

    嵌套路径可消除大量的 use 行.

    通过 glob 运算符将所有的公有定义引入作用域, 比如`use std::collections::*;`. 使用 glob 运算符时请多加小心！Glob 会使得难以推导作用域中有什么名称和它们是在何处定义的. glob 运算符经常用于测试模块 tests 中，这时会将所有内容引入作用域.

    Rust 中有两种简单的访问权：公共（public）和私有（private）. 默认情况下，如果不加修饰符，模块中的成员访问权将是私有的. 对于私有的模块，只有在与其平级的位置或下级的位置才能访问，不能从其外部访问.

    > `use iron::prelude::*;`: 把 iron::prelude 模块中所有的公有名称直接暴露在代码中

    将模块内容迁移到其它文件: 如果模块声明时模块名后边是`;`, 而不是代码块时, rust会从与模块同名的文件中加载内容.

    pub对结构体、枚举、trait或模块的作用:
    - pub对于一个结构：它使结构公开，但成员不是公开的。要想让一个成员公开，你也要为每个成员写pub。
    - pub 对于一个枚举或trait：所有的东西都变成了公共的。这是有意义的，因为traits是给事物赋予相同的行为。而枚举是值之间的选择，需要看到所有的枚举值才能做选择。
    - pub对于一个模块来说：一个顶层的模块会是pub，因为如果它不是pub，那么根本没有人可以使用里面的任何东西。但是模块里面的模块需要使用pub才能成为公共的

- 路径（path）: 一个命名例如结构体、函数或模块等项的方式

    模块与文件不是一码事, 但 Unix 文件系统的文件和目录结构与模块具有天然的对应关系.

    路径有两种形式：
    - 绝对路径（absolute path）从 crate 根开始，以 crate 名或者字面值 crate 开头
    - 相对路径（relative path）从当前模块开始，以 self、super 或当前模块的标识符开头

        - self: 指向与当前模块相关的元素。该前缀用于任何代码想要引用自身包含的模块时. 这主要用于在父模块中重新导出子模块中的元素
        - super: 可用于从父模块导入元素, 诸如 tests 这类子模块将使用它从父模块导入元素

        self 和 super 与特殊目录 `.`和`..`类似, 比如`use super::super::*;`

    绝对路径和相对路径都后跟一个或多个由双冒号（::）分割的标识符

包主要解决项目间代码共享的问题，而模块主要解决项目内代码组织的问题.

对于一个由一系列相互关联的包组合而成的超大型项目, Cargo 提供了`工作空间`这一解决方案.

> crate是rust最小的编译单元, package是若干crate的集合, 它们都可被称为包. 只在两者同时出现且需要区别对待时，将crate译为单元包，将package译为包.

Rust的代码从逻辑上是分 crate 和 mod 管理的. crate 可以理解为"项目". 每个 crate 是一个完整的编译单元，它可以生成为一个 lib 或者可执行文件. 而在 crate 内部则由 mod 管理, mod 大家可以理解为 namespace. 可以使用 use 语句把其他模块中的内容引入到当前模块中来, 这样就解决了局部模块路径过长的问题. 对于同名导入可用`use ... as ...`解决.

Rust 有一个极简标准库， 叫作 std ，除了极少数嵌入式系统下无法使用标准库之外，绝大部分情况下，我们都需要用到标准库里面的东西. 为了给大家减少麻烦， Rust 编译器对标准库有特殊处理. 默认情况下，用户不需要手动添加对标准库的依赖 ，编译器会自动引人对标准库的依赖. 除此之外 ，标准库中的某些 type, trait、 function, macro 等实在是太常用了, 每次都写 use 语句确实非常无聊，因此标准库提供了一个[`std::prelude`'](https://github.com/rust-lang/rust/blob/master/library/std/src/prelude/mod.rs)模块(doc在[这](https://doc.rust-lang.org/std/prelude/index.html))，在这个模块中导出了一些最常见的类型, trait 等东西, 编译器会为用户写的每个 crate 自动插入一句话：`use std::prelude::*;`, 这样，标准库里面的这些最重要的类型, trait 等名字就可以直接使用，而无须每次都写全称或者 use 语句. 如果需要的类型不在 prelude 中，则必须使用 use 语句显式地将其引入作用域.

> [std::prelude](https://doc.rust-lang.org/std/prelude/)目前的[mod.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/mod.rs) 中, 直接导出了 v1 和 `rust_<edition>` 模块中的内容, 而 [v1.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/v1.rs) 和`rust_<edition>`中则是编译器为我们自动导入的相关trait和类型.

### workspace
一个 workspace 可以包含一到多个 crates，当代码发生改变时，只有涉及的 crates 才需要重新编译.

### 使用第三方包
Rust可以在Cargo.toml中的[dependencies]下添加想依赖的包来使用第三方包. 然后在`src/main.rs`或`src/lib.rs`文件中, 使用`extern crate`命令声明引入该包即可使用.

值得注意的是, 使用extern crate声明包的名称是linked_list，用的是下划线`_`, 而在Cargo.toml中用的是连字符`-`, 其实Cargo默认会把连字符转换成下划线.

Rust也不建议以`-rs`或`_rs`为后缀来命名包名, 而且会强制性的将此后缀去掉.

## Cargo.toml
```toml
[package]
    name = "cargo-metadata-example"
    version = "1.2.3"
    description = "An example of Cargo metadata"
    documentation = "https://docs.rs/dummy_crate"
    license = "MIT"
    readme = "README.md"
    keywords = ["example", "cargo", "mastering"]
    authors = ["Jack Daniels <jack@danie.ls>", "Iddie Ezzard <iddie@ezzy>"]
    build = "build.rs"
    edition = "2018"

[package.metadata.settings]
    default-data-path = "/var/lib/example"

[features]

    default=["mysql"]

[build-dependencies]
    syntex = "^0.58"

[dependencies]
    serde = "1.0"
    serde_json = "1.0"
    time = { git = "https://github.com/rust-lang/time", branch = "master" }
    mysql = { version = "1.2", optional = true } # optional表示可选
    sqlite = { version = "2.5", optional = true }
```

`[features]、 [dependencies]及[build-dependencies]`会组合到一起使用.

`[package]`:
- description：它包含一个关于项目的、更长的、格式自由的文本字段。
- license：它包含软件许可证标识符
- readme：它允许你提供一个指向项目版本库某个文件的链接。这通常是项目简介的入口点
- documentation：如果这是一个程序库，那么其中包含指向程序库说明文档的链接
- keywords：它是一组单词列表，有助于用户通过搜索引擎或者 crates.io 网站发现你的项目
- authors：它列出了该项目的主要作者。
- build：它定义了一段 Rust 代码（通常是 build.rs），它在cargo build前运行. 这通常用于生成代码或者构建项目程序所依赖的第三方非rust代码.
- edition：它主要用于指定编译项目时使用的 Rust 版本
- `profile.release`: 自定义rustc编译配置, 还有`profile.dev`, `profile.bentch`, `profile.test`

    - `debug = true`: 控制 rustc 中的 -g 选项, `cargo build --release`会得到一个带有调试符号的二进制文件, 优化设置不受影响.

`[package.metadata.settings]`:

    通常， Cargo 会对它无法识别的键或属性向用户发出警告，但是包含元数据的部分是个例外。它们会被 Cargo 忽略，因此可以用于配置项目所需的任何键/值对.

`[features]`: 与条件编译功能相关
- `default = ["mysql"]`:  用户在构建程序时如果没有手动覆盖功能集，则只会引入 mysql 而不包括sqlite

`[lib]`: 表示最终编译目标的信息
- name
- crate-type : `["dylib","staticlib"]`表示同时生成动态库和静态库
- path : 库入口, 默认是`src/lib.rs`
- test : 可否使用单元测试
- bench : 可否使用性能基准测试

## 测试
cargo支持的测试:
- 单元测试: 通常编写在包含被测试代码的同一模块中, 需要添加`#[cfg(test)]`, 其只在cargo test时生效
- 集成测试: 集成测试在程序库根目录下的 tests/目录中单独编写. 它们被构造成本身就像是被测试程序库的使用者.

## 配置文件
Cargo 配置文件和 git 类似, 支持层级概念:
1. 所有用户的全局配置`/.cargo/config`
1. 当前用户的全局配置`$HOME/.cargo/config`
1. 根包 regex 的配置`/regex/.cargo/config`
1. 子包bench的配置`/regex/bench/.cargo/config`

cargo 会从上到下层层覆盖, 上下层的配置并不会相互影响.

> root crate下可使用rustfmt.toml自定义代码风格

## FAQ
### `cargo build`报`Blocking waiting for file lock on package cache`
原因: Cargo.lock被其他程序正在写入, 独占了. 我这里是vscode占用了Cargo.lock.

如果确定没有多个程序占用, 可以执行`rm -rf ~/.cargo/.package-cache`, 然后再执行即可

### 升级`Cargo.toml`的edition
`cargo fix`

### cargo添加mirror
ref:
- [Rust crates.io 索引镜像使用帮助](https://mirrors.tuna.tsinghua.edu.cn/help/crates.io-index.git/)
- [cargo config层级](https://doc.rust-lang.org/cargo/reference/config.html)

```bash
# vim $CARGO_HOME/config.toml
[source.crates-io]
replace-with = 'tuna'

[source.tuna]
registry = "https://mirrors.tuna.tsinghua.edu.cn/git/crates.io-index.git"
```

### Caret requirements(跳脱条件)
指定版本号范围的标记有以下几种:
- 补注号`^` : 允许新版本号在不修改`[major,minor,patch]` 中最左边非零数字的情况下更新
- 通配符`*` : 可以用在`[major,minor,patch]`的任何一个上面
- 波浪线`~` : 允许修改`[major,minor,patch]`中没有明确指定的版本号
- 手动指定 : 通过`>, >=, <, <=, =`来指定版本号

```
^1.2.3 := >=1.2.3 <2.0.0
^1.2 := >=1.2.0 <2.0.0
^1 := >=1.0.0 <2.0.0
^0.2.3 := >=0.2.3 <0.3.0
^0.2 := >= 0.2.0 < 0.3.0
^0.0.3 := >=0.0.3 <0.0.4
^0.0 := >=0.0.0 <0.1.0
^0 := >=0.0.0 <1.0.0
```
### 大小优化
- [min-sized-rust](https://github.com/johnthagen/min-sized-rust)

### `fatal error: 'stddef.h' file not found`
cargo build报`'stddef.h' file not found`, cargo使用llvm组建, 安装clang即可