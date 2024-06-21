# c++
ref:
- [C++ 从遗忘到入门](https://my.oschina.net/u/4662964/blog/11050396)
- [现代 C++ 教程](https://changkun.de/modern-cpp/)

传统 C++ : C++98 及其之前的 C++ 特性
现代 C++ :  C++11/14/17/20, 为传统 C++ 注入的大量特性使得整个 C++ 变得更加像一门现代化的语言.

C++ 是一种通用编程语言, 以其强大的功能、高效的性能和灵活性而著称. 以下是一些关键特点:
- 面向对象：C++ 支持面向对象编程（OOP）的四大特性：封装、继承、多态和抽象。通过类和对象，程序员能够创建模块化的代码，更容易地进行维护和扩展。
- 泛型编程：C++ 支持模板编程，允许编写与数据类型无关的代码。模板是实现泛型编程的关键工具，它们提高了代码的复用性。
- 直接内存管理：C++ 提供了对内存的直接操作能力，允许程序员手动管理内存分配和释放，这是 C++ 的一个强大特性，也是需要谨慎使用的地方，因为不当的内存管理可能会导致资源泄露和其他问题。
- 性能：C++ 编写的程序通常有很高的执行效率，这是因为 C++ 提供了与底层硬件直接对话的能力。这使得 C++ 成为开发要求性能的系统软件（如操作系统、游戏引擎）的理想选择。
- C 语言兼容：大部分 C 语言程序可以在 C++ 编译器上编译并运行。这一特性简化了从 C 到 C++ 的过渡。
- 多编程范式支持：除了面向对象和泛型编程外，C++ 还支持过程式编程和函数式编程等范式，使其成为一个多样化的工具，能适应不同的编程需求。

## base
变量的名称称为标识符, 其必须以字母或下划线开头, 且仅包含`字母, 数字, 下划线`. c++是大小写敏感的.

使用namespace的原因: c++有大多东西需要命名, 一个名称可能有多种定义, 为消除歧义, 将不同的项划分到不同的集合中以避免冲突.

cin运行输入时以一个或多个空格或一个换行符来分隔输入的多个参数, 它以`Enter键`表示结束.

箭头操作符`->`指定由一个指针变量指向的struct/class的对象的成员.

> 形参和实参的功能是作数据传送. 发生函数调用时， 主调函数把实参的值传送给被调函数的形参从而实现主调函数向被调函数的数据传送.

## header
```c++
#include <iostream> // 推荐, 新编译器支持
#include <iostream.h> // 不推荐, 旧编译器要求追加".h" 
```

list:
```c++
#include <array> // for 标准容器array
#include <cctype> // for 预定义字符函数
#include <cstdlib> // for exit
#include <fsstream> // for file I/O
#include <iostream> // 用于cout
#include <istream> // 通用型流参数
#include <regex> // for regexp
#include <thread> // for 线程/并发编程
#include <type_traits> // for decltype 用于类型推导，而 std::is_same 用于比较两个类型是否相等
```

大多c++编译器支持`#pragma once`确保源代码文件编译时只包含一次, 也可用它代替`#ifndef`.

## namespace
其实c++的每一句代码都在某个namespace, 未明确指定时是在默认的全局命名空间(global namespace).

using的作用域即它所在的代码块(从using指令开始到该代码块结束), 如果using指令在文件开头, 则作用域是到整个文件结束为止.

```c++
namespace XXX // 定义namespace
{
	...
}

using namespace XXX; // 使用namespace

// 在本文件使用无名命名空间成员时不必用命名空间限定. 其实无名命名空间和static是同样的道理，都是只在本文件内有效，无法被其它文件引用. 因此无名命名空间可以替代全局作用域的static数据成员.
namespace // 无名命名空间 它的作用域是无名命名空间声明开始到本文件结束
{
	
}
```

using 声明`using  std::cout`和using指令`using namespace std;`的区别:
1. using声明只让命名空间中的一个名称进入可用状态; using指令使那个命名空间中的所有名称都进入可用状态.
1. using声明在引入一个名称, 不允许该名称再有其他用途; using指令只是隐式引入了命名空间中的名称, 只要代码不实际使用冲突名称则没问题.

编译单元(compilation unit)是一个文件, 它可能是类的实现文件, 以及通过`#include`包含的其他所有文件, 比如类的接口文件(即头文件). 每个编译单元都有一个无名命名空间.

## 类型
C++ 是静态编译语言，所有变量在声明时都要指定具体的变量类型，或者能让编译器推导出具体的变量类型（比如使用 auto、decltype 关键字的场景），类型检查不通过将导致编译期出错.

c++11支持auto类型, 由右侧表达式推断变量类型, 但**不能用于传参和推导数组类型**.
c++11 支持`decltype (expr) value`, 可根据expr确定变量value的类型.
c++没有字符串的原生类型, 但可使用`#include<string>`来定义和操作字符串.
新版编译器支持bool类型, 而不是用1/0来表示true/false.
c++11支持枚举类`enum class Days {Sun,...}; Days d = Days::Sun; `, 以避免传统枚举的一些问题.
break用于退出while, do-while,for, switch. 
新版c++使用预定义的强制类型转换函数`static_cast<type>(xxx)`来取代旧有的`type(xxx)`.
向量作用和数组一样, 但支持运行时改变长度.
c++要求定义数组时，必须明确给定数组的大小(即需要常量表达式. 注:const常数不是常量表达式), 否则应使用new 动态定义数组.
**c++使用了new时就需自己管内存(即delete), 否则编译器处理**.

> C字符串(以`0`作为终止符的字符数组)的本质是字符数组.
> C字符串变量初始化时可省略数组长度, c++自动设为`字符串长度+1`.
> `char a[]={'a', 'b'};`不会追加`\0`.
> 使用操作符new创建的变量称为动态变量, 需初始化, 如果类型是带有构造函数的类, 那么创建动态变量时会调用默认的构造函数.
> 对一个指针进行delete操作后没有把它置空的指针称为虚悬指针. 当delete一个指针的时候，实际上仅让程序释放了它所指向的内存，标明这块内存区域可用，但指针本身仍然存在，并且指向原来的内存区域.
> 普通变量(自动变量): 函数中的**局部变量**，如果不用关键字static加以声明，编译系统对它们是动态地分配存储空间的, 在函数调用结束时就自动释放这些空间的变量.

```c++
int *p = new int(12);
delete p; // p指向的内存被释放, 且p变成未初始化的变量.
```

### 列表初始化
现代 C++ 提供了一种新的统一的变量初始化方式 - 列表初始化，推荐优先使用这种初始化方式，它能提供更加直观和统一的数据初始化方式.

列表初始化使用 {} 来初始化数据对象，包括基础类型、数组、结构体、类和容器等复杂的数据类型:
```c++
// 基础类型
int a{0};  
double b{3.14}; 

// 结构体
struct MyStruct {
    int x;
    double y;
};

MyStruct s{1, 2.0};

// 类
class MyClass {
public:
    MyClass(int a, double b) : a_(a), b_(b) {}
private:
    int a_;
    double b_;
};

MyClass obj{5, 3.14}; // MyClass 必须有一个匹配这个参数列表的构造函数

// 数组
int arr[3]{1, 2, 3};
int arr[] = {1, 2, 3, 4, 5};     // 数组大小为5，编译器自动确定
int arr[10] = {1, 2, 3};    // 数组前三项确定为1，2，3，其余被初始化为0
int arr[10] = {0};         // 整个数组全部为0


// 上面介绍的都是现代C++推荐写法，省略 = 
// 下面的2种写法绝大多数情况下是等价的
float arr[2]{1, 2};        // 写法1
float arr[2] = {1, 2};    // 写法2
// 编译器对这两种写法的处理是一致的，方法2并不会产生临时变量和拷贝赋值，包括类的声明
```

现代 C++ 推荐优先使用列表初始化来初始化变量，因为这种方式不允许进行窄化转换这能避免一些问题的发生：
```c++
int a = 7.7;   // 编译能通过，但是有warning
int b = {1.0}; // 编译器拒绝通过，因为浮点到整形的转换会丢失精度
```

列表初始化支持参数列表小于数据对象的个数，这种情况下会默认进行其他变量的零初始化.

### 基础类型
#### 数值
type | byte |  range  | format|
-|-|-|-|
short/short int | 2 | `-32,768 ~ 32,767` | %hi |
unsigned short | 2 | `0 ~ 65535` | %hu |
int | 4 | `-2,147,483,648 ~ -2,147,483,647` | %li |
unsigned int | 4 | `0 ~ 4,294,967,295` | %lu |
long/`long int` | 8 | `−9,223,372,036,854,775,808 ~ −9,223,372,036,854,775,807` | %lli |
unsigned long/`unsigned long int` | 8 | `0 ~ 18,446,744,073,709,551,615` | %llu |
long long|8|c++11新增||
unsigned long long|8|c++11新增||
定宽整数 (从 <cstdint> 导入)|8/16/32/64|int8_t, int16_t, int32_t, int64_t 等||
无符号定宽整数 (从 <cstdint> 导入)|8/16/32/64|uint8_t, uint16_t, uint32_t, uint64_t 等||
float | 4 | IEEE754, 精度7位 | %f |
double | 8 | IEEE754,精度15位 | %f |
long double | 实现依赖, 比如10 | 扩展精度浮点数，精度和大小由具体实现定义. 比如IEEE754, 精度19位 | %f |

**推荐使用[定宽整数类型 from C++11](https://zh.cppreference.com/w/cpp/types/integer)**, 它定义于头文件 `<cstdint>`里.

#### 布尔
type | byte |  range  | format|
-|-|-|-|
bool | 实现依赖 | | 表示布尔值 true 或 false |

#### 字符
type | byte |  range  | format|
-|-|-|-|
char | 1 | `-128 ~ 127` 或 `0 ~ 255` | %c, 可表示字符或小整数，有符号性由实现定义 |
unsigned char | 1 | `0 ~ 255` | %c， 明确的有符号字符类型 |
signed char | 1 | `-128 ~ 127` | %c, 无符号字符类型 |
char16_t|16|C++11 新增，用于 UTF-16 字符||
char32_t|32|C++11 新增，用于 UTF-32 字符||
wchar_t|实现依赖|用于宽字符集||

#### 特殊类型
type | byte |  desc|
-|-|-|
void|N/A|表示无类型，用于函数返回值|
nullptr_t|指针宽度（32/64）|C++11 新增，表示空指针 nullptr 的类型|

#### 自动类型
type | byte |  desc|
-|-|-|
auto|N/A|C++11 新增，允许编译器自动推导变量类型|

#### 指针和引用类型
type | byte |  desc|
-|-|-|
指针类型|指针宽度（32/64）|例如 int* 表示整数指针|
引用类型|一般和指针类型相同|例如 int& 表示整数引用|

#### 基础类型的隐式转换
编译器自动进行的类型转换，不需要程序员进行任何操作。这些转换通常在类型兼容的情况下发生，比如从小的整数类型转换到大的整数类型。下面是经常遇到的隐式类型转换：
- 安全的隐式转换：

	- 整型提升：小的整型（如 char、short）会自动转换成较大的整型（如 int）
	- 算术转换：例如，当 int 和 double 混合运算时，int 会转换为 double

- 存在隐患的隐式转换：
	- 窄化转换：大的整数类型转换到小的整数类型，或者浮点数转换到整数，可能会造成数据丢失或截断
	- 指针转换：例如，将 void* 转换为具体类型的指针时，如果转换不正确，会导致未定义行为

### 其他类型
#### struct
结构体是不同类型数据的集合，允许将数据组织成有意义的组合.

#### 枚举
枚举是一种用户定义的类型，它可以赋予一组整数值具有更易读的别名

C++11 引入了新的枚举类型 作用域枚举:
```c++
enum class Color {
    RED,
    GREEN,
    BLUE
};

Color myColor = Color::RED; // 使用作用域解析运算符(::)访问枚举值
```

作用域枚举解决了传统枚举可能导致命名冲突的问题，并提供了更强的类型检查.

#### union
联合体允许在相同的内存位置存储不同类型的数据，但一次只能使用其一.

#### class
类是 C++ 的核心特性，是面向对象的基础，允许将数据和操作这些数据的函数封装为一个对象.

```c++
// polardbx-engine/extra/IS/src/client/client_service.h
class ClientService {
 public:
  ClientService() = default; // 用default修饰, 使之成为缺省构造函数
  virtual ~ClientService() = default;

  const std::string &get(const std::string &key) { return map_[key]; }
  void set(const std::string &key, const std::string &val);
  const std::string set(const char *strKeyVal, uint64_t len);
  int serviceProcess(easy_request_t *r, void *args);

 protected:
  std::map<const std::string, const std::string> map_;

 private:
  ClientService(const ClientService &other);  // copy constructor. 拷贝构造函数是一种特殊的构造函数，它在创建对象时，是使用同一类中之前创建的对象来初始化新创建的对象
  const ClientService &operator=(
      const ClientService &other);            // assignment operator. 赋值运算符

};
```

#### 字符串
c++11支持原始字符串字面值(raw string literals),类似于golang中的<code>``</code>, 它适用于太多字符需要转义的场景, 比如`cout<<R"(c:\fiiles\)";`

```c++
#include<string>
#include<cstring>
using namespace std;

string str("abc"); // <=> string noun = "abc";
char *cstr;
cstr = new char [str.size()+1];
// char cstr[str.size()+1]={};
strncpy(cstr, str.c_str(),str.size()); // string -> cstring, 因为c_str()的返回类型是`const char*`
delete[] cstr; // 销毁动态数组, []不能忘, 表明是销毁动态数组 , 否则仅销毁一个char的内存空间.

string s1, s2;
cin >> s1; // c++实现会删除输入两边的所有空白符
cin >> s2;

string line;
getline(cin,line); // 历史原因, string类的getline是独立的函数
line.length();
line[i]; // 不检查非法索引, 可能引发无法预料的结果, 因此需要与`line.length()`联合使用
line.at(2) = 'x'; // 将line的第3个字符改为'x'

stoi("35"); // 35
stol("35"); // 35
stof("3.5"); // 3.5
stod("3.5"); // 3.5
to_string(3.5); // "3.500000"
```

#### 数组
C++ 的数组是一个固定大小的序列容器，它可以存储特定类型的元素的集合。数组中的元素在内存中连续存储，这允许快速的随机访问，即可以直接通过索引访问任何元素，而无需遍历数组.

注意:
1. 数组在声明后（无论静态声明还是动态声明），数组的大小即固定，不可更改
1. 数组不提供任何内置的方法来获取其大小，通常需要额外保存数组的大小，或者使用特殊标记结束元素（C 风格的字符串使用 '\0' 表示数组结束）
1. 数组不提供边界检查，越界访问的代码是可以通过编译的（静态数组编译器会给出警告），可能导致很多潜在问题

```c++
#include<iostream>
using namespace std;

/*
C++ 不会对形式参数执行边界检查 所有函数而言，数组的长度是无关紧要的
*/
int sum(int array[],int n) //传入数组首地址和长度,  推荐
// int sum(const int array[],int n) // 加const防止函数中array的内容被修改, 外层函数一旦使用了const修饰数组那么内层的函数也应该使用const修改该数组.
{
	int value=0;
	for(int i=0;i<n;i++)
	{
		value+=array[i];
	}
	return value;
}

int sum2(int array[3]) // 指定数组长度
{
	int value=0;
	for(int i=0;i<3;i++)
	{
		value+=array[i];
	}
	return value;
}

int sum3(int *p,int n) // 传入数组的首地址，即指针, 再使用指针遍历数组
{
	int value=0;
	for(int i=0;i<n;i++)
	{
		value+=*(p+i);
	}
	return value;
}

int main()
{
    int array[3]={1,2,3};
	cout<<sum(array,3)<<endl;
	cout<<sum2(array)<<endl;
    cout<<sum3(array,3)<<endl;

    return 0;
}
```

二位数组初始化:
```c++
// 完全初始化
int matrix[2][3] = {
    {1, 2, 3},
    {4, 5, 6}
};

// 部分初始化
int matrix[2][3] = {
    {1, 2}, // 第一行的最后一个元素将被初始化为 0
    {4}     // 第二行的第二个和第三个元素将被初始化为 0
};

// 单行初始化
int matrix[2][3] = {1, 2, 3}; // 只初始化第一行，其他行将默认初始化为0

// 自动推断，和一维数组一样，编译器会根据数组推断二维数组第一维的大小
int matrix[][3] = {
    {1, 2, 3},
    {4, 5, 6}
};
```

多维数组:
```c++
#include<iostream>
#include <typeinfo> // [typeid详解](http://www.cppblog.com/smagle/archive/2010/05/14/115286.aspx), 它的奇怪输出与编译器实现有关.
using namespace std;

 int sum(int array[][2],int row);//第二维必须指定并传入行数, 推荐
int sum(int array[][2],int row)
// int sum(const int array[][2],int row) // 也可使用const来修饰多维数组
{
	int value=0;
	for (int i=0;i<row;i++)
	{
		for(int j=0;j<2;j++)
		{
			value+=array[i][j];
		}
	}
	return value;
}

int sum2(int array[2][2]);//指定数组的行列
int sum2(int array[2][2])
{
	int value=0;
	for (int i=0;i<2;i++)
	{
		for(int j=0;j<2;j++)
		{
			value+=array[i][j];
		}
	}
	return value;
}

int sum3(int* p,int row,int col);//数组首地址和行列数
int sum3(int* p,int row,int col)
{
	int value=0;
	for (int i=0;i<row;i++)
	{
		for(int j=0;j<col;j++)
		{
			value+=*(p+col*i+j);
		}
	}
	return value;
}

int sum4(int (*array)[2],int row);//传入行指针和行数
int sum4(int (*array)[2],int row)
{
	int value=0;
	for (int i=0;i<2*row;i++)
	{
		value+=(*array)[i];
	}
	return value;
}

int main()
{
	int array[2][2]={{1,2},{3,4}};
    cout<<"array的类型是"<<typeid(array)<<endl;
	cout<<sum(array,2)<<endl;
	cout<<sum2(array)<<endl;
	 cout<<sum3(&array[0][0],2,2)<<endl;
    cout<<sum3(array[0],2,2)<<endl; // array[0]降维后sum3传参可参考一维数组的传参`array=&array[0]`
    cout<<sum4(&array[0],2)<<endl;

    return 0;
}
```

动态多维数组
```c++
const int M = 10, N = 5;//10行5列

//动态开辟空间
int ** a = new int *[M];
for(int i = 0; i < M; i ++)
{
	a[i] = new int[N];
}

// 释放开辟的资源
for(int i = 0; i < M; i ++)
{
	delete [] a[i];
}
delete []a;

vector<vector<int> > p(M,vector<int>(N));

for(int i = 0 ; i < M ; i++)
{
	for(int j = 0; j < N; j++)
	{
			p[i][j]=i+j;
	}
}
```

##### [向量](https://en.cppreference.com/w/cpp/container/vector)
```c++
#include <vector>
using namespace std;

vector<int> v = {1}; // vector<Base_Type>是一个模板类, 最终为Base_Type类型的向量生成一个类, 类名为`vector<int>`
vector<int> v2(10); // 会用零值初始化v2, 且v.size()==10.
v.size(); // 返回unsigned int
v.capacity(); // 容量, 向量当前实际分配了内存的元素的个数
v.push_back(2);
for (auto i : v)
{
	...
}
// 向量的size和capacity可参考golang的slice概念.
v.reserve(v.size()+10); // 修改capacity的值, 但是没有给这段内存进行初始化. reserve方法对于vector的大小(即size)没有任何影响.
// resize函数重新分配大小，改变容器的size:
// - 当n小于当前size()值时候，vector首先会减少size()值 保存前n个元素，然后将超出n的元素删除(remove and destroy)
// - 当n大于当前size()值时候，vector会插入相应数量的元素 使得size()值达到n，并对这些元素进行初始化, 指定val时，vector会用val来初始化这些新插入的元素
// - 当n大于capacity()值的时候，会自动分配重新分配内存存储空间
v.resize(20);
```

##### 数组的替代
数组本身是一种常见的 C++ 数据类型，使用范围很广，但是本身也存在局限性。因此为了提升开发效率，C++ 标准库中提供了更加灵活的数据容器供开发者使用：
- std::vector: 可变大小的数组。提供对元素的快速随机访问，并能高效地在尾部添加和删除元素
- std::list 双向链表。支持在任何位置快速插入和删除元素，但不支持快速随机访问
- std::deque: 双端队列。类似于 std::vector，但提供在头部和尾部快速添加和删除元素的能力
- std::array (C++11): 固定大小的数组。提供对元素的快速随机访问，并且其大小在编译时确定
- std::forward_list (C++11): 单向链表。提供在任何位置快速插入和删除元素，但不支持快速随机访问
- std::stack: 栈容器适配器。提供后进先出 (LIFO) 的数据结构
- std::queue: 队列容器适配器。提供先进先出 (FIFO) 的数据结构
- std::priority_queue: 优先队列容器适配器。元素按优先级出队，通常使用堆数据结构实现
- std::set: 一个包含排序唯一元素的集合。基于红黑树实现
- std::multiset: 一个包含排序元素的集合，元素可以重复。基于红黑树实现
- std::unordered_set (C++11): 一个包含唯一元素的集合，但不排序。基于散列函数实现
- std::unordered_multiset (C++11): 一个包含元素的集合，元素可以重复，但不排序。基于散列函数实现

## 函数
函数的主要目的是使代码更模块化、更易于管理，并且可以重用.

C++ 的完整函数定义包括以下几个要素:
1. 返回类型：函数可能返回的值的数据类型。如果函数不返回任何值，则使用关键字 void
1. 函数名：用于识别函数的唯一名称
1. 参数列表：括号内的变量列表，用于从调用者那里接收输入值。如果函数不接受任何参数，则参数列表为空
1. 函数体：花括号 {} 内包含的代码块，当函数被调用时将执行这些代码

部分数值处理的预定义函数:
![部分数值处理的预定义函数](/misc/img/c/截图_2020-01-08_18-00-46.png)

c++要求函数调用前要有具体的函数实现或存在函数声明. 一般函数声明都定义在头文件中，包含该头文件即可调用相关函数.

C++ 中函数的参数传递方式包含：
1. 值传递
1. 指针传递

	本质其实是指针类型的值传递
1. 引用传递（特指左值）

	传递方式中形参成为实参的别名（引用），所以任何对形参的操作实际上都是在实参上进行的。引用就是别名.
	编译器在底层可能会使用指针来实现引用，但会提供更严格的语义和更简单的语法.

	引用声明:
	```c++
	int a = 5；
	int & ra = a;  // ra的类型是 int&（引用），必须声明时立即初始化

	int b = 6;
	ra = b;      // 非法，引用变量不支持重新赋值
	```

	使用:
	```c++
	void swap(int & a, int & b) {
	    int tmp = a;
	    a = b;
	    b = tmp;
	}

	int x = 5;
	int y = 7;
	swap(x, y);      // x = 7 y = 5
	```

1. 右值传递

	右值传递同引用传递类似，是传递右值引用到函数内部的传递方式. 它主要被用来实现移动语义和完美转发.

	左值（lvalue）：
	左值是指表达式结束后依然存在的持久性对象，可以出现在赋值语句的左边。
	左值可以被取地址，即可以通过取地址运算符 & 获取其地址。
	通常，变量、数组元素、引用、返回左值引用的函数等都是左值。

	右值（rvalue）：
	右值是指表达式结束后不再存在的临时对象，不能出现在赋值语句的左边。
	右值不能被取地址，即不能通过取地址运算符 & 获取其地址。
	通常，字面量、临时对象、返回右值引用的函数等都是右值。
	
	C++11 引入了右值引用（rvalue reference）的概念，通过 && 来声明一个右值引用。右值引用可以绑定到临时对象，从而支持移动语义和完美转发。移动语义允许将资源（如动态分配的内存）从一个对象 “移动” 到另一个对象，而不是进行昂贵的复制操作。完美转发允许将参数以原样传递给其他函数，避免不必要的拷贝。

	总的来说，C++11 中的左值和右值概念更加严格和明确，为语言引入了更多的灵活性和性能优化的可能性.

参数修改保护: 对于使用指针传递方式和引用方式传递参数的函数，因为函数内部有修改外部变量数据的能力，因此使用不当可能出现问题。此时可以通过 const 关键字来修饰函数的参数，达到禁止函数修改参数的目的
```c++
// 下面指针传递示例
void printInfo(int arr[]， int size); // 内部可修改arr

void printInfo(const int arr[], int size); // 内部不可修改arr

// 下面是引用传递示例
void printInfo(std::string& info); // 内部可以修改info

void printInfo(const std::string& info); // 内部不可以修改info
```

传参方式选择原则: 在函数定义时，选择合适的参数传递方式对于代码的性能和可读性至关重要.

对于仅使用参数的值，并不会进行修改的函数而言，应尽量遵循下面的原则：
- 如果数据对象很小，如内置数据类型或者小型结构，这按值传递
- 如果数据对象是数组，这使用指针，因为这是唯一的选择，并将指针声明为常量指针
- 如果数据对象较大的结构，则使用常量指针或者 const 引用，可以节省复制结构所需要的时间和空间，提高程序的效率
- 如果数据对象是类对象，则使用 const 引用。类设计的语义常常要求使用引用。这是 C++ 增加引用的主要原因。因此传递类对参数的标准方式是按引用传递

而对于需要通过参数修改原来变量值的函数，应遵循下面的原则：
- 如果数据对象是内置数据类型，则使用指针
- 如果数据对象是数组，则只能使用指针
- 如果数据对象是结构，则使用指针或者引用
- 如果数据对象是类对象，则使用引用

C++ 引用传递:
形参相当于是**实参的"别名"**，对形参的操作其实就是对实参的操作，在引用传递过程中，被调函数的形式参数虽然也作为局部变量在栈中开辟了内存空间，但是这时存放的是由主调函数放进来的实参变量的地址. 被调函数对形参的任何操作都被处理成间接寻址，即通过栈中存放的地址访问主调函数中的实参变量. 正因为如此，被调函数对形参做的任何操作都影响了主调函数中的实参变量.

> 传引用效率高于传值, 因为传值参数是局部变量, 所以调用函数时会存在实参的两个拷贝,  而引用参数只是占位符, 会被实参取代, 所以只存在实参的一个拷贝. 特别是class作为参数时更明显.

引用传递和指针传递是不同的，虽然它们都是在被调函数栈空间上的一个局部变量，但是任何对于引用参数的处理都会通过一个间接寻址的方式操作到主调函数中的相关变量. 而对于指针传递的参数，如果改变被调函数中的指针地址，它将影响不到主调函数的相关变量.

操纵元(manipulator)是以非传统方式调用的函数, 位于插入操作符`<<`之后, 调用它后, 其本身又会调用一个成员函数, 比如`endl, setw,setprecision`.

流作为函数参数使用时, 必须**传引用**.

每个输入流中都有`get()`, 它会读取一个字符, 与提取运算符(`>>`)不同, **无论下个输入字符是什么, 它都会读取**.
输入流的`putback()`是将字符放回到输入流而不是输入文件, 因此原始输入文件的内容不变.

>`cin.getline(str, max_len+1);`会保证str必定以`\0`结尾.

```c++
double next = 0;
while (inStream >> next) // >>即是action也是bool条件, 表示是否满足执行action的条件.
{
    ...
}
```

### 函数重载(**个人不推荐使用**)
函数名相同,  参数列表(个数,类型,顺序)不同. 编译器通过查看函数的参数列表（也称之为函数签名）来区分重载的函数.

```c++
// 下面是一组重载函数，同样是计算两个数的和，针对不同类型提供了不同的定义
int add(int a, int b) {        // 版本1
    return a + b;
}      

float add(float a, float b) {    // 版本2
    return a + b;
}    

double add(double a, double b) {  // 版本3
    return a + b;
}  

add(1, 2);      // 匹配版本1
add(1.0f, 2.0f);   // 匹配版本2
add(1.0, 2.0);    // 匹配版本3

add(1.0f, 2.0)    // 匹配 ？？？（匹配版本3，原因可以搜索 ”重载解析“）
```

注意:
1. 返回值类型与函数重载无关
1. 调用函数时, 实参的隐式类型转换可能产生二义性.
1. 一个参数是否使用引用并不能作为签名不同的依据. 参数是否是 const 能作为不同的依据

本质: 采用了name mangling(也叫name decoration)的技术, 即c++编译器默认会对符号名(比如函数名)进行改编,修饰, 简单理解就是重命名, 且不同编译器有不同的重命名规则.

> 可通过编译时选禁止优化和不生成调试信息, 再通过反汇编验证.

函数重载遵循下面的原则：
1. 函数签名必须不同
1. 作用域必须相同：重载的函数必须处于同一个作用域，否则它们被视为不同作用域中的不相关函数。
1. 最佳实践是尽量保持重载函数的明确性，避免产生容易混淆的重载集合。

### 默认参数(**个人不推荐使用**)
注意:
1. 默认参数只能放在函数声明处**或者**定义处，能放在声明处就放在声明处
1. 如果某个参数是默认参数，那么它后面的参数必须都是默认参数

    因为非默认参数的参数必须要给出具体值，而调用函数传递参数的时候是从左到右的，所以非默认参数前面的都必须要传值进来. 那么默认参数后面的当然也得都为默认参数.

1. 不要重载一个带默认参数的函数

### 回调
函数名本身就是函数的指针。函数指针在定义时必须指明所指向函数的类型，包括返回类型和参数列表.

函数指针的定义语法：
```c++
// 返回类型 (*指针变量名)(参数列表);

// 示例
int add(int a, int b) {
    return a + b;
}

int (*pf)(int, int) = add;  // 可以这么理解定义：因为(*pf)表示函数，那么pf就是函数的指针

// 类似数组，函数指针也有两种使用方式
cout << (*pf)(2, 3) << endl;    // 5 指针使用方式
cout << pf(2, 3) << endl;    // 5 直接作为函数名使用

// 函数指针的定义一般都不怎么直接，使用也不方面

// 经典C++中可以使用typedef简化这个定义
typedef int (*p_fun)(int, int);    // 现在p_fun就是一种类型名称
p_fun pAdd = add;          // 精简很多

// 现代C++提供了 using 语法让这个过程更加直观，推荐使用
using p_fun = int (*)(int, int);   // 可读性更强
p_fun pAdd = add;

// auto大杀器
auto pAdd = add;          // 懒人利器
```

使用:
```c++
#include <iostream>

void callBack(int costTimeMs);
void work(void (*pf)(int));

int main() {
    work(callBack);
}

void callBack(int costTimeMs) {
    using namespace std;

    cout << "costTime:" << costTimeMs << endl; 
}

void work(void (*pf)(int)) {
    std::cout << "do some work" << std::endl;
    // ...
    pf(123);  // (*pf)(123) 也ok
}
```

`std::function`容器. 它是一个泛型函数封装器，其实例可以存储、复制和调用任何可调用对象，如普通函数、Lambda 表达式、函数对象（functors）以及其他函数指针:
```c++
#include <functional>
#include <iostream>

using namespace std;

void callBack(int costTimeMs) {
    cout << "costTime:" << costTimeMs << endl; 
}

void work(function<void(int)> callBack) {
    callBack(1234);
}

int main() {
    function<void(int)> func = callBack;
    work(func);
    return 0;
}
```

```c++
// 封装函数
void printHello() {
    std::cout << "Hello, World!" << std::endl;
}
std::function<void()> func = printHello;

// 封装Lambda表达式
std::function<int(int, int)> add = [](int a, int b) -> int {
    return a + b;
};

int sum = add(2, 3); // sum 的值为 5

// 封装成员函数
class MyClass {
public:
    void memberFunction() const {
        std::cout << "Member function called." << std::endl;
    }
};

MyClass obj;
std::function<void(const MyClass &)> f = &MyClass::memberFunction;
f(obj); // 输出: Member function called.

// 封装带有绑定参数的函数
void printSum(int a, int b) {
    std::cout << "Sum: " << a + b << std::endl;
}

int main() {
    using namespace std::placeholders; // 对于 _1, _2, _3...

    // 绑定第二个参数为 10，并将第一个参数留作后面指定
    std::function<void(int)> func = std::bind(printSum, _1, 10);
    func(5); // 输出: Sum: 15
    return 0;
}
```

## class
C++ 通过引入类支持了面向对象编程（OOP）。在 C++ 中，类是创建自定义数据类型的核心概念之一。类用于定义与特定类型相关的数据（成员变量）及操作这些数据的函数（成员函数）。通过类，可以实现面向对象编程（OOP）的基本原则，如封装、继承和多态

定义类:
```c++
class MyClass {
public:
    // 公共成员，通常的对外提供的方法定义
    void setMember(int member);
private:
    // 私有成员，成员变量，仅供内部调用函数
    int mMember;    // 集团规范推荐,使用m前缀
    void innerFunc();  // 函数一律小驼峰
protected:
    // 受保护成员，成员变量，供子类调用函数
};
```

访问控制
类成员的访问权限可以是 public、private 或 protected：
1. Public（公共）：公共成员可以在类的外部被访问
1. Private（私有）：私有成员只能在类的内部被访问
1. Protected（受保护）：受保护成员可以在类的内部以及其派生类中被访问

构造函数是用于初始化class对象的成员函数,声明类的对象时会自动调用构造函数, 要求:
1. 构造函数的名称与类的名称相同
1. 构造函数没有任何类型的返回值

构造函数可以被重载，以提供不同的初始化方式。成员初始化列表提供了初始化成员变量的一种更高效的方式，对于类中的常量成员、引用成员来说，成员初始化列表是必须的:
```c++
class MyClass {
public:
    MyClass(int m1, int m2, int m3) : mM1(m1), mM2(m2), mM3(m3) {}
private:
    int mM1;
    const int mM2；
    int & mM3;
};

// 类的初始化方式
MyClass a1(1, 2, 3);         // 传统构造函数
MyClass a1 = MyClass(1, 2, 3);   // 同上
MyClass a2 = {1, 2, 3};     // 列表初始化，会匹配最合适的构造函数
MyClass a3{1, 2, 3};        // 同上
```

默认构造函数: 没有参数列表或参数列表均有默认参数的构造函数. class没有构造函数时, 编译器会自动创建一个参数列表为空且什么事都不做的默认构造函数.

构造函数的初始化区域: 构造函数后紧跟的以`:`开始的部分,可初始化部分或全部成员变量. **它能初始化const修饰的成员变量**. **它初始化变量的顺序与class中变量声明的顺序有关,  与自身顺序无关, 因此顺序不一致可能导致未知结果**.

使用数据类型时访问不了值和操作的实现细节, 该数据类型就称为抽象数据类型(abstract data type, ADT). 它确保了类的接口和实现的分离, 是良好的编程实践.

> adt的类定义在头文件中(该文件称为接口文件), 实现则在同名(名称不一定相同但推荐相同)的实现文件中(扩展名通常是cpp).

使用ADT类需遵循以下规则:
1. 所有成员变量都是私有成员
1. 每个基本操作都设为类的公共成员函数
1. 任何辅助函数都是私有成员函数

> ADT类的公共成员函数, 友元函数, 普通函数或重载操作符

成员函数和非成员函数(比如友元函数)的选择:
1. 函数的执行只涉及一个对象, 就用成员函数.
1. 要执行的函数涉及多个对象, 就使用非成员函数.

**拷贝构造函数, 操作符=,析构函数统称为BigThree, 如果需要定义其中一个, 就必须定义全部, 否则编译器会自动创建, 可能导致未预料的结果.**
在默认情况下(用户没有定义，但是也没有显式的删除), 编译器会自动的隐式生成一个拷贝构造函数和赋值运算符. 但用户可以使用delete来指定不生成拷贝构造函数和赋值运算符，这样的对象就不能通过值传递，也不能进行赋值运算.

析构函数是类的一个特殊成员函数，它在类的对象生命周期结束时自动被调用以执行清理工作。主要用途是释放对象占用的资源，并执行一些必要的清理操作，例如释放动态分配的内存、关闭文件和数据库连接等.

自动调用析构函数的情况：
1. 局部对象：当局部对象的作用域结束时，例如函数结束时，其中的局部对象会被销毁，调用析构函数
2. 动态分配的对象：当使用 delete 操作符删除一个动态分配的对象时，析构函数会被调用
3. 静态和全局对象：当程序结束时，所有的静态和全局对象会被销毁，调用析构函数
4. 临时对象：当临时对象的生命周期结束时，例如临时对象作为函数参数传递，或者在它们创建的表达式结束后，析构函数会被调用
5. 通过 std::unique_ptr 或 std::shared_ptr 管理的对象：当智能指针销毁或被重新赋值，造成引用计数降为零时，析构函数会被调用

```c++
class Person
{
public:
    Person(const Person& p) = delete;
    Person& operator=(const Person& p) = delete;

private:
    int age;
    string name;
};
```

```c++
// 定义类时, 通常将所有成员变量设为私有, 再通过成员函数去访问.
// 通常会在赋值函数(变更类的私有成员的成员函数)前用缀`set`修饰
class Day
{
   public:
   friend bool equalYear(const Day& d1, Day d2); // 友元函数
   friend Day operator +(const Day&  d1, const Day&  d2);
   friend ostream& operator <<(ostream& outs, const Day& d);
   Day& operator =(const Day& d)//赋值运算符重载, 必须是成员函数
    {
        cout << "operator ====" << endl;
        if (this != &d) // 因为d是别名所以要取地址
        {
            this->year = d.year;
			this->month = d.month;
			this->day = d.day;
        }
        return *this;
    }
	// void operator =(const Day& d)//赋值运算符重载, 这种形式也可以 
    // {
    //     cout << "operator ====" << endl;
    //     if (this != &d)
    //     {
    //         this->year = d.year;
	// 		this->month = d.month;
	// 		this->day = d.day;
    //     }
    //     // return *this;
    // }
     Day(){} // 默认构造函数
	 Day(int year, int month, int day){
		 // year = year; // 错误写法. 构造函数参数和类的数据成员同名时, 必须使用`this`指明或用构造函数的初始化区域
		  this->year = year;
		  this->month = month;
		  this->day = day;
	 }
	  Day(int year):year(year+10){
	 }
	Day(int month,int day):Day(0,month,day){ // c++11支持构造函数委托, 即允许一个构造函数调用另一个构造函数.
	 }
      int month;
      int day;

     void output();//仅有函数声明
	  void printYear()
           {
                 cout<< "year = " << year << endl;
           }

	// 建议: 如果一个成员函数只在其他成员函数中作为辅助函数使用, 那么将其设为私有
	 private:// 仅成员函数中可见, 在其他地方都无法访问
        int year=10; // c++11支持成员初始化, 即为成员变量设置默认值, 可用带参数的构造函数覆盖.
		void print();
}; // 不要漏掉这个分号

// `::`是作用域解析操作符, 作用于类名
void Day::output() // 成员函数定义
{
	cout << "month = "<< month
		<<", day = " << day << endl;

		print();
}

void Day::print() // 成员函数定义
{
	cout << "year = "<< year
	<< ", month = "<< month
		<<", day = " << day << endl;
}

bool equalYear(const Day& d1, Day d2){
	// d1.month=2; // error: assignment of member ‘Day::month’ in read-only object

	return d1.year==d2.year;
}

Day operator +(const Day& d1,const Day& d2)
{
	return Day(d1.year+d2.year,d1.month+d2.month,d1.day+d2.day);
}

ostream& operator <<(ostream& outs, const Day& d){
	outs << "day:" << d.year << '-';
	outs <<  d.month << '-';
	outs <<  d.day << endl;

	return outs;
}

int main(){
Day day(2020,1,2);// 显示调用构造函数
Day day2; // 会调用默认构造函数

Day day3(2020);
day3.output(); // year = 2030

  cout<<equalYear(day,day3)<<endl;

  Day day4= day+day3;
day4.output();

	Day day5 = day+2020; // 编译器根据合适的构造函数, 自动完成类型转换
day5.output(); // year = 2030
// cout << "I have " << amount << " in my purse.\n";
// means the same as
// ((cout << "I have ") << amount) << " in my purse.\n";

cout<<day5<<endl;

Day day6; // Day day6=day; 不是赋值操作, 是通过拷贝构造函数进行的初始化, 因此`=`重载不执行
day6= day;
cout<<day6<<endl;

    return 0;
}
```

在 C++ 中，通常应用 “资源获取即初始化”（RAII）原则来管理资源。RAII 建议在构造函数中获取资源，并在析构函数中释放资源。这样，资源的生命周期就与包含它的对象的生命周期绑定在一起，简化了资源管理并防止了资源泄漏。

当正确使用 RAII 原则时，通常不需要手动调用析构函数，因为 C++ 会确保在对象生命周期结束时自动调用析构函数。然而，如果你使用 “裸” 指针手动管理资源，就必须非常小心地确保每个分配的资源最终都被释放，否则可能会导致资源泄漏。智能指针（如 std::unique_ptr 和 std::shared_ptr）是现代 C++ 推荐的资源管理方式，它们可以自动管理资源的生命周期，从而避免直接手动管理资源的复杂性和危险.

类可以重载各种运算符，以提供类似于内建类型的行为:
```c++
class MyClass {
public:
    MyClass() : data(new int[10]) { }   // 构造函数
    ~MyClass() { delete[] data; }     // 析构函数

    // 拷贝赋值运算符
    MyClass & operator=(const MyClass& other) {
        if (this != &other) { // 避免自赋值
            std::copy(other.data, other.data + 10, data);
        }
        return *this;
    }

private:
    int* data;
};

// 使用
MyClass a;
MyClass b = a;    // 默认的赋值操作是浅拷贝，这里因为重载了 = 运算符，变成深拷贝

// C++11开始可以删除默认的赋值操作符，从而防止因浅拷贝带来的风险
class MyClass2 {
    // ...
    MyClass2 & operator=(const MyClass2 & other) = delete; // 禁用赋值操作符
    // ...
};

MyClass2 a;
MyClass2 b = a; // 非法，MyClass2的 = 运算符被禁用
```

运算符重载注意:
1. 运算符重载并不改变运算符的优先级、结合性或操作数个数. 这些都是由语言规范定义的。
1. 不要滥用运算符重载。重载的运算符应该和它的原始意图保持相关性，否则可能导致代码难以阅读和理解。
1. 记得检查自赋值。特别是在重载赋值运算符时（如 operator=），要确保它能正确处理自赋值的情况。
1. 为了保持一致性，考虑重载对应的复合赋值运算符。例如，如果你重载了 operator+，那么也应该重载 operator+=
1. 当重载某些运算符，如 ==，通常也需要重载相应的运算符，如！=，以确保逻辑一致性。
1. 某些运算符最好重载为非成员函数。像 <<和>> 这类运算符，如果要用于输入输出流的话，通常作为非成员函数重载比较合适，因为它们的左操作数通常是流对象。


拷贝构造函数和拷贝赋值运算符
对象的赋值操作是常见的操作，应该尽量避免使用浅拷贝，因为这种方式存在潜在风向。为解决这个问题类可以定义专门的拷贝构造函数和拷贝赋值运算符，以控制对象如何被复制
```c++
#include <iostream>

class MyClass {
public:
    MyClass() : data(new int[10]) { } // 默认构造函数

    ~MyClass() { delete[] data; } // 析构函数

    // 拷贝构造函数
    MyClass(const MyClass & other) : data(new int[10]) {
        std::copy(other.data, other.data + 10, data);
        std::cout << "copy init" << std::endl;
    }

    // 拷贝赋值运算符
    MyClass & operator=(const MyClass & other) {
        if (this != &other) { // 避免自赋值
            std::copy(other.data, other.data + 10, data);
        }
        std::cout << "copy =" << std::endl;
        return *this;
    }

private:
    int* data;
};


int main() {
    MyClass a;
    MyClass b;
    MyClass c = a;
    c = b;
    return 0;
}

// 程序输出
// copy init
// copy =
```

移动构造函数和移动赋值运算符（C++11）
在 C++11 中引入了移动语义，允许从临时对象 “移动” 资源，而不是复制它们
```c++
#include <iostream>

using namespace std;

class BigMemoryPool {
    private:
        static const int POOL_SIZE = 4096;
        int* mPool;

    public:
        BigMemoryPool() : mPool(new int[POOL_SIZE]{0}) {
            cout << "call default init" << endl;
        }

        // 编译器会优化移动构造函数，正常情况可能不会被执行
        // 可以添加编译选项 “-fno-elide-constructors” 关闭优化来观察效果
        BigMemoryPool(BigMemoryPool && other) noexcept {
            mPool = other.mPool;
            other.mPool = nullptr;
            cout << "call move init" << endl;
        }

        BigMemoryPool & operator=(BigMemoryPool && other) noexcept {
            if (this != &other) {
                this->mPool = other.mPool;
                other.mPool = nullptr;
            }
            cout << "call op move" << endl;
            return *this;
        }
        void showPoolAddr() {
            cout << "pool addr:" << &(mPool[0]) << endl;
        }

        ~BigMemoryPool() {
            cout << "call destructor" << endl;
        }
};

BigMemoryPool makeBigMemoryPool() {
    BigMemoryPool x;  // 调用默认构造函数
    x.showPoolAddr();
    return x;         // 返回临时变量，属于右值
}

int main() {
    BigMemoryPool a(makeBigMemoryPool());  
    a.showPoolAddr();
    a = makeBigMemoryPool();  
    a.showPoolAddr();
    return 0;
}

// 输出内容
call default init
pool addr:0x152009600
instance addr:0x16fdfeda0
pool addr:0x152009600
instance addr:0x16fdfeda0  // 编译器优化，这里a和x其实是同一个实例，因此不会触发移动构造
call default init
pool addr:0x15200e600    // 新的临时变量，堆内存重新分配
instance addr:0x16fdfed88  // 临时变量对象地址
call op move        // 移动赋值
call destructor
pool addr:0x15200e600    // a的Pool指向的内存地址变成新临时对象分配的地址，完成转移
instance addr:0x16fdfeda0  // a对象的地址没有变化
call destructor
```

C++11 引入移动语义之前，类似的做法需要返回指针或者通过拷贝的方式来保存临时对象，前者会引入资源管理问题后者会有拷贝的性能损耗.

### 继承
继承是基于一个基类创建新类(派生类)的过程. 派生类自动拥有基类的所有成员变量和函数, 并可根据需要添加更多的成员函数/变量, 但**私有成员函数/构造函数/析构函数/赋值操作符=均不会被继承**.

C++ 继承方式有三种：
1. 公有继承（public）最常见的继承类型。在公有继承中，基类的公有成员和保护成员在派生类中保持其原有的访问级别，而基类的私有成员在派生类中是不可访问的
1. 保护继承（protected）基类的公有成员和保护成员都成为派生类的保护成员。这意味着它们只能被派生类或其进一步的派生类中的成员函数访问
1. 私有继承（private）私有继承会将基类的公有成员和保护成员都变成派生类的私有成员。这意味着这些成员只能被派生类的成员函数访问，而不能被派生类的派生类访问

**C++ 是支持多重继承的，即可以从多个类派生一个类，但是通常建议谨慎使用，因为多重继承可能会引起一些复杂的问题**

在类继承的场景中，基类的析构函数一般要声明为虚析构函数，这样才能保证在通过基类指针删除对象时，派生类的资源也能被正确的释放.

重载(overloading)：要求两个函数必须在同一个作用域, 函数名相同，函数的参数个数、参数类型或参数顺序三者中必须至少有一种不同. 函数返回值的类型可以相同，也可以不相同.
重定义(redefining)：也叫做隐藏，**子类重新定义父类中有相同名称的非虚函数 ( 参数列表可以不同 ) **，指派生类的函数屏蔽了与其同名的基类函数, 可以理解成发生在继承中的重载.
重写(override)：也叫做覆盖，发生在子类和父类继承关系之间, 要求父类的函数有virtual修饰,  **子类重新定义父类中有相同名称和参数的虚函数**.

> 对于程序员不区分redefine和override, 但编译器需区分.

对除派生类以外的对象来说基类的protected成员就等同于private.
当派生类存在与基类同名的成员变量时候，派生类的成员会隐藏基类成员，但派生类中存在基类成员的拷贝，继承限定允许的情况下可显示地通过BASE::member来访问.


当一个子类从父类继承时，父类的所有成员成为子类的成员，此时对父类成员的访问状态由继承时使用的继承限定符决定:
1.如果子类从父类继承时使用的继承限定符是public，那么
(1)父类的public成员成为子类的public成员，允许类以外的代码访问这些成员；
(2)父类的private成员仍旧是父类的private成员，子类成员不可以访问这些成员；
(3)父类的protected成员成为子类的protected成员，只允许子类成员访问；


2.如果子类从父类继承时使用的继承限定符是protected，那么
(1)父类的public成员成为子类的protected成员，只允许子类成员访问；
(2)父类的private成员仍旧是父类的private成员，子类成员不可以访问这些成员；
(3)父类的protected成员成为子类的protected成员，只允许子类成员访问


3.如果子类从父类继承时使用的继承限定符是private，那么
(1)父类的public成员成为子类的private成员，只允许子类成员访问；
(2)父类的private成员仍旧是父类的private成员，子类成员不可以访问这些成员；
(3)父类的protected成员成为子类的private成员，只允许子类成员访问;

> 继承限定可用父类可见性需降级到指定级别来理解.

切割问题(slicing problem): c++允许将派生类对象赋值给基类型的变量, 此时该基类型变量无法访问到派生类的成员, 但通过使用指向动态对象实例的指针赋值可以解决该问题., 同时此时需要使用virtual成员函数来访问.

头文件:
```c++
    class Employee
    {
    public:
        Employee( );
        Employee(string theName, string theSSN);
        string getName( ) const; // 加const表示禁止修改成员变量: 编译器会自动给每一个函数加一个this指针, 该形式实际上，也就是对这个this指针加上了const修饰.
        string getSSN( ) const;
        double getNetPay( ) const;
        void setName(string newName);
        void setSSN(string newSSN);
        void setNetPay(double newNetPay);
        void printcheck( ) const;
    private:
        string name;
        string ssn;
        double netPay;
    };

	    class HourlyEmployee : public Employee
    {
    public:
        HourlyEmployee( );
        HourlyEmployee(string theName, string theSSN,
                           double theWageRate, double theHours);
        void setRate(double newWageRate);
        double getRate( ) const;
        void setHours(double hoursWorked);
        double getHours( ) const;
        void printCheck( ) ; // redefining
    private:
        double wageRate;
        double hours;
    };
```
实现:
```c++
// 派生类对象拥有基类的所有成员变量, 调用派生类的构造函数时, 需要为这些成员变量分配内存并初始化. 因此派生类定义构造函数时, 应始终包含对基类构造函数的调用.
HourlyEmployee::HourlyEmployee(string theName, string theNumber,
                                   double theWageRate, double theHours)
    : Employee(theName, theNumber), wageRate(theWageRate), hours(theHours) // 派生类的构造函数
    {
        //deliberately empty
    }

```

### 重载
重载的赋值操作符必须定义为**类的成员函数**.

```c++
Derived& Derived::operator =(const Derived& object){ // 派生类中的赋值运算符重载
	Base::operater = (object)
	...
}

Derived::Derived(const Derived& object) : 	Base(object){ // 派生类中的拷贝构造函数
	...
}
```

### 友元函数
友元函数不是成员函数, 它本质还是普通函数, 但被特别授予了访问类的数据成员的权限, 包括私有成员.

友元函数是定义在类外部的普通函数，它被某个类声明为其 “友元”。这意味着友元函数可以访问该类的所有成员，包括私有和受保护的成员。友元函数不是类成员函数，也不受类的封装性约束。

友元函数的声明方式是在类的定义内部使用关键字 friend，后跟函数的原型，友元函数实现时不能加类名作用域限定.

```c++
#include <iostream>

// 声明 Vector2D 类
class Vector2D {
private:
    float x_;
    float y_;

public:
    Vector2D(float x = 0.0f, float y = 0.0f) : x_(x), y_(y) {}

    // 友元函数声明，用于重载 + 操作符
    friend Vector2D operator+(const Vector2D & a, const Vector2D & b);

    // 输出 Vector2D 对象的友元函数
    friend std::ostream & operator<<(std::ostream & out, const Vector2D & v);
};

// 重载 + 操作符的友元函数定义
Vector2D operator+(const Vector2D & a, const Vector2D & b) {
    return Vector2D(a.x_ + b.x_, a.y_ + b.y_);
}

// 重载 << 操作符的友元函数定义，用于输出 Vector2D 对象
std::ostream & operator<<(std::ostream & out, const Vector2D & v) {
    out << "(" << v.x_ << ", " << v.y_ << ")";
    return out;
}

int main() {
    Vector2D vec1(1.0, 2.0);
    Vector2D vec2(3.0, 4.0);
    Vector2D vec3;

    vec3 = vec1 + vec2; // 使用友元函数重载的 + 操作符

    std::cout << "vec1: " << vec1 << std::endl;
    std::cout << "vec2: " << vec2 << std::endl;
    std::cout << "vec3: " << vec3 << std::endl; // 输出: vec3: (4, 6)

    return 0;
}
```

友元类是一个允许特定类访问另一个类的私有和受保护成员的机制。在 C++ 中，通常情况下，一个类无法访问另一个类的私有（private）和受保护（protected）成员，即使它们需要彼此协作。友元类提供了一种方式，可以指定某些类之间有更紧密的关系，并允许它们访问对方的非公共接口.

```c++
#include <iostream>

class MyClass; // 前向声明

// 声明一个类（FriendClass），该类将访问MyClass的私有和受保护成员
class FriendClass {
public:
    void accessMyClass(MyClass & obj);
};

// 声明主类（MyClass）
class MyClass {
private:
    int secret;

public:
    MyClass(int val) : secret(val) {}

    // 声明FriendClass为MyClass的友元类
    friend class FriendClass;
};

// FriendClass成员函数实现
void FriendClass::accessMyClass(MyClass & obj) {
    // 可以访问MyClass的私有成员'secret'
    std::cout << "MyClass secret value is: " << obj.secret << std::endl;
}

int main() {
    MyClass obj(42);       // 创建MyClass对象
    FriendClass friendObj; // 创建FriendClass对象

    friendObj.accessMyClass(obj); // 访问MyClass的私有成员
    return 0;
}
```

使用友元可能会破坏类的封装性和数据隐藏原则，因为它们允许外部函数或者类直接访问类的私有成员。因此，建议谨慎使用友元，只在确实需要时才使用，并寻找是否有其他设计替代方案。在设计类时，应尽可能通过公共成员函数或成员函数的重载来提供类的行为和操作，而将友元作为特定情况下的解决方案.

### 析构函数
析构函数是成员函数, 用前缀`~`修饰, 在类的对象离开作用域时自动调用, 用于释放动态内存.

派生类析构时会自动调用基类的析构函数, 调用顺序与构造函数相反.

析构函数最好是虚函数, 因为`delete pBaseClass`时会从派生类开始释放内存, 而不是仅释放基类变量的派生类内存, 避免内存泄露.

### 拷贝构造函数
要求一个参数, 且类型与class相同, 该参数必须是引用, **通常会加const修饰**, 该类构造函数即拷贝**构造函数**.

深拷贝为相关涉及内存的完全拷贝, 源变化不会影响当前的拷贝.
浅拷贝为当前拷贝与源在内存引用上有重叠, 源的改变可能会影响当前的拷贝.

### 虚函数
通过运行时确定一个过程的具体实现的技术叫动态绑定.
多态是指借助动态绑定技术为一个函数名关联多种含义的能力.

虚函数是某种意义上能在定义前使用的函数.  它是c++提供动态绑定的一种具体实现.

虚函数的作用是允许在派生类中重新定义与基类同名的函数，并且可以通过基类指针或引用来访问基类和派生类中的同名函数.

一旦某个函数被声明成虚函数，则所有派生类中它都是虚函数.

```c++
virtual double biill() const; //  const用在函数上，说明这个函数不能修改类的成员变量
```

### 多态
多态允许派生类重写基类的虚拟函数，使得通过基类引用或指针调用这些函数时可以执行派生类的版本

```c++
#include <iostream>

class Base {
public:
    void baseMethod() {
        std::cout << "Base method" << std::endl;
    }

    virtual void polymorphicMethod() {
        std::cout << "Base polymorphic method" << std::endl;
    }

    virtual ~Base() {} // 虚析构函数，用于多态
};

// 公有继承派生类
class Derived : public Base {
public:
    // 重写基类的虚函数
    void polymorphicMethod() override {
        Base::polymorphicMethod();  // 可以通过添加限定域调用基类实现
        std::cout << "Derived polymorphic method" << std::endl;
    }
};

int main() {
    Derived d;
    d.baseMethod();           // 调用基类的方法
    d.polymorphicMethod();    // 调用派生类重写的方法

    Base *b = &d;
    b->polymorphicMethod();   // 通过基类指针调用派生类的方法，体现多态
    return 0;
}
```

### 抽象类和纯虚函数
如果一个类包含至少一个纯虚函数（以 = 0 结尾），则该类被认为是抽象类，不能直接实例化，只包含纯虚函数而没有成员变量的抽象类和 Java 中的接口（Interface）功能类似. 同时这个方法必须在派生类(derived class)中被实现.

```c++
// Interface in C++
class IShape {
public:
    virtual void draw() const = 0; // 纯虚函数
    virtual ~IShape() {} // 虚析构函数以确保派生类的析构函数被调用
};

class Circle : public IShape {
public:
    void draw() const override {
        // 实现绘制圆形的代码
    }
};

class Rectangle : public IShape {
public:
    void draw() const override {
        // 实现绘制矩形的代码
    }
};
```

### 模板类
C++ 模板类是一种强大的特性，它允许程序员编写泛型且可重用的代码。模板类可以用来定义在编译时可以指定类型参数的类，这意味着可以用相同的基本代码来处理不同的数据类型。可以这么说现代 C++ 的很多功能强大的特性都和模板技术有关系下面是模板类的一般定义语法：

```c++
template <typename T>
class MyTemplateClass {
    const T& getValue();
public:
    T myValue;
};
```

C++ 中常见的模板类应用如下:
1. 容器类
	C++ 标准库中提供一系列的泛型容器，前面提到过的 vector、list、stack 都是模板类实现的。

	相关容器的用法可以搜索对应的文档。

1. 智能指针

	智能指针，同样是利用模板类技术实现的，它们提供了自动内存管理功能，可以帮助避免内存泄漏。

### func list
- endl

从定义中看出，endl是一个函数模板，它实例化之后变成一个模板函数，作用是输出一个换行符，并立即刷新缓冲区.

## 异常处理
try-throw-catch.

抛出异常后, throw语句所在的try块会停止执行, 然后开始执行catch块的代码. 执行catch块的过程被称为捕捉异常, 它不是函数调用.

```c++
  try
    {
        ....
        throw donuts;
		....
    }
    catch(int e)// e叫catch块参数
    {
        cout << e << " donuts, and No Milk!\n"
              << "Go buy some milk.\n";
    }
```

异常类的本质是类.

c++异常处理支持多个catch块, 它们的顺序十分重要. `catch (...)`表示默认块, 应在所有其他catch之后.
c++支持在函数中抛出异常再由外层try-catch捕获处理.

函数调用可能抛出异常时, 应使用异常规范(exception specificaftion),即在函数声明和函数定义中列出异常, 函数的异常规范也被称为throw 列表. 异常规范中的所有异常都必须处理.
未处理异常:函数抛出异常规范未列出的异常, 该异常会导致程序终止.
表明函数的任何异常都已在函数内部处理完成, 没有必要抛出时可用空白的异常规范`throw()`.

在派生类重定义或重写函数定义时, 不可在异常规范中添加新异常, 但允许删减基类原有的异常.

异常能不用就尽量不用, 切记滥用.

## 模板
模板的哲学在于**将一切能够在编译期处理的问题丢到编译期进行处理**，仅在运行时处理那些最核心的动态服务，进而大幅优化运行期的性能. 因此模板也被很多人视作 C++ 的黑魔法之一.

函数模板通过模板参数化来实现，在实例化时，编译器根据传递给模板的实际参数类型生成具体的函数实例.

> 编译器在编译时会根据调用的参数类型生成对应的实际函数，这个过程被称为模板的实例化.

函数模板的定义实际是多个函数定义的"合集". 编译器仅针对**用到的每种类型**生成单独的函数定义, 未用到的类型不生成定义.

```c++
template<typename T> // `template<typename T>`为模板前缀, T为类型参数. 函数模板支持多个类型参数, 但通常只需一个.
void swapValues(T& variable1, T& variable2)
{
    T temp;
	...
}
```

例子:
```c++
#include <iostream>
using namespace std;

// 函数原型
template <typename T>
T add(T a, T b);    

int main() {
    cout << add(1, 2) << endl;       // 3
    cout << add(1.0f, 2.1f) << endl;   // 3.1
    cout << add(1.0, 3.2) << endl;    // 4.2

    return 0;
}

template <typename T>
T add(T a, T b) {
    return  a + b;
}
```

重载的模板:
```c++
// 函数原型
template <typename T>
T add(T a, T b);

template <class T>      // 声明模板时 typename 和 class 等价
T add(T a, T b, T c);

// 函数定义略
```

> typename 和 class 在模板参数列表中没有区别, 在 typename 这个关键字出现之前，都是使用 class 来定义模板参数的. 但在模板中定义有嵌套依赖类型的变量时，需要用 typename 消除歧义, 因此推荐使用typename.

函数模板一般不使用其声明(即不提前声明函数), 而是在使用前include其实现文件.

模板还支持类模板, 语法:
```c++
// 在类定义和成员函数的定义前用`template<class Type_Parameter>`作为开头
template<class Type_Parameter>
class xxx
{
	...
}

template<class Type_Parameter>
void xxx<Type_Parameter>::set(int index, Type_Parameter value)
{
	...
}

xxx<int> t; // 为具体类名来特化类模板;
```

返回类型后置:
```c++
// 正常函数声明
int add(int a, int b);
// 返回类型后置声明
auto add(int a, int b) -> int;

// 利用该语法可以这么声明上面的函数（推荐C+11中使用）
template <typename T1, typename T2>
auto funcName(T1 x, T2 y) -> decltype(x + y) {
    ...
    return x + y;
}

// C++14及以后得标准拓展了auto的类型推导能力
auto funcName(T1 x, T2 y) {
    ...
    return x + y;
}
```

## 标准模板库(standard template library, STL)
STL包含了栈, 队列和其他许多标准数据结构的实现, 是c++标准的一部分.

STL中的类都是模板类. STL容器类普遍使用了迭代器, 迭代器是一种特殊对象, 简化了遍历容器中所有数据的过程.

迭代器(iterator)是指针的泛化形式, 它通常通过指针来实现, 因此可隐藏实现细节, 提供在所有容器类中都一致的迭代器接口.
每种迭代器只能用于它自己的容器类.

常量迭代器(constant iterator): 提领操作符生成的元素为只读.
可变迭代器(mutable iterator): 提领操作符生成的元素可读写.

STL容器类是各种用于容纳数据的数据结构, 比如队列, 列表和栈.
STL最简单的列表是双链表(doubly linked list).

容器配接器(container adapters)是在其他类基础上实现的模板类, 比如stack模板类默认就是在deque模板类基础上实现的.

泛型算法(generic algorithm)即STL模板函数.

## lambda
Lambda 表达式，实际上就是提供了一个类似匿名函数的特性.

```code
[捕获列表](参数列表) mutable(可选) 异常属性 -> 返回类型 {
// 函数体
}
```

捕获列表也分为以下几种：
1. 值捕获

    与参数传值类似，值捕获的前提是变量可以拷贝，不同之处则在于, **被捕获的变量在 Lambda 表达式被创建时拷贝， 而非调用时才拷贝**
1.  引用捕获

    与引用传参类似，引用捕获保存的是引用，值会发生变化
1. 隐式捕获

    手动书写捕获列表有时候是非常复杂的，这种机械性的工作可以交给编译器来处理，这时候可以在捕获列表中写一个 & 或 = 向编译器声明采用引用捕获或者值捕获.
1. 表达式捕获

总结一下，捕获提供了 Lambda 表达式对外部值进行使用的功能，捕获列表的最常用的四种形式可以是：

    [] 空捕获列表
    [name1, name2, ...] 捕获一系列变量
    [&] 引用捕获, 让编译器自行推导引用列表
    [=] 值捕获, 让编译器自行推导值捕获列表

从 C++14 开始， Lambda 函数的形式参数可以使用 auto 关键字来产生意义上的泛型.
```c++
auto add = [](auto x, auto y) {
    return x+y;
};

add(1, 2);
add(1.1, 2.2);
```

## 右值引用
**左值(lvalue, left value)**，顾名思义就是赋值符号左边的值。准确来说， 左值是表达式（不一定是赋值表达式）后依然存在的持久对象。

**右值(rvalue, right value)**，右边的值，是指表达式结束后就不再存在的临时对象。

而 C++11 中为了引入强大的右值引用，将右值的概念进行了进一步的划分，分为：纯右值、将亡值。

**纯右值(prvalue, pure rvalue)**，纯粹的右值，要么是纯粹的字面量，例如 10, true； 要么是求值结果相当于字面量或匿名临时对象，例如 1+2. 非引用返回的临时变量、运算表达式产生的临时变量、 原始字面量、Lambda 表达式都属于纯右值.

需要注意的是，字符串字面量只有在类中才是右值，当其位于普通函数中是左值:
```c++
class Foo {
        const char*&& right = "this is a rvalue"; // 此处字符串字面量为右值
public:
        void bar() {
            right = "still rvalue"; // 此处字符串字面量为右值
        } 
};

int main() {
    const char* const &left = "this is an lvalue"; // 此处字符串字面量为左值
}
```

**将亡值(xvalue, expiring value)**，是 C++11 为了引入右值引用而提出的概念（因此在传统 C++ 中， 纯右值和右值是同一个概念），也就是即将被销毁、**却能够被移动的值**.

要拿到一个将亡值，就需要用到右值引用：`T &&`，其中 T 是类型. 右值引用的声明让这个临时值的生命周期得以延长、只要变量还活着，那么将亡值将继续存活.

C++11 提供了 std::move 这个方法将左值参数无条件的转换为右值.

## 指针
在 C++ 中，指针是一种基础数据类型，它存储了内存地址的值。通过指针，可以直接读取或修改相应内存地址处的数据。指针是 C/C++ 强大功能的一个关键组成部分，允许直接操作内存，这在底层编程和系统编程中非常有用，但这一切能力的代价就是指针操作的 高风险.

指针的定义语法：
```c++
Typename * ptrName;

// 指针定义风格，下面的声明都正确
int *p;   // C风格，旨在强调 （*p）是一个整形值
int* p;    // 经典C++风格，只在强调 p是一个整形指针类型（int*）

// 集团推荐的风格,指针、引用都是居中，两边留空格
int * p;     // 指针
int & a = xx;  // 左值引用
int && a = xx;  // 右值引用
```

不论指针的类型是什么，指针本身的内存占用是相同的，64 位系统占用 8 个字节。指针类型存储的是地址编号，本质上是整形，可以进行计算，但对地址的乘除法是没有意义的，加减法是有意义的，表示地址的偏移。

对指针进行 +1 操作，指针将会偏移其指向的类型所占用的字节数（编译器根据指针的类型确定偏移的字节数），下面有个实际例子：
```c++
int a = 123;  // 假设 a 地址为 0xfffff100
int * p = &a;  // 此时 p 中存储的值为 0xfffff100
p = p + 1;     // 此时 p 中存储的值为 0xfffff104 （0xfffff100偏移4个字节，即int变量占用的大小）
```

### 常量指针
常量指针指向一个常量值，不管指向的变量本身是否声明为常量都不能通过指针来修改指向的内容，但指针本身可以重新赋值指向新的地址.
```c++
int value = 5;
const int * p = &value;    // p是一个常量指针
int const * q = &value;   // 和上面的声明等价 
*p = 10;           // 非法，*p是常量不能修改

int a = 6;
p = &a;            // 合法，p本身不是常量，可以重新赋值
```
常量指针在函数传参时非常有用，它可以限制函数内部通过指针非法地修改原始内容.

### 指针常量
指针常量表示指针本身是常量，必须在声明时初始化，之后不能指向其他地址，但可以通过指针修改指向的内容。
```c++
int value = 5;
int * const p = &value;   // p是常量
*p = 6;     // 合法

int a = 7;
p = &a;    // 非法
```

要记住这两种声明的区别有个简单的方法：看 const 修饰是什么：
- `const int * p ：const 修饰 *p，即 *p 是常量`
- `int * const p ：const 修饰 p，即 p 是常量`

### 指针和数组名的异同
在 C++ 中，数组名在绝大多数场景下可以看做是指针，在这些场景下数组名和指向该数组首个元素的指针是等价的.
```c++
int arr[5] = {1, 2, 3, 4, 5};
int * p1 = arr;     // arr 被当做指向数组首元素的指针
int * p2 = &arr[0];    // 取arr首个元素的地址
// 这种情况下 p1 和 p2 是等价的
if (p1 == P2) {      // 检测会通过
    cout << "p1,p2是等价的" << endl;  
    cout << *p1 << endl;  // 打印 1
    cout << *p2 << endl;   // 打印 1
}

// 使用指针访问数组
// 指针方式
cout << *(p1 + 1) << endl;  // 访问数组第二个元素，这种方式符合指针的计算规则
// 类似数组名的使用方式
cout << p1[1] << endl;// p1虽然是指针，索引访问方式依然有效，本质是*(p1 + 1)的语法糖
```

指针和数组名有区别的地方：
```c++
int arr[5] = {1, 2, 3, 4, 5};
int * p1 = arr;

cout << sizeof(arr) << endl;  // 打印结果：20 
cout << sizeof(p1) << endl;    // 打印结果：8
// sizeof(arr)为数组本身的大小，这里是 5个int占用20字节
// sizeof(p1)为指针本身大小，64位系统中占用8个字节
```

此外 & 取地址运算符对于 指针和数组名的处理也是不同的:
```c++
cout << &arr << endl;      // 0x16b98aa40
cout << &arr + 1 << endl;    // 0x16b98aa54
cout << &arr[0] << endl;    // 0x16b98aa40
cout << &arr[0] + 1 << endl;  // 0x16b98aa44

// 可以看出 &arr 和 &arr[0] 的值是一样的，但是指针偏移1后
// (&arr + 1) 在 &arr 的基础上偏移了20（0x14）个字节
// (&arr[0] + 1) 在 &arr[0] 的基础上偏移了4个字节
```

对于数组名进行 & 取地址，得到的整个数组的地址，虽然值和首元素地址相同，但其指针类型是不同的。
- `&arr 得到的类型是 int (*)[5]` ，这是一个指向包含 5 个整数数组的指针
- `&arr[0] 得到的类型是 int *`，这是一个整型指针

数组和指针结合使用时会有一些容易出错的点：
```c++
int * p[10];  // p是一个包含10个int变量的数组
int (*p)[10];   // p是一个指向拥有10个int变量的数组的指针

// * [] 两个运算符的优先级不同，[]的优先级更高
// 第一个语句声明了 p[10], int * 是类型
// 第二个语句有括号改变了优先级，因此 p 是一个指针，剩下的部分定义了类型

// 下面函数指针也会有类似的定义
int (*pf)(int, int); // pf是指向形如 int func(int, int) 的函数指针
```

## 智能指针
c++11用新类shared_ptr简化了内存管理以及对象在内存中的共享.shared_ptr是一个模板, 是从自由存储分配的对象的包装器.
包装器通过引用计数来追踪其他还有多少个指针在引用对象., 计数器归零, 对象即可安全删除, 分配的内存归还给自由存储.

> 堆是操作系统维护的一块内存，而自由存储是C++中通过new与delete动态分配和释放对象的抽象概念. 堆与自由存储区并不等价, 但基本上所有的C++编译器默认使用堆来实现自由存储, 同时开发者也可以通过重载操作符，改用其他内存来实现自由存储，例如全局变量做的对象池，这时自由存储区就区别于堆了.

注意: shared_ptr类不是万能的, 循环引用列表会出问题, 因为引用计数永远不为0, 内存会一直无法回收. c++提供另外的weak_ptr类来解决该问题, 只要weak_ptr是唯一的对象引用, 该对象就会被销毁, 只要至少一个链接由weak_ptr连接, 整个循环列表最终都会被销毁.

c++11还提供了unique_ptr类, 是一种独占的智能指针，它禁止其他智能指针与其共享同一个对象，从而保证代码的安全, 因此不能把它赋给其他任何指针.

现代 C++ 提供的智能指针：
1. std::unique_ptr

	std::unique_ptr 是一个独有所有权的智能指针。它保证同一时间内只有一个智能指针实例可以拥有一个给定的对象。当 std::unique_ptr 被销毁时，它所拥有的对象也会被销毁。std::unique_ptr 通常用于对资源有独占所有权的情况，并且它是不可以被复制的，但可以被移动，以便所有权可以从一个 std::unique_ptr 转移到另一个。
1. std::shared_ptr

	std::shared_ptr 实现了共享所有权的概念。它通过内部的引用计数机制来跟踪有多少个 std::shared_ptr 实例共享同一个对象。当最后一个这样的指针被销毁时，所拥有的对象将会被删除。std::shared_ptr 适用于多个拥有者需要管理同一个对象的生命周期的情况。
1. std::weak_ptr

	std::weak_ptr 是一种非拥有（弱）引用的智能指针。它不会增加对象的引用计数，因此不会阻止所指向的对象被销毁。std::weak_ptr 主要用于解决 std::shared_ptr 之间可能出现的循环引用问题。通过 std::weak_ptr，你可以观察一个对象，但不会造成所有权关系。

```c++
// 简单示例
// 定义智能指针
// C++11语法
std::unique_ptr<MyClass> my_unique_ptr(new MyClass());
std::shared_ptr<MyClass> my_shared_ptr(new MyClass());
// C++14提供了更安全更现代的方法
auto my_unique_ptr = std::make_unique<MyClass>();
auto my_shared_ptr = std::make_shared<MyClass>(); // 可以按照构造函数的定义传参

// 调用类的方法和普通指针类似
my_unique_ptr->func();
my_shared_ptr->func();

// 在需要传对象指针和引用的场景
// 类指针类型
void testFunc1(MyClass * p);
testFunc1(my_unique_ptr.get()); // 通过get获取原始指针

// 引用类型
void testFunc2(MyClass & ref);
testFunc2(*my_unique_ptr);    // 通过*运算符获取对象的引用
```

## 废弃(c++11开始)
- 字符串字面值常量赋值和初始化应用`const char *`取代`char *`.
- C++98 异常说明、 unexpected_handler、set_unexpected() 等相关特性被弃用，应该使用 noexcept
- auto_ptr 被弃用，应使用 unique_ptr
- register 关键字被弃用，可以使用但不再具备任何实际含义
- 如果一个类有析构函数，为其生成拷贝构造函数和拷贝赋值运算符的特性被弃用了
- C 语言风格的类型转换被弃用（即在变量前使用 (convert_type)），应该使用 static_cast、reinterpret_cast、const_cast 来进行类型转换
- 在最新的 C++17 标准中弃用了一些可以使用的 C 标准库，例如 `<ccomplex>、<cstdalign>、<cstdbool> 与 <ctgmath>` 等
- 还有一些其他诸如参数绑定（C++11 提供了 std::bind 和 std::function）、export 等特性也均被弃用

## 与C兼容性
在编写 C++ 时，也应该尽可能的避免使用诸如 void* 之类的程序风格. 而在不得不使用 C 时，应该注意使用 extern "C" 这种特性，将 C 语言的代码与 C++代码进行分离编译，再统一链接这种做法. 参考如下做法

```c++
// foo.h
#ifdef __cplusplus
extern "C" {
#endif

int add(int x, int y);

#ifdef __cplusplus
}
#endif

// foo.c
int add(int x, int y) {
    return x+y;
}

// 1.1.cpp
#include "foo.h"
#include <iostream>
#include <functional>

int main() {
    [out = std::ref(std::cout << "Result from C code: " << add(1, 2))](){
        out.get() << ".\n";
    }();
    return 0;
}
```

```makefile
C = gcc
CXX = clang++

SOURCE_C = foo.c
OBJECTS_C = foo.o

SOURCE_CXX = 1.1.cpp

TARGET = 1.1
LDFLAGS_COMMON = -std=c++2a

all:
	$(C) -c $(SOURCE_C)
	$(CXX) $(SOURCE_CXX) $(OBJECTS_C) $(LDFLAGS_COMMON) -o $(TARGET)
clean:
	rm -rf *.o $(TARGET)
```

## modern c++ start from c++11
### 语言可用性强化
- nullptr 出现的目的是为了替代 NULL. 因为在某种意义上来说，传统 C++ 会把 NULL、0 视为同一种东西，这取决于编译器如何定义 NULL，有些编译器会将 NULL 定义为 ((void*)0)，有些则会直接将其定义为 0.

    ```c++
    #include <iostream>
    #include <type_traits>

    void foo(char *);
    void foo(int);

    int main() {
        if (std::is_same<decltype(NULL), decltype(0)>::value)
            std::cout << "NULL == 0" << std::endl;
        if (std::is_same<decltype(NULL), decltype((void*)0)>::value)
            std::cout << "NULL == (void *)0" << std::endl;
        if (std::is_same<decltype(NULL), std::nullptr_t>::value)
            std::cout << "NULL == nullptr" << std::endl;
        if (std::is_same<decltype(NULL), decltype(long(0))>::value)
            std::cout << "NULL == long(0)" << std::endl; // llvm 11.1.0 x86_64

        std::cout << "NULL's type name:" << typeid(NULL).name() << std::endl; // l = long

        foo(0);          // will call foo(int)
        // foo(NULL);    // doesn't compile
        foo(nullptr);    // will call foo(char*)
        return 0;
    }

    void foo(char *) {
        std::cout << "foo(char*) is called" << std::endl;
    }
    void foo(int i) {
        std::cout << "foo(int) is called" << std::endl;
    }
    ```
- C++11 提供了 constexpr 让用户显式的声明函数或对象构造函数在编译期会成为常量表达式.

    从 C++14 开始，constexpr 函数可以在内部使用局部变量、循环和分支等简单语句. C++17 将 constexpr 这个关键字引入到 if 语句中，允许在代码中声明常量表达式的判断条件.

    **const 常数与常量表达式不是一个概念**

    ```c++
    constexpr int len_foo_constexpr() {
    return 5;
    }

    int main() {
        const int len_2 = 1;
        constexpr int len_2_constexpr = 1 + 2 + 3;
        // char arr_4[len_2];                // 非法
        char arr_4[len_2_constexpr];         // 合法
        char arr_6[len_foo_constexpr() + 1]; // 合法
    }
    ```

## FAQ
### 编译器支持c++进度
- [C++ Standards Support in GCC](https://gcc.gnu.org/projects/cxx-status.html)
- [C++ Support in Clang](http://clang.llvm.org/cxx_status.html)
- [C++ 20的范围稳了，主流编译器的完整支持会在什么时候？](https://www.zhihu.com/question/313443905)

用法:
- `clang++ -std=c++2a`

### `gcc -lstdc++ slice.cpp`报: undefined reference to `std::cout'
明明使用了`use namespace std;`, 编译时还是报错.

使用`gcc slice.cpp -lstdc++`或`g++ slice.cpp`. 用gcc时要注意参数位置.

### 用GCC或者Clang观察预处理后的C++代码
`g++/clang++ -E -P -std=c++11 -Wall -DBOOST_LOG_DYN_LINK -c ./main.cc >> main.output`

### EXPORT_API
定义动态链接库的导出符号
