# rust
参考:
- [Rust 程序设计语言（第二版 & 2018 edition）](https://kaisery.github.io/trpl-zh-cn/)
- [Rust入门第一课](https://rust-book.junmajinlong.com/ch1/00.html)

程序的三大定律:
1. 程序必须正确
1. 程序必须可维护, 但不能违反第一定律
1. 程序必须高效, 但不能违反前两条定律

rust无疑是迈向这个目标的更近一步.

编程两大难题(本质是使用类型不安全的语言, 它们的内存管理机制不完善):
1. 编写内存安全的代码
1. 编写线程安全的代码

> 在之前的年代, 计算资源匮乏, 为追求性能, 牺牲了部分安全性.

rust是一门同时追求安全,并发和性能的现代系统级编程语言.
rust三大设计哲学:
1. 内存安全

    类型安全即类型系统可以保证程序的行为是意义明确, 不出错的. c/c++的类型系统不是类型安全的, 比如它们不检查数组越界.

    > 类型安全: 类型系统给内存中的plain data赋予了类型信息，如果不按照类型信息来解释内存中的数据，应该产生编译错误或者 产生well-specified的运行时错误.

    类型安全是内存安全的前提.

    内存安全: 不会出现内存访问错误, 只有当程序访问未定义内存时才会产生内存错误. 常见的场景有:
    1. 引用空指针
    1. 使用未初始化内存
    1. 释放后再次使用, 即虚悬指针
    1. 缓冲区溢出, 比如数组越界
    1. 非法释放已释放过的指针或未分配的指针, 即重复释放

    rust为了保证内存安全, 建立了严格的安全内存管理模型:
    - 所有权系统: 每个被分配的内存都有一个独占其所有权的指针. 只有当该指针被销毁时, 其对应的内存该会被随之释放.
    - 借用和生命周期: 每个变量都有其生命周期, 一旦超出生命周期, 变量就会被自动释放. 如果是借用, 则可以通过标记生命周期参数供编译器检查的方式, 防止出现虚悬指向即释放后使用的情况.

    > rust所有权系统还包括了从现代c++借鉴的RAII机制, 这是rust无gc但是可以安全管理内存的基石. 在 C++ 中，这种 item 在生命周期结束时释放资源的模式被称作 资源获取即初始化（Resource Acquisition Is Initialization (RAII)）.
    > 悬垂指针（dangling pointer, 也叫虚悬指针）是其指向的内存可能已经被分配给其它持有者.

    为了实现内存安全, rust还具备具有的特性:
    1. 仿射类型(Affine Type), 该类型可用来表达rust所有权中的Move语义.

    借助类型, rust可在编译阶段对类型进行检查是否满足安全内存模型, 有效地阻止未定义行为的发生.

    内存安全bug和并发安全的bug的产生内在原因均是内存的不正当访问造成的. 借助装备了所有权的强大类型系统, rust还解决了并发安全问题. 它通过静态检查分析, 在编译期就能检查出多线程并发代码中所有的数据竞争问题.
1. 零成本抽象即如果不使用某个抽象, 就不用为它付出开销.

    rust的绝大多数抽象并不存在运行时的开销, 其一切都是在编译期完成的.

    rust的零成本抽象的基石是泛型和trait.
1. 实用性

    为了保证支持硬实时, rust借鉴了c++的确定性析构, RAII和智能指针, 用于自动地, 确定地管理内存, 从而避免了gc的引入.

    为了保证程序的鲁棒性, rust重新审视了错误处理机制. rust针对三类非正常情况: 失败, 错误和异常, 提供了专门的处理方式:
    - 失败: 使用断言工具
    - 错误: 基于返回值的分层处理 
    - 异常: rust将其看作无法被合理解决的问题, 提供了线程恐慌机制, 发生异常时, 线程可以安全地退出.

    为了兼容现有生态, rust支持方便且零成本的FFI机制, 兼容C-ABI, 在语言架构层面上将rust分为safe rust和unsafe rust两部分. unsafe专门和外部生态打交道, 因为rust编译器检查和跟踪的范围有限, 不能检查到与其链接的其他生态接口, 因此这些生态由自身来保证安全性. 总结就是, safe rust由rust编译器在编译时保证安全, unsafe rust开发者让编译器信任自身有能力保证安全.

rust编程的哲学和golang相同: 组合优于继承. rust不提供类型层面上的继承, 所有的类型都是独立存在的.

从其他语言到Rust的最大思维转换就是变量的所有权和生命周期，这是几乎所有编程语言都未曾涉及的领域.

### 编译
rust编译器是LLVM编译器的前端, 它将代码编译成LLVM IR, 然后通过LLVM编译成对应架构的机器码.

rust源码经过分词和解析生成AST(抽象语法树), 再进一步简化处理为HIR(High-level IR, 方便编译器做类型检查), 再进一步编译为[MIR](https://play.rust-lang.org/)(middle IR, 在rust 1.12引入), 最后MIR被翻译为LLVM IR, 之后由LLVM编译成目标机器码.

引入MIR原因:
1. 缩短编译时间

    实现了增量编译, 仅重新编译更改过的部分.
1. 缩短执行时间

    进入llvm前实现更细颗粒度的优化, 单纯依赖llvm的优化颗粒度太粗, 增加了更多的优化空间
1. 更精确的类型检查

    实现更灵活的借用检查

MIR是基于控制流图（Control Flow Graph，CFG）的抽象数据结构，它用有向图（DAG）形式包含了程序执行过程中所有可能的流程, 所以将基于MIR的借用检查称为非词法作用域的生命周期, 因此确实不依赖词法作用域了.

MIR由一下关键部分组成：
- 基本块（Basic block，bb），他是控制流图的基本单位，
    
    - 语句（statement）
    - 终止句（Terminator）

- 本地变量, 栈中内存的位置，比如函数参数、局部变量等. 一般用下划线和数字作为标识(比如`_1`), 其中`_0`通常表示返回地址.
- 位置（Place），在内存中标识未知的额表达式。
- 右值（RValue），产生值的表达式。


> LLVM是构架编译器(compiler)的框架系统，以C++编写而成，用于优化以任意程序语言编写的程序的编译时间(compile-time), 链接时间(link-time), 运行时间(run-time)以及空闲时间(idle-time).

> 无疑，不同编译器的中间语言IR是不一样的，而IR可以说是集中体现了这款编译器的特征：它的算法，优化方式，汇编流程等等，想要完全掌握某种编译器的工作和运行原理，分析和学习这款编译器的中间语言无疑是重要手段. 由于中间语言相当于一款编译器前端和后端的“桥梁”，如果我们想进行基于LLVM的后端移植，无疑需要开发出对应目标平台的编译器后端，想要顺利完成这一工作，透彻了解LLVM的中间语言无疑是非常必要的工作. LLVM相对于gcc的一大改进就是大大提高了中间语言的生成效率和可读性, LLVM的中间语言是一种介于c语言和汇编语言的格式, 它既有高级语言的可读性, 又能比较全面地反映计算机底层数据的运算和传输的情况, 精炼而又高效.

rust语言组成:
- 语言规范
- 编译器
- 核心库

    定义了rust语言的核心, 不依赖于操作系统和网络等相关的库, 不负责堆分配,不提供并发和I/O.
    通过在模块顶部引入`#![no_std]`来使用核心库.
    组成:
    1. 基础trait, 如Copy, Debug, Display,Option等
    1. 基本原始类型, 如 bool, char,i8/u8,..i64/u64, isize/usize, f32/f64, str, array, slic, tuple, pointer等
    1. 常用功能型数据类型, 比如String, Vec, HashMap, Rc, Arc,Box等

        String 的类型是由标准库提供的，而没有写进核心语言部分，它是可增长的、可变的、有所有权的、UTF-8 编码的字符串类型.

        Rust 标准库中还包含一系列其他字符串类型，比如 OsString、OsStr、CString 和 CStr, 前缀(非String/非Str)对应着它们提供的所有权和可借用的字符串变体.

        使用 to_string 方法从字符串字面值创建 String: `"...".to_string()` <=> `String::from("...")`.

        **Rust 的字符串不支持索引**.

        遍历: `for c in "नमस्ते".chars()`或`for b in "नमस्ते".bytes()`.

        Rust 的原始字符串语法是在`r`后跟一个或多个`#`和一个双引号, 然后是字符串内容, 最后以另一个双引号及相同个数的`#`结尾. 原始字符串中的任何字符都无须转义，包括双引号. 为避免歧义，通过在双引号两侧添加更多的`#`，总是可以明确标识字符串的结束位置
    1. 常用的宏定义, 如println()!, assert!, panic!, vec!

        Rust 中的宏(`!`)与C/C＋＋ 中的宏是完全不一样的东西. 简单点说，可以把它理解为一种安全版的编译期语法扩展, 这里之所以使用宏，而不是函数，是因为标准输出宏可以完成编译期格式检查. 更加安全. 这个宏最终还是调用了`std::io`模块内提供的一些函数来完成的. 如果需要更精细地控制标准输出操作, 也可以直接调用标准库来完成.

        Rust 支持声明宏（declarative macro）和过程宏（procedure macro），其中过程宏又包含三种方式：函数宏（function macro），派生宏（derive macro）和属性宏（attribute macro）. println! 是函数宏，是因为 Rust 是强类型语言，函数的类型需要在编译期敲定，而 println! 接受任意个数的参数，所以只能用宏来表达.

    做嵌入式应用开发时, 核心库是必需的.
- 标准库

    提供应用程序开发所需的基础和跨平台支持.

    组成:
    1. 与核心库一样的基本trait, 原始数据类型, 功能型数据类型和常用的宏等, 以及与核心库几乎完全一致的API.
    1. 并发, I/O和运行时.
    
        如线程模块,  用于消息传递的通道类型, Sync trait等并发模块; 文件, tcp, udp, 管道, socket等常见I/O.
    1. 平台抽象.
        - os模块提供了许多与操作环境交互的基本功能, 包括程序参数, 环境变量, 目录导航
        - 路径模块封装了处理文件路径的平台特定规则.
    1. 底层操作接口

        如std::mem,std::ptr, std::intrinsic等, 操作内存,指针, 调用编译器固有函数.

    1. 可选和错误处理类型Option(表示一个值要么有值要么没值)和Result, 以及各种迭代器.

    Rust 并没有空值，不过它确实拥有一个可以编码存在或不存在概念的枚举, 这个枚举是 Option<T>，而且它定义于标准库中.

### Rust命令规范
> rust使用蛇形命名法(snake case)来规范函数和变量名称的风格: 只使用小写的字母进行命名, 并以`_`分割word.

- 函数： 蛇形命名法（snake_case），例如：func_name()
- 文件名： 蛇形命名法（snake_case），例如file_name.rs、main.rs
- 临时变量名：蛇形命名法（snake_case）
- 全局变量名：
- 结构体： 大驼峰命名法，例如：struct FirstName { name: String}
- enum类型: 大驼峰命名法
- 关联常量：常量名必须全部大写
- Cargo默认会把连字符`-`转换成下划线`_`
- Rust也不建议以`-rs`或`_rs`为后缀来命名包名, 而且会强制性的将此后缀去掉
  
### 语句和表达式
Rust 是一门基于表达式（expression-based）的语言. 语句（Statements）是执行一些操作但不返回值的指令. 表达式（Expressions）计算并产生一个值.

rust使用`{}`表示复杂表达式, 比如:
```rust
    // 在块中可以使用函数语句, 最后一个步骤是表达式, 此表达式的结果值是整个表达式块所代表的值, 这种表达式块叫做函数体表达式.
    // 函数体表达式并不能等同于函数体, 它不能使用 return 关键字
    let y = {
        let x = 3;
        x + 1
    };
```

> 使用 let 关键字创建变量并绑定一个值是一个语句, 因为语句不返回值， 因此，不能把 let 语句赋值给另一个变量. 在C语言或Ruby语言中的赋值语句会返回所赋的值, 因此`x = y = 6`是正确的, 它使得x和y变量同时拥有6这个值; 而rust不能这样写.

> Rust 源代码的后缀名使用`.rs`表示, 且必须使用 utf-8 编码.
rust类型没有“默认构造函数”，变量没有“默认值” .
Rust 里面的下划线是一个特殊 的标识符，在编译器内部它是被特殊处理的.

rust语法分语句(statement, 要执行的一些操作和产生副作用的表达式)和表达式(expression, 主要用于计算求值). 它们的区别是后面带不带分号.

语句分类:
- 声明语句(Declaration statement) : 用于声明各种语言项(item), 比如变量, 静态变量, 常量, 结构体, 函数等, 以及通过extern和use关键词引入的包和模块等.
- 表达式语句(expression statement) : 特指以分号结尾的表达式. 此类表达式求值结果会被舍弃, 并总是返回单元类型`()`.

表达式分类:
- 赋值表达式

    类型是unit, 即空的`tuple()`
- 字面量表达式
- 方法调用表达式
- 数组表达式
- 索引表达式
- 单目/双目运算符表达式等

表达式的其他分类:
- 左值 : 这个表达式可以表示一个内存地址, 因此它在赋值运算符的左边
- 右值 : 除左值外的值

rust编译器解析代码时, 如果遇到分号, 就会继续往后面执行; 如果碰到语句, 则执行语句; 如果碰到表达式, 则会对表达式求值, 如果分号后面什么都没有, 就会补上单元值;  当
遇到函数时, 会将函数体的花括号识别为块表达式(block expression, 由一对花括号和一系列表达式组成, 它总是返回块中最后一个表达式的值)

> 赋值号左边的部分是一个“模式”，`let (mut a, mut b) = (1, 2);`是对 tuple 的模式解构，`let Point { x : ref a, y : ref b} = p;`是对结构体的模式解构.

> 模式匹配可能导致变量所有权转移, 使用ref可避免该情况. ref是模式的一部分, 它只能出现在赋值操作符的左边. `&`是借用运算符, 是表达式的一部分, 只能出现在赋值操作符的右边.

> mut在模式匹配中是修饰变量绑定, `&mut`是表示引用, mut的意义完全不同.

let创建的变量一般称为绑定(bingding).

rust的表达式可分为位置表达式(place expression)和值表达式(value expression), 即其他语言中的左值和右值.

位置表达式就是表示内存位置的表达式, 分别有如下几类:
1. 本地变量
1. 静态变量
1. 解引用(*expr)
1. 数组索引(expr[expr])
1. 字段引用(expr.field)
1. 位置表达式组合

通过位置表达式可对某个数据单元的内存进行读写.
值表达式一般只引用了某个存储单元地址中的数据, 它相当于数据值, 只能进行读操作.

**从语法角度讲, 位置表达式代表了持久化数据, 值表达式代表了临时数据.**

表达式的求值过程可分为位置上下文(place context)和值上下文(value context).

let关键字声明位置表达式默认不可变, 即**不可变绑定**, 只能对对应的存储单元进行读取. 引用与变量一样，默认情况下 也是不可变的.
let mut声明的**可变绑定**可以对相应的存储单元进行写入.

> 在 Rust 中，一般把声明的局部变量并初始化的语句称为`变量绑定`, 强调的是“绑定”的含义，与 C/C＋＋ 中的“赋值初始化”语句有所区别: Rust中，每个变量必须被合理初始化之后才能被使用, 使用未初始化变量这样的错误，在Rust 中是不可能出现的 （利用 unsafe 做 hack 除外）. 编译器会帮我们做一个执行路径的静态分析，确保变量在使用前一定被初始化.
> 类型一定是跟在冒号`:`后面, 这样语法歧义更少，语法分析器更容易编写.

> rust支持定义一个与之前变量同名的新变量，而新变量会隐藏(shadow)之前的变量. 隐藏区别于mut: 隐藏会创建新变量, 同时可改变其类型.
> const(常量)可以在任意作用域进行定义，而定义的常量贯穿整个程序的生命周期. 在编译的时候，常量就能确定其值, 编译器不一定会给const分配内存空间, 同时编译器**会尽可能将其内联到代码中，所以在不同地方对同一常量的引用并不能保证引用到相同的内存地址**；对于变量出现重复的定义(绑定)会发生变量遮盖，而对于常量则是不允许出现重复的定义的.
> 全局(static)变量和常量类似，但static变量不会被内联，在整个程序中，全局变量只有一个实例, **必须是编译期可确定的常量**，也就是说所有的引用都会指向一个相同的地址. 因为全局变量可变，就会出被多个线程同时访问的情况，因而引发内存不安全的问题，所以对于全局可变(static mut)变量的访问和修改代码就必须在unsafe块中进行定义. 声明static变量时禁止调用普通函数, 但const fn可以, 因为const fn是编译时执行的.

在存储的数据比较大，需要引用地址或具有可变性的情况下使用静态变量. 否则，应该优先使用普通常量. 推荐使用lazy_static实现复杂的static变量初始化.

所有权转移: 当位置表达式出现在值上下文中, 表示将会把内存地址转移给另一个位置表达式.

**在语义上, 每个变量绑定实际上都拥有该存储单元的所有权, 这种转移内存地址的行为就是所有权(ownership)的转移, 在rust中称为移动(move)语义, 那种不转移的情况实际上是一种复制(copy)语义. **rust没有gc, 所以完全依靠所有权来进行内存管理.

rust提供引用操作符(&), 此时不转移所有权, 可直接获取表达式的存储单元地址(即内存位置), 可通过该地址对存储进行读取. 因此引用已被称为借用.

rust的作用域是静态作用域, 即词法作用域(lexical scope), 由一对花括号来开辟作用域, 该作用域在词法分析阶段就已确定, 不会动态改变.

Rust允许在同一个代码块中声明同样名字的变量的做法被称为变量遮蔽(variable shadow), **个人强烈不推荐使用**.

变量绑定的生命周期(lifetime): 从使用let声明创建变量绑定开始, 到超出词法作用域的范围时结束.

rust中函数是一等公民, 其自身可作为函数的参数和返回值使用.

rust编译器像C++或D语言一些, 支持编译时函数执行(compile-time function execution, CTFE), 该功能由miri(mir解析器, 已集成入rustc)来执行.

rust中固定长度的数组必须在编译期就知道长度, 否知会报错, 这就用了CTFE的能力.

`const fn`强制编译器在编译期执行函数, 其中const一般用于定义全局常量.

const 函数是纯函数，必须是可重现的。这意味着它们不能将可变参数带入任何类型，也不能包含动态的操作，例如堆分配.

当前rust支持的常量有: 字面量, 元组, 数组, 字段结构体, 枚举, 只包含单行代码的块表达式, 范围等.

函数是对代码中重复行为的抽象。在现代编程语言中，函数往往是一等公民，这意味着函数可以作为参数传递，或者作为返回值返回，也可以作为复合类型中的一个组成部分.

闭包与函数的区别: 闭包可捕获外部变量, 而函数不可以.

闭包是将函数，或者说代码和其环境一起存储的一种数据结构。闭包引用的上下文中的自由变量，会被捕获到闭包的结构中，成为闭包类型的一部分.

一般来说，如果一门编程语言，其函数是一等公民，那么它必然会支持闭包（closure），因为函数作为返回值往往需要返回一个闭包

rust闭包实际上由一个匿名结构体和trait来组合实现的, 语法是`|params| {expr}`.

```rust
// 一般情况下, 闭包默认会按引用捕获变量, 如果将此闭包返回, 则引用也会跟着返回, 但这里随着函数调用结束, 本地变量i会被销毁, 随闭包返回的i的引用就变成了虚悬指针.
// 因此需要使用move来将被引用的变量的所有权转移到闭包中, 不再按引用捕获变量, 这样闭包才可以安全地返回.
pub fn two_times_impl() -> impl Fn(i32) -> i32 {
    let i = 2;
    move |j| j * i
}
```

#### 格式化输出

打印操作由 std::fmt 里面所定义的一系列宏来处理，包括：
- format!：将格式化文本写到字符串（String）
- print!：与 format! 类似，但将文本输出到控制台（io::stdout）
- println!: 与 print! 类似，但输出结果追加一个换行符
- eprint!：与 format! 类似，但将文本输出到标准错误（io::stderr）
- eprintln!：与 eprint! 类似，但输出结果追加一个换行符

[println! format](https://doc.rust-lang.org/std/fmt/#formatting-traits):
- nothing(即默认)⇒ Display : `println!("{}",2)`
- ? ⇒ Debug : `println!("{:?}",2)`
- x? ⇒ Debug with lower-case hexadecimal integers
- X? ⇒ Debug with upper-case hexadecimal integers
- :#? ⇒ 带换行和缩进的Debug打印 : `println!("{:#?}",("t1","t2"))` 
- o ⇒ Octal : `println!("{:o}",2)`
- x ⇒ LowerHex : `println!("{:x}",2)`
- X ⇒ UpperHex : `println!("{:X}",2)`
- p ⇒ Pointer : `println!("{:p}",&2)`
- b ⇒ Binary :  : `println!("{:b}",2)`
- e ⇒ LowerExp : `println!("{:e}",2)`
- E ⇒ UpperExp : `println!("{:E}",2)`
- 命名参数       : `println!("{a} {b} {b}",a = "x", b="y")`

> 想把 a 输出两遍可用: `println!("a is {0}, a again is {0}", a);`

### 控制流
> 代码中的条件表达式必须产生一个bool类型的值，否则就会触发编译错误. 不像 Ruby 或 JavaScript 这样的语言，Rust 并不会尝试自动地将非布尔值转换为布尔值.

变量绑定支持if表达式. `if ... else if ... else`: `if n <  30 {} else if n > 50 else {}`, `let x = if condition {} else {}`

rust循环支持while, loop, for...in, break/continue. 无限循环需使用loop; for常用于循环遍历集合, for 循环的安全性和简洁性使得它成为 Rust 中使用最多的循环结构.
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
    
    for _i in 2..n {
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

Range是标准库提供的类型，用来生成从一个数字开始到另一个数字之前结束的所有数字的序列, 例如`1..4`.

rust支持match, 其使用了模式匹配(pattern matching)技术, 支持绑定模式(bingding mode, 即使用操作符@将模式中的值绑定给一个变量, 供分支右侧的代码使用). Rust 中的匹配是 穷尽的（exhaustive）：必须穷举到最后的可能性来使代码有效. 通常match使用`_`来兜底. Rust 的模式匹配被广泛应用在状态机处理、消息处理和错误处理中.

一个 match 表达式由 分支（arms） 构成. 一个分支包含一个 模式（pattern）和表达式开头的值与分支模式相匹配时应该执行的代码. Rust 获取提供给 match 的值并挨个检查每个分支的模式. match 结构和模式是 Rust 中强大的功能，它体现了代码可能遇到的多种情形，并**确保没有遗漏处理, 否则无法通过编译**.

match表达式要求所有的分支都必须返回相同的类型,且如果是一个单独的match表达式而不是赋值给变量时，每个分支必须返回()类型.
rust提供if let和while let表达式, 分别在某些场合代替match表达式.

对非枚举类进行分支选择时必须注意处理例外情况, 即使在例外情况下没有任何要做的事. 例外情况用`_`表示.

> if let 是 match 的一个语法糖，它当值匹配某一模式时执行代码而忽略所有其他值. 在 if let 中可包含一个 else, else 块中的代码与 match 表达式中的 `_` 分支块中的代码相同.

> 因为 if 是一个表达式，我们可以在 let 语句的右侧使用它.

如果把枚举类附加属性定义成元组，在 match 块中需要临时指定一个名字:
```rust
enum Book {
    Papery(u32),
    Electronic {url: String},
}
let book = Book::Papery(1001);

match book {
    Book::Papery(i) => { // i就是临时名字
        println!("{}", i);
    },
    Book::Electronic { url } => {
        println!("{}", url);
    }
}
```

### 操作符
rust操作符优先级与类C语言(c/c++)类似.

对于原始固定长度的数组, 只有实现了Copy trait的类型可作为其元素, 即只有在栈上存放的元素才可以存放在该类型的数组中.
slice是对一个数组(包括固定大小数组和动态数组)的**引用**片段. 在底层，切片代表一个指向数组起始位置的指针和数组长度. 如果用`[T]`类型表示连续序列，那么切片类型就是`&[T]`和`&mut[T]`.

unsafe块: rust编译器将内存安全交由开发者自行负责.

### 泛型和trait
每一个编程语言都有高效处理重复概念的工具. 在 Rust 中其工具之一就是 泛型（generics）. 泛型是具体类型或其他属性的抽象替代.

泛型允许开发者编写一些在使用时才制定类型的代码. rust编译器会在编译期间自动为具体类型生成实现代码, 即采用单态化（monomorphization）实现.

在 Rust 里, 数据的行为通过 trait 来定义. 一般用 impl  关键字为数据结构实现 trait, 但 Rust 贴心地提供了派生宏（derive macro）, 可以大大简化一些标准接口的定义. 比如`#[derive(Debug)]` 为数据结构实现了 Debug trait, 提供了 debug 能力，这样可以通过`{:?}/{:#?}`(后者会按字段换行, 适合更多字段的数据结构), 用`println!`打印出来.

> Clone 让数据结构可以被复制，而 Copy 则让数据结构可以在参数传递的时候自动按字节拷贝.

>  Clone 类型的不可变引用转换为自有值(owned values)，即 &T -> T。Clone 没有对这种转换的效率做出承诺，所以它可能是缓慢和昂贵的。为了快速地在一个类型上实现 Clone, 可以使用派生宏.

> ToOwned 是 Clone 的一个更通用的版本. Clone 允许把一个 &T 变成一个 T，但 ToOwned 允许把一个 &Borrowed 变成一个 Owned，其中 `Owned: Borrow<Borrowed>`. 换句话说就是，不能把一个 &str 克隆成一个 String，或者把一个 &Path 克隆成一个 PathBuf，或者把一个 &OsStr 克隆成一个 OsString，因为 clone 方法签名不支持这种跨类型克隆，而这正是 ToOwned 的用途.

trait 告诉 Rust 编译器某个特定类型拥有可能与其他类型共享的功能. 可以通过 trait 以一种抽象的方式定义共享的行为. 可以使用 trait bounds 指定泛型是任何拥有特定行为的类型.

trait借鉴了Haskell的Typeclass, 是rust实现零成本抽象的基石, 其机制如下:
- trait是rust唯一的接口抽象方法
- 可以静态生成 也可以动态调用
- 可以当做标记类型拥有某些特定行为的"标签"来使用

rust支持trait的默认实现: 有时为 trait 中的某些或全部方法提供默认的行为，而不是在每个类型的每个实现中都定义自己的行为是很有用的, 但明确定义该方法可覆盖其默认实现.

> Clone 宏让数据结构可以被复制，而 Copy 则让数据结构可以在参数传递的时候自动按字节拷贝

```rust
// `<T: Debug>`表示有trait限定(trait bound)的泛型, 即只有实现了Debug trait的类型才适用. 只有实现了Debug trait的类型才拥有使用`{:?}`格式化打印的行为
fn match_opton<T: Debug>(o: Option<T>) {
    match o {
        ...
    }
}

struct Duck;
trait Fly {
    fn fly(&self) -> bool;
}
impl Fly for Duck{
    fn fly(&self) -> bool {
        return true;
    }
}
fn fly_static<T: Fly>(s: T) -> bool {
    s.fly()
}

fn fly_static(s: &Fly) -> bool {
    s.fly()
}

fn main(){
    let duck = Duck;
    fly_static::<Duck>(duck); // 静态分发, rust编译器会为`fly_static::<Duck>(duck)`这个具体类型的调用生成特殊化的代码. 即对编译器而言,该中抽象并不存在, 它在编译阶段即可将泛型展成具体类型的代码
    fly_dyn(&Duck) // 动态分发, 它会在运行时查找对应类型的方法, 需一定的运行时开销(很小). 与golang的接口类似
}
```

rust没有传统面向对象编程语言中的继承概念. 它通过trait将类型和行为明确地进行了区分, 充分贯彻了组合优于继承和面向接口编程的编程思想.

rust的trait符合c++之父提出的零开销原则: 如果不使用某个抽象, 就不用为它付出开销(静态分发); 如果确实需要使用该抽象, 可以保证这是开销最小的使用方式(动态分发).

trait 和 trait bound 让开发者使用泛型类型参数来减少重复，并仍然能够向编译器明确指定泛型类型需要拥有哪些行为. 因为向编译器提供了 trait bound 信息，它就可以检查代码中所用到的具体类型是否提供了正确的行为. 在动态类型语言中，如果尝试调用一个类型并没有实现的方法，会在运行时出现错误. Rust 将这些错误移动到了编译时，甚至在代码能够运行之前就强迫我们修复错误. 另外，开发者也无需编写运行时检查行为的代码，因为在编译时就已经检查过了，这样相比其他那些不愿放弃泛型灵活性的语言有更好的性能.

#### trait 作为参数
impl Trait 语法适用于直观的例子，它实际上是一种较长形式语法的`trait bound`语法糖.

```rust
pub fn notify(item: impl Summary) {
    println!("Breaking news! {}", item.summarize());
}

// trait bound形式
pub fn notify<T: Summary>(item: T) { // 同上
    println!("Breaking news! {}", item.summarize());
}

pub fn notify(item: impl Summary + Display) {}
pub fn notify<T: Summary + Display>(item: T) {} // 同上

fn some_function<T: Display + Clone, U: Clone + Debug>(t: T, u: U) -> i32 {}
fn some_function<T, U>(t: T, u: U) -> i32
    where T: Display + Clone,
          U: Clone + Debug
{} // 同上
```

#### 使用 trait bound 有条件地实现方法
- 通过使用带有 trait bound 的泛型参数的 impl 块，可以有条件地只为那些实现了特定 trait 的类型实现方法

    ```rust
    use std::fmt::Display;

    struct Pair<T> {
        x: T,
        y: T,
    }

    impl<T> Pair<T> {
        fn new(x: T, y: T) -> Self {
            Self {
                x,
                y,
            }
        }
    }

    impl<T: Display + PartialOrd> Pair<T> {
        fn cmp_display(&self) {
            if self.x >= self.y {
                println!("The largest member is x = {}", self.x);
            } else {
                println!("The largest member is y = {}", self.y);
            }
        }
    }
    ```
- 对任何实现了特定 trait 的类型有条件地实现 trait. 对任何满足特定 trait bound 的类型实现 trait 被称为 blanket implementations, 它们被广泛的用于 Rust 标准库中. 例如，标准库为任何实现了 Display trait 的类型实现了 ToString trait: `impl<T: Display> ToString for T {}`.


## 错误处理
Rust 有一套独特的处理异常情况的机制，它并不像其它语言中的 try 机制那样.

Rust 将错误组合成两个主要类别：可恢复错误（recoverable）和 不可恢复错误（unrecoverable）. 可恢复错误通常代表向用户报告错误和重试操作是合理的情况，比如未找到文件. 不可恢复错误通常是 bug 的同义词, 比如尝试访问超过数组结尾的位置.

Rust 并没有异常, 但是有可恢复错误 Result<T, E> 和不可恢复(遇到错误时停止程序执行)错误 panic!.

### panic
panic是线程级别的, 是安全的, 它不违反 Rust 的任何安全规则.

展开栈是默认的诧异行为，但在两种情况下 Rust 不会展开栈:
1. 如果在 Rust 展开第一个函数之后的清理期间 .drop() 方法触发了第二个诧异，那么这个panic会被认为是致命的. Rust 会停止展开并中止整个进程。
2. Rust 的panic行为是可自定义的。如果编译时加上 -C panic=abort，那么编译后程序中的第一个panic就会立即中止进程.

    编译时加上这个选项， Rust 则无须知道如何展开栈，因此能够减少编译后代码的大小.

> 当出现 panic 时，程序默认会开始 展开（unwinding），这意味着 Rust 会回溯栈并清理它遇到的每一个函数的数据，不过这个回溯并清理的过程有很多工作. 另一种选择是直接 终止（abort），这会不清理数据就退出程序。那么程序所使用的内存需要由操作系统来清理. 如果需要项目的最终二进制文件越小越好, panic 时通过在 Cargo.toml 的 [profile] 部分增加 panic = 'abort'，可以由展开切换为终止.

> 设置 RUST_BACKTRACE 环境变量来得到一个 backtrace. backtrace 是一个执行到目前位置所有被调用的函数的列表. Rust 的 backtrace 跟其他语言中的一样：阅读 backtrace 的关键是从头开始读直到发现你编写的文件, 这就是问题的发源地. 这一行往上是代码所调用的代码；往下则是调用该代码的代码.

> 为了获取backtrace信息, 必须启用 debug 标识. 当不使用 --release 参数运行 cargo build 或 cargo run 时 debug 标识会默认启用.

### Result
rust中的错误处理是通过返回Result<T, E>类型的方式进行的. Result<T, E>类型是Option<T>类型的升级版本.

```rust
fn main() -> Result<(), std::io::Error> {
    let f = File::open("bar.txt")?; // ?是一个错误处理的语法糖, 它会自动在出现错误的情况下返回std::io::Error.
}
```

Result<T, E> 类型定义了很多辅助方法来处理各种情况:
- is_ok/is_error: 返回 bool 值，表明result 是成功的结果还是错误的结果
- ok: 返回Option<T>类型的成功值（如果有的话）。如果result是一个成功的结果，就返回 Some(success_value)；否则，返回 None，而丢弃错误值
- err: 返回 Option<E> 类型的错误值（如果有的话）
- unwrap : 如果 Result 值是成员 Ok，unwrap 会返回 Ok 中的值; 如果 Result 是成员 Err，unwrap 会调用 panic!
- unwrap_or(fallback): 返回成功值，如果 result 是成功的结果的话。否则，它返回 fallback，丢弃错误值

    这是对 .ok() 的一个完美替代，因为返回类型是 T 而非 Option<T>. 当然，只有在存在适当后备值的情况下才可以使用这个方法
- unwrap_or_else(): 只是传入的不是后备值，而是一个函数或闭包。这个方法适合计算后备值如果用不上会造成浪费的情况。只有在返回错误结果时才会调用 fallback_fn
- expect : expect 与 unwrap 的使用方式一样. 但expect 在调用 panic! 时使用的错误信息将是传递给 expect 的参数，而不像 unwrap 那样使用默认的 panic! 信息
- as_ref: 将 Result<T, E> 转换为 Result<&T, &E>，即借用现有 result 中成功或错误值的引用
- result.as_mut() : 类似as_ref, 只是借用了可修改引用, 返回类型为 Result<&mut T, &mut E>

当编写一个其实现会调用一些可能会失败的操作的函数时，除了在这个函数中处理错误外，还可以选择让调用者知道这个错误并决定该如何处理, 这被称为 传播（propagating）错误，这样能更好的控制代码调用，因为比起你代码所拥有的上下文，调用者可能拥有更多信息或逻辑来决定应该如何处理错误.

这种传播错误的模式在 Rust 中很常见，以至于 Rust 提供了 ? 问号运算符来使其更易于处理. ? 运算符可被用于返回值类型为 Result 的函数.

`?`的实际作用是将 Result 类非异常的值直接取出, 如果有异常就将异常 Result 返回出去. 所以, `?`仅用于返回值类型为 Result<T, E> 的函数，其中 E 类型必须和 ? 所处理的 Result 的 E 类型一致.

example:
```rust
fn f(i: i32) -> Result<i32, bool> {
    if i >= 0 { Ok(i) }
    else { Err(false) }
}

fn g(i: i32) -> Result<i32, bool> {
    let t = f(i)?; // 出现error时, ?已将异常 Result 返回了
    Ok(t) // 因此确定 t 不是 Err, t 在这里已经是 i32 类型
}

fn g2(i: i32) -> Result<i32, bool> { // 等价于g
    let t = f(i);

    match t {
        Ok(i) => Ok(i),
        Err(b) => Err(b)
    }
}

fn main() {
    let r = g(10000);
    if let Ok(v) = r {
        println!("Ok: g(10000) = {}", v);
    } else {
        println!("Err");
    }
}
```

match 表达式与`?`运算符所做的有一点不同：? 运算符所收到的错误值被隐式传递给了 from 函数，它定义于标准库的 From trait 中，其用来将错误从一种类型转换为另一种类型. 当 ? 运算符调用 from 函数时，收到的错误类型被转换为由当前函数返回类型所指定的错误类型. 这在当函数返回单个错误类型来代表所有可能失败的方式时很有用，即使其可能会因很多种原因失败. 只要每一个错误类型都实现了 from 函数来定义如何将自身转换为返回的错误类型，? 运算符会自动处理这些转换.

?运算符消除了大量样板代码并使得函数的实现更简单, 甚至可以在 ? 之后直接使用链式方法调用来进一步缩短代码.

panic场景选择:
- 示例、代码原型和测试都非常适合 panic
- 当开发者比编译器知道更多的情况

    `let home: IpAddr = "127.0.0.1".parse().unwrap();`, 确定`127.0.0.1`是一个有效的 IP 地址, 无需处理`Result`

判断 Result 的 Err 类型，获取 Err 类型的函数是`kind()`:
```rust
use std::io;
use std::io::Read;
use std::fs::File;

fn read_text_from_file(path: &str) -> Result<String, io::Error> {
    let mut f = File::open(path)?;
    let mut s = String::new();
    f.read_to_string(&mut s)?;
    Ok(s)
}

fn main() {
    let str_file = read_text_from_file("hello.txt");
    match str_file {
        Ok(s) => println!("{}", s),
        Err(e) => {
            match e.kind() {
                io::ErrorKind::NotFound => {
                    println!("No such file");
                },
                _ => {
                    println!("Cannot read the file");
                }
            }
        }
    }
}
```

### `std::error::Error`
`std::error::Error`方法:
- err.description() : 返回 &str 类型的错误消息
- err.cause() : 返回一个 Option<&Error>，这是触发 err 的底层错误（如果有的话）

## 类型系统
在类型系统中, 一切皆类型. 基于类型定义的一系列组合,运算和转换等方法, 可以看做类型的行为.

类型系统优势:
1. 排查错误 : 静态语言可在编译期排查的类型错误
1. 抽象 : 抽象有利于强化编程规范和工程化系统
1. 文档 : 明确的类型可表明程序的行为
1. 优化效率 : 通过类型检查来优化部分操作, 节省运行时的时间
1. 类型安全
    - 类型安全的语言可避免类型间的无效操作. 比如`3/"hello"`
    - 类型安全的语言可保证内存安全, 避免诸如空指针, 虚悬指针和缓冲区溢出等内存安全问题
    - 类型安全的语言可避免语义上的逻辑错误.

类型系统分类:
- 静态类型: 在编译期进行类型检查的语言
- 动态类型: 在运行期进行类型检查的语言
- 强类型: 不允许类型的自动隐式转换, 在强制转换前不同类型无法操作.
- 弱类型: 与强类型相对.

静态类型的语言能在编译期对代码进行静态分析, 依靠的就是类型系统.

如果一个类型系统允许在一段代码在不同的上下文中具有不同的类型, 该类型系统就叫作多态类型系统.
现代编程语言有三种多态形式:
1. 参数化多态(parametric polymorphism)即泛型

    很多时候函数或数据类型都需要适用多种类型, 以避免大量的重复性工作. 泛型使得语言极具表达力, 同时也能保证静态类型安全.

    在编译期, 都会被单态化(monomorshization, 编译器进行静态分发的一种测量, 即通过填充编译时使用的具体类型，将通用代码转换为特定代码的过程).

    单态化静态分发的好处是性能好, 没有运行时开销, 缺点是容易造成编译后生成的二进制文件膨胀. rust采用该
1. Ad-hoc(ad-hoc polymorphism), 也叫特定多态, 指同一种行为定义在不同的上下文中会响应不同的行为实现.

    haskell使用Typeclass来支持ad-hoc, **rust使用trait来支持ad-hoc多态**.
    函数重载是ad-hoc多态.
1. 子类型多态(subtype polymorphism)

    一般用在面向对象语言中, 比如java. 它代表一种包含关系, 父类型的值包含了子类型的值, 所以子类型的值有时也可以看作父类型的值, 反之则不然.

多态的上下文中的方法解析过程被称为分发, 调用该方法被称为分发化（ dispatching）. 因此按多态发生的时间又可分为:
1. 静多态(static polymorphism), 也叫静态分发

    发生在编译期.
    参数化多态和ad-hoc多态一般是静多态.
    静多态牺牲灵活性获取性能.
1. 动多态(dynamic polymorphism), 也叫动态分发, 类似golang的接口

    发生在运行时.
    子类型多态一般是动多态.
    动多态牺牲性能获取灵活性.
    动多态在运行时需要查表, 占用较多空间, 因此一般情况下都使用静多态.

rust同时支持静多态(参数化多态+ad-hoc多态, 即泛型和trait)和动多态, 其静多态就是一种零成本抽象.

rust中一切皆表达式, 表达式皆有值, 值皆有类型 => rust中一切皆类型.

编程语言中不同类型的本质是内存空间和编码方式的不同.

rust没有gc, 内存首先由编译器分配, rust代码被编译成llvm ir, 其中就携带了内存分配信息. 所以编译器需要事先知道类型, 才好分配合理的内存.

rust多大部分是在编译期可确定内存大小的类型(sized type), 但也支持少量的动态大小类型(Dynamic sized type, DST), 比如str类型的字符串字面量; 以及零大小类型(zero sized type, zst), 比如单元类型和单元结构体, 同时由该类型组成的数组大小也为零. ZST类型的特点是, 它们的值就是其本身, 运行时并不占用内存空间.

> Rust 的核心语言中只有一种字符串类型：`str`, 字符串 slice(是一些储存在别处的 UTF-8 编码字符串数据的引用)，它通常以被借用的形式出现即`&str`. `&str`是引用类型(包含指针和长度信息), 存储在栈上, str字符串是存储在堆上.

动态大小类型(Dynamic sized type, dst)是指编译阶段无法确定占用空间的类型, 为了安全, 指向dst的指针一般是胖指针. 胖指针的设计, 避免了数组类型作为参数传递时自动退化为裸指针类型，丢失了长度信息的问题, 保证了类型安全.

对于 DST 类型, Rust 有如下限制:
- 只能通过指针来间接创建和操作 DST 类型, `＆[T] Box<[T]`可以, [T]不可以
- 局部变量和函数参数的类型不能是 DST 类型, 因为局部变量和函数参数必须在编译阶段知道它的大小, 因为目前unsized rvalue 功能还没有实现
- enum不能包 DST类型, struct 中只有最后一个元素 以是 DST, 其他成员不行, 如果包含有 DST 类型, 那么这个结构体也就成了 DST 类型.

never类型(低类型)表示永远不可能有返回值的类型, 特点:
1. 没有值
1. 是其他任意类型的子类型, 因此**可强制转换为其他任何类型**

如果说ZST类型表示"空", 那么底类型表示"无". rust中使用`!`表示底类型.

rust中有很多情况确实没有值, 但为了类型安全, 必须把这些情况纳入类型系统进行统一处理, 包括:
1. 发散函数(diverging function)

    返回类型是`!`, 表示该代码永远不会返回. 发散类型可转为任意一种类型.
1. continue和break
1. loop的无限循环
1. 空枚举, 比如`enum Void{}`

发散函数有:
1. 导致线程崩溃的panic!("....")以及基于它实现的各种函数/宏, 比如unimplemented!, unreachable!
1. 死循环loop{}
1. 用于退出函数的`std::process::exit`, 这类函数永远没有返回值.

rust和go类似, 只能对局部变量/全局变量进行类型推导, 而函数签名等场景下是不允许的, 这是有意为之. 同时rust使用`as`显式转换类型.

`xxx::<i32>()`该形式是泛型函数标注类型, `::<>`的形式称为turbofish操作符, 它通常用于在表达式中为泛型类型、函数或方法指定参数.

## 数据类型
很多编程语言中的数据类型是分为两类：
- 值类型

    一般是指可以将数据都保存在同一位置的类型, 例如数值、布尔值、结构体等都是值类型

    值类型有:
    - 原生类型
    - 结构体
    - 枚举体

- 引用类型

    会存在一个指向实际存储区的指针. 比如通常一些引用类型会将数据存储在堆中，而栈中只存放指向堆中数据的地址（指针）

    引用类型有：
    - 普通引用类型(裸指针), 比如go的指针
    - 原生指针类型, 比如go的slice, map, chan

    Rust 引用永远不为空. 没有跟 C 的 NULL 或 C++ 的 nullptr 对应的东西存在。引用没有默认的初始值（无论什么类型的变量，在其初始化之前都不能使用）。而且， Rust（在 unsafe代码外部）不会将整数转换为引用，因此不能把 0 转换成引用.

随着类型越来越丰富, 值类型和引用类型难以描述全部情况, rust所以引入了：
- 值语义（Value Semantic）

    复制以后，两个数据对象拥有的存储空间是独立的，互不影响

    基本的原生类型都是值语义，这些类型也被称为 POD（Plain old data）. POD 类型都是值语义，但是值语义类型并不一定都是 POD 类型

    具有值语义的原生类型，在其作为右值进行赋值操作时，**编译器会对其进行按位复制**

- 引用语义（Reference Semantic）

    复制以后，两个数据对象互为别名. 操作其中任意一个数据对象，则会影响另外一个.

    智能指针Box<T>封装了原生指针，是典型的引用类型. Box<T>无法实现 Copy，意味着它被 rust 标记为了引用语义，禁止按位复制.

    引用语义类型不能实现 Copy，但可以实现 Clone 的 clone 方法，以实现深复制.


在 Rust 中，可以通过是否实现 Copy trait 来区分数据类型的值语义和引用语义. 但为了更加精准，Rust 也引用了新的语义：复制（Copy）语义和移动（Move）语义
- Copy语义：对应值语义, 即实现了 Copy 的类型在进行按位复制时是安全的
- Move语义：对应引用语义, 在 Rust 中不允许按位复制，只允许移动所有权

实现了 Copy trait的作用: 实现 Copy trait 的类型同时拥有复制语义，在进行赋值或者传入函数等操作时，默认会进行按位复制.

对于默认可以安全的在栈上进行按位复制的类型，就只需要按位复制，也方便管理内存.
对于默认只可在堆上存储的数据，必须进行深度复制. 深度复制需要在堆内存中重新开辟空间，这会带来更多的性能开销.

#### 哪些实现了 Copy
- 结构体 ：当成员都是复制语义类型时，不会自动实现 Copy
- 枚举体 ：当成员都是复制语义类型时，不会自动实现 Copy
- 元组类型 ：本身实现了 Copy. 如果元素均为复制语义类型，则默认是按位复制，否则执行移动语义
- 字符串字面量 &str： 支持按位复制

结构体 && 枚举体：
1. 所有成员都是复制语义类型时，需要添加属性#[derive(Debug,Copy,Clone)]来实现 Copy
1. 如果有移动语义类型的成员，则无法实现 Copy

#### 哪些未实现 Copy
String ：to_string() 可以将字符串字面量转换为String

#### 哪些实现了 Copy trait
原生整数类型

对于实现 Copy 的类型，其 clone 方法只需要简单的实现按位复制即可.

#### 哪些未实现 Copy trait
- Box<T>

## 类型
Rust 支持类型推导，在编译器能够推导类型的情况下，变量类型一般可以省略，但常量（const）和静态变量（static）必须声明类型.

静态值不像常量那样是内联的，当读取和写入静态值时，需要用到 unsafe 代码块. 静态值通常与同步原语搭配使用，它们还用于实现全局锁定，以及与 C 程序库集成.

通常，如果不需要依赖静态的单例属性及其预定义的内存位置，而只需要其具体值，那么应该更倾向于使用常量。它们允许编译器进行更好的优化，并且更易于使用.

全局值只能在初始化时声明非动态的类型，并且在编译期，它在堆栈上的大小是已知的。但 lazy_static!宏，可用于初始化任何能够从程序
中的任何位置全局访问的动态类型。

使用 lazy_static!宏声明的元素需要实现 Sync 特征。这意味着如果某个静态值可变，那么必须使用诸如 Mutex 或 RwLock 这样的多线程类型，而不是 RefCell.

Rust 函数参数的类型和返回值的类型都必须显式定义, 如果没有返回值可以省略, 返回`unit即空元组`. 函数内部如果提前返回, 需要用 return 关键字, 否则最后一个表达式就是其返回值. 如果最后一个表达式后添加了`;`, 隐含其返回值为 unit.

### 标量类型(scalar type)数据
#### 整型
#### 浮点型
#### 布尔类型
#### 字符类型

char : 单个字符, 大小为四个字节(four bytes)，并代表了一个 Unicode 标量值（Unicode Scalar Value）.

处于内存安全考虑, rust将字符串分为两种:
- [str](https://github.com/rust-lang/rust/blob/master/library/core/src/str/mod.rs#L122)字符串(也叫字符串切片, 内置类型, 是DST类型), 固定长度的字符串, 通常以不可变借用的形式存在(`&str`)

    rust字符串因为包含长度, 因此不是以`\0`表示结束.

    `&str`字符串类型由两部分组成：
    1. 指向字符串序列的指针
    2. 记录长度的值

    `&str`存储于栈上, str字符串序列存储于程序的静态只读数据段, 栈或者堆内存中. `&str`是一种胖指针.

    > str 是编译器能够识别的内置类型, 表示有限但大小未知的 UTF-8 编码的连续字节序列, 因此str在堆上, 并且不属于标准库.
- String字符串, 长度可变的字符串, 在堆上分配, 与`&str`的主要区别是它有管理内存空间的能力, 而`&str`没有.

    创建方法: `let s = String::from("hello")`, 两个冒号（::）是运算符，允许将特定的 from 函数置于 String 类型的命名空间（namespace）下，而不需要使用类似 string_from 这样的名字.

    String类型本质是一个成员变量为`Vec<u8>`类型的结构体，所以它是直接将字符内容存放于堆中的.

    String类型由三部分组成：
    - 执行堆中字节序列的指针（as_ptr方法）
    - 记录堆中字节序列的字节长度（len方法）
    - 堆分配的容量（capacity方法）

    基础类型变String:
    ```rust
    let one = 1.to_string();         // 整数到字符串
    let float = 1.3.to_string();     // 浮点数到字符串
    let slice = "slice".to_string(); // 字符串切片到字符串
    ```

    其他用法:
    ```rust
    fn main() {
        let s = String::from("hello中文");
        for c in s.chars() { // 等同go的 `range []rune(s)`
            println!("{}", c);
        }
    }
    ```

`& 'static str`是静态生命周期字符串. 静态生命周期即程序生命周期.

rust字符串的本质是一段有效的utf8字节序列.

Rust中的字符串不能使用索引访问其中的字符，但可以通过bytes和chars两个方法来分别返回按字节和按字符迭代的迭代器.

Rust提供了另外两种方法：get和get_mut来通过指定索引范围来获取字符串切片.

#### 引用/指针
引用是用`&`和`& mut`操作符来创建, 受Rust的安全检查规则的限制.

引用是Rust提供的一种无所有权的指针语义, 引用的生命期不能超过其引用的资源. 引用是基于指针的实现，它与指针的区别是：指针保存的是其指向内存的地址，而引用可以看做某块内存的别名（Alias）.

裸指针(原生指针):`*const T和*mut T`, 可以在unsafe块下任意使用, 不受Rust的安全检查规则的限制.

智能指针实际上是一种结构体，只是行为类似指针, 智能指针是对指针的一层封装，提供了一些额外的功能，比如自动释放堆内存.

> rust在堆上分配内存的唯一方法是通过智能指针类型.

> rustc自身使用jemalloc,而其构建的lib和bin使用系统内存分配器.

智能指针区别于常规结构体的特性在于：它实现了Deref和Drop这两个trait:
- Deref：提供了解引用能力
- Drop：提供了自动析构的能力

智能指针拥有资源的所有权，而普通引用只是对所有权的借用.

rust指针包括 引用(reference), 原生指针(raw pointer), 函数指针(fn pointer)和智能指针(smart pointer).
引用的本质是非空指针.

解引用deref: 解引用会获得所有权, 解引用操作符是`*`.

### 常用集合类型
Rust 标准库中包含一系列被称为 集合（collections）的非常有用的数据结构. 大部分其他数据类型都代表一个特定的值，不过集合可以包含多个值. 不同于内建的数组和元组类型，这些集合指向的数据是储存在堆上的，这意味着数据的数量不必在编译时就已知，并且还可以随着程序的运行增长或缩小.

std::collections提供了4种通用集合类型:
1. 线性序列: 向量(Vec, **堆上分配**的数组), 双端队列(VecDeque), 链表(LinkedList)

    向量也是一种数组, 但可动态增长.
    `vec!`是一个宏, 用来创建向量字面量.
    rust的VecDeque是基于可增长的RingBuffer算法实现的双端队列.
    通常最好使用Vec或VecDeque类型, 因为它们比链表更加快速, 内存访问效率更高, 并且可以更好地利用cpu缓存.

    > rust的String基于Vec, 可变字符串用String, 不可变字符串用`&str`.

    > 双端队列（Double-ended Queue，缩写Deque）是一种同时具有队列（先进先出）和栈（后进先出）性质的数据结构. 双端队列中的元素可以从两端弹出，插入和删除操作被限定在队列的两端进行.

    > Rust提供的链表是双向链表，允许在任意一端插入或弹出元素.

    > 当缓冲区达到其容量上限后，再给向量添加元素会导致一系列操作：分配一个更大的缓冲区，将现有内容复制过去，基于新缓冲区更新向量的指针和容量，最后释放旧缓冲区.

    ```rust
    // 声明
    let mut v2 = Vec::new();
    v2.push(2); // 现在 v2 的类型是 Vec<i32>
    let v3 = Vec::<u8>::new(); // 使用 turbofish 符号, 泛型函数中的 turbofish 运算符出现在函数名之后和圆括号之前

    let mut v: Vec<i32> = Vec::new();
    // let v = vec![1, 2, 3]; // 使用vec!宏, 让rust自动推导类型
    v.push(5); // 在 vector 的结尾增加新元素时，在没有足够空间将所有所有元素依次相邻存放的情况下，可能会要求分配新内存并将老的元素拷贝到新的空间中

    let vector = vec![1, 2, 4, 8];     // 通过数组创建向量

    let third: &i32 = &v[2]; // 当引用一个不存在的元素时 Rust 会造成 panic
    println!("The third element is {}", third);

    match v.get(2) { // 2 is index. 当 get 方法被传递了一个数组外的索引时，它不会 panic 而是返回 None
        Some(third) => println!("The third element is {}", third),
        None => println!("There is no third element."),
    }

    for i in &v {
        println!("{}", i);
    }
    ```
1. Key-Value映射: 无序哈希表(HashMap), 有序哈希表(BTreeMap)

    HashMap要求key是必须可哈希的类型，BTreeMap的key必须是可排序的
    value必须是在编译期已知大小的类型

    ```rust
    use std::collections::HashMap;
    let mut scores = HashMap::new();

    let teams  = vec![String::from("Blue"), String::from("Yellow")];
    let initial_scores = vec![10, 50];

    let scores: HashMap<_, _> = teams.iter().zip(initial_scores.iter()).collect(); // 使用下划线占位是因为 Rust 能够根据 vector 中数据的类型推断出 HashMap 所包含的类型

    let team_name = String::from("Blue");
    let score = scores.get(&team_name);

    for (key, value) in &scores {
        println!("{}: {}", key, value);
    }

    scores.insert(String::from("Blue"), 11); // 覆盖value
    scores.entry(String::from("Yellow")).or_insert(50); // entry 函数的返回值是一个枚举，Entry，它代表了可能存在也可能不存在的值. Entry 的 or_insert 方法在键对应的值存在时就返回这个值的可变引用，如果不存在则将参数作为新值插入并返回新值的可变引用.

    let mut map = HashMap::new();
    map.insert(1, "a");
   
    if let Some(x) = map.get_mut(&1) {
        *x = "b";
    }
    ```

    > 对于像 i32 这样的实现了 Copy trait 的类型，其值可以拷贝进哈希 map。对于像 String 这样拥有所有权的值，其值将被移动而哈希 map 会成为这些值的所有者

    > 如果将值的引用插入哈希 map，这些值本身将不会被移动进哈希 map。但是这些引用指向的值必须至少在哈希 map 有效时也是有效的.

    HashMap 默认使用一种 “密码学安全的”（“cryptographically strong” ）哈希函数，它可以抵抗拒绝服务（Denial of Service, DoS）攻击. 然而这并不是可用的最快的算法，不过为了更高的安全性值得付出一些性能的代价. 如果性能监测显示此哈希函数非常慢，以致于无法接受，可以指定一个不同的 hasher 来切换为其它函数. hasher 是一个实现了 BuildHasher trait 的类型.
1. 集合类型: 无序集合(HashSet), 有序集合(BTreeMap)

    集合类型实际就是把Key-Value映射的Value设置成空元组.
1. 优先队列: 基于二叉最大堆(BinaryHeap)实现.

无论是Vec还是HashMap，使用这些集合容器类型，最重要的是理解容量（Capacity）和大小（Size/Len）:
- 容量是指为集合容器分配的内存容量
- 大小是指集合中包含的元素数量

### 方法
不过方法与函数是不同的, 因为它们在结构体的上下文中被定义(或者是枚举或 trait 对象的上下文).

impl 块的另一个有用的功能是: 允许在 impl 块中定义不以 self 作为参数的函数, 这被称为 关联函数（associated functions, 在主流的编程语言中, 这也被称为静态方法）, 因为它们与结构体相关联. 它们仍是函数而不是方法, 因为它们并不作用于一个结构体的实例. 比如`String::from`关联函数. 它类似于面向对象编程语言中的静态方法. 这些方法在类型自身上即可调用, 并且不需要类型的实例来调用, 调用方法是`<类型名>::<函数名>`

> 结构体允许拥有多个 impl 块.

impl实例方法的变体, 根据限制由少到多排列的:
- &self 作为第一个参数, 此方法仅提供对类型实例的读取访问权限
- &mut self 作为第一个参数, 此方法提供对类型实例的可变访问
- self 作为第一个参数, 这些方法拥有调用它的实例的所有权，并且类型在后续调用时将失效

### enumerations、enums枚举
> Rust 的枚举与 F#、OCaml 和 Haskell 这样的函数式编程语言中的 代数数据类型（algebraic data types）最为相似

### 泛型
泛型就是把一个泛化的类型作为参数.

> C++ 语言中用"模板"来实现泛型. 泛型机制是编程语言用于表达类型抽象的机制，一般用于功能确定、数据类型待定的类，如链表、映射表等.

与枚举类型额函数一样,结构体名称旁边的`<T>`叫做泛型声明. 泛型只有被声明之后才可以实现.

```rust
struct Point<T> { x: T, y: T}
impl<T> Point<T> {
    fn new(x: T,y: T) -> Self {
        Point{x:x, y:y}
    }
}
fn main() {
    let p1 = Point::new(1,2);
    let p1 = Point::new("1","2");   
}
```

### trait
trait定义了一组类型可以选择性实现的`契约`或共享行为.

> trait和泛型通过单态化（静多态）或运行时多态（动多态）提供了两种代码复用的方式.

接口是一个软件系统开发的核心部分，它反映了系统的设计者对系统的抽象理解。作为一个抽象层，接口将使用方和实现方隔离开来，使两者不直接有依赖关系，大大提高了复用性和扩展性.

很多编程语言都有接口的概念，允许开发者面向接口设计，比如 Java 的 interface、Elixir 的 behaviour、Swift 的 protocol 和 Rust 的 trait.

当在运行期使用接口来引用具体类型的时候，代码就具备了运行时多态的能力. 在运行时，一旦使用了关于接口的引用，变量原本的类型被掩藏，此时需要一个胖指针获取这个引用具备什么样的能力, 即在生成这个引用的时候，需要构建胖指针，除了指向数据本身外，还需要指向一张涵盖了这个接口所支持方法的列表, 这个列表，就是熟知的虚表（virtual table）.

trait是在行为上对类型的约束即对类型行为的抽象, 是Rust实现零成本抽象的基石，它有如下机制：
- trait是Rust唯一的接口抽象方式
- 可以静态分发，也可以动态分发
- 可以当做标记类型拥有某些特定行为的"标签"来使用

impl格式为`impl <特性名> for <所实现的类型名>`:
```rust
trait Descriptive {
    fn describe(&self) -> String;
}

struct Person {
    name: String,
    age: u8
}

impl Descriptive for Person {
    fn describe(&self) -> String {
        format!("{} {}", self.name, self.age)
    }
}
```

Rust 同一个类可以实现多个特性，每个 impl 块只能实现一个.

默认trait: 接口只能规范方法而不能定义方法，但特性可以定义方法作为默认方法，因为是"默认"，所以对象既可以重新定义方法，也可以不重新定义方法使用默认的方法. 这是trait与接口的不同点

trait有4中用法:
- 接口抽象

    接口是对类型行为的统一约束, 是trait最基础的用法, 特点:
    - 接口中定义方法, 并支持默认实现
    - 接口中不能实现另一个接口, 但是接口间可以继承
    - 同一个接口可以同时被多个类型实现, 但不能被同一个类型实现多次
    - 使用trait关键字来定义接口
    - 使用impl关键字为类型实现接口方法

    ## 关联类型
    关联类型（associated types）是一个将类型占位符与 trait 相关联的方式，这样 trait 的方法签名中就可以使用这些占位符类型. trait 的实现者会针对特定的实现在这个类型的位置指定相应的具体类型.
    **rust很多操作符都是基于trait来实现的.**, 比如加法操作符:
    ```rust
    pub trait Add<RHS = Self> { // `Add<RHS = Self>`表示参数类型RHS为Self, Self是每个trait都带有的隐式类型参数, 代表实现当前trait的具体类型.
        type Output; // 关联类型. Output是一个占位类型, trait的实现者会指定 Output的具体类型.
        fn add(self, rhs: RHS) -> Self::Output;
    }

    1+2 // =>`1.add(2)` 代码中出现`+`时, rust就会自动调用操作符左侧的操作数对应的add()方法去实现具体的操作, 即`+`操作与调用`add()`等价.
    // 可看rust源码为u32实现的Add trait(用宏完成的).
    ```

    在语义层面上, 使用关联类型增强了trait表示行为的语义, 因为它表示了和某个行为(trait)相关联的类型. 在工程上, 也体现出了高内聚的特点.

    ## trait一致性
    孤儿原则(Orphan Rule): impl块要么与trait的声明在同一个crate中, 要么与类型的声明在同一个crate中. 其可限制代码被破坏性改写, 导致出现难以预料的bug.

    同理匿名impl必须与类型本身在同一个crate中.

    ## trait继承
    rust不支持传统面向对象的继承, 但支持trait继承.

    ```rust
    trait Paginatge: Page + PerPage { // 用冒号表示继承其他trait.
        ...
    }
    impl <T: Page + PerPage>Paginate for T{

    }
    ```
    
- 泛型约束

    trait限定(trait bound) : 泛型的行为被trait限定在更有限的范围内, 多个trait限定用`+`连接.

    ```rust
    fn notify(item: impl Summary + Display) // 使用 impl 特征语法`impl Summary + Display`
    fn notify<T: Summary + Display>(item: T) // 等价于同上

    fn some_function<T: Display + Clone, U: Clone + Debug>(t: T, u: U)
    fn some_function<T, U>(t: T, u: U) -> i32 // 等价于同上
    where T: Display + Clone,
          U: Clone + Debug
    ```

    ```rust
    use std::ops::Add;
    // 表示sum函数的参数必须实现Add trait
    fn sum<T: Add<T, Output=T>>(a: T, b:T) -> T{ // Add<T, Output=T>: Add表示和T类型相加, 并将产生T类型作为输出即输入输出具有相同类型. T 是实现此特征的类型别名
        a+b
    }
    fn foo<T, K, R>(a: T, b:K, c:R) where T: A, K:B+C, R:D {...}// where用于为泛型增加较多的trait限定, 提高代码可读性.

    #[derive(Default, Debug, PartialEq, Copy, Clone)]
    struct Complex<T> {
        //实部
        re: T,
        //虚部
        im: T
    }
    impl<T> Complex<T> {
        fn new(re: T, im: T) -> Self {
            Complex { re, im }
        }
    }

    impl<T: Add<T, Output=T>> Add for Complex<T> { // T:Add 表示必须实现 Add trait. 如果没有实现, 那么无法对`T`使用`+`运算符
        type Output = Complex<T>;
        fn add(self, rhs: Complex<T>) -> Self::Output {
            Complex { re: self.re + rhs.re, im: self.im + rhs.im }
        }
    }

    let a = Complex::new(1,-2);
    let b = Complex::default();
    let res = a + b;
    ```

    类型上的泛型约束:
    ```rust
    use std::fmt::Display;
    struct Foo<T: Display> {
        bar: T
    }

    struct Bar<F> where F: Display {
        inner: F
    }
    ```

    不推荐在类型上使用泛型约束， 因为它对类型自身施加了限制. 通常, 希望类型尽可能是泛型，从而允许使用任何类型创建实例.

    trait限定给予了开发者更大的自由度, 因为不再需要类型间的继承, 也简化了编译器的检查操作. 包含trait限定的泛型属于**静态分发**, 在编译期通过单态化分别生成具体类型的实例, 所以调用trait限定中的方法也都是运行时零成本的, 因为不需要在运行时再进行方法查找.
- 抽象类型

    在运行时作为一种间接的抽象类型来使用, 动态地分发给具体的类型.
- 标签trait

    对类型的约束, 即直接作为一种"标签"使用.

    Sized 特征是一种标记性特征，用于表示编译期已知大小的类型, 并且可以将该类型的实例放在栈上.

rust规定函数在参数传递, 返回值传递中类型必须是编译阶段可确定大小的, 而trait的大小在编译时是不固定的, 因此它无法作为实例变量, 参数, 返回值, 这与go的interface不同.

特性做返回值:
```rust
// 会报错: 特性做返回值只接受实现了该特性的对象做返回值且在同一个函数中所有可能的返回值类型必须完全一样, 比如即使结构体 A 与结构体 B 都实现了Descriptive
fn some_function(bl: bool) -> impl Descriptive {
    if bl {
        return A {};
    } else {
        return B {};
    }
}
```

有条件实现方法:
```rust
struct A<T> {}

impl<T: B + C> A<T> { // 声明了 A<T> 类型必须在 T 已经实现 B 和 C 特性的前提下才能有效实现此 impl 块
    fn d(&self) {}
}
```

### 常见trait
- std::ops 模块的 Add : 允许使用`+`运算符将两个复数相加。
- std::convert 模块的 Into 和 From : 使用户能够根据其他类型创建复数类型。
- Display : 使用户能够输出人类可读版本的复数类型

### 宏
宏语句可以使用圆括号, 中括号, 花括号, 一般使用中括号表示数组.

### 语法
字符串 slice（string slice）是 String 中一部分值的引用. 字符串 slice range 的索引必须位于有效的 UTF-8 字符边界内，如果尝试从一个多字节字符的中间位置创建字符串 slice，则程序将会因错误而退出.
字符串字面值就是 slice, 是一个指向二进制程序特定位置的 slice.
通用slice跟字符串 slice 的工作方式一样，通过存储第一个集合元素的引用和一个集合总长度.

#### 函数
在Rust中，函数定义以fn关键字开始并紧随函数名称(按小写字母以下划线分割)与一对圆括号，另外还有一对花括号用于标识函数体开始和结尾的地方. 使用函数名加圆括号的方式来调用函数.

> 函数体由一系列的语句和一个可选的结尾表达式构成.

> Rust不关心在何处定义函数，只要这些定义对于使用区域是可见的即可, 这与go相同, 与c不同.

在函数签名中，必须声明每个参数的类型. 函数不支持自动返回值类型判断, 如果没有明确声明函数返回值的类型, 函数将被认为是"纯过程"即不允许产生返回值, 因此return 后面不能有返回值表达式.

> **参数变量和传入的具体参数值有自己分别对应的名称parameter和argument, 即形参和实参**.

rust并不对返回值命名，但要在`箭头（->）`后声明它的类型.
**在 Rust 中，函数的返回值等同于函数体最后一个表达式的值, 因此rust允许省略return**.
rust不支持多返回值, 但可以利用元组来返回多个值.

```rust
fn main() {
    let x = plus_one(5);
}
fn plus_one(x: i32) -> i32 {
    x + 1
}
// 错误: 在包含 x + 1 的行尾加上一个分号，把它从表达式变成语句. 语句并不会返回值，使用空元组 () 表示不返回值. 因为不返回值与函数定义相矛盾，从而出现一个错误.
// fn plus_one(x: i32) -> i32 {
//     x + 1;
// }
```

函数的第一个参数如果是Self相关的类型, 且命名为`self`, 那么这个参数就是receiver. 有receiver的函数即为方法(method), 用`变量实例.方法名`来调用.
没有receiver参数的函数是静态函数(static function), 通过类型加`::`的方式调用.

#### 所有权(ownership)
所有权（系统）是 Rust 最为与众不同的特性，它让 Rust 无需垃圾回收（garbage collector）即可保障内存安全. **所有权的存在就是为了管理堆数据**.

所有运行的程序都必须管理其使用计算机内存的方式. 一些语言中具有垃圾回收机制(go, java等)，在程序运行时不断地寻找不再使用的内存；在另一些语言中(c, c++等)，开发者必须亲自分配和释放内存. Rust 则选择了第三种方式：通过所有权系统管理内存，编译器在**编译时**会根据一系列的规则进行检查, 在运行时，所有权系统的任何功能都不会减慢程序.

所有权本质上就是在语言层面禁止了同一个可变数据会有多个变量引用的情况, 一旦作为参数传递了, 就会发生所有权的移动（Move）或借用（Borrow）, 即从根本上杜绝了并发情景下的数据共享冲突.

> 变量范围是变量的一个属性, 其代表变量的可行域, 默认从声明变量开始有效直到变量所在域结束.

> Rust 之所以没有明示释放变量的步骤是因为在变量范围结束的时候, 其编译器自动添加了调用释放资源函数的步骤.

> rust的栈中的所有数据都必须占用已知且固定的大小, 在编译时大小未知或大小可能变化的数据会分配到堆上.

Rust中分配的每块内存都有其所有者，所有者负责该内存的释放和读写权限，并且每次每个值只能有唯一的所有者.

在进行赋值操作时，对于可以实现Copy的复制语义类型，所有权并未改变. 对于复合类型来说，是复制还是移动，取决于其成员的类型.

所有权的规则:
- Rust 中的每一个值都有一个被称为其所有者（owner）的变量
- 值在任一时刻有且只有一个所有者
- 当所有者（变量）离开作用域，这个值将被丢弃

    内存在拥有它的变量离开作用域后就被自动释放(调用`drop`方法)

    在 C++ 中，这种 item 在生命周期结束时释放资源的模式有时被称作 资源获取即初始化（Resource Acquisition Is Initialization (RAII)）, 即本质是将资源的生命周期和对象的生命周期绑定.

    > 在推断所有权规则时，作用域是一个非常重要的属性, 它也会被用来推断借用和生命周期

这三条规则很好理解, 核心就是保证单一所有权. 其中第二条规则讲的所有权转移是 Move 语义, Rust 从 C++ 那里学习和借鉴了这个概念.

rust 内存回收策略：内存在拥有它的变量离开作用域后就被自动释放.

所有权规则解决了谁真正拥有数据的生杀大权问题，让堆上数据的多重引用不复存在，这是它最大的优势.

```rust
#[derive(Debug)]
struct Foo(u32);

fn main() {
    let foo = Foo(2048);
    let bar = foo;
    println!("Foo is {:?}", foo); // 报错: foo已被move给bar
    println!("Bar is {:?}", bar);
}
```

```rust
fn main() {
    let foo = 4623;
    let bar = foo;
    println!("{:?} {:?}", foo, bar); // 正常
}
```

```rust
#[derive(Copy, Clone, Debug)]
struct Dummy;

fn main() {
    let a = Dummy;
    let b = a;
    println!("{}", a); // 正常, 因此Dummy有实现Copy
    println!("{}", b);
}
```

所有权也让代码变得复杂, 如果要**避免所有权转移之后不能访问**的情况, 一般需要手动复制, 编写麻烦且效率不高, 因此rust提供了两种解决方法:
1. Rust 提供了 Copy 语义. 如果一个数据结构实现了 Copy trait，那么它就会使用 Copy 语义. 这样, **在赋值或者传参**时, 值会自动**按位拷贝**（浅拷贝）.

    Copy: 默认情况下，通过变量分配或访问，以及从函数返回时复制的值（例如按位复制）具有复制语义. 这意味着该值可以使用任意次数，每个值都是全新的. 

    等同于go的值传递.

    > Clone是Copy的父trait.

    > 默认情况下, C++具有复制语义. 后来的C++ 11 版本提供了对移动语义的支持.

1. 无法使用 Copy 语义，那可以`借用`数据

变量与数据交互的方式:
1. move

    > 当以转移所有权的方式给函数传参时，称其为传值（by value）。如果传给函数的是对值的引用，则称其为传引用（by reference）.

    在 Rust 中，对多数类型而言，给变量赋值、给函数传值或从函数返回值这样的操作不会复制值, 而是转移（move）值.

    移动: 通过变量访问或重新分配给变量时移动到接收项的值表示移动语义. 由于Rust 的仿射类型系统，它默认会采用移动语义. 仿射类型系统的一个突出特点是值或资源只能使用一次，而 Rust 通过所有权规则展示此属性.

    默认情况下，所有类型都有“移动语义”，但是一旦一个类型实现了 `Copy'，它就会得到`复制语义`.

    设计选择: **Rust 永远也不会自动创建数据的 “深拷贝”. 因此，任何 自动 的复制可以被认为对运行时性能影响较小.**

    仅在栈中的基本数据类型的数据的"移动"方式是直接复制, 这不会花费更长的时间或更多的存储空间, "基本数据"类型有这些：

    - 所有整数类型: i32 、 u32 、 i64 等
    - 布尔类型 bool: true 或 false
    - 所有浮点类型: f32 和 f64
    - 字符类型 char
    - 仅包含以上类型数据的元组(Tuples)

    Rust会尽可能地降低程序的运行成本, 所以默认情况下, 长度较大的数据存放在堆中, 且采用移动的方式进行数据交互.

    > 在 match 表达式中，移动类型默认也会被移动.
1. clone

    克隆仅在需要复制的情况下使用, 毕竟复制数据会花费更多的时间.
    由于其运行时消耗，许多 Rustacean 之间有一个趋势是倾向于避免使用 clone 来解决所有权问题.

Copy是一种自动化特征，大多数堆栈上的数据类型都自动实现了它. 它通常用于数据完全在栈上的变量的复制, 比如基元类型和不可变引用(&T); 否则Copy开销很大, 因为它需要从堆中复制数据. Copy 特征复制类型的方式与 C 语言中的 memcpy 函数类似，后者用于按位复制值。默认情况下不会为自定义类型实现 Copy 特征，因为 Rust 希望显式指定复制操作，并且要求开发人员必须选择实现该特征. 没有实现 Copy 特征的类型包括`Vec<T>、 String 和可变引用`等.

> 有一条经验规则，就是任何在值被清除后需要特殊处理的类型都不能是 Copy 类型。比如，Vec 需要释放其元素， File 需要关闭其文件勾柄，而 MutexGuard 需要解锁其互斥量。

> 用户定义的类型默认属于非 Copy 类型。如果自定义结构体的所有字段本身都是 Copy 类型，那可以在定义上方添加 #[derive(Copy, Clone)] 属性把这个类型标注成 Copy 类型;对于并非所有字段都是 Copy 类型的结构体，就算加这个属性也不管用, 编译器会报错.

Clone 特征用于显式复制, 并附带 clone 方法， 类型可以实现该方法以获取自身的副本.

Clone 与 Copy 的不同之处在于, 其中的赋值操作时Copy是隐式复制值, 要复制 Clone值，就必须显式调用 clone 方法. clone 方法是一种更通用的复制机制, Copy 是它的一个
特例，即总是按位复制.

String 和 Vec 这类元素很难进行复制, 只实现了 Clone 特征. 智能指针类型也实现了 Clone 特征, 它只是在指向堆上相同数据的同时复制指针和额外的元数据（例如引用
计数）.

涉及函数入参的所有权机制:
```rust
fn main() {
    let s = String::from("hello");
    // s 被声明有效

    takes_ownership(s);
    // s 的值被当作参数传入函数
    // 所以可以当作 s 已经被移动，从这里开始已经无效

    let x = 5;
    // x 被声明有效

    makes_copy(x);
    // x 的值被当作参数传入函数
    // 但 x 是基本类型，依然有效
    // 在这里依然可以使用 x 却不能使用 s

} // 函数结束, x 无效, 然后是 s. 但 s 已被移动, 所以不用被释放


fn takes_ownership(some_string: String) {
    // 一个 String 参数 some_string 传入，有效
    println!("{}", some_string);
} // 函数结束, 参数 some_string 在这里释放

fn makes_copy(some_integer: i32) {
    // 一个 i32 参数 some_integer 传入，有效
    println!("{}", some_integer);
} // 函数结束, 参数 some_integer 是基本类型, 无需释放
```

涉及函数返回值的所有权机制:
```rust
fn main() {
    let s1 = gives_ownership();
    // gives_ownership 移动它的返回值到 s1

    let s2 = String::from("hello");
    // s2 被声明有效

    let s3 = takes_and_gives_back(s2);
    // s2 被当作参数移动, s3 获得返回值所有权
} // s3 无效被释放, s2 被移动, s1 无效被释放.

fn gives_ownership() -> String {
    let some_string = String::from("hello");
    // some_string 被声明有效

    return some_string;
    // some_string 被当作返回值移动出函数
}

fn takes_and_gives_back(a_string: String) -> String {
    // a_string 被声明有效

    a_string  // a_string 被当作返回值移出函数
}
```

Rust 有一个叫做 Copy trait(类似深拷贝)的特殊注解，可以用在类似整型这样的**存储在栈上的类型上. 如果一个类型拥有 Copy trait，一个旧的变量在将其赋值给其他变量后仍然可用**. Rust 不允许自身或其任何部分实现了 Drop trait 的类型使用 Copy trait.

作为一个通用的规则，任何简单标量值的组合可以是 Copy 的，不需要分配内存或某种形式资源的类型是 Copy 的. 默认支持 Copy 的类型有：
- 所有整数类型，比如 u32
- 布尔类型，bool，它的值是 true 和 false
- 所有浮点数类型，比如 f64
- 字符类型, char
- 元组，当且仅当其包含的类型也都是 Copy 的时候. 比如 (i32, i32) 是 Copy 的，但 (i32, String) 就不是.

将值传递给函数在语义上与给变量赋值相似: 向函数传递值可能会移动或者复制，就像赋值语句一样. 返回值也可以转移所有权.

变量的所有权总是遵循相同的模式：将值赋给另一个变量时**移动**它 . 当持有堆中数据值的变量离开作用域时，其值将通过 drop 被清理掉，除非数据被移动为另一个变量所有.

引用(`&`)语法可创建一个 指向 值 的引用，但是**并不拥有它**, **因为并不拥有这个值，当引用离开作用域时其指向的值也不会被丢弃**. 实质上"引用"是变量的间接访问方式, 且"引用"并没有在栈中复制变量的值. 引用本身也是一个类型并具有一个值，这个值记录的是别的值所在的位置，但引用不具有所指值的所有权.

引用举例:
```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = &s1; // s2是s1的引用`&String`
    s2.push_str("oob"); // 错误，禁止修改租借的值
    let s3 = s1; //  报错. s2 租借的 s1 已经将所有权移动到 s3，所以 s2 将无法继续租借使用 s1 的所有权。如果需要使用 s2 使用该值，必须重新租借
    println!("{}", s2);
}

// 改为:
fn main() {
    let s1 = String::from("hello");
    let mut s2 = &s1;
    let s3 = s1;
    s2 = &s3; // 使用前重新从 s3 租借所有权
    println!("{}", s2);
}

fn main() {
    let mut s1 = String::from("run");
    // s1 是可变的

    let s2 = &mut s1;
    // s2 是可变的引用

    s2.push_str("oob");
    println!("{}", s2);
}
```

可变引用与不可变引用相比除了权限不同以外, **可变引用不允许多重引用**, 但不可变引用可以.

Rust 对可变引用的这种设计主要出于对并发状态下发生数据访问碰撞的考虑，在编译阶段就避免了这种事情的发生. 由于发生数据访问碰撞的必要条件之一是数据被至少一个使用者写且同时被至少一个其他使用者读或写，所以在一个值被可变引用时不允许再次被任何引用.

垂悬引用（Dangling References）: 没有实际指向一个真正能访问的数据的指针（注意: 不一定是空指针，还有可能是已经释放的资源）, 比如:
```rust
fn main() {
    // 伴随着 dangle 函数的结束，其局部变量的值本身没有被当作返回值, 被释放了, 但它的引用却被返回, 这个引用所指向的值已经不能确定的存在, 因此会报错
    let reference_to_nothing = dangle();
}

fn dangle() -> &String {
    let s = String::from("hello");

    &s
}
```

**获取引用作为函数参数称为`借用(borrowing)`**.

引用默认不允许修改引用的值, 允许可变引用(`&mut`), 但可变引用有一个很大的限制：**在特定作用域中的特定数据有且只有一个可变引用. 这个限制的好处是 Rust 可以在编译时就避免数据竞争.数据竞争（data race）类似于竞态条件**，它可由这三个行为造成：
- 两个或更多指针同时访问同一数据
- 至少有一个指针被用来写入数据
- 没有同步数据访问的机制

数据竞争会导致未定义行为，难以在运行时追踪，并且难以诊断和修复；Rust 避免了这种情况的发生，因为它甚至不会编译存在数据竞争的代码！

引用的规则:
1. 在任意给定时间，要么只能有一个可变引用，要么只能有多个不可变引用. 即任意时刻不能在拥有不可变引用的同时拥有该变量的可变引用.
2. 引用必须总是有效的.

除了引用, 另一个没有所有权的数据类型是 slice. slice 允许你引用集合中一段连续的元素序列，而不用引用整个集合. 字符串 slice（string slice）是 String 中一部分值的引用. 字符串字面值就是 slice.

> 字符串 slice range 的索引必须位于有效的 UTF-8 字符边界内，如果尝试从一个多字节字符的中间位置创建字符串 slice，则程序将会因错误而退出.

使用可变借用的前提是：出借所有权的绑定变量必须是一个可变绑定.

在所有权系统中，引用&x也可以称为x的借用（Borrowing）. 通过&操作符来完成所有权租借。所以引用并不会造成绑定变量所有权的转移.

引用在离开作用域之时，就是其归还所有权之时:
- 不可变借用（引用）不能再次出借为可变借用
- 不可变借用可以被出借多次
- 可变借用只能出借一次
- 不可变借用和可变借用不能同时存在，针对同一个绑定而言
- 借用的生命周期不能长于出借方的生命周期

核心原则：**共享不可变，可变不共享**

因为解引用操作会获得所有权，所以在需要对移动语义类型（如&String）进行解引用时需要特别注意.

> rust 有一个叫 自动引用和解引用（automatic referencing and dereferencing）的功能, 这与golang相同.

## slice
切片（Slice）是对数据值的部分引用. 切片用`&[T]`表示，其中 T 表示任意类型. 切片的引用是一个胖指针（fat pointer）.

```rust
x..y // [x, y) 的数学含义

..y // 等价于 0..y
x.. // 等价于位置 x 到数据结束
.. // 等价于位置 0 到结束
```

> str 是 Rust 核心语言类型, 就是字符串切片（String Slice）, 常常以引用的形式出现（&str）. String 类型是 Rust 标准公共库提供的一种数据类型, 它有所有权, 它的功能更完善——它支持字符串的追加、清空等实用的操作. String 和 str 除了同样拥有一个字符开始位置属性和一个字符串长度属性以外还有一个容量（capacity）属性. String 和 str 都支持切片，切片的结果是 &str 类型的数据.

> 与go不同, rust slice没有cap, 且没法直接创建, 必须依赖数组或 Vec并通过引用来创建; Go 中使用`:`来引用片段, 而 Rust 使用 `..`.

## 生命周期
Rust 生命周期机制是与所有权机制同等重要的资源管理机制, 引入这个概念主要是应对复杂类型系统中资源管理的问题.

生命周期用于处理引用, 即Rust 中的所有**引用**都附加了生命周期信息. 生命周期定义了引用相对值的原始所有者的生存周期，以及引用
作用域的范围.

> 引用是对待复杂类型时必不可少的机制，毕竟复杂类型的数据不能被处理器轻易地复制和计算. 但引用往往导致极其复杂的资源管理问题.

大多数情况下它是隐式的, 编译器通过分析代码来确定变量的生命周期. 在某些情况下, 编译器却不能确定变量的生命周期而需要开发者指明.

Rust 中的每一个引用都有其 生命周期（lifetime），也就是引用保持有效的作用域. 大部分时候生命周期是隐含并可以推断的，正如大部分时候类型也是可以推断的一样. 但当引用的生命周期可能以一些不同方式相互关联时，Rust 需要开发者使用泛型生命周期参数来注明它们的关系，这样就能确保运行时实际使用的引用绝对是有效的.

它是一类允许开发者向编译器提供引用如何相互关联的泛型. Rust 的生命周期功能允许在很多场景下借用值的同时仍然使编译器能够检查这些引用的有效性. 它是 Rust 最与众不同的功能.

**生命周期纯粹是一个编译期构造**, 它可以帮助编译器确定某个引用有效的作用域, 并确保它遵循借用规则. 它可以跟踪诸如引用的来源，以及它们是否比借用值生命周期更长
这类事情. Rust 中的生命周期能够确保引用的存续时间不超过它指向的值.

Rust 编译器有一个 借用检查器（borrow checker），它比较作用域来确保所有的借用都是有效的.

#### 生命周期注解语法
生命周期注释是描述引用生命周期的办法.

**生命周期注解并不改变任何引用的生命周期的长短**. 与当函数签名中指定了泛型类型参数后就可以接受任何类型一样，当指定了泛型生命周期后函数也能接受任何生命周期的引用. 生命周期注解描述了多个引用生命周期相互的关系，而不影响其生命周期.

生命周期注解有着一个不太常见的语法：生命周期参数名称必须以撇号（'）开头，其名称通常全是小写，类似于泛型其名称非常短. `'a` 是大多数人默认使用的名称. 生命周期参数注解位于引用的 & 之后，并有一个空格来将引用类型与生命周期注解分隔开.
```rust
&i32        // 引用
&'a i32     // 带有显式生命周期的引用
&'a mut i32 // 带有显式生命周期的可变引用
```

example:
```rust
fn longer<'a>(s1: &'a str, s2: &'a str) -> &'a str { // 让函数返回值的生命周期将与两个参数的生命周期一致
    if s2.len() > s1.len() {
        s2
    } else {
        s1
    }
}

fn main() {
    let r;
    {
        let s1 = "rust";
        let s2 = "ecmascript";
        r = longer(s1, s2);
        println!("{} is longer", r);
    }
}
```

单个的生命周期注解本身没有多少意义，因为生命周期注解告诉 Rust 多个引用的泛型生命周期参数如何相互联系的.

泛型、特性与生命周期协同:
```rust
use std::fmt::Display;

fn longest_with_an_announcement<'a, T>(x: &'a str, y: &'a str, ann: T) -> &'a str
    where T: Display
{
    println!("Announcement! {}", ann);
    if x.len() > y.len() {
        x
    } else {
        y
    }
}
```

**生命周期也是泛型**.

#### 词法作用域（生命周期）
match、for、loop、while、if let、while let、花括号、函数、闭包都会创建新的作用域，相应绑定的所有权会被转移.

函数体本身是独立的词法作用域：
- 当复制语义类型作为函数参数时，会按位复制
- 如果是移动语义作为函数参数，则会转移所有权

借用规则： 借用方的生命周期不能长于出借方的生命周期.

非词法作用域生命周期(Non-Lexical Lifetime，NLL)

#### 生命周期参数
编译器的借用检查机制无法对跨函数的借用进行检查，因为当前借用的有效性依赖于词法作用域. 所以，需要开发者显式的对借用的生命周期参数进行标注.

显式生命周期参数
- 生命周期参数必须是以单引号开头；
- 参数名通常都是小写字母，例如：'a；
- 生命周期参数位于引用符号&后面，并使用空格来分割生命周期参数和类型

标注生命周期参数是由于borrowed pointers导致的. 因为有borrowed pointers，当函数返回borrowed pointers时，为了保证内存安全，需要关注被借用的内存的生命周期(lifetime).

标注生命周期参数并不能改变任何引用的生命周期长短，它只用于编译器的借用检查，来防止悬垂指针. 即：生命周期参数的目的是帮助借用检查器验证合法的引用，消除悬垂指针.

##### 函数签名中的生命周期参数
就像泛型类型参数，泛型生命周期参数需要声明在函数名和参数列表间的尖括号中.

下文的`'a`的实际含义是 foo 函数返回的引用的生命周期与传入该函数的引用的生命周期的较小者一致(即作用域相重叠的那一部分). 这就是开发者告诉 Rust 需要其保证的约束条件. 记住通过在函数签名中指定生命周期参数时， 并没有改变任何传入值或返回值的生命周期，而是指出任何不满足这个约束条件的值都将被借用检查器拒绝.

当在函数中使用生命周期注解时，这些注解出现在函数签名中，而不存在于函数体中的任何代码中. 这是因为 Rust 能够分析函数中代码而不需要任何协助，不过当函数引用或被函数之外的代码引用时，让 Rust 自身分析出参数或返回值的生命周期几乎是不可能的. 这些生命周期在每次函数被调用时都可能不同. 这也就是为什么需要手动标记生命周期.

```rust
fn foo<'a>(s: &'a str, t: &'a str) -> &'a str;
```
函数名后的<'a>为生命周期参数的声明. 函数或方法参数的生命周期叫做输入生命周期（input lifetime），而返回值的生命周期被称为输出生命周期（output lifetime）.

规则：
- 禁止在没有任何输入参数的情况下返回引用，因为会造成悬垂指针
- 从函数中返回（输出）一个引用，其生命周期参数必须与函数的参数（输入）相匹配，否则，标注生命周期参数也毫无意义.

对于多个输入参数的情况，也可以标注不同的生命周期参数.

```rust
fn main() {
    let string1 = String::from("long string is long");
    let result;
    {
        let string2 = String::from("xyz");
        result = longest(string1.as_str(), string2.as_str());
    }
    println!("The longest string is {}", result); // 报错, 因为string2较长时, 到println!时, string2已被释放.
}

fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}
```

当从函数返回一个引用，**返回值的生命周期参数需要与一个参数的生命周期参数相匹配**. 如果返回的引用 没有 指向任何一个参数，那么唯一的可能就是它指向一个函数内部创建的值，它将会是一个悬垂引用，因为它将会在函数结束时离开作用域. 比如:
```rust
fn longest<'a>(x: &str, y: &str) -> &'a str {
    let result = String::from("really long string");
    result.as_str()
}
```

针对这种情况，最好的解决方案是返回一个有所有权的数据类型而不是一个引用， 这样函数调用者就需要负责清理这个值了.

##### 方法定义中的生命周期参数
结构体中包含引用类型成员时，需要标注生命周期参数，则在impl关键字之后也需要声明生命周期参数，并在结构体名称之后使用.

```rust
impl<'a> Foo<'a> {
    fn split_first(s: &'a str) -> &'a str {
        …
    }
}

struct Decoder<'a, 'b, S, R> {
    schema: &'a S,
    reader: &'b R
}

impl<'a, 'b, S, R> Decoder<'a, 'b, S, R>
    where 'a: 'b { // 'a:'b表示'a 的生命周期比'b 长
}
```

在添加生命周期参数'a之后，结束了输入引用的生命周期长度要长于结构体Foo实例的生命周期长度.

##### 结构体定义中的生命周期参数
结构体在含有引用类型成员的时候也需要标注每一个引用的生命周期参数，否则编译失败.

```rust
struct Foo<'a> {
    part: &'a str,
    ...
}
```

struct的生命周期参数标记，实际上是和编译器约定了一个规则：结构体实例的生命周期应短于或等于任意一个成员的生命周期.

> 注：枚举体和结构体对生命周期参数的处理方式是一样的

##### 静态生命周期参数
静态生命周期 'static：是Rust内置的一种特殊的生命周期. **'static生命周期存活于整个程序运行期间**. 所有的字符串字面量都有生命周期，类型为`& 'static str`

字符串字面量是全局静态类型，他的数据和程序代码一起存储在可执行文件的数据段中，其地址在编译期是已知的，并且是只读的，无法更改.

##### 省略生命周期参数
满足以下三条规则时，可以省略生命周期参数. 该场景下，是将其硬编码到Rust编译器中，以便编译期可以自动补齐函数签名中的生命周期参数

生命周期省略规则：
1. 每一个是引用的参数都有它自己唯一的生命周期参数
1. 如果只有一个输入生命周期参数（无论省略还是没省略），则该生命周期都将分配给输出生命周期参数
1. 如果有多个输入生命周期参数，而其中包含 &self 或者 &mut self，那么所有输出生命周期参数被赋予 self 的生命周期. 针对这条真正能够适用的就只有方法签名

举例:
```rust
fn first_word(s: &str) -> &str { // 原始
fn first_word<'a>(s: &'a str) -> &str { // 编译器应用第一条规则，也就是每个引用参数都有其自己的生命周期
fn first_word<'a>(s: &'a str) -> &'a str { // 对于第二条规则，因为这里正好只有一个输入生命周期参数所以是适用的. 第二条规则表明输入参数的生命周期将被赋予输出生命周期参数. 现在这个函数签名中的所有引用都有了生命周期，如此编译器可以继续它的分析而无须程序员标记这个函数签名中的生命周期
```

如果编译器检查完这三条规则后仍然存在没有计算出生命周期的引用，编译器将会停止并生成错误.

> 被编码进 Rust 引用分析的模式被称为 生命周期省略规则（lifetime elision rules）

##### 生命周期限定

生命周期参数可以向trait那样作为泛型的限定，有以下两种形式：

T: 'a，表示T类型中的任何引用都要“获得”和'a一样长。
T: Trait + 'a，表示T类型必须实现Trait这个trait，并且T类型中任何引用都要“活的”和'a一样长

##### 高阶生命周期
Rust还提供了高阶生命周期（Higher-Ranked Lifetime）方案，该方案也叫高阶trait限定（Higher-Ranked Trait Bound，HRTB）。该方案提供了for<>语法.

for<>语法整体表示此生命周期参数只针对其后面所跟着的“对象”.

## 并发安全与所有权
如果类型T实现了Send： 就是告诉编译器该类型的实例可以在线程间安全传递所有权.
如果类型T实现了Sync：就是向编译器表明该类型的实例在多线程并发中不可能导致内存不安全，所以可以安全的跨线程共享.

## 函数
函数是编程语言的基本要素，它是对完成某个功能的一组相关语句和表达式的封装。函数也是对代码中重复行为的抽象.

### 函数参数
- 当函数参数按值传递时，会转移所有权或者执行复制（Copy）语义
- 当函数参数按引用传递时，所有权不会发生变化，但是需要有生命周期参数（符合规则时不需要显示的标明）

### 函数参数模式匹配
- ref ：使用模式匹配来获取参数的不可变引用
- ref mut ：使用模式匹配来获取参数的可变引用
- 除了ref和ref mut，函数参数也可以使用通配符来忽略参数

对于引用, Rust 支持两种模式: ref 模式和 `&`模式. 前者借用匹配值的元素, 后者匹配引用.

### 泛型函数
函数参数并未指定具体的类型，而是用了泛型T，对T只有一个Mult trait限定，即只有实现了Mul的类型才可以作为参数，从而保证了类型安全.

泛型函数并未指定具体类型，而是靠编译器来进行自动推断的。如果使用的都是基本原生类型，编译器推断起来比较简单。如果编译器无法自动推断，就需要显式的指定函数调用的类型。

### 法和函数
方法代表某个实例对象的行为，函数只是一段简单的代码，它可以通过名字来进行调用。方法也是通过名字来进行调用，但它必须关联一个方法接受者。

### 高阶函数
高阶函数是指以函数作为参数或返回值的函数，它是函数式编程语言最基础的特性

## 闭包(closure)
闭包通常是指词法闭包，是一个持有外部环境变量的匿名函数.

> 闭包是将函数，或者说代码和其环境一起存储的一种数据结构。闭包引用的上下文中的自由变量，会被捕获到闭包的结构中，成为闭包类型的一部分.

外部环境是指闭包定义时所在的词法作用域.

外部环境变量，在函数式编程范式中也被称为自由变量，是指并不是在闭包内定义的变量.

**将自由变量和自身绑定的函数就是闭包**.

> 闭包的大小在编译期是未知的.

Rust 的 闭包（closures）是可以保存进变量或作为参数传递给其他函数的匿名函数. 可以在一个地方创建闭包，然后在不同的上下文中执行闭包运算. 不同于函数，闭包允许捕获调用者作用域中的值.

### 闭包的基本语法
闭包由管道符（两个对称的竖线）和花括号（或圆括号）组成.

- 管道符里是闭包函数的参数，可以向普通函数参数那样在冒号后添加类型标注，也可以省略

    例如：let add = |a, b| -> i32 { a + b };

- 花括号里包含的是闭包函数执行体，花括号和返回值也可以省略。

    例如：let add = |a, b| a + b;

- 当闭包函数没有参数只有捕获的自由变量时，管道符里的参数也可以省略

    例如： let add = || a + b;

**如果尝试对同一闭包使用不同类型则会得到类型错误**.

**闭包不要求像 fn 函数那样在参数和返回值上注明类型(但允许自行标注类型)**. 函数中需要类型注解是因为它们是暴露给用户的显式接口的一部分. 严格的定义这些接口对于保证所有人都认同函数使用和返回值的类型来说是很重要的,  但是闭包并不用于这样暴露在外的接口：它们储存在变量中并被使用，不用命名他们或暴露给库的用户调用.

闭包通常很短，并只关联于小范围的上下文而非任意情境. 在这些有限制的上下文中，编译器能可靠的推断参数和返回值的类型.

### 闭包的实现
闭包是一种语法糖. 闭包不属于Rust语言提供的基本语法要素，而是在基本语法功能之上又提供的一层方便开发者编程的语法.

闭包和普通函数的差别就是闭包可以捕获环境中的自由变量.

闭包可以作为函数参数，这一点直接提升了Rust语言的抽象表达能力. 当它作为函数参数传递时，可以被用作泛型的trait限定，也可以直接作为trait对象来使用.

闭包无法直接作为函数的返回值，如果要把闭包作为返回值，必须使用trait对象.

### 闭包与所有权
闭包可以通过三种方式捕获其环境，他们直接对应函数的三种获取参数的方式：获取所有权，可变借用和不可变借用. 即闭包表达式会由编译器自动翻译为结构体实例，并为其实现Fn、FnMut、FnOnce三个trait中的一个:
- `FnOnce`：会转移方法接收者的所有权。没有改变环境的能力，只能调用一次。

    FnOnce 需要**取得其参数的所有权**，只能调用一次.

    从执行环境中获取数据的所有权的闭包实现了 FnOnce 特征. 该名称表示此闭包只能被调用一次。因此，相关的变量只能使用一次。这是构造和使用闭包最不推荐的方法，因为后续不能使用其引用的变量

    ```rust
    fn main() {
        let mut a = Box::new(23);
        let call_me = || {
            let c = a;
            _ = c;
        };
        call_me(); // 成功
        // call_me(); // 报错
    }

    fn main() {
        let range = 0..10;
        let get_range_count = || range.count();
        assert_eq!(get_range_count(), 10); // ✅
        get_range_count(); // ❌
    }
    ```
- FnMut:会对方法接收者进行可变借用。有改变环境的能力，可以多次调用。

    FnMut 只需要**取得可变的引用**，可以多次调用.
    
    当编译器检测出闭包改变了执行环境中引用的某个值时，它实现了 FnMut 特征.

    FnMut 可以在任何可以使用 FnOnce 的地方使用.

    ```rust
    fn main() {
        let mut a = String::from("Hey!");
        let mut fn_mut_closure = || {
            a.push_str("Alice");
        };
        fn_mut_closure();
        println!("Main says: {}", a);
    }
    ```
- Fn:会对方法接收者进行不可变借用。没有改变环境的能力，可以多次调用。

    Fn只需要不可变的引用并可多次调用, 且不改变它从环境中捕获的任何变量, 即 Fn 闭包没有副作用或无状态.

    Fn 可以用在任何可以使用 FnMut 的地方，包括可以使用 FnOnce 的地方.

    仅为读取访问变量的闭包实现 Fn 特征。它们访问的任何值都是引用类型（ &T）。这是使用闭包的默认模式.

    ```rust
    fn main() {
        let a = String::from("Hey!");
        let fn_closure = || {
            println!("Closure says: {}", a);
        };
        fn_closure();
        println!("Main says: {}", a);
    }
    ```

如果要实现Fn，就必须实现FnMut和FnOnce
如果要实现FnMut，就必须实现FnOnce
如果要实现FnOnce，就不需要实现FnMut和Fn

由于所有闭包都可以被调用至少一次，所以所有闭包都实现了 FnOnce. 那些并没有移动被捕获变量的所有权到闭包内的闭包也实现了 FnMut, 而不需要对被捕获的变量进行可变访问的闭包则也实现了 Fn.

如果希望强制闭包获取其使用的环境值的所有权, 可以在参数列表前使用 **move 关键字. 这个技巧在将闭包传递给新线程以便将数据移动到新线程中时最为实用**.

#### 捕获环境变量的方式
- 对于复制语义类型，以不可变引用（&T）来进行捕获
- 对于移动语义类型，执行移动语义，转移所有权来进行捕获
- 对于可变绑定，并且在闭包中包含对其进行修改的操作，则以可变引用（&mut T）来进行捕获

Rust使用move关键字来强制让闭包所定义环境中的自由变量转移到闭包中.

#### 规则总结
- 如果闭包中没有捕获任何环境变量，则默认自动实现Fn
- 如果闭包中捕获了复制语义类型的环境变量，则：

    - 如果不需要修改环境变量，无论是否使用move关键字，均会自动实现Fn。
    - 如果需要修改环境变量，则自动实现FnMut。

- 如果闭包中捕获了移动语义类型的环境变量，则：

    - 如果不需要修改环境变量，而且没有使用move关键字，则会自动实现FnOnce。
    - 如果不需要修改环境变量，而且使用move关键字，则会自动实现Fn。
    - 如果需要修改环境变量，则自动实现FnMut。

- FnMut的闭包在使用move关键字时，如果捕获变量是复制语义类型的，则闭包会自动实现Copy/Clone。如果捕获变量是移动语义类型的，则闭包不会自动实现Copy/Clone

## 迭代器
迭代器模式允许对一个序列的项进行某些处理. 迭代器（iterator）负责遍历序列中的每一项和决定序列何时结束的逻辑, 是一种高效访问集合类型元素的方法.

> 在 Rust 中，**迭代器是 惰性的（lazy）**, 这意味着在调用方法使用迭代器之前它都不会有效果, 仅在需要时对集合中的元素进行求值或访问.

迭代器都实现了一个叫做 Iterator 的定义于标准库的 trait. next 是 Iterator 实现者被要求定义的唯一方法. next 一次返回迭代器中的一个项，封装在 Some 中，当迭代器结束时，它返回 None.

> iter()生成一个不可变引用的迭代器. 如果需要一个获取所有权并返回拥有所有权的迭代器，则可以调用 into_iter. 类似的, 如果希望迭代可变引用，则可以调用 iter_mut().

<table>
<thead>
<tr>
<th><code>Vec&lt;T&gt;</code> 方法</th>
<th>返回</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>.iter()</code></td>
<td><code>Iterator&lt;Item = &amp;T&gt;</code></td>
</tr>
<tr>
<td><code>.iter_mut()</code></td>
<td><code>Iterator&lt;Item = &amp;mut T&gt;</code></td>
</tr>
<tr>
<td><code>.into_iter()</code></td>
<td><code>Iterator&lt;Item = T&gt;</code></td>
</tr>
</tbody>
</table>

这些调用 next 方法的方法被称为 消费适配器（consuming adaptors），因为调用它们会消费迭代器.

Iterator trait 中定义了另一类方法，被称为 迭代器适配器（iterator adaptors），允许开发者将当前迭代器变为不同类型的迭代器, 可以链式调用多个迭代器适配器. 不过因为所有的迭代器都是惰性的，必须调用一个消费适配器方法以便获取迭代器适配器调用的结果.

```rust
let v1: Vec<i32> = vec![1, 2, 3];
// v1.iter().map(|x| x + 1); // 指定的闭包从未被调用过: 因为迭代器适配器是惰性的必须有消费
let v2: Vec<_> = v1.iter().map(|x| x + 1).collect(); // 调用 collect 方法消费新迭代器
```

Rust 的 for 循环可以用于任何实现了  IntoIterator trait 的数据结构.

在执行过程中，IntoIterator 会生成一个迭代器，for 循环不断从迭代器中取值，直到迭代器返回 None 为止。因而，for 循环实际上只是一个语法糖，编译器会将其展开使用 loop 循环对迭代器进行循环访问，直至返回 None.

Rust使用的是外部迭代器，也就是for循环. 外部迭代器：外部可以控制整个遍历过程.

Rust中使用了trait来抽象迭代器模式. Iterator trait是Rust中对迭代器模式的抽象接口.

迭代器主要包含：
- next方法：迭代其内部元素
- 关联类型Item
- size_hint方法：返回类型是一个元组，该元组表示迭代器剩余长度的边界信息

Iter类型迭代器，next方法返回的是Option<&[T]>或Option<&mut [T]>类型的值. for循环会自动调用迭代器的next方法. for循环中的循环变量(是引用)则是通过模式匹配，从next返回的Option<&[T]>或Option<&mut [T]>类型中获取&[T]或&mut [T]类型的值.

IntoIter类型的迭代器的next方法返回的是Option<T>类型，在for循环中产生的循环变量是值，而不是引用.

> 迭代器可能比循环快的原因: 它`展开(unroll)`了循环. 展开是一种移除循环控制代码的开销并替换为每个迭代中的重复代码的优化.

### IntoIterator trait
如果想要迭代某个集合容器中的元素，必须将其转换为迭代器才可以使用.

Rust提供了FromIterator和IntoIterator两个trait，他们互为反操作:
- FromIterator ：可以从迭代器转换为指定类型
- IntoIterator ：可以从指定类型转换为迭代器

Intoiter可以使用into_iter之类的方法来获取一个迭代器. into_iter的参数时self，代表该方法会转移方法接收者的所有权. 而还有其他两个迭代器不用转移所有权
- Iter ：获取不可变借用，对应&self
- IterMut ：获得可变借用，对应&mut slef

### 哪些实现了Iterator的类型？
只有实现了Iterator的类型才能作为迭代器.

实现了IntoIterator的集合容器可以通过into_iter方法来转换为迭代器.

实现了IntoIterator的集合容器有：
- `Vec<T>`
- `&'a [T]`
- `&'a mut [T] => 没有为[T]类型`实现IntoIterator

### 迭代器适配器
通过适配器模式可以将一个接口转换成所需要的另一个接口. 适配器模式能够使得接口不兼容的类型在一起工作. 适配器也叫包装器(Wrapper).

迭代器适配器，都定义在std::iter模块中：
- Map ：通过对原始迭代器中的每个元素调用指定闭包来产生一个新的迭代器。
- Chain ：通过连接两个迭代器来创建一个新的迭代器。
- Cloned ：通过拷贝原始迭代器中全部元素来创建新的迭代器。
- Cycle ：创建一个永远循环迭代的迭代器，当迭代完毕后，再返回第一个元素开始迭代。
- Enumerate ：创建一个包含计数的迭代器，它返回一个元组（i,val），其中i是usize类型，为迭代的当前索引，val是迭代器返回的值。
- Filter ：创建一个机遇谓词判断式过滤元素的迭代器。
- FlatMap ：创建一个类似Map的结构的迭代器，但是其中不会包含任何嵌套。
- FilterMap ：相当于Filter和Map两个迭代器一次使用后的效果。
- Fuse ：创建一个可以快速遍历的迭代器。在遍历迭代器时，只要返回过一次None，那么之后所有的遍历结果都为None。该迭代器适配器可以用于优化。
- Rev ：创建一个可以反向遍历的迭代器

Rust可以自定义迭代器适配器

### 消费器
迭代器不会自动发生遍历行为，需要调用next方法去消费其中的数据. 最直接消费迭代器数据的方法就是使用for循环.

Rust提供了for循环之外的用于消费迭代器内数据的方法，叫做消费器（Consumer）.

Rust标准库std::iter::Iterator中常用的消费器：
- any ：可以查找容器中是否存在满足条件的元素
- fold ：该方法接收两个参数，第一个为初始值，第二个为带有两个参数的闭包。其中闭包的第一个参数被称为累加器，它会将闭包每次迭代执行的结果进行累计，并最终作为fold方法的返回值
- collect ：专门用来将迭代器转换为指定的集合类型
- `all`
- `for_each`
- `position`

## 锁
- RwLock读写锁：是多读单写锁，也叫共享独占锁. 它允许多个线程读，单个线程写。但是在写的时候，只能有一个线程占有写锁；而在读的时候，允许任意线程获取读锁。读锁和写锁不能被同时获取。
- Mutex互斥锁：只允许单个线程读和写

## 内存管理
drop-flag：在函数调用栈中为离开作用域的变量自动插入布尔标记，标注是否调用析构函数，这样，在运行时就可以根据编译期做的标记来调用析构函数.

实现了Copy的类型，是没有析构函数的. 因为实现了Copy的类型会复制，其生命周期不受析构函数的影响.

## 异步
在异步操作里, 异步处理完成后的结果, 一般用 Promise 来保存, 它是一个对象，用来描述在未来的某个时刻才能获得的结果的值, 一般存在三个状态:
- 初始状态，Promise 还未运行
- 等待（pending）状态, Promise 已运行, 但还未结束
- 结束状态, Promise 成功解析出一个值，或者执行失败

一般而言, async 定义了一个可以并发执行的任务，而 await 则触发这个任务并发执行. 大多数语言中, async/await 是一个语法糖（syntactic sugar）, 它使用状态机将 Promise 包装起来, 让异步调用的使用感觉和同步调用非常类似, 也让代码更容易阅读.

## 属性
属性是 Rust 中写给编译器看的各种指令和建议的普适语法.

Rust 代码中的属性是指元素的注释. 属性通常是编译器内置的，不过也可以由用户通过编译器插件创建。它们指示编译器为其下显示的元素注入额外的代码或含义.

属性:
- `#[<name>]`：这适用于每个元素，通常显示在它们定义的上方
- `#![<name>]`：这适用于每个软件包. 它通常位于用户软件包根目录的最顶端部分
- 其他属性

    - `#[cfg(test)]`: 此属性添加在测试模块之上，以提示编译器有条件地编译模块，但仅在测试模式下有效

常见属性:
- `#[repr(C)]` : 要求 Rust 以兼容 C 和 C++ 的方式在内存中存储结构体

    与 C 和 C++ 不同, Rust 不保证结构体的字段或元素在内存中会以某种顺序存储, 但保证把字段的值直接存储在结构体的内存块中.

## test
Rust 的单元测试一般放在和被测代码相同的文件中，使用条件编译  #[cfg(test)] 来确保测试代码只在测试环境下编译.

集成测试一般放在 tests 目录下，和 src 平行. 和单元测试不同, 集成测试只能测试 crate 下的公开接口，编译时编译成单独的可执行文件.

Rust 中的测试函数是用来验证非测试代码是否按照期望的方式运行的. 测试函数体通常执行如下三种操作：
1. 设置任何所需的数据或状态
1. 运行需要测试的代码
1. 断言其结果是我们所期望的

Rust 社区倾向于根据测试的两个主要分类来考虑问题：单元测试（unit tests）与 集成测试（integration tests）. 单元测试倾向于更小而更集中，在隔离的环境中一次测试一个模块，或者是测试私有接口. 而集成测试对于开发者的库来说则完全是外部的. 它们与其他外部代码一样，通过相同的方式使用开发者的代码，只测试公有接口而且每个测试都有可能会测试多个模块.

在编写代码方面, 编写集成测试和单元测试没有太大的区别, 唯一的区别是目录结构和其中的项目需要公开, 开发人员已经根据软件包的设计原则公开了这些项目.

Rust 的私有性规则确实允许你测试私有函数.

Rust 二进制项目的结构明确采用 src/main.rs 调用 src/lib.rs 中的逻辑的方式: 因为通过这种结构，集成测试 就可以 通过 extern crate 测试库 crate 中的主要功能了, 根本原因是只有库 crate 才会向其他 crate 暴露了可供调用和使用的函数, 而二进制 crate 只意在单独运行.

### 单元测试
单元测试的目的是在与其他部分隔离的环境中测试每一个单元的代码，以便于快速而准确的某个单元的代码功能是否符合预期。单元测试与他们要测试的代码共同存放在位于 src 目录下相同的文件中。规范是在每个文件中创建包含测试函数的 tests 模块，并使用 cfg(test) 标注模块.

测试模块的 #[cfg(test)] 注解告诉 Rust 只在执行 cargo test 时才编译和运行测试代码，而在运行 cargo build 时不这么做.

cfg 属性代表 configuration ，它告诉 Rust 其之后的项只应该被包含进特定配置选项中, 通常用于条件编译，但不限于测试代码, 比如它可以为不同体系结
构或配置标记引用或排除某些代码. 通常配置选项是 test，即 Rust 所提供的用于编译和运行测试的配置选项. 通过使用 cfg 属性，Cargo 只会在我们主动使用 cargo test 运行测试时才编译测试代码.

### 集成测试
在 Rust 中，集成测试对于需要测试的库来说完全是外部的. 同其他使用库的代码一样使用库文件，也就是说它们只能调用一部分库中的公有 API 。集成测试的目的是测试库的多个部分能否一起正常工作. 一些单独能正确运行的代码单元集成在一起也可能会出现问题，所以集成测试的覆盖率也是很重要的. 为了创建集成测试，需要先创建一个 tests 目录, 与 src 同级. Cargo 知道如何去寻找这个目录中的集成测试文件。接着可以随意在这个目录中创建任意多的测试文件，Cargo 会将每一个文件当作单独的 crate 来编译.

> tests 目录中的子目录不会被作为单独的 crate 编译或作为一个测试结果部分出现在测试输出中, 但可将其作为模块以便在任何集成测试文件中使用.

集成测试不需要将任何代码标注为`#[cfg(test)]`. tests 文件夹在 Cargo 中是一个特殊的文件夹， Cargo 只会在运行 cargo test 时编译这个目录中的文件.

可以通过指定测试函数的名称作为 cargo test 的参数来运行特定集成测试.

也可以使用 cargo test 的 --test 后跟文件的名称来运行某个特定集成测试文件中的所有测试, 比如`cargo test --test integration_test`, 只运行了 tests 目录中指定的文件 integration_test.rs 中的测试.

### 测试函数
作为最简单例子，Rust 中的测试就是一个带有 test 属性注解的函数. 属性（attribute）是关于 Rust 代码片段的元数据.

为了将一个函数变成测试函数，需要在 fn 行之前加上`#[test]`.

assert! 宏由标准库提供，在希望确保测试中一些条件为 true 时非常有用. 因此需要向 assert! 宏提供一个求值为布尔值的参数.

测试功能的一个常用方法是将需要测试代码的值与期望值做比较，并检查是否相等. 可以通过向 assert! 宏传递一个使用 == 运算符的表达式来做到。不过这个操作实在是太常见了，以至于标准库提供了一对宏来更方便的处理这些操作 —— assert_eq! 和 assert_ne!, 这两个宏分别比较两个值是相等还是不相等.


> `assert!(a == b)`中的`==`, 实际上会转变成一个方法调用`a.eq(&b)`, eq 方法来自特征 PartialEq. PartialEq 定义了
局部排序，而 Eq 需要全局排序.

> assert_eq! 和 assert_ne! 宏在底层分别使用了 == 和 !=。当断言失败时，这些宏会使用调试格式打印出其参数，这意味着被比较的值必需实现了 PartialEq 和 Debug trait.

> debug_assert!: 类似assert!, 仅在debug模式中, 主要用于代码运行时，对应该保存的任何契约或不变性进行断言的情况, 有助于在调试模式下运行代码时捕获断言异常. 类似的还有debug_assert_eq!和 debug_assert_ne!.

还可以向 assert!、assert_eq! 和 assert_ne! 宏传递一个可选的失败信息参数， 以便在测试失败时将自定义失败信息一同打印出来.

```rust
pub fn greeting(name: &str) -> String {
    format!("Hello {}!", name)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn greeting_contains_name() {
        let result = greeting("Carol");
        assert!(
            result.contains("Carol"),
            "Greeting did not contain name, value was `{}`", result
        );
    }
}
```

属性`#[should_panic]`在函数中的代码 panic 时会通过，而在其中的代码没有 panic 时失败. 还,可给 should_panic 属性增加一个可选的 expected 参数, 测试工具会确保错误信息中包含其提供的文本.

```rust
impl Guess {
    pub fn new(value: i32) -> Guess {
        if value < 1 {
            panic!("Guess value must be greater than or equal to 1, got {}.",
                   value);
        } else if value > 100 {
            panic!("Guess value must be less than or equal to 100, got {}.",
                   value);
        }

        Guess {
            value
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    #[should_panic(expected = "Guess value must be less than or equal to 100")]
    fn greater_than_100() {
        Guess::new(200);
    }
}
```

可以使用 Result<T, E> 编写测试. 但不能对这些使用 Result<T, E> 的测试使用 #[should_panic] 注解, 相反应该在测试失败时直接返回 Err 值.

```rust
#[cfg(test)]
mod tests {
    #[test]
    fn it_works() -> Result<(), String> {
        if 2 + 2 == 4 {
            Ok(())
        } else {
            Err(String::from("two plus two does not equal four"))
        }
    }
}
```

`#[ignore]`: 使用#[ignore]属性标记告知测试工具在执行 cargo test 命令时忽略此类测试功能, 然后可以向测试工具或 cargo test 命令传递`--ignored` 参数来单独运行这些测试.

### 执行测试
`rustc --test first_unit_test.rs`即测试first_unit_test.rs, 默认测试都是并行的, 除非设置了`RUST_TEST_THREADS=1`.

cargo test 在测试模式下编译代码并运行生成的测试二进制文件， 但可以指定命令行参数来改变 cargo test 的默认行为. cargo test 生成的二进制文件的默认行为是并行的运行所有测试，并截获测试运行过程中产生的输出, 这与golang类同.

当运行多个测试时， Rust 默认使用线程来并行运行. `cargo test -- --test-threads=1`可传递 --test-threads 参数和希望使用线程的数量给测试二进制文件, 这里将测试线程设置为 1，告诉程序不要使用任何并行机制.

如果希望能看到通过的测试中打印的值，截获输出的行为可以通过 --nocapture 参数来禁用, 比如` cargo test test_with_fixture --
--nocapture`, test_with_fixture是测试函数.

可以向 cargo test 传递任意测试的名称来只运行这个测试, 比如`cargo test one_hundred`.

指定部分测试的名称，任何名称匹配这个名称的测试会被运行, 比如`cargo test add`, 只运行了所有名字中带有 add 的测试.

可以使用 ignore 属性来标记耗时的测试并排除他们. 当需要运行 ignored 的测试时，可以执行`cargo test -- --ignored`.

### 基准测试
`cargo bench`

基准测试基于:
1. 函数上方的#[bench]注释，这表示该函数是一个基准测试
1. 内部编译器软件包 libtest 包含一个 Bencher 类型，基准函数通过它在多次迭代中运行相同的基准代码，此类型是针对编译器内部的，只适用于测试模式

```rust
// bench_example/src/lib.rs
#![feature(test)]
extern crate test;

use test::Bencher;

pub fn do_nothing_slowly() {
    print!(".");
    for _ in 1..10_000_000 {};
}

pub fn do_nothing_fast() {
}

#[bench]
fn bench_nothing_slowly(b: &mut Bencher) {
    b.iter(|| do_nothing_slowly());
}

#[bench]
fn bench_nothing_fast(b: &mut Bencher) {
    b.iter(|| do_nothing_fast());
}
```

在标有#[bench]注释的函数内部， iter 的参数是一个没有参数的闭包函数.

输出格式是每次迭代花费的时间，括号内的数字表示每次运行之间的差异. 性能较差的实现的运行速度非常慢, 并且运行时间不固定（用+/−符号所示）.

criterion-rs可生成比内置基准测试框架更多的统计报告, 并使用 gnuplot 生成实用的图形和报表, 使用户更容易理解, 使用方法是:
```
# cat Cargo.toml
...
[dev-dependencies]
criterion = "0.1"

[[bench]]
name = "fibonacci"
harness = false
...

# cat criterion_demo/benches/fibonacci.rs
#[macro_use] // 意味着要使用来自此软件包的任何宏时，我们需要使用此属性来选择它，因为默认情况下它们是非公开的
extern crate criterion;
extern crate criterion_demo;

use criterion_demo::{fast_fibonacci, slow_fibonacci};
use criterion::Criterion;

fn fibonacci_benchmark(c: &mut Criterion) { // c可保存基准代码的运行状态
    c.bench_function("fibonacci 8", |b| b.iter(|| slow_fibonacci(8)));
}

criterion_group!(fib_bench, fibonacci_benchmark); // 将fib_bench 的基准名称分配给基准组
criterion_main!(fib_bench);
```
添加了一个名为`[[bench]]`的新属性， 它向 Cargo 表明我们有一个名为 fibonacci的新基准测试，并且它不使用内置的基准测试工具（ harness=false）, 因为我们正在使用
criterion 软件包的测试工具.

使用

## FAQ
### 并发与并行
并发是一种能力，而并行是一种手段. 因为并行处理是基于硬件的，而并发处理是可以通过设计编码进行提高的.

在代码的运行方式中，并发是并行的基础，是同时与多个任务打交道的能力；并行是并发的体现，是同时处理多个任务的手段.

很多拥有高并发处理能力的编程语言，会在用户程序中嵌入一个 M:N 的调度器，把 M 个并发任务，合理地分配在 N 个 CPU core 上并行运行，让程序的吞吐量达到最大, 比如go.

### 同步和异步
同步是指一个任务开始执行后，后续的操作会阻塞，直到这个任务结束.

同步执行保证了代码的因果关系（causality），是程序正确性的保证.

然而在遭遇 I/O 处理时，高效 CPU 指令和低效 I/O 之间的巨大鸿沟，成为了软件的性能杀手.

异步是指一个任务开始执行后，与它没有因果关系的其它任务可以正常执行，不必等待前一个任务结束.

在异步操作里，异步处理完成后的结果，一般用 Promise 来保存，它是一个对象，用来描述在未来的某个时刻才能获得的结果的值，一般存在三个状态:
1. 初始状态，Promise 还未运行
1. 等待（pending）状态，Promise 已运行，但还未结束
1. 结束状态， Promise 成功解析出一个值，或者执行失败

在很多支持异步的语言中，Promise 也叫 Future / Delay / Deferred 以及 async/await.

一般而言，async 定义了一个可以并发执行的任务，而 await 则触发这个任务并发执行。大多数语言中，async/await 是一个语法糖（syntactic sugar），它使用**状态机**将 Promise 包装起来，让异步调用的使用感觉和同步调用非常类似，也让代码更容易阅读.

### 各编程语言中的类型系统
静态类型语言：在编译阶段确定所有变量的类型.
动态类型语言：在执行阶段确定所有变量的类型.
不允许隐式转换的是强类型，允许隐式转换的是弱类型.

- 没类型： 比如汇编语言，没有类型的概念，所有都只是一个数字
- 弱静态类型： 比如C/C++语言，可以定义类型，但是不强制执行，在不同类型之间自动转换
- 强静态类型： 比如Java，定义类型，并且用虚拟机检查类型
- 强动态类型： 比如Python和Ruby，动态推断类型而不需要定义，然后解释器会强制执行.
- 弱动态类型： Perl/PHP,js

> ruby和Python在运行时通过Duck Typing来进行运行时类型检查, 以保证类型安全.

![](/misc/img/rust/9zv7ejia98.jpeg)

### 函数式编程（functional programming）
函数式编程风格通常包含将函数作为参数值或其他函数的返回值、将函数赋值给变量以供之后执行等等.

Rust 的一些在功能上与其他被认为是函数式语言类似的特性:
- 闭包（Closures），一个可以储存在变量里的类似函数的结构
- 迭代器（Iterators），一种处理元素序列的方式