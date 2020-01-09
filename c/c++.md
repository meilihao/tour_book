# c++
变量的名称称为标识符, 其必须以字母或下划线开头, 且仅包含`字母, 数字, 下划线`. c++是大小写敏感的.

使用namespace的原因: c++有大多东西需要命名, 一个名称可能有多种定义, 为消除歧义, 将不同的项划分到不同的集合中以避免冲突.

cin运行输入时以一个或多个空格或一个换行符来分隔输入的多个参数, 它以`Enter键`表示结束.

## header
```c++
#include <iostream> // 推荐, 新编译器支持
#include <iostream.h> // 不推荐, 旧编译器要求追加".h" 
```

list:
```c++
#include <fsstream> // for file I/O
#include <iostream> // 用于cout
#include <istream> // 通用型流参数
#include <cstdlib> // for exit
#include <cctype> // for 预定义字符函数
```

## 类型
c++11支持auto类型, 由右侧表达式推断变量类型.
c++11 支持`decltype (expr) value`, 可根据expr确定变量value的类型.
c++没有字符串的原生类型, 但可使用`#include<string>`来定义和操作字符串.
新版编译器支持bool类型, 而不是用1/0来表示true/false.
c++11支持枚举类`enum class Days {Sun,...}; Days d = Days::Sun; `, 以避免传统枚举的一些问题.
break用于退出while, do-while,for, switch. 
新版c++使用预定义的强制类型转换函数`static_cast<type>(xxx)`来取代旧有的`type(xxx)`.
向量作用和数组一样, 但支持运行时改变长度.
c++要求定义数组时，必须明确给定数组的大小, 否则应使用new 动态定义数组.
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

**推荐使用[定宽整数类型 from C++11](https://zh.cppreference.com/w/cpp/types/integer)**, 它定义于头文件 `<cstdint>`里.

## 字符串
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

## 数组
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

多维数组:
```c++
#include<iostream>
#include <typeinfo>
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

## [向量](https://en.cppreference.com/w/cpp/container/vector)
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

## 函数
部分数值处理的预定义函数:
![部分数值处理的预定义函数](/misc/img/c/截图_2020-01-08_18-00-46.png)

c++要求函数调用前要有具体的函数实现或存在函数声明.

C++ 引用传递:
形参相当于是**实参的"别名"**，对形参的操作其实就是对实参的操作，在引用传递过程中，被调函数的形式参数虽然也作为局部变量在栈中开辟了内存空间，但是这时存放的是由主调函数放进来的实参变量的地址. 被调函数对形参的任何操作都被处理成间接寻址，即通过栈中存放的地址访问主调函数中的实参变量. 正因为如此，被调函数对形参做的任何操作都影响了主调函数中的实参变量.

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
函数名相同,  参数列表(个数,类型,顺序)不同

注意:
1. 返回值类型与函数重载无关
1. 调用函数时, 实参的隐式类型转换可能产生二义性.

本质: 采用了name mangling(也叫name decoration)的技术, 即c++编译器默认会对符号名(比如函数名)进行改编,修饰, 简单理解就是重命名, 且不同编译器有不同的重命名规则.

> 可通过编译时选禁止优化和不生成调试信息, 再通过反汇编验证.

### 默认参数(**个人不推荐使用**)
注意:
1. 默认参数只能放在函数声明处**或者**定义处，能放在声明处就放在声明处
1. 如果某个参数是默认参数，那么它后面的参数必须都是默认参数

    因为非默认参数的参数必须要给出具体值，而调用函数传递参数的时候是从左到右的，所以非默认参数前面的都必须要传值进来. 那么默认参数后面的当然也得都为默认参数.

1. 不要重载一个带默认参数的函数

## class
```c++
class Day
{
   public:
     void output();//仅有函数声明
      int month;
      int day;

	  void printYear()
           {
                 cout<< "year = " << year << endl;
           }

	 private:// 仅成员函数中可见, 在其他地方都无法访问
                    int year;
};

// `::`是作用域解析操作符, 作用于类名
void Day::output() // 成员函数定义
{
	cout << "month = "<< month
		<<", day = " << day << endl;
}
```

### func list
- endl

从定义中看出，endl是一个函数模板，它实例化之后变成一个模板函数，作用是输出一个换行符，并立即刷新缓冲区.

## FAQ
### 编译器支持c++进度
- [C++ Standards Support in GCC](https://gcc.gnu.org/projects/cxx-status.html)
- [C++ Support in Clang](http://clang.llvm.org/cxx_status.html)
- [C++ 20的范围稳了，主流编译器的完整支持会在什么时候？](https://www.zhihu.com/question/313443905)
