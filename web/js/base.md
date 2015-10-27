## 立即执行函数[表达式]

立即执行函数表达式（IIFE,Immediately-Invoked Function Expression）,也叫「自执行匿名函数」（self-executing anonymous function）,推荐前面那种叫法,原因是`(function foo(){ /* code */ }());
`,该语句有函数名称.

IIFE常见形式：
```js
//形式A
(function(){
    console.log("test");
})();

//形式B
(function(){
    console.log("test");
}()); //jslint推荐的写法，好处是，能提醒阅读代码的人，这段代码是一个整体.推荐这种写法.
```
### 要理解立即执行函数，需要先理解一些函数的基本概念.

定义一个函数有三种方法:函数声明,函数表达式,Function构造函数.

- 函数声明(function declaration),它有一个重要特征就是**函数声明提升**:
```js
function name([param,[, param,[..., param]]]) {
   [statements]
}
```
```js
hoisted(); // logs "foo",因为‘提升’了函数声明，函数调用可在函数声明之前
function hoisted() {
  console.log("foo");
}
function fnName(){
    alert('Hello World');
}();//报错，因后面的括号(这里的括号表示分组操作符,其里面内容不允许为空)导致javascript引擎无法解析函数声明
function funcName(){
  console.log("balabalabala");
}("bala")//这里不报错但不会调用，因为()在这里并不是立即调用函数的意思
//上面和下面的相同
function funcName(){
  console.log("balabalabala");
}
("bala")
```
>当括号放置在函数表达式后，此括号即为括号运算符，表示调用函数；然而当括号放置在语句之后意味着分离括号前面与括号中的内容，此时括号仅仅做为分组运算符（即用于改变运算的优先关系）,此时括号内容不能为空.

- 函数表达式(function expression),**不会提升**.函数表达式非常类似于函数声明，并且拥有几乎相同的语法,区别如下,1. 函数名称(最主要区别)，在函数表达式中可忽略它，从而创建匿名函数（anonymous functions）;2. 函数表达式可用于赋值.
```js
function [name]([param1[, param2[, ..., paramN]]]) {
   statements
}
```
```js
var x = function(y) {
   return y * y;
};
//语法也允许为其指定任意一个函数名，当写递归函数时可以调用它自己
var f = function fact(x) {
  if (x < = 1) return 1;
  else return x*fact(x-1);
};
```
```js
x();
var x=function(){...}//报错，被提升的仅仅是变量名x,此时它的定义依然停留在原处。因此在执行 x() 之前，作用域只知道 x 的命名，不知道它到底是什么，所以执行会报错（通常会是：{x} is not a function）.函数表达式只有命名会被提升，定义的函数体则不会，所以函数调用必须在函数表达式之后.
var fnName=function(){
    alert('Hello World');
}();//函数表达式后面加括号，当javascript引擎解析到此处时能立即调用函数<=>先定义fnName,在调用fnName()
function(){
    console.log('Hello World');
}();//语法错误，JavaScript解释器会在默认的情况下把遇到的function关键字当作是函数声明语句(statement)来进行解释.虽然匿名函数属于函数表达式，但是未进行赋值操作，所以报错：要求需要一个函数名(SyntaxError: Unexpected token ())
```

- Function构造函数(不常用),没有函数名
```js
new Function ([arg1[, arg2[, ...argN]],] functionBody)
```
```js
var adder = new Function("a", "b", "return a + b"); // 创建了一个能返回两个参数和的函数
adder(2, 6); // 调用函数
```

