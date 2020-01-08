# c++
变量的名称称为标识符, 其必须以字母或下划线开头, 且仅包含`字母, 数字, 下划线`. c++是大小写敏感的.

使用namespace的原因: c++有大多东西需要命名, 一个名称可能有多种定义, 为消除歧义, 将不同的项划分到不同的集合中以避免冲突.

cin运行输入时以一个或多个空格或一个换行符来分隔输入的多个参数, 它以`Enter键`表示结束.

## header
```c++
#include <iostream> // 推荐, 新编译器支持
#include <iostream.h> // 不推荐, 旧编译器要求追加".h" 
```

## 类型
c++11支持auto类型, 由右侧表达式推断变量类型.
c++11 支持`decltype (expr) value`, 可根据expr确定变量value的类型.
c++没有字符串的原生类型, 但可使用`#include<string>`来定义和操作字符串.
新版编译器支持bool类型, 而不是用1/0来表示true/false.
c++11支持枚举类`enum class Days {Sun,...}; Days d = Days::Sun; `, 以避免传统枚举的一些问题.
break用于退出while, do-while,for, switch. 
新版c++使用预定义的强制类型转换函数`static_cast<type>(xxx)`来取代旧有的`type(xxx)`.

### 数值
type | byte |  range  | format|
-|-|-|-|
char | 1 | `-128 ~ 127` 或 `0 ~ 255` | %c |
unsigned char | 1 | `0 ~ 255` | %c |
signed char | 1 | `-128 ~ 127` | %c |
short/short int | 2 | `-32,768 ~ 32,767` | %hi |
unsigned short | 2 | `0 ~ 65535` | %hu |
int | 4 | `-2,147,483,648 ~ -2,147,483,647` | %li |
unsigned int | 4 | `0 ~ 4,294,967,295` | %lu |
long/`long int` | 8 | `−9,223,372,036,854,775,808 ~ −9,223,372,036,854,775,807` | %lli |
unsigned long/`unsigned long int` | 8 | `0 ~ 18,446,744,073,709,551,615` | %llu |
float | 4 | IEEE754, 精度7位 | %f |
double | 8 | IEEE754,精度15位 | %f |
long double | 10 | IEEE754, 精度19位 | %f |

推荐适用[定宽整数类型 from C++11](https://zh.cppreference.com/w/cpp/types/integer), 它定义于头文件 `<cstdint>`里.

## 字符串
c++11支持原始字符串字面值(raw string literals),类似于golang中的<code>``</code>, 它适用于太多字符需要转义的场景, 比如`cout<<R"(c:\fiiles\)";`

## 函数
部分数值处理的预定义函数:
![部分数值处理的预定义函数](/misc/img/c/截图_2020-01-08_18-00-46.png)

c++要求函数调用前要有具体的函数实现或存在函数声明.

C++ 引用传递:
形参相当于是**实参的"别名"**，对形参的操作其实就是对实参的操作，在引用传递过程中，被调函数的形式参数虽然也作为局部变量在栈中开辟了内存空间，但是这时存放的是由主调函数放进来的实参变量的地址. 被调函数对形参的任何操作都被处理成间接寻址，即通过栈中存放的地址访问主调函数中的实参变量. 正因为如此，被调函数对形参做的任何操作都影响了主调函数中的实参变量.

引用传递和指针传递是不同的，虽然它们都是在被调函数栈空间上的一个局部变量，但是任何对于引用参数的处理都会通过一个间接寻址的方式操作到主调函数中的相关变量. 而对于指针传递的参数，如果改变被调函数中的指针地址，它将影响不到主调函数的相关变量.

### 函数重载
函数名相同,  参数列表(个数,类型,顺序)不同

注意:
1. 返回值类型与函数重载无关
1. 调用函数时, 实参的隐式类型转换可能产生二义性.

本质: 采用了name mangling(也叫name decoration)的技术, 即c++编译器默认会对符号名(比如函数名)进行改编,修饰, 简单理解就是重命名, 且不同编译器有不同的重命名规则.

> 可通过编译时选禁止优化和不生成调试信息, 再通过反汇编验证.

### func list
- endl

从定义中看出，endl是一个函数模板，它实例化之后变成一个模板函数，作用是输出一个换行符，并立即刷新缓冲区.

## FAQ
### 编译器支持c++进度
- [C++ Standards Support in GCC](https://gcc.gnu.org/projects/cxx-status.html)
- [C++ Support in Clang](http://clang.llvm.org/cxx_status.html)
- [C++ 20的范围稳了，主流编译器的完整支持会在什么时候？](https://www.zhihu.com/question/313443905)
