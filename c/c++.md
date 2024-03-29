# c++
传统 C++ : C++98 及其之前的 C++ 特性
现代 C++ :  C++11/14/17/20, 为传统 C++ 注入的大量特性使得整个 C++ 变得更加像一门现代化的语言.

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
构造函数是用于初始化class对象的成员函数,声明类的对象时会自动调用构造函数, 要求:
1. 构造函数的名称与类的名称相同
1. 构造函数没有任何类型的返回值

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

### 继承
继承是基于一个基类创建新类(派生类)的过程. 派生类自动拥有基类的所有成员变量和函数, 并可根据需要添加更多的成员函数/变量, 但**私有成员函数/构造函数/析构函数/赋值操作符=均不会被继承**.

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

函数模板的定义实际是多个函数定义的"合集". 编译器仅针对**用到的每种类型**生成单独的函数定义, 未用到的类型不生成定义.

```c++
template<typename T> // `template<typename T>`为模板前缀, T为类型参数. 函数模板支持多个类型参数, 但通常只需一个.
void swapValues(T& variable1, T& variable2)
{
    T temp;
	...
}
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

## 智能指针
c++11用新类shared_ptr简化了内存管理以及对象在内存中的共享.shared_ptr是一个模板, 是从自由存储分配的对象的包装器.
包装器通过引用计数来追踪其他还有多少个指针在引用对象., 计数器归零, 对象即可安全删除, 分配的内存归还给自由存储.

> 堆是操作系统维护的一块内存，而自由存储是C++中通过new与delete动态分配和释放对象的抽象概念. 堆与自由存储区并不等价, 但基本上所有的C++编译器默认使用堆来实现自由存储, 同时开发者也可以通过重载操作符，改用其他内存来实现自由存储，例如全局变量做的对象池，这时自由存储区就区别于堆了.

注意: shared_ptr类不是万能的, 循环引用列表会出问题, 因为引用计数永远不为0, 内存会一直无法回收. c++提供另外的weak_ptr类来解决该问题, 只要weak_ptr是唯一的对象引用, 该对象就会被销毁, 只要至少一个链接由weak_ptr连接, 整个循环列表最终都会被销毁.

c++11还提供了unique_ptr类, 是一种独占的智能指针，它禁止其他智能指针与其共享同一个对象，从而保证代码的安全, 因此不能把它赋给其他任何指针.

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