# this

## this指向
1. 函数在被直接调用的时候，其中的this指针永远指向全局对象
2. 匿名函数this总是指向全局对象
3. 谁执行函数，this就指向谁
4. 如果函数new了一下，那么就会创建一个对象，并且this指向新创建的对象(new时必定返回一个对象.如果构造函数有return且返回值不是对象或为null则忽略return)

> 全局对象在browser中指window;node REPL中指global;node环境里执行的JS脚本中，顶级的this其实是个空对象，有别于global.

**箭头函数对上下文的绑定([用于对函数内部的上下文this绑定是定义函数所在的作用域的上下文即箭头函数里面根本没有自己的this，而是引用外层的this](http://www.open-open.com/lib/view/open1447222864319.html))
是强制性的，无法通过apply或call方法改变其上下文**

> 由于在严格模式中，函数内部的this不能指向全局对象，默认等于undefined，导致不加new调用会报错（JavaScript不允许对undefined添加属性）

## 绑定 this 的方法

包括:
- call
- apply
- bind

```
// call方法的参数，应该是一个对象。如果参数为空、null和undefined，则默认传入全局对象
// call方法实质上是调用Function.prototype.call
func.call(thisValue, arg1, arg2, ...)
// apply与call类似, 唯一的区别就是，它接收一个数组作为函数执行时的参数
// 通过apply方法，利用Array构造函数将数组的空元素变成undefined
func.apply(thisValue, [arg1, arg2, ...])
// bind与apply,call类似, 但bind比call方法和apply方法更进一步的是，除了绑定this以外，还可以绑定原函数的参数
// bind方法每运行一次，就返回一个新函数
function.bind(thisValue,arg1, arg2,...,argn-1)(argn)
```

```js
[1, 2, 3].slice(0, 1)
// 等同于
Array.prototype.slice.call([1, 2, 3], 0, 1)
// 等同于
var slice = Function.prototype.call.bind(Array.prototype.slice);
slice([1, 2, 3], 0, 1)

// ---理解---
var slice = Function.prototype.call.bind(Array.prototype.slice)
等同于
var slice=fn.bind(obj)
fn =  Function.prototype.call
obj = Array.prototype.slice

fn.bind(obj) 等同于 obj.fn // fn 绑定了 obj ， fn 中的 this 就指向了 obj
// 因此
var slice=Array.prototype.slice.call
```
