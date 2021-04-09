# cargo
```sh
$ cargo new xxx # 创建一个编写可执行文件的项目
$ cargo new --lib xxx # 创建用于编写库的项目
$ cargo build # 编译项目
$ cargo run # 运行项目
```

## 代码组织
Rust的代码从逻辑上是分 crate 和 mod 管理的. crate 可以理解为"项目". 每个 crate 是一个完整的编译单元，它可以生成为一个 lib 或者可执行文件. 而在 crate 内部则由 mod 管理, mod 大家可以理解为 namespace. 可以使用 use 语句把其他模块中的内容引入到当前模块中来.

Rust 有一个极简标准库， 叫作 std ，除了极少数嵌入式系统下无法使用标准库之外，绝大部分情况下，我们都需要用到标准库里面的东西. 为了给大家减少麻烦， Rust 编译器对标准库有特殊处理. 默认情况下，用户不需要手动添加对标准库的依赖 ，编译器会自动引人对标准库的依赖. 除此之外 ，标准库中的某些 type, trait、 function, macro 等实在是太常用了, 每次都写 use 语句确实非常无聊，因此标准库提供了一个[`std::prelude`'](https://github.com/rust-lang/rust/blob/master/library/std/src/prelude/mod.rs)模块，在这个模块中导出了一些最常见的类型, trait 等东西, 编译器会为用户写的每个 crate 自动插入一句话：`use std: :prelude::*;`, 这样，标准库里面的这些最重要的类型, trait 等名字就可以直接使用，而无须每次都写全称或者 use 语句.

> [std::prelude](https://doc.rust-lang.org/std/prelude/)目前的[mod.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/mod.rs) 中, 直接导出了 v1 模块中的内容, 而 [v1.rs](https://github.com/rust-lang/rust/blob/master/library/core/src/prelude/v1.rs) 中则是编译器为我们自动导人的相关trait和类型.