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
> 整型类型转换时, 仅转换变量类型, **原先类型底层存储的补码不变**

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

### 函数
```c
// 如果函数类型是 void,则该函数最后面就不再需要 return 语句进行返回了, 即该函数没有返回值
datatype function_name(datatype parameters1 ,  datatype parameters2 , . . .)
{
   body of the function
}
```

### 指针
```c
type *var-name;
```

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

## typeof
### 修饰符
1. static
1. const

    不允许再次改变

## FAQ
### 行末尾的`\`
如果一行代码有很多内容，导致太长影响阅读，可以通过在结尾加`\`的方式实现换行，编译时会忽略\及其后的换行符，并当做一行处理.