# rust
rust是一门系统编程语言, 特点:
1. 高性能: 没有gc
1. 高可靠: 所有权模型保证了内存安全和线程安全
1. 生产力: 拥有出色的文档、友好的编译器和清晰的错误提示信息, 还集成了一流的工具 cargo, 智能地自动补全和类型检验的多编辑器支持, 以及自动格式化代码等等

Rust 中最大的思维转换就是变量的所有权和生命周期，这是几乎所有编程语言都未曾涉及的领域.

> 系统编程语言是相对于应用级编程语言, 它处于更底层, 更接近硬件层次, 其特点有:
> 1. 可在资源非常受限环境运行
> 1. 运行时开销小, 非常高效
> 1. 很小的运行库, 甚至没有
> 1. 可以允许直接的内存操作
>
> 其他系统编程语言有c,c++等.


Rust 编译器的版本号采用了`语义化版本号（ Semantic Versioning ）`规划, 版本格式为: `主版本号．次版本号.修订号`.

为了兼顾更新速度以及稳定性, Rust 使用了多渠道发布的策略:
- nightly

	nightly 版本是每天在主版本上自动创建出来的版本，这个版本上的功能最多，更新最快，但是某些功能存在问题的可能性也更大 因为新功能会首先在这个版本上开启，供用户试用.

	nightly 版本中使用试验性质的功能，必须手动开启 featur gate 也就是说要在当前项目的入口文件中加入`#![feature (name,...)]`语句, 否则是编译不过的. 等到这个功能最终被稳定了, 再用新版编译器编译的时候，它会警告你这个 feature gate 现在是多余的了，可以去掉了.

- beta

	beta 版本是每隔 段时间，将一些在 nightly 版本中验证过的功能开放给用户使可以被看作 stable 版本的“预发布”版本
- stable

	stable 版本则是正式版，它定期发布一个新版本，一些实验性质的新功能在此版本上无法使用, 因此它也是最稳定、最可靠的, 保证向前兼容的版本. 


