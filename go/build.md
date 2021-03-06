
# 构建

## 条件编译
Go语言没有用预处理器、宏定义或＃define声明来控制指定平台，相反，Go语言标准库提供了go/build工具，该工具支持Go语言的构建标签（Build Tag）机制来构建约束条件（Build Constraint）. Go语言编译时的约束条件，其也被称为条件编译.

Go语言的条件编译有两种定义方法，分别介绍如下:
- 构建标签 ：在源码里添加注释信息，比如//+build linux，该标签决定了源码文件只在Linux平台上才会被编译
- 文件后缀 ：改变Go语言代码文件的后缀，比如foo_linux.go，该后缀决定了源码文件只在Linux平台上才会被编译

### 编译标签(`build tag`)：

```
// +build { GOOS }, { GOOS }, { !GOOS }
```
```
// +build (linux AND 386) OR (darwin AND (NOT cgo))
```
```
// +build linux darwin
// +build 386
```

1. 以`// +build`开头,一个源文件里可以有多个编译标签，多个编译标签之间是逻辑"与"的关系
1. 支持 GOOS 与 GOARCH 并可以具有多个值，用`,`分割， 例如：`// +build linux, darwin, freebsd`
1. 支持 不等条件 `!` ， 例如：`\\ +build !windows`,即不在windows环境下时，均可编译此文件。
1. 支持 与或非 逻辑， AND OR NOT
1. 条件编译需要前后空一行，否则无法识别

### 文件后缀：

```
xxx_{ GOOS }.go xxx_{ GOOS }_{ GOARCH }.go
```

1. 支持 GOOS ，例如： curl_windows.go
1. 支持 GOARCH， 例如： curl_386.go
1. 支持 上述两种叠加，但不可调换顺序 xxx_{ GOOS }_{ GOARCH }.go ，例如： curl_windows_amd64.go

### 如何选择

这两者可以叠加使用，但注意不要出现冗余，如：curl_windows.go 里面写`// +build windows`则重复了.
如果编译的文件是一一对应关系的话，使用文件后缀更简单些，如对每个 GOOS 生成一个文件.
如果有复杂条件的话，可以使用标签编译方式。如：
```
curl_windows.go 对应 windows 平台。
curl_others.go 里面写 \\ +build !windows 对应 非windows 平台
```

## 交叉编译

Golang 1.5及以上修改`GOOS和GOARCH`即可.如：

```
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test
$ file test
```

[GOOS 与 GOARCH 支持的参数](https://golang.org/doc/install/source#environment)

# ast
Go语言的优势在于它是一个静态类型语言，语法很简单，与动态类型语言相比更简单一些, 且Go语言标准库支持代码解析功能. 代码解析流程可分为3步:

1. 通过标准库go/tokens提供的Lexer词法分析器对代码文本进行词法分析，最终得到Tokens

    Go语言标准库提供了go/tokens词法分析器（Lexical Analyzer，简称Lexer，也被称为扫描器）. 词法分析是将字符序列转换为Tokens（或称Token序列、单词序列）的过程. 其工作原理是对输入的代码文本进行词法分析，将一个个字符以从左到右的顺序读入，根据构词规则识别单词，最终得到Token（单词）. Token是语言中的最小单位，它可以是变量、函数、运算符或数字.

2. 通过标准库go/parser和go/ast将Tokens构建为抽象语法树（AST）

    通过Lexer词法分析器得到Token序列以后，它将被传递给Parser解析器. 解析器是编译器的一个阶段，它将Token序列转换为抽象语法树（AST，Abstract Syntax Tree）. 抽象语法树也被称为语法树（Syntax Tree），是编程语言源码的抽象语法结构的树状表现形式，树上的每个节点都表示源码中的一种结构. 抽象语法树是源码的结构化表示.

3. 通过标准库go/types下的Check方法进行抽象语法树类型检查，完成代码解析过程

    通过Parser解析器得到抽象语法树之后，需要对抽象语法树中定义和使用的类型进行Type-Checking检查. 对每一个抽象语法树节点进行遍历，在每个节点上对当前子树的类型进行验证，进而保证不会出现类型错误. 通过Go语言标准库go/types下的Check方法进行抽象语法树检查. 另外，抽象语法树一般有多种遍历方式，比如深度优先搜索（DFS）遍历和广度优先搜索（BFS）遍历等.