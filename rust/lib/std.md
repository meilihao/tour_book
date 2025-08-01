# std
ref:
- [inside-rust-std-library](https://github.com/Warrenren/inside-rust-std-library)

RUST语言的设计目标是系统编程语言. 它需要考虑操作系统内核与用户态两种模型, 不像其他语言的标准库仅需要支持用户态模型. C语言在解决这个问题的方法是只提供用户态的标准库，os内核的库由各os自行实现. RUST的现代语言特性决定了标准库无法象C语言那样把操作系统内核及用户态程序区分成完全独立的两个部分，所以只能更细致的设计，做模块化的处理.

RUST标准库体系分为三个模块:
1. 语言核心库core

    core库适用于os内核及用户态编程, 是与cpu arch无关的可移植库, 主要内容：
    - 编译器内置固有(intrinsic)函数

        包括内存操作函数，数学函数，位操作函数，原子变量操作函数等， 这些函数通常与CPU硬件架构紧密相关，且一般需要汇编来提供最佳性能. intrinsic函数实际上也是对CPU指令的屏蔽层
    - 基本特征(trait)

        包括运算符(OPS)Trait, 编译器Marker Trait, 迭代器(Iterator)Trait, 类型转换Trait等
    - Option/Result类型
    - 基本数据类型

        包括整数类型, 浮点类型, bool, 字符类型和单元类型, 以及为这些类型实现基本特征和一些特有函数
    - 数组, slice和range类型

        为数组, slice和range类型实现基本特征和一些特有函数
    - 内存操作

        包括alloc模块, mem模块, ptr模块. rust中90%的unsafe语法都可归结到这3个模块.
    - 字符串及格式化
        
        为字符串类型实现基本特征和一些特有函数
    - 内部可变性类型

        包括`UnsafeCell<T>, Cell<T>, RefCell<T>`等, 以及为这些类型实现基本特征和一些特有函数
    - etc..

        包括FFI, 时间, 异步库等

1. 智能指针库alloc

    alloc库的所有类型都是基于堆, 包括智能指针类型，集合类型, 容器类型. 这些类型与为它们实现的函数和trait组成了alloc库的主体. alloc库仅依赖core库. alloc库适用于操作系统内核及用户态程序.

    主要内容:
    1. 内存申请和释放: Allocator Trait及其实现者Global单元类型
    1. 基础智能指针类型: `Box<T>, Rc<T>`
    1. 动态数组智能指针类型: `RawVec<T>, Vec<T>`
    1. 字符串智能指针类型: String
    1. 并发安全基础指针类型: `Arc<T>`
    1. 集合类型: `LinkList<T>, VecQueue<T>, BTreeSet<T>, BTreeMap<T>`等

1. 用户态std库

    std库建立在os syscall基础上, 仅适用于用户态编程. 它主要的工作是针对os资源设计rust的类型, trait和函数.

    主要内容:
    1. 对core库和alloc库的内容进行映射
    1. 实现进程管理和进程间通信
    1. 实现线程管理, 线程间互斥锁, 消息通信和其他线程相关内容
    1. 实现文件, 目录及os env
    1. 实现输入, 输出
    1. 实现网络通信

core库是基础, alloc库和std库基于core库.