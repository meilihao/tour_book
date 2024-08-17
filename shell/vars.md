# 变量

变量名必须以字母或下划线字符开头,其余的字符可以是字母、数字(0~9)或下划线字符,同时变量名是大小写敏感的.

> 约定: 环境变量和从全局配置文件里读取的变量使用大写,其余均用小写,并使用下划线分隔变量名的各部分.

变量的赋值方式为：`变量名称=值`,当所赋的值中间含有空格时，要加上引号.
赋值号两边**没有任何空格**,当想取shell变量的值时，需要在变量名前加上$字符.

一般用花括号括起变量名,以便于理解和阅读,但在双引号字符串中扩展变量时不可省略,如下:

    $ test=ll
    $ echo "he${test}o"
    hello
    $ echo "he$testo" #被提前截断
    he

shell以相似的方式处理单/双引号括起的字符串,区别在于双引号里的变量可以进行替换和变量扩展,单引号里的内容为原样输出.

    $ test=world
    $ echo "hello ${test}"
    hello world
    $ echo 'hello ${test}'
    hello ${test}

反引号里的内容会被当做shell命令执行,输出结果替代其内容的位置.

    echo "kernel is `uname -s`."

`export varname`可将一个shell变量提升为环境变量.
`${#var}`,可获得变量的长度(字符数).

## 变量作用域

在脚本里的变量是全局变量,但在函数中可用local声明局部变量(会屏蔽同名的局部变量).

## 特殊变量

- `$0`：当前脚本的文件名,比如`./hello.sh`
- `$num`：num为从1开始的数字，$1是第一个参数，$2是第二个参数，${10}是第十个参数
- `$#` : 传给脚本的参数个数(不包括脚本自身路径)
- `$@` : 将每个参数作为单元返回一个参数列表
- `$*` : 保存传给脚本的所有参数,将所有参数作为一个整体返回（字符串）, 不实用, **推荐使用`$@`**
- `$?` : 上次执行命令的退出状态,0表示正常退出,非0为异常退出.
- `$$` :  当前Shell进程ID. 对于 Shell 脚本，就是这些脚本所在的进程ID.

# 变量替换
ref:
- [shell(bash)替换字符串大全](https://blog.csdn.net/coraline1991/article/details/120235471)

## ${var/Pattern/Replacement}

在变量var中第一个匹配Pattern的字符串会被Replacement代替. 如果省略了Replacement,则第一个匹配Pattern的字符串会被删除.

## ${var//Pattern/Replacement}
全局替换replacement.所有在变量var中被Pattern匹配到的都由Replacement代替.同上类似,如果省略了Replacement,则所有匹配Pattern的字符串会被删除.


ps:

参考:http://www.360doc.com/content/12/0419/14/496343_204908029.shtml
http://www.cnblogs.com/barrychiao/archive/2012/10/22/2733210.html

# 数组

bash支持`普通数组`和`关联数组`.普通数组只能用整数做索引,但关联数组可用**字符串作索引,类似其他语言的map**.

## 普通数组

在Shell中，用括号来表示数组(下标从`0`开始)，数组元素用“空格”符号分隔, 如果元素中有空格需要用`'`包裹.定义数组的一般形式为：
```shell
  array_name=(value1 ... valuen)
```
```shell
# 可以不使用连续的下标，而且下标的范围没有限制
# 方式1
array_name=(value0 value1 value2 value3)
# 方式2
array_name=(
value0
value1
value2
value3
)
# 方式3
array_name[0]=value0
array_name[1]=value1
array_name[2]=value2
```

### 读取数组

读取数组元素值的一般格式是：
```shell
    ${array_name[index]}
```
```shell
valuen=${array_name[2]}
echo ${array_name[*]} # 输出数组列表
echo ${array_name[@]} # 和上面相同
echo ${#array_name[*]} # 数组长度
echo ${!array_name[@]} # 索引列表
```

## 关联数组(字符串做索引)

其需要使用单独的声明语句将一个变量声明为关联数组:

    # declare -A ass_array # 定义
    # ass_array=([index1]=value1 [index2]=value2) # 赋值
    # ass_array[index3]=value3 # 赋值
    # echo ${ass_array[index1]}

## alias
```
alias new_command='command sequence'
```

忽略别名:
```
# \command # 在不可信环境下执行特权命令时, 可避免潜在的安全问题即利用别名伪装的特权命令.
```

## 参数扩展
ref:
- [Shell 截取文件名和后缀](http://zuyunfei.com/2016/03/23/Shell-Truncate-File-Extension/)

参数形式    扩展后
x{y,z}  xy xz
${x}{y, z}  ${x}y ${x}z
${x}{y, $z}     ${x}y ${x}${z}
${param#pattern}    从param前面删除pattern的最小匹配
${param##pattern}   从param前面删除pattern的最大匹配
${param%pattern}    从param后面删除pattern的最小匹配
${param%%pattern}   从param后面删除pattern的最大匹配
${param/pattern/string}     从param中用string替换pattern的第一次匹配，string可为空
${param//pattern/string}    从param中用string替换pattern的所有匹配，string可为空
${param:3:2}    截取$param中索引3开始的2个字符
${param:3}  截取$param中索引3至末尾的字符
${@:3:2}    截取参数列表$@中第3个开始的2个参数
${param:-word}  若$param为空或未设置，则参数式返回word，$param不变
${param:+word}  若$param为非空，则参数式返回word，$param不变
${param:=word}  若$param为空或为设置，则参数式返回word，同时$param设置为word
${param:?message}   若$param为空或为设置，则输出错误信息message，若包含空白符，则需引号

## FAQ
### 获取当前工作目录的绝对路径
```bash
SHELL_FOLDER=$(dirname $(readlink -f "$0")) # 推荐
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
```