Rust 相对重大的设计必须经过 RFC(Request For Comments ）设计步骤. 这个步骤主要是用于讨论如何“设计”语言. [这个项目](https://github.com/rust-lang/rfcs)旨在于所有大功能必须先写好设计文挡，讲清楚设计的目标、实现方式 优缺点等，让整个社区参与讨论，然后由“核心组”(Core Team)的成员参与定夺是否接受这个设计. 许多深层次的设计思想问题可以在这个项目中找到答案.

rust的`RFC -> Nightly -> Beta -> Stable`策略成功实践了快速迭代、敏捷交付以及 视用户反馈的特点，同时也保证了核心设计的稳定性--用户可以根据自己的要和风险偏好，选择合适的版本.

2017 年下半年, rust 设计组又提出了一个基于 poch 演进策略（后来也被称为edition). 它要解决的问题是`如何让 Rust 平稳地进化`. 简单来说就是让 Rust 的兼容性保证是一个有时限的长度, 而不是永久.

## 内存安全
ref:
- [Rust、Go、C ，哪个才是“内存管理大师”？](https://mp.weixin.qq.com/s?__biz=MjM5ODI5Njc2MA==&mid=2655873042&idx=1&sn=db17488d9ea741c280cf26ee6fb4b8ea)

Rust通过使用借用检查器(borrow checker)、所有权(ownership)、借用(borrow)这三个概念来管理和确保跨堆栈和堆的内存安全来管理内存，从而实现内存管理.

# base
## 代码
rust代码使用`.rs`扩展名, 且必须是utf-8编码.

注释支持:
1. `//` : 行注释
1. `/**/`: 块注释
1. 文档注释

## 语句/表达式
Rust 程序里, 表达式（ Expressio ）和语句（ Statement ）是完成流程控制、计算求值的主要工具. 在Rust 里, 表达式可以是语句的一部分，反过来，语句也可以是表达式的一部分. 一个表达式总是会产生 个值，因此它必然有类型; 语句不产生值，它的类型永远是`()`. 如果把一个表达式加上分号，那么它就变成了一个语句；如果把语句放到一个语句块中包起来, 那么它就可以被当成一个表达式使用.

rust中语句块也可以是表达式的一部分. 语句和表达式的区分是后者不带`;`. 如果带了分号, 意味着这是一条语句，它的类型是`()`; 如果不带分号，它的类型就是表达式的类型.

Rust 中的表达式语法具有非常好的“一致性”，每种表达式都可以嵌入到另一种表达式中，组成更强大的表达式.

Rust 表达式包括字面量表达式、方法调用表达式、数组表达式、索引表达式、单目运算符表达式、双目运算符表达式等:
- 运算表达式

	比较运算符的两边必须是同类型的, 并满足 PartialEq 约束.

	**Rust里面的运算符优先级与C语言里面的运算符优先级设置是不一样的，有些细微的差别. 建议: 如果碰到复杂一点的表达式, 尽量用小括号明确表达计算顺序, 避免依赖语言默认的运算符优先级**.
- 赋值表达式

	Rust规定, 赋值表达式的类型为`unit`即`()`, 以避免像C语言那样的连续赋值, 或者误将`==`写成`=`

	Rust 不支持`++`,`--`运算符，请使用`+=1`,`-=1`替代

	在 rust 中没有必要专门设计像 C/C＋＋ 那样的三元运算符`?:`, 因为通过现有的设计可以轻松实现同样的功能, 且可读性更佳.

Rust 表达式又可以分为‘左值’ （lvalue ）和‘右值’（rvalue)两类. 所谓左值是这个表达式可以表达一个内存地址，它们可以放到赋值运算符左边使用, 其他的都是右值.

## 变量, 函数和数据结构
- 变量

	- 不可变: `let x:T;`
	- 可变: `let mut x:T;`

		```rust
		let (mut a, mut b) = (1, 2); 
		let Point { x : ref a, y : ref b} = p;
		let mut v:Vec<u8> = Vec::new();
		```

	实际上, let 语句引入了一个模式解构, 不能把 let mut 视为一个组合, 而应该将 `mut x` 视为一个组合.

	Rust中，一般把声明的局部变量并初始化的语句称为“变量绑定”，强调的是“绑定”的含义，与 C/C＋＋ 中的“赋值初始化”语句有所区别.

	Rust 中，每个变量必须被合理初始化之后才能被使用, 使用未初始化变量这样的错误，在Rust 是不可能出现的 （利用 unsafe hack 除外）.

	Rust 允许在同一个代码块中声明同样名字的变量, 如果这样做, 后面声明的变量会将前面声明的变量“遮蔽”(Shadowing)起来. 实际上，传统编程语 C/C＋＋ 中也存在类似的功能，只不过它们只允许嵌套的代码块内部变量出现遮蔽, 而 Rust 在这方面放得稍微宽一点, 同一语句块内部声明的变量也可以发生遮蔽.

	rust变量声明的类型后置更方便类型推导.

	Rust 只允许**`局部变量/全局变量`实现类型推导，而函数签名等场景则是不允许, 这是特意这样设计的**:
	局部变量只有局部的影响; 全局变量必须当场初始化; 而函数名具有全局性影响. 函数签名如果使用自动挡类型推导, 可能导致某个调用的地方使用方式发生变化, 它的参数、返回值类型就发生了变化, 进而导致其他调用地方的编译错误，这是设计者不希望看到的情况.

- 静态变量

	- 不可变: `static X:T = T::new();`
	- 可变: `static mut X:T = T::new();`

	静态变量和常量一样全局可访问, 它也被编入可执行文件的数据段. 静态变量可以被声明为可变.

	static 语句同样也是一个模式匹配, 与let 不同的是，用static 声明的变量的生命周期是整个程序，从启动到退出, static变量的生命周期永远是`'static`, 它占用的内存也不会在执行过程中被回收. 这也是 rust 中唯一的声明全局的方法.

	在使用静态变量时, 有一些限制, 可用lazy_static检查:
	1. 全局变量必须在声明的时候马上初始化

		**局部变量声明后只要在使用前初始化即可**.
	1. 全局变量的初始化必须是编译期可确定的常量，不能包括执行期才能确定的表达式、语句和函数调用
	1. 带有 mut 修饰的全局变量，在使用的时候必须使用 unsafe 关键字

	rust禁止在声明static变量的时候调用普通函数, 或者利用语句块调用其他非const代码. const fn是允许的, 因为它是在编译期执行的.

- 常量: 不允许用mut修饰.

	`const X:T = <value>;`

	常量是一个右值, 它不能被修改. 常量编译后会被放在可执行文件的数据段, 全局可访问.

	常量的初始化表达式也一定是一个编译期常量, 不能是运行期的.

	它与 static 大区别在于: 编译器并不一定会给const常量分配内存空间, 在编译过程中，它很可能会被内联优化. 因此, 千万不要用 hack 的方式, 通过 unsafe 代码去修改常量的值, 这么做是没有意义的. const也不具备类似 let 模式匹配功能.

- 函数: `func x(a1:T1,...)-> T{}`

	rust函数使用fn标识, 其参数列表与let一样也可模式解构.
	在rust中， 如果函数没有返回值， 那么其返回值是unit即空元组`()`

	Rust 中, 每个函数具有自己单独的类型，但是这个类型可以自动转换成fn类型. 因此两个有同样的参数类型和同样的返回值类型的函数, 但它们是不同类型因而不能赋值给相同变量, 解决方案是让先转为通用的fn类型即可.

	Rust 支持一种特殊的发散函数（ Diverging functions ）, 它的返回类型是`!`. 发散类型的最大特点就是，它可以被转换为任意一个类型.

	Rust 中，有以下这些情况永远不会返回，它们的类型就是!:
	1. panic 以及基于它实现的各种函数/宏，比如`unimplemented!, unreachable!` ; 
	1. 死循环 `loop {}`
	1. 进程退出函数 `std::process::exit` 以及类似的 libc 中的 exec 一类函数

	在大部分主流操作系统上，一个进程开始执行的时候可以接受一系列的参数，退出的时候也可以返回一个错误码. 许多编程语言也因此为 main 函数设计了参数和返回值类型, rust和go不同, 传递参数和返回状态码都由单独的 API 来完成.

	Rust 设计组扩展了 main 函数的签名，使它变成了一个泛型函数，这个函数的返回类型可以是任何一个满足 Terminationtrait的类型，其中`（）,booL Result` 是满足这个约束的，它们都可以作为 main 函数的返回
类型.

	函数可以用 co st 键字修饰，这样的函数可以在编译阶段被编译器执行，返回值也被视为编译期常量.

- 元组: 它通过圆括号包含一组表达式构成

	如果元组中只包含一个元素, 应该在后面添加一个逗号, 以区分括号表达式和元组. 元组内部也可以一个元素都没有 这个类型单独有一个名字， unit （单元类型）.
	访问元组内部元素有两种方法, 一种是“模式匹配”（ pattern destructuring ）, 另外一种是“数字索引”.

	unit 类型是 Rust 最简单的类型之一， 是占用空间最小的类型之一. 空元组与空结构体 struct Foo 一样，都是占用0字节空间. `std::mem::size_of`函数可以计算一个类型所占用的内存空间.

- 结构体: `struct S{...}`, 不能使用自动类型推导功能, 必须显式指定.

	三种形式:
	1. 空结构体, 不占用任何空间, 比如`struct Marker;`
	1. 元组结构体(`tuple struct `), struct的每个成员都是匿名的, 可通过索引范围, 比如`struct Color(u8,u8,u8);`
	1. 普通结构体, struct的每个成员都有名字, 可通过名字访问

		```rust
		struct Person{
			name: String,
			age: u8,
		}
		```

	Rust 允许 struct 类型的初始化使用一种简化的写法: 如果有局部变量名字和成员变量名字恰好一致, 那么可以省略掉重复的冒号初始化.

	Rust 设计了一个语法糖，允许用一种简化的语法赋值使用另外一个 struct 的部分成员: `..expr`这样的语法, 只能放在初始化表达式中, 所有成员的最后最多只能有一个.

	`struct Fool;`/`struct Foo();`/`struct Foo{}`其实是同一个东西.

	tuple struct 一个特别有用的场景，那就是当它只包含一个元素的时候，就是所谓的newtype idiom. 因为它实际上让我们非常方便地在一个类型的基础上创建了一个新的类型.

	通过关键字 type, 可创建一个新的类型名称，但是这个类型不是全新的类型，而只是一个具体类型的别名, 在编译器看来, 这个别名与原先的具体类型是一模一样. 而使用 tuple struct 做包装，则是 创造了一个全新的类型，它跟被包装的类型不能发生隐式类型转换, 可以具有不同的方法, 满足不同的 trait, 完全按需而定.

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
		// 由于它实在是太常用, 标准库将 Option 以及它的成员 Some,None 都加入到了Prelude, 它表示的含义是`要么存在、要么不存在`
		// Rust 在语言层面彻底不允许空值 null 的存在，但无奈null 可以高效地解决少量的问题，所以 Rust 引入了 Option 枚举类. Option 是 Rust 标准库中的枚举类，这个类用于填补 Rust 不支持 null 引用的空白.
		// 由于 Option 是 Rust 编译器默认引入的，在使用时可以省略 `Option::` 直接写 None 或者 Some().
		enum Option<T>{
			Some(T),
			None,
		}
		```

		> 很多语言默认不允许 null，但在语言层面支持 null 的出现（常在类型前面用`?`符号修饰）.

	如果说`tuple, struct, tuple struct`, 在Rust 中代表的是多个类型的“与”关系，那么 enum 类型在 Rust 中代表的就是多个类型的“或”关系.

	与C/C＋＋ 中的枚举相比, Rust 中的 enum 要强大得多，它可以为每个成员指定附属的类型信息. 它是一种更安全的类型, 可以被称为`tagged union`.

	Rust enum 类型的变量需要区分它里面的数据究竟是哪种变体，所以它包含了一个内部的`tag 标记`来描述当前变量属于哪种类型, 这个标记对用户是不可见的，通过恰当的语法设计，保证标记与类型始终是匹配的，以防止用户错误地使用内部数据.

	Rust 里面也支持 union 类型，这个类型与C语言中的 union 完全一致, 但在 Rust 里面，读取它内部的值被认为是 unsafe 行为, 一般情况下我们不使用这种类型. 它存在的主要目的是为了方便与C语言进行交互.

	```rust
	enum Number {
		Int(i32),
		Float(f32),
	}
	```
	等价于:
	```c
	struct Number {
		enum {Int, Float} tag;
		union {
			int32_t int_value;
			float float_value;
		} value;
	};
	```

	在实际中, enum 的内存布局未必是这个样子， 编译器有许多优化，可以保证语义正确的同时减少内存使用，并加快执行速度. 如果是在 FFI 下， 要保证 Rust 里面的 enum 的内存布局和C语言兼容的话, 可以给这个 enum 添加一个`#[repr(C, Int)]`属性标签.

	> rust的enum其实是一种代数类型系统(Algebraic Data Type, ADT), 即enum内部的variant(用于区分enum里面的数据类型)的类型是函数类型. 因此Some可以当成函数作为参数传递给迭代器的`map()`.

Rust 里的合法标识符（包括变量名、函数名、 trait 名等） 必须由数字、字母、 下划线组成， 且不能以数字开头, 这个规定和许多现有的编程语言是一样.

rust的`_`和go类似, 是一个特殊的标识符，在编译器内部它是被特殊处理的, 其的含义是`忽略这个变量绑定, 后面不会再用到了`.

Rust 中, enum和struct均为内部成员创建了新的名字空间, 如果要访问内部成员，可
以使用`::`, 不同的 enum 中重名的元素也不会互相冲突.

rust复合类型支持递归定义, 但需要使用指针, 否则计算其大小时因递归而无解.

## 基本数据类型
- bool : true, false
- char : 单个字符, 大小为四个字节(four bytes)，并代表了一个 Unicode 标量值（Unicode Scalar Value）, 等价于go的rune. char 由单引号包裹, 不同于字符串使用双引号.

	对于ASCII字符用`u8`表示.
	
	> Unicode 标量值包含从 U+0000 ~ U+D7FF 和 U+E000 ~ U+10FFFF 在内的值.
- 整数

	整数是一个没有小数部分的数字. 有符号整数范围: -2^(n-1)~2^(n-1)-1; 无符号范围: 0~2^n-1

	> 如果一个变量是有符号类型，那么它的最高位的那一个 bit 就是“符号位”，表示该数为正值还是负值; 如果个变量是无符号类型，那么它的最高位和其他位一样，表示该数的大小.

    大小  有符号     无符号
    8 bit   i8  u8
    16 bit  i16     u16
    32 bit  i32     u32
    64 bit  i64     u64
    128 bit     i128    u128
    Arch    isize   usize // arch 是由 CPU 构架决定的大小的整型类型, 与指针占用的空间大小一致, 在 x86 机器上为 32 位，在 x64 机器上为 64 位. 即isize和usize是自适应类型, 它们主要作为某些集合的索引.

    > 所有数值字面量支持任意位置添加`_`以方便阅读, 并且支持后缀表示类型, 比如`0x_ff_u8`

    > 整数自动推导时**默认是i32**

    > 字面量后面可以跟后缀，可代表该数字的具体类型，从而省略掉显示类型标记, **个人不推荐**.

    > rust不支持`++/--`

    c中无符号算术运算永远不会overflow, 如果超出范围则自动舍弃高位数据; 有符号如果发生了overflow, 则是undefined behavior, 由编译器处理. 

    未定义行为有利于编译器做一些更激进的性能优化，但是这样的规定有可能导致在程序员不知情的某些极端场景下, 产生诡异 bug.

    Rust 的设计思路则倾向于预防 bug, 而不是无条件压榨效率, Rust 设计者希望能尽量减少`未定义行为`. 因此rust在debug模式下编译器自动插入溢出检查, 一旦溢出就panic; 在release下, 不检查整数溢出, 而是自动舍弃高位即二进制补码包装（two’s complement wrapping）的操作. rustc可使用`-C overflow-checks=no/yes`决定是否开启溢出检查.

    开发者可以调用标准库中的`checked_*, saturating_*, wrapping_*`系列函数更精细地自主控制整数溢出的行为:
	- `checked_*`系列函数返回的类型是`Option<_>`, 当出现溢出的时候，返回None
	- `saturating_`系列函数返回类型是整数, 如果溢出，则给出该类型可表示范围的`最大/最小`值
	- `wrapping_*`系列函数则是直接抛弃已经溢出的最高位, 将剩下的部分返回.

	在对安全性要求非常高的情况下, 强烈建议用户尽量使用这几个方法替代默认的算术运算符来做数学运算. rust 标准库中就大量使用了这几个方法. ，标准库还提供了一个叫作`std::num::Wrapping<T>`的类型, 它重载了基本的运算符， 可以被当成普通整数使用, 凡是被它包起来的整数, 任何时候出现溢出都是截断行为.

	```rust
	fn main() { 
	    let i = 10 0 i8; 
	    println! ("checked { : ?}", i.checked_add(i)); 
	    println! ("saturating {:?}", i.saturating_add(i)); 
	    println! ("wrapping {:?}", i. wrapping_add( i));
	}

	输出结果为:
	checked None
	saturating 127
	wrapping - 56
	```
- 浮点型

	rust提供基于IEEE-754-2008标准的浮点类型: f32/f64. `std::num::FpCategory`可表示浮点状态(`Nan/Infinite/Zero/Subnormal/Normal`), 默认是f64.

	> 非0数除以0是inf; 0除以0是Nan; `inf * 0.0`=Nan; `inf/inf`=NaN.
- 指针

	无GC 的编程语言， C/C++以及 Rust, 对数据的组织操作有更多的自由度, 表现为:
	- 同一个类型，某些时候可以指定它在栈上, 某些时候可以指定它在堆上. 内存分配方式可以取决于使用方式, 与类型本身无关.
	- 既可以直接访问数据 ，也可以通过指针间接访问数据. 可以针对任何一个对象取得指向它的指针
	- 既可以在复合数据类型中直接嵌入别的类型的实体, 也可以使用指针, 间接指向别的类型
	- 甚至可能在复合数据类型末尾嵌入不定长数据构造出不定长的复合数据类型

	Rust 有不止一种指针类型, 常见的几种指针类型:
	- `Box<T>` : 指向类型T的, 具有所有权的指针, 有权释放内存

    	Rust中的值默认被分配到栈内存, 可通过`Box<T>`将值装箱(在堆内存中分配). 可通过解引用来获取`Box<T>`中的T. 因为`Box<T>`的行为像引用, 并且可以自动释放内存, 因此将其称为智能指针.

    	String类型和Vec类型的值都是被分配到堆内存并返回指针的，通过将返回的指针封装来实现Deref和Drop.

    	Box<T>是指向类型为T的堆内存分配值的智能指针. 当Box<T>超出作用域范围时，将调用其析构函数，销毁内部对象，并自动释放堆中的内存.
	- `&T` : 指向类型T的借用指针, 也称为引用, 无权释放内存, 无权写数据
	- `&mnut T` : 指向类型T的mut型借用指针, 无权释放内存, 有权写数据
	- `*const T` : 指向类型T的只读裸指针, 没有生命周期信息, 无权写数据
	- `*mut T` : 指向类型T的可读写裸指针, 没有生命周期信息, 有权写数据

	此之外，在标准库中还有一种封装起来的可以当作指针使用的类型, 即智能指针(smart pointer, 来自c++):
	- `Rc<T>` : 指向类型T的引用计数指针, 共享所有权, 线程不安全

	    通过clone方法共享的引用所有权称为强引用，RC<T>是强引用.
	- `Arc<T>` : 指向类型T的原子型引用计数指针, 共享所有权, 线程安全
	- `Cow<’a, T>` : Clone-on-write, 写时复制指针. 可能是借用指针, 也可能是具有所有权的指针 

    	Cow<T>的功能是：以不可变的方式访问借用内容，以及在需要可变借用或所有权的时候再克隆一份数据. Cow<T>旨在减少复制操作，提高性能，一般用于读多写少的场景. Cow<T>的另一个用处是统一实现规范.

Rust使用as用于类型转换, 前提是编译器认为是合理的转换.

## 流程控制
Rust 的循环和大部分语言都一致, 支持死循环`loop {}`、条件循环`while expr {}`，以及对迭代器的循环`for x in iter {}`. 循环可以通过 break 提前终止，或者 continue 来跳到下一轮循环.

> 可用比如`let number = if a > 0 { 1 } else { -1 };`的 if-else 结构实现类似于三元条件运算表达式 (A ? B : C) 的效果

> 在 C 语言中 for 循环使用三元语句控制循环，但是 Rust 中没有这种用法，需要用 while 循环来代替. 且没有 do-while 的用法(do 被规定为保留字)

> loop 循环可以通过 break 关键字类似于 return 一样使整个循环退出并给予外部一个返回值.

> 可在loop/while/for 循环前面加上`生命周期标识符`(该标识符以单引号开头), 其在内部的循环中可以使用 break/continue 语句来选择跳出到哪一层.

> 如果一个 loop 永远不返回，那么它的类型就是“发散类型”. 编译器可以判断出发散类型, 其后代码是永远不会执行的死代码.

> `loop{}`和`while true{}` 循环有何区别, 为什么 Rust 设计了loop, 难道不是完全多余的吗？实际上不是, 主要原因在于, 相比于其他的许多语言, Rust 要做更多的静态分析它俩在运行时是没有什么区别, 它们主要是会影响编译器内部的静态分析. `let x; loop { x = 1; break; }`可以执行, `let x; while true { x = 1; break ; }`会报错, 因为编译器认为while语句的执行跟条件表达式在运行阶段的值有关, 因此不确定x是否一定会初始化而报错.

> for 循环的主要用处是利用迭代器对包含同样类型的多个元素的容器执行遍历，如数组,链表,HashMap,HashSet等.

> rust不支持switch, 而是使用match. 很多语言摒弃 switch 的原因都是因为 switch 容易存在因忘记添加 break 而产生的串接运行问题，Java 和 C# 这类语言通过安全检查杜绝这种情况出现.

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

## 宏
rust宏和c/c++中的宏完全不是一个概念. 它是一种安全版的编译期语法扩展, 之所以使用宏, 而不是函数, 是因为宏可以完成编译期格式检查, 更加安全.

> 函数则不具备字符串格式化的静态检查功能，如果出现了不匹配的情况, 只能是运行期报错.

> `format!, write!`最终还是调用`std::io`模块提供的一些函数来完成的. 如果用户需要更精细地控制标准输出操作, 也可以直接调用标准库来完成.

## 代码管理
rust支持使用mod 来组织代码.

使用方法: 在项目的入口文件`lib.rs/main.rs`里, 用 mod 来声明要加载的其它代码文件. 如果模块内容比较多, 可以放在一个目录下, 再在该目录下放一个 mod.rs 引入该模块的其它文件, mod.rs 和 Python 的 `__init__.py` 有异曲同工之妙.

在 Rust 里, 一个项目也被称为一个 crate. crate 可以是可执行项目，也可以是一个库.

在一个 crate 下，除了项目的源代码，单元测试和集成测试的代码也会放在 crate 里. Rust 的单元测试一般放在和被测代码相同的文件中，使用条件编译  #[cfg(test)] 来确保测试代码只在测试环境下编译.

集成测试一般放在 tests 目录下，和 src 平行. 和单元测试不同，**集成测试只能测试 crate 下的公开接口，编译时编译成单独的可执行文件**. 在 crate 下，如果要运行测试用例，可以使用`cargo test`.

当代码规模继续增长，把所有代码放在一个 crate 里就不是一个好主意了，因为任何代码的修改都会导致这个 crate 重新编译，这样效率不高. 此时可以使用 workspace, 一个 workspace 可以包含一到多个 crates，当代码发生改变时，只有涉及的 crates 才需要重新编译. 当要构建一个 workspace  时，需要先在某个目录下生成一个 Cargo.toml，包含 workspace 里所有的 crates，然后可以  cargo new 生成对应的 crates.

## trait
所有的 trait 中都有一个隐藏的类型 Self （大写），代表当前这个实现了此 trait 的具体类型. trait 中定义的函数，也可以称作关联函数（ associated function). 函数的第一个参数如果是 Self 相关的类型，且命名为 self（小写），这个参数可以被称为“receiver ”（接收者）. 具有 receiver 参数的函数，称为“方法”（method), 可以通过变量实例使用小数点来调用. 没有 receiver 参数的函数，称为“静态函数”（static function ），可以通过类型加`::`的方式来调用.

