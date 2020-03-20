# python

[书单](https://zhuanlan.zhihu.com/p/34378860):
- 入门
    - Python基础教程（第3版）
    - Python编程：从入门到实践
    - Python编程快速上手
- 进阶
    - Python高性能编程
    - 流畅的Python
    - Python Cookbook（第3版）中文版
    - 编写高质量代码: 改善Python程序的91个建议

## 环境
### 安装
```
$ sudo apt-get install python3.8
$ sudo apt install python3-pip
$ sudo apt install python-pip # pip2 for python2.7
```

> [pip安装, 未测试](https://pip.pypa.io/en/stable/installing/)

> 作为过渡, python3有很多特性被移植到了python2.7(将于2020.1.1终止支持), 因此如果程序可在python2.7运行就可通过python3自带的转换工具2to3(`python -m pip install  2to3`/`sudo apt install 2to3`)无缝迁移到Python3.

python解析器:
1. CPython : 使用C实现的解析器, 最常用的解析器, 通常说的python解析器就是指它.
1. PyPy : 使用Python实现的解析器

### 切换python版本
- [update-alternatives](https://blog.csdn.net/White_Idiot/article/details/78240298)
- alias
    ```
    $ vim ~/.bashrc
    $ alias python='/usr/bin/python3.8'
    ```

### 其他
```sh
$ python -m pip -V # 检查是否安装pip成功
$ mkdir -p ~/.pip
$ vim ~/.pip/pip.conf # [为pip换源](https://blog.csdn.net/xuezhangjun0121/article/details/81664260)
[global]
index-url = https://pypi.tuna.tsinghua.edu.cn/simple
```

> pip配置查找: `python -m  pip config list -v`

其他pip.conf:
```
[global]
index-url = http://user:user@192.168.0.236:8081/repository/pypi-group/simple

[install]
trusted-host=192.168.0.236 # 因为未使用https
```

> python module可在https://pypi.org上查找.

vscode:
- 安装插件`Python`(ms-python.python), 但它的智能提示和跳转无法在pygame上生效, 但pycharm社区版正常.

其他安装:
```
$ python -m pip install -U pygame --user # 安装pygame, 无法安装成功, 缺少依赖. 但系统自带的python3.7可以安装成功并运行demo
$ python -m pygame.examples.aliens # 运行pygame demo
$ python -m pip  show pygame # 查找安装位置
$ pip install --target /usr/lib/python3.7/dist-packages netifaces # 指定pip安装目录
```

> 可使用`apt-file search "sdl-config"`查找pygame的依赖

### idle
- `>>>`为提示符
- 连续两个回车表示结束.
- 中断程序执行: ctrl+c
- 退出idle:
    - 输入`quit()`
    - `ctrl+d`
    - `ctrl+z`
    - 直接关闭idle窗口
- F1可查看python document

## 摘要
1. python程序的构成
    
    1. python程序由模块组成.
    1. 模块由语句组成.
1. 执行顺序

    python的执行顺序与golang初始化类似, module 对应了 package.
    C++中一main函数为执行的起点; Python中首先执行最先出现的非函数定义和非类定义的没有缩进的代码，会从上而下顺序执行.
    程序中为了区分主动执行还是被调用，Python引入了变量__name__，当文件是被调用时，__name__的值为模块名，当文件被执行时，__name__为'__main__'.
1. 没有接口
1. python一切皆对象. 每个对象都有标识(id), 类型, 值.
    ```python
    >>> a = 3 # 内存形式: a(id: 8791360980704. 在栈) 引用了 object(id: 8791360980704, type: int, value: 3. 在堆). 因为object包含类型(因此python是强类型语言), 因此python变量不需要显示声明类型, 可由解析器推导.
    >>> a
    3
    >>> id(a) # id() 函数返回对象的唯一标识符，标识符是一个整数. CPython 中 id() 函数用于获取对象的内存地址.
    8791360980704
    >>> type(a)
    <class 'int'>
    >>> a = "string" # Python是动态类型，即在编程中允许随意改变变量的类型，这个过程称为 "绑定（Binding）". "绑定"只存在于动态类型语言中；对于静态语言如C，是"强制类型转换".
    >>> a
    'string'
    ```

    在 python 中赋值语句总是建立对象的引用值，而不是复制对象. 因此，python 变量更像是指针，而不是数据存储区域.

### [装饰器](https://www.liaoxuefeng.com/wiki/1016959663602400/1017451662295584)
- @try_except


## 变量和简单数据类型
变量名只能使用字母,数字和下划线, 且不能以数字开头. 不能使用python保留用于特殊用途的单词(关键字和函数名, 在idle中可用`help() + keywords`查看)作为变量名.

如果想要在函数内使用全局作用域的变量时，需要加上global修饰符, 即提示python 解释器，表明被其修饰的变量是全局变量.

> 事实上，很多文章推崇另外的一种方法来使用全局变量：使用单独的global文件.

> Python 的全局变量是模块 (module) 级别的.

> 双下划线开头和结尾的名称通常具有特殊含有, 尽量避免这种写法.

> del可删除变量(对象还在, 但对象没有引用后会被gc回收).

> `x=y=123` => `x=123;y=123`: 链式赋值, 用于将一个对象赋值给多个变量. python也支持类似go的`a,b,c=4,5,6`和`a,b=b,a`

> **python不支持常量**.

python中用引号(单引号/双引号)包裹的都是字符串.

方法是Python可对数据执行的操作.

> python3的print是函数需要括号, Python2的print括号则可有可无.

Python使用`**`表示乘方运算.
Python将带小数点的数字都称为浮点数.
Python的`"""`(docstring, 文档字符串)类似golang的"``", 输出时原样输出, Python使用它们来生成有关程序中函数的文档.
Python使用`#`作为注释标识, 但也可将`"""`用作多行注释.

在Python中，用方括号`[]`来表示列表(其他语言的数组)，并用逗号来分隔其中的元素.  **list中元素的类型可以不同**.

> 列表和元组的主要不同在于，列表是可以修改的，而元组不可以.

Python为反向访问一个列表元素提供了一种特殊语法, 即通过将索引指定为`-n`, `-1`表示最后一个元素.

Python根据缩进来判断代码行与前一个代码行的关系.

列表的部分元素——Python称之为切片, 它通过复制来创建.
创建切片时需指定要使用的第一个元素和最后一个元素的索引.

Python将不能修改的值称为不可变的，而不可变的列表被称为元组, 其用圆括号包裹, 想修改元组中的元素时必须给存储元组的变量重新赋值.
元组作用:
1. 用作映射中的键（以及集合的成员），而列表不行
1. 有些内置函数和方法返回元组

Python的布尔表达式的结果要么为True，要么为False.
Python并不要求if-elif结构后面必须有else代码块. 在有些情况下，else代码块很有用; 而在其他一些情况下，使用一条elif语句来处理特定的情形更清晰.
在if语句中将列表名用在条件表达式中时，Python将在列表至少包含一个元素时返回True, 并在列表为空时返回False.
在Python中，字典是一系列键—值对, 用`{}`包裹, 每个键都与一个值相关联，可以使用键来访问与之相关联的值, Python可将任何Python对象用作字典中的值.
遍历字典时, Python会默认遍历所有的键. 字典的`update()`方法用于更新字典中的键/值对

条件测试:
- and : 其他语言的`&&`
- or : 其他语言的`||`
- in : 判断特定的值是否已包含在列表中
- not in : 检查特定值是否不包含在列表中

函数`input("提示")`让程序暂停运行，等待用户输入一些文本, 获取用户输入后，Python将其存储在一个变量中，以方便之后使用.
函数`int()`将数字的字符串表示转换为数值表示.

> Python 2.7也包含函数input()，但它将用户输入解读为Python代码，并尝试运行它们, 其用函数`raw_input()`来提示用户输入.

> any(iterable) 函数用于判断给定的可迭代参数 iterable 是否全部为 False，则返回 False，如果有一个为 True，则返回 True.

for循环是一种遍历列表的有效方式，但在**for循环中不应修改列表**，否则将导致Python难以跟踪其中的元素. 要在遍历列表的同时对其进行修改，可使用while循环.

### 引用
list引用:
```py
def add_list(p):
    p = p + [1]
p1 = [1,2,3]
add_list(p1)
print(p1)
>>> [1, 2, 3]

def add_list(p):
    p += [1]
p2 = [1,2,3]
add_list(p2)
print(p2)
>>>[1, 2, 3, 1]

a = []
b = {'num':0, 'sqrt':0}
resurse = [1,2,3]
for i in resurse:
  b['num'] = i
  b['sqrt'] = i * i
  a.append(b)
print(a) # 这是由于a中的元素就是b的引用
>>> [{'num': 3, 'sqrt': 9}, {'num': 3, 'sqrt': 9}, {'num': 3, 'sqrt': 9}]

a = []
resurse = [1,2,3]
for i in resurse:
   a.append({"num": i, "sqrt": i * i})
>>> [{'num': 1, 'sqrt': 1}, {'num': 2, 'sqrt': 4}, {'num': 3, 'sqrt': 9}]

>>> values = [0, 1, 2]
>>> values[1] = values
>>> values
[0, [...], 2]       # 实际结果. 可以说 Python 没有赋值，只有引用.  这样相当于创建了一个引用自身的结构，所以导致了无限循环.
[0, [0, 1, 2], 2]   # 预想结果

>>> values = [0, 1, 2]
>>> values[1] = values[:] # values[:] 生成对象的拷贝或者是复制序列，不再是引用和共享变量，但此法只能顶层复制. 深复制的方法是`copy.deepcopy(a)`
>>> values
[0, [0, 1, 2], 2]
```

这区别主要是由于`=`操作符会新建一个新的变量保存赋值结果，然后再把引用名指向`=`左边，即修改了原来的p引用，使p成为指向新赋值变量的引用. 而+=不会，直接修改了原来p引用的内容. **事实上+=和=在python内部使用了不同的实现函数**.

## 模块
import语句允许在当前运行的程序文件中使用模块中的代码.

只需编写一条import语句并在其中指定模块名，就可在程序中使用该模块中的所有函数, 调用形式如下:`module_name.function_name()`.
导入模块中的特定函数：`from module_name import function_0, function_1, function_2`, 若使用这种语法，调用函数时就无需使用句点.
导入时给函数指定别名：`from module_name import function_name as fn`
导入时给模块指定别名：`import module_name as mn`
导入模块中的所有函数: `from pizza import *`, 此时可通过名称来调用每个函数，而无需使用句点表示法. **不推荐**: 如果模块中有函数的名称与你的项目中使用的名称相
同, 可能导致意想不到的结果, 因为Python可能遇到多个名称相同的函数或变量, 进而被覆盖.

> 实际写代码的实践中，`import *`的做法是严格被禁止的，它容易造成包中模块名与当前命名空间的名称冲突.

`__name__`是python的一个内置属性，用来表示当前模块的名字. `__main__`是顶层代码执行作用域的名字. 因此`if __name__ == '__main__'`表示当模块被直接运行时，代码块将被运行，当模块是被导入时，代码块不被运行.

module搜索顺序:
- 内建模块 : 除了sys.builtin_module_names 列出的内置模块之外，还会加载其他一些标准库，都存放在sys.modules字典中
- sys.path
- 环境变量PYTHONPATH

> 路径应为**module的上级目录**.

### __init__.py
__init__.py 文件的作用是将所在的文件夹变为一个Python package, 因此Python 中的每个package中，都有__init__.py 文件, 且python在导入一个包时，实际上是导入了它的__init__.py文件.

__init__.py中的__all__变量，可用于模块导入时限制，比如`from module import *`, 此时被导入模块若定义了__all__变量，则只有__all__内指定的变量、方法、类可被导出; 若没定义，导出将按照以下规则执行：
1. 此 package 被导入，并且执行 __init__.py 中可被执行的代码
1. __init__.py 中定义的 variable 被导入
1. __init__.py 中被显式导入的 module 被导入

__all__对于 `from <module> import <member>`导入方式并没有影响.

### __import__()
参考：
- [Python中__import__()的fromlist参数用法](https://docs.lvrui.io/2017/10/13/Python%E4%B8%AD-import-%E7%9A%84fromlist%E5%8F%82%E6%95%B0%E7%94%A8%E6%B3%95/)

当使用import导入Python模块的时候，默认调用的是__import__()函数. 直接使用该函数的情况很少见，一般用于动态加载模块.

参数fromlist指明需要导入的子模块名，level指定导入方式（相对导入或者绝对导入， 默认两者都支持）．

### 关于.pyc 文件 与 .pyo 文件
.py文件的汇编,只有在import语句执行时进行，当.py文件第一次被导入时，它会被汇编为字节代码，并将字节码写入同名的.pyc文件中. 后来每次导入操作都会直接执行.pyc 文件（当.py文件的修改时间发生改变，这样会生成新的.pyc文件），在解释器使用-O选项时，将使用同名的.pyo文件，这个文件去掉了断言（assert）、断行号以及其他调试信息，体积更小，运行更快.（使用-OO选项，生成的.pyo文件会在`-O`的基础上再去除__doc__ string(文档信息)).

pyc的生成时机是在执行了 import 指令之后. import时已经存在 pyc 的话，就可以直接载入而省去编译过程, python还会在 pyc 文件中存储的创建时间信息来保证pyc文件是最新的. 当执行 import 指令的时候，如果已存在 pyc 文件，Python 会检查创建时间是否晚于代码文件的修改时间，这样就能判断是否需要重新编译，还是直接载入了; 如果不存在 pyc 文件，就会先将 py 文件编译.

## package
参考:
- [彻底明白Python package和模块](https://www.jianshu.com/p/178c26789011)

python 是通过module组织代码的，每一个module就是一个python文件(一个模块可以import其它模块)，但是package是通过modules来组织的. 因为package 是一个包含 __init__.py 的文件夹, 因此package 是至少包含一个 __init__.py 的 module.

Python 的 package 以及 package 中的 __init__.py 共同决定了 package 中的 module 是如何被外界访问的.

package 的初始化工作: 一个 package 被导入，不管在什么时候 __init__.py 中的代码只执行一次.

> 在Python解释器运行中，一个模块只可以被import一次，除非使用importlib.reload()

> 包名构建了一个Python模块的命名空间. 比如，模块名A.B表示A包中名为B的子模块.

> 想使用子package的内容，但是在父package的__init__.py的文件内并没有导入，此时需要手动导入

## 函数
使用`def`定义函数.
**关键字实参**是传递给函数的名称—值对, 此时无需考虑函数调用中的实参顺序，还清楚地指出了函数调用中各个值的用途.
Python支持指定默认值, 给形参指定默认值时，等号两边不要有空格.

将列表传递给函数后，函数就可对其进行修改, 在函数中对这个列表所做的任何修改都是永久性的.
向函数`function_name(list_name[:])`传递列表的副本而不是原件; 这样函数所做的任何修改都只影响副本，而丝毫不影响原件.
`def make_pizza(*toppings): `的形参名*toppings中的星号让Python创建一个名为toppings的空元组, 并将收到的所有值都封装到这个元组中, 以实现像函数传递任意数量的实参.
如果要让函数接受不同类型的实参，必须在函数定义中将接纳任意数量实参的形参放在最后, 比如`def make_pizza(size, *toppings): `.
`def build_profile(first, last, **user_info):`的形参**user_info中的两个星号让Python创建一个名为user_info的空字典，并将收到的所有名称—值对都封装到这个字典中.

## 类
```python
class Dog(): # 在Python中，类名的首字母应大写, 这个类定义中的括号是空的.
    """A simple attempt to model a dog."""
    
    def __init__(self, name, age):
        """Initialize name and age attributes."""
        self.name = name
        self.age = age
        self.owner = "chen" #提供初始值
        
    def sit(self):
        """Simulate a dog sitting in response to a command."""
        print(self.name.title() + " is now sitting.")

    def roll_over(self):
        """Simulate rolling over in response to a command."""
        print(self.name.title() + " rolled over!")
        

my_dog = Dog('willie', 6)
your_dog = Dog('lucy', 3)

print("My dog's name is " + my_dog.name.title() + ".")
print("My dog is " + str(my_dog.age) + " years old.")
my_dog.sit()
```

类中的函数称为方法. `__init__()`是一个特殊的方法，每当根据类创建新实例时，Python都会自动运行它. 在这个方法的名称中，开头和末尾各有两个下划线，这是一种约
定，旨在避免Python默认方法与普通方法发生名称冲突. 在这个方法的定义中，形参self必不可少，还必须位于其他形参的前面, 因为每个与类相关联的方法
调用都自动传递实参self，它是一个指向实例本身的**引用**，让实例能够访问类中的属性和方法, 同时self会自动传递，因此不需要显示传递它.

通过实例访问的变量称为属性, 比如`self.name = name`的`self.name`.

> 类方法的第一个参数必定是指向实例本身的**引用**, 通常命名为self, 但也可是其他名称(**不推荐**).

> 在Python 2.7中创建类时，需要在类名后面的括号内包含单词object.

类中的每个属性都必须有初始值，在方法`__init__()`内指定这种初始值是可行的(等于c++的构造方法)；如果对某个属性这样做了，就无需包含为它提供初始值的形参.

判断是否有属性或方法`hasattr(对象, name)`, 获取属性或方法`getattr(对象, name[, default])`, 设置属性或方法`setattr(对象,name,value)`.

```python
"""A class that can be used to represent a car."""
class Car():
    """A simple attempt to represent a car."""

    def __init__(self, manufacturer, model, year):
        """Initialize attributes to describe a car."""
        self.manufacturer = manufacturer
        self.model = model
        self.year = year
        self.odometer_reading = 0
        
    def get_descriptive_name(self):
        """Return a neatly formatted descriptive name."""
        long_name = str(self.year) + ' ' + self.manufacturer + ' ' + self.model
        return long_name.title()
    
    def read_odometer(self):
        """Print a statement showing the car's mileage."""
        print("This car has " + str(self.odometer_reading) + " miles on it.")
        
    def update_odometer(self, mileage):
        """
        Set the odometer reading to the given value.
        Reject the change if it attempts to roll the odometer back.
        """
        if mileage >= self.odometer_reading:
            self.odometer_reading = mileage
        else:
            print("You can't roll back an odometer!")
    
    def increment_odometer(self, miles):
        """Add the given amount to the odometer reading."""
        self.odometer_reading += miles

"""A set of classes that can be used to represent electric cars."""

class Battery():
    """A simple attempt to model a battery for an electric car."""

    def __init__(self, battery_size=60):
        """Initialize the batteery's attributes."""
        self.battery_size = battery_size

    def describe_battery(self):
        """Print a statement describing the battery size."""
        print("This car has a " + str(self.battery_size) + "-kWh battery.")  
        
    def get_range(self):
        """Print a statement about the range this battery provides."""
        if self.battery_size == 60:
            range = 140
        elif self.battery_size == 85:
            range = 185
            
        message = "This car can go approximately " + str(range)
        message += " miles on a full charge."
        print(message)
    
        
class ElectricCar(Car): # 
    """Models aspects of a car, specific to electric vehicles."""

    def __init__(self, manufacturer, model, year):
        """
        Initialize attributes of the parent class.
        Then initialize attributes specific to an electric car.
        """
        super().__init__(manufacturer, model, year) # 这行代码让Python调用ElectricCar的父类的方法__init__(), 让ElectricCar实例包含父类的所有属性
        self.battery = Battery()

my_tesla = ElectricCar('tesla', 'model s', 2016) 
print(my_tesla.get_descriptive_name()) 
my_tesla.battery.get_range() 
```

创建子类时，**父类必须包含在当前文件中，且位于子类前面**.
**定义子类时，必须在括号内指定父类的名称**, 且子类方法`__init__()`接受创建父类实例所需的信息. 子类会继承父类的属性和方法.

`super()`是一个特殊函数，帮助Python将父类和子类关联起来. 父类也称为超类（superclass）, 名称super因此而得名. `issubclass(子类名称, 父类名称)`可判断是否继承关系, 魔法属性`类名.__bases__`可得到父类的名称(可能有多个父类); `isinstance(对象, 类名称/父类名称)`判断是否类的实例对象.

> python支持多继承. 父类属性或方法名相同时, 靠前的覆盖靠后的, 同理子类会覆盖父类的.

可在子类中定义一个与要重写的父类方法同名的方法, 这样, Python将不会考虑这个父类方法，而只关注你在子类中定义的相应方法.

属性和方法清单以及文件都越来越长时, 可能需要将类的一部分作为一个独立的类提取出来, 这样可以将大型类拆分成多个协同工作的小类, 比如上面的`Battery类`.

从一个模块中导入多个类:`from car import Car, ElectricCar `.
导入整个模块:` import car`
导入模块中的所有类:`from module_name import * `, 不推荐, 原因和导入module相同.

属性和方法默认是可直接访问的, 在属性或方法名前面添加两个下划线`__`, 成员私有化的作用，确保外部代码不能随意修改对象内部的状态，增加了代码的安全性.
```python
class Person:
    def __init__(self, name):
        self.__name=name
    def __sayName(self):
        print("name", self.__name)

p=Person("chen")
p.__name="hao" # 等同添加属性
#p.sayName() # 不能调用
p._Person__sayName() # 能调用, 在python中, 私有化的类方法在被编译时, 真实的的名称是`"_"+类名+私有化的方法名`,同理私有属性也是这样.

import inspect
print(inspect.getmembers(p,predicate=inspect.ismethod)) # 验证上面的说法
print(inspect.getmembers(p))

# output:
# [('_Person__sayName', <bound method Person.__sayName of <__main__.Person object at 0x0000000002D17C70>>), ('__init__', <bound method Person.__init__ of <__main__.Person object at 0x0000000002D17C70>>)]
# ...
# ('__dict__', {'_Person__name': 'chen', '__name': 'hao'})
# ...
```

类代码块:
```python
class myclass:

    # class块中的语句(类定义到类第一个函数直接的代码块即为类代码块)，会立刻执行
    
    print('myclass')
    count = 0 # 变为类的属性

    def counter(self):
        self.count += 1

my = myclass()  # 实例化时会立即执行类代码块

my.counter()        # 调用counter方法
print(my.count)     # 输出结果：1

my.count = 'abc'    # 将count变量改变成字符串类型
print(my.count)     # 输出结果：abc

my.name = 'hello'   # 向my对象动态添加name变量
print(my.name)      # 输出结果：hello
```

类名应采用驼峰命名法，即将类名中的每个单词的首字母都大写，而不使用下划线. 实例名和模块名都采用小写格式，并在单词之间加上下划线.

对于每个类，都应紧跟在类定义后面包含一个文档字符串. 这种文档字符串简要地描述类的功能，并遵循编写函数的文档字符串时采用的格式约定.

```python
# 在Python 2.7中，继承语法稍有不同
class Car(object): 
    def __init__(self, make, model, year): 
        ...
class ElectricCar(Car): 
    def __init__(self, make, model, year): 
        super(ElectricCar, self).__init__(make, model, year)  # 函数super()需要两个实参：子类名和对象self, 为帮助Python将父类和子类关联起来
        ...
```

父类初始化方法中有调用了父类中定义的方法，恰好这个方法又被子类所覆盖，则`super().__init__(xing, age)`调用的父类初始化方法中调用的方法将是被覆盖后的方法:
```python
class A():
    def __init__(self, xing, gender):          # ！#1
        self.namea = "aaa"                     # ！#2
        self.xing = xing                       # ！#3
        self.gender = gender                   # ！#4
        self.funca()
 
    def funca(self):
        print("function a : %s" % self.namea)
 
 
class B(A):
    def __init__(self, xing, age):             # ！#5
        super().__init__(xing, age)     # ！#6（age处应为gender）
        self.nameb = "bbb"                     # ！#7
        ##self.namea="ccc"                     # ！#8
        ##self.xing = xing.upper()             # ！#9
        self.age = age                         # ！#10
 
    def funcb(self):
        print("function b : %s" % self.nameb)
 
    def funca(self):
        print("(override)function a : %s" % self.namea)
 
 
b = B("lin", 22)                               # ！#11
print(b.nameb)
print(b.namea)
print(b.xing)                                  # ！#12
print(b.age)
print(b.gender)                                # ！#13
b.funcb()
b.funca()

# output:
# (override)function a : aaa
# bbb
# aaa
# lin
# 22
# 22
# function b : bbb
# (override)function a : aaa
```

### 单例模式
参考:
- [Python中的单例模式实现](http://whosemario.github.io/2016/01/22/pattern-singleton/)

利用Python的metaclass实现单例(**推荐**, 其他方法有缺陷).

## 异常
```python
# 在 Windows系统中，在文件路径中使用反斜杠（\）而不是斜杠（/）.
with open('pi_digits.txt') as file_object:  # open('pi_digits.txt')返回一个表示文件pi_digits.txt的对象；然后将这个对象存储在后面使用的变量file_object中
    contents = file_object.read()  # 读取这个文件的全部内容，并将其作为一个长长的字符串存储在变量contents中
    print(contents)

with open(filename) as file_object: 
    for line in file_object: # 逐行读取, 不会清除行尾的换行
        print(line) 

with open(filename) as file_object: 
    lines = file_object.readlines() # 创建一个包含文件各行内容的列表

for line in lines: 
    print(line.rstrip()) 

filename = 'programming.txt' 
# 打开文件时，可指定读取模式（'r'）、写入模式（'w'）、附加模式（'a'）或能够读取和写入文件的模式('r+'). 如果你省略了模式实参，Python将以默认的只读模式打开文件
# 以写入（'w'）模式打开文件时: 如果要写入的文件不存在，函数open()将自动创建它; 如果指定的文件已经存在，Python将在返回文件对象前清空该文件
with open(filename, 'w') as file_object:  # 写入空文件
    file_object.write("I love programming.\n") 
    file_object.write("I love creating new games.") 
```

关键字with在不再需要访问文件后将其关闭.

Python使用被称为异常的特殊对象来管理程序执行期间发生的错误. 每当发生让Python不知所措的错误时，它都会创建一个异常对象. 如果未对异常进行处理，程序将停止，并显示一个traceback，其中包含有关异常的报告; 如果编写了处理该异常的代码，程序将继续运行.

异常是使用try-except代码块处理的, 异常会逐层向上抛出直至被捕获处理. 依赖于try代码块成功执行后(即未发生异常)的代码都应放到else代码块中.

Python有一个pass语句，可在代码块中使用它来让Python什么都不要做.
raise Exception: 手动抛出异常, 由外层的except处理, Exception是异常类.
自定义异常类需要继承Exception.

```python
try: 
    print(5/0) 
except ZeroDivisionError: 
    print("You can't divide by zero!") 

###
def MachineID():
    try:
        with open("/etc/machine-id", 'r') as f:
            return f.read().strip()
    except Exception as e:
        raise e
###
filename = 'alice.txt' 
try: 
    with open(filename) as f_obj: 
        contents = f_obj.read()
    # return
except FileNotFoundError: # 处理已知异常
    msg = "Sorry, the file " + filename + " does not exist." 
    print(msg) 
except XXXError as a: # 取别名 
    msg = "Sorry, the file " + filename + " does not exist." 
    print(msg) 
except (XXXError2 as a2, YYYError): # 处理多个异常 
    msg = "Sorry, the file " + filename + " does not exist." 
    print(msg) 
except: # 处理未知异常 
    msg = "Sorry, the file " + filename + " does not exist." 
    print(msg) 
else: 
    print("OK") 
finally: # 不管try的代码块是否存在return, finally都会执行, 存在return时会在return前执行.
    print"finally")
```

## 测试
Python标准库中的模块unittest提供了代码测试工具. 单元测试用于核实函数的某个方面没有问题；测试用例是一组单元测试，这些单元测试一起核实函数在各种情形下的行为都符合要求.

```python
class AnonymousSurvey():
    """Collect anonymous answers to a survey question."""
    
    def __init__(self, question):
        """Store a question, and prepare to store responses."""
        self.question = question
        self.responses = []
        
    def show_question(self):
        """Show the survey question."""
        print(self.question)
        
    def store_response(self, new_response):
        """Store a single response to the survey."""
        self.responses.append(new_response)
        
    def show_results(self):
        """Show all the responses that have been given."""
        print("Survey results:")
        for response in self.responses:
            print('- ' + response)

### 
import unittest
from survey import AnonymousSurvey

# 运行test_name_function.py时，所有以test_打头的方法都将自动运行.
class TestAnonymousSurvey(unittest.TestCase):
    """Tests for the class AnonymousSurvey."""
    
    def setUp(self):
        """
        Create a survey and a set of responses for use in all test methods.
        """
        question = "What language did you first learn to speak?"
        self.my_survey = AnonymousSurvey(question)
        self.responses = ['English', 'Spanish', 'Mandarin']
        
    
    def test_store_single_response(self):
        """Test that a single response is stored properly."""
        self.my_survey.store_response(self.responses[0])
        self.assertIn(self.responses[0], self.my_survey.responses) # 使用了unittest类最有用的功能之一：一个断言方法. 断言方法用来核实得到的结果是否与期望的结果一致
        
        
    def test_store_three_responses(self):
        """Test that three individual responses are stored properly."""
        for response in self.responses:
            self.my_survey.store_response(response)
        for response in self.responses:
            self.assertIn(response, self.my_survey.responses)
            

unittest.main() # 让python运行这个文件中的测试
```

Python在unittest.TestCase类中提供了很多断言方法. 断言方法检查你认为应该满足的条件是否确实满足.

unittest.TestCase中的断言方法:
- assertEqual(a, b) 核实a == b
-  assertNotEqual(a, b) 核实a != b
-  assertTrue(x) 核实x为True
-  assertFalse(x) 核实x为False
-  assertIn(item, list) 核实item在list中
-  assertNotIn(item, list) 核实item不在list中

如果在TestCase类中包含了方法setUp()，Python将先运行它，再运行各个以test_打头的方法.

## lambda
Python中，lambda函数也叫匿名函数，及即没有具体名称的函数，它允许快速定义单行函数，类似于C语言的宏，可以用在任何需要函数的地方, 这区别于def定义的函数.

lambda语法格式：`lambda 变量 : 要执行的语句`, 比如`lambda x : x ** 2`.

## [推导式](https://blog.csdn.net/yjk13703623757/article/details/79490476)
推导式comprehensions（又称解析式），是Python的一种独有特性. 推导式是可以从一个数据序列构建另一个新的数据序列的结构体. 共有三种推导式，在Python2和3中都有支持：

- 列表(list)推导式

    1. 使用[]生成list

    基本格式：
    ```python
    variable = [out_exp_res for out_exp in input_list if out_exp == 2]
    ```
    out_exp_res：列表生成元素表达式，可以是有返回值的函数
    for out_exp in input_list：迭代input_list将out_exp传入out_exp_res表达式中
    if out_exp == 2：根据条件过滤哪些值可以

    1. 使用()生成generator

    multiples = (i for i in range(30) if i%3 == 0)
    print(type(multiples))

- 字典(dict)推导式

    字典推导和列表推导的使用类似，只不过中括号改成大括号, 基本格式：`{ key_expr: value_expr for value in collection if condition }`
- 集合(set)推导式

    集合推导式跟列表推导式也是类似的。 唯一的区别在于它使用大括号{ }. 基本格式：`{ expr for value in collection if condition }`

## Python脚步
```text
#!/usr/bin/env python
....
```


### demo
```python
### 编码
# 使用ASCII、UTF-8和UTF-32编码将字符串转换为bytes
>>> "Hello, world!".encode("ASCII") 
b'Hello, world!' 
>>> "Hello, world!".encode("UTF-8") 
b'Hello, world!' 
>>> "Hello, world!".encode("UTF-32") 
b'\xff\xfe\x00\x00H\x00\x00\x00e\x00\x00\x00l\x00\x00\x00l\x00\x00\x00o\x00\x00\x00,\x00\ 
x00\x00 \x00\x00\x00w\x00\x00\x00o\x00\x00\x00r\x00\x00\x00l\x00\x00\x00d\x00\x00\x00!\x00\ 
x00\x00'
>>> "Hællå, wørld!".encode("ASCII", "ignore")  # 第二个参数: 如何处理错误, 这个参数默认为strict
b'Hll, wrld!' 
>>> x = bytearray(b"Hello!") # bytearray是bytes的可变版
>>> x[1] = ord(b"u") 
### 解构(唯一的前提是变量的数量必须跟序列元素(任何的序列或可迭代对象)的数量是一样的)
>>> s = 'Hello'
>>> a, b, c, d, e = s
>>> a
'H'
>>> data = [ 'ACME', 50, 91.1, (2012, 12, 21) ]
>>> _, shares, price, _ = data
>>> shares
50
### 
>>> record = ('ACME', 50, 123.45, (12, 18, 2012))
>>> name, *_, (*_, year) = record
>>> name
'ACME'
>>> year
2012
>>>
### 
def drop_first_last(grades):
    first, *middle, last = grades # 排除掉第一个和最后一个分数, 剩余的所有分数
    return avg(middle)
###
>>> record = ('Dave', 'dave@example.com', '773-555-1212', '847-555-1212')
>>> name, email, *phone_numbers = record
>>> name
'Dave'
>>> email
'dave@example.com'
>>> phone_numbers
['773-555-1212', '847-555-1212']
###
>>> line = 'nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false'
>>> uname, *fields, homedir, sh = line.split(':')
###
>>> items = [1, 10, 7, 4, 5, 9]
>>> def sum(items):
...     head, *tail = items
...     return head + sum(tail) if tail else head # 递归算法
...
>>> sum(items)
### 
previous_lines = deque(maxlen=history) # 构造函数会新建一个固定大小的队列. 当新的元素加入并且这个队列已满的时候， 最老的元素会自动被移除掉. 不设置maxlen时, 容量为不限制.
 previous_lines.append("x")
 ### 找最大或最小的 N 个元素
 import heapq
nums = [1, 8, 2, 23, 7, -4, 18, 23, 42, 37, 2]
print(heapq.nlargest(3, nums)) # Prints [42, 37, 23]
print(heapq.nsmallest(3, nums)) # Prints [-4, 1, 2]

portfolio = [
    {'name': 'IBM', 'shares': 100, 'price': 91.1},
    {'name': 'AAPL', 'shares': 50, 'price': 543.22},
    {'name': 'FB', 'shares': 200, 'price': 21.09}
]
cheap = heapq.nsmallest(2, portfolio, key=lambda s: s['price'])
### 堆排序
>>> heap = list([1, 8, 2, 23, 7, -4, 18, 23, 42, 37, 2])
>>> heapq.heapify(heap) # 找最小值
>>> heap
>>> heapq.heappop(heap) #将第一个元素弹出来, 重新找找出最小值放入heap[0]

# 当要查找的元素个数相对比较小的时候，函数 nlargest() 和 nsmallest() 是很合适的.
# 如果仅仅想查找唯一的最小或最大（N=1）的元素的话，那么使用 min() 和 max() 函数会更快些.
# 类似的，如果 N 的大小和集合大小接近的时候，通常先排序这个集合然后再使用切片操作会更快点 （ sorted(items)[:N] 或者是 sorted(items)[-N:] ）.
# 需要在正确场合使用函数 nlargest() 和 nsmallest() 才能发挥它们的优势 （如果 N 快接近集合大小了，那么使用排序操作会更好些）.
### 实现一个键对应多个值的字典
from collections import defaultdict

d = defaultdict(list)
d['a'].append(1)
d['a'].append(2)
d['b'].append(4)
d['b'].append(4)

d = defaultdict(set) # 值不重复
d['a'].add(1)
d['a'].add(2)
d['b'].add(4)
d['b'].add(4)
# output: defaultdict(<class 'set'>, {'a': {1, 2}, 'b': {4}})

# 普通做法, 推荐使用
d = {} # 一个普通的字典defaultdict
d.setdefault('a', []).append(1) <=> d["a"] = []; d["a"].append(1)
d = defaultdict(list)
d["a"].append(1)

### 字典排序
from collections import OrderedDict
d = OrderedDict() # 保证插入时的顺序, 与key和value无关
d['foo'] = 1

# OrderedDict 内部维护着一个根据键插入顺序排序的双向链表.
# 每次当一个新的元素插入进来的时候， 它会被放到链表的尾部, 对于一个已经存在的键的重复赋值不会改变键的顺序.
### 反转字典的kv
prices = {
    'ACME': 45.23,
    'AAPL': 612.78,
    'IBM': 205.55,
    'HPQ': 37.20,
    'FB': 10.75
}
# zip : 当多个实体拥有相同的值的时候，键会决定返回结果
min_price = min(zip(prices.values(), prices.keys())) # 需要注意的是 zip() 函数创建的是**一个只能访问一次的迭代器**
# output: min_price is (10.75, 'FB')
prices_sorted = sorted(zip(prices.values(), prices.keys())) # 类似的，可以使用 zip() 和 sorted() 函数来排列字典数据, 然后再取值
min(prices) # Returns 'AAPL' # 仅仅作用于键，而不是值
min(prices.values()) # Returns 10.75
min(prices, key=lambda k: prices[k]) # Returns 'FB'
### 把一个字符串变成 Unicode 码位的列表
>>> symbols = '$¢£¥€¤'
>>> codes = []
# 写法1:
>>> for symbol in symbols:
            # ord() : 是 chr() 函数（对于8位的ASCII字符串）或 unichr() 函数（对于Unicode对象）的配对函数，它以一个字符（长度为1的字符串）作为参数，返回对应的 ASCII 数值，或者 Unicode 数值，如果所给的 Unicode 字符超出了你的 Python 定义范围，则会引发一个 TypeError 的异常
            codes.append(ord(symbol)) 
# 写法2: 
>>>  codes = [ord(symbol) for symbol in symbols] # 个人不推荐
### 
symbols = '$¢£¥€¤'
beyond_ascii = [ord(s) for s in symbols if ord(s) > 127]
beyond_ascii = list(filter(lambda c: c > 127, map(ord, symbols))) # 功能与上面相同, 但推荐使用上面的写法
```

## 字符串
字符串格式设置中的类型说明符:
b 将整数表示为二进制数
c 将整数解读为Unicode码点
d 将整数视为十进制数进行处理，这是整数默认使用的说明符
e 使用科学表示法来表示小数（用e来表示指数）
E 与e相同，但使用E来表示指数
f 将小数表示为定点数
F 与f相同，但对于特殊值（nan和inf），使用大写表示
g 自动在定点表示法和科学表示法之间做出选择。这是默认用于小数的说明符，但在默认情况下至少有1位小数
G 与g相同，但使用大写来表示指数和特殊值
n 与g相同，但插入随区域而异的数字分隔符
o 将整数表示为八进制数
s 保持字符串的格式不变，这是默认用于字符串的说明符
x 将整数表示为十六进制数并使用小写字母
X 与x相同，但使用大写字母
% 将数表示为百分比值（乘以100，按说明符f设置格式，再在后面加上%）

```python
>>> "{3} {0} {2} {1} {3} {0}".format("be", "not", "or", "to") # => 'to be or not to be'
>>> "{name} is approximately {value}.".format(value=pi, name="π")  #=> 'π is approximately 3.141592653589793.' 
>>> "{foo} {} {bar} {}".format(1, 2, bar=4, foo=3) # 未命名参数按按顺序将字段和参数配对
'3 1 4 2'
>>> fullname = ["Alfred", "Smoketoomuch"] 
>>> "Mr {name[1]}".format(name=fullname) 
'Mr Smoketoomuch' 
>>> print("{pi!s} {pi!r} {pi!a}".format(pi="π")) # `s、r和a`指定分别使用str、repr和ascii进行转换
>>> "The number is {num:b}".format(num=42) 
'The number is 101010' 
>>> "{pi:10.2f}".format(pi=pi)  # 宽度+精度
'      3.14' 
>>> '{:010.2f}'.format(pi) # 在指定宽度和精度的数前面，可添加一个标志. 这个标志可以是零、加号、减号或空格，其中零表示使用0来填充数字.
'0000003.14' 
>>> 'One googol is {:,}'.format(10**100) # 可使用逗号来指出要添加千位分隔符
'One googol is 10,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,00 
0,000,000,000,000,000,000,000,000,000,000,000,000,000,000'
>>> print('{0:<10.2f}\n{0:^10.2f}\n{0:>10.2f}'.format(pi)) # 指定左对齐、右对齐和居中，可分别使用<、>和^
3.14 
    3.14 
        3.14 
>>> print('{0:+.2}\n{1:+.2}'.format(pi, -pi)) 
+3.1 
-3.1 
>>> print('{0: .2}\n{1: .2}'.format(pi, -pi)) # 如果将符号说明符指定为空格，会在正数前面加上空格而不是+
 3.1 
-3.1
>>> "{:#b}".format(42) # 在符号说明符和宽度之间, 它会触发格式转换
'0b101010' 
###
>>> from string import Template 
>>> tmpl = Template("Hello, $who! $what enough for ya?") 
>>> tmpl.substitute(who="Mars", what="Dusty")  # => 'Hello, Mars! Dusty enough for ya?'
### 
>>> from math import e 
>>> f"Euler's constant is roughly {e}."  # 如果变量与替换字段同名，还可使用一种简写。在这种情况下，可使用f字符串——在字符串前面加上f
"Euler's constant is roughly 2.718281828459045." 
###
>>> "The Middle by Jimmy Eat World".center(39, "*")  # 在两边添加填充字符（默认为空格）让字符串居中
'*****The Middle by Jimmy Eat World*****'
>>> title = "Monty Python's Flying Circus"  # 在字符串中查找子串. 如果找到，就返回子串的第一个字符的索引，否则返回-1
>>> title.find('Monty') 
>>> title.find('!!!', 0, 16) # 同时指定了起点和终点或仅起点
-1
>>> dirs = '', 'usr', 'bin', 'env' 
>>> '/'.join(dirs) 
'/usr/bin/env'
>>> '1+2+3+4+5'.split('+') 
['1', '2', '3', '4', '5'] 
# 很多字符串方法都以is打头，如isspace、isdigit和isupper，它们判断字符串是否具有特定的性质（如包含的字符全为空白、数字或大写）
```


```python
-10 // 3=>-4 # 对于整除运算，需要明白的一个重点是它向下圆整结果
 -10 %  3  =>  2 # x y % 等价于x - ((x // y) * y)
1 / 2 => 0.5 # Python2.7的结果是`0`.
1 // 2 => 0
x = input("x: ")  # 获取用户输入
'xxx'.upper() # 转大写
' xxx '.strip() # 删除两边空格, 单边用`lstrip()/rstrip()`
'x y z'.split() # 以空格为分隔符将字符串分拆成多个部分，并将这些部分都存储到一个列表中
str(12) # 将非字符串转换成字符串
':'.join(["a","B"]) # "a:B"
endings = ['st', 'nd', 'rd'] + 2 * ['th']  # => `['st', 'nd', 'rd', 'th', 'th']`
motorcycles = ['honda', 'yamaha', 'suzuki'] 
motorcycles.append('ducati')  # 追加一个元素
motorcycles.insert(0, 'ducati') # 在指定索引前面插入元素
del motorcycles[0]  # 删除元素
popped_motorcycle = motorcycles.pop() # 类似栈的pop()操作
popped_motorcycle = motorcycles.pop(1) # 弹出指定索引位置的元素
motorcycles.remove('ducati') # 根据值来删除列表, 仅删除第一个匹配到的值
motorcycles.clear() # 就地清空列表
motorcycles.pop() # 从列表中删除一个元素, 默认从末尾开始删除或指定要删除元素的索引，并返回这一元素. 
motorcycles.index("honda") # 在列表中查找指定值第一次出现的索引
b = motorcycles.copy() # 复制list, 效果与用a[:]或list(a)类似
[[1, 2], 1, 1, [2, 1, [1, 2]]] .count(1) # =>2.  计算指定的元素在列表中出现了多少次
motorcycles.extend([1,2]) # 能够同时将多个值附加(复制)到列表末尾,
motorcycles.sort() # 就地排序, 调用后motorcycles即变成了sorted列表
 motorcycles.sort(key=len) # 法sort接受两个可选参数：key和reverse, key用于排序的函数, 用该函数的结果作为排序依据. `reverse=True`表示逆向排序.
motorcycles.reverse() # 反向排序
list(reversed(motorcycles))  # 用list将返回的对象(迭代器)转换为列表
sorted(motorcycles) # motorcycles不变, 输出排序后的内容. 该方法支持传递参数`reverse=True`进行反向排序.
len(motorcycles) # 获取列表的长度
list('Hello')  # list() 方法用于将元组转换为列表
tuple([1, 2, 3])  # 将list转为元组
sorted(motorcycles) # 返回排序后的列表
for magician in magicians: # 不能忘记末尾的`:`
    print(magician) 
for value in range(1,5): # 生成一系列的数字, 范围[start, end), 这里即1~4
 print(value) 
numbers = list(range(1,6)) # 用函数list()将range()的结果直接转换为列表
even_numbers = list(range(2,11,3)) # 可指定步长为3
squares = [] # 创建了一个空list
digits = [1, 2, 3, 4, 5, 6, 7, 8, 9, 0]
min(digits)
max(digits)
sum(digits)
[0]*3 = [0] + [0] + [0] = [0, 0, 0] 
squares = [value** for value in range(1,11)] # 列表解析(**不推荐**)将for循环和创建新元素的代码合并成一行, 并自动附加新元素.
print(motorcycles[0:3]) # 创建切片, 此时是通过复制元素的方式创建, 因此修改切片不影响源列表
digits[0:10:2] # => [1, 4, 7, 0], `2`为步长, 可以为负数，即从右向左提取元素, 此时切片的start:end也要互换
digits[:5:-2] # => [0,8],
### 给切片赋值
>>> name = list('Perl') 
>>> name 
['P', 'e', 'r', 'l'] 
>>> name[2:] = list('ar') 
>>> numbers  = [1, 2, 3, 4, 5] 
>>>  numbers[1:4] = [] => [1, 5]
print(motorcycles[:]) # motorcycles[0:len(motorcycles)-1]
print(motorcycles[:4]) # 即motorcycles[0:4]
print(motorcycles[2:]) # 即motorcycles[2:len(motorcycles)-1]
motorcycles[-3:] # 输出最后三个元素
for player in players[:3]:  # 遍历切片
    print(player.title()) 
###
cars = ['audi', 'bmw', 'subaru', 'toyota'] 
for car in cars: 
    if car == 'bmw':  # 不能忘记结尾的冒号
        print(car.upper()) 
    else: 
        print(car.title()) 
###
banned_users = ['andrew', 'carolina', 'david'] 
user = 'marie' 
    if user not in banned_users: 
        print(user.title() + ", you can post a response if you wish.") 
### if-elif-else
age = 12 
if age < 4: 
    print("Your admission cost is $0.") 
elif age < 18: 
    print("Your admission cost is $5.") 
else: 
    print("Your admission cost is $10.") 
### 
alien_0 = {'color': 'green', 'points': 5} # key不可重复, 但仅保留最新的一个
del alien_0['points'] # 删除字典元素

for key, value in alien_0.items():  # 遍历字典
    print("\nKey: " + key) 
    print("Value: " + value)
for name in alien_0.keys(): # 遍历所有key
    print(name.title()) 

for name in sorted(alien_0.keys()): # 按顺序遍历所有键
    print(name.title()) 
for language in alien_0.values():  # 遍历所有值
    print(language.title()) 
for language in set(alien_0.values()):  # set类似于列表, 但每个元素都是唯一的. 在这里可用于剔除重复项
    print(language.title())
### while
current_number = 1 
while current_number <= 5:
    print(current_number) 
    current_number += 1
### 函数
def greet_user(username, other ="王"): # 指定默认值
    """显示简单的问候语""" # 
    print("Hello!" + username) 
 
greet_user("chen", "陈") 
greet_user(username='chen', other='陈') 

### json marshal
import json

numbers = [2, 3, 5, 7, 11, 13]

filename = 'numbers.json'
with open(filename, 'w') as file_object:
    json.dump(numbers, file_object)

### json unmarshal
import json

filename = 'numbers.json'
with open(filename) as file_object:
    numbers = json.load(file_object)
    
print(numbers)
```

## FAQ
### PyCharm 设置和更换 Python 版本
通过`Settings -> Project:xxx -> Project Interpreter`为每个项目选择指定的默认编译器版本. 没有时可从系统环境导入(前提是系统已安装该版本的python), 步骤:
    1. 选择`Project Interpreter`列表框右侧`设置`图标中的添加按钮
    1. 选择tab `System Interpreter`, 选择指定版本即可
### pycharm 导入模块时提示`unresolved reference`
本质均是修改环境变量PYTHONPATH.

共2种方法:
1. File->setting->Project:xxx->Project Interpreter, 选择相应的环境, 再点击后面的图标选择"show all", 选中对应的环境, 再点击右侧工具栏中的"树形"图标, 追加需导入module的上级目录的路径即可(**推荐**).
1. 简单项目只需将需导入module的上级目录的路径设为`source`即可, 具体方法: 找到该目录, 右键选中`Mark Directory as` -> "Sources Root".
### refresh import after pip install in pycharm
File -> Invalidate Caches/Restart...
### `pip的pip.conf`和`pypirc`的区别
`.pypirc`是多个工具使用的文件标准，但不是由pip使用的. 例如，easy_install工具reads that file和twine一样, 它包含在**发布包**时如何访问特定pypi索引服务器的配置. 因此如果不发布包，则不需要.pypirc文件, 即不能使用它为pip配置索引服务器.
但`pip.conf`仅由pip工具使用，**pip从不发布包，而是从中下载包**. 因此，它从不查看.pypirc文件.

对于--index-url和--index开关，它们用于不同的pip命令.
--index-url是处理安装包的几个pip命令中的一个通用开关（pip install，pip download，pip list，和pip wheel），它是一组交换机的一部分（连同--extra-index-url，--no-index，--find-links和--process-dependency-links和一些不推荐的开关），它们一起配置包发现如何工作. url必须指向PEP 503 Simple Repository API位置，默认为https://pypi.org/simple.
--index仅由pip search使用；它只需要这一条信息. 它是单独命名的，因为它应该指向公共搜索Web界面，而不是简单的存储库！对于https://pypi.org，这是https://pypi.org/pypi.
###  SyntaxError: Missing parentheses in call to 'exec'
`sudo python -m pip install Pyro`报错.  Pyro是python2的.

### Python 函数参数前面一个星号（*）和两个星号（**）的区别
单星号(*agrs) : 将所有参数以元组(tuple)的形式导入
星号（**kwargs）: 将参数以字典的形式导入

### setup.py
setuptools 是一个优秀的，可靠的 Pthon 包安装与分发工具.

Python 库打包的格式包括 Wheel 和 Egg. Egg 格式是由 setuptools 在 2004 年引入，而 Wheel 格式是由 PEP427 在 2012 年定义. 使用 Wheel 和 Egg 安装都不需要重新构建和编译，其在发布之前就应该完成测试和构建.

Egg 和 Wheel 本质上都是一个 zip 格式包，Egg 文件使用 .egg 扩展名，Wheel 使用 .whl 扩展名. **Wheel 的出现是为了替代 Egg，其现在被认为是 Python 的二进制包的标准格式**.

### ImportError: No module named Cython.Build


## 用法
- ConfigParser : 解析ini配置文件, 加载多个配置时, 后加载的相同键会覆盖前面的.
- io.BytesIO() : BytesIO实现了在内存中读写bytes(二进制数据), 可像读文件一样读取. 类似的有StringIO, 它用于字符串的操作.
- [open()](https://www.runoob.com/python/file-methods.html)
- os.environ.get("PATH") # os.environ(x [,x]) raises an exception if the environmental variable does not exist
- os.getenv('MYENV', "abc") # os.getenv(x) does not raise an exception ,but returns None. 但有默认值时返回默认值
- os.path.expanduser("~/g/my.conf") # Expand the user's home directory
- SafeConfigParser : SafeConfigParser类实现了ConfigParser相同的接口, 但新增了`set(section, option, value)`(如果给定的section存在，给option赋值；否则抛出NoSectionError异常。Value值必须是字符串（str或unicode）；如果不是，抛出TypeError异常)


# lib
## Pyro
 - [Pyro简单使用(一)](https://www.cnblogs.com/flyingzl/articles/1870799.html)

 >[Pyro4/5 is for python3.x](https://python-parallel-programmning-cookbook.readthedocs.io/zh_CN/latest/chapter5/06_Remote_Method_Invocation_with_Pyro4.htmls)

Pyro即Python Remote Object，类似于Hessian，可以进行远程对象(方法)调用. 不过和Hessian不一样，Pyro利用python本身的pickle模块进行序列化、反序列化，更pythonic，也更方便.

```python
# python2
# server:
import Pyro.core

class JokeGen(Pyro.core.ObjBase):
        def __init__(self):
                Pyro.core.ObjBase.__init__(self)
        def joke(self, name):
                return "Sorry "+name+", I don't know any jokes."

Pyro.core.initServer()
daemon=Pyro.core.Daemon()
uri=daemon.connect(JokeGen(),"jokegen")

print "The daemon runs on port:",daemon.port
print "The object's uri is:",uri

daemon.requestLoop()

# client:
import Pyro.core

# you have to change the URI below to match your own host/port.
jokes = Pyro.core.getProxyForURI("PYROLOC://localhost:7766/jokegen")

print jokes.joke("Irmen")
```


## [multiprocessing](https://docs.python.org/zh-cn/3/library/multiprocessing.html)
multiprocessing和multiprocessing.dummy(复制了 multiprocessing 的 API，不过是在 threading 模块之上包装了一层)是Python下两个常用的多进程和多线程模块

> 其中对于多进程，multiprocessing.freeze_support()语句在windows系统上是必须的，这是因为windows的API不包含fork()等函数. 因此当子进程被创建后，会重新执行一次全局变量的初始化, 这可能就会覆盖掉主进程中已经被修改了的值.

> **主进程会等待直到multiprocessing的子进程退出才退出**.

启动多进程(fork):
```python
import os
import multiprocessing

def foo(i):
    # 同样的参数传递方法
    print("这里是 ", multiprocessing.current_process().name)
    print('模块名称:', __name__)
    print('父进程 id:', os.getppid())  # 获取父进程id
    print('当前子进程 id:', os.getpid())  # 获取自己的进程id
    print('------------------------')

if __name__ == '__main__':

    for i in range(5):
        p = multiprocessing.Process(target=foo, args=(i,)) # 因为args要求是可迭代的
        p.start()
```

### 要在进程之间进行数据共享可以使用Queues、Value/Array(共享内存 by mmap)和Manager(服务进程)这三个multiprocess模块提供的类.

> Queues并发安全, Value、Array不是.

1. 使用Manager共享数据:
```python
# Manager()返回的manager对象提供一个服务进程，使得其他进程可以通过代理的方式操作Python对象。manager对象支持 list, dict, Namespace, Lock, RLock, Semaphore, BoundedSemaphore, Condition, Event, Barrier, Queue, Value ,Array等多种格式
from multiprocessing import Process, Manager

def func(i, dic):
    dic["num"] = 100+i
    print(dic.items())

if __name__ == '__main__':
    dic = Manager().dict()
    for i in range(10):
        p = Process(target=func, args=(i, dic))
        p.start()
        p.join() # join方法将阻塞，直到调用 join() 方法的进程终止
```


1. 使用queues的Queue类共享数据
```python
from multiprocessing import Process, queues

def func(i, q):
    ret = q.get()
    print("进程%s从队列里获取了一个%s，然后又向队列里放入了一个%s" % (i, ret, i))
    q.put(i)

if __name__ == "__main__":
    lis = queues.Queue(20, ctx=multiprocessing)
    lis.put(0)
    for i in range(10):
        p = Process(target=func, args=(i, lis,))
        p.start()
```

### 进程锁
为了防止和多线程一样的出现数据抢夺和脏数据的问题，同样需要设置进程锁. 与threading类似，在multiprocessing里也有同名的锁类RLock，Lock，Event，Condition和 Semaphore.
```python
from multiprocessing import Process, Array, RLock, Lock, Event, Condition, Semaphore
import time

def func(i,lis,lc):
    lc.acquire()
    lis[0] = lis[0] - 1
    time.sleep(1)
    print('say hi', lis[0])
    lc.release()

if __name__ == "__main__":
    array = Array('i', 1)
    array[0] = 10
    lock = RLock()
    for i in range(10):
        p = Process(target=func, args=(i, array, lock))
        p.start() # 非阻塞
```

### 进程池Pool类
进程启动的开销比较大，过多的创建新进程会消耗大量的内存空间. 仿照线程池的做法，我们可以使用进程池控制内存开销.

进程池内部维护了一个进程序列，需要时就去进程池中拿取一个进程，如果进程池序列中没有可供使用的进程，那么程序就会等待，直到进程池中有可用进程为止. maxtasksperchild参数可以控制一条进程最多处理的任务数，当超过这个数量时该进程会退出且pool会重新启动一个新的进程.

进程池中常用的方法：

apply() 同步执行（串行）
apply_async() 异步执行（并行）
terminate() 立刻关闭进程池
join() 主进程等待所有子进程执行完毕. 必须在close或terminate()之后.
close() 调用close后不能添加新的进程
imap() 是进程函数名不变，改变的传递进去的参数，结果是进程返回的结果，有先后顺序


```python
from multiprocessing import Pool
import time

def func(args):
    time.sleep(1)
    print("正在执行进程 ", args)

if __name__ == '__main__':
    p = Pool(5)     # 创建一个包含5个进程的进程池

    for i in range(30):
        p.apply_async(func=func, args=(i,))

    p.close()           # 等子进程执行完毕后关闭进程池
    # time.sleep(2)
    # p.terminate()     # 立刻关闭进程池
    p.join()
```

### Manager
参考:
- [multiprocessing模块实现](http://www.calvinneo.com/2019/11/23/multiprocessing-implement/)

Manager对象会使用一个Server进程来维护需要共享的对象，而**其他进程需要通过Proxy来访问这些共享对象**. 比如当我们调用multiprocessing.Manager()时，实际上创建了一个SyncManager，并调用它的start. 在start中，会创建一个子进程来维护共享对象.

> 使用服务进程的管理器比使用共享内存对象更灵活，因为它们可以支持任意对象类型. 此外，单个管理器可以通过网络由不同计算机上的进程共享. 但是，它们比使用共享内存慢.

> manager对象仅能传播一个可变对象本身所做的修改，如果一个manager.list()对象，管理列表本身的任何更改会传播到所有其他进程，但是如果容器对象内部还包括可修改对象，则内部可修改对象的任何更改都不会传播到其他进程.

sever:
```python
from multiprocessing.managers import BaseManager

BaseManager.register('get_task', callable=gettask) # register的typeid是用于标识特定类型的共享对象的“类型标识符”。这必须是字符串
```

client:
```python
from multiprocessing.managers import BaseManager
BaseManager.register('get_task') # 注册要用的资源(比如函数, Class), 但不能注册服务端没有注册的资源
```

## apscheduler
[APScheduler](https://zhuanlan.zhihu.com/p/46948464)是一个python的第三方库，提供了非常丰富而且方便易用的定时任务接口. 包含四个组件，分别是：
- triggers： 任务触发器组件，提供任务触发方式

    触发器包含调度逻辑，描述一个任务何时被触发，按**日期或按时间间隔或按 cronjob 表达式三种方式触发**. 每个作业都有它自己的触发器，除了初始配置之外，触发器是完全无状态的

    ```python
    # 在2019年11月6日16:30:05
    sched.add_job(my_job, 'date', run_date='2009-11-06 16:30:05', id='my_job_id') # = run_date=datetime(2009, 11, 6, 16, 30, 5)

    # 每两小时执行一次
    sched.add_job(job_function, 'interval', hours=2)
    # 在2012年10月10日09:30:00 到2014年6月15日11:00:00的时间内，每两小时执行一次
    sched.add_job(job_function, 'interval', hours=2, start_date='2012-10-10 09:30:00', end_date='2014-06-15 11:00:00'

    # 在6、7、8、11、12月的第三个周五的00:00, 01:00, 02:00和03:00 执行
    sched.add_job(job_function, 'cron', month='6-8,11-12', day='3rd fri', hour='0-3')
    # 在2014年5月30日前的周一到周五的5:30执行
    sched.add_job(job_function, 'cron', day_of_week='mon-fri', hour=5, minute=30, end_date='2014-05-30')
    ```
- job stores： 任务store组件，提供任务保存方式

    作业存储器指定了作业被存放的位置，默认情况下作业保存在内存，也可将作业保存在各种数据库中，当作业被存放在数据库中时，它会被序列化，当被重新加载时会反序列化。作业存储器充当保存、加载、更新和查找作业的中间人. **在调度器之间不能共享作业存储**.

    ```python
    MemoryJobStore # 默认内存存储
    scheduler.add_jobstore('mongodb', collection='example_jobs')
    scheduler.add_jobstore('redis', jobs_key='example.jobs', run_times_key='example.run_times')
    scheduler.add_jobstore('sqlalchemy', url=url)
    ```
- executors： 任务执行组件，提供任务执行方式

    执行器是将指定的作业（调用函数）提交到线程池或进程池中运行，当任务完成时，执行器通知调度器触发相应的事件
- schedulers： 任务调度组件，提供任务调度方式

    任务调度器，属于控制角色，**通过它配置作业存储器、执行器和触发器，添加、修改和删除任务**。调度器协调触发器、作业存储器、执行器的运行，通常只有一个调度程序运行在应用程序中，开发人员通常不需要直接处理作业存储器、执行器或触发器，配置作业存储器和执行器是通过调度器来完成的.

    常用的有BackgroundScheduler（后台运行）和BlockingScheduler(阻塞式).

    ```python
    from pytz import utc
    from apscheduler.schedulers.background import BackgroundScheduler
    from apscheduler.jobstores.mongodb import MongoDBJobStore
    from apscheduler.jobstores.sqlalchemy import SQLAlchemyJobStore
    from apscheduler.executors.pool import ThreadPoolExecutor, ProcessPoolExecutor


    jobstores = {
        'mongo': MongoDBJobStore(),
        'default': SQLAlchemyJobStore(url='sqlite:///jobs.sqlite')
    }
    executors = {
        'default': ThreadPoolExecutor(20),
        'processpool': ProcessPoolExecutor(5)
    }
    job_defaults = {
        'coalesce': False, # 当由于某种原因导致某个job积攒了好几次没有实际运行（比如说系统挂了5分钟后恢复，有一个任务是每分钟跑一次的，按道理说这5分钟内本来是“计划”运行5次的，但实际没有执行），如果coalesce为True，下次这个job被submit给executor时，只会执行1次，也就是最后这次，如果为False，那么会执行5次（不一定，因为还有其他条件，看后面misfire_grace_time的解释）
        'max_instances': 3, # 同一个job同一时间最多有几个实例在跑
        'misfire_grace_time': 30 # 设想和上述coalesce类似的场景，如果一个job本来14:00有一次执行，但是由于某种原因没有被调度上，现在14:01了，这个14:00的运行实例被提交时，会检查它预订运行的时间和当下时间的差值（这里是1分钟），大于我们设置的30秒限制，那么这个运行实例不会被执行
    }
    scheduler = BackgroundScheduler(jobstores=jobstores, executors=executors, job_defaults=job_defaults, timezone=utc)

    scheduler.add_job(myfunc, 'interval', minutes=2, id='my_job_id')  # 添加任务    
    scheduler.remove_job('my_job_id')  # 删除任务
    scheduler.pause_job('my_job_id')  # 暂定任务
    scheduler.resume_job('my_job_id')  # 恢复任务

    job = scheduler.add_job(myfunc, 'interval', minutes=2)  # 添加任务
    job.modify(max_instances=6, name='Alternate name') # 修改
    scheduler.reschedule_job('my_job_id', trigger='cron', minute='*/5') # 重设

    scheduler.start()
    scheduler.shotdown(wait=True | False)
    scheduler.pause()
    scheduler.resume()

    def my_listener(event):
    if event.exception:
        print('The job crashed :(')
    else:
        print('The job worked :)')

    scheduler.add_listener(my_listener, EVENT_JOB_EXECUTED | EVENT_JOB_ERROR) # [监听](https://apscheduler.readthedocs.io/en/v3.3.0/modules/events.html#module-apscheduler.events)
    ```