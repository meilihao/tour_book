# base
通常约定使用 TypeScript 编写的文件以`.ts`为后缀; 用 TypeScript 编写 React 时, 以.tsx`为后缀.

## 类型
原始类型:
```ts
// boolean
let isDone: boolean = false;

let createdByNewBoolean: Boolean = new Boolean(1); // 使用构造函数 Boolean 创造的Boolean对象
let createdByBoolean: boolean = Boolean(1); // 直接调用 Boolean 也可以返回一个 boolean 类型, 不推荐

// number
let decLiteral: number = 6;
let hexLiteral: number = 0xf00d;
// ES6 中的二进制表示法
let binaryLiteral: number = 0b1010;
// ES6 中的八进制表示法
let octalLiteral: number = 0o744;
let notANumber: number = NaN;
let infinityNumber: number = Infinity;

// string
let myName: string = 'Tom';
let myAge: number = 25;

// 模板字符串
let sentence: string = `Hello, my name is ${myName}.
I'll be ${myAge + 1} years old next month.`;

// 在 TypeScript 中，可以用 void 表示没有任何返回值的函数
function alertName(): void {
    alert('My name is Tom');
}

// Null 和 Undefined
// 与 void 的区别是，undefined 和 null 是所有类型的子类型。也就是说 undefined 类型的变量，可以赋值给其他类型的变量, 比如number 类型
let u: undefined = undefined;
let n: null = null;
```

其他:
```ts
// any
// 任意值（Any）用来表示允许赋值为任意类型. 如果是一个普通类型，在赋值过程中改变类型是不被允许的; 但如果是 any 类型，则允许被赋值为任意类型.
// 在任意值上访问任何属性都是允许的, 也允许调用任何方法. 对任意值的任何操作，返回的内容的类型都是任意值.
// 变量如果在声明的时候，未指定其类型，那么它会被识别为任意值类型
let something; // is any. 如果定义的时候没有赋值，不管之后有没有赋值，都会被推断成 any 类型而完全不被类型检查
something = 'seven';

let myFavoriteNumber = 'seven'; // is string by 类型推导

// 联合类型（Union Types）表示取值可以为多种类型中的一种
// 当 TypeScript 不确定一个联合类型的变量到底是哪个类型的时候, 就只能访问此联合类型的所有类型里共有的属性或方法
// 联合类型的变量在被赋值的时候，会根据类型推论的规则推断出一个类型
let myFavoriteNumber: string | number;
myFavoriteNumber = 'seven';
myFavoriteNumber = 7;

// interface
// 赋值的时候，变量的形状必须和接口的形状(属性)保持一致
// `[propName: string]: any`: Person 类型的对象可以有任意属性签名, prop 类似于函数的形参，是可以取其他名字的, string 指的是对象的键都是字符串类型的, any 则是指定了属性值的类型
// 任意属性有两种定义的方式：一种属性签名是 string 类型的, 比如对象的属性; 另一种属性签名是 number 类型的, 比如数组下标.
interface Person {
	readonly id: number; // readonly: 对象中的一些字段只能在创建的时候被赋值
    name: string;
    age?: number; // 可选属性的含义是该属性可以不存在.
    [propName: string]: any; // 一个接口中只能定义一个相同类型的任意属性, 它可用于限制新增属性的值的类型. 如果接口中有多个类型的属性，则可以在任意属性中使用联合类型. 一旦定义了任意属性，那么其他属性(确定属性、可选属性、只读属性等)的类型都必须是它的类型的子集, 任意属性的值的类型必须涵盖其他属性的值的类型. 比如number 不是 string 的子属性(string 也不是 number 的子集), 因此不能将`[propName: string]: any`替换为`[propName: string]: string`
}

let tom: Person = {
	id: 1,
    name: 'Tom',
    age: 25,
    gender: 'male'
};

// 一个接口可以同时定义这两种任意属性，但是 number 类型的签名指定的值类型必须是 string 类型的签名指定的值类型的子集
interface C {
    [prop: string]: object;
    [index: number]: Function; // Function 是 object 的子集
}

// number 类型的任意属性签名不会影响其他 string 类型的属性签名, 可用将其他 string 类型的属性签名理解为该对象的属性.
type Arg = {
    [index: number]: number;
    length: string;
};

let k:Arg = {
	0:100,
	1:101,
	length:2
}
```