# prototype

JavaScript的每个`对象`都继承另一个`对象`，后者称为“原型”（prototype）对象.只有null除外，它没有自己的原型对象.

`原型对象上的所有属性和方法，都能被派生对象共享`. 这就是JavaScript继承机制的基本设计.

`原型对象的属性不是实例对象自身的属性`. 只要修改原型对象，变动就立刻会体现在所有实例对象上.

当实例对象本身没有某个属性或方法的时候，它会到构造函数的prototype属性指向的对象，去寻找该属性或方法;
如果实例对象自身就有某个属性或方法，它就不会再去原型对象寻找这个属性或方法. 这就是原型对象的特殊之处即原型链.

> “原型链”的作用是，读取对象的某个属性时，JavaScript引擎先寻找对象本身的属性，如果找不到，就到它的原型去找，如果还是找不到，就到原型的原型去找。
> 如果直到最顶层的Object.prototype还是找不到，则返回undefined

Object.prototype对象(`Object {}`)的原型是没有任何属性和方法的null对象，而null对象没有自己的原型.

`Object.getPrototypeOf`方法返回一个对象的原型.

`__proto__`属性指向当前对象的原型对象，即构造函数的prototype属性.

获取实例对象obj的原型对象，有三种方法:

1. obj.__proto__ // 不可靠,最新的ES6标准规定，`__proto__`属性只有浏览器才需要部署，其他环境可以不部署
1. obj.constructor.prototype // 不可靠,obj.constructor.prototype在手动改变原型对象时，可能会失效
1. Object.getPrototypeOf(obj) // 推荐

**Object.create()不会调用构造函数**

# constructor

prototype对象有一个constructor属性，`默认指向prototype对象所在的构造函数`,用于分辨原型对象到底属于哪个构造函数.

instanceof运算符返回一个布尔值，表示指定对象是否为某个构造函数的实例,`instanceof运算符只能用于对象，不适用原始类型的值`.

undefined和null不是对象，所以instanceOf运算符总是返回false.

```js
v instanceof Vehicle
// 等同于
Vehicle.prototype.isPrototypeOf(v)
```

new命令通过构造函数新建实例对象，实质就是将实例对象的原型，指向构造函数的prototype属性，然后在实例对象上执行构造函数.
