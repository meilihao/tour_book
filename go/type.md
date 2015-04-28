## 接口

参考:[golang反射规则使用详解](http://www.sharejs.com/codes/go/7309)

接口类型的变量存储了两个内容：***赋值给变量实际的值**和**这个值的类型描述**。更准确的说，值是底层实现了接口的实际数据内容，而类型描述了这个项目完整的类型。例如下面，
```go
var r io.Reader
tty, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil { return nil, err }
r = tty
```

用模式的形式来表达 r 包含了的是 (value, type) 对，如 (tty, *os.File)。

注意: 类型 *os.File 除了 Read 方法还实现了其他方法：尽管接口值仅仅提供了访问 Read 方法的可能（即通过r 只能访问Read方法），但是内部包含了这个值的完整的类型信息（反射的依据）。
这也就是为什么可以这样做：
```go
var w io.Writer
w = r.(io.Writer) //接口查询
```

在这个赋值中的断言是一个类型断言：它断言了 r 内部的条目同时也实现了 io.Writer，因此可以赋值它到 w。在赋值之后，w 将会包含 (tty, *os.File)，跟在 r 中保存的一致。接口的静态类型决定了哪个方法可以通过接口变量调用，即便内部实际的值可能有一个更大的方法集。

接下来，可以这样做：
```go
view sourceprint?
var empty interface{}
empty = w
```

而空接口值 e 也将包含同样的 (tty, *os.File)。这很方便：空接口可以保存任何值同时保留关于那个值的所有信息。

（这里无需类型断言，因为 w 是肯定满足空接口的。在这个例子中，将一个值从 Reader 变为 Writer，由于 Writer 的方法不是 Reader 的子集，所以就必须明确使用类型断言。）

**一个很重要的细节是接口内部的对总是 (value, 实际类型) 的格式，而不会有 (value, 接口类型) 的格式。接口不能保存接口值。**

## 接口查询

接口查询是否成功，要在运行期才能够确定。他不像接口的赋值，编译器只需要通过静态类型检查即可判断赋值是否可行。
```go
var file1  Writer = ...
if file5,ok := file1.(two.IStream);ok {
...
}
```
这个if语句检查file1接口指向的对象实例是否实现了two.IStream接口，如果实现了，则执行特定的代码。

在Go语言中，你可以询问它指向的对象是否是某个类型，比如，
```go
var file1 Writer = ...
if file6,ok := file1.(*File);ok {
...
}
```
这个if语句判断file1接口指向的对象实例是否是*File类型，如果是则执行特定的代码。
```go
slice := make([]int, 0)
slice = append(slice, 1, 2, 3)

var I interface{} = slice


if res, ok := I.([]int)；ok {
    fmt.Println(res) //[1 2 3]
}
这个if语句判断接口I所指向的对象是否是[]int类型，如果是的话输出切片中的元素。
```
```go
func Sort(array interface{}, traveser Traveser) error {

  if array == nil {
    return errors.New("nil pointer")
  }
  var length int //数组的长度
  switch array.(type) {
  case []int:
    length = len(array.([]int))
  case []string:
    length = len(array.([]string))
  case []float32:
    length = len(array.([]float32))

  default:
    return errors.New("error type")
  }

  if length == 0 {
    return errors.New("len is zero.")
  }

  traveser(array)

  return nil
}
```
通过使用.(type)方法可以利用switch来判断接口存储的类型。

小结: 查询接口所指向的对象是否为某个类型的这种用法可以认为是接口查询的一个特例.

## 类型转换（Conversions）

类型转换的语法：
```go
Conversion = Type "(" Expression [ "," ] ")"
```
当以运算符*或者<-开始，有必要加上括号避免模糊：

```go
*Point(p)        // same as *(Point(p))
(*Point)(p)      // p is converted to *Point
<-chan int(c)    // same as <-(chan int(c))
(<-chan int)(c)  // c is converted to <-chan int
func()(x)        // function signature func() x
(func())(x)      // x is converted to func()
(func() int)(x)  // x is converted to func() int
func() int(x)    // x is converted to func() int (unambiguous)
```

## Type switch与Type assertions

在Go语言中，我们可以使用type switch语句查询接口变量的真实数据类型，语法如下：

```go
//x必须是接口类型
switch x.(type) {
// cases
}
```

来看一个详细的示例：

```go
type Stringer interface {
    String() string
}

//语句switch中的value必须是接口类型，变量str的类型为转换后的类型
var value interface{} // Value provided by caller.
switch str := value.(type) {
case string:
    return str //type of str is string
case Stringer: //type of str is Stringer
    return str.String()
}
```

If the switch declares a variable in the expression, the variable will have the corresponding type in each clause. It's also idiomatic to reuse the name in such cases, in effect declaring a new variable with the same name but a different type in each case.
如果我们只关心一种类型该如何做？如果我们知道值为一个string，只是想将它抽取出来该如何做？只有一个case的类型switch是可以的，不过也可以用类型断言（type assertions）。

类型断言接受一个接口值，从中抽取出显式指定类型的值。其语法借鉴了类型switch子句，不过是使用了显式的类型，而不是type关键字，如下：

```go
//同样，x必须是接口类型
x.(T)
```

str := value.(string)
上面的转换有一个问题，如果该值不包含一个字符串，则程序会产生一个运行时错误。为了避免这个问题，可以使用“comma, ok”的习惯用法来安全地测试值是否为一个字符串：

```go
//如果类型断言失败，则str将依然存在，并且类型为字符串，不过其为零值
str, ok := value.(string)
if ok {
    fmt.Printf("string value is: %q\n", str)
} else {
    fmt.Printf("value is not a string\n")
}
```

我们可以使用类型断言来实现type switch的中例子：
```go
//这种做法没有多大实用价值
if str, ok := value.(string); ok {
    return str
} else if str, ok := value.(Stringer); ok {
    return str.String()
}
```