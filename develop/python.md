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
## 变量和简单数据类型
变量名只能使用字母,数字和下划线, 且不能以数字开头. 不能使用python保留用于特殊用途的单词(关键字和函数名)作为变量名.

python中用引号(单引号/双引号)包裹的都是字符串.

方法是Python可对数据执行的操作.

python3的print是函数需要括号, Python2的print括号则可有可无.

Python使用`**`表示乘方运算.
Python将带小数点的数字都称为浮点数.
Python使用`#`作为注释标识.

在Python中，用方括号`[]`来表示列表(其他语言的数组)，并用逗号来分隔其中的元素.
Python为反向访问一个列表元素提供了一种特殊语法, 即通过将索引指定为`-n`, `-1`表示最后一个元素.

Python idle中连续两个回车表示结束.

Python根据缩进来判断代码行与前一个代码行的关系.

列表的部分元素——Python称之为切片
创建切片时需指定要使用的第一个元素和最后一个元素的索引.

### demo
```python
'xxx'.upper() # 转大写
' xxx '.strip() # 删除两边空格, 单边用`lstrip()/rstrip()`
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
```