```rust
trait T { 
	fn methodl(self: Self); 
	fn method2(self: &Self); 
	fn method3 (self: &mut Self); 
}
// 上下两种写法是完全一样的
trait T { 
	fn methodl (self) ; 
	fn method2(&self); 
	fn method3(&mut self); 
}
```

直接对它 impl 来增加成员方法, 无须 trait 名字, 比如：
```rust
impl Circle { 

	fn get radius(&self) -> f64 { self.radius } 
}
```
可以把这段代码看作是为 Circle 类型 impl 了一个匿名的 trait. 用这种方式定义的方法叫作这个类型的`内在方法`（ inherent methods).

> 结构体方法的第一个参数必须是 &self，不需声明类型，因为 self 不是一种风格而是关键字. 而在调用结构体方法的时候不需要填写 self, 这是出于对使用方便性的考虑.

> 结构体关联函数: 在 impl 块中却没有 &self 参数, 这种函数不依赖实例，但是使用它需要声明是在哪个 impl 块中的, 比如`String::from`.

trait 中可以包含方法的默认实现, 如果需要针对特殊类型作特殊处理，也可以选择重新实现来`override`默认的实现方式.

impl 的对象甚至可以是 trait, 如下:
```rust
trait Shape { 
	fn area(&self) - > f64;
}

trait Round { 
	fn get_radius(&self) - > f64;
}

struct Circle { 
	radius: f64, 
}

impl Round for Circle {
	fn get_radius(&self) -> f64 { self.radius } 
}

// impl Trait for Trait 
impl Shape for Round { 
	fn area(&self) -> f64 { 
		std::f64::consts::PI * self.get_radius() * self.get_radius() 
	}
}

fn main() { 
	let  c =Circle { radius : 2f64}; 
	// build err
	// c. area ( ) ; 

	let b = Box::new(Circle {radius : 4f64}) as Box<Round>;
	b.area();
}
```

