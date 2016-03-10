# js学习笔记

参考：
- [JavaScript教程](http://www.liaoxuefeng.com/wiki/001434446689867b27157e896e74d51a89c25cc8b43bdb3000)
- [ECMAScript 6 入门](http://es6.ruanyifeng.com/)

工具:
- [caniuse # 了解各个浏览器以及它们的不同版本对HTML5、CSS3特性的支持情况](http://caniuse.com/)
- [ECMAScript兼容性(ES6)](http://kangax.github.io/compat-table/es6/)

# 基础

## 数据类型

JavaScript语言的每一个值，都属于某一种数据类型。JavaScript的数据类型，共有五种(ES6又新增了第六种Symbol类型的值)

- 数值（number）：整数和小数（比如1和3.14）
- 字符串（string）：字符组成的文本（比如”Hello World”）
- 布尔值（boolean）：true（真）和false（假）两个特定值
- undefined：表示“未定义”或不存在，即此处目前没有任何值
- null：表示空缺，即此处应该有一个值，但目前为空

对象（object）：各种值组成的集合,不是数据类型
通常将数值、字符串、布尔值称为原始类型（primitive type）的值，即它们是最基本的数据类型，不能再细分了.


### null和undefined

JavaScript的设计者希望用null表示一个空的值，而undefined表示值未定义。事实证明，区分两者的意义不大。**大多数情况下，我们都应该用null。undefined仅仅在判断函数参数是否传递的情况下有用**。null的特殊之处在于，JavaScript把它包含在对象类型（object）之中.

null表示空值，即该处的值现在为空。典型用法是：
- 作为函数的参数，表示该函数的参数是一个没有任何内容的对象。
- 作为对象原型链的终点。

undefined表示不存在值，就是此处目前不存在任何值。典型用法是：
- 变量被声明了，但没有赋值时，就等于undefined。
- 调用函数时，应该提供的参数没有提供，该参数等于undefined。
- 对象没有赋值的属性，该属性的值为undefined。
- 函数没有返回值时，默认返回undefined。

### boolean

除了下面六个值被转为false，其他值都视为true:
```
undefined
null
false
0
NaN
""（空字符串）
```

### 数组

**JavaScript的数组可以包括任意数据类型;访问时,索引超出了范围，返回undefined**.

    // 创建数组[1, 2, 3]
    new Array(1, 2, 3);
    [1,2,3] //推荐(代码的可读性)

给Array的length赋一个新的值会导致Array大小的变化:截断数组或扩充数组(使用undefined填充).

如果通过索引赋值时，索引超过了范围，同样会引起Array大小的变化(即扩充了数组)．

大多数其他编程语言不允许直接改变数组的大小，越界访问索引会报错。然而，**JavaScript的Array却不会有任何错误。在编写代码时，不建议直接修改Array的大小，访问索引时要确保索引不会越界**。

#### indexOf
Array也可以通过indexOf()来检索一个指定元素的下标，未找到则返回-1．

#### slice
slice()就是对应String的substring()版本，它截取Array的部分元素，然后返回一个新的Array：

    var arr = ['A', 'B', 'C', 'D', 'E', 'F', 'G'];
    arr.slice(0, 3); // 从索引0开始，到索引3结束，但不包括索引3: ['A', 'B', 'C']
    arr.slice(3); // 从索引3开始到结束: ['D', 'E', 'F', 'G']

**如果不给slice()传递任何参数，它就会从头到尾截取所有元素。利用这一点，我们可以很容易地复制一个Array**．

#### push和pop
push()向Array的末尾添加若干元素并返回新数组的长度，pop()则把Array的最后一个元素删除掉并返回该元素;空数组继续pop不会报错，而是会返回undefined．

#### unshift和shift
如果要往Array的头部添加若干元素，使用unshift()方法会返回新数组的长度，shift()方法则把Array的第一个元素删掉并返回该元素；空数组继续shift不会报错，而是会返回undefined．

#### sort
sort()可以对当前Array进行排序(默认升序)，它会直接修改当前Array的元素位置.

#### reverse
reverse()把整个Array的元素给掉个个，也就是反转.

#### splice
splice()方法是修改Array的“万能方法”，它可以从指定的索引开始删除若干元素并返回删除的元素，然后再从该位置添加若干元素.

```js
var arr = ['Microsoft', 'Apple', 'Yahoo', 'AOL', 'Excite', 'Oracle'];
// 从索引2开始删除3个元素,然后再添加两个元素:
arr.splice(2, 3, 'Google', 'Facebook'); // 返回删除的元素 ['Yahoo', 'AOL', 'Excite']
arr; // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
// 只删除,不添加:
arr.splice(2, 2); // ['Google', 'Facebook']
arr; // ['Microsoft', 'Apple', 'Oracle']
// 只添加,不删除:
arr.splice(2, 0, 'Google', 'Facebook'); // 返回[],因为没有删除任何元素
arr; // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
```
#### concat
concat()方法把当前的Array和另一个Array连接起来，并返回一个新的Array.实际上，concat()方法可以接收任意个元素和Array，并且自动把Array拆开，然后全部添加到新的Array里.

var arr = ['A', 'B', 'C'];
var added = arr.concat(4,[1, 2, 3]);
added; // ['A', 'B', 'C', 4, 1, 2, 3]
arr; // ['A', 'B', 'C']

请注意，concat()方法并没有修改当前Array，而是返回了一个新的Array。
#### join

join()方法把当前Array的每个元素都用指定的字符串连接起来，然后返回连接后的字符串：

    var arr = ['A', 'B', 'C', 1, 2, 3];
    arr.join('-'); // 'A-B-C-1-2-3'
如果Array的元素不是字符串，将自动转换为字符串后再连接.

### 对象

JavaScript的对象是一组由键-值组成的无序集合.JavaScript对象的键都是字符串类型，值可以是任意数据类型．

JavaScript中有两种方式来访问对象的属性，`点操作符`和`中括号操作符`．两种语法是等价的，但是中括号操作符在下面两种情况下依然有效：
- 动态设置属性　//obj['string'+variable]
- 属性名不是一个有效的变量名．//即不可省略属性的引号时．

>当属性名满足下面条件之一时，不能省去引号：
- 当属性名为JavaScript的保留字时
- 当属性名含有空格或特殊字符时（除了字母，数字和下划线外的字符）
- 属性名以数字开头

**当对象属性是固定的，且是有效变量名时推荐使用点操作符．同时我们在编写JavaScript代码的时候，属性名尽量使用标准的变量名，这样就可以直接通过object.prop的形式访问一个属性了**．

由于JavaScript的对象是动态类型，使用`delete`可删除对象的属性，且删除一个不存在的school属性也不会报错．

使用in操作符可检测对象是否拥有某一属性．
```js
var xiaoming = {
    name: '小明',
    score: null
};

'name' in xiaoming; // true
'grade' in xiaoming; // false
```
**不过要小心，如果in判断一个属性存在，这个属性不一定是xiaoming的，它可能是xiaoming继承得到的**：

    'toString' in xiaoming; // true
因为toString定义在object对象中，而所有对象最终都会在原型链上指向object，所以xiaoming也拥有toString属性。

要判断一个属性是否是xiaoming自身拥有的，而不是继承得到的，可以用hasOwnProperty()方法：

```js
var xiaoming = {
    name: '小明'
};
xiaoming.hasOwnProperty('name'); // true
xiaoming.hasOwnProperty('toString'); // false
```

### 变量

变量在JavaScript中就是用一个变量名表示，变量名是大小写英文、数字、$和_的组合，且不能用数字开头．如果一个变量没有通过var申明就被使用，那么该变量就自动被申明为全局变量．

**使用var申明的变量则不是全局变量，它的范围被限制在该变量被申明的函数体内，同名变量在不同的函数体内互不冲突**．

为了修补JavaScript这一严重设计缺陷，ECMA在后续规范中推出了strict模式，在strict模式下运行的JavaScript代码，强制通过var申明变量，未使用var申明变量就使用的，将导致运行错误。

启用strict模式的方法是在JavaScript代码的第一行写上：

    'use strict';
这是一个字符串，不支持strict模式的浏览器会把它当做一个字符串语句执行，支持strict模式的浏览器将开启strict模式运行JavaScript。

**不用var申明的变量会被视为全局变量，为了避免这一缺陷，推荐所有的JavaScript代码都应该使用strict模式．**

### 字符串

JavaScript的字符串就是用''或""括起来的字符表示．如果'本身也是一个字符，那就可以用""括起来；如果字符串内部既包含'又包含"则需用转义字符\来标识．

ASCII字符可以以\x##形式的十六进制表示:

    '\x41'; // 完全等同于 'A'
可以用\u####表示一个Unicode字符:

    '\u4e2d\u6587'; // 完全等同于 '中文'

由于多行字符串用\n写起来比较费事，所以最新的ES6标准新增了一种多行字符串的表示方法，用两个反引号包裹内容来表示．

**字符串是不可变的，如果对字符串的某个索引赋值，不会有任何错误，但是，也没有任何效果**．

常用操作：

```js
var s = 'Hello, world!中';
s.length;           // 获取字符串长度，14，中文算一个字符．
s.toUpperCase();    // 字符转大写
s.toLowerCase();　  // 字符转小写
s.indexOf('world'); // 检索指定字符串出现的起始位置，返回7，没找到时返回-1
s.substring(0, 5);  // 返回指定索引区间的子串，索引0开始到5（不包括5），返回'Hello'
s.substring(7);     // 从索引7开始到结束，返回'world!中'
//----其他方法
Method    描述
charAt()    返回指定索引位置的字符
charCodeAt()    返回指定索引位置字符的 Unicode 值
codePointAt 会正确返回四字节的UTF-16字符的码点。对于那些两个字节储存的常规字符，它的返回结果与charCodeAt方法相同。
concat()    连接两个或多个字符串，返回连接后的字符串
fromCharCode()    将字符转换为 Unicode 值(\u0000——\uFFFF)
fromCodePoint() 可以识别0xFFFF的字符，弥补了String.fromCharCode方法的不足.
indexOf()    返回字符串中检索指定字符第一次出现的位置
includes() 返回布尔值，表示是否找到了参数字符串。
startsWith() 返回布尔值，表示参数字符串是否在源字符串的头部。
endsWith() 返回布尔值，表示参数字符串是否在源字符串的尾部。
lastIndexOf()    返回字符串中检索指定字符最后一次出现的位置
localeCompare()    用本地特定的顺序来比较两个字符串
match()    找到一个或多个正则表达式的匹配
replace()    替换与正则表达式匹配的子串
repeat() 返回一个新字符串，表示将原字符串重复n次
search()    检索与正则表达式相匹配的值
slice()    提取字符串的片断，并在新的字符串中返回被提取的部分
split()    把字符串分割为子字符串数组
substr()    从起始索引号提取字符串中指定数目的字符
substring()    提取字符串中两个指定的索引号之间的字符
toString()    返回字符串对象值
toUpperCase()    把字符串转换为大写
trim()    移除字符串首尾空白
valueOf()    返回某个字符串对象的原始值
```
更多方法参考[JavaScript标准库之String](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/String)

### FAQ

#### map和set

JavaScript的默认对象表示方式`{}`可以视为其他语言中的`Map`或`Dictionary`的数据结构，即一组键值对.

但是JavaScript的对象有个小问题，就是键必须是字符串。但实际上Number或者其他数据类型作为键也是非常合理的。

为了解决这个问题，最新的ES6规范引入了新的数据类型`Map`和`Set`.

```js
'use strict';
var m = new Map();
var s = new Set();
alert('你的浏览器支持Map和Set！');
```
```js
var m = new Map([['Michael', 95], ['Bob', 75], ['Tracy', 85]]);
m.set('Adam', 67); // 添加新的key-value
m.has('Adam'); // 是否存在key 'Adam': true
m.get('Adam'); // 67
m.delete('Adam'); // 删除key 'Adam'
m.get('Adam'); // undefined
```
Set和Map类似，也是一组key的集合，但不存储value。由于key不能重复，所以，在Set中，没有重复的key。add(key)方法可以添加元素到Set中，可以重复添加，但不会有效果;通过delete(key)方法可以删除元素.

## 逻辑控制

### 条件判断

如果条件判断语句块只包含一条语句，那么可以省略`{}`：
```js
var age = 20;
if (age >= 18)
    alert('adult');
else
    alert('teenager');
```
**省略`{}`的危险之处在于，如果后来想添加一些语句，却忘了写{}，就改变了if...else...的语义,因此不推荐省略.**

在if...else if...else...语句中，如果某个条件成立，则后续就不再继续判断了.

**JavaScript把null、undefined、0、NaN和空字符串''视为false，其他值一概视为true.**

### 循环

`for ... in`循环是for循环的一个变体.
```js
var o = {
    name: 'Jack',
    age: 20,
    city: 'Beijing'
};
for (var key in o) {
    if (o.hasOwnProperty(key)) {
        alert(key); // 'name', 'age', 'city'
    }
}
//由于Array也是对象，而它的每个元素的索引被视为对象的属性，因此，for ... in循环可以直接循环出Array的索引：

var a = ['A', 'B', 'C'];
for (var i in a) {
    alert(i); // '0', '1', '2'
    alert(a[i]); // 'A', 'B', 'C'
}
```
请注意，**for ... in对Array的循环得到的是String而不是Number**.

`do { ... } while()`循环和while循环的唯一区别在于，不是在每次循环开始的时候判断条件，而是在每次循环完成的时候判断条件.

### iterable

遍历Array可以采用下标循环，遍历Map和Set就无法使用下标。为了统一集合类型，ES6标准引入了新的iterable类型，Array、Map和Set都属于iterable类型。

具有iterable类型的集合可以通过新的`for ... of`循环来遍历。
```js
var a = [1, 2, 3];
  for (var x of a) {
}
```
for ... in循环由于历史遗留问题，它遍历的实际上是对象的属性名称。一个Array数组实际上也是一个对象，它的每个元素的索引被视为一个属性。

当我们手动给Array对象添加了额外的属性后，for ... in循环将带来意想不到的意外效果：
```js
var a = ['A', 'B', 'C'];
a.name = 'Hello';
for (var x in a) {
    alert(x); // '0', '1', '2', 'name'
}
```
for ... in循环将把name包括在内，但Array的length属性却不包括在内。

for ... of循环则完全修复了这些问题，它只循环集合本身的元素：
```js
var a = ['A', 'B', 'C'];
a.name = 'Hello';
for (var x of a) {
    alert(x); 'A', 'B', 'C'
}
```
这就是为什么要引入新的for ... of循环。

然而，更好的方式是直接使用iterable内置的forEach方法，它接收一个函数，每次迭代就自动回调该函数。以Array为例：
```js
var a = ['A', 'B', 'C'];
a.forEach(function (element, index, array) {
    // element: 指向当前元素的值
    // index: 指向当前索引
    // array: 指向Array对象本身
    alert(element);
});
```
注意，forEach()方法是ES5.1标准引入的，你需要测试浏览器是否支持。

Set与Array类似，但Set没有索引，因此其forEach的回调函数的element与index相等.
如果对某些参数不感兴趣，由于JavaScript的函数调用不要求参数必须一致，因此可以忽略它们。例如，只需要获得Array的element：
```js
var a = ['A', 'B', 'C'];
a.forEach(function (element) {
    alert(element);
});
```

### FAQ

#### 相等运算符

JavaScript在设计时，有两种比较相等运算符：

1. `==` : 它会自动转换数据类型再比较，很多时候，会得到非常诡异的结果.
1. `===` : 它不会自动转换数据类型，如果数据类型不一致，返回false，如果一致，再比较。

由于JavaScript这个设计缺陷，**不要使用==比较，始终坚持使用===比较**。

另一个例外是**NaN这个特殊的Number与所有其他值都不相等，包括它自己**：

    NaN === NaN; // false
唯一能判断NaN的方法是通过isNaN()函数：

    isNaN(NaN); // true

因此，不等运算符也有两个`!=`和`!==`,区别与相等运算符类试．

## 函数

于JavaScript的函数也是一个对象，而函数名可以视为指向该函数的变量。

由于JavaScript允许传入任意个参数而不影响调用，因此传入的参数比定义的参数多也没有问题，虽然函数内部并不使用这些参数；传入的参数比定义的少也没有问题，此时参数将收到undefined，要避免收到undefined，可以对参数进行检查。

#### arguments参数
arguments，JavaScript的关键字，它只在函数内部起作用，并且永远指向当前函数的调用者传入的所有参数。arguments类似Array但它不是一个Array，即利用arguments，可以获得调用者传入的所有参数．实际上arguments最常用于判断传入参数的个数．

#### rest参数
ES6标准引入,用于获取已定义参数外的其他参数.rest参数只能写在最后，前面用...标识，从运行结果可知，传入的参数先绑定a、b，多余的参数以数组形式交给变量rest;如果传入的参数连正常定义的参数都没填满，rest参数会接收一个空数组.

chrome 44还未支持rest,firefox最新版已支持.
```js
function foo(a, b, ...rest) {
    console.log('a = ' + a);
    console.log('b = ' + b);
    console.log(rest);
}

foo(1, 2, 3, 4, 5);
// 结果:
// a = 1
// b = 2
// Array [ 3, 4, 5 ]

foo(1);
// 结果:
// a = 1
// b = undefined
// Array []
```

### 变量作用域
在JavaScript中，用var申明的变量实际上是有作用域的。

规则:

- 如果一个变量在函数体内部申明，则该变量的作用域为整个函数体，在函数体外不可引用该变量.
- 如果两个不同的函数各自申明了同一个变量，那么该变量只在各自的函数体内起作用。换句话说，不同函数内部的同名变量互相独立，互不影响.
- 由于JavaScript的函数可以嵌套，此时，内部函数可以访问外部函数定义的变量，反过来则不行.
- JavaScript的函数在查找变量时从自身函数定义开始，从"内"向"外"查找。如果内部函数定义了与外部函数重名的变量，则内部函数的变量将"屏蔽"外部函数的变量。

#### 变量提升

**JavaScript的函数定义有个特点，它会先扫描整个函数体的语句，把所有申明的变量“提升”到函数顶部(提升只对var声明有效,否则仅将其视为对顶层对象的属性的赋值)**：

```js
'use strict';
function foo() {
    var x = 'Hello, ' + y;
    alert(x);
    var y = 'Bob';
}

foo();
```
尽管是strict模式，但语句var x = 'Hello, ' + y;并不报错，原因是变量y在稍后申明了。但是alert显示`Hello, undefined`，说明变量y的值为undefined。这正是因为**JavaScript引擎自动提升了变量y的声明，但不会提升变量y的赋值**。

对于上述foo()函数，JavaScript引擎看到的代码相当于：
```js
function foo() {
    var y; // 提升变量y的申明
    var x = 'Hello, ' + y;
    alert(x);
    y = 'Bob';
}
```
由于JavaScript的这一怪异的“特性”，我们**在函数内部定义变量时，请严格遵守“在函数内部首先申明所有变量”这一规则**。最常见的做法是用一个var申明函数内部用到的所有变量：
```js
function foo() {
    var
        x = 1, // x初始化为1
        y = x + 1, // y初始化为2
        z, i; // z和i为undefined
    // 其他语句:
    for (i=0; i<100; i++) {
        ...
    }
}
```
#### 全局作用域
不在任何函数内定义的变量就具有全局作用域。实际上，JavaScript默认有一个全局对象`window`，全局作用域的变量实际上被绑定到window的一个属性.

顶层函数的定义也被视为一个全局变量，并绑定到window对象.常用的`alert()`函数其实也是`window`的一个变量.

>JavaScript实际上只有一个全局作用域。任何变量（函数也视为变量），如果没有在当前函数作用域中找到，就会继续往上查找，最后如果在全局作用域中也没有找到，则报ReferenceError错误。

>ES6规定，var命令和function命令声明的全局变量，属于全局对象的属性；let命令、const命令、class命令声明的全局变量，不属于全局对象的属性。
>
>     let b = 1;
>     window.b // undefined

#### 命名空间
全局变量会绑定到window上，不同的JavaScript文件如果使用了相同的全局变量，或者定义了相同名字的顶层函数，都会造成命名冲突，并且很难被发现。

**减少冲突的一个方法是把自己的所有变量和函数全部绑定到一个变量(通常是全局变量)中,该变量就叫命名空间**.
#### 局部作用域
由于JavaScript的变量作用域实际上是函数内部，我们在for循环等语句块中是无法定义具有局部作用域的变量的：
```js
'use strict';
function foo() {
    for (var i=0; i<100; i++) {
        //
    }
    i += 100; // 仍然可以引用变量i
}
```
为了解决块级作用域，ES6引入了新的关键字let，用let替代var可以申明一个块级作用域的变量,但let不允许在相同作用域内，重复声明同一个变量否则会报错：
```js
'use strict';
function foo() {
    var sum = 0;
    for (let i=0; i<100; i++) {
        sum += i;
    }
    i += 1; // SyntaxError
}

//let不像var那样，会发生“变量提升”现象
function do_something() {
  console.log(foo); // ReferenceError
  let foo = 2;
}
//---
var a = [];
for (var i = 0; i < 10; i++) {
  a[i] = function () {
    console.log(i);
  };
}
a[6](); // 10

//---使用let，声明的变量仅在块级作用域内有效
var a = [];
for (let i = 0; i < 10; i++) {
  a[i] = function () {
    console.log(i);
  };
}
a[6](); // 6
```
只要块级作用域内存在let命令，它所声明的变量就“绑定”（binding）这个区域，不再受外部的影响。
```js
var tmp = 123;
if (true) {
  tmp = 'abc'; // ReferenceError
  let tmp;
}
```
上面代码中，存在全局变量tmp，但是块级作用域内let又声明了一个局部变量tmp，导致后者绑定这个块级作用域，所以在let声明变量前，对tmp赋值会报错。

**ES6明确规定，如果区块中存在let和const命令，这个区块对这些命令声明的变量，从一开始就形成了封闭作用域。凡是在声明之前就使用这些命令，就会报错。**

**总之，在代码块内，使用let命令声明变量之前，该变量都是不可用的**。这在语法上，称为“暂时性死区”（temporal dead zone，简称TDZ）。

chrome 44还不支持`let`,firefox最新版已支持.
#### 常量
由于var和let申明的是变量，如果要申明一个常量，在ES6之前是不行的，我们通常用全部大写的变量来表示“这是一个常量，不要修改它的值”：

    var PI = 3.14;
ES6标准引入了新的关键字`const`来定义常量，const也不存在提升，只能在声明的位置后面使用,与let一样不可重复声明.**const与let都具有块级作用域**：
```js
'use strict';

const PI = 3.14;
PI = 3; // chrome 44不报错,但没有效果,firefox最新版报错.
PI; // 3.14
```
### 方法
在一个对象的属性上绑定函数，称为这个对象的方法。

**在一个方法内部，this是一个特殊变量，它始终指向当前对象.如果以对象的方法形式调用，比如`obj.xxx()`，该函数的this指向被调用的对象，这是符合我们预期的;如果单独调用函数，比如`xxx()`，此时，该函数的this指向全局对象，也就是window**。
```js
function getAge() {
    var y = new Date().getFullYear();
    return y - this.birth;
}

var xiaoming = {
    name: '小明',
    birth: 1990,
    age: getAge
};

xiaoming.age(); // 25, 正常结果
getAge(); // NaN

var fn = xiaoming.age; // 先拿到xiaoming的age函数
fn(); // NaN,因此要保证this指向正确，必须用obj.xxx()的形式调用.
```
由于这是一个巨大的设计错误，要想纠正可没那么简单。ECMA决定，在strict模式下让函数的this指向undefined并报错(chrome 44仍指向全局对象`window`,不报错;firefox最新版已支持).这个决定只是让错误及时暴露出来，并没有解决this应该指向的正确位置.

有些时候，喜欢重构的你把方法重构了一下：
```js
'use strict';
var xiaoming = {
    name: '小明',
    birth: 1990,
    age: function () {
        function getAgeFromBirth() {
            var y = new Date().getFullYear();
            return y - this.birth;
        }
        return getAgeFromBirth();
    }
};

xiaoming.age(); // Uncaught TypeError: Cannot read property 'birth' of undefined
```
结果又报错了！原因是this指针只在age方法的函数内指向xiaoming，在函数内部定义的函数，this又指向undefined了（在非strict模式下，它重新指向全局对象window）!firefox已支持该方式;chrome 44在非strict和struct都指向window.

修复的办法也不是没有，我们用一个that变量首先捕获this：
```js
'use strict';
var xiaoming = {
    name: '小明',
    birth: 1990,
    age: function () {
        var that = this; // 在方法内部一开始就捕获this
        function getAgeFromBirth() {
            var y = new Date().getFullYear();
            return y - that.birth; // 用that而不是this
        }
        return getAgeFromBirth();
    }
};

xiaoming.age(); // 25
```
#### apply

要指定函数的this指向哪个对象，可以用函数本身的apply方法,即apply() 方法在指定this值和参数的情况下调用某个函数.它接收两个参数，第一个参数就是需要绑定的this变量，第二个参数是Array，表示函数本身的参数。

用apply修复getAge()调用：
```js
function getAge() {
    var y = new Date().getFullYear();
    return y - this.birth;
}

var xiaoming = {
    name: '小明',
    birth: 1990,
    age: getAge
};

xiaoming.age(); // 25
getAge.apply(xiaoming, []); // 25, this指向xiaoming, 参数为空
```
另一个与apply()类似的方法是call()，唯一区别是：

- apply()把参数打包成Array再传入；
- call()把参数按顺序传入。

比如调用Math.max(3, 5, 4)，分别用apply()和call()实现如下：
```js
Math.max.apply(null, [3, 5, 4]); // 5
Math.max.call(null, 3, 5, 4); // 5
```
对普通函数调用，我们通常把this绑定为null。

### 高阶函数

JavaScript的函数其实都指向某个变量。既然变量可以指向函数，函数的参数能接收变量，那么一个函数就可以接收另一个函数作为参数，这种函数就称之为高阶函数.

#### map/reduce
map()方法定义在JavaScript的Array中，我们调用Array的map()方法，传入我们自己的函数，就得到了一个新的Array作为结果：
```js
function pow(x) {
    return x * x;
}

var arr = [1, 2, 3, 4, 5, 6, 7, 8, 9];
arr.map(pow); // [1, 4, 9, 16, 25, 36, 49, 64, 81]
```
Array的reduce()把一个函数作用在这个Array的[x1, x2, x3...]上，这个函数**必须接收两个参数**，reduce()把结果继续和序列的下一个元素做累积计算，其效果就是：

    [x1, x2, x3, x4].reduce(f) = f(f(f(x1, x2), x3), x4)
比方说对一个Array求和，就可以用reduce实现：
```js
var arr = [1, 3, 5, 7, 9];
arr.reduce(function (x, y) {
    return x + y;
}); // 25
```
当Array.length===1时,Array.reduce()返回Array[0],其实是reduce()的回调函数根本就没有执行.

#### filter
用于把Array的某些元素过滤掉，然后返回剩下的元素.

Array的filter()也接收一个函数。和map()不同的是，filter()把传入的函数依次作用于每个元素，然后根据返回值是true还是false决定保留还是丢弃该元素。

例如，把一个Array中的空字符串删掉，可以这么写：
```js
var arr = ['A', '', 'B', null, undefined, 'C', '  '];
arr.filter(function (s) {
    return s && s.trim(); // 注意：IE9以下的版本没有trim()方法
}); // ['A', 'B', 'C']
```
#### sort
默认情况下，对字符串排序，是按照ASCII的大小比较的;数字排序时是所有元素先转换为String再排序.注意,**sort()方法会直接对Array进行修改，它返回的结果仍是当前Array.**

sort()方法也是一个高阶函数，它还可以接收一个比较函数来实现自定义的排序。比较的过程必须通过函数抽象出来,通常规定，对于两个元素x和y，如果认为x < y，则返回-1，如果认为x == y，则返回0，如果认为x > y，则返回1.
```js
var arr = [10, 20, 1, 2];
arr.sort(function (x, y) {
    if (x < y) {
        return -1;
    }
    if (x > y) {
        return 1;
    }
    return 0;
}); // [1, 2, 10, 20]
```
### 闭包
```js
function lazy_sum(arr) {
    var sum = function () {
        return arr.reduce(function (x, y) {
            return x + y;
        });
    }
    return sum;
}
var f = lazy_sum([1, 2, 3, 4, 5]); // function sum()
f(); // 15
var f2 = lazy_sum([1, 2, 3, 4, 5]);
f === f2; // false,每次调用都会返回一个新的函数，即使传入相同的参数
```
在函数lazy_sum中又定义了函数sum，并且，内部函数sum可以引用外部函数lazy_sum的参数和局部变量，当lazy_sum返回函数sum时，相关参数和变量都保存在返回的函数中，这种形式称为"闭包(Closure)".
```js
function count() { //每次循环，都创建了一个新的函数，然后，把创建的3个函数都添加到一个Array中返回
    var arr = [];
    for (var i=1; i<=3; i++) {
        arr.push(function () {
            return i * i; //返回的函数并没有立刻执行，而是直到调用时才执行
        });
    }
    return arr;
}

var results = count();
var f1 = results[0];
var f2 = results[1];
var f3 = results[2];
f1(); // 16
f2(); // 16
f3(); // 16
```
全部都是16！原因就在于返回的函数引用了变量i，但它并非立刻执行。等到3个函数都返回时，它们所引用的变量i已经变成了4，因此最终结果为16。

因此,**返回闭包时牢记的一点就是：返回函数不要引用任何循环变量，或者后续会发生变化的变量。**

如果一定要引用循环变量怎么办？方法是再创建一个函数，用该函数的参数绑定循环变量当前的值，无论该循环变量后续如何更改，已绑定到函数参数的值不变：
```js
function count() {
    var arr = [];
    for (var i=1; i<=3; i++) {
        arr.push((function (n) {
            return function () {
                return n * n;
            }
        })(i));
    }
    return arr;
}

var results = count();
var f1 = results[0];
var f2 = results[1];
var f3 = results[2];

f1(); // 1
f2(); // 4
f3(); // 9
```
这里用了一个“创建一个匿名函数并立刻执行”的语法：
```js
(function (x) {
    return x * x;
})(3); // 9
```
闭包其他作用:
```js
function create_counter(initial) {
    var x = initial || 0;
    return {
        inc: function () {
            x += 1;
            return x;
        }
    }
}
var c2 = create_counter(10);
c2.inc(); // 11
c2.inc(); // 12
```
在返回的对象中，实现了一个闭包，该闭包携带了局部变量x，并且，从外部代码根本无法访问到变量x,相当于C++里的私有变量(private修饰一个成员变量)。换句话说，闭包就是携带状态的函数，并且它的状态可以完全对外隐藏起来。
#### 箭头函数

ES6标准新增了一种新的函数：Arraw Function（箭头函数）.chrome44不支持,firefox最新版支持.
```js
x => x * x
//上面的箭头函数相当于：
function (x) {
    return x * x;
}
```
箭头函数类似于匿名函数，并且简化了函数定义。箭头函数有两种格式，一种像上面的，只包含一个表达式，连{ ... }和return都省略掉了。还有一种可以包含多条语句，这时候就不能省略{ ... }和return.而且如果参数不是一个，就需要用括号()括起来：
```js
// 两个参数:
(x, y) => x * x + y * y

// 无参数:
() => 3.14

// 可变参数:
(x, y, ...rest) => {
    var i, sum = x + y;
    for (i=0; i<rest.length; i++) {
        sum += rest[i];
    }
    return sum;
}
```
如果要返回一个对象，就要注意：
```js
// SyntaxError,因为和函数体的{ ... }有语法冲突:
x => { foo: x }

// ok:
x => ({ foo: x })
```
箭头函数看上去是匿名函数的一种简写，但实际上，箭头函数和匿名函数有个明显的区别：箭头函数内部的this是词法作用域，由上下文确定。即**箭头函数完全修复了this的指向，this总是指向词法作用域，也就是外层调用者obj**.

由于this在箭头函数中已经按照词法作用域绑定了，所以，用call()或者apply()调用箭头函数时，无法对this进行绑定，即传入的第一个参数被忽略：
```js
var obj = {
    birth: 1990,
    getAge: function (year) {
        var b = this.birth; // 1990
        var fn = (y) => y - this.birth; // this.birth仍是1990
        return fn.call({birth:2000}, year);
    }
};
obj.getAge(2015); // 25
```
### generator
generator（生成器）是ES6标准引入的新的数据类型。一个generator看上去像一个函数，但可以返回多次.generator和函数不同的是，generator由`function*`定义（注意多出的*号），并且，除了return语句，还可以用yield返回多次。

产生斐波那契数列的函数:
```js
function fib(max) { //函数只能返回一次，所以必须返回一个Array
    var
        t,
        a = 0,
        b = 1,
        arr = [0, 1];
    while (arr.length < max) {
        t = a + b;
        a = b;
        b = t;
        arr.push(t);
    }
    return arr;
}

// 测试:
fib(5); // [0, 1, 1, 2, 3]
fib(10); // [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
//---
function* fib(max) { //创建了一个generator对象，还没有去执行它
    var
        t,
        a = 0,
        b = 1,
        n = 1;
    while (n < max) {
        yield a;
        t = a + b;
        a = b;
        b = t;
        n ++;
    }
    return a;
}
var f = fib(5);
f.next(); // {value: 0, done: false}
f.next(); // {value: 1, done: false}
f.next(); // {value: 1, done: false}
f.next(); // {value: 2, done: false}
f.next(); // {value: 3, done: true}

for (var x of fib(5)) {
    console.log(x); // 依次输出0, 1, 1, 2, 3
}
```
调用generator对象有两个方法:
1. 不断地调用generator对象的next()方法,next()方法会执行generator的代码，然后，每次遇到yield x;就返回一个对象{value: x, done: true/false}，然后“暂停”。返回的value就是yield的返回值，done表示这个generator是否已经执行结束了。如果done为true，则value就是return的返回值,当执行next()超过可执行次数时将返回`{done=true,  value=undefined}`
2. 直接用`for ... of`循环迭代generator对象，这种方式不需要我们自己判断done.

generator可以在执行过程中多次返回，所以它看上去就像一个可以记住执行状态的函数，利用这一点，写一个generator就可以实现需要用面向对象才能实现的功能:
```js
//编写自增的ID
function* next_id() {
    var count = 0;

    while(true){
        count++;
        yield count;
    }
}
```
generator还有另一个巨大的好处，就是把异步回调代码变成“同步”代码。这个好处要等到后面学了AJAX以后才能体会到。

没有generator之前的黑暗时代，用AJAX时需要这么写代码：
```js
ajax('http://url-1', data1, function (err, result) {
    if (err) {
        return handle(err);
    }
    ajax('http://url-2', data2, function (err, result) {
        if (err) {
            return handle(err);
        }
        ajax('http://url-3', data3, function (err, result) {
            if (err) {
                return handle(err);
            }
            return success(result);
        });
    });
});
```
回调越多，代码越难看。

有了generator的美好时代，用AJAX时可以这么写：
```js
try {
    r1 = yield ajax('http://url-1', data1);
    r2 = yield ajax('http://url-2', data2);
    r3 = yield ajax('http://url-3', data3);
    success(r3);
}
catch (err) {
    handle(err);
}
```
看上去是同步的代码，实际执行是异步的。
## javascript标准对象
在JavaScript的世界里，一切都是对象。为了区分对象的类型，我们用typeof操作符获取对象的类型，它总是返回一个字符串：
```js
typeof 123; // 'number'
typeof NaN; // 'number'
typeof 'str'; // 'string'
typeof true; // 'boolean'
typeof undefined; // 'undefined'
typeof Math.abs; // 'function'
typeof null; // 'object'
typeof []; // 'object'
typeof {}; // 'object'
```
可见，number、string、boolean、function和undefined有别于其他类型。特别注意null的类型是object，Array的类型也是object，如果我们**用typeof将无法区分出null、Array和通常意义上的object——{}**。
#### 包装对象
除了这些类型外，JavaScript还提供了包装对象，类似于Java的int和Integer.

number、boolean和string都有包装对象。包装对象用new创建：
```js
var n = new Number(123); // 123,生成了新的包装类型
var b = new Boolean(true); // true,生成了新的包装类型
var s = new String('str'); // 'str',生成了新的包装类型
```
虽然包装对象看上去和原来的值一模一样，显示出来也是一模一样，但他们的类型已经变为object了！所以，包装对象和原始值用===比较会返回false：
```js
typeof new Number(123); // 'object'
new Number(123) === 123; // false

typeof new Boolean(true); // 'object'
new Boolean(true) === true; // false

typeof new String('str'); // 'object'
new String('str') === 'str'; // false
```
所以**不要使用包装对象！尤其是针对string类型**.

如果我们在使用Number、Boolean和String时，此时，Number()、Boolean和String()被当做普通函数，把任何类型的数据转换为number、boolean和string类型：
```js
var n = Number('123'); // 123，相当于parseInt()或parseFloat()
typeof n; // 'number'

var b = Boolean('true'); // true
typeof b; // 'boolean'

var b2 = Boolean('false'); // true! 'false'字符串转换结果为true！因为它是非空字符串！
var b3 = Boolean(''); // false

var s = String(123.45); // '123.45'
typeof s; // 'string'
```
总结一下，有这么几条规则需要遵守：

- 不要使用new Number()、new Boolean()、new String()创建包装对象；
- 用parseInt()或parseFloat()来转换任意类型到number；
- 用String()来转换任意类型到string，或者直接调用某个对象的toString()方法；
- 通常不必把任意类型转换为boolean再判断，因为可以直接写if (myVar) {...}；
- typeof操作符可以判断出number、boolean、string、function和undefined；
- 判断Array要使用Array.isArray(arr)；
- 判断null请使用myVar === null；
- 判断某个全局变量是否存在用typeof window.myVar === 'undefined'；
- 函数内部判断某个变量是否存在用typeof myVar === 'undefined'。
- 任何对象都有toString()方法吗？null和undefined就没有！确实如此，这两个特殊值要除外，虽然null还伪装成了object类型。
- 数值直接调用toString()报SyntaxError.
      ```js
      123.toString(); // SyntaxError,javascript的解析器试图将点操作符解析为浮点数字面值的一部分

      //遇到这种情况，要特殊处理一下：
      123..toString(); // '123', 注意是两个点！
      (123).toString(); // '123'
      ```
### Date
在JavaScript中，Date对象用来表示日期和时间。
```js
var now = new Date(); // 获取系统当前时间
now; // Wed Jun 24 2015 19:49:22 GMT+0800 (CST)
now.getFullYear(); // 2015, 年份
now.getMonth(); // 5, 月份，注意月份范围是0~11，5表示六月
now.getDate(); // 24, 表示24号
now.getDay(); // 3, 表示星期三
now.getHours(); // 19, 24小时制
now.getMinutes(); // 49, 分钟
now.getSeconds(); // 22, 秒
now.getMilliseconds(); // 875, 毫秒数
now.getTime(); // 1435146562875, 以number形式表示的Unix时间戳(精确到毫秒,GMT时区)<=>Date.now()// 老版本IE没有now()方法
```
注意，当前时间是浏览器从本机操作系统获取的时间，所以不一定准确，因为用户可以把当前时间设定为任何值。

如果要创建一个指定日期和时间的Date对象，可以用：
```js
var d = new Date(2015, 5, 19, 20, 15, 30, 123);//最后一个参数是毫秒数
d; // Fri Jun 19 2015 20:15:30 GMT+0800 (CST)
```
**JavaScript的月份范围用整数表示是0~11，0表示一月，1表示二月……，所以要表示6月，我们传入的是5**.

第二种创建一个指定日期和时间的方法是解析一个符合`ISO 8601`格式的字符串：
```js
var d = Date.parse('2015-06-24T19:49:22.875+08:00');
d; // 1435146562875 //但它返回的不是Date对象，而是一个时间戳。不过有时间戳就可以很容易地把它转换为一个Date：
var d = new Date(1435146562875);
d; // Wed Jun 24 2015 19:49:22 GMT+0800 (CST)
```
#### 时区
Date对象表示的时间总是按浏览器所在时区显示的，不过我们既可以显示本地时间，也可以显示调整后的UTC时间：
```js
var d = new Date(1435146562875);
d.toLocaleString(); // '2015/6/24 下午7:49:22'，本地时间（北京时区+8:00），显示的字符串与操作系统设定的格式有关
d.toUTCString(); // 'Wed, 24 Jun 2015 11:49:22 GMT'，UTC时间，与本地时间相差8小时
```
那么在JavaScript中如何进行时区转换呢？实际上，只要我们传递的是一个number类型的时间戳，我们就不用关心时区转换。**任何浏览器都可以把一个时间戳正确转换为本地时间**。
### RegExp
正则表达式是一种用来匹配字符串的强有力的武器。它的设计思想是用一种描述性的语言来给字符串定义一个规则，凡是符合规则的字符串，我们就认为它“匹配”了，否则，该字符串就是不合法的.

因为正则表达式也是用字符串表示的，所以，我们要首先了解如何用字符来描述字符。

在正则表达式中，如果直接给出字符，就是精确匹配。用`\d`可以匹配一个数字，`\w`可以匹配一个字母或数字，所以：

- `'00\d'`可以匹配`'007'`，但无法匹配`'00A'`；
- `'\d\d\d'`可以匹配`'010'`；
- `'\w\w'`可以匹配`'js'`；

`.`可以匹配任意字符，所以：

 -`'js.'`可以匹配`'jsp'`、`'jss'`、`'js!'`等等。

要匹配变长的字符，在正则表达式中，用`{n}`表示`n`个字符，用`{n,m}`表示`n-m`个字符：

来看一个复杂的例子：`\d{3}\s+\d{3,8}`,我们来从左到右解读一下：
1. `\d{3}`表示匹配3个数字，例如`'010'`；
1. `\s`可以匹配一个空格（也包括Tab等空白符），所以\s+表示至少有一个空格，例如匹配`' '`，`'\t\t'`等；
1. `\d{3,8}`表示3-8个数字，例如`'1234567'`。

综合起来，上面的正则表达式可以匹配以任意个空格隔开的带区号的电话号码。

如果要匹配`010-12345`这样的号码呢？由于`-`是特殊字符，在正则表达式中，要用`\`转义，所以，上面的正则是`\d{3}\-\d{3,8}`。

但是，仍然无法匹配`010 - 12345`，因为带有空格。所以我们需要更复杂的匹配方式。
#### 特殊字符
参考:[正则表达式语法](https://msdn.microsoft.com/zh-cn/library/ae5bf541(v=vs.100).aspx)

| 特殊字符 | 说明 |
|--------|--------|
|`.`| 匹配除换行符 \n之外的任何单字符|
|`*`|任意个字符（包括0个）|
|`+`|至少一个字符|
|`?`|0个或1个字符,或指明一个非贪婪限定符|
|`^`|匹配输入字符串的开始位置,除非在方括号表达式中使用，此时它表示不接受该字符集合|
|`$`|匹配输入字符串的结尾位置|
|`()`|标记一个**子表达式**的开始和结束位置|
|`[]`|标记一个中括号表达式的开始和结束|
|`{}`|标记限定符表达式的开始和结束|
|`|`|指明两项或多项之间的一个选择|
#### 进阶
要做更精确地匹配，可以用`[]`表示范围，比如：
- `[0-9a-zA-Z\_]`可以匹配一个数字、字母或者下划线；
- `[0-9a-zA-Z\_]+`可以匹配至少由一个数字、字母或者下划线组成的字符串，比如'a100'，'0_Z'，'js2015'等等；
- `[a-zA-Z\_\$][0-9a-zA-Z\_\$]*`可以匹配由字母或下划线、$开头，后接任意个由一个数字、字母或者下划线、$组成的字符串，也就是JavaScript允许的变量名；
- `[a-zA-Z\_\$][0-9a-zA-Z\_\$]{0, 19}`更精确地限制了变量的长度是1-20个字符（前面1个字符+后面最多19个字符）。

`A|B`可以匹配A或B，所以`[J|j]ava[S|s]cript`可以匹配'JavaScript'、'Javascript'、'javaScript'或者'javascript'。

`^\d`表示必须以数字开头。
`\d$`表示必须以数字结束。

你可能注意到了，`js`也可以匹配'jsp'，但是加上`^js$`就变成了整行匹配，就只能匹配'js'了。
#### RegExp
JavaScript有两种方式创建一个正则表达式：
- 直接通过`/正则表达式/`写出来，
- 通过`new RegExp('正则表达式(字符串)')`创建一个RegExp对象。ES6允许RegExp构造函数接受正则表达式作为参数，这时会返回一个原有正则表达式的拷贝。

```js
var re1 = /ABC\-001/;
var re2 = new RegExp('ABC\\-001');

re1; // /ABC\-001/
re2; // /ABC\-001/
```
注意，如果使用第二种写法，因为字符串的转义问题，字符串的两个`\\`实际上是一个`\`。

RegExp对象的test()方法用于测试给定的字符串是否符合条件:
```js
var re = /^\d{3}\-\d{3,8}$/;
re.test('010-12345'); // true
re.test('010-1234x'); // false
re.test('010 12345'); // false
```
#### ES6
ES6对正则表达式添加了`u`修饰符，含义为“Unicode模式”，用来正确处理大于`\uFFFF`的Unicode字符。也就是说，会正确处理四个字节的UTF-16编码。
```js
/^\uD83D/u.test('\uD83D\uDC2A')
// false
/^\uD83D/.test('\uD83D\uDC2A')
// true
```
上面代码中，“\uD83D\uDC2A”是一个四个字节的UTF-16编码，代表一个字符。但是，ES5不支持四个字节的UTF-16编码，会将其识别为两个字符，导致第二行代码结果为true。加了u修饰符以后，ES6就会识别其为一个字符，所以第一行代码结果为false。

除了u修饰符，ES6还为正则表达式添加了y修饰符，叫做“粘连”（sticky）修饰符。

y修饰符的作用与g修饰符类似，也是全局匹配，后一次匹配都从上一次匹配成功的下一个位置开始。不同之处在于，g修饰符只要剩余位置中存在匹配就可，而y修饰符确保匹配必须从剩余的第一个位置开始，这也就是“粘连”的涵义.**如果同时使用g修饰符和y修饰符，则y修饰符覆盖g修饰符**。
```js
var s = "aaa_aa_a";
var r1 = /a+/g;
var r2 = /a+/y;

r1.exec(s) // ["aaa"]
r2.exec(s) // ["aaa"]

r1.exec(s) // ["aa"]
r2.exec(s) // null
```
上面代码有两个正则表达式，一个使用g修饰符，另一个使用y修饰符。这两个正则表达式各执行了两次，第一次执行的时候，两者行为相同，剩余字符串都是"_aa_a"。由于g修饰没有位置要求，所以第二次执行会返回结果，而y修饰符要求匹配必须从头部开始，所以返回null。

与y修饰符相匹配，ES6的正则对象多了sticky属性，表示是否设置了y修饰符。

ES6为正则表达式新增了flags属性，会返回正则表达式的修饰符。
```js
// ES5的source属性
// 返回正则表达式的正文
/abc/ig.source
// "abc"

// ES6的flags属性
// 返回正则表达式的修饰符
/abc/ig.flags
// 'gi'
```
#### 切分字符串
用正则表达式切分字符串比用固定的字符更灵活，请看正常的切分代码：
```js
'a b   c'.split(' '); // ['a', 'b', '', '', 'c']
```
嗯，无法识别连续的空格，用正则表达式：
```js
'a b   c'.split(/\s+/); // ['a', 'b', 'c']
```
无论多少个空格都可以正常分割。加入,试试：
```js
'a,b, c  d'.split(/[\s\,]+/); // ['a', 'b', 'c', 'd']
```
再加入;试试：
```js
'a,b;; c  d'.split(/[\s\,\;]+/); // ['a', 'b', 'c', 'd']
```
如果用户输入了一组标签，下次记得用正则表达式来把不规范的输入转化成正确的数组。
#### 分组
正则表达式还有提取子串的强大功能。用`()`表示的就是要提取的分组（Group）。比如：

`^(\d{3})-(\d{3,8})$`分别定义了两个组，可以直接从匹配的字符串中提取出区号和本地号码：
```js
var re = /^(\d{3})-(\d{3,8})$/;
re.exec('010-12345'); // ['010-12345', '010', '12345']
re.exec('010 12345'); // null
```
如果正则表达式中定义了组，就可以在RegExp对象上用exec()方法提取出子串来。

exec()方法在匹配成功后，会返回一个Array，第一个元素始终是原始字符串本身，后面的字符串表示匹配成功的子串;exec()方法在匹配失败时返回`null`。

提取子串非常有用。来看一个更凶残的例子：
```js
var re = /^(0[0-9]|1[0-9]|2[0-3]|[0-9])\:(0[0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9]|[0-9])\:(0[0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9]|[0-9])$/;
re.exec('19:05:30'); // ['19:05:30', '19', '05', '30']
```
这个正则表达式可以直接识别合法的时间。但是有些时候，用正则表达式也无法做到完全验证，比如识别日期：
```js
var re = /^(0[1-9]|1[0-2]|[0-9])-(0[1-9]|1[0-9]|2[0-9]|3[0-1]|[0-9])$/;
```
对于'2-30'，'4-31'这样的非法日期，用正则还是识别不了，或者说写出来非常困难，这时就需要程序配合识别了。
####贪婪匹配

需要特别指出的是，正则匹配默认是**贪婪匹配**，也就是匹配尽可能多的字符。举例如下，匹配出数字后面的0：
```js
var re = /^(\d+)(0*)$/;
re.exec('102300'); // ['102300', '102300', '']
```
由于`\d+`采用贪婪匹配，直接把后面的0全部匹配了，结果0*只能匹配空字符串了。

必须让`\d+`采用非贪婪匹配（也就是尽可能少匹配），才能把后面的0匹配出来，加个?就可以让\d+采用非贪婪匹配：
```js
var re = /^(\d+?)(0*)$/;
re.exec('102300'); // ['102300', '1023', '00']
```
####全局搜索

JavaScript的正则表达式还有几个特殊的标志，最常用的是g，表示全局匹配：
```js
var r1 = /test/g;
var r2 = new RegExp('test', 'g');//r2<=>r1
```
全局匹配可以多次执行exec()方法来搜索一个匹配的字符串。当我们指定g标志后，每次运行exec()，正则表达式本身会更新lastIndex属性，表示上次匹配到的最后索引：
```js
var s = 'JavaScript, VBScript, JScript and ECMAScript';
var re=/[a-zA-Z]+Script/g;

// 使用全局匹配:
re.exec(s); // ['JavaScript']
re.lastIndex; // 10

re.exec(s); // ['VBScript']
re.lastIndex; // 20

re.exec(s); // ['JScript']
re.lastIndex; // 29

re.exec(s); // ['ECMAScript']
re.lastIndex; // 44

re.exec(s); // null，直到结束仍没有匹配到
```
全局匹配类似搜索，因此不能使用/^...$/，那样只会最多匹配一次。

正则表达式还可以指定`i`标志，表示**忽略大小写**，`m`标志，表示**执行多行匹配**。
### json
JSON是JavaScript Object Notation的缩写，它是一种**数据交换格式**,用于取代复杂的xml。在JSON中，一共就这么几种数据类型：

- number：和JavaScript的number完全一致；
- boolean：就是JavaScript的true或false；
- string：就是JavaScript的string；
- null：就是JavaScript的null；
- array：就是JavaScript的Array表示方式——[]；
- object：就是JavaScript的{ ... }表示方式。

以及上面的任意组合。

并且，**JSON还定死了字符集必须是`UTF-8`**，表示多语言就没有问题了。为了统一解析，**JSON的字符串规定必须用双引号`""`，Object的键也必须用双引号`""`**。
#### 序列化
```js
var xiaoming = {
    name: '小明',
    age: 14,
    gender: true,
    height: 1.65,
    grade: null,
    'middle-school': '\"W3C\" Middle School',
    skills: ['JavaScript', 'Java', 'Python', 'Lisp']
};

JSON.stringify(xiaoming); // '{"name":"小明","age":14,"gender":true,"height":1.65,"grade":null,"middle-school":"\"W3C\" Middle School","skills":["JavaScript","Java","Python","Lisp"]}'
//要输出得好看一些，可以加上参数，按缩进输出,第二个参数用于控制如何筛选对象的键值，如果我们只想输出指定的属性，可以传入Array：
JSON.stringify(xiaoming, null, '  '); //null表示不筛选
JSON.stringify(xiaoming, ['name', 'skills'], '  ');
//还可以传入一个函数，这样对象的每个键值对都会被函数先处理：
function convert(key, value) {
    if (typeof value === 'string') {
        return value.toUpperCase();
    }
    return value;
}
JSON.stringify(xiaoming, convert, '  ');
//想要精确控制如何序列化，可以定义一个toJSON()的方法，直接返回JSON应该序列化的数据
var xiaoming = {
    name: '小明',
    age: 14,
    gender: true,
    height: 1.65,
    grade: null,
    'middle-school': '\"W3C\" Middle School',
    skills: ['JavaScript', 'Java', 'Python', 'Lisp'],
    toJSON: function () {
        return { // 只输出name和age，并且改变了key：
            'Name': this.name,
            'Age': this.age
        };
    }
};

JSON.stringify(xiaoming); // '{"Name":"小明","Age":14}'
```
#### 反序列化
直接用JSON.parse()把它变成一个JavaScript对象.
```js
JSON.parse('[1,2,3,true]'); // [1, 2, 3, true]
JSON.parse('{"name":"小明","age":14}'); // Object {name: '小明', age: 14}
JSON.parse('true'); // true
JSON.parse('123.45'); // 123.45
//JSON.parse()还可以接收一个函数，用来转换解析出的属性：
JSON.parse('{"name":"小明","age":14}', function (key, value) {
    if (key === 'name') {
        return value + '同学';
    }
    return value;
}); // Object {name: '小明同学', age: 14}
```
## 面向对象
**JavaScript的原型链和Java的Class区别就在，它没有“Class”的概念，所有对象都是实例，所谓继承关系不过是把一个对象的原型指向另一个对象而已**。
```js
// 基于Obj原型创建一个新对象:
var s = Object.create(Obj); //是E5中提出的一种新的对象创建方式
```
### 创建对象
JavaScript对每个创建的对象都会设置一个原型，指向它的原型对象。

当我们用obj.xxx访问一个对象的属性时，JavaScript引擎先在当前对象上查找该属性，如果没有找到，就到其原型对象上找，如果还没有找到，就一直上溯到Object.prototype对象，最后，如果还没有找到，就只能返回undefined。

#### 构造函数
除了直接用`{...}`创建一个对象外，JavaScript还可以用一种构造函数的方法来创建对象。它的用法是，先定义一个构造函数：
```js
function Student(name) {
    this.name = name;
    this.hello = function () {
        alert('Hello, ' + this.name + '!');
    }
}

var xiaoming = new Student('小明'); //用关键字new来调用这个函数，并返回一个对象
xiaoming.name; // '小明'
xiaoming.hello(); // Hello, 小明!
```

注意，如果不写new，这就是一个普通函数，它返回undefined。但是，如果写了new，它就变成了一个构造函数，它绑定的this指向新创建的对象，并默认返回this，也就是说，不需要在最后写`return this;`。

用`new Student()`创建的对象还从原型上获得了一个constructor属性，它指向函数Student本身.
```js
xiaoming.constructor === Student.prototype.constructor; // true
Student.prototype.constructor === Student; // true

Object.getPrototypeOf(xiaoming) === Student.prototype; // true

xiaoming instanceof Student; // true
```
Student.prototype指向的对象就是xiaoming、xiaohong的原型对象，这个原型对象自己还有个属性constructor，指向Student函数本身。

另外，函数Student恰好有个属性prototype指向它的原型，但是xiaoming、xiaohong这些对象可没有prototype这个属性，不过可以用__proto__这个非标准用法来查看。

现在我们就认为xiaoming、xiaohong这些对象“继承”自Student。
```js
xiaoming.name; // '小明'
xiaohong.name; // '小红'
xiaoming.hello; // function: Student.hello()
xiaohong.hello; // function: Student.hello()
xiaoming.hello === xiaohong.hello; // false,xiaoming和xiaohong各自的hello是一个函数，但它们是两个不同的函数，虽然函数名称和代码都是相同的
```
要让创建的对象共享一个hello函数，根据对象的属性查找原则，我们只要把hello函数移动到xiaoming、xiaohong这些对象共同的原型上就可以了，也就是Student.prototype：
```js
function Student(name) {
    this.name = name;
}
Student.prototype.hello = function () {
    alert('Hello, ' + this.name + '!');
};
```
如果一个函数被定义为用于创建对象的构造函数，但是调用时忘记了写new怎么办？

在strict模式下，this.name = name将报错，因为this绑定为undefined，在非strict模式下，this.name = name不报错，因为this绑定为window，于是无意间创建了全局变量name，并且返回undefined，这个结果更糟糕。

所以，调用构造函数千万不要忘记写new。**为了区分普通函数和构造函数，按照约定，构造函数首字母应当大写，而普通函数首字母应当小写，这样，一些语法检查工具如jslint将可以帮你检测到漏写的new**。

最后，我们还可以编写一个createStudent()函数，在内部封装所有的new操作。一个常用的编程模式像这样：
```js
function Student(props) {
    this.name = props.name || '匿名'; // 默认值为'匿名'
    this.grade = props.grade || 1; // 默认值为1
}
Student.prototype.hello = function () {
    alert('Hello, ' + this.name + '!');
};
function createStudent(props) {
    return new Student(props || {})
}
```
这个createStudent()函数有几个巨大的优点：一是不需要new来调用，二是参数非常灵活，可以不传，也可以这么传：
```js
var xiaoming = createStudent({
    name: '小明'
});

xiaoming.grade; // 1
```
如果创建的对象有很多属性，我们只需要传递需要的某些属性，剩下的属性可以用默认值。由于参数是一个Object，我们无需记忆参数的顺序。如果恰好从JSON拿到了一个对象，就可以直接创建出xiaoming。
## 原型继承
基于Student扩展出PrimaryStudent，可以先定义出PrimaryStudent：
```js
function PrimaryStudent(props) {
    // 调用Student构造函数，绑定this变量:
    Student.call(this, props);
    this.grade = props.grade || 1;
}
```
但是，调用了Student构造函数不等于继承了Student，PrimaryStudent创建的对象的原型是：

    new PrimaryStudent() ----> PrimaryStudent.prototype ----> Object.prototype ----> null
必须想办法把原型链修改为：

    new PrimaryStudent() ----> PrimaryStudent.prototype ----> Student.prototype ----> Object.prototype ----> null
这样，原型链对了，继承关系就对了。新的基于PrimaryStudent创建的对象不但能调用PrimaryStudent.prototype定义的方法，也可以调用Student.prototype定义的方法。

想要实现正确的原型链,我们必须借助一个中间对象来实现正确的原型链，这个中间对象的原型要指向Student.prototype。为了实现这一点，中间对象可以用一个空函数F来实现：
```js
// PrimaryStudent构造函数:
function PrimaryStudent(props) {
    Student.call(this, props);
    this.grade = props.grade || 1;
}

// 空函数F:
function F() {
}

// 把F的原型指向Student.prototype:
F.prototype = Student.prototype;

// 把PrimaryStudent的原型指向一个新的F对象，F对象的原型正好指向Student.prototype:
PrimaryStudent.prototype = new F();

// 把PrimaryStudent原型的构造函数修复为PrimaryStudent:
PrimaryStudent.prototype.constructor = PrimaryStudent;

// 继续在PrimaryStudent原型（就是new F()对象）上定义方法：
PrimaryStudent.prototype.getGrade = function () {
    return this.grade;
};

// 创建xiaoming:
var xiaoming = new PrimaryStudent({
    name: '小明',
    grade: 2
});
xiaoming.name; // '小明'
xiaoming.grade; // 2

// 验证原型:
xiaoming.__proto__ === PrimaryStudent.prototype; // true
xiaoming.__proto__.__proto__ === Student.prototype; // true

// 验证继承关系:
xiaoming instanceof PrimaryStudent; // true
xiaoming instanceof Student; // true
```
注意，函数F仅用于桥接，我们仅创建了一个new F()实例，而且，没有改变原有的Student定义的原型链。

如果把继承这个动作用一个extends()函数封装起来，还可以隐藏F的定义，并简化代码：
```js
function extends(Child, Parent) { //这个extends()函数可以复用
    var F = function () {};
    F.prototype = Parent.prototype;
    Child.prototype = new F();
    Child.prototype.constructor = Child;
}

function Student(props) {
    this.name = props.name || 'Unnamed';
}

Student.prototype.hello = function () {
    alert('Hello, ' + this.name + '!');
}

function PrimaryStudent(props) {
    Student.call(this, props);
    this.grade = props.grade || 1;
}

// 实现原型继承链:
extends(PrimaryStudent, Student);

// 绑定其他方法到PrimaryStudent原型:
PrimaryStudent.prototype.getGrade = function () {
    return this.grade;
};
```

小结,JavaScript的原型继承实现方式就是：
1. 定义新的构造函数，并在内部用call()调用希望“继承”的构造函数，并绑定this；
1. 借助中间函数F实现原型链继承，最好通过封装的extends函数完成；
1. 继续在新的构造函数的原型上定义新方法。

## 浏览器
### 浏览器对象

`window`对象不但充当全局作用域，而且表示浏览器窗口。window对象有innerWidth和innerHeight属性，可以获取浏览器窗口的内部宽度和高度。内部宽高是指除去菜单栏、工具栏、边框等占位元素后，用于显示网页的净宽高。对应的，还有一个outerWidth和outerHeight属性，可以获取浏览器窗口的整个宽高。

`navigator`对象表示浏览器的信息，最常用的属性包括：
- navigator.appName：浏览器名称；
- navigator.appVersion：浏览器版本；
- navigator.language：浏览器设置的语言；
- navigator.platform：操作系统类型；
- navigator.userAgent：浏览器设定的User-Agent字符串。

`screen`对象表示屏幕的信息，常用的属性有：
- screen.width：屏幕宽度，以像素为单位；
- screen.height：屏幕高度，以像素为单位；
- screen.colorDepth：返回颜色位数，如8、16、24。

`location`对象表示当前页面的URL信息:
- location.protocol; // 'http'
- location.host; // 'www.example.com'
- location.port; // '8080'
- location.pathname; // '/path/index.html'
- location.search; // '?a=1&b=2'
- location.hash; // 'TOP'

 要加载一个新页面，可以调用location.assign()。如果要重新加载当前页面，调用location.reload()方法非常方便

>location.href是一个属性，要这样使用：
location.href='http://www.example.com'
而location.assign('http://www.example.com') 就是 location.href='http://www.example.com'

>至于,location.replace('http://www.example.com')与前两者的区别是，在replace之后，浏览历史就被清空了（href与assign方法会产生历史记录）。

>建议：使用replace。因为前两者会产生历史记录，而浏览者如果点‘后退’按钮，就会产生'redirection loop'，会被浏览器禁止。

`document`对象表示当前页面。由于HTML在浏览器中以DOM形式表示为树形结构，document对象就是整个DOM数的根节点。用document对象提供的getElementById()和getElementsByTagName()可以按ID获得一个DOM节点和按Tag名称获得一组DOM节点.

JavaScript可以通过document.cookie读取到当前页面的Cookie.如果引入的第三方的JavaScript中存在恶意代码，则www.foo.com网站将直接获取到www.example.com网站的用户登录信息。为了解决这个问题，服务器在设置Cookie时可以使用`httpOnly`，设定了httpOnly的Cookie将不能被JavaScript读取。

`history`对象保存了浏览器的历史记录，JavaScript可以调用history对象的back()或forward ()，相当于用户点击了浏览器的“后退”或“前进”按钮。这个对象属于历史遗留对象，对于现代Web页面来说，由于大量使用AJAX和页面交互,history对象已经不适用.任何情况，你都不应该使用history这个对象了。

### 操作DOM
```js
// 返回ID为'test'的节点：
var test = document.getElementById('test');

// 先定位ID为'test-table'的节点，再返回其内部所有tr节点：
var trs = document.getElementById('test-table').getElementsByTagName('tr');

// 先定位ID为'test-div'的节点，再返回其内部所有class包含red的节点：
var reds = document.getElementById('test-div').getElementsByClassName('red');

// 获取节点test下的所有直属子节点:
var cs = test.children;

// 获取节点test下第一个、最后一个子节点：
var first = test.firstElementChild;
var last = test.lastElementChild;

//第二种方法是使用querySelector()和querySelectorAll()，需要了解selector语法，然后使用条件来获取节点，更加方便：

// 通过querySelector获取ID为q1的节点：
var q1 = document.querySelector('#q1');

// 通过querySelectorAll获取q1节点内的符合条件的所有节点：
var ps = q1.querySelectorAll('div.highlighted > p');
```
#### 更新DOM
一种是修改innerHTML属性，这个方式非常强大，不但可以修改一个DOM节点的文本内容，还可以直接通过HTML片段修改DOM节点内部的子树：
```js
// 获取<p id="p-id">...</p>
var p = document.getElementById('p-id');
// 设置文本为abc:
p.innerHTML = 'ABC'; // <p id="p-id">ABC</p>
// 设置HTML:
p.innerHTML = 'ABC <span style="color:red">RED</span> XYZ';
// <p>...</p>的内部结构已修改
```
用innerHTML时要注意，是否需要写入HTML。如果写入的字符串是通过网络拿到了，要注意对字符编码来避免XSS攻击。

第二种是修改innerText或textContent属性，这样可以自动对字符串进行HTML编码，保证无法设置任何HTML标签：
```js
// 获取<p id="p-id">...</p>
var p = document.getElementById('p-id');
// 设置文本:
p.innerText = '<script>alert("Hi")</script>';
// HTML被自动编码，无法设置一个<script>节点:
// <p id="p-id">&lt;script&gt;alert("Hi")&lt;/script&gt;</p>
```
两者的区别在于读取属性时，innerText不返回隐藏元素的文本，而textContent返回所有文本。另外注意IE<9不支持textContent。
#### 插入DOM
如果这个DOM节点是空的，例如，`<div></div>`，那么，直接使用innerHTML = '<span>child</span>'就可以修改DOM节点的内容，相当于“插入”了新的DOM节点。

如果这个DOM节点不是空的:
```js
var
    js = document.getElementById('js'),
    list = document.getElementById('list');
list.appendChild(js); //因为我们插入的js节点已经存在于当前的文档树，因此这个节点首先会从原先的位置删除，再插入到新的位置

var
    list = document.getElementById('list'),
    haskell = document.createElement('p'); //创建一个新的节点
haskell.id = 'haskell';
haskell.innerText = 'Haskell';
list.appendChild(haskell);

var d = document.createElement('style');
d.setAttribute('type', 'text/css');
d.innerHTML = 'p { color: red }';
document.getElementsByTagName('head')[0].appendChild(d);

var
    list = document.getElementById('list'),
    ref = document.getElementById('python'),
    haskell = document.createElement('p');
haskell.id = 'haskell';
haskell.innerText = 'Haskell';
list.insertBefore(haskell, ref);　//子节点haskell会插入到referenceElement(ref节点)之前
```

#### 删除DOM
```js
// 拿到待删除节点:
var self = document.getElementById('to-be-removed');
// 拿到父节点:
var parent = self.parentElement;
// 删除:
var removed = parent.removeChild(self);
removed === self; // true
```
注意到删除后的节点虽然不在文档树中了，但其实它还在内存中，可以随时再次被添加到别的位置。

当你遍历一个父节点的子节点并进行删除操作时，要注意，children属性是一个只读属性，并且它在子节点变化时会实时更新。

### 操作Form
```js
// <input type="text" id="email">
var input = document.getElementById('email');
input.value; // '用户输入的值'

//.value这种方式可以应用于text、password、hidden以及select。但是，对于单选框和复选框，value属性返回的永远是HTML预设的值，而我们需要获得的实际是用户是否“勾上了”选项，所以应该用checked判断
// <label><input type="radio" name="weekday" id="monday" value="1"> Monday</label>
// <label><input type="radio" name="weekday" id="tuesday" value="2"> Thurthday</label>
var mon = document.getElementById('monday');
var tue = document.getElementById('tuesday');
mon.value; // '1'
tue.value; // '2'
mon.checked; // true或者false
tue.checked; // true或者false

//设置值和获取值类似，对于text、password、hidden以及select，直接设置value就可以：

// <input type="text" id="email">
var input = document.getElementById('email');
input.value = 'test@example.com'; // 文本框的内容已更新
//对于单选框和复选框，设置checked为true或false即可。
```

JavaScript可以以两种方式来处理表单的提交（AJAX方式在后面章节介绍）。

方式一是通过<form>元素的submit()方法提交一个表单，例如，响应一个`<button>`的click事件，在JavaScript代码中提交表单：
```js
<!-- HTML -->
<form id="test-form">
    <input type="text" name="test">
    <button type="button" onclick="doSubmitForm()">Submit</button>
</form>

<script>
function doSubmitForm() {
    var form = document.getElementById('test-form');
    // 可以在此修改form的input...
    // 提交form:
    form.submit();
}
</script>
```
这种方式的缺点是扰乱了浏览器对form的正常提交。浏览器默认点击`<button type="submit">`时提交表单，或者用户在最后一个输入框按回车键。因此，第二种方式是响应`<form>`本身的onsubmit事件，在提交form时作修改：
```js
<!-- HTML -->
<form id="test-form" onsubmit="return checkForm()">
    <input type="text" name="test">
    <button type="submit">Submit</button>
</form>

<script>
function checkForm() {
    var form = document.getElementById('test-form');
    // 可以在此修改form的input...
    // 继续下一步:
    return true;
}
</script>
```
注意要return true来告诉浏览器继续提交，如果return false，浏览器将不会继续提交form，这种情况通常对应用户输入有误，提示用户错误信息后终止提交form。

### 操作文件

在HTML表单中，可以上传文件的唯一控件就是`<input type="file">`。

>注意：当一个表单包含`<input type="file">`时，表单的enctype必须指定为`multipart/form-data` ，method必须指定为post，浏览器才能正确编码并以multipart/form-data格式发送表单的数据。

上传的文件都由后台服务器处理，JavaScript可以在提交表单时对文件扩展名做检查，以便防止用户上传无效格式的文件：
```js
var f = document.getElementById('test-file-upload');
var filename = f.value; // 'C:\fakepath\test.png'
if (!filename || !(filename.endsWith('.jpg') || filename.endsWith('.png') || filename.endsWith('.gif'))) {
    alert('Can only upload image file.');
    return false;
}
```
HTML5的File API提供了File和FileReader两个主要对象，可以获得文件信息并读取文件.
```js
var
    fileInput = document.getElementById('test-image-file'),
    info = document.getElementById('test-file-info'),
    preview = document.getElementById('test-image-preview');
// 监听change事件:
fileInput.addEventListener('change', function () {
    // 清除背景图片:
    preview.style.backgroundImage = '';
    // 检查文件是否选择:
    if (!fileInput.value) {
        info.innerHTML = '没有选择文件';
        return;
    }
    // 获取File引用:
    var file = fileInput.files[0];
    // 获取File信息:
    info.innerHTML = '文件: ' + file.name + '<br>' +
                     '大小: ' + file.size + '<br>' +
                     '修改: ' + file.lastModifiedDate;
    if (file.type !== 'image/jpeg' && file.type !== 'image/png' && file.type !== 'image/gif') {
      alert('不是有效的图片文件!');
        return;
    }
    // 读取文件:
    var reader = new FileReader();
    reader.onload = function(e) {
        var
            data = e.target.result; // 'data:image/jpeg;base64,/9j/4AAQSk...(base64编码)...'
        preview.style.backgroundImage = 'url(' + data + ')';
    };
    // 以DataURL的形式读取文件:
    reader.readAsDataURL(file);
});
```
上面的代码演示了如何通过HTML5的File API读取文件内容。以DataURL的形式读取到的文件是一个字符串，类似于`data:image/jpeg;base64,/9j/4AAQSk...(base64编码)...`，常用于设置图像。如果需要服务器端处理，把字符串base64,后面的字符发送给服务器并用Base64解码就可以得到原始文件的二进制内容。

#### 回调

上面的代码还演示了JavaScript的一个重要的特性就是单线程执行模式。在JavaScript中，浏览器的JavaScript执行引擎在执行JavaScript代码时，总是以单线程模式执行，也就是说，任何时候，JavaScript代码都不可能同时有多于1个线程在执行。

在JavaScript中，执行多任务实际上都是异步调用，比如上面的代码：

    reader.readAsDataURL(file);
就会发起一个异步操作来读取文件内容。因为是异步操作，所以我们在JavaScript代码中就不知道什么时候操作结束，因此需要先设置一个回调函数：

    reader.onload = function(e) {
      // 当文件读取完成后，自动调用此函数:
    };
当文件读取完成后，JavaScript引擎将自动调用我们设置的回调函数。执行回调函数时，文件已经读取完毕，所以我们可以在回调函数内部安全地获得文件内容。

### AJAX
AJAX(Asynchronous JavaScript and XML)不是JavaScript的规范，是指用JavaScript执行异步网络请求。

在现代浏览器上写AJAX主要依靠｀XMLHttpRequest｀对象：
```js
function success(text) {
    var textarea = document.getElementById('test-response-text');
    textarea.value = text;
}

function fail(code) {
    var textarea = document.getElementById('test-response-text');
    textarea.value = 'Error code: ' + code;
}

var request = new XMLHttpRequest(); // 新建XMLHttpRequest对象

request.onreadystatechange = function () { // 状态发生变化时，函数被回调
    if (request.readyState === 4) { // 成功完成
        // 判断响应结果:
        if (request.status === 200) {
            // 成功，通过responseText拿到响应的文本:
            return success(request.responseText);
        } else {
            // 失败，根据响应码判断失败原因:
            return fail(request.status);
        }
    } else {
        // HTTP请求还在继续...
    }
}

// 发送请求:
request.open('GET', '/api/categories');
request.send();

alert('请求已发送，请等待响应...');
```
当创建了XMLHttpRequest对象后，要先设置onreadystatechange的回调函数。在回调函数中，通常我们只需通过readyState === 4判断请求是否完成，如果已完成，再根据status === 200判断是否是一个成功的响应。

XMLHttpRequest对象的open()方法有3个参数，第一个参数指定是GET还是POST，第二个参数指定URL地址，第三个参数指定是否使用异步，默认是true，所以不用写。

>注意，千万不要把第三个参数指定为false，否则浏览器将停止响应，直到AJAX请求完成。如果这个请求耗时10秒，那么10秒内你会发现浏览器处于“假死”状态。

最后调用send()方法才真正发送请求。GET请求不需要参数，POST请求需要把body部分以字符串或者FormData对象传进去。

#### 安全限制
上面代码的URL使用的是相对路径。如果你把它改为`http://www.sina.com.cn/`，再运行，肯定报错。在Chrome的控制台里，还可以看到错误信息。

这是因为浏览器的**同源策略**导致的。默认情况下，JavaScript在发送AJAX请求时，URL的域名必须和当前页面完全一致。

完全一致的意思是，域名要相同（www.example.com和example.com不同），协议要相同（http和https不同），端口号要相同（默认是:80端口，它和:8080就不同）。

解决方法：
1. 通过在同源域名下架设一个代理服务器来转发，JavaScript负责把请求发送到代理服务器:

       'proxy?url=http://www.sina.com.cn'
 代理服务器再把结果返回，这样就遵守了浏览器的同源策略。这种方式麻烦之处在于需要服务器端额外做开发。

1. JSONP，它有个限制，只能用GET请求，并且要求返回JavaScript。这种方式跨域实际上是利用了浏览器允许跨域引用JavaScript资源.

#### JSONP
```js
function refreshPrice(data) {
    var p = document.getElementById('test-jsonp');
    p.innerHTML = '当前价格：' +
        data['0000001'].name +':' +
        data['0000001'].price + '；' +
        data['1399001'].name + ': ' +
        data['1399001'].price;
}

function getPrice() {
    var
        js = document.createElement('script'),
        head = document.getElementsByTagName('head')[0];
    js.src = 'http://api.money.126.net/data/feed/0000001,1399001?callback=refreshPrice';
    head.appendChild(js);
}
```
JSONP通常以函数调用的形式返回，例如，返回JavaScript内容如下：

    foo('data');
这样一来，我们如果在页面中先准备好foo()函数，然后给页面动态加一个`<script>`节点，相当于动态读取外域的JavaScript资源，最后就等着接收回调了。

#### CORS
如果浏览器支持HTML5，那么就可以一劳永逸地使用新的跨域策略：CORS了。

CORS全称Cross-Origin Resource Sharing，是HTML5规范定义的如何跨域访问资源。

**Origin表示本域**，也就是浏览器当前页面的域。当JavaScript向外域（如sina.com）发起请求后，浏览器收到响应后，首先检查Access-Control-Allow-Origin是否包含本域，如果是，则此次跨域请求成功，如果不是，则请求失败，JavaScript将无法获取到响应的任何数据。

假设本域是my.com，外域是sina.com，只要响应头Access-Control-Allow-Origin为http://my.com，或者是*，本次请求就可以成功。

可见，跨域能否成功，取决于对方服务器是否愿意给你设置一个正确的Access-Control-Allow-Origin，**决定权始终在对方手中**。

上面这种跨域请求，称之为“简单请求”。简单请求包括GET、HEAD和POST（POST的Content-Type类型
仅限application/x-www-form-urlencoded、multipart/form-data和text/plain），并且不能出现任何自定义头（例如，X-Custom: 12345），通常能满足90%的需求。

无论你是否需要用JavaScript通过CORS跨域请求资源，你都要了解CORS的原理。最新的浏览器全面支持HTML5。在引用外域资源时，**除了JavaScript和CSS外，都要验证CORS**。例如，当你引用了某个第三方CDN上的字体文件时：
```css
/* CSS */
@font-face {
  font-family: 'FontAwesome';
  src: url('http://cdn.com/fonts/fontawesome.ttf') format('truetype');
```
如果该CDN服务商未正确设置Access-Control-Allow-Origin，那么浏览器无法加载字体资源。

**对于PUT、DELETE以及其他类型如application/json的POST请求，在发送AJAX请求之前，浏览器会先发送一个OPTIONS请求（称为preflighted请求）到这个URL上，询问目标服务器是否接受**：
```html
OPTIONS /path/to/resource HTTP/1.1
Host: bar.com
Origin: http://bar.com
Access-Control-Request-Method: POST
```
服务器必须响应并明确指出允许的Method：
```html
HTTP/1.1 200 OK
Access-Control-Allow-Origin: http://foo.com
Access-Control-Allow-Methods: POST, GET, PUT, OPTIONS
Access-Control-Max-Age: 86400
```
浏览器确认服务器响应的Access-Control-Allow-Methods头确实包含将要发送的AJAX请求的Method，才会继续发送AJAX，否则，抛出一个错误。

由于以POST、PUT方式传送JSON格式的数据在REST中很常见，所以要跨域正确处理POST和PUT请求，服务器端必须正确响应OPTIONS请求。

### Promise
这种“承诺将来会执行”的对象在JavaScript中称为Promise对象,在ES6中被统一规范，由浏览器直接支持。
```js
function test(resolve, reject) {
    var timeOut = Math.random() * 2;
    log('set timeout to: ' + timeOut + ' seconds.');
    setTimeout(function () {
        if (timeOut &lt; 1) {
            log('call resolve()...');
            resolve('200 OK');
        }
        else {
            log('call reject()...');
            reject('timeout in ' + timeOut + ' seconds.');
        }
    }, timeOut * 1000);
}

new Promise(test).then(function (result) {
    console.log('成功：' + result);
}).catch(function (reason) {
    console.log('失败：' + reason);
});
//下面代码等价于上面的Promise对象的连缀写法
var p1 = new Promise(test);
var p2 = p1.then(function (result) {
    console.log('成功：' + result);
});
var p3 = p2.catch(function (reason) {
    console.log('失败：' + reason);
});
```
可见Promise最大的好处是在异步执行的流程中，把执行代码和处理结果的代码清晰地分离了.

romise还可以做更多的事情，比如，有若干个异步任务，需要先做任务1，如果成功后再做任务2，任何任务失败则不再继续并执行错误处理函数。

要串行执行这样的异步任务，不用Promise需要写一层一层的嵌套代码。有了Promise，我们只需要简单地写：

    job1.then(job2).then(job3).catch(handleError);

了串行执行若干异步任务外，Promise还可以并行执行异步任务,用`Promise.all()`实现:
```js
var p1 = new Promise(function (resolve, reject) {
    setTimeout(resolve, 500, 'P1');
});
var p2 = new Promise(function (resolve, reject) {
    setTimeout(resolve, 600, 'P2');
});
// 同时执行p1和p2，并在它们都完成后执行then:
Promise.all([p1, p2]).then(function (results) {
    console.log(results); // 获得一个Array: ['P1', 'P2']
});
```
有些时候，多个异步任务是为了容错。比如，同时向两个URL读取用户的个人信息，只需要获得先返回的结果即可。这种情况下，用`Promise.race()`实现：
```js
var p1 = new Promise(function (resolve, reject) {
    setTimeout(resolve, 500, 'P1');
});
var p2 = new Promise(function (resolve, reject) {
    setTimeout(resolve, 600, 'P2');
});
Promise.race([p1, p2]).then(function (result) {
    console.log(result); // 'P1'
});
```
由于p1执行较快，Promise的then()将获得结果'P1'。p2仍在继续执行，但执行结果将被丢弃。

如果我们组合使用Promise，就可以把很多异步任务以并行和串行的方式组合起来执行。

### canvas
Canvas是HTML5新增的组件，它就像一块幕布，可以用JavaScript在上面绘制各种图表、动画等。

Canvas除了能绘制基本的形状和文本，还可以实现动画、缩放、各种滤镜和像素转换等高级操作。如果要实现非常复杂的操作，考虑以下优化方案：
- 通过创建一个不可见的Canvas来绘图，然后将最终绘制结果复制到页面的可见Canvas中；
- 尽量使用整数坐标而不是浮点数；
- 可以创建多个重叠的Canvas绘制不同的层，而不是在一个Canvas中绘制非常复杂的图；
- 背景图片如果不变可以直接用`<img>`标签并放到最底层。
