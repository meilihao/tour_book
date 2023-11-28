# python

缺点:
1. json序列化/反序列化没有go简单
1. 不能静态检查
1. Python try...catch有时不容易定位错误位置

[书单](https://zhuanlan.zhihu.com/p/34378860):
- 入门
    - Python基础教程（第3版）
    - Python编程：从入门到实践
    - Python编程快速上手
    - [python3-cookbook](https://python3-cookbook.readthedocs.io/zh_CN/latest/index.html)
- 进阶
    - Python高性能编程
    - 流畅的Python
    - Python Cookbook（第3版）中文版
    - 编写高质量代码: 改善Python程序的91个建议


## 环境
[virtualenv](https://www.liaoxuefeng.com/wiki/1016959663602400/1019273143120480)用于建立一个隔离的python运行环境, 一个专属于项目的python环境, 比如针对不同的python版本. 因此用virtualenv 来保持一个干净的环境非常有用.

### 安装
```
$ sudo apt-get install python3.8
$ sudo apt install python3-pip
$ sudo apt install python-pip # pip2 for python2.7
$ sudo pip install robotframework==2.8.7 # 安装指定版本
# --- ubuntu20.04开始默认不安装py2.7, 且没有python-pip, 因此要使用get-pip.py安装pip2, 或通过python2-xxx安装package
$ curl https://bootstrap.pypa.io/get-pip.py --output get-pip.py
$ sudo python2 get-pip.py
```

`pipdeptree -p xxxx [-r]`显示xxx所有的依赖包及其子包(`-r`反向查询即谁依赖xxx), 其他的有`pip show xxx`

`pip install wheel`后安装依赖时也会构建出相应的whl包, 此时[`pip install --no-cache-dir`](https://stackoverflow.com/questions/35169608/when-does-pip-install-build-a-wheel)可禁用构建whl包.

> [pip安装, 未测试](https://pip.pypa.io/en/stable/installing/)

> 作为过渡, python3有很多特性被移植到了python2.7(将于2020.1.1终止支持), 因此如果程序可在python2.7运行就可通过python3自带的转换工具2to3(`python -m pip install  2to3`/`sudo apt install 2to3`)无缝迁移到Python3.

python解析器:
1. CPython : 使用C实现的解析器, 最常用的解析器, 通常说的python解析器就是指它.
1. PyPy : 使用Python实现的解析器

> IPython是一个python交互shell，它比默认的python shell更易于使用. 它支持自动变量完成、自动缩进、bash shell命令，并且内置了许多有用的函数和函数.

### 切换python版本
- [update-alternatives](https://blog.csdn.net/White_Idiot/article/details/78240298)
- alias
    ```
    $ vim ~/.bashrc
    $ alias python='/usr/bin/python3.8'
    ```

### pip离线部署
```bash
# --- **`pandas` 和`-r requirements_full.txt` 不要混用, 否则可能会因计算出的依赖版本不同导致未知错误**
pip download -d ./packages pandas
pip download -d ./packages -r requirements_full.txt

pip install --no-index --find-links=./packages pandas
pip install --no-index --find-links=./packages -r requirements.txt

pip freeze > requirements.txt
pip download -r requirements.txt -d /packs -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
```

> 其中`--no-index`代表忽视pip 忽视默认的依赖包索引, 即忽略已安装的包, 而是仅从`--find-links`获取; `--find-links`代表从指定的目录寻下找离线包

[`pip install --download`已被废弃](https://github.com/pypa/pip/pull/3085), `pip download`后安装可能会报错, 可将依赖分批次进行安装, 尽可能避免该问题.

原因是pip download不会考虑pip installed的包, 或者是因为pip download没有满足依赖的包它就重新下载了.

有些包下载时就会检查依赖关系, 因为后一个包依赖前一个包, 而前一个包还没有安装, 此时会报错, 因此它们离线部署时必须下一个安装一个, 比如PasteDeploy和PasteScript

### style
[git hook, 脚本是Python3的, 需要注意pycodestyle的路径, 这里是`args = ['/home/ubuntu/.local/lib/python2.7/site-packages/pycodestyle.py']`](https://github.com/cbrueffer/pep8-git-hook/blob/master/pre-commit)

```bash
$ pycodestyle --ignore E501 *.py # pycodestyle原名pep8, python code style checker
$ autopep8 --in-place --aggressive --aggressive *.py # style 修正
```

pycharm设置pep8: [Pycharm配置autopep8教程，让Python代码更符合pep8规范](https://segmentfault.com/a/1190000005816556)

### 利用Cython将Python项目编译为`.so`
参考:
- [尝试利用Cython将Python项目转化为单个.so](https://paper.seebug.org/1139/)

编译Cython module的主要过程:
1. Cython compiler将.py/.pyx文件编译为C/C++文件
2. C compiler再将C/C++编译为.so(windows 为.pyd)

```bash
$ mkdir test && cd test
$ cat fn.py 
def say_hello_to(name):

    print("Hello %s!" % name)
$ cat main.py 
import fn

fn.say_hello_to('xxx')
$ cat setup.py # 包含了构建指令
from distutils.core import setup
from Cython.Build import cythonize

setup(name='test',
   ext_modules=cythonize("fn.py"))
$ python3 setup.py build_ext --inplace # --inplace 参数让 Cython 在当前目录中构建编译模块，而不是在一个独立的构建
目录中
$ python3 main.py
```

[`Cython.Build.cythonize(module_list, exclude=None, nthreads=0, aliases=None, quiet=False, force=False, language=None, exclude_failures=False, show_all_warnings=False, **options)`参数](https://blog.csdn.net/jay_yxm/article/details/106679075):
- module_list : 作为模块列表，传递一个 glob 模式，一个 glob 模式列表或一个扩展对象列表。后者允许您通过常规 distutils 选项单独配置扩展。您还可以传递具有 glob 模式作为其源的扩展对象。然后，cythonize 将解析模式并为每个匹配的文件创建扩展的副本
- exclude : 当将glob模式传递作为module_list时，可以通过将某些模块名称传递到exclude选项中来显式排除某些模块名称
- nthreads : 并行编译的并发构建数（需要multiprocessing模块）
- alias
- quiet : 如果为True，则Cython在编译过程中不会打印错误，警告或状态消息。
- force : 强制重新编译Cython模块，即使时间戳不表明需要重新编译也是如此。
- language : 要全局启用C ++模式，可以通过language='c++'。否则，这将基于编译器指令在每个文件级别确定。这仅影响基于文件名找到的模块。传入的扩展实例cythonize()将不会更改。建议使用编译器指令而不是此选项。# - distutils: language = c++
- exclude_failures :对于广泛的“尝试编译”模式，该模式将忽略编译失败并仅排除失败的扩展，请通过exclude_failures=True。请注意，这仅对编译.py文件有意义，这些文件也可以不经编译而使用。
- show_all_warnings : 默认情况下，并非所有Cython警告都会被打印。设置为true以显示所有警告。
- annotate : 如果True，将为每个编译的.pyx或.py文件生成一个 HTML 文件。与普通的 C 代码相比，HTML 文件指示了每个源代码行中的 Python 交互量。它还允许您查看为每行 Cython 代码生成的 C / C ++代码。当优化速度函数以及确定何时 释放 GIL时，此报告非常有用
- compiler_directives -允许集合中的编译器指令setup.py是这样的：compiler_directives={'embedsignature': True}

注意:
- [celery 4.3开始才支持Cythonized Tasks](https://docs.celeryproject.org/en/stable/history/whatsnew-4.3.html)

### 其他
```sh
$ pip --version
$ pip config list -v # 查看正在使用的pip conf. ubuntu 16.04的pip没有config子命令.
$ pip config set global.index-url https://mirrors.aliyun.com/pypi/simple # 推荐, 避免存在多份pip conf导致有效配置被覆盖
$ pip config set global.trusted-host "mirrors.aliyun.com"
$ python -m pip -V # 检查是否安装pip成功
$ mkdir -p ~/.pip
$ vim ~/.config/pip/pip.conf # [为pip换源](https://blog.csdn.net/xuezhangjun0121/article/details/81664260), 会用到sudo时建议添加到`/etc/pip.conf`. `~/.pip/pip.conf`有时失效, 可能与pip版本有关
[global]
index-url = https://pypi.tuna.tsinghua.edu.cn/simple
$ vim ~/.config/pip/pip.conf # aliyun pip mirror
[global]
trusted-host =  mirrors.aliyun.com
index-url = https://mirrors.aliyun.com/pypi/simple
$ python -m site # pip 软件包的安装位置
$ sudo vim /usr/lib/python[2.7|3.8]/site.py # 这里也支持修改 USER_SITE, USER_BASE 
ENABLE_USER_SITE = False # 将该值设置为 False 即可, 顺便可到`/home/${USER}/.local/lib`下清理已下载的package
$ pip install -r requirements.txt
$ pip freeze > requirements.txt
$ pip install pip-20.x.tar.gz # 升级pip
$ pip install setuptools== # 查询可用的软件版本
$ pip show setuptools # 查看setuptools信息
$ pip3 show -f pycryptodome # 显示pycryptodome包的文件
$ python3.9 -m pip3 install cython # 为指定版本的python安装包
$ pipdeptree # 显示依赖树
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

python3.8/3.9源码安装:
```bash
# tar -xf Python-3.9.10.tgz
# cd Python-3.9.10
# ./configure --enable-optimizations
# make
# make install # 有些blog这里使用`make altinstall`. make install会执行commoninstall、bininstall、maninstall三个过程，而make altinstall只执行commoninstall过程即不在`/usr/local/bin`下创建软链和安装man手册
# ln -s /usr/local/bin/python3.9 /usr/bin/python3
# ln -s /usr/local/bin/pip3.9 /usr/bin/pip3
# echo "/usr/local/lib" > /etc/ld.so.conf.d/python3.9.conf
# python3 -V
# pip3 # 会报`subprocess.CalledProcessError: Command '('lsb_release', '-a')' returned non-zero exit status 1`, `mv /usr/bin/lsb_release /usr/bin/lsb_release.bak`可解决但会丢失lsb_release命令. 原因: python路径下缺少 'lsb_release.py' 模块, 最佳解决方法: 1. 查找到lsb_release模块所在的目录: `find / -name 'lsb_release.py'`; 2. 将其复制到设置python3.8的系统模块加载位置，也就是报错处subprocess.py所在的目录`cp  /usr/lib/python3/dist-packages/lsb_release.py /usr/local/lib/python3.8/`
```

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
1. python一切皆对象. 每个对象都有标识(id), 类型, 值. 因为该机制可以认为所有原生对象(非 C、Cython 等扩展)都在`堆`上分配

    > 每个 Python 对象都带有一个引用计数，记录当前正在引用它的值的数量.
    > 内置函数`dir()`可获取对象的所有属性.
    > 内置`help()`可查看对象属性的方法doc, 比如`help(''.replace)`

    ```python
    >>> a = 3 # 内存形式: a(id: 8791360980704. 在栈) 引用了 object(id: 8791360980704, type: int, value: 3. 在堆). 因为object包含类型(因此python是强类型语言), 因此python变量不需要显示声明类型, 可由解析器推导. 事实上,每个对象都很 “重”。即便是简单的整数类型,都有一个标准对象头,保存类型指针和引用计数等信息。如果是变长类型(比如 str、list 等),还会记录数据项长度,然后才是对象状态数据. object大小可用`sys.getsizeof(x)`查看
    >>> a
    3
    >>> id(a) # id() 函数返回对象的唯一标识符，标识符是一个整数. CPython 中 id() 函数用于获取对象的内存地址. 这意味着它只能保证在某个时间,在所有存活对象里是唯一的, 但不保证整个进程生命周期內是唯一的,因为内存地址会被复用. 如此,ID 也就不适合作为全局身份标识
    8791360980704
    >>> type(a)
    <class 'int'>
    >>> a = "string" # Python是动态类型，即在编程中允许随意改变变量的类型，这个过程称为 "绑定（Binding）". "绑定"只存在于动态类型语言中；对于静态语言如C，是"强制类型转换".
    >>> a
    'string'
    >>> isinstance(1, int) # 判断实例是否属于特定类型
    True
    >>> isinstance(int, type) # 所有类型对象都是 type 的实例,这与继承关系无关
    True
    >>> x = 1234
    >>> y = x # 总是按引用用传递
    >>> x is y # is用于判断是否引用用同一一对象, 这里也可用id(x)=id(y)来判断
    True
    ```

    在 python 中赋值语句总是建立对象的引用值，而不是复制对象. 因此，python 变量更像是指针，而不是数据存储区域.

    在静态编译型语言里,变量是一段有特定格式的可读写/只读内存,变量名则是这段内存的别名, 其编译器会使用直接或间接寻址方式替代作为变量名符号, 也就是说变量名不参与执行过程,是可被剔除的. 但在 Python 这类动态语言里,名字和对象是两个实体. 名字(没有类型)不但要分配内存,还会介入实际执行过程, 甚至可以说,名字才是其动态执行模型的基础.

    名字空间是上下文环境里专门用来存储名字和目标引用关联的容器即作用域. 名字空间默认使用 dict 数据结构,由多个键值对(key/value)组成。每个 key 总是唯一, 即变量名的字符串形式.

    对 Python 而言,每个模块(源码文件)都有一个全局名字空间(globals)。而根据代码作用域,又有本地名字空间(locals)一说。如果直接在模块级别执行,那么当前名字空间(`locals()`)和全局名字空间(`globals()`)相同。但在某个函数內,当前名字空间就专指函数执行栈帧(stack frame)作用域.

    > globals 总是指向所在模块名字空间,而 locals 则指向当前作用域环境

### [装饰器](https://www.liaoxuefeng.com/wiki/1016959663602400/1017451662295584)
- @try_except

### gc
python总是通过名字来完成 “引用传递”(pass-by-reference). 名字关联会增加计数,反之减少, 如删除全部名字关联,那么该对象引用计数归零,会被系统自动回收. 这就是默认的引用计数垃圾回收机制(Reference Counting Garbage Collector).

引用计数可通过`sys.getrefcount(b)`查看. `del b`会删除名称与对象的关联, 因此计数会-1.

名字与目标对象关联构成强引用关系,会增加引用计数,会影响其生命周期.

弱引用(weak reference)就是减配版,在保留引用的前提下,不增加计数,也不阻止目标被回收。不过,并不是所有类型都支持弱引用. 弱引用，与强引用相对，是指不能确保其引用的对象不会被垃圾回收器回收的引用。**一个对象若只被弱引用所引用，则可能在任何时刻被回收**。弱引用的主要作用就是**减少循环引用，减少内存中不必要的对象存在的数量**.

针对循环引用python有一套专门用来处理循环引用的垃圾回收器作为补充, 简称 gc, 可用`gc.disable(), 关闭gc; gc.enable(), 启用gc; gc.collect(), 主动启动执行一次回收操作,循环引用用对象被正确回收`来控制.

在python进程启动时,gc 默认被打开,并跟踪所有可能造成循环引用的对象。相比引用计数,gc 是一种**延迟回收方式**。只有当内部预设的阈值条件满足时,才会在后台启动。虽然可忽略该条件,强制执行回收,但不建议频繁使用.

> CPython 使用引用计数,PyPy 却使用了标记清理.

### 编译
源码须编译成字节码(byte code, 与平台无关)指令后,才能由解释器(interpreter)解释执行. 正常情况下,源码文件在被导入(import)时完成编译. 这也是 Python 性能为人诟病的一个重要原因.

Python 3 使用专门的 __pycache__ 目录保存字节码缓存文件(.pyc)。这样在程序下次启动时,可避免再次编译,以提升导入速度。标准 pyc 文件大体上由两部分组成。头部存储有编译相关信息,在启动时,可判断源码文件是否被更新,是否需要重新编译.

> Lib/importlib/text_bootstrap_external.py 文件里的 _code_to_bytecode 代码,可看到字节码头部信息构成.
> 字节码最后由python虚拟机(pvm, python virtual machine)执行.

除内置 compile 函数外,标准库还有编译源码文件的相关模块:
```python
>>> import py_compile, compileall
>>> py_compile.compile("main.py")
'__pycache__/main.cpython-36.pyc'
>>> compileall.compile_dir(".")
Listing '.'...
Compiling './test.py'...
True
```

### 转换
数值转换:

    内置多个函数将整数转换为指定进制的字符串,反向操作用 int. 它默认识别为十进制, 会忽略空格、制表符等多余字符. 如指定进制参数,还可省略字符串前缀.
    ```python
    >>> bin(100)
    '0b1100100'
    >>> oct(100)
    '0o144'
    >>> hex(100)
    '0x64'
    >>> int("0x64", 16)
    100
    >>> int("64", 16) # 省略略进制前缀
    100
    >>> int(" 100 ") # 忽略略多余空白白字符
    100
    >>> b = x.to_bytes(2, sys.byteorder) # sys.byteorder 获取当前系统字节序
    >>> b.hex()
    '3412'
    >>> hex(int.from_bytes(b, sys.byteorder))
    '0x1234'
    ```

    将整数或字符串转换为浮点数很简单,且能自动处理字符串内正负符号和空白符。只是超出有效精度时,结果与字符串内容存在差异.
    反过来,将浮点数转换为整数时,有多种不同方案可供选择。通过`from math import trunc, floor, ceil`可截掉(int, trunc)小数部分,或分别往大(ceil)、小(floor)方向取临近整数.

### 运算符
单斜线称作 True Division,无论是否整除,总是返回浮点数; 而双斜线 Floor Division 会截掉小数部分,仅返回整数结果.

格式化: f-strings

Python 3 对运算符做了些调整:
- 移除 `<>`,统一使用 `!=` 运算符
- 移除 cmp 函数,自行重载相关运算符方法
- 除法 `/` 表示 True Division,**总是返回浮点数**
- 不再支持反引号 repr 操作,调用同名函数
- 不再支持非数字类型混合比较,可自定义相关方法
- 不再支持字典相等以外的比较操作

`T if X else F : 当条件 X 为真时,返回 T,否则返回 F. 等同 X ? T : F`

## 变量和简单数据类型
变量名只能使用字母,数字和下划线, 且不能以数字开头. 不能使用python保留用于特殊用途的单词(关键字和函数名, 在idle中可用`help() + keywords`查看)作为变量名.

如果想要在函数内使用全局作用域的变量时，需要加上global修饰符, 即提示python 解释器，表明被其修饰的变量是**全局变量(即位于模块文件内顶层的变量名)**.

python内置的基本数据类型: 数字(int, float), 序列(str, bytes(只读), bytearray, list, tuple), 映射(dict), 集合(set, frozenset)等几类.

> 事实上，很多文章推崇另外的一种方法来使用全局变量：使用单独的global文件.

> Python 的全局变量是模块 (module) 级别的.

> 双下划线开头和结尾的名称通常具有特殊含有, 尽量避免这种写法.

> del可删除变量(对象还在, 但对象没有引用后会被gc回收).

> `x=y=123` => `x=123;y=123`: 链式赋值, 用于将一个对象赋值给多个变量. python也支持类似go的`a,b,c=4,5,6`和`a,b=b,a`

> **python不支持常量**.

> Python 3 将 2.7 里的 int、long 两种整数类型合并为 int,默认采用变长结构(支持无穷大), 按需分配内存. 对于较长的数字,为了便于阅读, 可用下划线分隔,且不限分隔位数. 内置函数hex(), oct(), bin()可把一个整数转换成对应进制的字符串, 而int(str,base)可把进制字符串转成一个整数.

> 默认 float 类型存储双精度(double)浮点数, 可通过 `float.hex()` 方法输出实际存储值的十六进制格式字符串. 相比 float 基于硬件的二进制浮点类型,decimal.Decimal 是十进制实现,最高可提供 28 位
有效精度, 能准确表达十进制数和运算,**不存在二进制近似值问题**. python也支持分数`Fraction(x, y)`

以下划线开头的名字,有特殊含义:
- 模块成员以单下划线开头(_x),属私有成员,不会被`import *`语句导入
- 类型成员以双下划线开头,但无结尾(__x),属私有成员,会被自动重命名
- 以双下划线开头和结尾(__x__),通常是系统成员,避免使用
- 交互模式(shell)下,单下划线(_)返回最后一个表达式结果

Python 语言规范里并没有枚举(enum)定义,而是采用标准库`import enum`实现, 并没有规定枚举值就必须是整数, 枚举类型内部以字典方式实现,每个枚举值都有 name 和 value 属性.

python中用引号(单引号/双引号)包裹的都是字符串. `hex(ord("汉"))`可获得code point.

Python2 和 Python3 中默认编码的差异:
- 脚本字符编码：就是解释器解释脚本文件时使用的编码格式，python2以 ASCII 为默认编码,如果源码里出现 Unicode 字符,就会导致其无法正确解析, 可以通过 `# -*- coding: utf-8 -*-` 显式指定; Python 3 默认编码已改成 UTF-8
- 解释器字符编码：解释器内部逻辑过程中对 str 类型进行处理时使用的编码格式
- Python2 中默认把脚步文件使用 ASCII 来处理(历史原因)
- Python2 中字符串除了 str 还有 Unicode，可以用 decode 和 encode 相互转换
- Python3 中默认把脚步文件使用 UTF-8 来处理
- Python3 中文本字符和二进制分别使用 str 和 bytes 进行区分，也是使用 decode 和 encode 进行相互转换

> 乱码的终极原因就是：对同一个字符串的 encode 和 decode 编码格式不一致.

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
创建切片时需指定要使用的第一个元素和最后一个元素的索引. 完整的切片操作由三个参数构成: `start : stop : step`.

> 返回切片时会创建新列表对象,并复制相关指针数据到新的数组, **除部分引用目标相同外,对列表自身的修改(插入、删除等)互不影响**.
> 编译器将 “+=” 运算符处理成 INPLACE_ADD 操作,也就是修改原数据,而非新建对象, 因此`a = [1, 2];b = a`时`a += [3, 4]`和`a = a + [3, 4]`操作的结果分别是`b为[1, 2, 3, 4]`和`b为[1, 2]`.
> tuple支持与列表类似的运算符操作,但没有 INPLACE,**总是返回新对象**.
> 切片索引负值表示从右至左反向步进. 以切片方式进行序列局部赋值,相当于先删除,后插入.

Python将不能修改的值称为不可变的，而不可变的列表被称为元组, 其用圆括号包裹, 想修改元组中的元素时必须给存储元组的变量重新赋值.
元组作用:
1. 用作映射中的键（以及集合的成员），而列表不行
1. 有些内置函数和方法返回元组

数组(array)与列表、元组的区别在于:**元素单一类型和内容嵌入**.

Python的布尔表达式的结果要么为True，要么为False. 布尔转换时,数字零、空值(None)、空序列和空字典被视作 False,反之为 True.
Python并不要求if-elif结构后面必须有else代码块. 在有些情况下，else代码块很有用; 而在其他一些情况下，使用一条elif语句来处理特定的情形更清晰.
在if语句中将列表名用在条件表达式中时，Python将在列表至少包含一个元素时返回True, 并在列表为空时返回False.
在Python中，字典是一系列键—值对, 用`{}`包裹, 每个键都与一个值相关联，可以使用键来访问与之相关联的值, Python可将任何Python对象用作字典中的值.
遍历字典时, Python会默认遍历所有的键. 字典的`update()`方法用于更新字典中的键/值对, `pop(key)`可删除key.
支持使用`dict(zip([...],[...]))`创建dict, `zip()`会将对象中对应的元素打包成一个个元组，然后返回由这些元组组成的列表.

> zip 方法在 Python 2 和 Python 3 中的不同：在 Python 3.x 中为了减少内存，zip() 返回的是一个对象. 如需展示列表，需手动 list() 转换.

字典(dict)是内置类型中唯一的映射(mapping)结构,基于哈希表存储键值对数据. 值(value)可以是任意数据,但主键(key)必须是可哈希类型. 常见的可变类型,如列
表、集合等都不能作为主键使用, 可通过`hash(xxx)`来验证key.

条件测试:
- and : 其他语言的`&&`
- or : 其他语言的`||`
- in : 判断特定的值是否已包含在列表中
- not in : 检查特定值是否不包含在列表中

函数`input("提示")`让程序暂停运行，等待用户输入一些文本, 获取用户输入后，Python将其存储在一个变量中，以方便之后使用.
函数`int()`将数字的字符串表示转换为数值表示.

> Python 2.7也包含函数input()，但它将用户输入解读为Python代码，并尝试运行它们, 其用函数`raw_input()`来提示用户输入.

> any(iterable) 函数用于判断给定的可迭代参数 iterable 是否全部为 False，则返回 False，如果有一个为 True，则返回 True.

> 函数总是返回一个结果: 即便什么都不做,也会返回 None.

> 在函数内访问变量,会以特定顺序依次查找不同层次作用域(LEGB: locals -> nonlocal(外层函数) -> globals -> builtins)

> 除非调用`locals()`,否则函数执行期间,根本不会创建所谓名字空间字典. 也就是说,其返回的字典是按需延迟创建(by 复制).

for循环是一种遍历列表的有效方式，但在**for循环中不应修改列表**，否则将导致Python难以跟踪其中的元素. 要在遍历列表的同时对其进行修改，可使用while循环.


### 迭代器/生成器
迭代是指重复从对象中获取数据,直至结束. 迭代器有两个基本的方法：iter() 和 next(), 迭代器只能往前不会后退. 而所谓迭代协议,概括起来就是用 __iter__ 方法返回一个实现了 __next__ 方法的迭代器对象.

实现 __iter__ 方法,表示目标为可迭代(iterable)类型. 字符串，列表或元组对象都可用于创建迭代器.

```python
M=[[1,2,3],[4,5,6],[7,8,9]]                                             
G=(sum(row) for row in M)                                              
next(G)                                                                
6
```

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

复制还分浅拷贝(shallow copy)和深度拷贝(deep copy)两类. **对于对象内部成员,浅拷贝仅复制名字引用,而后者会递归复制所有引用成员**.

```python
>>> class X: pass
>>> x = X() # 创建实例例
>>> x.data = [1, 2] # 成员 data 引用一个列列表

>>> import copy
>>> x2 = copy.copy(x) # 浅拷⻉
>>> x2 is x # 复制成功
False
>>> x2.data is x.data # 但成员 data 依旧指向原列表,仅仅复制了引用
True

>>> x3 = copy.deepcopy(x) # 深拷⻉
>>> x3 is x # 复制成功
False
>>> x3.data is x.data # 成员 data 列表同样被复制
False
>>> x3.data.append(3) # 向复制的 data 列表追加数据, 不会影响原列表
>>> x3.data
[1, 2, 3]
>>> x.data
[1, 2]
```

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
- sys.path : sys.path[0]为空字符串, 表示启动python解析器时的当前目录.

    执行的py文件与其导入的lib重名会报`No module named xxx`
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
.py文件的字节码形式,只有在import语句执行时进行，当.py文件**第一次被导入**时，它会被编码为字节代码，并将字节码写入同名的.pyc文件中. 后来每次导入操作都会直接执行.pyc 文件（当.py文件的修改时间发生改变，这样会生成新的.pyc文件），在解释器使用-O选项时，将使用同名的.pyo文件，这个文件去掉了断言（assert）、断行号以及其他调试信息，体积更小，运行更快.（使用-OO选项，生成的.pyo文件会在`-O`的基础上再去除__doc__ string(文档信息)).

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

python 使用 lambda 来创建匿名函数, 格式:`lambda [arg1 [,arg2,.....argn]]:expression`.

### async/await
async和await是针对coroutine的新语法, 从python 3.5开始.

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
> issubclass(int, object) : 所有类型都有一个共同祖先类型 object,它为所有类型提供原始模版,以及系统所需的基本操作方式.

可在子类中定义一个与要重写的父类方法同名的方法, 这样, Python将不会考虑这个父类方法，而只关注你在子类中定义的相应方法.

属性和方法清单以及文件都越来越长时, 可能需要将类的一部分作为一个独立的类提取出来, 这样可以将大型类拆分成多个协同工作的小类, 比如上面的`Battery类`.

从一个模块中导入多个类:`from car import Car, ElectricCar `.
导入整个模块:` import car`
导入模块中的所有类:`from module_name import *`, 不推荐, 原因和导入module相同. 此时以`_`开头的全局变量不会被导入, 但可用明确指定名称的方式来导入, 比如`from module_name import _x`.

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

单双下划线:
1. _xxx     不能用于’from module import *’ 以单下划线开头的表示的是protected类型的变量. 即保护类型只能允许其本身与子类进行访问.
2. __xxx    双下划线的表示的是私有类型的变量. 只能是允许这个类本身进行访问了, 连子类也不可以
3. __xxx__ 定义的是python的特殊方法, 像__init__之类的

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

### metaclass
metaclass允许创建类或者修改类.

当传入关键字参数metaclass时，魔术就生效了，它指示Python解释器在创建MyList时，要通过ListMetaclass.__new__()来创建，在此，开发者可以修改类的定义，比如加上新的方法，然后返回修改后的定义.

`__new__()`方法接收到的参数依次是：
- 当前准备创建的类的对象
- 类的名字
- 类继承的父类集合
- 类的方法集合

```python
# metaclass是类的模板，所以必须从`type`类型派生：
class ListMetaclass(type):
    def __new__(cls, name, bases, attrs):
        attrs['add'] = lambda self, value: self.append(value)
        return type.__new__(cls, name, bases, attrs)

class MyList(list, metaclass=ListMetaclass):
    pass
```


truenas scale example:
```python
class ServiceBase(type):
    """
    Metaclass of all services

    This metaclass instantiates a `_config` attribute in the service instance
    from options provided in a Config class, e.g.

    class MyService(Service):

        class Meta:
            namespace = 'foo'
            private = False

    Currently the following options are allowed:
      - datastore: name of the datastore mainly used in the service
      - datastore_extend: datastore `extend` option used in common `query` method
      - datastore_prefix: datastore `prefix` option used in helper methods
      - service: system service `name` option used by `SystemServiceService`
      - service_model: system service datastore model option used by `SystemServiceService` (`service` if used if not provided)
      - service_verb: verb to be used on update (default to `reload`)
      - namespace: namespace identifier of the service
      - namespace_alias: another namespace identifier of the service, mostly used to rename and
                         slowly deprecate old name.
      - private: whether or not the service is deemed private
      - verbose_name: human-friendly singular name for the service
      - thread_pool: thread pool to use for threaded methods
      - process_pool: process pool to run service methods
      - cli_namespace: replace namespace identifier for CLI
      - cli_private: if the service is not private, this flags whether or not the service is visible in the CLI
    """

    def __new__(cls, name, bases, attrs):
        print("---1---", cls)
        print("---2---", name)
        print("---3---", bases)
        print("---4---", attrs)
        
        super_new = super(ServiceBase, cls).__new__
        if name == 'Service' and bases == ():
            print("---Service---")
            return super_new(cls, name, bases, attrs)

        config = attrs.pop('Config', None)
        klass = super_new(cls, name, bases, attrs)

        if config:
            klass._config_specified = {k: v for k, v in config.__dict__.items() if not k.startswith('_')}
        else:
            klass._config_specified = {}

        print("---5:_config_specified---", klass._config_specified)

        klass._config = service_config(klass, klass._config_specified)

        print("---6:_config---", klass._config.__dict__)
        return klass


def service_config(klass, config):
    namespace = klass.__name__
    if namespace.endswith('Service'):
        namespace = namespace[:-7]
    namespace = namespace.lower()

    config_attrs = {
        'datastore': None,
        'datastore_prefix': '',
        'datastore_extend': None,
        'datastore_extend_context': None,
        'datastore_primary_key': 'id',
        'datastore_primary_key_type': 'integer',
        'event_register': True,
        'event_send': True,
        'service': None,
        'service_model': None,
        'service_verb': 'reload',
        'service_verb_sync': True,
        'namespace': namespace,
        'namespace_alias': None,
        'private': False,
        'thread_pool': None,
        'process_pool': None,
        'cli_namespace': None,
        'cli_private': False,
        'cli_description': None,
        'verbose_name': klass.__name__.replace('Service', ''),
    }
    config_attrs.update({
        k: v
        for k, v in list(config.items())
        if not k.startswith('_')
    })

    return type('Config', (), config_attrs)


class Service(object, metaclass=ServiceBase):
    """
    Generic service abstract class

    This is meant for services that do not follow any standard.
    """
    def __init__(self):
        print("Init Service")

class CoreService(Service):

    class Config:
        cli_namespace = 'system.core'

c = CoreService()
```

执行:
```bash
# python3.9 p.py 
---1--- <class '__main__.ServiceBase'>
---2--- Service
---3--- (<class 'object'>,)
---4--- {'__module__': '__main__', '__qualname__': 'Service', '__doc__': '\n    Generic service abstract class\n\n    This is meant for services that do not follow any standard.\n    ', '__init__': <function Service.__init__ at 0x7f0b3a03a790>}
---5:_config_specified--- {}
---6:_config--- {'datastore': None, 'datastore_prefix': '', 'datastore_extend': None, 'datastore_extend_context': None, 'datastore_primary_key': 'id', 'datastore_primary_key_type': 'integer', 'event_register': True, 'event_send': True, 'service': None, 'service_model': None, 'service_verb': 'reload', 'service_verb_sync': True, 'namespace': '', 'namespace_alias': None, 'private': False, 'thread_pool': None, 'process_pool': None, 'cli_namespace': None, 'cli_private': False, 'cli_description': None, 'verbose_name': '', '__module__': '__main__', '__dict__': <attribute '__dict__' of 'Config' objects>, '__weakref__': <attribute '__weakref__' of 'Config' objects>, '__doc__': None}
---1--- <class '__main__.ServiceBase'>
---2--- CoreService
---3--- (<class '__main__.Service'>,)
---4--- {'__module__': '__main__', '__qualname__': 'CoreService', 'Config': <class '__main__.CoreService.Config'>}
---5:_config_specified--- {'cli_namespace': 'system.core'}
---6:_config--- {'datastore': None, 'datastore_prefix': '', 'datastore_extend': None, 'datastore_extend_context': None, 'datastore_primary_key': 'id', 'datastore_primary_key_type': 'integer', 'event_register': True, 'event_send': True, 'service': None, 'service_model': None, 'service_verb': 'reload', 'service_verb_sync': True, 'namespace': 'core', 'namespace_alias': None, 'private': False, 'thread_pool': None, 'process_pool': None, 'cli_namespace': 'system.core', 'cli_private': False, 'cli_description': None, 'verbose_name': 'Core', '__module__': '__main__', '__dict__': <attribute '__dict__' of 'Config' objects>, '__weakref__': <attribute '__weakref__' of 'Config' objects>, '__doc__': None}
Init Service
```

### 单例模式
参考:
- [Python中的单例模式实现](http://whosemario.github.io/2016/01/22/pattern-singleton/)

利用Python的metaclass实现单例(**推荐**, 其他方法有缺陷).

## 异常
```python
filename = "/etc/default/corosync"
tempfile = "/etc/default/corosync.tmp"
with open(filename) as f, open(tempfile, "w") as working:
    for line in f:
        if "START" in line:
            working.write("START=no")
        else:
            working.write(line)
os.rename(tempfile, filename)

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
推导式comprehensions（又称解析式），是Python的一种独有特性, 格式: `[输出表达式 数据源迭代 过滤表达式(可选)]`

> 列表推导式更快一点是因为它针对 Python 解释器做了优化.


推导式是可以从一个数据序列构建另一个新的数据序列的结构体. 共有三种推导式，在Python2和3中都有支持：

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

## 调试
参考:
- [Python 代码调试技巧](https://www.ibm.com/developerworks/cn/linux/l-cn-pythondebugger/index.html)

## 性能
对于 Python 来说,充分利用多核性能的阻碍主要在于 Python 的 全局解释器锁 (GIL). GIL 确保 Python 进程一次只能执行一条指令,无论当前有多少个核心. 该问题可以通过一些方法来避免,比如标准库的 multiprocessing,或 numexpr、Cython 等技术,或分布式计算模型等.

可选的编译器总结
Cython Shed Skin Numba Pythran PyPy
成熟度 Y Y _ _ Y
广泛使用性 Y _ _ _ _
支持 Numpy Y _ Y Y _
没有间断的代码改动 _ Y Y Y Y
需要有 C 的知识 Y _ _ _ _
支持 OpenMP Y _ Y Y Y

> multiprocessing多线程遇到syscall阻塞时会切换线程使得在GIL存在的情况下, 也有较小的概率提升性能.

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
### 解包操作优先保障没有`*`前缀变量的赋值,所以右值元素不能少于此数量. 另外,星号只能有一个,否则无法界定收集边界. 星号不能单独出现,要么与其他名字一起,要么放入列表或元组内.
>>> a, *b, c = 1, 2
>>> a, b, c
1, [], 2
### list赋值
l1=[2,3,4]
l2=l1 # l1, l2引用了同一个对象
l1[0]=24
l2
[24, 3, 4]
l1==l2 # 判断两个被引用的对象是否有相同的值
True
l1 is l2 # 判断两个变量名是否引用了相同对象
True

l2=l1[:] # 浅拷贝整个列表
l1[0]=25 # l1[0]绑定了新对象25

l2 # l2[0]引用原有对象
[24, 3, 4]

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

# set相减是集合的求差集

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
i 整数
n 与g相同，但插入随区域而异的数字分隔符
o 将整数表示为八进制数
r 同s, 但使用repr(将对象转化为供解释器读取的形式), 而不是str
s 保持字符串的格式不变，这是默认用于字符串的说明符
x 将整数表示为十六进制数并使用小写字母
X 与x相同，但使用大写字母
u 无符号整数
% 将数表示为百分比值（乘以100，按说明符f设置格式，再在后面加上%）

```python
>>> "%(n)d %(x)s" % {"n":1, "x":"spam"}" # 基于dict的字符串格式化, 但推荐使用**format**
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
>>> a=r'xxx' #r表示raw字符串, 即不支持转义的字符串
>>> ord('s') # 将单个字符转为其ascii码
115
>>> chr(115) # 将ascii码转成字符
's'
```


```python
0 < timeout == count => 0 < timeout && timeout == count
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
motorcycles.sort() # 就地排序, 调用后motorcycles即变成了sorted列表, 没有返回值
 motorcycles.sort(key=len) # 法sort接受两个可选参数：key和reverse, key用于排序的函数, 用该函数的结果作为排序依据. `reverse=True`表示逆向排序.
motorcycles.reverse() # 反向排序
list(reversed(motorcycles))  # 用list将返回的对象(迭代器)转换为列表
sorted(motorcycles) # motorcycles不变, 返回排序后的内容. 该方法支持传递参数`reverse=True`进行反向排序. 
len(motorcycles) # 获取列表的长度
len("a我") # python2.7返回字节数, python3返回字符数
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

## setup.py
setuptools 是一个优秀的，可靠的 Pthon 包安装与分发工具.

Python 库打包的格式包括 Wheel 和 Egg. Egg 格式是由 setuptools 在 2004 年引入，而 Wheel 格式是由 PEP427 在 2012 年定义. 使用 Wheel 和 Egg 安装都不需要重新构建和编译，其在发布之前就应该完成测试和构建.

Egg 和 Wheel 本质上都是一个 zip 格式包，Egg 文件使用 .egg 扩展名，Wheel 使用 .whl 扩展名. **Wheel 的出现是为了替代 Egg，其现在被认为是 Python 的二进制包的标准格式**.

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
multiprocessing和multiprocessing.dummy(复制了 multiprocessing 的 API，不过是在 threading 模块之上包装了一层)是Python下两个常用的多进程和多线程模块.

> 其中对于多进程，multiprocessing.freeze_support()语句在windows系统上是必须的，这是因为windows的API不包含fork()等函数. 因此当子进程被创建后，会重新执行一次全局变量的初始化, 这可能就会覆盖掉主进程中已经被修改了的值.

> **主进程会等待直到multiprocessing的子进程退出才退出**.


multiprocessing的主要组件是:
- 进程

    一个当前进程的派生(forked)拷贝,创建了一个新的进程标识符,并且任务在操作系统中以一个独立的子进程运行. 可以启动并查询进程的状态并给它提供一个目标方法来运行.
- 池

    包装了进程或线程. 在一个方便的工作者线程池中共享一块工作并返回聚合的结果.
- 队列
    
    一个先进先出(FIFP)的队列允许多个生产者和消费者.
- 管理者

    一个单向或双向的在两个进程间的通信渠道
- ctypes

    允许在进程派生(forked)后,在父子进程间共享原生数据类型(例如,整型数、浮点数和字节数)
- 同步原语

    锁和信号量在进程间同步控制流

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

### threading
使用方法:
1. 使用Thread类
```python
# 导入Python标准库中的Thread模块 
from threading import Thread 
# 创建一个线程
# function_name: 需要线程去执行的方法名
# args: 线程执行方法接收的参数，该属性是一个元组，如果只有一个参数也需要在末尾加逗号
mthread = threading.Thread(target=function_name, args=(function_parameter1, function_parameterN)) 
# 启动刚刚创建的线程 
mthread .start() # 非阻塞方法
```

1. 重新写一个类，继承threading.Thread
```python
# 重写了父类threading.Thread的run()方法，其他方法（除了构造函数)都不应在子类中被重写. 使用线程的时候生成一个子线程类的对象，然后对象调用start()方法就可以运行线程
import threading, time
class MyThread(threading.Thread):
    def __init__(self):
        threading.Thread.__init__(self)
    def run(self):
        global n, lock
        time.sleep(1)
        if lock.acquire():
            print n , self.name
            n += 1
            lock.release()
if "__main__" == __name__:
    n = 1
    ThreadList = []
    lock = threading.Lock()
    for i in range(1, 200):
        t = MyThread()
        ThreadList.append(t)
    for t in ThreadList:
        t.start() # 会执行run()
    for t in ThreadList:
        t.join() # 线程等待，我们的主线程不会等待子线程执行完毕再结束自身. 可以使用Thread类的join()方法来子线程执行完毕以后，主线程再关闭
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

## demo
### http_auth.py
twisted 主要概念:
- Site : 负责创建HTTPChannel实例来解析HTTP请求，并开始对象查找过程. 它们包含根资源，即代表该网站上的URL资源.

    - requestFactory=Request ：指定了请求报文处理工厂
- Resource : 代表一个单一的URL段. IResource接口描述了资源对象必须实现的方法，以便参与对象发布过程.
- Resource trees : 将资源对象安排成一个资源树. 从根资源对象开始，资源对象树定义了所有有效的URL.
- .rpy script : 是python脚本，twisted.web静态文件服务器会执行它，就像CGI一样. 但是，与CGI不同的是，它们必须创建一个Resource对象，当URL被访问时，该对象将被渲染.
- Resource rendering : 当Twisted Web定位到一个叶子资源对象时，资源渲染就会发生. Resource可以返回一个html字符串，也可以写到请求对象.
- Session : 允许在多个请求中存储信息. 每个使用系统的浏览器都有一个唯一的Session实例.

```python
# from https://gist.github.com/mrchrisadams/169102
# Copyright (c) 2008 Twisted Matrix Laboratories.
# See LICENSE for details.

# 参考:
# - [Configuring and Using the Twisted Web Server](https://twistedmatrix.com/documents/current/web/howto/using-twistedweb.html)
# - [Configuring and Using the Twisted Web Server翻译](https://blog.csdn.net/xiarendeniao/article/details/9844117)
# - [Python的Twisted框架中使用Deferred对象来管理回调函数](https://www.jb51.net/article/85041.htm)
# curl  http://192.168.0.112:8080/ab -u joe:blow
import sys

from zope.interface import implementer

from twisted.python import log
from twisted.internet import reactor
from twisted.web import server, resource, guard
from twisted.cred.portal import IRealm, Portal
#from twisted.cred.checkers import FilePasswordDB
from twisted.cred.checkers import InMemoryUsernamePasswordDatabaseDontUse

class GuardedResource(resource.Resource): # 通过验证后执行: getChild() -> render_${Method}()/render(), render()用于兜底, 否则未匹配时会返回501
    """
    A resource which is protected by guard and requires authentication in order
    to access.
    """
    def getChild(self, path, request):
        return self

    def render(self, request):
        # is served on root
        return b"Authorized!"
    def render_GET(self, request):
        return b"Authorized! GET"

class SimpleRealm(object):
    """
    A realm which gives out L{GuardedResource} instances for authenticated
    users.
    """
    implementer(IRealm)
    # requestAvatar supplies the username, and checks against the corresponding password 
    def requestAvatar(self, avatarId, mind, *interfaces):
        if resource.IResource in interfaces:
            # somewhat confused here...
            return resource.IResource, GuardedResource(), lambda: None
        raise NotImplementedError()



# compare password, 
def cmp_pass(uname, password, storedpass):
   return crypt.crypt(password, storedpass[:2])

# checker  opens a file called htpasswd, and passing in the hash as defined by the method cmp_pass
# checkers使用requestAvatarId()校验用户授权. 这里'blow'是joe的密码.
checkers = [InMemoryUsernamePasswordDatabaseDontUse(joe='blow')]# [FilePasswordDB(path_to_htpasswd, hash=cmp_pass)]

# guard acts like middleware, forcing all incoming requests to 'yoursite.com' be checked the file defined in checkers
wrapper = guard.HTTPAuthSessionWrapper(Portal(SimpleRealm(), checkers), [guard.BasicCredentialFactory('auth')])

# serves the this as a resource on port 8080
reactor.listenTCP(8080, server.Site(resource=wrapper))
reactor.run()
```

# 调用链:
# 1. SimpleRealm() -> GuardedResource(), 开始等待请求
# 1. guard.HTTPAuthSessionWrapper.render() ->  guard.HTTPAuthSessionWrapper._authorizedResource()
#                                               -> guard.BasicCredentialFactory('auth').decode()
#                                               -> guard.HTTPAuthSessionWrapper._login()
#                                                   -> Portal().login()
#                                                    -> checkers.requestAvatarId
#                                                      -> credentials.checkPassword()
#                                                    -> SimpleRealm.requestAvatar()
#                                               -> util.DeferredResource(self._login(credentials)).render()
#                                                       util.DeferredResource()._cbChild()
#
# 1. GuardedResource.getChild()
# 1. /usr/local/lib/python3.8/dist-packages/twisted/web/resource.py#Resource.render # 按照Resource的具体Method进行处理
# 1. GuardedResource.render_${Method}()/render()

## FAQ
### PyCharm 设置和更换 Python 版本
通过`Settings -> Project:xxx -> Project Interpreter`为每个项目选择指定的默认编译器版本. 没有时可从系统环境导入(前提是系统已安装该版本的python), 步骤:
    1. 选择`Project Interpreter`列表框右侧`设置`图标中的添加按钮
    1. 选择tab `System Interpreter`, 选择指定版本即可
### pycharm 导入模块时提示`unresolved reference`
本质均是修改环境变量PYTHONPATH(需精确到package的上级路径).

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
双星号（**kwargs）: 将参数以字典的形式导入

```python2
import concurrent.futures import ThreadPoolExecutor
import functools import wraps

def record_logger_decorator():
    def _record_logger_decorator(f):
        @wraps(f)
        def decorated(*args, **kwargs):
            print args
            print kwargs
            return 200, "ok"
        return kwargs
    return _record_logger_decorator

@record_logger_decorator()
def spawn_volgroup(group_uuid, group_name, triger=None, user=None):
    return 200

data = []
data.append((("id","name"), ({"triger":None, "user":None}))) # 当p[0]只有一个参数时, 需使用`("id",)`, 即末尾的`,`不能省略, 否则程序会panic
with ThreadPoolExecutor(max_workers=len(data)) as executor:
    for future in executor.map(lambda p: spawn_volgroup(*p[0], **p[1]), data)
```

### ImportError: No module named license.LicenseManager
明明`license.LicenseManager.py`却提示找不到, 因为LicenseManager.py同目录的`__init__.py`被删除了.

### ImportError: No module named Cython.Build
`pip install Cython`

### 编译成`.so`的源py文件被修改并重启应用后代码未生效
应先删除`.so`, 否则应用还是用旧的`.so`代码来运行

### 被调函数输出的日志信息(by logging)中的函数名和行数都是其装饰器函数的信息
原先代码被编译成了so, 源码更新后, 同名的`.so`未删除导致.

解决方法: 删除同名的`.so`

> 日志与代码无法对应, 根源在于运行中的代码不是最新代码, 检查看看是否有旧进程, `so`等在跑, 或更新错环境了等等.

### TypeError: Class advice impossible in Python3.  Use the @implementer class decorator instead.
python3使用`from zope.interface import implementer`, 而python2.7使用`from zope.interface import implements`


### Ubuntu 16.04构建的应用(py已编译成so)在oracle linux 7.9上报`undefined symbol: PyFPE_jbuf`
ref:
- [undefined symbol: PyFPE_jbuf 问题分析并处理](https://blog.csdn.net/u011728480/article/details/104817741)

```bash
# readefl -a app_xxx.so |grep PyFPE_ # 应用的so, 有PyFPE_jbuf
# readefl -a libpython3.5m.so |grep PyFPE_ # Ubuntu 16.04有PyFPE_jbuf, 而oracle linux 7.9没有
```

参考ref获悉: `PyFPE_jbuf`来自构建python3时的`--with-fpectl`(处理浮点数异常)选项, 使用后会启用pyfpe.h的WANT_SIGFPE_HANDLER, 从而拥有了PyFPE_jbuf.

### [Python 中如何将字节 bytes 转换为整数 int](https://www.delftstack.com/zh/howto/python/how-to-convert-bytes-to-integers/)
参考:
- [Python中struct.pack()和struct.unpack()用法](https://cloud.tencent.com/developer/article/1406350)

python2.7:
```python
import struct

# native byteorder
buffer = struct.pack("ihb", 1, 2, 3) # "ihb"表示有三个成员, 分别指明它们的编码方式
print(repr(buffer))
'''
b'\x01\x00\x00\x00\x02\x00\x03'
'''
print(struct.unpack("ihb", buffer))
'''
(1, 2, 3)
'''
```

python3的: `int.from_bytes() `

### subprocess.Popen执行`echo -e "123456\n123456" | pdbedit -a -t -u chen`报错, 未能从stdin读入密码
subprocess.Popen默认使用sh(dash)执行命令, 而dash不支持`echo -e`.

解决方法: 换用bash, 即`Popen(cmd, stdout=PIPE, stderr=PIPE, shell=True, executable='/bin/bash')`

### `'NoneType' object has no attribute '__getitem__'`
读取了None值的属性

### subprocess.Popen(cmd)卡住
使用`Popen(cmd, stdout=PIPE, stderr=PIPE)`

### [`if not someobj`和`if someobj == None`的过程](https://stackoverflow.com/questions/100732/why-is-if-not-someobj-better-than-if-someobj-none-in-python)
`if not someobj`:
1. If the object has a __nonzero__ special method (as do numeric built-ins, int and float), it calls this method. It must either return a bool value which is then directly used, or an int value that is considered False if equal to zero.
1. Otherwise, if the object has a __len__ special method (as do container built-ins, list, dict, set, tuple, ...), it calls this method, considering a container False if it is empty (length is zero).
1. Otherwise, the object is considered True unless it is None in which case, it is considered False.

`if someobj == None`:
1. If the object has a __eq__ method, it is called, and the return value is then converted to a boolvalue and used to determine the outcome of the if.
1. Otherwise, if the object has a __cmp__ method, it is called. This function must return an int indicating the order of the two object (-1 if self < other, 0 if self == other, +1 if self > other).
1. Otherwise, the object are compared for identity (ie. they are reference to the same object, as can be tested by the is operator).

### 将log同时输出到console和日志文件
```python
import logging

logger = logging.getLogger()
logger.setLevel('DEBUG') # 从哪个level开始输出
BASIC_FORMAT = "%(asctime)s:%(levelname)s:%(message)s"
DATE_FORMAT = '%Y-%m-%d %H:%M:%S'
formatter = logging.Formatter(BASIC_FORMAT, DATE_FORMAT)
chlr = logging.StreamHandler() # 输出到控制台的handler
chlr.setFormatter(formatter)
chlr.setLevel('INFO')  # 也可以不设置，不设置就默认用logger的level
fhlr = logging.FileHandler('example.log') # 输出到文件的handler
fhlr.setFormatter(formatter)
logger.addHandler(chlr)
logger.addHandler(fhlr)
logger.info('this is info')
logger.debug('this is debug')
```

log color可用coloredlogs.

### [python3: Python.h: No such file or directory](https://stackoverflow.com/questions/21530577/fatal-error-python-h-no-such-file-or-directory)

### pep8
#### E127：continuation line over-indented for visual indent
在括号内的参数很多的时候, 为了满足每一行的字符不超过79个字符, 需要将参数换行编写, 这个时候换行的参数应该与上一行的括号对齐.
或者将所有参数换行编写, 此时第一行不能有参数, 即第一行的最后一个字符一定要是(, 换行后需要有一个缩进. 类似的规则也用在[], {}上.

```python3
# Aligned with opening delimiter.
foo = long_function_name(var_one, var_two,
                         var_three, var_four)

# Hanging indents should add a level.
foo = long_function_name(
    var_one, var_two,
    var_three, var_four)
```

#### E128: continuation line under-indented for visual indent
```python3
urlpatterns = patterns('',
    url(r'^$', listing, name='investment-listing'),
)
```

改进: 将行缩进到左括号, 或者不在起始行上放置任何参数，然后缩进到统一级别.
```python3
urlpatterns = patterns('',
                       url(r'^$', listing, name='investment-listing'),
)

urlpatterns = patterns(
    '',
    url(r'^$', listing, name='investment-listing'),
)

urlpatterns = patterns(
    '', url(r'^$', listing, name='investment-listing'))
```

#### W291 trailing whitespace
行尾有多余的空格

#### AttributeError: '_ssl._SSLSocket' object has no attribute '_sslobj'
使用了[bareos example](https://pypi.org/project/python-bareos/)中的`Use JSON objects of the API mode 2`, 并用python3.8.5执行时报错.

```bash
wget https://github.com/drbild/sslpsk/archive/master.zip
tar -xvf master
cd sslpsk-master

python setup.py build
python setup.py install
```

#### vscode python自定义包无法跳转报: `Import "xxx" could not be resolved Pylance(reportMissingImports)`
在settings.json文件中添加:
```
"python.analysis.extraPaths": [
        "/home/chen/test/truenas/src/middlewared" # 路径必须正确
    ]
```

或`export PYTHONPATH="$PYTHONPATH:xxx"`

#### `pip3 install -r uvicorn`时报与urllib3相关的`temporary failure in name resolution`
无法访问到pip源即dns故障, 检查是否配置了`~/.pip/pip.conf`

### 安装setuptools-50.2.0.tar.gz`(python3 setup.py install )`报"RuntimeError: cannot build setuptools without metadata. Run `bootstrap.py`"
使用`pip3 install setuptools==50.2.0`

### `No module named 'gi'`
`yum install python3-gobject-base`

### [flask请求上下文流程](https://zhuanlan.zhihu.com/p/353187030)

### uwsgi 2.0.20(use http) +golang 1.17.2 + `guonaihong/gout v0.2.4`请求丢失
gout发给uwsgi的请求可能丢失, 表现为go报错:`EOF`或`read: connection reset by peer`.

解决方法: uwsgi前加nginx, uwsgi与nginx用socket通信.
nginx:
```conf
{
    listen 8359;
    server_name _;
    charset utf-8;

    location / {
        include /etc/nginx/uwsgi_params;
        uwsgi_pass 127.0.0.1:8361;
    }
}
```

uwsgi.ini:
```conf
[uwsgi]
socket=127.0.0.1:8361
...
```

### date format
- [Python datetime 格式化字符串：strftime()](https://blog.csdn.net/shomy_liu/article/details/44141483)

### 打印对象
ref:
- [Python打印对象的全部属性](https://blog.51cto.com/steed/2046408)

`print('\n'.join(('%s:%s' % item for item in dict1.items())))`  # 每行一对key和value，中间是分号

### `pyrasite-shell <pid>`卡住
ref:
- [python内存泄露诊断过程记录pyrasite](https://www.cnblogs.com/shengulong/p/8044132.html)

解决方法:
1. 安装依赖

    1. `pip2 install pyrasite meliae`
    1. `dnf install gdb python2-debuginfo`或`apt install python-dbg gdb`
1. 开启进程调试

    `echo 0|sudo tee /proc/sys/kernel/yama/ptrace_scope` # 关闭kernel对进程调试的检查, 否则会导致程序卡住. kylin v10 arm64没有该文件

### 安装高版本python
```bash
$ sudo apt install software-properties-common
$ sudo add-apt-repository ppa:deadsnakes/ppa
$ sudo apt update
$ sudo apt install python3.10
```

### python3-aiohttp报`type object '_asyncio.Task' has no attribute 'all_tasks'`
asyncio.Task.all_tasks() is fully moved to asyncio.all_tasks() starting with 3.9. Also applies to current_task. 升级aiohttp即可.


### ModuleNotFoundError: No module named 'xxx'
已检查PYTHONPATH下存在package xxx.

经检查, 当前python脚本的命名是'xxx.py', 与package xxx重名, 导致脚本中的`import xxx`将当前脚本当做package xxx来import, 从而导致出错.

### ModuleNotFoundError: No module named 'PyQt5'
`apt install python3-pyqt5`

### `pip install -r prod.txt`报`Cannot uninstall 'pyparsing'. It is a distutils installed project and thus we cannot accurately determine which files belong to it which would lead to only a partial uninstall.`
os已有rpm安装的`pyparsing 2.0.3`

解决方法:
`pip install --ignore-installed pyparsing==2.1.0`

### `pip install pyvmomi`报`'extras_require' must be a dictionary whose values are strings or lists of strings containing valid project/version requirement specifiers`
编译pyvmomi的setuptools太旧, 解决方法: `pip install --upgrade setuptools`

### `pip download -r req.txt`报`rsa requires Python '>=3.5, < 4' but the running Python is 2.7.16`

之前已通过`pip install rsa=3.4.2`, 当req.txt没有rsa=3.4.2时, oauth2client==4.0.0去拉取了rsa最新的4.0版本, 该版本需要python3.

解决方法: 将rsa=3.4.2加入到req.txt即可.

### `pip donwnload scipy==1.2.2 statsmodels==0.5.0`报`Command "python setup.py egg_info"`

需要先下载scipy并安装后才能安装statsmodels, 可能是scipy需要编译的缘故???.

### pip install报`Could not find a version that satisfies the requirement wheel (from version:)`

源或网络问题, 切换到国内源或多尝试几次. 网上也有提示说是先升级pip到最新, 可以获取更精准的日志.

也可能是setuptool的版本不够导致.

`pip install psutil`正常. 但`pip download`+`pip install`/`requirements.txt`+`pip download`+`pip install`一起安装`wrapt和psutils`就报错, 但先安装wrapt再安装psutils又是好的, 已检查psutils没有依赖, py真是一门神奇的语言.

pip安装的某个软件版本不对(我这里是过高)也会报该错, 但pip没有提示哪个组件的版本有问题, 怀疑是使用了多个requirements.txt, pip计算出来的版本有差异导致该问题.

### apscheduler报"Unable to determine the name of the local timezone -- you must explicitly specify the name of the local timezone"
装了个tzlocal, 用get_localzone 这个方法来打印自己机器的timezone, 发现结果就很奇怪, 系统显示的是正常的`Asia/Beijing`, 但是python获取的却是`local?`

```bash
rm -rf /etc/localtime
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 如果是systemd系统就更简单了,单单执行下面命令就OK了
timedatectl set-timezone Asia/Shanghai
```

### paramiko使用ed25519
[paramiko>=2.2.0](https://stackoverflow.com/questions/60660919/paramiko-ssh-client-is-unable-to-unpack-ed25519-key)

### paramiko只使用密码
```python
client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect(hostname="host.sftp.com", username="user", password="xxx", look_for_keys=False) # 没有look_for_keys=False时, 即使指定了password, 也是优先使用pkey
```

### [paramiko - connect with private key - not a valid OPENSSH private/public key file](https://stackoverflow.com/questions/45829838/paramiko-connect-with-private-key-not-a-valid-openssh-private-public-key-fil)

To convert "BEGIN OPENSSH PRIVATE KEY" to "BEGIN RSA PRIVATE KEY":
```
ssh-keygen -p -m PEM -f ~/.ssh/id_rsa
```

### paramiko log
paramiko.util.log_to_file("<log_file_path>", level = "WARN") # level in [DEBUG, WARN]

### 使用linux下gdb来调试python程序
**gdb调试Python没有pdb那么方便, 主要是没法直接给python代码打断点**

前提条件:
1. 确保gdb版本>=7.0
2. 安装python-debuginfo包

    apt install python<3.8>-dbg/python3-dbg

    > centos需先执行`yum install yum-utils && debuginfo-install glibc`再安装gdb和python-debuginfo

    其他方式安装python-debuginfo:
    ```bash
    # debuginfo-install python # for python2.7
    # debuginfo-install python3 libgcc # for python3
    ```

假设python进程是1000:
1. 执行`gdb -p 1000`(需手动载入libpython)/`gdb python 1000`(已安装libpython, 自动载入)
2. 载入[libpython : 找对应python版本的](https://github.com/python/cpython/blob/main/Tools/gdb/libpython.py)

    如果gdb是redhat或fedora等厂商修改过的，会有--python选项，使用此选项即可指定gdb启动时载入的Python扩展脚本: `gdb --python /path/to/libpython .py -p 1000`

    安装的是GNU的gdb，就需要打开gdb后手动载入libpython.py脚本:
    ```gdb
    (gdb) python
    >import sys
    >sys.path.insert(0, '/path/to/libpython.py' )
    >import libpython
    >end
    ```

    网上也有其他载入libpython方式:
    ```bash
    $ gdb python <PID>
    (gdb) source /usr/lib/debug/usr/lib64/libpython2.7.so.1.0.debug-gdb.py # from python2.7
    (gdb) source /usr/lib/debug/usr/lib64/libpython3.6dm.so.1.0-3.6.8-18.el7.x86_64.debug-gdb.py # python3
    ```
3. 开始调试

    使用py-bt命令打印当前线程的Python traceback了, 其他命令:
    - py-print : 打印变量
    - py-locals : 打印所有本地变量
    - py-bt  : 当前Py调用栈
    - py-bt-full: 输出Python调用栈
    - py-list  : 当前py代码位置
    - py-up : 查看上层调用方的信息
    - py-down : 返回之前的栈

    命令详细可打开libpython.py查看

    常用其他gdb命令:
    ```
    bt    # 当前C调用栈
    info thread   # 线程信息
    info threads  # 所有线程信息
    thread <id>   # 切换到某个线程
    thread apply all py-list  # 查看所有线程的py代码位置
    ctrl-c  # 中断
    ```

    为了不影响运行中的进程，可以通过生成 core file 的方式来保存进程的当前信息:
    ```bash
    (gdb) generate-core-file
    (gdb) quit
    $ gdb python core.6489
    ```

### pdb
```
pdb 命令行：
    1）进入命令行 Debug 模式，python -m pdb xxx.py
    2）h：（help）帮助
    3）w：（where）打印当前执行堆栈
    4）d：（down）执行跳转到在当前堆栈的深一层（个人没觉得有什么用处）
    5）u：（up）执行跳转到当前堆栈的上一层
    6）b：（break）添加断点
                 b: 列出当前所有断点，和断点执行到统计次数
                 b line_no：当前脚本的line_no行添加断点
                 b filename:line_no：脚本 filename 的 line_no 行添加断点
                 b function：在函数 function 的第一条可执行语句处添加断点
    7）tbreak：（temporary break）临时断点
                 在第一次执行到这个断点之后，就自动删除这个断点，用法和b一样
    8）cl：（clear）清除断点
                cl 清除所有断点
                cl bpnumber1 bpnumber2... 清除断点号为bpnumber1,bpnumber2...的断点
                cl lineno 清除当前脚本lineno行的断点
                cl filename:line_no 清除脚本filename的line_no行的断点
    9）disable：停用断点，参数为 bpnumber，和 cl 的区别是，断点依然存在，只是不启用
    10）enable：激活断点，参数为 bpnumber
    11）s：（step）执行下一条命令
                如果本句是函数调用，则s会执行到函数的第一句
    12）n：（next）执行下一条语句
                如果本句是函数调用，则执行函数，接着执行当前执行语句的下一条。
    13）r：（return）执行当前运行函数到结束
    14）c：（continue）继续执行，直到遇到下一条断点
    15）l：（list）列出源码
                 l: 列出当前执行语句周围11条代码
                 l first: 列出first行周围11条代码
                 l first, second: 列出first--second范围的代码，如果second<first，second将被解析为行数
    16）a：（args）列出当前执行函数的参数
    17）p expression：（print）输出expression的值
    18）pp expression：好看一点的p expression
    19）run：重新启动debug，相当于restart
    20）q：（quit）退出debug
    21）j lineno：（jump）设置下条执行的语句函数
                只能在堆栈的最底层跳转，向后重新执行，向前可直接执行到行号
    22）unt：（until）执行到下一行（跳出循环），或者当前堆栈结束
    23）condition bpnumber conditon，给断点设置条件，当参数condition返回True的时候bpnumber断点有效，否则bpnumber断点无效
```

注意：
1. 直接输入 Enter，会执行上一条命令
1. 输入 PDB 不认识的命令，PDB 会把他当做 Python 语句在当前环境下执行

## 调试
### [madbg](https://stackoverflow.com/questions/25308847/attaching-a-process-with-pdb)
```bash
madbg attach <pid>
ipdb> help # 查看帮助
```

## celery
### redis keys
- celery：表示当前正在队列中的 task，等待被 worker 所接收. celery 消失说明任务已经被启动的 worker 接收了即当前没有等待被接收的任务
- `_kombu.binding.celery`: celery 使用 kombu 维护消息队列
- `_kombu.binding.celery.pidbox`: kombu 维护
- `_kombu.binding.celeryev`: kombu 维护, 用来记下当前连接的 worker
- unacked：被 worker 接收了但是还没开始执行的 task 列表
- unacked_index：用户标记上面 unacked 的任务的 id，理论上应该与 unacked 一一对应的