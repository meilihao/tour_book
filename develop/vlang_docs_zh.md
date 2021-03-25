# V Documentation
> version: 21.3.23@35c60cf

## 介绍

V是一种静态类型的编译型编程语言，旨在构建可维护的软件.

它与Go类似，其设计也受到Oberon, Rust, Swift, Kotlin 和 Python 的影响.

V是一种非常简单的语言. 看完这篇文档大约需要一个小时，结束后你就会把这门语言学得差不多了.

本语言提倡用最小的抽象来编写简单清晰的代码.

尽管简单，但V给了开发者很大的权力. 任何你在其他语言中可以做到的事情，你也都可以在V中做到.

## 从源码安装
获得最新, 最完善的V的主要方式是__从源码开始__. 这很__简单__，通常__只需要几秒钟__.

### Linux, macOS, FreeBSD等:
你需要`git`和一个C编译器, 比如`tcc`, `gcc` 或 `clang`, 以及`make`:
```bash
git clone https://github.com/vlang/v
cd v
make
```

### Windows:
你需要`git`和一个C编译器, 比如`tcc`, `gcc`, `clang` 或 `msvc`:
```bash
git clone https://github.com/vlang/v
cd v
make.bat -tcc
```
NB: 如果你喜欢使用不同的C编译器，也可以将`-gcc`, `-msvc`, `-clang`中的一个作为参数传入make.bat，但是`-tcc`体积小, 速度快, 安装方便(因为V会自动下载一个已预制的二进制文件).

建议将这个文件夹添加到环境变量的PATH中. 这可以通过命令`v.exe symlink`来完成.