上面的`impl Shape for Round`和`impl<T: Round> Shape for T`是不一样的, 在前一种写法中, self 是`&Round`类型, 它是一个 trait object ，是胖指针. 在后一种写法中, self 是&T, T是具体类型 前一种写法是为 trait object增加一个成员方法; 而后一种写法是为所有的满足`T: Round`的具体类型增加一个成员方法. 所以上面的示例中，我们只能构造一个 trait object 之后才能调用 area()成员方法.

Rust 2018 edition开始, trait object 的语法会被要求加上 dyn 关键字即`impl Shape for dyn Round`.

## 面向对象
### 封装
封装就是对外显示的策略，在 Rust 中可以通过模块的机制来实现最外层的封装，并且每一个 Rust 文件都可以看作一个模块，模块内的元素可以通过 pub 关键字对外明示.

```rust
// --- second.rs
pub struct ClassName {
    field: i32,
}

impl ClassName {
    pub fn new(value: i32) -> ClassName {
        ClassName {
            field: value
        }
    }

    pub fn public_method(&self) {
        println!("from public method");
        self.private_method();
    }

    fn private_method(&self) {
        println!("from private method");
    }
}

// --- main.rs
mod second;
use second::ClassName;

fn main() {
    let object = ClassName::new(1024);
    object.public_method();
}
```

