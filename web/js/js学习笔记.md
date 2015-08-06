## 基础

## 数据类型和变量

### null和undefined

JavaScript的设计者希望用null表示一个空的值，而undefined表示值未定义。事实证明，区分两者的意义不大。**大多数情况下，我们都应该用null。undefined仅仅在判断函数参数是否传递的情况下有用**。

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
```

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

Set与Array类似，但Set没有索引，因此回调函数最多两个参数(element,set本身).
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