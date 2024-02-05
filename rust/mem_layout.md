# mem_layout
ref:
- [用了这么多年Rust终于搞明白了内存分布](https://zhuanlan.zhihu.com/p/624552143)

env:
- rust: Edition 2021, 1.77.0
- os: 64bit

## tupe
```rust
use std::mem;

fn main() {
    let _a: (char, u8, i32) = ('a', 7, 354); // Stack(4B, 1B, 4B) // Rust会选择Tuple中对齐值最大的元素为_a该元组的对齐值, 因此alignment是4

    println!("{}", mem::size_of::<(char, u8, i32)>()); // 12
    println!("{}", mem::align_of::<(char, u8, i32)>()); // 4
}
```

由于Rust有多种数据排布风格（默认的Rust风格，还有C语言风格，primitive和transparent风格），在Rust风格中，Rust可以对元组中的元素做任意重排，也包括padding的位置，Rust是根据其优化算法做出其认为最优的排序，对最终排序结果并没有统一规则.

因此可能的布局: 4B+1B+3B(padding)+4B

## Reference
```rust
use std::mem;

fn main() {
    let a: i8 = 6; // Stack(1B)
    let _b : &i8 = &a; // Stack(8B)

    println!("{}", mem::size_of::<i8>()); // 1
    println!("{}", mem::align_of::<i8>()); // 1

    println!("{}", mem::size_of::<&i8>()); // 8
    println!("{}", mem::align_of::<&i8>()); // 8
}
```

`&T和&mut T`在内存分布上规则一致，他们的区别是在使用方式和编译器处理方式上.

## Array/Vector
```rust
use std::mem;

fn main() {
    let _a: [i8; 3] = [1, 2, 3]; // 变量, 值都在Stack
    let _b: Vec<i8> = vec![1, 2, 3]; // _b在Stack, 值在Heap

    println!("{}", mem::size_of::<[i8; 3]>()); // 3
    println!("{}", mem::align_of::<[i8; 3]>()); // 1

    println!("{}", mem::size_of::<Vec<i8>>()); // 24
    println!("{}", mem::align_of::<Vec<i8>>()); // 8
}
```

数组Array是固定大小的，所以在创建的时候都指定好了长度；动态数组Vector可以自由伸缩的.

`_b`内存布局:
1. pointer(8B): 指向heap的值
1. cap(8B)
1. len(8B)

## slice
```rust
use std::mem;

fn main() {
    let _a: [i8; 3] = [1, 2, 3];
    let _b: Vec<i8> = vec![1, 2, 3];

    let _slice_1: &[i8] = &_a[0..2]; // 指向_a的值
    let _slice_2: &[i8] = &_b[0..2]; // 指向_b的值

    println!("{}", mem::size_of::<[i8; 3]>()); // 3
    println!("{}", mem::align_of::<[i8; 3]>()); // 1

    println!("{}", mem::size_of::<Vec<i8>>()); // 24
    println!("{}", mem::align_of::<Vec<i8>>()); // 8

    println!("{}", mem::size_of::<&[i8]>()); // 16
}
```



`_slice_<x>`(是胖指针)内存布局:
1. pointer(8B): 指向值
1. len(8B)

## String, str, &str
```rust
use std::mem;

fn main() {
    let s1: String = String::from("HELLOS"); // s1在Stack, 值在heap
    let _s2: &str = "ЗдP"; // д -> Russian Language // 这个string数据不会存储在堆heap上，而是会直接存在编译后的二进制中，同时具有static生命周期, 在程序的只读内存里
    let _s3: &str = &s1[1..3];

    println!("{}", mem::size_of::<String>()); // 24
    println!("{}", mem::align_of::<String>()); // 8

    println!("{}", mem::size_of::<&str>()); // 16
    println!("{}", mem::align_of::<&str>()); // 8
}
```

String是Vec的一个封装, String是保证UTF-8兼容的.

`_s2和_s3`是string slice

## Struct
### 1. unit-like Struct
由于并没有定义Data结构体的细节，Rust也不会为其分配任何内存

### 2. Struct with named fields && tuple-like struct
```rust
use std::mem;

struct Data {
    nums: Vec<usize>,
    dimension: (usize, usize),
}

fn main() {
    let _s1: Data = Data {
        nums: vec![1, 2, 3],
        dimension: (1, 2),
    }; // _s1在stack, 值在heap

    println!("{}", mem::size_of::<Data>()); // 40
    println!("{}", mem::align_of::<Data>()); // 8
}
```

## Enum
```rust
use std::mem;

enum HTTPStatus {
    OK,
    NOTFOUND,
}

enum Data {
    Empty,
    Number(i32),
    Array(Vec<i32>),
}

fn main() {
    let _s1: HTTPStatus = HTTPStatus::OK;

    println!("{}", mem::size_of::<HTTPStatus>()); // 1
    println!("{}", mem::align_of::<HTTPStatus>()); // 1

    println!("{}", mem::size_of::<Data>()); // 24
    println!("{}", mem::align_of::<Data>()); // 8
}
```

对于C-style enum, 在内存中, rust会根据该enum中最大的数来选择内存占用. 此例中没有指定就会默认Ok为0, NotFound为1, Rust选择占用1 byte的i8来存储enum.

Data占用24是rust优化了tag, 判断Data的前8B, 可以判断Array还是非Array, 再通过其他区分其他两个字段.

## Option
```rust
use std::mem::{size_of, align_of};

fn main() {
    println!("{}", size_of::<Option<Box<i32>>>()); // 8, rust优化了Enum tag: Rust对于类似Box这样的不允许为null的SmartPointer进行了优化, null和非null代替了tag
    println!("{}", align_of::<Option<Box<i32>>>()); // 8
}
```

Box会将原来的i32从栈放到堆, 然后Box会是一个指针指向原来的i32的堆地址

## Box
```rust
let t: (i32, String) = (5, "Hello".to_string); // 32B, t在stack,  "Hello"在heap
```

```rust
let t: (i32, String) = (5, "Hello".to_string);
let mut b = Box::new(t); // b在stack, t和其值在heap
```