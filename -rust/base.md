# rust
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

    > rust所有权系统还包括了从现代c++借鉴的RAII机制, 这是rust无gc但是可以安全管理内存的基石.

    为了实现内存安全, rust还具备具有的特性:
    1. 仿射类型(Affine Type), 该类型可用来表达rust所有权中的Move语义.

    借助类型, rust可在编译阶段对类型进行检查是否满足安全内存模型, 有效地阻止未定义行为的发生.

    内存安全bug和并发安全的bug的产生内在原因均是内存的不正当访问造成的. 借助装备了所有权的强大类型系统, rust还解决了并发安全问题. 它通过静态检查分析, 在编译期就能检查出多线程并发代码中所有的数据竞争问题.
1. 零成本抽象即代码表达能力不存在运行时开销.

    rust的抽象并不存在运行时的开销, 其一切都是在编译期完成的.

    rust的零成本抽象的基石是泛型和trait.
1. 实用性

    为了保证支持硬实时, rust借鉴了c++的确定性析构, RAII和智能指针, 用于自动地, 确定地管理内存, 从而避免了gc的引入.

    为了保证程序的鲁棒性, rust重新审视了错误处理机制. rust针对三类非正常情况: 失败, 错误和异常, 提供了专门的处理方式:
    - 失败: 使用断言工具
    - 错误: 基于返回值的分层处理 
    - 异常: rust将其看作无法被合理解决的问题, 提供了线程恐慌机制, 发生异常时, 线程可以安全地退出.

    为了兼容现有生态, rust支持方便且零成本的FFI机制, 兼容C-ABI, 在语言架构层面上将rust分为safe rust和unsafe rust两部分. unsafe专门和外部生态打交道, 因为rust编译器检查和跟踪的范围有限, 不能检查到与其链接的其他生态接口, 因此这些生态由自身来保证安全性. 总结就是, safe rust由rust编译器在编译时保证安全, unsafe rust开发者让编译器信任自身有能力保证安全.

### 编译
rust编译器是LLVM编译器的前端, 它将代码编译成LLVM IR, 然后通过LLVM编译成对应架构的机器码.

rust源码经过分词和解析生成AST(抽象语法树), 再进一步简化处理为HIR(High-level IR, 方便编译器做类型检查), 再进一步编译为MIR(middle IR, 在rust 1.12引入), 最后MIR被翻译为LLVM IR, 之后由LLVM编译成目标机器码.

引入MIR原因:
1. 缩短编译时间

    实现了增量编译, 仅重新编译更改过的部分.
1. 缩短执行时间

    进入llvm前实现更细颗粒度的优化, 单纯依赖llvm的优化颗粒度太粗, 增加了更多的优化空间
1. 更精确的类型检查

    实现更灵活的借用检查

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
    1. 常用的宏定义, 如println()!, assert!, panic!, vec!.

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
    1. 可选和错误处理类型Option和Result, 以及各种迭代器.
    
### 语句和表达式
rust语法分语句(statement, 要执行的一些操作和产生副作用的表达式)和表达式(expression, 主要用于计算求值).

语句又分:
- 声明语句(Declaration statement) : 用于声明各种语言项(item), 比如变量, 静态变量, 常量, 结构体, 函数等, 以及通过extern和use关键词引入的包和模块等.
- 表达式语句(expression statement) : 特指以分号结尾的表达式. 此类表达式求值结果会被舍弃, 并总是返回单元类型`()`.

rust编译器解析代码时, 如果遇到分号, 就会继续往后面执行; 如果碰到语句, 则执行语句; 如果碰到表达式, 则会对表达式求值, 如果分号后面什么都没有, 就会补上单元值;  当
遇到函数时, 会将函数体的花括号识别为块表达式(block expression, 由一对花括号和一系列表达式组成, 它总是返回块中最后一个表达式的值)

let创建的变量一般称为绑定(bingding).

rust的表达式可分为位置表达式(place expression)和值表达式(value expression), 即其他语言中的左值和右值.

通过位置表达式可对某个数据单元的内存进行读写.
值表达式一般只引用了某个存储单元地址中的数据, 它相当于数据值, 只能进行读操作.

**从语法角度讲, 位置表达式代表了持久化数据, 值表达式代表了临时数据.**

## FAQ
### 各编程语言中的类型系统
静态类型语言：在编译阶段确定所有变量的类型.
动态类型语言：在执行阶段确定所有变量的类型.
不允许隐式转换的是强类型，允许隐式转换的是弱类型.

- 没类型： 比如汇编语言，没有类型的概念，所有都只是一个数字
- 弱静态类型： 比如C/C++语言，可以定义类型，但是不强制执行，在不同类型之间自动转换
- 强静态类型： 比如Java，定义类型，并且用虚拟机检查类型
- 强动态类型： 比如Python和Ruby，动态推断类型而不需要定义，然后解释器会强制执行.
- 弱动态类型： Perl/PHP,js

![](/misc/img/rust/9zv7ejia98.jpeg)