# c
特点:
1. 效率高 : 编译型语言
1. 灵活度高 : 指针/指针运算
1. 可移植性高 : 靠近硬件

注释:
1. 单行 : `//`
1. 多行 : `/* ... */`

> 编译器剔除注释的方法是用空格代替注释.
> `y=x/*p`, 从`/*`开始会被认为是注释, 可用`y=x/(*p)`代替
> 头文件仅用于声明

三目操作符: `? :`, 比如`int b = a>4 ? 5 : 6 ;`

style: [Linux kernel coding style](https://www.kernel.org/doc/Documentation/process/coding-style.rst), 其翻译版在[Linux内核编码风格 Linux kernel coding style（中英对照）](http://iyu.is-programmer.com/posts/30315.html).

## 关键词
由ANSI标准定义的C语言关键词共32个 :
```
    auto double int struct break else long switch
    case enum register typedef char extern return union
    const float short unsigned continue for signed void
    default goto sizeof volatile do if while static
```

### 作用分类

根据关键词的作用，可以将关键词分为数据类型关键词和流程控制关键词两大类.

#### 数据类型关键词
A 基本数据类型（5个）

    1.void [vɔɪd]：声明函数无返回值或无参数，声明无类型指针，显式丢弃运算结果
    2.char [tʃɑr]：字符型类型数据，属于整型数据的一种
    3.int [ɪnt]：整型数据，通常为编译器指定的机器字长
    4.float [flot]：单精度浮点型数据，属于浮点数据的一种
    5.double [‘dʌb!]：双精度浮点型数据，属于浮点数据的一种

> void的真正作用是: 对函数返回的限定; 对函数参数的限定.

> 在c中, 凡是不加返回值类型限定的函数, 都会被编译器作为返回整型值处理, 因此函数没有返回值时一定要加void限定.

> 如何函数无参数, 那么应声明其参数为void, 即`func(void)`

> void是一种抽象的需要, 因此无法定义变量, 因为定义变量时需要分配具体的内存空间.

B 类型修饰关键词（4个）

    1.short [ʃɔrt]：修饰int，短整型数据，可省略被修饰的int
    2.long [lɔŋ]：修饰int，长整形数据，可省略被修饰的int
    3.signed [saɪnd]：修饰整型数据，有符号数据类型
    4.unsigned [ʌn’saɪnd]：修饰整型数据，无符号数据类型

C 复杂类型关键词（5个）

    1.struct [strʌkt]：结构体声明
    2.union [‘junjən]：共用体声明
    3.enum [i.nju:mə]：枚举声明
    4.typedef [taɪpdɛf]：声明类型别名
    5.sizeof [saɪzɑv]：得到特定类型或特定类型变量的大小

D 存储级别关键词（6个）

    1.auto [‘ɔto]：指定为自动变量，由编译器自动分配及释放。通常在栈上分配
    2.static [‘stætɪk]：指定为静态变量，分配在静态变量区，修饰函数时，指定函数作用域为文件内部
    3.register [‘rɛdʒɪstɚ]：指定为寄存器变量，**建议(不是绝对)**编译器将变量存储到寄存器中使用，也可以修饰函数形参，建议编译器通过寄存器而不是堆栈传递参数.
    4.extern [‘ɛkstɝn]：指定对应变量为外部变量，即在另外的目标文件中定义，可以认为是约定由另外文件声明的。
    5.const [‘kɔnstənt]：与volatile合称“cv特性”，指定变量不可被当前线程/进程改变（但有可能被系统或其他线程/进程改变）
    6.volatile [‘vɑlət!]：与const合称“cv特性”，指定变量的值有可能会被系统或其他进程/线程改变，强制编译器每次从内存中取得该变量的值

> register变量大小必须小于等于寄存器大小, 且是单个值; 因为register变量可能不在内存中, 因此不能使用`&`操作符.

#### 流程控制关键词
A 跳转结构（4个）

    1.return [rɪ’tɝn]：用在函数体中，返回特定值（或者是void值，即不返回值）
    2.continue [kən’tɪnju]：结束当前循环，开始下一轮循环
    3.break [brek]：跳出当前循环或switch结构
    4.goto ：无条件跳转语句

B 分支结构（5个）

    1.if [ɪf]：条件语句
    2.else [ɛls]：条件语句否定分支（与if连用）
    3.switch [swɪtʃ]：开关语句（多重分支语句）
    4.case [kes]：开关语句中的分支标记, 后面只能是`整型或字符型的常量`或`常量表达式`
    5.default [dɪ’fɔlt]：开关语句中的“其他”分治，可选。
   
> 建议switch必须带default分支(即使程序不需要), 以避免忘记default的处理; switch case组合中禁止使用return.

C 循环结构（3个）

    1.for ：for循环结构，for(1;2;3)4;的执行顺序为1->2->4->3->2…循环，其中2为循环条件
    2.do ：do循环结构，do 1 while(2); 的执行顺序是 1->2->1…循环，2为循环条件
    3.while [hwaɪl]：while循环结构，while(1) 2; 的执行顺序是1->2->1…循环，1为循环条件
    以上循环语句，当循环条件表达式为真则继续循环，为假则跳出循环

> 多重循环时, 建议将最长的循环放在最内层, 最短循环放在最外层, 以减少cpu跨切循环层的次数.

> 不能在循环中修改循环变量, 防止循环失控.

> 建议循环嵌套不超过3层, 避免难以理解; 超过时可以子函数的形式代替.

> 不要在循环的控制表达式中使用浮点类型, 避免精度问题导致不可预期的结果.

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

> 任何以0开头的数字(`0x`除外)都会被认为是八进制数.

> 所有无符号常量都应该带有字母后缀`U`.

> 布尔类型bTest使用`if(bTest)`写法好, 因为:`if (bTest==0)`,bTest会被误认为整型; `if(bTest == TRUE)`, TRUE不一定是1, 比如Visual Basic就把TRUE定义为`-1`.

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

> `~`和`<<`适用于无符号字符类型或无符号整型, 避免符号位的干扰.

> c规定(贪心法): 每个符号应该包含尽可能多的字符, 因此`a+++b`==`(a++)+b`, 推荐使用`()`取代贪心法.

[C语言运算符优先级 详细列表](https://cloud.tencent.com/developer/article/1171615),总结:
- 同一优先级的运算符，运算次序由结合方向所决定
- 简单记就是：！ > 算术运算符 > 关系运算符 > && > || > 赋值运算符

优先级易错情况归纳:
1. `.`的优先级高于`*` : `*p.f` == `*(p.f)`
1. `[]`高于`*` :  `int*ap[]` == `int* (ap[])` // ap是元素为int*的数组.
1. `函数()`高于`*` : `int *fp()` == `int *(fp())` // fp是函数, 返回`int*`
1.  ==和!=高于位操作: `(val & mask !=0)`==`val & (mask!=0)`
1.  ==和!=高于赋值: `c=getchart()!=EOF`==`c=(getchart()!=EOF)`
1. 算数运行符高于位运算: `msb <<4 + lsb` == `msb <<(4+lsb)`
1. 逗号运算符在所有运算符中 最低: `i=1,2` == `(i=1),2`

忘记运算符优先级时用`()`代替准没错.

### enum
```c
// 第一个枚举值默认是0, 没有指定值的枚举元素，其值为前一元素加 1. 也就说 spring 的值为 0，summer 的值为 3，autumn 的值为 4，winter 的值为 5
enum season {spring, summer=3, autumn, winter};
```

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

union的物理存储是多个成员变量共同占用同一块内存, 因此任何时候**只能使用一个成员变量**.

> union占用的内存是足够存储union中最大成员的大小.

> union需考虑大小端的影响, 因为对union的成员的存取都是从该union的基地址的偏移量为0处开始的即首地址开始的.

```c
int IsLittleEndian()//原理：联合体union的存放顺序是所有成员都从低地址开始存放，而且所有成员共享存储空间
{
	union temp
	{
		short a;
		char b;
	}temp;
	temp.a = 0x1234;
	return temp.b != 0x12 //低字节存的是数据的低字节数据
}
```

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
- buf是数组名称
- buf等价于`&buf[0]`

> buf 与 &buf 的值相等，但是含义完全不同, 比如`printf("%p\n",buf+1) 与 printf("%p\n",&buf+1) `的结果完全不同，因为它们的含义不同，buf 表示数组的第一个元素的首字节地址，加 1 加的时一个元素空间的大小；&buf 表示的是整个数组的首地址，加 1 加的是整个数组空间大小，数组首地址主要用于构建多维数组，对于一维数组来说，数组首地址没有太大的实用意义.

> 数组只有地址传递

`&a[0]`与`&a`的区别:
a[0]是一个元素, a是整个数组, 虽然`&a[0]`与`&a`的值一样, 但意义完全不同, 前者是数组首元素的首地址, 而后者是数组的首地址.

`char *p = "abc"`和`char a[] = "abc"`区别:
相同点: 都可用`*(x+i)`或`x[i](编译器还是会转成*(x+i))`的形式访问, 本质是`x+i*sizeof(x所指向的数据的数据类型)`.

不同点: 它们是不同类型的变量, p是指针, a是数组.
- `char *p="abc"` : 定义了一个字符串常量, p指向的内存在静态区, 是全局的, 且其内容只读. 新版编译器推荐写成`const char * p = "abcdef";`.
- `char a[]`:  定义了一个数组a, 需分配内存, 可修改.

```c
#include <stdio.h>

int main(void)
{    
	int a[5]={1,2,3,4,5};
	int* ptr=(int *)(&a+1); // 表示&a+5*sizeof(int), 因为`&a+1`=`&a+1*sizeof(a)`

	printf("%d , %d\n",*(a+1),*(ptr -1)); // 同上注释, `ptr -1` = ptr - 1*sizeof*(int) = a[5]
    ptr[-1] // =>*(ptr-1)
	return 0;
}
```

#### 二维数组
二维数组`int a[M][N]`是一维数组a[M]的扩展, a[M]的类型是`int[N]`.
参考一维数组可得a=`&a[0]`, 因此数组名 a 的类型为`int(*)[N]`, 而a[0]=`&a[0][0]`(把a[0]视为一个整体buf可轻松推导得到), 所以a = `&(&a[0][0])`

```c
char a[3][2]={{0,1},{2,3},{4,5}}; // `a[i][j]=*(&a[i][j])=*(a[i]+j) =*(&a[i][0]+j) =*(*(a+i)+j)` // 二维数组里a[i]也是个地址,而不是类似一维数组里的值. 一维数组中`a[i]=*(a+i)`
```

```c
int a[5][5];
int (*p)[4];
p = (int (*)[4])a;

&p[4][2] // = (int *)(p+4)+2
```

### struct
```c
// tag、member-list、variable-list 这 3 部分至少要出现 2 个
struct tag { // tag 是结构体标签
    member-list
    member-list 
    member-list  
    ...
} variable-list ; // variable-list 定义struct变量，最后一个分号(不能省略)之前，可以指定一个或多个struct变量
```

> 通常编译器将没有字段的struct的sizeof()默认为1, 但gcc中是0.

虽然c语言是面向过程的, 但也可以使用面向对象的思想来编写程序, 比如linux kernel:
```c
struct {
   int age;
   void (*pFunc)(void); // 函数指针, 指向void func(void), 类似class 的成员方法
}
```

#### 柔性数组(flexible array)
c99中, strut的最后一个元素(它前面需至少有一个其他成员)允许是未知长度的数组, 该数组即是柔性数组, sizeof返回的struct大小不包括柔性数组, 该数组用malloc()函数进行动态分配.

```c
typedef struct st_type
{
   int i;
   int a[];
}type_a;

type_a *p = (type_a *)malloc(sizeof(type_a)+100*sizeof(int));
sizeof(*p) // 4, 因为sizeof(type_a)为4即type_a大小不包括柔性数组
```

### 函数
```c
// 如果函数类型是 void,则该函数最后面就不再需要 return 语句进行返回了, 即该函数没有返回值
// 如何函数参数是指针, 且仅作为输入参数用, 则应该在类型前加const, 避免该指针在函数内被意外修改
datatype function_name(datatype parameters1 ,  datatype parameters2 , . . .)
{
   body of the function
}
```

> 函数名与数组名的最大区别: 函数名做右值时加不加`&`效果和意义都一样, 但数组名则不同.

> 在arm下传参不推荐操过4个(推测与寄存器数量有关), 超出时推荐使用结构体.

> 递归调用必须注意: 收敛和栈溢出.

> 库函数是预先写好的函数集合, 供人复用.

> `-lxxx`可指定链接的具体库(即libxxx), `-L`可附加so的查找位置.

函数指针是一个指针, 指向一个函数, 比如
```c
char* (*func)(char* p1, char* p2); // func为函数指针
```

```c
void function(){
    ...
}

int main(){
    void(*p)(); // 定义一个函数指针, 其参数列表和返回值均是void
    *(int *)&p=(int)function; // 将函数入口地址赋值给指针p. `(int)function`将函数入口地址强制转换成int类型的值.`*(int *)&p`:`&p`取指针p的地址p', 将p'转换成指向int类型数据的指针, 再进行赋值
    (*p)(); // 对函数的调用
    return 0;
}
```

`(*(void(*)())0)();`理解:
1. void(*)() : 是函数指针
1. `(void(*)())0`: 将0强制转换为函数指针类型,0为一个保存函数入口的地址
1. `(*(void(*)())0)`: 取到函数入口
1. 最后调用

函数指针数组:
```c
char* (*pf[3])(cha* p); // pf为函数指针数组
```

函数指针数组的指针:
```c
char* (*(*pf)[3])(cha* p); // pf为函数指针数组的指针
```

传递数组:
```c
//  [二维数组作为函数参数传递剖析(C语言)(6.19更新第5种)](https://www.cnblogs.com/wuyuegb2312/archive/2013/06/14/3135277.html)
#include <stdio.h>
#include <stdlib.h>

void myFunc1(int length, int *arr) {
    for (int i = 0; i < length; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
    return;
}

void myFunc2(int length, int arr[]) {
    for (int i = 0; i < length; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
    return;
}

void myFunc3(int length, int arr[length]) {
    for (int i = 0; i < length; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
    return;
}

void myFunc4(int row, int column, int *arr2) {
    for (int i = 0; i < row; i++) {
        for (int j = 0; j < column; j++) {
            printf("arr2[%d][%d] = %d   ", i, j, *(arr2 + i*column + j));
        }
        printf("\n");
    }
    return;
}

void myFunc5(int row, int column, int arr2[][column]) {
    for (int i = 0; i < row; i++) {
        for (int j = 0; j < column; j++) {
            printf("arr2[%d][%d] = %d   ", i, j, arr2[i][j]);
        }
        printf("\n");
    }
    return;
}

int main() {
    int arr[5] = {1, 2, 3, 4, 5};
     
    printf("Parameter is pointer:\n");
    myFunc1(5, &arr[0]);    // ok with arr
     
    printf("Parameter is array[]:\n");
    myFunc2(5, &arr[0]);    // ok with arr
     
    printf("Parameter is array[length]:\n");
    myFunc3(5, &arr[0]);    // ok with arr
     
    int arr2[2][3] = {{1, 2, 3}, {4, 5, 6}};
     
    printf("Parameter is pointer:\n");
    myFunc4(2, 3, arr2[0]); // &arr2[0][0] is ok, but arr2 is not ok
                            // *arr2 is ok.
                            // but actually they have the same address
     
    printf("Parameter is pointer:\n");
    myFunc5(2, 3, arr2);    // not ok with arr2[0]
                            // not ok with &arr2[0][0]
     
    printf("\narr2                           = %p\n", arr2);
    printf("arr2[0]                        = %p\n", arr2[0]);
    printf("*arr2 = arr2[0]                = %p\n", *arr2);
    printf("arr2[0][0]                     = %p\n", &arr2[0][0]);
    printf("**arr2 = *arr2[0] = arr2[0][0] = %p\n\n", &arr2[0][0]);
     
     
    printf("address of this two-dimentional array:\n");
    printf("arr2           = %p\n", arr2);
    for (int i = 0; i < 2; i++) {
        printf("  arr2[%d]      = %p\n", i, arr2[0]);
        for (int j = 0; j < 3; j++) {
            printf("    arr2[%d][%d] = %p   ", i, j, &arr2[i][j]);
        }
        printf("\n");
    }
     
    printf("we can use the unexisted element, interesting!\n");
    printf("*arr2 means arr2[0], the first row.\n");
    printf("**arr2 means *arr2[0] means arr2[0][0], the first element of the first row.\n");
    printf("arr2[1][0]               = %d\n", arr2[1][0]);
    printf("arr2[0][3]               = %d\n", arr2[0][3]);
    printf("*(*(arr2 + 1) + 0)       = %d\n", *(*(arr2 + 1) + 0));
    printf("*(*(arr2 + 0) + 3)       = %d\n", *(*(arr2 + 0) + 3));
    printf("*(arr2[0] + 1*3 + 0)     = %d\n", *(arr2[0] + 1*3 + 0));
    printf("*(*arr2 + 1*3 + 0)       = %d\n", *(*arr2 + 1*3 + 0));
    printf("*(&arr2[0][0] + 1*3 + 0) = %d\n", *(&arr2[0][0] + 1*3 + 0));
    printf("*(&arr2[0][0] + 3)       = %d\n", *(&arr2[0][0] + 3));
     
    return 0;
}
```

### 指针
```c
type *var-name;
```

> 指针是c语音的精髓.

> 指针与零值的比较:`if(p==NULL)`

> ANSI标准规定: 进行算法操作的指针必须确定它指向数据类型的大小, 因此void指针无法进行算法操作; 但gnu认为`void *`的算法操作与`char *`一致.

> 如果函数的参数是任意类型的指针, 那应声明其参数为`void *`.

```c
int *p1, p2; // p1是指针, p2是int型变量 => 该demo可体现golang类型后置的优势: 类型明确, 无歧义.
void (*p)(void); // p表示函数指针`void (*)(void)`
```

```c
int i = 10;
int* p= &i;
*p = NULL; // 因为`#define NULL 0`, 但大部分是`#define NULL ( (void *) 0)`
```

野指针是指向一个不确定的地址或引用的空间不确定, 危害:
1. 引发段错误

   段错误即地址错误, 是对程序和系统的保护性措施, 发生时程序立即终止, 避免雪崩式错误.

   > 段错误分类: 大段错误是指针指向的地址不存在; 小段错误: 指针指向地址存在, 但对该空间的操作权限受到限制.
1. 产生不可预知的结果/错误
1. 引发程序的连环错误

避免野指针的方法: 使用时初始化为NULL,  用完后再设为NULL.

### 字符串
在 C 语言中，字符串是使用 null 字符(`'\0'`)表示终止的一维字符数组, 是通过字符指针间接实现的.

`char a[10]={ 'A', 'B', 'C', 'D', 'E'};`<=>`char a[10]={ "ABCDE" };` // `{}`在这里与if{}的作用一样, 是打包, 使之成为一个整体, 并与外界绝缘.

> 双引号表示字符串常量, 单引号表示字符常量.

> go的string底层包含length, 不用null, 更严谨.

### 宏
所有的预处理器命令都是以`#`开头:
- #define :	定义宏
- #undef	: 取消已定义的宏, **不推荐使用**
- #include	: 包含另一个源代码文件, 它支持相对路径
- #if	 :   如果给定条件为真，则编译它与`#endif`间的代码, 否则跳过这些代码
- #else	 :   #if 的替代方案,与c语言中的else类似
- #elif	 :   与c语言的`else if`类似, 如果前面的 #if 给定条件不为真，当前条件为真，则编译下面代码
- #endif	 : 标识一个 #if……#else 条件编译块的结束
- #ifdef	: 如果宏已经定义，则返回真
- #ifndef :	如果宏没有定义，则返回真
- #if defined, 支持比`#ifdef`更复杂的预编译条件
- #if !defined
- #line : 改变当前行数和文件的名称, 格式为`#line number["filename"]`, 常用于编译器生成的中间文件中, 用于保证代码位置是固定的, 不会被替换, 便于分析.
- #error	: 编译程序时只要遇到`#error`，就输出错误消息, 并停止编译
- #pragma	使用标准化方法，向编译器发布特殊的命令到编译器中, 用于设置编译器的状态或指示编译器完成一些特定的动作.

ANSI标准的预定义宏:
- __DATE__	: 表示编译时刻的当前日期，一个以 "MMM DD YYYY" 格式表示的字符常量
- __TIME__	: 表示编译时刻的当前时间，一个以 "HH:MM:SS" 格式表示的字符常量
- __FILE__	: 正在编译的文件的当前文件名，一个字符串常量
- __LINE__	: 正在编译的文件的当前行号，一个十进制常量
- __STDC__	当编译器以 ANSI 标准编译时，则定义为 1

可使用`gcc -E xxx.c`查看

**用宏定义表达式时不能吝啬括号, 特别是最外层的括号**.
**函数宏有参数的话, 调用时就不能缺少参数**
**定义函数宏时, 每个参数实例都必须用小括号包裹, 除非它们是作为`#`或`##`的操作数**
**预处理指令中的所有宏标识符在使用前必须已定义, 除了`#ifdef`,`#ifndef`,`defined()操作符`**
**`#ifdef ...#elif...#else...#endif`关联指令应该放在同一个文件中便于维护**

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
#define tokenpaster(n) printf ("token" #n " = %d", token##n) // `token##n`表示拼接的变量名, 因此需多一个`#`

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

#### pragma
```c
#ifdef _X86
#Pragma message("_86 arch activated!") // 提示定义了特定的宏
#endif
```

## 修饰符
### typedef
给已存在的数据类型取别名, 可理解为`typerename`.

优点:
1. 简化类型, 让程序更易理解和书写
1. 定义平台无关类型, 便于移植

理解typedef: 去掉typedef, 再将typedef定义的类型看做变量声明, 查看声明了什么变量即可.

```c
// 一次定义两个类型
typedef struct node{} Node, *pNode; // = typedef 类型 Node + typedef 类型* pNode
// struct node n == Node n
// struct node* m == pNode m == Node* m 
typedef int *pint
const pint p2; // = int *const p2;
pint const p1; // = int *const p1;
```

记忆方法: 去掉别名和typedef看, 剩余部分的类型.

> typedef在语法上是一个存储类的关键字(如auto, extern, static, register), 而变量只能被一种存储类的关键词修饰.

> 想一次性定义多个指针变量, 需使用typedef, 否则会发生歧义.

### static
1. 变量
静态全局变量: 作用域是从定义之处开始到文件结尾处, 其他文件即使用extern声明都没法使用它.
静态局部变量: 虽然该变量在静态存储区, 但作用域仅限该函数.


```c
#include <stdio.h>

int fun(void){
    static int count = 10;    // 仅程序初始化时运行过一次, 之后此赋值语句就不再执行过(即函数fun中没有该赋值语句了). 该变量在gcc汇编中会被重命名为`count.${N}`以避免和全局变量count冲突.
	count--;
		
	printf("---%p\n",&count);    
	
    return count;
}

int count = 1;

int main(void)
{    
    printf("global\t\tlocal static\n");
    for(; count <= 10; ++count)
        printf("@@@%p,%d\t\t%d\n",&count, count, fun());    
   
    return 0;
}
```

2. 函数
函数前加static使得函数变成静态函数, 作用域是本文件即不同文件使用相同名称的静态函数互不影响, 因此它也被称为内部函数.

### const
修饰的变量不可变, 可将const理解为readonly.

> 编译器通常不为普通的const变量分配存储空间, 而是将它们保存在符号表中, 使得它成功一个编译期间的值.

const实现:
1. 将const修饰的变量放在代码段, 因为代码段只读. 常见于嵌入式程序.
1. const是编译器实现, 放在数据段, 在编译时检查排错, 因此运行时可通过指针修改.

> `int const a=10` <=> `const int a=10`

> 函数传参使用const指针: 效率 + 不能修改指向空间的内容.

const修饰指针的三种形式:
1. `const int *p`<=>`int const *p`, const修饰`*p`即`*p(p指向的空间内容)`是常量(即不能通过`*p`赋值), 但p可变
1. `int *const p`, p不可变, 但p指向的空间内容可变
1. `int const* const p`<=>`const int *const p`, p不可变, p指向的空间也不可变

记忆方法: 先忽略类型, 再看const离哪个近, 离谁近就修饰谁. 比如`int const* const p`=>`const * const p`, 前一个const修饰了`*p`, 后一个const修改了p.

const也可修饰函数的返回值, 即返回值不可改变.

引用另一文件的const变量时使用`extern const int i`.

> go的const是`本行的赋值表达式 = 上一行的赋值表达式`.

### volatile
用volatile修饰的变量表示可以被某些编译器未知的因素更改, 比如操作系统, 硬件, 或其他线程. 因此使用volatile修饰的变量, 编译器对访问该变量的代码不做优化, 从而可以提供对特殊地址的稳定访问.

### void
万能类型.

> malloc返回的是`void *`类型指针, 需自行强制转换类型; 其失败会返回`NULL`. gcc中malloc默认最小是以16B为单位进行分配.

### extern
修饰变量或函数, 表示该对象的定义在其他文件中, 提示链接器遇到该对象时去其他模块中解析此标识符.

`extern char a[]`和`extern char a[100]`没有区别, 因为仅是声明不分配空间, 因此编译器无需知道数组的长度.

## FAQ
### 行末尾的`\`(断行)
如果一行代码有很多内容，导致太长影响阅读，可以通过在结尾加`\`的方式实现换行，编译时会忽略\及其后的换行符，并当做一行处理.

### 内存对齐
内存对齐是硬件问题: 对齐访问更符合硬件规律, 因此效率更高.

在c语言中可通过`#prgama pack(N)`来自定义对齐. gcc中推荐使用`_attribute_((packed))`来取消对齐和`_attribute_((aligned(N)))`来自定义对齐. **通常不推荐改动对齐方式**.

### 函数名的本质
函数是一段代码的封装, 其函数名是指向这段代码的首地址, 即函数名的本质是内存地址.

### 指针类型
指针类型是该指针指向的内存空间的解析规则.

### 堆
heap是一种动态内存管理方法, 通过malloc和free来使用, 其生命周期由free或delete决定.

特点:
1. 容量不限, 动态分配
1. 申请和释放需要手动

### 栈
保存局部变量.

### 静态存储区
静态存储区保存静态(static)变量和全局变量, 在整个程序的生命周期内都存在, 由编译器在编译时分配,在程序刚运行时初始化一次且仅一次.

在静态数据区，内存中所有的字节默认值都是 0x00, 即static变量都有默认初始值`0`.


```c
struct S4{ 
    char a; 
    long b; 
    static long c; // 静态变量存放在全局数据区内，而sizeof计算栈中分配的空间的大小，故不计算在内，S4的大小为4+4=8.
}; 
```

### 进程中的变量/常量
1. 变量保存在数据区(.data和 .bbs), 栈, 堆等位置, 可读可写
1. 常量保持在.ro.data中, 只读

### typedef和`#define`区别
两者都可以定义别名, 但`#define`只是简单的宏替换, 而typedef不是;  `#define`是预编译时处理, 而typedef是编译时处理.

### const与`#define区别`
1. 编译器处理方式不同
   - define宏是在**预编译**阶段展开
   -  const常量是编译运行阶段使用

1. 类型和安全检查不同
   - define宏没有类型，不做任何类型检查，仅仅是展开
   - const常量有具体的类型，在编译阶段会执行类型检查

为了安全, 定义宏常数时用const代替, 让编译器做类型校验, 减少错误的可能. 注意: const修饰的是readonly的变量, 因此不能作为数组的维度, 也不能跟在case语句的后面.

1. 存储方式不同
   -  define宏仅仅是展开，有多少地方使用，就展开多少次，不会分配内存（宏定义不分配内存，变量定义分配内存)
   - const常量会在内存中分配(可以是堆中也可以是栈中)

(4)const  可以节省空间，避免不必要的内存分配
```c
        #define PI 3.14159 //常量宏 
        const doulbe Pi=3.14159; //此时并未将Pi放入ROM中 ...... 
        double i=Pi; //此时为Pi分配内存，以后不再分配！ 
        double I=PI; //编译期间进行宏替换，分配内存 
        double j=Pi; //没有内存分配 
        double J=PI; //再进行宏替换，又一次分配内存!
```

const定义常量从汇编的角度来看，只是给出了对应的内存地址，而不是象#define一样给出的是立即数，所以，const定义的常量在程序运行过程中只有一份拷贝（因为是全局的只读变量，存在静态区），而 #define定义的常量在内存中有若干个拷贝. 

1. 提高了效率
编译器通常不为普通const常量分配存储空间，而是将它们保存在符号表中，这使得它成为一个编译期间的常量，没有了存储与读内存的操作，使得它的效率也很高.

1. 宏替换只作替换，不做计算，不做表达式求解;
宏预编译时就替换了，程序运行时，并不分配内存

### enum与`#define`区别
1. 编译器处理方式不同
   - define宏是在**预编译**阶段展开
   -  enum是编译运行阶段使用

2. 调试方式
一般在调试器中可调试enum常量, 但不能调试宏常量.

3. 定义的变量数量
   - enum可以一次定义多个相关常量
   - define宏一次只能定义一个

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

### 带参数的宏定义, 函数和内联函数的区别
宏定义是在预处理阶段处理, 是原地展开, 没有调用开销. 它也不检查参数的类型, 返回值也不会附带类型, 编译器不做静态类型检查.

函数是编译阶段处理, 是通过函数跳转实现的. 它有明确的参数类型和返回值类型, 编译器会做静态类型检查.

内联函数通过在函数定义前加`inline`来声明, 其本质是函数, 但也同时具备带参宏的有点即原地展开, 缺点是不适合用于较长代码.

总结: 代码较多时用函数合适; 对于那些仅一两句的函数适合用带参宏.

### 定义与声明的区别
定义:  (编译器)创建一个对象, 为这个对象分配一块内存并给该段内存取个名称. 一个对象在一定区域内(比如函数内, 全局等)只能定义一次.
声明: 告诉编译器, 该名称已被使用. 声明可以出现多次, 且不会分配内存.

```c
extern int i;       //声明，不是定义
int i;                      //定义会包含声明
```

### struct和class区别
struct的成员默认的属性是public, 而class成员的默认是private.

### 指针与数组的区别
1. 存储内容
- 指针: 保存数据的地址, 任何存入指针变量p的数据都会被当做地址来处理
- 数组: 保存数据, 数组名a代表数组首元素的地址, 而不是数组的首地址, `&a`才是数组的首地址, 尽管a和`&a`的值相同, 但意义完全不同.

1. 场景
- 指针: 常用于动态数据结构
- 数组: 用于存储成员个数固定且数据类型相同的元素

### 指针数组和数组指针
指针数组: 是数组, 成员是指针, 比如`int* p1[10]`.
数组指针: 是指针, 指向一个数组, 比如`int (*p2)[10]`.

记忆: `[]`的优先级比`*`高.