### 继承
继承是多态（Polymorphism）思想的实现, 多态指的是编程语言可以处理多种类型数据的代码. 在 Rust 中，通过特性（trait）实现多态.

总结地说，Rust 没有提供跟继承有关的语法糖，也没有官方的继承手段（完全等同于 Java 中的类的继承），但灵活的语法依然可以实现相关的功能.

## 闭包
闭包是可以保存进变量或作为参数传递给其他函数的匿名函数, 闭包相当于 Rust 中的 Lambda 表达式:
```rust
|参数1, 参数2, ...| -> 返回值类型 {
    // 函数体
}
```

## 并发
move:
```rust
use std::thread;

fn main() {
    let s = "hello";
   
    let handle = thread::spawn(move || {
        println!("{}", s);
    });

    handle.join().unwrap();
}
```

channel:
```rust
// 子线程获得了主线程的发送者 tx，并调用了它的 send 方法发送了一个字符串，然后主线程就通过对应的接收者 rx 接收到了
use std::thread;
use std::sync::mpsc;

fn main() {
    let (tx, rx) = mpsc::channel();

    thread::spawn(move || {
        let val = String::from("hi");
        tx.send(val).unwrap();
    });

    let received = rx.recv().unwrap();
    println!("Got: {}", received);
}
```

## unsafe
unsafe不过是把 Rust 编译器在编译器做的严格检查退步成为 C++ 的样子, 由开发者自己为其所撰写的代码的正确性做担保.