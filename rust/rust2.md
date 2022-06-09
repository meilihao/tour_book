# rust

# base
## 变量, 函数和数据结构
- 变量

	- 不可变: `let x:T;`
	- 可变: `let mut x:T;`

		`let mut v:Vec<u8> = Vec::new();`
- 常量: `const X:T = <value>;`

	常量是一个右值, 它不能被修改. 常量编译后会被放在可执行文件的数据段, 全局可访问.
- 静态变量

	- 不可变: `static X:T = T::new();`
	- 可变: `static mut X:T = T::new();`

	静态变量和常量一样全局可访问, 它也被编入可执行文件的数据段. 静态变量可以被声明为可变.

	在使用静态变量时, 有一些限制, 可用lazy_static检查
- 函数: `func x(a1:T1,...)-> T{}`

	在rust中， 如果函数没有返回值， 那么其返回值是unit即空元组`()`
- 结构体: `struct S{...}`

	三种形式:
	1. 空结构体, 不占用任何空间, 比如`struct Marker;`
	1. 元组结构体, struct的每个成员都是匿名的, 可通过索引范围, 比如`struct Color(u8,u8,u8);`
	1. 普通结构体, struct的每个成员都有名字, 可通过名字访问

		```rust
		struct Person{
			name: String,
			age: u8,
		}
		```
- enum: `enum E{...}`

	两种形式:
	1. 枚举

		```rust
		enum Status {
			Ok = 0,
			Bad = 1,
			NotFound = 2,
			...
		}
		```
	1 标签联合

		enum可承载多个不同的数据结构中的一种.

		```rust
		enum Option<T>{
			Some(T),
			None,
		}
		```
## 流程控制
Rust 的循环和大部分语言都一致, 支持死循环`loop {}`、条件循环`while expr {}`，以及对迭代器的循环`for x in iter {}`. 循环可以通过 break 提前终止，或者 continue 来跳到下一轮循环.

Rust 的 for 循环可以用于任何实现了  IntoIterator trait 的数据结构. 在执行过程中，IntoIterator 会生成一个迭代器，for 循环不断从迭代器中取值，直到迭代器返回 None 为止。因而，**for 循环实际上只是一个语法糖，编译器会将其展开使用 loop 循环对迭代器进行循环访问，直至返回 None**

通过斐波那契数列, 使用 if 和 loop / while / for 这几种循环，来实现程序的基本控制流程:
```rust

fn fib_loop(n: u8) {
    let mut a = 1;
    let mut b = 1;
    let mut i = 2u8;
    
    loop {
        let c = a + b;
        a = b;
        b = c;
        i += 1;
        
        println!("next val is {}", b);
        
        if i >= n {
            break;
        }
    }
}

fn fib_while(n: u8) {
    let (mut a, mut b, mut i) = (1, 1, 2);
    
    while i < n {
        let c = a + b;
        a = b;
        b = c;
        i += 1;
        
        println!("next val is {}", b);
    }
}

fn fib_for(n: u8) {
    let (mut a, mut b) = (1, 1);
    
    for _i in 2..n { // Range 的下标上标都是 usize 类型，不能为负数
        let c = a + b;
        a = b;
        b = c;
        println!("next val is {}", b);
    }
}

fn main() {
    let n = 10;
    fib_loop(n);
    fib_while(n);
    fib_for(n);
}
```

满足某个条件时会跳转, Rust 支持
- 分支跳转: `if/else`
- 模式匹配: Rust 的模式匹配可以通过匹配表达式或者值的某部分的内容，来进行分支跳转
	
	需要根据表达式所有可能的值进行匹配, 并进行相应的处理.

	`match expr {}`或`if let pat = expr {}`, `if let`是match的简写, 表示仅关心某种模式匹配的情况

	Rust 的模式匹配吸取了函数式编程语言的优点，强大优雅且效率很高. 它可以用于 struct / enum 中匹配部分或者全部内容.

- 错误跳转: 在错误跳转中，当调用的函数返回错误时，Rust 会提前终止当前函数的执行，向上一层返回错误

	`expr?`, 比如`fs::write("/tmp/1.log", b"hello")?;`
- 异步跳转: 在 Rust 的异步跳转中, 当 async 函数执行 await 时, 程序当前上下文可能被阻塞, 执行流程会跳转到另一个异步任务执行, 直至 await 不再阻塞.

	`expr.await`, 比如`socket.write(data).await?`

## 错误处理
Rust 没有沿用 C++/Java 等诸多前辈使用的异常处理方式, 而是借鉴 Haskell，把错误封装在  `Result<T, E>` 类型中, 同时提供了`?`操作符来传播错误, 方便开发. `Result<T, E>` 类型是一个泛型数据结构，T 代表成功执行返回的结果类型, E 代表错误类型.

## 代码管理
rust支持使用mod 来组织代码.

使用方法: 在项目的入口文件`lib.rs/main.rs`里, 用 mod 来声明要加载的其它代码文件. 如果模块内容比较多, 可以放在一个目录下, 再在该目录下放一个 mod.rs 引入该模块的其它文件, mod.rs 和 Python 的 `__init__.py` 有异曲同工之妙.

在 Rust 里, 一个项目也被称为一个 crate. crate 可以是可执行项目，也可以是一个库.

在一个 crate 下，除了项目的源代码，单元测试和集成测试的代码也会放在 crate 里. Rust 的单元测试一般放在和被测代码相同的文件中，使用条件编译  #[cfg(test)] 来确保测试代码只在测试环境下编译.

集成测试一般放在 tests 目录下，和 src 平行. 和单元测试不同，**集成测试只能测试 crate 下的公开接口，编译时编译成单独的可执行文件**. 在 crate 下，如果要运行测试用例，可以使用`cargo test`.

当代码规模继续增长，把所有代码放在一个 crate 里就不是一个好主意了，因为任何代码的修改都会导致这个 crate 重新编译，这样效率不高. 此时可以使用 workspace, 一个 workspace 可以包含一到多个 crates，当代码发生改变时，只有涉及的 crates 才需要重新编译. 当要构建一个 workspace  时，需要先在某个目录下生成一个 Cargo.toml，包含 workspace 里所有的 crates，然后可以  cargo new 生成对应的 crates.