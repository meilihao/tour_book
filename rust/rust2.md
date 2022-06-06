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

- 错误跳转: 在错误跳转中，当调用的函数返回错误时，Rust 会提前终止当前函数的执行，向上一层返回错误

	`expr?`, 比如`fs::write("/tmp/1.log", b"hello")?;`
- 异步跳转: 在 Rust 的异步跳转中, 当 async 函数执行 await 时, 程序当前上下文可能被阻塞, 执行流程会跳转到另一个异步任务执行, 直至 await 不再阻塞.

	`expr.await`, 比如`socket.write(data).await?`