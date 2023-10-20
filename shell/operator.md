## 逻辑运算
### 文件系统
- `[ -b file ]` : file存在且是块设备文件
- `[ -c file ]` : file存在且是字符设备文件
- `[ -d file ]` : file存在且是目录
- `[ -e file ]` : file存在
- `[ -f file ]` : file存在且是普通文件
- `[ -h file ]` : 如果文件是软链接，则为真
- `[ -L file ]` : file存在且是符号连接
- `[ -s file ]` : file存在且大小不为零
- `[ -r file ]` : file存在且用户可读
- `[ -w file ]` : file存在且用户可写
- `[ -x file ]` ：file存在且用户可执行
- `[ -z "$APP_PATH" ]` : 判断字符串长度是否为0
- `[ file1 -nt file2 ]` : file1比file2新(file1存在且file2不存在为真)
- `[ file1 -ot file2 ]` : file1比file2旧(file2存在且file1不存在为真)
- `[ file1 -ef file2 ]` : file1和file2是否硬连接到同一个文件

测试:
```sh
$ [ -e /dev/cdrom ] && echo "Exist" 
Exist
```
### 数值
- `-eq` : 数值等于
- `-ne` : 数值不相等
- `-lt` : 数值小于
- `-le` : 数值小于等于
- `-gt` : 数值大于
- `-ge` : 数值大于等于
### 字符串(**使用双中括号包裹**, `[]`会报错)
- `!=` : 字符串不相等
- `==`,`=` : 字符串等于,推荐使用`==`,因为`=`后必须紧跟一个空格,容易忘记.
- `<` : 字符串小于
- `<=` : 字符串小于等于
- `>` : 字符串大于
- `>=` : 字符串大于等于
- `[[ -n $var ]]` ：$var不为空
- `[[ -z $var ]]` ：$var为空
- `[[ ! $strA =~ $strB ]]` : strA不包含strB
### 逻辑
- `-a` : 逻辑与，操作符两边均为真，结果为真，否则为假
- `-o` : 逻辑或，操作符两边一边为真，结果为真，否则为假
- `!` : 逻辑否，条件为假，结果为真
- `if [ "${STATUS}" != 200 ] && [ "${STRING}" != "${VALUE}" ]; then`

# `(())`

描述:

双括号`(())`结构语句是对shell中算数及赋值运算的扩展.`(())`里的内容只是执行，并不会返回值,需加前缀`$`才可取到值.

语法：

    （（表达式1,表达式2…））

特点：

1. 在双括号结构中，所有表达式可以像c语言一样，如：a++,b--等
1. 在双括号结构中，所有变量可以不加前缀`$`
1. 双括号可以进行逻辑运算，四则运算
1. 双括号结构 扩展了for，while,if条件测试运算
1. 支持多个表达式运算，各个表达式之间用","分开

# $[...]

整数扩展(integer expansion)表达式,在方括号里面执行整数运算,`$[]`会返回里面表达式的值.

# `()`

命令组（Command group）.由一组圆括号括起来的命令是命令组，命令组中的命令实在子shell（subshell）中执行。因为是在子shell内运行，因此在括号外面是没有办法获取括号内变量的值，但反过来，命令组内是可以获取到外面的值，这点有点像局部变量和全局变量的关系.
在实作中，如果碰到要cd到子目录操作，并在操作完成后要返回到当前目录的时候，可以优先考虑使用subshell来处理.

> 子shell是嵌在圆括号()内部的命令序列，子Shell内部定义的变量为局部变量.

```shell
cmd1 | ( cmd2; cmd3; cmd4 ) | cmd5
```
如果cmd2 是cd /，那么就会改变子Shell的工作目录，这种改变只是局限于子shell内部，cmd5则完全不知道工作目录发生的变化.

## FAQ
### 判断一个命令是否存在
ref:
- [是否存在命令](https://stackoverflow.com/questions/592620/how-can-i-check-if-a-program-exists-from-a-bash-script)

```bash
command -v xxx >/dev/null 2>&1

# fail on non-zero return value
if [ "$?" -ne 0 ]; then
    return 1
fi
```

```bash
if ! [ -x "$(command -v git)" ]; then
  echo 'Error: git is not installed.' >&2
  exit 1
fi
```

### 判断文件是否存在
```bash
# 方法1, 推荐:
if compgen -G "${PROJECT_DIR}/*.png" > /dev/null; then
    echo "pattern exists!"
fi

# 方法2:
# ls returns non-zero when the files do not exist
if ls /path/to/your/files* 1> /dev/null 2>&1; then
    echo "files do exist"
else
    echo "files do not exist"
fi

# 方法3:
if stat --printf='' /path/to/your/files* 2>/dev/null
then
    echo found
else
    echo not found
fi
```