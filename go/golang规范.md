# 序言

为了代码美观.

## 格式化规范
go默认已经有了gofmt工具，但是我们强烈建议使用**goimport**工具，这个在gofmt的基础上增加了自动删除和引入包.

go get golang.org/x/tools/cmd/goimports
不同的编辑器有不同的配置, sublime的配置教程：http://michaelwhatcott.com/gosublime-goimports/

LiteIDE默认已经支持了goimports，如果你的不支持请点击属性配置->golangfmt->勾选goimports,之后编辑器保存之前自动fmt代码.

## 行长约定

一行最长不超过80个字符，超过的请使用换行展示，尽量保持格式优雅.

### 长句子打印或者调用，使用参数进行格式化分行

我们在调用fmt.Sprint或者log.Sprint之类的函数时，有时候会遇到很长的句子，我们需要在参数调用处进行多行分割

## import 规范
引入了三种类型的包，标准库包，程序内部包，第三方包，建议采用如下方式进行组织你的包：
- 有顺序的引入包，不同的类型采用空格分离，第一种标准库，第二是第三方包,第三是项目包.
- 在项目中不要使用相对路径引入包

## 代码文件名
尽量用一个语义话的单词命名,实在不行时使用下划线分隔.

## 变量声明

变量名采用驼峰标准，不要使用_来命名函数内的变量，多个相关的变量申明应放在一起.

常量,首字母大写(不可导出的常量前面加`_`);变量,首字母小写;参数传递：驼峰式，小写字母开头.

变量命名基本上遵循相应的英文翻译。
在相对简单的环境（对象数量少、针对性强）中，可以将一些名称由完整单词简写为单个字母，例如：
user 可以简写为 u
userId 可以简写 uid
若变量类型为 bool 类型，则名称应以 Has, Is, Can 或 Allow 开头.

### 变量命名惯例

代表某个用户：u
代表某个用户 ID：uid
代表某个索引：idx
代表某个值：val

## 参数传递

对于少量数据，不要传递指针
对于大量数据的struct可以考虑使用指针
传入参数是map，slice，chan不要传递指针
因为map，slice，chan是引用类型，不需要传递指针的指针

## 错误处理

错误处理的原则就是不能丢弃任何有返回err的调用，不要采用_丢弃，必须全部处理。接收到错误，要么返回err，要么实在不行就panic，或者使用log记录下来

## struct

struct申明和初始化格式采用多行.

## recieved

到底是采用值类型还是指针类型主要参考如下原则：

func(w Win) Tally(playerPlayer)int    //w不会有任何改变
func(w *Win) Tally(playerPlayer)int    //w会改变数据
更多的请参考：https://code.google.com/p/go-wiki/wiki/CodeReviewComments#Receiver_Type

**带mutex的struct必须是指针receivers**

指针receivers统一采用单字母'p'而不是this，me或者self.

## 函数或方法

函数或方法的参数排列顺序遵循以下几点原则：
1. 参数的重要程度与逻辑顺序(依赖关系)
2. 简单类型优先于复杂类型
3. 尽可能将同种类型的参数放在相邻位置，则只需写一次类型

如果一个结构拥有对应操作函数，大体上按照 CRUD 的顺序放置结构定义之后.

若函数或方法为判断类型（返回值主要为 bool 类型），则名称应以 Has, Is, Can 或 Allow 等判断性动词开头.

## 接口

单个函数的接口名以"er"作为后缀，如Reader,Writer

## range
如果只需要第一项（key），就丢弃第二个：

for key := range m {
    if key.expired() {
        delete(m, key)
    }
}
如果只需要第二项，则把第一项置为下划线

sum := 0
for _, value := range array {
    sum += value
}

## 项目
以下为一般项目结构，根据不同的 Web 框架习惯，可使用括号内的文字替换；根据不同的项目类型和需求，可自由增删某些结构：

- view：用于存放模板文件
- public(static)：用于存放静态文件
    - css：用于存放 CSS 文件
    - fonts：用于存放字体文件
    - img：用于存放图像文件
    - js：用于存放 JavaScript 文件
    - lib :第三方引用
- router(controller)：用于存放路由文件
- model(datatype)：用于存放业务逻辑层文件(某些统一的struct和相关方法),统一存放可避免循环引用
- module：用于存放子模块文件
    - setting：用于存放应用配置存取文件
- cmd：用于存放命令行程序命令文件
- conf：用于存放配置文件
    - locale: 用于存放 i18n 本地化文件
- data：用于存放应用生成数据文件
- log：用于存放应用生成日志文件
命令行应用

当应用类型为命令行应用时，需要将命令相关文件存放于 /cmd 目录下，并**为每个命令创建一个单独的源文件**.