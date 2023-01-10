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
cargo build --release # 这种模式会以更长的编译时间为代价来优化代码, 从而使代码拥有更好的运行时性能, 但不会显示 panic backtrace 的具体行号
cargo run [--release] # 运行项目
cargo install --path
cargo uninstall first_pro_create
cargo fmt: 类似gofmt, 格式化代码
cargo clippy: 类似eslint, 检查代码规范
cargo tree: 查看第三方库的版本和依赖关系
cargo bench: 运行benchmark(基准测试,性能测试)
cargo udeps(第三方): 检查项目中未使用的依赖
```

## 代码组织
Rust 有许多功能可用于组织和管理代码, 包括哪些内容可以被公开, 哪些内容作为私有部分, 以及程序每个作用域中的名字. 这些功能有时被称为`模块系统(the module system)`包括:
- 包（Packages） : Cargo 的一个功能，它允许构建、测试和分享 crate

    一个包会包含有一个 Cargo.toml 文件，描述了包的基本信息以及依赖项以及如何去构建这些 crate.

    包中所包含的内容由几条规则来确立:
    1. 一个包中至多只能包含一个库 crate(library crate)
    1. 包中可以包含任意多个二进制 crate(binary crate)
    1. 包中至少包含一个 crate，无论是库的还是二进制的

    执行`cargo new my-project`, 查看 Cargo.toml 的内容，会发现并没有提到 src/main.rs, 因为 Cargo 遵循的一个约定: src/main.rs 就是一个与包同名(这里就是my-project)的二进制 crate 的 crate 根. 同样的, Cargo 知道如果包目录中包含 src/lib.rs, 则 src/lib.rs 就是一个与包同名的 lib crate的 crate 根. crate 根文件将由 Cargo 传递给 rustc 来实际构建库或者二进制项目.

    如果一个包同时含有 src/main.rs 和 src/lib.rs，则它有两个 crate：一个库和一个二进制项，且名字都与包相同. 通过将文件放在 src/bin 目录下，一个包可以拥有多个二进制 crate：每个 src/bin 下的文件都会被编译成一个独立的二进制 crate.

- Crates : 一个模块的树形结构，它形成了库或二进制项目, 存在于`包`中.

    crate root 是一个源文件, Rust 编译器以它为起始点, 同时也是 crate 的根模块.
- 模块（Modules即mod）和 use : 允许控制作用域和路径的私有性

    > Java 组织功能模块的主要单位是类; JavaScript 组织模块的主要方式是 function.

    模块可以将一个 crate 中的代码进行分组, 以提高可读性与重用性. 模块还可以控制项的 私有性，即项是可以被外部代码使用的（public），还是作为一个内部实现的内容，不能被外部代码使用（private）.

    Rust 中**默认所有项（函数、方法、结构体、枚举、模块和常量）**都是私有的。 父模块中的项不能使用子模块中的私有项，但是子模块中的项可以使用他们父模块中的项. 这是因为子模块封装并隐藏了他们的实现详情，但是子模块可以看到他们定义的上下文.

    模块上的 pub 关键字只允许其父模块引用它.

    使用 super 开头可以用来构建从父模块开始的相对路径.

    如果在一个结构体定义的前面使用了 pub, 这个结构体会变成公有的，但是这个结构体的字段仍然是私有的.
    如果将枚举设为公有，则它的所有成员都将变为公有.
    在作用域中增加 use 和路径类似于在文件系统中创建软连接（符号连接，symbolic link）.
    使用 use 将函数的父模块引入作用域，就必须在调用函数时指定父模块，这样可以清晰地表明函数不是在本地定义的，同时使完整路径的重复度最小化. 另一方面，使用 use 引入结构体、枚举和其他项时，习惯是指定它们的完整路径. 这种习惯用法背后没有什么硬性要求：它只是一种惯例，人们已经习惯了以这种方式阅读和编写 Rust 代码.

    使用 use 将两个同名类型引入同一作用域时可在这个类型的路径后面使用 as 指定一个新的本地名称或者别名.

    当使用 use 关键字将名称导入作用域时，在新作用域中可用的名称是私有的. 如果为了让调用你编写的代码的代码能够像在自己的作用域内引用这些类型，可以结合 pub 和 use. 这个技术被称为 “重导出（re-exporting）”，因为这样做将项引入作用域并同时使其可供其他代码引入自己的作用域.

    在 Cargo.toml 中加入 xxx 依赖告诉了 Cargo 要从 [crates.io](https://crates.io/) 下载 xxx 和其依赖，并使其可在项目代码中使用.

    注意标准库（std）对于你的包来说也是外部 crate。因为标准库随 Rust 语言一同分发，无需修改 Cargo.toml 来引入 std，不过需要通过 use 将标准库中定义的项引入项目包的作用域中来引用它们, 比如`use std::collections::HashMap;`, 这是一个以标准库 crate 名 std 开头的绝对路径.

    嵌套路径可消除大量的 use 行.

    通过 glob 运算符将所有的公有定义引入作用域, 比如`use std::collections::*;`. 使用 glob 运算符时请多加小心！Glob 会使得难以推导作用域中有什么名称和它们是在何处定义的. glob 运算符经常用于测试模块 tests 中，这时会将所有内容引入作用域.

    Rust 中有两种简单的访问权：公共（public）和私有（private）. 默认情况下，如果不加修饰符，模块中的成员访问权将是私有的. 对于私有的模块，只有在与其平级的位置或下级的位置才能访问，不能从其外部访问.

- 路径（path）: 一个命名例如结构体、函数或模块等项的方式

    路径有两种形式：
    - 绝对路径（absolute path）从 crate 根开始，以 crate 名或者字面值 crate 开头。
    - 相对路径（relative path）从当前模块开始，以 self、super 或当前模块的标识符开头。

    绝对路径和相对路径都后跟一个或多个由双冒号（::）分割的标识符

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