[三者区别](http://www.ibm.com/developerworks/cn/web/1406_dengxb_jsfunction/):
1. 从作用域上来说，函数声明和函数表达式使用的是局部变量，而 Function()构造函数却是全局变量.
2. 从性能上来说，Function()构造函数的效率要低于其他两种方式，尤其是在循环体中.使用 Function()构造函数来定义函数允许我们动态的来定义和编译一个函数，而不是限定在 function 预编译的函数体中.但这同时也会带来负面影响，因为每次调用这个函数都要解析函数主体，并创建一个新的函数对象，对性能有一定的影响.
3. 从加载顺序上来说，函数声明式在JavaScript编译的时候就加载到作用域中(函数声明提升),而其他两种方式则是在代码执行的时候加载，如果在定义之前调用，则会返回报错.

### 立即执行函数的两种形式
在JavaScript中()之间只能包含表达式（expression），解释器把()中的内容当作表达式（expression）而不是语句(statement)来执行,同时括号在这种情况下表示的函数会被立即被执行.

形式A是括号使函数表达式立即执行再返回一个函数引用，再通过通过后面的括号调用函数引用来执行函数从而得到返回值;形式B是括号使函数调用立即执行,后得到返回值.两者实现的功能都是相同,没有区别.

其实在function前面加一元操作符号,都将函数声明转换成函数表达式(但返回值有所不同,代码可读性没用括号好),所以同加()的效果一样：
```js
(function (arg){
   console.log(arg);   //输出bala,使用()运算符
})("bala");
(function(arg){
   console.log(arg);   //输出lala，使用()运算符
}("lala"));
!function(arg){
   console.log(arg);   //输出bala,使用！运算符
}("bala");
+function(arg){
    console.log(arg);   //输出lala,使用+运算符
}("lala");
-function(arg){
    console.log(arg);   //输出byebye,使用-运算符
}("byebye");
//还可以用void操作符，~操作符
```
## 模块模式(Module Pattern)
在JavaScript中，每一个函数被调用时都会创建一个执行上下文（execution context）。定义在函数内部的变量和函数都只能在这个执行上下文的内部访问到，所以函数提供了一种创建**私有成员**的便捷的方法。
```js
// makeCounter返回了另外一个内部函数。
// 而这个内部函数可以访问私有变量i，所以这个内部函数实际上拥有一个特权（访问内部私有变量）
function makeCounter() {
    // `i`变量仅在`makeCounter`函数内部有效
    var i = 0;

    return function () {
        console.log(++i);
    };
}

// 注意，`counter` 和 `counter2` 拥有各自的`i`
var counter = makeCounter();
counter(); // logs: 1
counter(); // logs: 2

var counter2 = makeCounter();
counter2(); // logs: 1
counter2(); // logs: 2

i; // ReferenceError: i未定义（i仅在makeCounter中有效）
```

若将上例返回值改成一个对象，也就通常实现单例模式（Singleton Pattern）的方法，如下代码所示：
```js
// 创建一个立即执行的匿名函数表达式，并将函数的返回值赋予一个变量。
// 与上例相比，这个方法略去了`makeWhatever`中间函数。

// 就像上面重要提示中所述，尽管在这个例子中，外面的括号是非必需的。
// 但加上括号可以明确这是以立即执行的函数，将函数的结果赋予变量，而非将函数赋予变量。
var counter = (function () {
    var i = 0;

    return {
        get: function () {
            return i;
        },
        set: function (val) {
            i = val;
        },
        increment: function () {
            return ++i;
        }
    };
}());

// `counter`是一个带有成员的对象，在此例中她的成员都是函数。

counter.get(); // 0
counter.set(3);
counter.increment(); // 4
counter.increment(); // 5

counter.i; // undefined (`i` 并非`counter`的成员)
i; // ReferenceError: i未定义 (仅存在于匿名函数表达式形成的私有作用域中，即闭包)
```
受到上面例子的启发，CHristian Heilmann提出了Revealing Module Pattern(透露模块模式)。他的方法是将所有方法定义为私有变量，也就是说，不在return中定义，但是在那里暴露给用户，如下所示：
```js
var jspy = (function() {
  var _count = 0;
  var incrementCount = function() {
    _count++;
  };
  var resetCount = function() {
    _count = 0;
  };
  var getCount = function() {
    return _count;
  };
  return {
    add: incrementCount,
    reset: resetCount,
    get: getCount
  };
})();
```
这种设计模式有两个好处：

- 首先，它使我们更容易的了解暴露的函数。当你不在return中定义函数时，我们能轻松的了解到每一行就是一个暴露的函数，这时我们阅读代码更加轻松。
- 其次，你可以用简短的名字（例如 add）来暴露函数，但在定义的时候仍然可以使用冗余的定义方法（例如 incrementCount）。

**模块模式**被定义为给类提供私有和公共封装的一种方法，也就是我们常说的“模块化”.

js的模块模式使用了JavaScript中的一个很棒的特性即**闭包**,会**创建对象**.它不仅仅强大并且简洁明了,使用很少的代码就可以有效地将**方法和属性**封装起来，与此同时不污染全局命名空间以及创建私有作用域.上面两个例子就是模块模式.

### 高级模式

#### 拓展
幸运的是我们有一个很好的方式来拓展modules。首先我们导入一个module，然后加属性，最后将它导出。这里的这个例子，就是用上面所说的方法来拓展MODULE。
```js
var MODULE = (function (my) {
    my.anotherMethod = function () {
      // ...
    };
    return my;
}(MODULE));
```
虽然不必要，但是为了一致性 ，我们再次使用var关键字。然后代码执行，module会增加一个叫做MODULE.anotherMethod的公有方法。这个拓展文件同样也维持着它私有的内部状态和导入。

#### 松拓展
上面module模式(拓展)的一个局限性就是整个module必须是写在一个文件里面的。每个进行过大规模代码开发的人都知道将一个文件分离成多个文件的重要性。

我们上面的那个例子需要我们先创建module，然后在对module进行拓展，这并不是必须的,且异步加载脚本是提升 Javascript 应用性能的最佳方式之一。通过松拓展，我们创建灵活的，可以以任意顺序加载的，分成多个文件的module。每个文件的结构大致如下
```js
var MODULE = (function (my) {
    // add capabilities...
    return my;
}(MODULE || {}));
```
在这种模式下，var语句是必须。如果导入的module并不存在就会被创建。这意味着你可以用类似于RequireJS的工具来并行加载这些module的文件。

#### 紧拓展
虽然松拓展已经很棒了，但是它也给你的module增添了一些局限。最重要的一点是，你没有办法安全的重写module的属性，在初始化的时候你也不能使用其他文件中的module属性（但是你可以在初始化之后运行中使用）。紧拓展包含了一定的载入顺序，但是支持重写，下面是一个例子（拓展了我们最初的MODULE）。
```js
var MODULE = (function (my) {
    var old_moduleMethod = my.moduleMethod;
    my.moduleMethod = function () {
       // method override, has access to old through old_moduleMethod...
    };
    return my;
}(MODULE));
```
这里我们已经重写了MODULE.moduleMethod，还按照需求保留了对原始方法的引用。

#### 跨文件的私有状态
把一个module分成多个文件有一很大的局限，就是每一个文件都在维持自身的私有状态，而且没有办法来获得其他文件的私有状态。这个是可以解决的，下面这个松拓展的例子，可以在不同文件中维持私有状态。
```js
var MODULE = (function (my) {
    var _private = my._private = my._private || {},
        _seal = my._seal = my._seal || function () {
            delete my._private;
            delete my._seal;
            delete my._unseal;
        },
        _unseal = my._unseal = my._unseal || function () {
            my._private = _private;
            my._seal = _seal;
            my._unseal = _unseal;
        };
    // permanent access to _private, _seal, and _unseal
    return my;
}(MODULE || {}));
```
每一个文件可以为它的私有变量_private设置属性，其他文件可以立即调用。当module加载完毕，程序会调用MODULE._seal(),让外部没有办法接触到内部的   _.private。如果之后module要再次拓展，某一个属性要改变。在载入新文件前,每一个文件都可以调用_unseal(),，在代码执行之后再调用_seal()。

#### Sub-modules
最后一个高级模式实际上是最简单的，有很多创建子module的例子，就像创建一般的module一样的。
```js
MODULE.sub = (function () {
    var my = {};
    // ...
    return my;
}());
```
子module有一般的module所有优质的特性，包括拓展和私有状态。