### Android
通过[vab](https://github.com/vlang/vab)也可以在Android上运行V图形应用.

V Android 依赖: **V**, **Java JDK** >= 8, Android **SDK + NDK**.

  1. Install dependencies (see [vab](https://github.com/vlang/vab))
  2. Connect your Android device
  3. Run:
  ```bash
  git clone https://github.com/vlang/vab && cd vab && v vab.v
  ./vab --device auto run /path/to/v/examples/sokol/particles
  ```
更多细节和故障排查, 请浏览[vab GitHub repository](https://github.com/vlang/vab).

## Table of Contents

<table>
    <tr><td width=33% valign=top>

* [Hello world](#hello-world)
* [Running a project folder](#running-a-project-folder-with-several-files)
* [Comments](#comments)
* [Functions](#functions)
    * [Returning multiple values](#returning-multiple-values)
* [Symbol visibility](#symbol-visibility)
* [Variables](#variables)
* [Types](#types)
    * [Strings](#strings)
    * [Numbers](#numbers)
    * [Arrays](#arrays)
    * [Fixed size arrays](#fixed-size-arrays)
    * [Maps](#maps)
* [Module imports](#module-imports)
* [Statements & expressions](#statements--expressions)
    * [If](#if)
    * [In operator](#in-operator)
    * [For loop](#for-loop)
    * [Match](#match)
    * [Defer](#defer)
* [Structs](#structs)
    * [Embedded structs](#embedded-structs)
    * [Default field values](#default-field-values)
    * [Short struct literal syntax](#short-struct-initialization-syntax)
    * [Access modifiers](#access-modifiers)
    * [Methods](#methods)

</td><td width=33% valign=top>

* [Unions](#unions)
* [Functions 2](#functions-2)
    * [Pure functions by default](#pure-functions-by-default)
    * [Mutable arguments](#mutable-arguments)
    * [Variable number of arguments](#variable-number-of-arguments)
    * [Anonymous & high order functions](#anonymous--high-order-functions)
* [References](#references)
* [Constants](#constants)
* [Builtin functions](#builtin-functions)
* [Printing custom types](#printing-custom-types)
* [Modules](#modules)
* [Types 2](#types-2)
    * [Interfaces](#interfaces)
    * [Enums](#enums)
    * [Sum types](#sum-types)
    * [Option/Result types & error handling](#optionresult-types-and-error-handling)
* [Generics](#generics)
* [Concurrency](#concurrency)
    * [Spawning Concurrent Tasks](#spawning-concurrent-tasks)
    * [Channels](#channels)
    * [Shared Objects](#shared-objects)
* [Decoding JSON](#decoding-json)
* [Testing](#testing)
* [Memory management](#memory-management)
* [ORM](#orm)

</td><td valign=top>

* [Writing documentation](#writing-documentation)
* [Tools](#tools)
    * [v fmt](#v-fmt)
    * [Profiling](#profiling)
* [Advanced Topics](#advanced-topics)
    * [Memory-unsafe code](#memory-unsafe-code)
    * [Structs with reference fields](#structs-with-reference-fields)
    * [sizeof and __offsetof](#sizeof-and-__offsetof)
    * [Calling C from V](#calling-c-from-v)
    * [Debugging generated C code](#debugging-generated-c-code)
    * [Conditional compilation](#conditional-compilation)
    * [Compile time pseudo variables](#compile-time-pseudo-variables)
    * [Compile-time reflection](#compile-time-reflection)
    * [Limited operator overloading](#limited-operator-overloading)
    * [Inline assembly](#inline-assembly)
    * [Translating C to V](#translating-c-to-v)
    * [Hot code reloading](#hot-code-reloading)
    * [Cross compilation](#cross-compilation)
    * [Cross-platform shell scripts in V](#cross-platform-shell-scripts-in-v)
    * [Attributes](#attributes)
    * [Goto](#goto)
* [Appendices](#appendices)
    * [Keywords](#appendix-i-keywords)
    * [Operators](#appendix-ii-operators)

</td></tr>
</table>

<!--
NB: there are several special keywords, which you can put after the code fences for v:
compile, live, ignore, failcompile, oksyntax, badsyntax, wip, nofmt
For more details, do: `v check-md`
-->

## Hello World


```v
fn main() {
	println('hello world')
}
```

将上面的片段保存到`hello.v`, 再执行`v run hello.v`.

> 上面假设你已按照[这里](https://github.com/vlang/v/blob/master/README.md#symlinking)所述, 用`v symlink`为V建立了软连接. 如果还没有，则必须手动输入V的路​​径.

祝贺你: 你刚刚编写并执行了你的第一个V程序!

你也可以通过`v hello.v`实现仅编译而不执行. 查看`v help`可获得所有受支持的命令信息.

从上面的示例中, 你可以看到是使用`fn`关键字来声明函数. 返回类型在函数名称之后指定. 在这个例子中, `main`不返回任何内容, 因此没有返回类型.

与许多其他语言(例如C, Go和Rust)一样, `main`是程序的入口.

`println`是为数不多的内置函数之一. 它将传递给它的值打印到stdout.

可以在一个文件程序中舍弃`fn main()`声明. 在编写小型程序，"脚本"或仅学习语言时很有用. 为简便起见, 本教程中将跳过`fn main()`.

这意味着V中的`hello world`程序非常简单.

```v
println('hello world')
```

## 运行包含多个文件的项目文件夹

假设你有一个包含多个`.v`文件的文件夹, 其中一个文件包含`main()`函数, 其他文件具有其他辅助函数. 它们可能是按功能进行组织的, 但仍不足以使其成为各自独立的可重用模块，同时你希望将它们全部编译到一个程序中.

在其他语言中, 你将必须使用include或构建系统来枚举所有文件, 将它们分别编译为目标文件, 然后将它们链接为一个最终的可执行文件.

但是, 在V中仅使用`v run`即可一起编译并运行整个文件夹中的.v文件. 该命令也支持传递参数, 因此可执行操作: `v run . --yourparam some_other_stuff`.

上面的命令首先会将这些文件编译为一个程序(以文件夹/项目命名), 然后执行该程序时会将`--yourparam some_other_stuff`作为CLI参数传递给它.

你的程序可以如下方式使用CLI参数： 
```v
import os

println(os.args)
```
NB: 运行成功后, V会删除生成的可执行文件. 如果你想保留它, 可以使用`v -keepc run .`代替, 或者直接用`v .`手动编译.

NB: 任何V编译器的标志都应该在`run`命令之前传递. 源文件/文件夹之后的所有参数内容都将原封不动地传递给编出的程序，V不会对其进行处理.

## 注释

```v
// This is a single line comment.
/*
This is a multiline comment.
   /* It can be nested. */
*/
```

## 函数

```v
fn main() {
	println(add(77, 33))
	println(sub(100, 50))
}

fn add(x int, y int) int {
	return x + y
}

fn sub(x int, y int) int {
	return x - y
}
```

同样，类型也是在参数名称之后.

就像在Go和C中一样，函数不能被重载. 这简化了代码, 提高了可维护性和可读性.

函数可以在声明之前使用: `add`和`sub`在`main`之后声明，但仍然可以在`main`中调用. 这对V中的所有声明都是如此, 而且不需要头文件或考虑文件和声明的顺序.

### 多值返回

```v
fn foo() (int, int) {
	return 2, 3
}

a, b := foo()
println(a) // 2
println(b) // 3
c, _ := foo() // ignore values using `_`
```

## 符号可见性

```v
pub fn public_function() {
}

fn private_function() {
}
```

函数默认是私有的(不导出). 要允许其他模块使用它们, 请在前面加上`pub`. 这同样适用于常量和类型.

注意: `pub`只能在一个命名的模块中使用. 关于创建模块的信息，请参见[Modules](#modules).

## 变量

```v
name := 'Bob'
age := 20
large_number := i64(9999999999)
println(name)
println(age)
println(large_number)
```

变量是用`:=`来声明和初始化的, 这是V中声明变量的唯一方式, 这意味着变量总是有一个初始值.

变量的类型是由右侧的值推断出来的. 要转换不同的类型, 可以使用类型转换: 表达式T(v)可将值v转换为类型T.

与大多数其他语言不同, V只允许在函数中定义变量. 全局(模块级)变量是不允许的. 在V中没有全局状态(详见[默认情况下的纯函数](#默认情况下的纯函数))

为了在不同的代码库中保持一致, 所有的变量和函数名都必须使用`snake_case`风格, 而类型名则必须使用`PascalCase`.

### 可变变量

```v
mut age := 20
println(age)
age = 21
println(age)
```

改变变量的值可使用`=`. 在V中，变量默认是不可改变的. 为了能够改变变量的值，你必须用`mut`声明它.

把第一行的`mut`去掉再试着编译上面的程序.

### 初始化 vs 赋值

注意: `:=` 和 `=` 的重大区别: `:=`用于声明和初始化, `=`用于赋值.

```v failcompile
fn main() {
    age = 21
}
```

由于未声明变量`age`, 因此该代码将无法编译. 在V中所有变量都需要声明.

```v
fn main() {
	age := 21
}
```

可以在一行中更改多个变量的值. 这样, 可以在没有中间变量的情况下交换它们的值.

```v
mut a := 0
mut b := 1
println('$a, $b') // 0, 1
a, b = b, a
println('$a, $b') // 1, 0
```

### 错误的声明

在开发模式下, 编译器将警告你尚未使用的变量(你将收到"unused variable"的警告). 在生产模式下(通过将`-prod`标志传递给v, 比如`v -prod foo.v`), 它根本不会编译(就像在Go中一样).

```v failcompile
fn main() {
    a := 10
    if true {
        a := 20 // error: redefinition of `a`
    }
    // warning: unused variable `a`
}
```

与大多数语言不同，不允许使用变量覆盖. 声明一个在父作用域中已声明的同名变量将导致编译错误.

但是，可以对导入的模块进行覆盖处理，因为在某些情况下它非常有用： 
```v ignore
import ui
import gg

fn draw(ctx &gg.Context) {
    gg := ctx.parent.get_ui().gg
    gg.draw_rect(10, 10, 100, 50)
}
```

## 类型

### 基础类型

```v ignore
bool

string

i8    i16  int  i64      i128 (soon)
byte  u16  u32  u64      u128 (soon)

rune // represents a Unicode code point

f32 f64

byteptr, voidptr, charptr, size_t // these are mostly used for C interoperability

any // similar to C's void* and Go's interface{}
```

请注意, 与C和Go不同, `int`总是一个32位的整数.

V中的所有运算符的两边必须是相同类型的值这一规则有一个例外: 如果一边的基础类型完全适合于另一边类型的数据范围, 就可以自动推导. 下面是允许的可能性:

```v ignore
   i8 → i16 → int → i64
                  ↘     ↘
                    f32 → f64
                  ↗     ↗
 byte → u16 → u32 → u64 ⬎
      ↘     ↘     ↘      ptr
   i8 → i16 → int → i64 ⬏
```

例如, 一个`int`值可以自动提升到`f64`或`i64`，但不能提升到u32(u32意味着负值的符号丢失). 然而, 从`int`到`f32`的提升目前是自动完成的(但对于大值来说可能会导致精度损失).

像`123`或`4.56`这样的字面量会以特殊的方式处理. 它们不会导致类型提升, 但是当它们的类型需要确定时, 它们会分别默认为int和f64.


```v nofmt
u := u16(12)
v := 13 + u    // v is of type `u16` - no promotion
x := f32(45.6)
y := x + 3.14  // x is of type `f32` - no promotion
a := 75        // a is of type `int` - default for int literal
b := 14.7      // b is of type `f64` - default for float literal
c := u + a     // c is of type `int` - automatic promotion of `u`'s value
d := b + x     // d is of type `f64` - automatic promotion of `x`'s value
```

### 字符串

```v
name := 'Bob'
println(name.len)
println(name[0]) // indexing gives a byte B
println(name[1..3]) // slicing gives a string 'ob'
windows_newline := '\r\n' // escape special characters like in C
assert windows_newline.len == 2
```

在V中, 字符串是一个只读的字节数组. 字符串数据使用UTF-8编码. 字符串的值是不可改变的, 你不能对元素进行修改:

```v failcompile
mut s := 'hello 🌎'
s[0] = `H` // not allowed
```
> error: cannot assign to `s[i]` since V strings are immutable

请注意: 索引一个字符串将产生一个`byte`, 而不是一个`rune`. 索引对应的是字符串中的字节，而不是Unicode码点.

字符的类型是`rune`. 要表示它们, 使用使用"`"包裹:

```v
rocket := `🚀`
assert 'aloha!'[0] == `a`
```

单引号和双引号都可以用来表示字符串. 为了保持一致性, vfmt会将双引号转换为单引号, 除非字符串中包含一个单引号字符.

对于raw字符串, 请在前面加上r, raw字符串就不会被转义:

```v
s := r'hello\nworld'
println(s) // "hello\nworld"
```

字符可以很容易地转换为整数:

```v
s := '42'
n := s.int() // 42
```

### 字符串插值

基础的插值语法非常简单: 在变量名前使用`$`即可. 变量将被转换为一个字符串并嵌入到字面量中:
```v
name := 'Bob'
println('Hello, $name!') // Hello, Bob!
```
它也适用于字段: `'age = $user.age'`.
如果你需要更复杂的表达式请使用`${}`: `'can register = ${user.age > 13}'`.

也支持类似于 C 语言的`printf()`的格式指定符. `f`, `g`, `x`等是可选的, 用于指定了输出格式. 编译器会考虑到存储大小, 所以没有`hd`和`llu`.

```v
x := 123.4567
println('x = ${x:4.2f}')
println('[${x:10}]') // pad with spaces on the left => [   123.457]
println('[${int(x):-10}]') // pad with spaces on the right => [123       ]
println('[${int(x):010}]') // pad with zeros on the left => [0000000123]
```

### 字符串操作符

```v
name := 'Bob'
bobby := name + 'by' // + is used to concatenate strings
println(bobby) // "Bobby"
mut s := 'hello '
s += 'world' // `+=` is used to append to a string
println(s) // "hello world"
```

V中的所有运算符的两边必须具有相同类型. 你不能将整数与字符串连接起来:

```v failcompile
age := 10
println('age = ' + age) // not allowed
```
> error: infix expr: cannot use `int` (right expression) as `string`

我们可使用将`age`转成`string`:

```v
age := 11
println('age = ' + age.str())
```

或使用字符串内插法(首选):

```v
age := 12
println('age = $age')
```

### 数值

```v
a := 123
```

这将把123的值分配给`a`. 默认情况下, `a`的值为`int`类型.

你也可以用十六进制, 二进制或八进制来表示整数:

```v
a := 0x7B
b := 0b01111011
c := 0o173
```

所有这些都将被分配相同的值123. 它们的类型都是`int`，不管你用什么写法.

V还支持用`_`作为分隔符写数字:

```v
num := 1_000_000 // same as 1000000
three := 0b0_11 // same as 0b11
float_num := 3_122.55 // same as 3122.55
hexa := 0xF_F // same as 255
oct := 0o17_3 // same as 0o173
```

如果你想要一个不同类型的整数, 你可以使用类型转换:

```v
a := i64(123)
b := byte(42)
c := i16(12345)
```

浮点数的赋值方法也是一样的:

```v
f := 1.0
f1 := f64(3.14)
f2 := f32(3.14)
```

如果你没有明确指定类型, 默认情况下, 浮点字面量将是`f64`的类型.

### 数组

```v
mut nums := [1, 2, 3]
println(nums) // "[1, 2, 3]"
println(nums[1]) // "2"
nums[1] = 5
println(nums) // "[1, 5, 3]"
println(nums.len) // "3"
nums = [] // The array is now empty
println(nums.len) // "0"
// Declare an empty array:
users := []int{}
```

数组的类型由第一个元素决定:
* `[1, 2, 3]` 是int类型的数组 (`[]int`).
* `['a', 'b']` 是string类型的数组 (`[]string`).

用户可以明确指定第一个元素的类型：`[byte(16)，32，64，128]`. V数组是同质的(所有元素必须具有相同的类型).
这意味着像`[1, 'a']`这样的代码将无法编译.

`.len`字段会返回数组的长度. 注意这是一个只读字段, 并且用户不能修改. 在V中, 导出的字段默认为只读. 参考[访问修改器](#访问修改器)

#### 数值操作符

```v
mut nums := [1, 2, 3]
nums << 4
println(nums) // "[1, 2, 3, 4]"
// append array
nums << [5, 6, 7]
println(nums) // "[1, 2, 3, 4, 5, 6, 7]"
mut names := ['John']
names << 'Peter'
names << 'Sam'
// names << 10  <-- This will not compile. `names` is an array of strings.
println(names.len) // "3"
println('Alex' in names) // "false"
```

`<<<`是一个运算符, 它将一个值追加到数组的末尾, 它也可以追加整个数组.

`val in array`表示如果数组中包含`val`, 则返回true. 参见[`in`运算符](#in-运算符).

#### 初始化数组属性

在初始化过程中, 你可以指定数组的容量(`cap`), 初始长度(`len`) 和默认元素(`init`):

```v
arr := []int{len: 5, init: -1}
// `[-1, -1, -1, -1, -1]`
```

设置容量可以提高插入的性能, 因为它减少了所需的重新分配次数:

```v
mut numbers := []int{cap: 1000}
println(numbers.len) // 0
// Now appending elements won't reallocate
for i in 0 .. 1000 {
	numbers << i
}
```
注意: 上面的代码使用了[range `for`](#range-for)语句.

#### 数组方法

所有的数组都可以很容易地用`println(arr)`打印出来, 并用`s :=arr.str()`转换为一个字符串.

用`.clone()`可复制数组中的数据:

```v
nums := [1, 2, 3]
nums_copy := nums.clone()
```

数组可以通过`.filter()'和`.map()'有效地过滤和映射:

```v
nums := [1, 2, 3, 4, 5, 6]
even := nums.filter(it % 2 == 0)
println(even) // [2, 4, 6]
// filter can accept anonymous functions
even_fn := nums.filter(fn (x int) bool {
	return x % 2 == 0
})
println(even_fn)
words := ['hello', 'world']
upper := words.map(it.to_upper())
println(upper) // ['HELLO', 'WORLD']
// map can also accept anonymous functions
upper_fn := words.map(fn (w string) string {
	return w.to_upper()
})
println(upper_fn) // ['HELLO', 'WORLD']
```

`it`是一个内置的变量, 它指的是当前在filter/map方法中处理的元素.

此外, `.any()`和`.all()`可以用来方便地测试满足条件的元素.

```v
nums := [1, 2, 3]
println(nums.any(it == 2)) // true
println(nums.all(it >= 2)) // false
```

#### 多维数组

数组可以有多个维度.

二维数组的例子:
```v
mut a := [][]int{len: 2, init: []int{len: 3}}
a[0][1] = 2
println(a) // [[0, 2, 0], [0, 0, 0]]
```

3维数组的例子:
```v
mut a := [][][]int{len: 2, init: [][]int{len: 3, init: []int{len: 2}}}
a[0][1][1] = 2
println(a) // [[[0, 0], [0, 2], [0, 0]], [[0, 0], [0, 0], [0, 0]]]
```

#### 数组排序

对各种数组进行排序是非常简单和直观的. 特殊变量`a`和`b`可用于自定义排序的条件.

```v
mut numbers := [1, 3, 2]
numbers.sort() // 1, 2, 3
numbers.sort(a > b) // 3, 2, 1
```

```v
struct User {
	age  int
	name string
}

mut users := [User{21, 'Bob'}, User{20, 'Zarkon'}, User{25, 'Alice'}]
users.sort(a.age < b.age) // sort by User.age int field
users.sort(a.name > b.name) // reverse sort by User.name string field
```

#### 数组slice

Slice是数组的部分, 它们表示两个用`...`运算符分隔的索引之间的每个元素. 右边的指数必须大于或等于左边的索引.

如果没有右侧的索引, 则假定为数组的长度. 如果一个左侧指数不存在, 则假设为0.

```v
nums := [0, 10, 20, 30, 40]
println(nums[1..4]) // [10, 20, 30]
println(nums[..4]) // [0, 10, 20, 30]
println(nums[1..]) // [10, 20, 30, 40]
```

所有的数组操作都适用于切片.
分片可以被追加到同一类型的数组上:

```v
array_1 := [3, 5, 4, 7, 6]
mut array_2 := [0, 1]
array_2 << array_1[..3]
println(array_2) // [0, 1, 3, 5, 4]
```

### Fixed size arrays

V还支持固定大小的数组. 与普通数组不同，它们的长度是固定的. 你不能给它们追加元素，也不能缩小它们. 你只能修改它们的元素.

不过, 与普通数组不同, 访问固定大小的数组元素的效率更高, 它们比普通数组占用更少的内存. 它们的数据在堆栈上, 所以你可把它们作为缓冲区使用, 而不需要额外的堆分配.

大多数方法都被定义为在普通数组上，而不是在固定大小的数组上. 但你可以通过分片将固定大小的数组转换为普通数组:
```v
mut fnums := [3]int{} // fnums is a fixed size array with 3 elements.
fnums[0] = 1
fnums[1] = 10
fnums[2] = 100
println(fnums) // => [1, 10, 100]
println(typeof(fnums).name) // => [3]int

anums := fnums[0..fnums.len]
println(anums) // => [1, 10, 100]
println(typeof(anums).name) // => []int
```
请注意: 切片会导致固定大小的数组的数据会被复制到新创建的普通数组中.

### Map

```v
mut m := map[string]int{} // a map with `string` keys and `int` values
m['one'] = 1
m['two'] = 2
println(m['one']) // "1"
println(m['bad_key']) // "0"
println('bad_key' in m) // Use `in` to detect whether such key exists
m.delete('two')
```
map可使用string, rune, integer, float 或 voidptr作为key. 

整个map可以使用这个简短的语法来初始化:
```v
numbers := map{
	1: 'one'
	2: 'two'
}
println(numbers)
```

如果一个key没找到就会返回对应值的零值:

```v
sm := map{
	'abc': 'xyz'
}
val := sm['bad_key']
println(val) // ''
```
```v
intm := map{
	1: 1234
	2: 5678
}
s := intm[3]
println(s) // 0
```

也可以使用`or {}`代码块来处理丢失的key:

```v
mm := map[string]int{}
val := mm['bad_key'] or { panic('key not found') }
```

同样的可选检查也适用于数组:

```v
arr := [1, 2, 3]
large_index := 999
val := arr[large_index] or { panic('out of bounds') }
```

## Module导入

关于创建module可参考 [Modules](#modules).

Module导入使用`import`关键词:

```v
import os

fn main() {
	// read text from stdin
	name := os.input('Enter your name: ')
	println('Hello, $name!')
}
```
这个程序可以使用`os`模块中的任何公共定义，如`input`函数. 见[标准库](https://modules.vlang.io/)文档中的常用模块及其公共符号的列表.

默认情况下, 每次调用外部函数时都必须指定模块前缀. 这在一开始可能会显得很啰嗦，但它使代码更易读, 并且更容易理解 - 总是很清楚地知道哪个函数从
哪个模块中来. 这在大型代码库中特别有用.

循环导入是不允许的, 这与Go一样.

### 选择性导入

你也可以直接从模块中导入特定的函数和类型:

```v
import os { input }

fn main() {
	// read text from stdin
	name := input('Enter your name: ')
	println('Hello, $name!')
}
```
注意: 常量不允许这样做, 它们必须总是有前缀.

你可以同时导入几个特定的符号:

```v
import os { input, user_os }

name := input('Enter your name: ')
println('Name: $name')
os := user_os()
println('Your OS is ${os}.')
```

### Module导入支持alias

任何导入的模块名都可以使用`as`关键字进行重命名.

注意：除非你创建了`mymod/sha256.v`, 否则这个例子不会被编译.
```v failcompile
import crypto.sha256
import mymod.sha256 as mysha256

fn main() {
    v_hash := sha256.sum('hi'.bytes()).hex()
    my_hash := mysha256.sum('hi'.bytes()).hex()
    assert my_hash == v_hash
}
```

你不能对一个导入的函数或类型进行别名. 但是, 你可以重新声明一个类型.

```v
import time
import math

type MyTime = time.Time

fn (mut t MyTime) century() int {
	return int(1.0 + math.trunc(f64(t.year) * 0.009999794661191))
}

fn main() {
	mut my_time := MyTime{
		year: 2020
		month: 12
		day: 25
	}
	println(time.new_time(my_time).utc_string())
	println('Century: $my_time.century()')
}
```

## 语句和表达式

### If

```v
a := 10
b := 20
if a < b {
	println('$a < $b')
} else if a > b {
	println('$a > $b')
} else {
	println('$a == $b')
}
```

`if`语句非常直接, 与其他大多数语言类似. 与其他类似C语言不同的是, 条件周围没有括号, 而且总是需要括号(包裹代码块).

`if`可以作为表达式使用:

```v
num := 777
s := if num % 2 == 0 { 'even' } else { 'odd' }
println(s)
// "odd"
```

#### 类型检查和强制转换
You can check the current type of a sum type using `is` and its negated form `!is`.
你可以使用`is`和它的否定形式`!is`来检查一个和类型的当前类型.

你可以在`if`中进行:
```v
struct Abc {
	val string
}

struct Xyz {
	foo string
}

type Alphabet = Abc | Xyz

x := Alphabet(Abc{'test'}) // sum type
if x is Abc {
	// x is automatically casted to Abc and can be used here
	println(x)
}
if x !is Abc {
	println('Not Abc')
}
```
或使用`match`:
```v oksyntax
match x {
	Abc {
		// x is automatically casted to Abc and can be used here
		println(x)
	}
	Xyz {
		// x is automatically casted to Xyz and can be used here
		println(x)
	}
}
```

这也适用于struct的字段:
```v
struct MyStruct {
	x int
}

struct MyStruct2 {
	y string
}

type MySumType = MyStruct | MyStruct2

struct Abc {
	bar MySumType
}

x := Abc{
	bar: MyStruct{123} // MyStruct will be converted to MySumType type automatically
}
if x.bar is MyStruct {
	// x.bar is automatically casted
	println(x.bar)
}
match x.bar {
	MyStruct {
		// x.bar is automatically casted
		println(x.bar)
	}
	else {}
}
```

可变的变量可以发生变化, 进行类型转换是不安全的. 然而, 有时尽管变量是可变的, 但还是需要进行类型转换. 在这种情况下, 开发者必须用`mut`关键字来标记表达式, 来告诉编译器你知道自己在做什么.

它的方式是这样的:
```v oksyntax
mut x := MySumType(MyStruct{123})
if mut x is MyStruct {
	// x is casted to MyStruct even if it's mutable
	// without the mut keyword that wouldn't work
	println(x)
}
// same with match
match mut x {
	MyStruct {
		// x is casted to MyStruct even it's mutable
		// without the mut keyword that wouldn't work
		println(x)
	}
}
```

### In操作符

`in`允许检查一个数组或一个map是否包含一个元素.
相反操作用`!in`.

```v
nums := [1, 2, 3]
println(1 in nums) // true
println(4 !in nums) // true
m := map{
	'one': 1
	'two': 2
}
println('one' in m) // true
println('three' !in m) // true
```

它对于书写布尔表达式也很有用, 使其更清晰, 更紧凑:

```v
enum Token {
	plus
	minus
	div
	mult
}

struct Parser {
	token Token
}

parser := Parser{}
if parser.token == .plus || parser.token == .minus || parser.token == .div || parser.token == .mult {
	// ...
}
if parser.token in [.plus, .minus, .div, .mult] {
	// ...
}
```

V优化了这种表达方式. 所以上面两个`if`语句产生的机器代码是一样的, 都没有创建数组.

### For循环

V只有一个循环关键词: `for`, 但有多种形式.

#### `for`/`in`

这是最常见的形式. 你可以在数组, map或数值范围中使用它.

##### Array `for`

```v
numbers := [1, 2, 3, 4, 5]
for num in numbers {
	println(num)
}
names := ['Sam', 'Peter']
for i, name in names {
	println('$i) $name')
	// Output: 0) Sam
	//         1) Peter
}
```

`for value in arr`形式用于遍历一个数组的元素. 如果需要索引, 可以使用另一种形式`for index, value in arr`.

注意: 这个值是只读的. 如果你需要在循环时修改数组, 你需要将元素声明为可变的.

```v
mut numbers := [0, 1, 2]
for mut num in numbers {
	num++
}
println(numbers) // [1, 2, 3]
```
当一个标识符只是一个下划线时, 它将被忽略.

##### Map `for`

```v
m := map{
	'one': 1
	'two': 2
}
for key, value in m {
	println('$key -> $value')
	// Output: one -> 1
	//         two -> 2
}
```

通过使用一个下划线作为标识符, 可以忽略任何一个键或值.
```v
m := map{
	'one': 1
	'two': 2
}
// iterate over keys
for key, _ in m {
	println(key)
	// Output: one
	//         two
}
// iterate over values
for _, value in m {
	println(value)
	// Output: 1
	//         2
}
```

##### Range `for`

```v
// Prints '01234'
for i in 0 .. 5 {
	print(i)
}
```
`low..high`指的是一个*排他性*的范围, 代表从`low`开始到`high`(不包括`high`)中的所有数值.

#### Condition `for`

```v
mut sum := 0
mut i := 0
for i <= 100 {
	sum += i
	i++
}
println(sum) // "5050"
```

这种形式的循环类似于其他语言中的`while`循环. 一旦布尔条件值为false, 循环将停止迭代. 同样地, 条件周围没有括号, 而且代码块总是需要括号.

#### Bare `for`

```v
mut num := 0
for {
	num += 2
	if num >= 10 {
		break
	}
}
println(num) // "10"
```

这个条件可以省略, 会导致无限循环.

#### C `for`

```v
for i := 0; i < 10; i += 2 {
	// Don't print 6
	if i == 6 {
		continue
	}
	println(i)
}
```

最后是传统的C风格的`for`循环. 它比`while`形式更安全. 因为使用后者, 很容易忘记更新计数器, 而导致卡在一个无限循环中.

这里`i`不需要用`mut`来声明, 因为根据定义, 它是可变的.

#### 带标签的break/continue

`break`和`continue`默认控制的是最里面的`for`循环. 你也可以使用`break` 和 `continue` 后面的标签名来跳转到外部的`for`循环:
循环。

```v
outer: for i := 4; true; i++ {
	println(i)
	for {
		if i < 7 {
			continue outer
		} else {
			break outer
		}
	}
}
```
label必须紧接在外部循环之前.
上面的代码会打印:
```
4
5
6
7
```

### Match

```v
os := 'windows'
print('V is running on ')
match os {
	'darwin' { println('macOS.') }
	'linux' { println('Linux.') }
	else { println(os) }
}
```

match语句是编写一系列`if - else`语句的较短方法. 当找到一个匹配的分支时, 将执行其中的语句块. 当没有其他分支匹配时则执行else分支.

```v
number := 2
s := match number {
	1 { 'one' }
	2 { 'two' }
	else { 'many' }
}
```

match表达式从匹配分支返回最终表达式的值.

```v
enum Color {
	red
	blue
	green
}

fn is_red_or_blue(c Color) bool {
	return match c {
		.red, .blue { true } // comma can be used to test multiple values
		.green { false }
	}
}
```

match语句也可以使用简写的`.variant_here`语法作为`enum`变体的分支. 当所有的分支都是无穷尽的时候, 此时不允许使用`else`分支.


```v
c := `v`
typ := match c {
	`0`...`9` { 'digit' }
	`A`...`Z` { 'uppercase' }
	`a`...`z` { 'lowercase' }
	else { 'other' }
}
println(typ)
// 'lowercase'
```

你也可以使用范围作为`match`模式. 如果数值在分支的范围内, 则该分支将被执行.

请注意: 范围使用`...`(三点)而不是`...`(两点). 这就是因为范围是*包含*最后一个元素的, 而不是排他性的(比如`..`范围). 在匹配分支中使用`.`将引发一个错误.

注意: `match`作为表达式不能用于`for`循环和`if`语句.

### Defer

defer语句会推迟执行一组语句, 直到外围函数返回.

```v
import os

fn read_log() {
	mut ok := false
	mut f := os.open('log.txt') or { panic(err.msg) }
	defer {
		f.close()
	}
	// ...
	if !ok {
		// defer statement will be called here, the file will be closed
		return
	}
	// ...
	// defer statement will be called here, the file will be closed
}
```

## Structs

```v
struct Point {
	x int
	y int
}

mut p := Point{
	x: 10
	y: 20
}
println(p.x) // Struct fields are accessed using a dot
// Alternative literal syntax for structs with 3 fields or fewer
p = Point{10, 20}
assert p.x == 10
```

### Heap structs

struct是在堆栈上分配的. 要在堆上分配一个struct并获得对它的引用需使用`&`前缀:

```v
struct Point {
	x int
	y int
}

p := &Point{10, 10}
// References have the same syntax for accessing fields
println(p.x)
```

`p`的类型是`&Point`. 它是对`Point`的[引用](#references). 引用类似于Go指针和C++引用.

### 嵌入式struct

V不允许子类，但它支持嵌入式struct:

```v
struct Widget {
mut:
	x int
	y int
}

struct Button {
	Widget
	title string
}

mut button := Button{
	title: 'Click me'
}
button.x = 3
```
如果没有嵌入, 我们就必须给`Widget`字段命名, 然后做以下操作:

```v oksyntax
button.widget.x = 3
```

### 字段默认值

```v
struct Foo {
	n   int    // n is 0 by default
	s   string // s is '' by default
	a   []int  // a is `[]int{}` by default
	pos int = -1 // custom default value
}
```

在创建结构的过程中, 所有的结构字段默认为零值. 数组和map字段会被分配. 但也可以定义自定义的默认值.

### Required fields

```v
struct Foo {
	n int [required]
}
```

你可以用`[required]`属性标记一个结构体字段, 告诉V当创建该结构的实例时, 必须初始化该字段.

由于字段`n`没有被显式初始化, 这个例子将无法编译:
```v failcompile
_ = Foo{}
```

<a id='short-struct-initialization-syntax' />

### 简短的struct字面量语法

```v
struct Point {
	x int
	y int
}

mut p := Point{
	x: 10
	y: 20
}
// you can omit the struct name when it's already known
p = {
	x: 30
	y: 4
}
assert p.y == 4
```

省略struct名也可以用于返回一个struct的字面量或作为函数参数传递一个struct字面量.

#### 尾随struct的字面量参数

V doesn't have default function arguments or named arguments, for that trailing struct
literal syntax can be used instead

V没有默认的函数参数或命名的参数, 但可用尾随struct的字面量参数来代替:

```v
struct ButtonConfig {
	text        string
	is_disabled bool
	width       int = 70
	height      int = 20
}

struct Button {
	text   string
	width  int
	height int
}

fn new_button(c ButtonConfig) &Button {
	return &Button{
		width: c.width
		height: c.height
		text: c.text
	}
}

button := new_button(text: 'Click me', width: 100)
// the height is unset, so it's the default value
assert button.height == 20
```

如你所见, struct名称和花括号都可以省略, 而不是:

```v oksyntax nofmt
new_button(ButtonConfig{text:'Click me', width:100})
```

这仅适用于为最后一个参数是struct的函数.

### 访问修改器

struct字段默认是私有的, 不可变的(使得结构也是不可变的). 它们的访问修饰符可以用`pub`和`mut`, 总共有5个可能的选项.

```v
struct Foo {
	a int // private immutable (default)
mut:
	b int // private mutable
	c int // (you can list multiple fields with the same access modifier)
pub:
	d int // public immutable (readonly)
pub mut:
	e int // public, but mutable only in parent module
__global:
	// (not recommended to use, that's why the 'global' keyword starts with __)
	f int // public and mutable both inside and outside parent module
}
```

例如, 这里是在`builtin`模块中定义的`string`类型:

```v ignore
struct string {
    str byteptr
pub:
    len int
}
```

从这个定义中不难看出, `string`是一个不可改变的类型, 它包含字符串数据的字节指针在`builtin`之外根本无法访问. `len`字段是公共的, 但是不可变的:
```v failcompile
fn main() {
    str := 'hello'
    len := str.len // OK
    str.len++      // Compilation error
}
```

这意味着在V中定义公开可读字段非常容易, 不需要使用getters/setters或属性.

## Methods

```v
struct User {
	age int
}

fn (u User) can_register() bool {
	return u.age > 16
}

user := User{
	age: 10
}
println(user.can_register()) // "false"
user2 := User{
	age: 20
}
println(user2.can_register()) // "true"
```

V没有类, 但你可以在类型上定义方法. 一个方法是一个带有特殊接受者参数的函数, 接受者会出现在它自己的参数列表中，位于`fn`关键字和方法名之间.
接收者出现在它自己的参数列表中，位于`fn`关键字和方法名之间。
方法必须与接受者类型在同一个模块中.

在这个例子中, `can_register`方法有一个名为`u`, 类型是`User`的接收器. 惯例是不使用诸如`self`或`this`这样的接收者名称, 
但名字要简短，最好是一个字母长度.

## Unions

就像struct一样，union支持嵌入.

```v
struct Rgba32_Component {
	r byte
	g byte
	b byte
	a byte
}

union Rgba32 {
	Rgba32_Component
	value u32
}

clr1 := Rgba32{
	value: 0x008811FF
}

clr2 := Rgba32{
	Rgba32_Component: {
		a: 128
	}
}

sz := sizeof(Rgba32)
unsafe {
	println('Size: ${sz}B,clr1.b: $clr1.b,clr2.b: $clr2.b')
}
```

Output: `Size: 4B, clr1.b: 136, clr2.b: 0`

union成员的访问必须在一个`unsafe`块中进行.

请注意: 嵌入的struct参数不一定按所列顺序存储.

## Functions 2

### 默认是纯函数

V函数默认为纯函数, 这意味着它们的返回值只由它们的参数决定，而且函数的计算求值时没有任何副作用(除了I/O).
参数，而且它们的评估没有任何副作用（）。

这是因为缺少全局变量和所有函数参数在默认情况下是不可变的, 即使在传递[引用](#引用)时也是如此.

然而, V并不是一种纯粹的函数式语言.

有一个编译器标志来启用全局变量(`--enable-globals`), 但这是一个很重要的标志, 用于低级应用, 如内核和驱动程序.

### 可变参数

可以通过使用关键字 `mut`来修改函数参数:

```v
struct User {
	name string
mut:
	is_registered bool
}

fn (mut u User) register() {
	u.is_registered = true
}

mut user := User{}
println(user.is_registered) // "false"
user.register()
println(user.is_registered) // "true"
```

在本例中, receiver(它只是第一个参数)被标记为可变的，因此register()可以更改user对象. 同样适用于没有接收者参数的函数:

```v
fn multiply_by_2(mut arr []int) {
	for i in 0 .. arr.len {
		arr[i] *= 2
	}
}

mut nums := [1, 2, 3]
multiply_by_2(mut nums)
println(nums)
// "[2, 4, 6]"
```

注意: 调用这个函数时, 必须在`nums`前面加上`mut`. 很明显, 这使得被调用的函数可修改值.

最好是返回值而不是修改参数. 修改参数应该只在应用程序的性能关键部分进行, 以减少分配和复制.

出于这个原因, V不允许修改基础类型的参数(例如整数). 只有更复杂的类型, 如数组和map才可以被修改.

使用`user.register()` 或 `user = register(user)` 代替 `register(mut user)`.

#### struct更新语法

V使它很容易返回一个对象的修改版本:

```v
struct User {
	name          string
	age           int
	is_registered bool
}

fn register(u User) User {
	return {
		...u
		is_registered: true
	}
}

mut user := User{
	name: 'abc'
	age: 23
}
user = register(user)
println(user)
```

### 可变的参数数量

```v
fn sum(a ...int) int {
	mut total := 0
	for x in a {
		total += x
	}
	return total
}

println(sum()) // 0
println(sum(1)) // 1
println(sum(2, 3)) // 5
// using array decomposition
a := [2, 3, 4]
println(sum(...a)) // <-- using prefix ... here. output: 9
b := [5, 6, 7]
println(sum(...b)) // output: 18
```

### 匿名和高阶函数

```v
fn sqr(n int) int {
	return n * n
}

fn cube(n int) int {
	return n * n * n
}

fn run(value int, op fn (int) int) int {
	return op(value)
}

fn main() {
	// Functions can be passed to other functions
	println(run(5, sqr)) // "25"
	// Anonymous functions can be declared inside other functions:
	double_fn := fn (n int) int {
		return n + n
	}
	println(run(5, double_fn)) // "10"
	// Functions can be passed around without assigning them to variables:
	res := run(5, fn (n int) int {
		return n + n
	})
	println(res) // "10"
	// You can even have an array/map of functions:
	fns := [sqr, cube]
	println(fns[0](10)) // "100"
	fns_map := map{
		'sqr':  sqr
		'cube': cube
	}
	println(fns_map['cube'](2)) // "8"
}
```

## 引用

```v
struct Foo {}

fn (foo Foo) bar_method() {
	// ...
}

fn bar_function(foo Foo) {
	// ...
}
```

如果一个函数参数是不可改变的(比如上面例子中的`foo`), V可以通过值或引用来传递, 这由编译器决定, 而开发者不需要考虑这个问题.

你不再需要考虑是否应该通过值还是引用来传递struct.

你可以通过增加`&`确保结构体总是通过引用来传递的:

```v
struct Foo {
	abc int
}

fn (foo &Foo) bar() {
	println(foo.abc)
}
```

`foo`仍然是不可改变的. 为此必须使用`(mut foo Foo)`.

一般来说，V的引用类似于Go指针和C++引用. 例如, 一个通用的树结构的定义是这样的:

```v wip
struct Node<T> {
    val   T
    left  &Node
    right &Node
}
```

## 常亮

```v
const (
	pi    = 3.14
	world = '世界'
)

println(pi)
println(world)
```

常量用`const`来声明. 它们只能被定义为在模块级(需在函数外).
常量值永远不能被改变. 你也可以在模块中声明单个常量:

```v
const e = 2.71828
```

V常量比大多数语言更灵活. 你可以分配更复杂的值:

```v
struct Color {
	r int
	g int
	b int
}

fn rgb(r int, g int, b int) Color {
	return Color{
		r: r
		g: g
		b: b
	}
}

const (
	numbers = [1, 2, 3]
	red     = Color{
		r: 255
		g: 0
		b: 0
	}
	// evaluate function call at compile-time*
	blue    = rgb(0, 0, 255)
)

println(numbers)
println(red)
println(blue)
```
\* WIP - 目前函数调用在程序启动时进行评估.

全局变量通常是不被允许的, 所以这可能真的很有用.

### 所需模块前缀

在命名常量时, 必须使用`snake_case`. 为了区分常量和局部变量, 必须指定consts的完整路径. 例如要访问PI常量, 必须在`math`之外使用完整的`math.pi`名称. 这一限制仅对`main`模块放宽(包含你的`fn main()`), 在这里你可以使用非限定名称的常量, 即`numbers`, 而不是`main.numbers`.

vfmt会处理这个规则, 所以你可以在`math`模块中使用`println(pi)`, 而vfmt会自动更新为`println(math.pi)`.

<!--
Many people prefer all caps consts: `TOP_CITIES`. This wouldn't work
well in V, because consts are a lot more powerful than in other languages.
They can represent complex structures, and this is used quite often since there
are no globals:

```v oksyntax
println('Top cities: ${top_cities.filter(.usa)}')
```
-->

## 内置函数

有些函数是内置的, 如`println`. 以下是完整的列表:

```v ignore
fn print(s string) // print anything on sdtout
fn println(s string) // print anything and a newline on sdtout

fn eprint(s string) // same as print(), but use stderr
fn eprintln(s string) // same as println(), but use stderr

fn exit(code int) // terminate the program with a custom error code
fn panic(s string) // print a message and backtraces on stderr, and terminate the program with error code 1
fn print_backtrace() // print backtraces on stderr
```

`println`是一个简单而强大的内置函数, 它可以打印任何东西: string, number, array, map, struct.

```v
struct User {
	name string
	age  int
}

println(1) // "1"
println('hi') // "hi"
println([1, 2, 3]) // "[1, 2, 3]"
println(User{ name: 'Bob', age: 20 }) // "User{name:'Bob', age:20}"
```

<a id='custom-print-of-types' />

## 打印自定义类型

如果你想为你的类型定义一个自定义打印值, 只需定义一个简单的`.str() string`方法即可:

```v
struct Color {
	r int
	g int
	b int
}

pub fn (c Color) str() string {
	return '{$c.r, $c.g, $c.b}'
}

red := Color{
	r: 255
	g: 0
	b: 0
}
println(red)
```

## 模块

文件夹目录下的每个文件都是同一个模块的一部分. 简单的程序不需要指定模块名, 在这种情况下, 它默认为'main'.

V是一种非常模块化的语言. 我们鼓励创建可重复使用的模块, 并且很容易做到. 要创建一个新的模块, 可用模块名称创建一个目录, 其中包含了
带代码的.v文件即可:

```shell
cd ~/code/modules
mkdir mymodule
vim mymodule/myfile.v
```
```v failcompile
// myfile.v
module mymodule

// To export a function we have to use `pub`
pub fn say_hi() {
    println('hello from mymodule!')
}
```

现在你可用在代码中使用`mymodule`了:

```v failcompile
import mymodule

fn main() {
    mymodule.say_hi()
}
```

* 模块名称应简短, 10个字符以下。
* 模块名称必须使用`snake_case`
* 不允许循环导入
* 一个模块中可以有任意多的.v文件
* 你可以在任何地方创建模块
* 所有的模块都可被静态地编译成一个可执行文件

### `init` 函数

如果你想让一个模块在导入时自动调用一些设置/初始化代码, 那么你可以使用模块的`init`函数.

```v
fn init() {
	// your setup code here ...
}
```

`init`函数不能是公开的 - 它将被自动调用. 这一特点对初始化C库特别有用.

## Types 2

### Interfaces

```v
struct Dog {
	breed string
}

struct Cat {
	breed string
}

fn (d Dog) speak() string {
	return 'woof'
}

fn (c Cat) speak() string {
	return 'meow'
}

// unlike Go and like TypeScript, V's interfaces can define fields, not just methods.
interface Speaker {
	breed string
	speak() string
}

dog := Dog{'Leonberger'}
cat := Cat{'Siamese'}

mut arr := []Speaker{}
arr << dog
arr << cat
for item in arr {
	println('a $item.breed says: $item.speak()')
}
```

一个类型通过实现其方法和字段来实现一个接口. 不需要明确的意图声明, 没有"implements"关键字.

#### 接口断言

我们可以使用`is`操作符来测试一个接口的底层类型:
```v oksyntax
interface Something {}

fn announce(s Something) {
	if s is Dog {
		println('a $s.breed dog') // `s` is automatically cast to `Dog` (smart cast)
	} else if s is Cat {
		println('a $s.breed cat')
	} else {
		println('something else')
	}
}
```
更多信息参考 [动态断言](#动态断言).

#### 接口方法的定义

同样与Go不同, 一个接口可以实现一个方法. 这些方法不是由实现了该接口的struct实现的.

当一个struct被封装在一个已经实现了该方法的接口中时. 与这个struct所实现的名称相同，此时只有接口上的该方法会被调用.

```v
struct Cat {}

fn (c Cat) speak() string {
	return 'meow!'
}

interface Adoptable {}

fn (a Adoptable) speak() string {
	return 'adopt me!'
}

fn new_adoptable() Adoptable {
	return Cat{}
}

fn main() {
	cat := Cat{}
	assert cat.speak() == 'meow!'
	a := new_adoptable()
	assert a.speak() == 'adopt me!'
	if a is Cat {
		println(a.speak()) // meow!
	}
}
```

### Enum

```v
enum Color {
	red
	green
	blue
}

mut color := Color.red
// V knows that `color` is a `Color`. No need to use `color = Color.green` here.
color = .green
println(color) // "green"
match color {
	.red { println('the color was red') }
	.green { println('the color was green') }
	.blue { println('the color was blue') }
}
```

枚举匹配必须是详尽的, 或者有一个`else`分支. 这确保了如果增加了一个新的枚举字段, 它在代码中的所有地方都会被处理.

枚举字段不能使用保留关键字. 然而, 保留的关键字可以用`@`转义而是被使用.

```v
enum Color {
	@none
	red
	green
	blue
}

color := Color.@none
println(color)
```

Integers may be assigned to enum fields.

```v
enum Grocery {
	apple
	orange = 5
	pear
}

g1 := int(Grocery.apple)
g2 := int(Grocery.orange)
g3 := int(Grocery.pear)
println('Grocery IDs: $g1, $g2, $g3')
```

Output: `Grocery IDs: 0, 5, 6`.

不允许在枚举变量上进行操作, 它们必须被明确地转换为`int`.

### 和类型

一个和类型的实例可以容纳几个不同类型的值. 使用`type`关键字来声明一个和类型:

```v
struct Moon {}

struct Mars {}

struct Venus {}

type World = Mars | Moon | Venus

sum := World(Moon{})
assert sum.type_name() == 'Moon'
println(sum)
```

内置方法`type_name`返回当前持有的类型.

使用和类型, 你可以建立递归结构, 并在其上写出简洁但强大的代码.
```v
// V's binary tree
struct Empty {}

struct Node {
	value f64
	left  Tree
	right Tree
}

type Tree = Empty | Node

// sum up all node values
fn sum(tree Tree) f64 {
	return match tree {
		Empty { f64(0) } // TODO: as match gets smarter just remove f64()
		Node { tree.value + sum(tree.left) + sum(tree.right) }
	}
}

fn main() {
	left := Node{0.2, Empty{}, Empty{}}
	right := Node{0.3, Empty{}, Node{0.4, Empty{}, Empty{}}}
	tree := Node{0.5, left, right}
	println(sum(tree)) // 0.2 + 0.3 + 0.4 + 0.5 = 1.4
}
```

#### 动态断言

要检查一个和类型实例是否拥有某个类型，使用`sum is Type`.
要将一个和类型转换为它的一个变体, 使用`sum as Type`.

```v
struct Moon {}

struct Mars {}

struct Venus {}

type World = Mars | Moon | Venus

fn (m Mars) dust_storm() bool {
	return true
}

fn main() {
	mut w := World(Moon{})
	assert w is Moon
	w = Mars{}
	// use `as` to access the Mars instance
	mars := w as Mars
	if mars.dust_storm() {
		println('bad weather!')
	}
}
```

如果`w`没有持有`Mars`实例, `as`会panic. 更安全的方法是使用智能转换.

#### 智能转换

```v oksyntax
if w is Mars {
	assert typeof(w).name == 'Mars'
	if w.dust_storm() {
		println('bad weather!')
	}
}
```
if`语句里的`w`内有`Mars`类型. 这就是所谓的 *流敏感类型*.
因为`w`是一个可变的标识符, 如果编译器智能转换它而不发出警告, 那将是不安全的.
这就是为什么你必须在`is`表达式之前声明一个`mut`.

```v ignore
if mut w is Mars {
	assert typeof(w).name == 'Mars'
	if w.dust_storm() {
		println('bad weather!')
	}
}
```
否则`w`将保持其基础类型.
> 这既适用于简单的变量, 也适用于复杂的表达式, 如`user.name`.

#### 匹配和类型

你也可以使用 `match` 来确定变种:

```v
struct Moon {}

struct Mars {}

struct Venus {}

type World = Mars | Moon | Venus

fn open_parachutes(n int) {
	println(n)
}

fn land(w World) {
	match w {
		Moon {} // no atmosphere
		Mars {
			// light atmosphere
			open_parachutes(3)
		}
		Venus {
			// heavy atmosphere
			open_parachutes(1)
		}
	}
}
```

`match`必须为每一个变体提供一个模式或有一个`else`分支.

```v ignore
struct Moon {}
struct Mars {}
struct Venus {}

type World = Moon | Mars | Venus

fn (m Moon) moon_walk() {}
fn (m Mars) shiver() {}
fn (v Venus) sweat() {}

fn pass_time(w World) {
    match w {
        // using the shadowed match variable, in this case `w` (smart cast)
        Moon { w.moon_walk() }
        Mars { w.shiver() }
        else {}
    }
}
```

### Option/Result 类型 和 error处理

Option类型用`?Type`定义:
```v
struct User {
	id   int
	name string
}

struct Repo {
	users []User
}

fn (r Repo) find_user_by_id(id int) ?User {
	for user in r.users {
		if user.id == id {
			// V automatically wraps this into an option type
			return user
		}
	}
	return error('User $id not found')
}

fn main() {
	repo := Repo{
		users: [User{1, 'Andrew'}, User{2, 'Bob'}, User{10, 'Charles'}]
	}
	user := repo.find_user_by_id(10) or { // Option types must be handled by `or` blocks
		return
	}
	println(user.id) // "10"
	println(user.name) // "Charles"
}
```

V将`Option`和`Result`合并为一种类型, 因此你无需决定要使用哪种类型.

将一个函数"升级"为可选函数所需的工作量是很小的. 你必须在返回类型中添加`？`, 并在出现问题时返回错误.

如果你不需要返回错误消息, 则只需`return none`(这与`return error("")`等效).

这是V中错误处理的主要机制. 它们仍然是值, 类似于Go, 但优点是error必须处理，而且处理起来也不冗长. 与其他语言不同, V不用`throw/try/catch`块来处理异常.

err是在or代码块中定义的, 并设置为传递字符串消息
给`error()`函数. 如果返回了`none`则`err`为空.

```v oksyntax
user := repo.find_user_by_id(7) or {
	println(err) // "User 7 not found"
	return
}
```

### 处理可选

有四种方法可以处理一个可选. 第一种方法是传播错误:

```v
import net.http

fn f(url string) ?string {
	resp := http.get(url) ?
	return resp.text
}
```

`http.get`返回`?http.Response`. 因为`?`跟在调用后面, 所以导致了error将被传播给`f`的调用者. 当使用`?`后, 函数调用产生一个可选函数, 则外层函数必须返回也是一个可选项. 如果在`main()`中使用了错误传播, 它将`panic`, 因为此时错误不能进一步传播.

`f`的body基本上是以下内容的浓缩版:

```v ignore
    resp := http.get(url) or { return err }
    return resp.text
```

---
第二种方法是提前脱离执行:

```v oksyntax
user := repo.find_user_by_id(7) or { return }
```

在这里，你可以调用`panic()`或`exit()`, 这样整个程序就会停止执行. 或使用控制流语句(`return`, `break`, `continue`等)跳出当前的代码块. 注意`break`和`continue`只能在`for`中使用.

V没有办法强制"unwrap"一个optional(像其他语言那样, 例如Rust的`unwrap()`或Swift的`!`). 要做到这一点, 可以使用`or { panic(err.msg) }`代替.

---
第三种方法是在`or`代码块的末尾提供一个默认值. 如果出现错误, 将以该值代替. 所以它必须与被处理的`Option`的内容具有相同的类型.

```v
fn do_something(s string) ?string {
	if s == 'foo' {
		return 'foo'
	}
	return error('invalid string') // Could be `return none` as well
}

a := do_something('foo') or { 'default' } // a will be 'foo'
b := do_something('bar') or { 'default' } // b will be 'default'
println(a)
println(b)
```

---
第四种方法是使用`if`拆包:

```v
import net.http

if resp := http.get('https://google.com') {
	println(resp.text) // resp is a http.Response, not an optional
} else {
	println(err)
}
```
上面`http.get`返回一个`?http.Response`. `resp`只在`if`分支作用域内, 而`err'只属于`else'分支的作用域.

## 泛型

```v wip

struct Repo<T> {
    db DB
}

struct User {
	id   int
	name string
}

struct Post {
	id   int
	user_id int
	title string
	body string
}

fn new_repo<T>(db DB) Repo<T> {
    return Repo<T>{db: db}
}

// This is a generic function. V will generate it for every type it's used with.
fn (r Repo<T>) find_by_id(id int) ?T {
    table_name := T.name // in this example getting the name of the type gives us the table name
    return r.db.query_one<T>('select * from $table_name where id = ?', id)
}

db := new_db()
users_repo := new_repo<User>(db) // returns Repo<User>
posts_repo := new_repo<Post>(db) // returns Repo<Post>
user := users_repo.find_by_id(1)? // find_by_id<User>
post := posts_repo.find_by_id(1)? // find_by_id<Post>
```

目前通用函数定义必须声明其类型参数, 但在未来V可以在运行时从单字母类型名推断出通用类型参数. 这就是为什么`find_by_id`可以省略`<T>`, 因为
接收器参数`r`使用通用类型`T`.

另一个例子:
```v
fn compare<T>(a T, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compare<int>
println(compare(1, 0)) // Outputs: 1
println(compare(1, 1)) //          0
println(compare(1, 2)) //         -1
// compare<string>
println(compare('1', '0')) // Outputs: 1
println(compare('1', '1')) //          0
println(compare('1', '2')) //         -1
// compare<f64>
println(compare(1.1, 1.0)) // Outputs: 1
println(compare(1.1, 1.1)) //          0
println(compare(1.1, 1.2)) //         -1
```


## 并发
### 生成并发任务
V的并发模型和Go的模型非常相似. 在V中要在其他线程并发执行`foo()`, 用`go foo()`即可:

```v
import math

fn p(a f64, b f64) { // ordinary function without return value
	c := math.sqrt(a * a + b * b)
	println(c)
}

fn main() {
	go p(3, 4)
	// p will be run in parallel thread
}
```

有时需要等待一个并行线程完成. 这可以通过给启动的线程分配一个*handle*, 并调用它的`wait()`方法来实现:

```v
import math

fn p(a f64, b f64) { // ordinary function without return value
	c := math.sqrt(a * a + b * b)
	println(c) // prints `5`
}

fn main() {
	h := go p(3, 4)
	// p() runs in parallel thread
	h.wait()
	// p() has definitely finished
}
```

这种方法也可以用来从一个在并行线程中运行的函数中获取返回值, 而并发调用时不需要修改函数本身.

```v
import math { sqrt }

fn get_hypot(a f64, b f64) f64 { //       ordinary function returning a value
	c := sqrt(a * a + b * b)
	return c
}

fn main() {
	g := go get_hypot(54.06, 2.08) // spawn thread and get handle to it
	h1 := get_hypot(2.32, 16.74) //   do some other calculation here
	h2 := g.wait() //                 get result from spawned thread
	println('Results: $h1, $h2') //   prints `Results: 16.9, 54.1`
}
```

如果有大量的任务, 使用线程数组来管理它们可能会更容易.

```v
import time

fn task(id int, duration int) {
	println('task $id begin')
	time.sleep(duration * time.millisecond)
	println('task $id end')
}

fn main() {
	mut threads := []thread{}
	threads << go task(1, 500)
	threads << go task(2, 900)
	threads << go task(3, 100)
	threads.wait()
	println('done')
}

// Output:
// task 1 begin
// task 2 begin
// task 3 begin
// task 3 end
// task 1 end
// task 2 end
// done
```

另外对于返回相同类型的线程, 在线程数组所在线程上调用`wait()`将返回所有计算值.

```v
fn expensive_computing(i int) int {
	return i * i
}

fn main() {
	mut threads := []thread int{}
	for i in 1 .. 10 {
		threads << go expensive_computing(i)
	}
	// Join all tasks
	r := threads.wait()
	println('All jobs finished: $r')
}

// Output: All jobs finished: [1, 4, 9, 16, 25, 36, 49, 64, 81]
```

### Channels
通道是coroutine之间的首选通信方式. V的通道的工作原理基本上就像Go. 你可以在一端将对象推入一个通道, 并从另一端弹出对象.
通道可以是缓冲的或无缓冲的, 并且可以用`select`对多个通道进行选择.

#### 语法和用法
通道的类型为`chan objtype`. 一个可选的缓冲区长度可以在声明中指定为`cap`属性:

```v
ch := chan int{} // unbuffered - "synchronous"
ch2 := chan f64{cap: 100} // buffer length 100
```

channel不必声明为`mut`. 缓冲区长度不是类型的一部分, 而是单个channel对象的一个属性. channel可以像普通的一样传递给coroutine变量:

```v
fn f(ch chan int) {
	// ...
}

fn main() {
	ch := chan int{}
	go f(ch)
	// ...
}
```

可以使用箭头操作符将对象推送到通道. 同样的操作符可以用来从另一端弹出对象:

```v
ch := chan int{}
ch2 := chan f64{}
n := 5
x := 7.3
ch <- n
// push
ch2 <- x
mut y := f64(0.0)
m := <-ch // pop creating new variable
y = <-ch2 // pop into existing variable
```

关闭channel表示不能再推入其他对象. 任何这样的尝试都会导致运行时的panic(除了`select`和`try_push()`--见下文). 如果相关的channel已经关闭并且缓冲区是空的, 那么弹出的尝试将立即返回. 这种情况可以使用or分支来处理(参见[处理选项](#处理选项)).

```v wip
ch := chan int{}
ch2 := chan f64{}
// ...
ch.close()
// ...
m := <-ch or {
    println('channel has been closed')
}

// propagate error
y := <-ch2 ?
```

#### Channel选择

`select`命令允许在没有明显的CPU负载的情况下同时监控几个通道. 它由一个可能的传输列表和相关的语句分支组成--类似于[match](#match)命令:
```v wip
import time
fn main () {
  c := chan f64{}
  ch := chan f64{}
  ch2 := chan f64{}
  ch3 := chan f64{}
  mut b := 0.0
  // ...
  select {
    a := <-ch {
        // do something with `a`
    }
    b = <-ch2 {
        // do something with predeclared variable `b`
    }
    ch3 <- c {
        // do something if `c` was sent
    }
    > 500 * time.millisecond {
        // do something if no channel has become ready within 0.5s
    }
  }
}
```
超时分支是可选的. 如果没有超时分支, 则`select`等待的时间不受限制. 如果在调用`select`时没有通道准备好, 也可以立即进行, 增加一个`else { .... }`分支即可. `else`和`> timeout`是排斥的.

`select`命令可以作为`bool`类型的*表达式*使用, 如果所有通道都关闭, 则会变为false:
```v wip
if select {
    ch <- a {
        // ...
    }
} {
    // channel was open
} else {
    // channel is closed
}
```

#### 特殊的Channel功能

对于特殊用途, 有一些内置的属性和方法:
```v
struct Abc {
	x int
}

a := 2.13
ch := chan f64{}
res := ch.try_push(a) // try to perform `ch <- a`
println(res)
l := ch.len // number of elements in queue
c := ch.cap // maximum queue length
is_closed := ch.closed // bool flag - has `ch` been closed
println(l)
println(c)
mut b := Abc{}
ch2 := chan Abc{}
res2 := ch2.try_pop(b) // try to perform `b = <-ch2`
```

`try_push/pop()`方法将立即返回其中一个结果: `.success`, `.not_ready`或`.closed` - 取决于对象是否已被转移, 或为什么不这样做.

不建议在生产中使用这些方法和属性 - 基于它们的算法往往受制于竞赛条件, 特别是`.len`和 `.closed`不应被用来做决定.

使用`or`分支, 错误传播或`select`代替(见[语法和用法](#语法和用法) 和 上面的[通道选择](#通道选择)).

### 共享对象

数据可以通过共享变量在coroutine和调用线程之间共享.

这样的变量应该创建为`shared`, 并且也以这样的方式传递给coroutine.

底层的`struct`包含一个隐藏的*mutex*, 允许锁定并发访问: 使用`rlock`代表只读, 使用`lock`代表读/写访问.

```v
struct St {
mut:
	x int // data to shared
}

fn (shared b St) g() {
	lock b {
		// read/modify/write b.x
	}
}

fn main() {
	shared a := St{
		x: 10
	}
	go a.g()
	// ...
	rlock a {
		// read a.x
	}
}
```
共享变量必须是struct, array 或 map.

## 解析JSON

```v
import json

struct Foo {
	x int
}

struct User {
	name string
	age  int
	// Use the `skip` attribute to skip certain fields
	foo Foo [skip]
	// If the field name is different in JSON, it can be specified
	last_name string [json: lastName]
}

data := '{ "name": "Frodo", "lastName": "Baggins", "age": 25 }'
user := json.decode(User, data) or {
	eprintln('Failed to decode json')
	return
}
println(user.name)
println(user.last_name)
println(user.age)
// You can also decode JSON arrays:
sfoos := '[{"x":123},{"x":456}]'
foos := json.decode([]Foo, sfoos) ?
println(foos[0].x)
println(foos[1].x)
```

由于JSON的普遍性，V中直接内置了对它的支持.

`json.decode`函数有两个参数: 第一个是JSON值应该被解码成的类型, 第二个是包含JSON数据的字符串.

V会生成JSON编码和解码的代码, 没有使用运行时反射, 这可以有更好的性能.

## 测试

### Assert

```v
fn foo(mut v []int) {
	v[0] = 1
}

mut v := [20]
foo(mut v)
assert v[0] < 4
```
`assert`语句检查其表达式是否为`true`. 如果断言失败, 程序将被中止. 断言只能用于检测编程错误. 当一个assert失败后, 会被报告给*stderr*, 而且比较操作符(如`<`、`==`)两边的数值将尽可能地被打印出来. 这有利于轻松找到一个意外值. 断言语句可以在任何函数中使用.

### 测试文件

```v
// hello.v
module main

fn hello() string {
	return 'Hello world'
}

fn main() {
	println(hello())
}
```

```v failcompile
module main
// hello_test.v
fn test_hello() {
    assert hello() == 'Hello world'
}
```
使用`v hello_test.v`执行上述测试. 这将检查函数`hello`是否是产生正确的输出. V会执行文件中的所有测试函数.

* 所有的测试函数必须在测试文件中, 文件名必须以`_test.v`结尾
* 测试函数的名字必须以`test_`开头，以标记它们的执行
* 普通函数也可以在测试文件中定义，并应手动调用. 其它符号也可以在测试文件中定义, 例如类型.
* 有两种测试：外部测试和内部测试
* 内部测试必须*声明*它们的模块，就像就像来自同一模块的所有其他.v文件一样. 内部测试甚至可以调用同一模块中的私有函数.
* 外部测试必须*导入*需要测试的模块. 它们不能访问模块的私有函数/类型. 它们只能测试模块提供的外部/公共 API.

在上面的例子中, `test_hello`是一个内部测试, 它调用私有函数`hello()`, 因为`hello_test.v`有`module main`.
就像`hello.v`一样，即两者都是同一个模块的一部分, 还请注意因为`module main`和其他模块一样是一个常规模块, 所以内部测试可以用来测试主程序.v文件中的私有函数.

你也可以在测试文件中定义特殊的测试函数:
* `testsuite_begin`将在所有其他测试函数之前运行
* `testsuite_end`将在所有其他测试函数之后运行

#### 执行测试

要在单个测试文件中运行测试函数, 使用`v foo_test.v`.

要测试整个模块, 使用`v test mymodule`. 你也可以使用`v test .`来测试你当前文件夹(和子文件夹)内的所有内容. 你可以通过`-stats`选项来查看关于单个测试运行的更多细节.

## Memory management

V通过使用值类型和字符串缓冲区, 首先避免了做不必要的分配, 促进了简单的无抽象代码风格.

大多数对象(约90-100%)都被V的自动释放引擎释放了: 编译器在编译过程中插入了自动进行必要的free调用, 剩余的小部分
的对象是通过引用计数释放的.

开发者不需要在他们的代码中改变任何东西. "它只是工作", 就像在
Python, Go或Java, 除了没有繁重的GC跟踪所有对象, 也没有为每个对象提供昂贵的RC.

### 控制

你可以利用V的自动释放引擎, 在自定义数据类型上定义一个`free()`方法即可:

```v
struct MyType {}

[unsafe]
fn (data &MyType) free() {
	// ...
}
```

就像编译器使用C的`free()`释放C数据类型一样, 它会在每个变量的生存期末为数据类型静态插入`free()`调用.

对于愿意进行更多低级控制的开发人员, 可以使用`-manualfree`禁用自动释放, 或在要手动管理其内存的每个函数上添加一个`[manualfree]`来禁用自动释放, 请参阅[属性](#属性).

_注意：现在, 自动释放隐藏在`-autofree`标志的后面. 默认情况下, 它将在V 0.3中启用. 如果不使用autofree, 则V程序将会出现内存泄漏.

### 例子

```v
import strings

fn draw_text(s string, x int, y int) {
	// ...
}

fn draw_scene() {
	// ...
	name1 := 'abc'
	name2 := 'def ghi'
	draw_text('hello $name1', 10, 10)
	draw_text('hello $name2', 100, 10)
	draw_text(strings.repeat(`X`, 10000), 10, 50)
	// ...
}
```

字符串不会转义`draw_text`, 因此当函数退出时它们会被清除.

实际上, 使用`-prealloc`标志, 前两个调用根本不会导致任何分配. 这两个字符串很小, 因此V将为它们使用预分配的缓冲区.

```v
struct User {
	name string
}

fn test() []int {
	number := 7 // stack variable
	user := User{} // struct allocated on stack
	numbers := [1, 2, 3] // array allocated on heap, will be freed as the function exits
	println(number)
	println(user)
	println(numbers)
	numbers2 := [4, 5, 6] // array that's being returned, won't be freed here
	return numbers2
}
```

## ORM

(目前仍处于alpha状态)

V内置了ORM(对象关系映射), 支持SQLite, 并将很快支持MySQL, Postgres, MS SQL和Oracle.

V的ORM提供了许多好处:

- 所有的SQL方言都用一种语法 (在数据库之间的迁移变得更加容易)
- 使用V的语法构建查询 (不需要学习另一种语法)
- 安全性 (所有的查询都会自动处理, 以防止SQL注入)
- 编译时的检查 (这可以防止只有在运行时才能发现的排版错误)
- 可读性和简单性 (你不需要手动解析查询的结果, 也不需要从解析结果中手动构造对象)

```v
import sqlite

struct Customer {
	// struct name has to be the same as the table name (for now)
	id        int // a field named `id` of integer type must be the first field
	name      string
	nr_orders int
	country   string
}

db := sqlite.connect('customers.db') ?
// select count(*) from Customer
nr_customers := sql db {
	select count from Customer
}
println('number of all customers: $nr_customers')
// V syntax can be used to build queries
// db.select returns an array
uk_customers := sql db {
	select from Customer where country == 'uk' && nr_orders > 0
}
println(uk_customers.len)
for customer in uk_customers {
	println('$customer.id - $customer.name')
}
// by adding `limit 1` we tell V that there will be only one object
customer := sql db {
	select from Customer where id == 1 limit 1
}
println('$customer.id - $customer.name')
// insert a new customer
new_customer := Customer{
	name: 'Bob'
	nr_orders: 10
}
sql db {
	insert new_customer into Customer
}
```

更多例子见 <a href='https://github.com/vlang/v/blob/master/vlib/orm/orm_test.v'>vlib/orm/orm_test.v</a>.

## 文档编写

它的工作方式与Go非常相似. 很简单: 不需要为你的代码单独写文档, vdoc会根据源代码中的docstrings生成文档.

每个函数/类型/const的文档必须放在声明之前:

```v
// clearall clears all bits in the array
fn clearall() {
}
```

注释必须以定义的名称开始.

有时一行不足以解释一个函数的作用, 在这种情况下, 注释应该使用单行注释跨越到文档中的函数前:

```v
// copy_all recursively copies all elements of the array by their value,
// if `dupes` is false all duplicate values are eliminated in the process.
fn copy_all(dupes bool) {
	// ...
}
```

按照惯例, 最好用*现在时*写评论.

模块的概述必须放在模块名称之后的第一条评论中.

要生成文档，请使用vdoc, 例如`v doc net.http`.

## 工具

### v fmt

你不需要担心你的代码格式化或设置风格准则. `v fmt`会处理这些问题:

```shell
v fmt file.v
```

建议设置你的编辑器, 在每次保存时执行`v fmt -w`. vfmt执行成本通常很便宜(需要<30ms).

在推送代码之前, 一定要运行`v fmt -w file.v`.

### 剖析

V对程序剖析有很好的支持: `v -profile profile.txt run file.v`. 这将产生一个profile.txt文件, 你可以对其进行分析.

生成的profile.txt文件会有4列:
a) 一个函数的调用次数
b) 一个函数总共需要多少时间(毫秒)
c) 调用一个函数平均需要多少时间(纳秒)
d) v函数的名称

你可以使用以下方法对第3列(每个函数的平均时间)进行排序:
`sort -n -k3 profile.txt|tail`

你也可以使用秒表来精确地测量代码的一部分:
```v
import time

fn main() {
	sw := time.new_stopwatch({})
	println('Hello world')
	println('Greeting the world took: ${sw.elapsed().nanoseconds()}ns')
}
```

# 高级功能

## 非内存安全的代码
有时为了效率, 你可能会想写一些低级别的代码, 尽管这些代码可能会破坏内存或容易被安全漏洞利用. V支持编写这样的代码, 但不是默认的.

V要求有意标记任何潜在的不安全的内存操作. 标记这些操作也向任何阅读代码的人表明, 这里违反了内存安全, 可能导致错误.

潜在的内存不安全操作的例子有:
* 指针运算
* 指针索引
* 从不兼容类型转换为指针
* 调用某些C函数，如`free`, `strlen`和`strncmp`

要标记潜在的不安全内存操作, 请将其放在`unsafe`块中:

```v wip
// allocate 2 uninitialized bytes & return a reference to them
mut p := unsafe { malloc(2) }
p[0] = `h` // Error: pointer indexing is only allowed in `unsafe` blocks
unsafe {
    p[0] = `h` // OK
    p[1] = `i`
}
p++ // Error: pointer arithmetic is only allowed in `unsafe` blocks
unsafe {
    p++ // OK
}
assert *p == `i`
```

最好的做法是避免将内存安全表达式放在`unsafe` 块中, 以便尽可能明确使用`unsafe`的原因. 一般来说, 任何代码你认为是内存安全的就不应该放在一个`unsafe` 块中, 因为编译器可以验证它.

如果你怀疑你的程序违反了内存安全规定, 你就有了一个好的开端找出原因: 查看`unsafe`代码块(以及它们是如何与周边相互作用的).

* 注: 这是一项正在进行的工作.

### 带参考字段的struct

带有引用的struct需要明确地将初始值设置为一个引用值, 除非struct已经定义了自己的初始值.

零值引用或者nil指针, 未来将**不**支持. 目前依赖于可以使用值'0'的引用字段的数据结构Linked Lists或Binary Tree将被理解为不安全, 并且会引起panic.

```v
struct Node {
	a &Node
	b &Node = 0 // Auto-initialized to nil, use with caution!
}

// Reference fields must be initialized unless an initial value is declared.
// Zero (0) is OK but use with caution, it's a nil pointer.
foo := Node{
	a: 0
}
bar := Node{
	a: &foo
}
baz := Node{
	a: 0
	b: 0
}
qux := Node{
	a: &foo
	b: &bar
}
println(baz)
println(qux)
```

## sizeof 和 __offsetof

* `sizeof(Type)` 返回一个类型大小.
* `__offsetof(Struct, field_name)` 返回struct字段的偏移量

```v
struct Foo {
	a int
	b int
}

assert sizeof(Foo) == 8
assert __offsetof(Foo, a) == 0
assert __offsetof(Foo, b) == 4
```

## 在 V 中调用 C

### 例子

```v
#flag -lsqlite3
#include "sqlite3.h"
// See also the example from https://www.sqlite.org/quickstart.html
struct C.sqlite3 {
}

struct C.sqlite3_stmt {
}

type FnSqlite3Callback = fn (voidptr, int, &charptr, &charptr) int

fn C.sqlite3_open(charptr, &&C.sqlite3) int

fn C.sqlite3_close(&C.sqlite3) int

fn C.sqlite3_column_int(stmt &C.sqlite3_stmt, n int) int

// ... you can also just define the type of parameter and leave out the C. prefix
fn C.sqlite3_prepare_v2(&C.sqlite3, charptr, int, &&C.sqlite3_stmt, &charptr) int

fn C.sqlite3_step(&C.sqlite3_stmt)

fn C.sqlite3_finalize(&C.sqlite3_stmt)

fn C.sqlite3_exec(db &C.sqlite3, sql charptr, cb FnSqlite3Callback, cb_arg voidptr, emsg &charptr) int

fn C.sqlite3_free(voidptr)

fn my_callback(arg voidptr, howmany int, cvalues &charptr, cnames &charptr) int {
	unsafe {
		for i in 0 .. howmany {
			print('| ${cstring_to_vstring(cnames[i])}: ${cstring_to_vstring(cvalues[i]):20} ')
		}
	}
	println('|')
	return 0
}

fn main() {
	db := &C.sqlite3(0) // this means `sqlite3* db = 0`
	// passing a string literal to a C function call results in a C string, not a V string
	C.sqlite3_open('users.db', &db)
	// C.sqlite3_open(db_path.str, &db)
	query := 'select count(*) from users'
	stmt := &C.sqlite3_stmt(0)
	// NB: you can also use the `.str` field of a V string,
	// to get its C style zero terminated representation
	C.sqlite3_prepare_v2(db, query.str, -1, &stmt, 0)
	C.sqlite3_step(stmt)
	nr_users := C.sqlite3_column_int(stmt, 0)
	C.sqlite3_finalize(stmt)
	println('There are $nr_users users in the database.')
	//
	error_msg := charptr(0)
	query_all_users := 'select * from users'
	rc := C.sqlite3_exec(db, query_all_users.str, my_callback, 7, &error_msg)
	if rc != C.SQLITE_OK {
		eprintln(cstring_to_vstring(error_msg))
		C.sqlite3_free(error_msg)
	}
	C.sqlite3_close(db)
}
```

### 传递C编译参数

在V文件的顶部添加`＃flag`指令, 以提供C编译标志, 例如：

-`-I'用于添加C包含文件搜索路径
-`-l`用于添加要链接的C库名称
-`-L`用于添加C库文件的搜索路径
-`-D`设置编译时间变量

你可以(可选)对不同的目标使用不同的标志. 当前支持`linux`, `darwin`, `freebsd`和`windows`标志.

注意：每个标志必须用自己的行（暂时） 

```v oksyntax
#flag linux -lsdl2
#flag linux -Ivig
#flag linux -DCIMGUI_DEFINE_ENUMS_AND_STRUCTS=1
#flag linux -DIMGUI_DISABLE_OBSOLETE_FUNCTIONS=1
#flag linux -DIMGUI_IMPL_API=
```

在控制台构建命令中, 你可以使用:
* `-cflags`传递自定义标志给后端C语言编译器
* `-cc`来改变默认的C语言后端编译器
* 例如: `-cc gcc-9 -cflags -fsanitize=thread`

你可以在你的终端中定义一个 `VFLAGS` 环境变量来存储你的`-cc`和`-cflags`设置, 而不用每次在编译命令中包含它们.

### #pkgconfig

添加`#pkgconfig`指令是用来告诉编译器应该使用哪些模块来编译, 并使用各自依赖提供的pkg-config文件进行链接.

在`#flag`中不能使用backticks，而且出于安全和可移植性的考虑, 不希望产生进程，V使用自己的pkgconfig库, 它与标准的freedesktop库兼容.

如果没有传递flags, 它会添加`--cflags`和`--libs`, 下面两行都是一样的:

```v oksyntax
#pkgconfig r_core
#pkgconfig --cflags --libs r_core
```

pkgconfig会从一个硬编码的默认 pkg-config 路径列表查找`.pc`文件. 用户可以通过使用`PKG_CONFIG_PATH`环境变量来添加额外的路径. pkgconfig可以传入多个模块.

### 包含C代码

你也可以直接在你的V模块中包含C代码.
例如, 假设你的C代码位于你的模块文件夹中的一个名为'c'的文件夹中. 那么, 我们可以这样做:

* 在top level文件夹中放一个v.mod 文件(如果你使用`v new`创建模块, 则你已经有了v.mod文件), 例如:
```v ignore
Module {
	name: 'mymodule',
	description: 'My nice module wraps a simple C library.',
	version: '0.0.1'
	dependencies: []
}
```


* 在module的顶部添加以下几行:
```v oksyntax
#flag -I @VROOT/c
#flag @VROOT/c/implementation.o
#include "header.h"
```
NB: @VROOT 将被 V 替换为 *最接近的父文件夹，那里有 v.mod 文件*. v.mod文件所在文件夹同级或下面的任何`.v`文件均可以使用`#flag @VROOT/abc`来引用这个文件夹. @VROOT文件夹也是模块查找路径的*前缀*, 所以你可以在你的@VROOT下*导入*其他模块, 只需给它们命名即可.

上面的说明将使V在你的@VROOT中寻找一个编译过的`.o`文件, 路径是你的模块`folder/c/implementation.o`. 
如果V找到它, `.o`文件将被链接到使用该模块的主可执行文件. 如果没有找到, V会认为有一个`@VROOT/c/implementation.c`文件, 先尝试将其编译成.o文件, 然后再使用该文件.

这允许你把C代码包含在V模块中, 这样它的发布就更容易了. 你可以在这里看到一个完整的在V包装模块中使用C代码的最小例子[project_with_c_code](https://github.com/vlang/v/tree/master/vlib/v/tests/project_with_c_code). 另一个例子, 演示了将struct从C到V再从V传递回去
[interoperate between C to V to C](https://github.com/vlang/v/tree/master/vlib/v/tests/project_with_c_code_2).

### C类型
普通的c字符串可用`unsafe { charptr(cstring).vstring() }`转成V string, 如果已知道其长度还可用`unsafe { charptr(cstring).vstring_with_len(len) }`.

NB: .vstring()和.vstring_with_len()方法不创建`cstring`的副本. 所以你不应该在调用`.vstring()`方法后释放它. 如果你需要复制一个C语言字符串(一些libc API, 比如`getenv`就需要这样做, 因为它们返回指向内部libc内存的指针), 你可以使用`cstring_to_vstring(cstring)`.

在Windows上，C API经常返回所谓的`wide`字符串(utf16编码). 这些字符串可以用`string_from_wide(&u16(cwidestring))`转换为V字符串.

V 字符串有这些类型是为了更容易与 C 字符串互操作:

- `voidptr` for C's `void*`,
- `byteptr` for C's `byte*` and
- `charptr` for C's `char*`.
- `&charptr` for C's `char**`

要将`voidptr`转为V引用, 可使用`user := &User(user_void_ptr)`.

`voidptr`也可以解引用到V struct中: `user := User(user_void_ptr)`.

[一个在V中调用C代码的例子](https://github.com/vlang/v/blob/master/vlib/v/tests/project_with_c_code/mod1/wrapper.v)

### C 声明

C标识符的访问方式与特定模块标识符的访问方式类似，使用`C`前缀. 函数在使用前必须在V中重新声明. 任何C类型都可以在`C`前缀后面使用, 但类型必须在V中重新声明才能访问类型成员.

要重新声明复杂的类型, 例如在下面的C代码中:

```c
struct SomeCStruct {
	uint8_t implTraits;
	uint16_t memPoolData;
	union {
		struct {
			void* data;
			size_t size;
		};

		DataView view;
	};
};
```

子数据结构的成员可以直接在包含的struct中声明, 如下所示:

```v
struct C.SomeCStruct {
	implTraits  byte
	memPoolData u16
	// These members are part of sub data structures that can't currently be represented in V.
	// Declaring them directly like this is sufficient for access.
	// union {
	// struct {
	data voidptr
	size size_t
	// }
	view C.DataView
	// }
}
```

V知道了数据成员的存在, 可以不完全重新创建原始结构而使用它们.

另外, 可以将子数据结构[嵌入](#嵌入struct), 以保持一个并行的代码结构.

## 调试生成的C代码

为了调试生成的C代码中的问题，你可以传递这些标志:

- `-g` : 产生一个优化程度较低的可执行文件, 其中包含更多的调试信息.

    V将强制执行堆栈跟踪中.v文件的行号，以便可执行文件panic时可看到. 通常最好是通过`-g, 除非是你正在编写低级代码, 在这种情况下, 使用下一个选项`-cg`
- `-cg` : 产生一个优化程度较低的可执行文件，其中包含更多的调试信息.


    在这种情况下, 可执行文件将使用C源代码行号, 它经常与`-keepc`结合使用, 这样你就可以检查生成的或者让你的调试器(`gdb`, `lldb`等) 可以向你展示生成的C源代码
- `-showcc` : 打印用于构建程序的C命令
- `-show-c-output` : 打印编译程序时C编译器产生的输出
- `-keepc` : 在编译成功后不删除生成的C源代码文件. 同时继续使用相同的文件路径, 这样更稳定. 并且更容易在编辑器/IDE中使用.

为了获得最佳的调试体验, 如果你正在封装一个现有的C库，你可以同时传递几个这样的标志: `v -keepc -cg -showcc yourprogram.v`, 然后运行你的调试器(gdb/lldb)或IDE.

如果你只是想检查生成的C代码，那么无需进一步编译, 你也可以使用`-o`标志(例如`-o file.c`), 这将使V产生`file.c`后停止.

如果你仅想看C源代码生成的某个C函数，例如`main`, 你可以使用`-o`标志(例如`-o file.c`), 即`-printfn main -o file.c`.

要查看V支持的所有标志的详细列表，可使用`v help`, `v help build`和`v help build-c`.

## 有条件的汇编

### 编译时代码

`$` 被用作编译时操作的前缀.

#### $if
```v
// Support for multiple conditions in one branch
$if ios || android {
	println('Running on a mobile device!')
}
$if linux && x64 {
	println('64-bit Linux.')
}
// Usage as expression
os := $if windows { 'Windows' } $else { 'UNIX' }
println('Using $os')
// $else-$if branches
$if tinyc {
	println('tinyc')
} $else $if clang {
	println('clang')
} $else $if gcc {
	println('gcc')
} $else {
	println('different compiler')
}
$if test {
	println('testing')
}
// v -cg ...
$if debug {
	println('debugging')
}
// v -prod ...
$if prod {
	println('production build')
}
// v -d option ...
$if option ? {
	println('custom option')
}
```

如果你想让一个`if`在编译时被评估, 它必须在前面加上`$`符号. 现在它可以用来检测操作系统, 编译器, 平台或编译选项. `$if debug`是一个特殊的选项, 像`$if windows`或`$if x32`, 如果你使用的是自定义的ifdef, 那么你确实需要`$if选项? {}`并使用`v -d option`编译.

完整的内置选项列表:
| OS                            | Compilers         | Platforms             | Other                     |
| ---                           | ---               | ---                   | ---                       |
| `windows`, `linux`, `macos`   | `gcc`, `tinyc`    | `amd64`, `aarch64`    | `debug`, `prod`, `test`   |
| `mac`, `darwin`, `ios`,       | `clang`, `mingw`  | `x64`, `x32`          | `js`, `glibc`, `prealloc` |
| `android`,`mach`, `dragonfly` | `msvc`            | `little_endian`       | `no_bounds_checking`      |
| `gnu`, `hpux`, `haiku`, `qnx` | `cplusplus`       | `big_endian`          | |
| `solaris`, `linux_or_macos`   | | | |

#### $embed_file

```v ignore
import os
fn main() {
	embedded_file := $embed_file('v.png')
	os.write_file('exported.png', embedded_file.to_string()) ?
}
```

V可以用`$embed_file(<path>)`将任意文件嵌入到可执行文件中, 是编译时调用的, 路径可以是源文件的绝对路径或相对路径.

当你不使用`-prod`时, 文件将不会被嵌入. 相反，它将在你的程序运行时第一次调用`f.data()`时被加载, 使得更容易在外部编辑程序中进行修改, 而不需要重新编译可执行文件.

当你用`-prod`编译时, 该文件*会被嵌入*到可执行文件，增加二进制文件的大小, 但使它更加自由, 从而更容易分发. 在这种情况下，`f.data()`没有io操作, 且始终返回相同的数据.

#### $tmpl: 内嵌和解析V template文件

V有一个简单的模板语言, 用于文本和html模板, 它们可以通过`$tmpl('path/to/template.txt')`轻松嵌入:

```v ignore
fn build() string {
	name := 'Peter'
	age := 25
	numbers := [1, 2, 3]
	return $tmpl('1.txt')
}

fn main() {
	println(build())
}
```

1.txt:

```
name: @name

age: @age

numbers: @numbers

@for number in numbers
  @number
@end
```

output:

```
name: Peter

age: 25

numbers: [1, 2, 3]

1
2
3
```




#### $env

```v
module main

fn main() {
	compile_time_env := $env('ENV_VAR')
	println(compile_time_env)
}
```

V可以在编译时从环境变量中引入值. `$env('ENV_VAR')`也可以用在顶级的`#flag`和`#include`语句中: `#flag linux -I $env('JAVA_HOME')/include`.

### 特定环境文件

如果一个文件有一个特定环境的后缀, 它将只针对该环境进行编译.

- `.js.v` => 将仅由JS后端使用. 这些文件可以包含JS代码
- `.c.v` => 仅供C后端使用, 这些文件可以包含C代码
- `.x64.v` => 仅供V的x64后端使用
- `_nix.c.v` => 仅供Unix系统(非Windows)使用
- `_${os}.c.v` => 将只在特定的`os`系统上使用.

    例如, `_windows.c.v`将只在Windows系统上编译时使用, 或者使用`-os windows`
- `_default.c.v` => 只有在没有更具体的平台文件时才会使用.

    例如，如果你同时有`file_linux.c.v`和`file_default.c.v, 而你是为linux编译的, 那么只需要使用`file_linux.c.v`, 而`file_default.c.v`将被忽略

这里有一个更完整的例子:
main.v:
```v ignore
module main
fn main() { println(message) }
```

main_default.c.v:
```v ignore
module main
const ( message = 'Hello world' )
```

main_linux.c.v:
```v ignore
module main
const ( message = 'Hello linux' )
```

main_windows.c.v:
```v ignore
module main
const ( message = 'Hello windows' )
```

使用上面的例子:
- 当你为windows编译时, 你会得到'Hello windows'
- 当你为linux编译时，你会得到'Hello linux'
- 当你在其他平台上编译时, 你会得到的是非特定的'Hello world'信息
- `_d_customflag.v` => 只有当你把`-d customflag`传给V时才会使用

    这相当于`$if customflag ? {}`, 但是对于整个文件, 而不仅仅是一个代码块. `customflag`应该是一个snake_case标识符, 不能包含任意字符(只能是小写拉丁字母+数字+`_`).

    NB: 组合式的`_d_customflag_linux.c.v` postfix将无法工作. 如果你确实需要一个自定义的flag文件, 其中有依赖于平台的代码, 请使用`_d_customflag.v`, 然后在里面使用plaftorm依赖编译时的条件块, 即`$if linux {}`等.
- `_notd_customflag.v` => 类似于_d_customflag.v, 但会被用于 *只有当你不向V传递`-d customflag`时才会出现

## 编译时伪变量

V还允许你的代码访问一组伪字符串变量, 其在编译时会被替换.

- `@FN` => 用当前V函数的名称代替
- `@METHOD` => 用ReceiverType.MethodName替换
- `@MOD` => 用当前V模块的名称代替
- `@STRUCT` => 用当前V结构的名称代替
- `@FILE` => 用V源文件的路径替换
- `@LINE` => 替换为出现的V行号(作为字符串)
- `@COLUMN` => 替换为出现的列号(作为一个字符串)
- `@VEXE` => 替换为V编译器的路径
- `@VHASH` => 替换为V编译器的缩短提交哈希值(作为一个字符串)
- `@VMOD_FILE` => 用最近的v.mod文件的内容(作为一个字符串)替换

这允许你做下面的例子, 在调试/记录/跟踪你的代码时很有用:
```v
eprintln('file: ' + @FILE + ' | line: ' + @LINE + ' | fn: ' + @MOD + '.' + @FN)
```

另一个例子是，如果你想把v.mod的版本/名称嵌入到你的可执行文件中:
```v ignore
import v.vmod
vm := vmod.decode( @VMOD_FILE ) or { panic(err.msg) }
eprintln('$vm.name $vm.version\n $vm.description')
```

## 性能调整

当编译代码时使用`-prod`, 生成的C代码一般都是够快的. 但在某些情况下, 你可能添加额外的提示给编译器, 这样它就可以进一步优化某些代码块.

NB: 这些是很*少*用到, 甚至不应该使用, 除非 *剖析代码*后看到它们有显著的好处. 引用gcc的文档: "程序员在预测他们的程序是如何实际执行的方面是出了名的差劲".

`[inline]` - 你可以用`[inline]`标记函数, 这样C编译器会尽量将其内嵌, 在某些情况下, 可能对性能有利. 但可能会影响你的可执行文件的大小.

`[direct_array_access]` - 在带有`[direct_array_access]`标记的函数中, 编译器会将数组操作直接翻译成C数组操作(省略边界检查), 这可能会在迭代一个数组的函数中节省很多时间, 但代价是使函数不安全, 除非由用户自己检查边界.

`if _likely_(bool表达式) {`这暗示了C编译器布尔表达式大概率为真, 所以它生成的汇编代码, 分支误判的几率较小. 在JS后端, 什么都不做.

`if _unlikely_(bool表达式) {`类似于`_likely_(x)`, 但它暗示了的布尔表达式是非常不可能执行到的. 在JS后端, 这什么也做不了.

<a id='Reflection via codegen'>

## 编译时反思

拥有内置的JSON支持是不错的, 但V也允许你为任何数据格式创建高效的序列化器. V有编译时的`if`和`for`结构:

```v wip
// TODO: not fully implemented

struct User {
    name string
    age  int
}

// Note: T should be passed a struct name only
fn decode<T>(data string) T {
    mut result := T{}
    // compile-time `for` loop
    // T.fields gives an array of a field metadata type
    $for field in T.fields {
        $if field.typ is string {
            // $(string_expr) produces an identifier
            result.$(field.name) = get_string(data, field.name)
        } $else $if field.typ is int {
            result.$(field.name) = get_int(data, field.name)
        }
    }
    return result
}

// `decode<User>` generates:
fn decode_User(data string) User {
    mut result := User{}
    result.name = get_string(data, 'name')
    result.age = get_int(data, 'age')
    return result
}
```

## 有限的操作符重载

```v
struct Vec {
	x int
	y int
}

fn (a Vec) str() string {
	return '{$a.x, $a.y}'
}

fn (a Vec) + (b Vec) Vec {
	return Vec{a.x + b.x, a.y + b.y}
}

fn (a Vec) - (b Vec) Vec {
	return Vec{a.x - b.x, a.y - b.y}
}

fn main() {
	a := Vec{2, 3}
	b := Vec{4, 5}
	mut c := Vec{1, 2}
	println(a + b) // "{6, 8}"
	println(a - b) // "{-2, -2}"
	c += a
	println(c) // "{3, 5}"
}
```

操作符超载违背了V的简单和可预测性的理念. 但由于科学和图形应用是V的领域之一. 为了提高可读性，运算符重载是一个重要的特征.

`a.add(b).add(c.mul(d))`比`a + b + c * d`可读性要差得多.

为了提高安全性和可维护性, 对操作符的重载进行了限制:

- 只能重载`+, -, *, /, %, <, >, ==, !=, <=, >=`运算符
- `==`和`!=`由编译器自行生成, 但可以重载
- 不允许在运算符函数里面调用其他函数
- 运算符函数不能修改其参数
- 当使用`<`和`==`运算符时, 返回类型必须是`bool`
- 当定义了`==`和`<`时，`！=`、`>`、`<=`和`>=`会自动生成
- 两个参数必须具有相同的类型(就像V中的所有操作符一样)
- 赋值运算符(`*=`, `+=`, `/=`等)在定义运算符时, 会自动生成, 但它们必须返回相同的类型.

## 内联编译
<!-- ignore because it doesn't pass fmt test (why?) --> 
```v ignore
a := 100
b := 20
mut c := 0
asm amd64 {
    mov eax, a
    add eax, b
    mov c, eax
    ; =r (c) as c // output 
    ; r (a) as a // input 
      r (b) as b
}
println('a: $a') // 100 
println('b: $b') // 20 
println('c: $c') // 120
```

更多例子可见 [github.com/vlang/v/tree/master/vlib/v/tests/assembly/asm_test.amd64.v](https://github.com/vlang/v/tree/master/vlib/v/tests/assembly/asm_test.amd64.v)

## 将C翻译成V

TODO：在V 0.3中可以将C语言翻译成V语言.

V可以将你的C代码翻译成人类可读的V代码，并在C库之上生成V包装器.

我们先创建一个简单的程序`test.c`:

```c
#include "stdio.h"

int main() {
	for (int i = 0; i < 10; i++) {
		printf("hello world\n");
	}
        return 0;
}
```

执行`v translate test.c`, 然后V会生成 `test.v`:

```v
fn main() {
	for i := 0; i < 10; i++ {
		println('hello world')
	}
}
```

要在C库的基础上生成一个封装器, 请使用以下命令:

```bash
v wrapper c_code/libsodium/src/libsodium
```

这将生成一个带有V模块的目录`libsodium`.

C2V生成的libsodium包装器的例子, 见https://github.com/medvednikov/libsodium

<br>

什么时候应该翻译C代码, 什么时候应该简单地从V中调用C代码?

如果你的C代码写得很好，经过很好的测试, 那么当然你可以一直简单地从V中调用这个C代码.

将它翻译成V给你带来了几个好处:

- 如果你打算开发那个代码库, 你现在已经在一种语言中得到了所有的东西, 比C语言更安全，更容易开发
- 交叉编译变得更加容易. 你根本不用担心这个问题.
- 也没有更多的构建标志和包含文件

## 热更新

```v live
module main

import time
import os

[live]
fn print_message() {
	println('Hello! Modify this message while the program is running.')
}

fn main() {
	for {
		print_message()
		time.sleep(500 * time.millisecond)
	}
}
```

用`v -live message.v`构建这个例子.

你想重载的函数在定义前必须有`[live]`属性.

目前, 在程序运行时还不能修改类型.

更多的例子(包括一个图形应用程序)可见
[github.com/vlang/v/tree/master/examples/hot_code_reload](https://github.com/vlang/v/tree/master/examples/hot_reload)

## 交叉汇编

要交叉编译你的项目, 只需运行:

```shell
v -os windows .
```

或

```shell
v -os linux .
```

(macOS的交叉编译暂时无法实现。)

如果你没有任何C语言的依赖, 那就只需要这样做. 这甚至在使用`ui`模块编译GUI应用程序或使用`gg`编译图形应用程序时也可这样.

你需要安装Clang, LLD链接器, 并根据V提供的一个链接下载一个包含以下内容的zip文件, 它包含了为Windows和Linux提供支持的lib和include文件.

## V中的跨平台shell脚本

V可以作为Bash的替代品来编写部署脚本, 构建脚本等.

使用V的优势在于语言的简单性和可预测性, 以及跨平台支持. "V script"既可以在类似Unix的系统上运行, 也可以在Windows上运行.

使用`.vsh`文件扩展名, 并将`os`模块中的所有函数成为全局函数(这样你就可以使用`mkdir()`而不是`os.mkdir()`).

`deploy.vsh`例子:
```v wip
#!/usr/bin/env -S v run
// The shebang above associates the file to V on Unix-like systems,
// so it can be run just by specifying the path to the file
// once it's made executable using `chmod +x`.

// Remove if build/ exits, ignore any errors if it doesn't
rmdir_all('build') or { }

// Create build/, never fails as build/ does not exist
mkdir('build') ?

// Move *.v files to build/
result := exec('mv *.v build/') ?
if result.exit_code != 0 {
	println(result.output)
}
// Similar to:
// files := ls('.') ?
// mut count := 0
// if files.len > 0 {
//     for file in files {
//         if file.ends_with('.v') {
//              mv(file, 'build/') or {
//                  println('err: $err')
//                  return
//              }
//         }
//         count++
//     }
// }
// if count == 0 {
//     println('No files')
// }
```

现在你可以像编译普通的V程序一样编译这个程序, 然后得到一个可执行文件, 并可以在任何地方部署和运行:
`v deploy.vsh && ./deploy`

或者就像传统的Bash脚本一样运行它:
`v run deploy.vsh`

在类似Unix的平台上, 使用`chmod +x`使文件可执行后, 可直接运行:
`./deploy.vsh`

## 属性
V有几个属性可以修改函数和struct的行为.

属性是指在`[]`内指定的编译器指令，它位于function/struct/enum声明, 并且只适用于以下声明.

```v
// Calling this function will result in a deprecation warning
[deprecated]
fn old_function() {
}

// It can also display a custom deprecation message
[deprecated: 'use new_function() instead']
fn legacy_function() {}

// This function's calls will be inlined.
[inline]
fn inlined_function() {
}

// The following struct must be allocated on the heap. Therefore, it can only be used as a
// reference (`&Window`) or inside another reference (`&OuterStruct{ Window{...} }`).
[heap]
struct Window {
}

// V will not generate this function and all its calls if the provided flag is false.
// To use a flag, use `v -d flag`
[if debug]
fn foo() {
}

fn bar() {
	foo() // will not be called if `-d debug` is not passed
}

// Calls to following function must be in unsafe{} blocks.
// Note that the code in the body of `risky_business()` will still be
// checked, unless you also wrap it in `unsafe {}` blocks.
// This is usefull, when you want to have an `[unsafe]` function that
// has checks before/after a certain unsafe operation, that will still
// benefit from V's safety features.
[unsafe]
fn risky_business() {
	// code that will be checked, perhaps checking pre conditions
	unsafe {
		// code that *will not be* checked, like pointer arithmetic,
		// accessing union fields, calling other `[unsafe]` fns, etc...
		// Usually, it is a good idea to try minimizing code wrapped
		// in unsafe{} as much as possible.
		// See also [Memory-unsafe code](#memory-unsafe-code)
	}
	// code that will be checked, perhaps checking post conditions and/or
	// keeping invariants
}

// V's autofree engine will not take care of memory management in this function.
// You will have the responsibility to free memory manually yourself in it.
[manualfree]
fn custom_allocations() {
}

// For C interop only, tells V that the following struct is defined with `typedef struct` in C
[typedef]
struct C.Foo {
}

// Used in Win32 API code when you need to pass callback function
[windows_stdcall]
fn C.DefWindowProc(hwnd int, msg int, lparam int, wparam int)

// Windows only:
// If a default graphics library is imported (ex. gg, ui), then the graphical window takes
// priority and no console window is created, effectively disabling println() statements.
// Use to explicity create console window. Valid before main() only.
[console]
fn main() {
}
```

## Goto

V允许用`goto`无条件地跳转到一个标签. 标签名称必须与`goto`语句包含在同一个函数中. 程序可以goto到当前作用域之外或更深的地方. `goto`允许跳过变量初始化或跳回访问已经释放的内存的代码, 所以它需要`unsafe`.

```v ignore
if x {
	// ...
	if y {
		unsafe {
			goto my_label
		}
	}
	// ...
}
my_label:
```
应避免使用`goto`, 特别是在可以使用 `for`的情况下. [带标签的break/continue](#带标签的break/continue)可以用来脱离嵌套循环, 这些不会有违反内存安全的风险.

# 附录

## 附录一：关键词

V有41个保留关键词(3个是字词):

```v ignore
as
asm
assert
atomic
break
const
continue
defer
else
embed
enum
false
fn
for
go
goto
if
import
in
interface
is
lock
match
module
mut
none
or
pub
return
rlock
select
shared
sizeof
static
struct
true
type
typeof
union
unsafe
__offsetof
```
可见[类型](#类型).

## 附录二: 运算符

这里只列出了[基础类型](#基础类型)的运算符.

```v ignore
+    sum                    integers, floats, strings
-    difference             integers, floats
*    product                integers, floats
/    quotient               integers, floats
%    remainder              integers

~    bitwise NOT            integers
&    bitwise AND            integers
|    bitwise OR             integers
^    bitwise XOR            integers

!    logical NOT            bools
&&   logical AND            bools
||   logical OR             bools
!=   logical XOR            bools

<<   left shift             integer << unsigned integer
>>   right shift            integer >> unsigned integer


Precedence    Operator
    5             *  /  %  <<  >>  &
    4             +  -  |  ^
    3             ==  !=  <  <=  >  >=
    2             &&
    1             ||


Assignment Operators
+=   -=   *=   /=   %=
&=   |=   ^=
>>=  <<=
```
