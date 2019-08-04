# c
特点:
1. 效率高 : 编译型语言
1. 灵活度高 : 指针/指针运算
1. 可移植性高 : 靠近硬件

注释:
1. 单行 : `//`
1. 多行 : `/* ... */`

三目操作符: `? :`, 比如`int b = a>4 ? 5 : 6 ;`

style: [Linux kernel coding style](https://www.kernel.org/doc/Documentation/process/coding-style.rst), 其翻译版在[Linux内核编码风格 Linux kernel coding style（中英对照）](http://iyu.is-programmer.com/posts/30315.html).

## 变量类型
**c声明变量时没有零值, 需显式初始化**.

### 数值
C数据类型的大小 :　现今所有64位的类Unix平台均使用LP64数据模型，而64位Windows使用LLP64数据模型. 但为了避免依赖OS的数据模型和编译器决定数值长度的怪异行为, iso c99引入了[固定长度的数值类型即Fixed width integer types, 在stdint.h里,**推荐使用**](https://en.cppreference.com/w/c/types/integer).

LP64:
type | byte |  range  | format|
-|-|-|-|
char | 1 | `-128 ~ 127` 或 `0 ~ 255` | %c |
unsigned char | 1 | `0 ~ 255` | %c |
signed char | 1 | `-128 ~ 127` | %c |
short | 2 | `-32,768 ~ 32,767` | %hi |
unsigned short | 2 | `0 ~ 65535` | %hu |
int | 4 | `-2,147,483,648 ~ -2,147,483,647` | %li |
unsigned int | 4 | `0 ~ 4,294,967,295` | %lu |
long/`long long` | 8 | `−9,223,372,036,854,775,808 ~ −9,223,372,036,854,775,807` | %lli |
unsigned long/`unsigned long long` | 8 | `0 ~ 18,446,744,073,709,551,615` | %llu |
float | 4 | IEEE754 | %f |
double | 8 | IEEE754 | %f |

其他format:
- `%d` : 数值
- `%e` : 科学计数
- `%s` : 字符串
- `%o` : 无符号八进制
- `%x` : 16进制
- `%X` : 大写16进制
- `\t` : 输出一个制表符
- `%n` : 输出一个换行符

> 为了得到某个类型或某个变量在特定平台上的准确大小，可使用`sizeof(type)`得到对象或类型的存储字节大小.

> 整型类型转换时, 仅转换变量类型, **原先类型底层存储的补码不变**; 转换中发生空间收缩时, 保留的字节由机器的大小端确定.

#### 运算符
非0为true, 否则为false.
定位清零用`&`: a &= 0xFFFF00FF
定位置1用`|`: a |= 0x0000FF00
定位取反用`^`: a ^= 0x0000FF00

构造特定位置为0/1的二进制数: `移位`+`|`+`取反`.

sizeof:
```c
char str[]="hello";
sizeof(str); // 6, 包括`\0`
strlen(str); // 5, 因为不包括`\0`
char *p = str;
sizeof(*p); // 1, `*p`即`str[0]`
strlen(p); // 5, 因为p是字符串的首地址
```

### enum
```c
// 没有指定值的枚举元素，其值为前一元素加 1. 也就说 spring 的值为 0，summer 的值为 3，autumn 的值为 4，winter 的值为 5
enum season {spring, summer=3, autumn, winter};
```

> go的const是`本行的赋值表达式 = 上一行的赋值表达式`

### union
```c
union [tag]
{
   member definition;
   member definition;
   ...
   member definition;
} [one or more union variables];
```

union的物理存储是多个成员变量共同占用同一块内存, 因此任何时候只能使用一个成员变量.

> union占用的内存是足够存储union中最大成员的大小.

### 数组
```c
// 省略掉了数组的大小时, 数组的大小为初始化时元素的个数
type arrayName [ arraySize ];
```

优点: 定义简单, 访问方便
缺点:
1. 数组中元素类型必须一致, 大小必须在定义时给出.
1. 空间连续

`int buf[100] = {0}`解读:
- buf: 数组名称; 等价于`&buf[0]`; 等价于`&buf`

> 数组只有地址传递

> 二维数组: a[0]=`&a[0][0]`, 又因为a等价于`&a[0]`, 因此a = `&&a[0][0]`

### struct
```c
// tag、member-list、variable-list 这 3 部分至少要出现 2 个
struct tag { 
    member-list
    member-list 
    member-list  
    ...
} variable-list ;
```

虽然c语言是面向过程的, 但也可以使用面向对象的思想来编写程序, 比如linux kernel:
```c
struct {
   int age;
   void (*pFunc)(void); // 函数指针, 指向void func(void), 类似class 的成员方法
}
```

### 函数
```c
// 如果函数类型是 void,则该函数最后面就不再需要 return 语句进行返回了, 即该函数没有返回值
datatype function_name(datatype parameters1 ,  datatype parameters2 , . . .)
{
   body of the function
}
```

> 函数名与数组名的最大区别: 函数名做右值时加不加`&`效果和意义都一样, 但数组名则不同.

### 指针
```c
type *var-name;
```

> 指针是c语音的精髓.

```c
int *p1. p2; // p1是指针, p2是int型变量 => 该demo可体现golang类型后置的优势: 类型明确, 无歧义.
void (*p)(void); // p表示函数指针`void (*)(void)`
```

野指针是指向一个不确定的地址或引用的空间不确定, 危害:
1. 引发段错误

   段错误即地址错误, 是对程序和系统的保护性措施, 发生时程序立即终止, 避免雪崩式错误.

   > 段错误分类: 大段错误是指针指向的地址不存在; 小段错误: 指针指向地址存在, 但对该空间的操作权限受到限制.
1. 产生不可预知的结果/错误
1. 引发程序的连环错误

### 字符串
在 C 语言中，字符串是使用 null 字符(`'\0'`)表示终止的一维字符数组.

> go的string底层包含length, 不用null, 更严谨.

### 宏
所有的预处理器命令都是以`#`开头:
- #define	定义宏
- #include	包含一个源代码文件
- #undef	取消已定义的宏
- #ifdef	如果宏已经定义，则返回真
- #ifndef	如果宏没有定义，则返回真
- #if	    如果给定条件为真，则编译下面代码
- #else	    #if 的替代方案
- #elif	    如果前面的 #if 给定条件不为真，当前条件为真，则编译下面代码
- #endif	结束一个 #if……#else 条件编译块
- #error	当遇到标准错误时，输出错误消息
- #pragma	使用标准化方法，向编译器发布特殊的命令到编译器中

预定义宏:
- __DATE__	当前日期，一个以 "MMM DD YYYY" 格式表示的字符常量
- __TIME__	当前时间，一个以 "HH:MM:SS" 格式表示的字符常量
- __FILE__	这会包含当前文件名，一个字符串常量
- __LINE__	这会包含当前行号，一个十进制常量
- __STDC__	当编译器以 ANSI 标准编译时，则定义为 1

可使用`gcc -E xxx.c`查看

#### 预处理器运算符
- # : 在宏定义中，把一个宏的参数转换为字符串常量
- ## : 允许在宏定义中两个独立的标记被合并为一个标记

```c
#define  message_for(a, b)  \
    printf(#a " and " #b ": We love you!\n")

int main(void)
{
   message_for(Carole, Debra); // -> printf("Carole" " and " "Debra" ": We love you!\n");
   return 0;
}
```

```c
#define tokenpaster(n) printf ("token" #n " = %d", token##n)

int main(void)
{
   int token34 = 40;

   tokenpaster(34); // -> printf ("token" "34" " = %d", token34); 
   return 0;
}
```

#### 参数化的宏
```c
#define MAX(x,y) ((x) > (y) ? (x) : (y))

int main(void)
{
   printf("Max between 20 and 10 is %d\n", MAX(10, 20)); // -> printf("Max between 20 and 10 is %d\n", ((10) > (20) ? (10) : (20)));
   return 0;
}
```

## 修饰符
### typedef
给类型取别名.

优点:
1. 简化类型, 让程序更易理解和书写
1. 定义平台无关类型, 便于移植

理解typedef: 去掉typedef, 再将typedef定义的类型看做变量声明, 查看声明了什么变量即可.

```c
// 一次定义两个类型
typedef struct node{} Node, *pNode; // = typedef 类型 Node + typedef 类型 *pNode
typedef int *pint
const pint p2; // = int *const p2;
pint const p1; // = int *const p1;
```

> typedef在语法上是一个存储类的关键字(如auto, extern, static, register), 而变量只能被一种存储类的关键词修饰.

> 想一次性定义多个指针变量, 需使用typedef, 否则会发生歧义.
### static
### const
修饰的变量不可变

> const是编译器实现, 运行时可通过指针修改.

> `int const a=10` <=> `const int a=10`

> 函数传参使用const指针: 效率 + 不能修改指向空间的内容.

const修饰指针的三种形式:
1. `const int *p`<=>`int const *p`, 表示p指向的空间是常量(即不能通过`*p`赋值), 但p可变
1. `int *const p`, p不可变, 但p指向的空间内容可变
1. `int const* const p`, p不可变, p指向的空间也不可变

## FAQ
### 行末尾的`\`
如果一行代码有很多内容，导致太长影响阅读，可以通过在结尾加`\`的方式实现换行，编译时会忽略\及其后的换行符，并当做一行处理.

### 内存对齐
内存对齐是硬件问题: 对齐访问更符合硬件规律, 因此效率更高.

### 函数名的本质
函数是一段代码的封装, 其函数名是指向这段代码的首地址, 即函数名的本质是内存地址.

### 指针类型
指针类型是该指针指向的内存空间的解析规则.

### 堆
heap是一种动态内存管理方法, 通过malloc和free来使用.

特点:
1. 容量不限, 动态分配
1. 申请和释放需要手动

### 静态存储区
静态存储区保存静态局部(static)变量和全局变量.

### 程序中的变量/常量
1. 变量保存在.data, .bbs, 栈, 堆等位置, 可读可写
1. 常量保持在.ro.data中, 只读

### `#define`和`typedef`区别
两者都可以定义别名, 但`#define`只是简单的宏替换, 而typedef不是;  `#define`是预编译时处理, 而typedef是编译时处理.

### 复杂表达式的拆解
拆解方法:
1. 确定核心
1. 找结合: 谁跟核心最近, 谁先跟核心结合

   如果核心和`*`号结合表示核心是指针，如果核心和`[]`结合表示核心是数组，如果核心和`()`结合表示核心是函数
1. 继续向外结合直至整个表达式介绍

举例:
- `int *p[5]` : p是核心, `[]`比`*`优先级更高, 因此p是数组，数组中的5个元素都是指针，指针指向int型，所以`int *p[5]`是一个指针数组
- `int (*p)[5]` : p是核心, 因为`()`的优先级变更, 因此p是一个指针，指向一个数组，数组有5个元素都是int类型，所以`int (*p)[5]`是一个数组指针
- `int *(p[5])`: 是一个指针数组，结合方式同第一个一样, 这里的`()`可忽略