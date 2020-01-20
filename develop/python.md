# python

[书单](https://zhuanlan.zhihu.com/p/34378860):
- 入门
    - Python编程：从入门到实践
    - Python编程快速上手
    - Python基础教程（第3版）
- 进阶
    - 流畅的Python
    - Python Cookbook（第3版）中文版
    - 编写高质量代码: 改善Python程序的91个建议

## 环境
### 安装
```
$ sudo apt-get install python3.8
$ sudo apt install python3-pip
```

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

vscode:
- 安装插件`Python`(ms-python.python), 但它的智能提示和跳转无法在pygame上生效, 但pycharm社区版正常.

其他安装:
```
$ python -m pip install -U pygame --user # 安装pygame, 无法安装成功, 缺少依赖. 但系统自带的python3.7可以安装成功并运行demo
$ python -m pygame.examples.aliens # 运行pygame demo
$ python -m pip  show pygame # 查找安装位置
```

> 可使用`apt-file search "sdl-config"`查找pygame的依赖

## 变量和简单数据类型
变量名只能使用字母,数字和下划线, 且不能以数字开头. 不能使用python保留用于特殊用途的单词(关键字和函数名)作为变量名.

python中用引号(单引号/双引号)包裹的都是字符串.

方法是Python可对数据执行的操作.

python3的print是函数需要括号, Python2的print括号则可有可无.

Python使用`**`表示乘方运算.
Python将带小数点的数字都称为浮点数.
Python的`"""`(docstring, 文档字符串)类似golang的"``", 输出时原样输出, Python使用它们来生成有关程序中函数的文档.
Python使用`#`作为注释标识, 但也可将`"""`用作多行注释.

在Python中，用方括号`[]`来表示列表(其他语言的数组)，并用逗号来分隔其中的元素.  **list中元素的类型可以不同**.
Python为反向访问一个列表元素提供了一种特殊语法, 即通过将索引指定为`-n`, `-1`表示最后一个元素.

Python idle中连续两个回车表示结束.

Python根据缩进来判断代码行与前一个代码行的关系.

列表的部分元素——Python称之为切片, 它通过复制来创建.
创建切片时需指定要使用的第一个元素和最后一个元素的索引.

Python将不能修改的值称为不可变的，而不可变的列表被称为元组, 其用圆括号包裹, 想修改元组中的元素时必须给存储元组的变量重新赋值.
Python的布尔表达式的结果要么为True，要么为False.
Python并不要求if-elif结构后面必须有else代码块. 在有些情况下，else代码块很有用; 而在其他一些情况下，使用一条elif语句来处理特定的情形更清晰.
在if语句中将列表名用在条件表达式中时，Python将在列表至少包含一个元素时返回True, 并在列表为空时返回False.
在Python中，字典是一系列键—值对, 用`{}`包裹, 每个键都与一个值相关联，可以使用键来访问与之相关联的值, Python可将任何Python对象用作字典中的值 .
遍历字典时, Python会默认遍历所有的键.

条件测试:
- and : 其他语言的`&&`
- or : 其他语言的`||`
- in : 判断特定的值是否已包含在列表中
- not in : 检查特定值是否不包含在列表中

函数`input("提示")`让程序暂停运行，等待用户输入一些文本, 获取用户输入后，Python将其存储在一个变量中，以方便之后使用.
函数`int()`将数字的字符串表示转换为数值表示.

> Python 2.7也包含函数input()，但它将用户输入解读为Python代码，并尝试运行它们, 其用函数`raw_input()`来提示用户输入.

for循环是一种遍历列表的有效方式，但在**for循环中不应修改列表**，否则将导致Python难以跟踪其中的元素. 要在遍历列表的同时对其进行修改，可使用while循环.

## 函数
使用`def`定义函数.
**关键字实参**是传递给函数的名称—值对, 此时无需考虑函数调用中的实参顺序，还清楚地指出了函数调用中各个值的用途.
Python支持指定默认值, 给形参指定默认值时，等号两边不要有空格.

将列表传递给函数后，函数就可对其进行修改, 在函数中对这个列表所做的任何修改都是永久性的.
向函数`function_name(list_name[:])`传递列表的副本而不是原件; 这样函数所做的任何修改都只影响副本，而丝毫不影响原件.
`def make_pizza(*toppings): `的形参名*toppings中的星号让Python创建一个名为toppings的空元组, 并将收到的所有值都封装到这个元组中, 以实现像函数传递任意数量的实参.
如果要让函数接受不同类型的实参，必须在函数定义中将接纳任意数量实参的形参放在最后, 比如`def make_pizza(size, *toppings): `.
`def build_profile(first, last, **user_info):`的形参**user_info中的两个星号让Python创建一个名为user_info的空字典，并将收到的所有名称—值对都封装到这个字典中.
import语句允许在当前运行的程序文件中使用模块中的代码.

只需编写一条import语句并在其中指定模块名，就可在程序中使用该模块中的所有函数, 调用形式如下:`module_name.function_name()`.
导入模块中的特定函数：`from module_name import function_0, function_1, function_2`, 若使用这种语法，调用函数时就无需使用句点.
导入时给函数指定别名：`from module_name import function_name as fn`
导入时给模块指定别名：`import module_name as mn`
导入模块中的所有函数: `from pizza import *`, 此时可通过名称来调用每个函数，而无需使用句点表示法. **不推荐**: 如果模块中有函数的名称与你的项目中使用的名称相
同, 可能导致意想不到的结果, 因为Python可能遇到多个名称相同的函数或变量, 进而被覆盖.

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
调用都自动传递实参self，它是一个指向实例本身的引用，让实例能够访问类中的属性和方法, 同时self会自动传递，因此不需要显示传递它.

通过实例访问的变量称为属性, 比如`self.name = name`的`self.name`.

> 在Python 2.7中创建类时，需要在类名后面的括号内包含单词object.

类中的每个属性都必须有初始值，在方法`__init__()`内指定这种初始值是可行的；如果对某个属性这样做了，就无需包含为它提供初始值的形参.

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

创建子类时，父类必须包含在当前文件中，且位于子类前面.
定义子类时，必须在括号内指定父类的名称, 且子类方法`__init__()`接受创建父类实例所需的信息.

`super()`是一个特殊函数，帮助Python将父类和子类关联起来. 父类也称为超类（superclass）, 名称super因此而得名.

可在子类中定义一个与要重写的父类方法同名的方法, 这样, Python将不会考虑这个父类方法，而只关注你在子类中定义的相应方法.

属性和方法清单以及文件都越来越长时, 可能需要将类的一部分作为一个独立的类提取出来, 这样可以将大型类拆分成多个协同工作的小类, 比如上面的`Battery类`.

从一个模块中导入多个类:`from car import Car, ElectricCar `.
导入整个模块:` import car`
导入模块中的所有类:`from module_name import * `, 不推荐, 原因和导入module相同.

在属性或方法名前面添加两个下划线`__`, 成员私有化的作用，确保外部代码不能随意修改对象内部的状态，增加了代码的安全性.

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

异常是使用try-except代码块处理的, 而依赖于try代码块成功执行的代码都应放到else代码块中.

Python有一个pass语句，可在代码块中使用它来让Python什么都不要做.

```python
try: 
    print(5/0) 
except ZeroDivisionError: 
    print("You can't divide by zero!") 
###
filename = 'alice.txt' 
try: 
    with open(filename) as f_obj: 
        contents = f_obj.read() 
except FileNotFoundError: 
    msg = "Sorry, the file " + filename + " does not exist." 
    print(msg) 
else: 
    print("OK") 
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

### demo
解构(唯一的前提是变量的数量必须跟序列元素(任何的序列或可迭代对象)的数量是一样的):
```python
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
```


```python
'xxx'.upper() # 转大写
' xxx '.strip() # 删除两边空格, 单边用`lstrip()/rstrip()`
'x y z'.split() # 以空格为分隔符将字符串分拆成多个部分，并将这些部分都存储到一个列表中
str(12) # 将非字符串转换成字符串
motorcycles = ['honda', 'yamaha', 'suzuki'] 
motorcycles.append('ducati')  # 追加元素
motorcycles.insert(0, 'ducati') # 在指定索引前面插入元素
del motorcycles[0]  # 删除元素
popped_motorcycle = motorcycles.pop() # 类似栈的pop()操作
popped_motorcycle = motorcycles.pop(1) # 弹出指定索引位置的元素
motorcycles.remove('ducati') # 根据值来删除列表, 仅删除第一个匹配到的值
motorcycles.sort() # 排序, 调用后motorcycles即变成了sorted列表
motorcycles.reverse() # 反向排序
sorted(motorcycles) # motorcycles不变, 输出排序后的内容. 该方法支持传递参数`reverse=True`进行反向排序.
len(motorcycles) # 获取列表的长度
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
squares = [value** for value in range(1,11)] # 列表解析(**不推荐**)将for循环和创建新元素的代码合并成一行, 并自动附加新元素.
print(motorcycles[0:3]) # 创建切片, 此时是通过复制元素的方式创建, 因此修改切片不影响源列表
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