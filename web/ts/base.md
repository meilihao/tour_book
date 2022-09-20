# base
通常约定使用 TypeScript 编写的文件以`.ts`为后缀; 用 TypeScript 编写 React 时, 以.tsx`为后缀.

## 类型
原始类型:
```ts
// --- boolean
let isDone: boolean = false;

let createdByNewBoolean: Boolean = new Boolean(1); // 使用构造函数 Boolean 创造的Boolean对象
let createdByBoolean: boolean = Boolean(1); // 直接调用 Boolean 也可以返回一个 boolean 类型, 不推荐

// --- number
let decLiteral: number = 6;
let hexLiteral: number = 0xf00d;
// ES6 中的二进制表示法
let binaryLiteral: number = 0b1010;
// ES6 中的八进制表示法
let octalLiteral: number = 0o744;
let notANumber: number = NaN;
let infinityNumber: number = Infinity;

// --- string
let myName: string = 'Tom';
let myAge: number = 25;

// 模板字符串
let sentence: string = `Hello, my name is ${myName}.
I'll be ${myAge + 1} years old next month.`;

// 在 TypeScript 中，可以用 void 表示没有任何返回值的函数
function alertName(): void {
    alert('My name is Tom');
}

// --- Null 和 Undefined
// 与 void 的区别是，undefined 和 null 是所有类型的子类型。也就是说 undefined 类型的变量，可以赋值给其他类型的变量, 比如number 类型
let u: undefined = undefined;
let n: null = null;
```

其他:
```ts
// --- any
// 任意值（Any）用来表示允许赋值为任意类型. 如果是一个普通类型，在赋值过程中改变类型是不被允许的; 但如果是 any 类型，则允许被赋值为任意类型.
// 在任意值上访问任何属性都是允许的, 也允许调用任何方法. 对任意值的任何操作，返回的内容的类型都是任意值.
// 变量如果在声明的时候，未指定其类型，那么它会被识别为任意值类型
let something; // is any. 如果定义的时候没有赋值，不管之后有没有赋值，都会被推断成 any 类型而完全不被类型检查
something = 'seven';

let myFavoriteNumber = 'seven'; // is string by 类型推导

// --- 联合类型（Union Types）表示取值可以为多种类型中的一种
// 当 TypeScript 不确定一个联合类型的变量到底是哪个类型的时候, 就只能访问此联合类型的所有类型里共有的属性或方法
// 联合类型的变量在被赋值的时候，会根据类型推论的规则推断出一个类型
let myFavoriteNumber: string | number;
myFavoriteNumber = 'seven';
myFavoriteNumber = 7;

// --- interface
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

// number 类型的任意属性签名不会影响其他 string 类型的属性签名. 它约束当索引的类型是数字时, 值的类型必须是数字之外, 也约束了它还要有一个 length 属性.
type Arg = {
    [index: number]: number;
    length: string;
};

let k:Arg = {
	0:100,
	1:101,
	length:2
}

// --- array
let fibonacci: number[] = [1, 1, 2, 3, 5];
let fibonacci: Array<number> = [1, 1, 2, 3, 5]; // 数组泛型
interface NumberArray {
    [index: number]: number;
}
let fibonacci: NumberArray = [1, 1, 2, 3, 5];
let list: any[] = ['xcatliu', 25, { website: 'http://xcatliu.com' }];

// --- function
// 函数声明（Function Declaration）
function sum(x: number, y: number): number { // 输入多余的（或者少于要求的）参数，是不被允许的
    return x + y;
}
sum(1, 2, 3);

// 函数表达式（Function Expression）
let mySum = function (x: number, y: number): number { // 这里的代码只对等号右侧的匿名函数进行了类型定义，而等号左边的 mySum，是通过赋值操作进行类型推论而推断出来的
    return x + y;
};
等价于:
let mySum: (x: number, y: number) => number = function (x: number, y: number): number { // 在 TypeScript 的类型定义中，=> 用来表示函数的定义，左边是输入类型，需要用括号括起来，右边是输出类型
    return x + y;
};

// 用接口定义函数的形状
interface SearchFunc {
    (source: string, subString: string): boolean;
}

let mySearch: SearchFunc = function(source: string, subString: string) {
    return source.search(subString) !== -1;
}

function buildName(firstName: string, lastName?: string) { // 可选参数后面不允许再出现必需参数
    if (lastName) {
        return firstName + ' ' + lastName;
    } else {
        return firstName;
    }
}
let tomcat = buildName('Tom', 'Cat');
let tom = buildName('Tom');

// TypeScript 会将添加了默认值的参数识别为可选参数, 此时不受「可选参数必须接在必需参数后面」的限制
function buildName(firstName: string = 'Tom', lastName: string) {
    return firstName + ' ' + lastName;
}
let tomcat = buildName('Tom', 'Cat');
let cat = buildName(undefined, 'Cat');

// rest 参数只能是最后一个参数
function push(array: any[], ...items: any[]) {
    items.forEach(function(item) {
        array.push(item);
    });
}

let a = [];
push(a, 1, 2, 3);

// 使用重载定义多个 reverse 的函数类型, 前几次都是函数定义，最后一次是函数实现. TypeScript 会优先从最前面的函数定义开始匹配，所以多个函数定义如果有包含关系，需要优先把精确的定义写在前面.
function reverse(x: number): number;
function reverse(x: string): string;
function reverse(x: number | string): number | string | void {
    if (typeof x === 'number') {
        return Number(x.toString().split('').reverse().join(''));
    } else if (typeof x === 'string') {
        return x.split('').reverse().join('');
    }
}

// --- 类型断言: 值 as 类型
// 类型断言只能够「欺骗」TypeScript 编译器，无法避免运行时的错误. 使用类型断言时一定要格外小心，尽量避免断言后调用方法或引用深层属性，以减少不必要的运行时错误
interface Cat {
    name: string;
    run(): void;
}
interface Fish {
    name: string;
    swim(): void;
}

function isFish(animal: Cat | Fish) {
    if (typeof (animal as Fish).swim === 'function') {
        return true;
    }
    return false;
}

// 将一个父类断言为更加具体的子类
class ApiError extends Error {
    code: number = 0;
}
class HttpError extends Error {
    statusCode: number = 200;
}

function isApiError(error: Error) { // 用类型断言，通过判断是否存在 code 属性，来判断传入的参数是不是 ApiError
    if (typeof (error as ApiError).code === 'number') {
        return true;
    }
    return false;
}

// 作用同上. 但是有的情况下 ApiError 和 HttpError 不是一个真正的类，而只是一个 TypeScript 的接口（interface），接口是一个类型，不是一个真正的值，它在编译结果中会被删除, 当然就无法使用 instanceof 来做运行时判断了
function isApiError(error: Error) {
    if (error instanceof ApiError) {
        return true;
    }
    return false;
}

// 将一个变量断言为 any 可以说是解决 TypeScript 中类型问题的最后一个手段. 它极有可能掩盖了真正的类型错误，所以**如果不是非常确定，就不要使用 as any**
(window as any).foo = 1; // 在 any 类型的变量上，访问任何属性都是允许的

// 将 any 断言为一个具体的类型. any用于处理由于第三方库未能定义好自己的类型，也有可能是历史遗留的或其他人编写的烂代码，还可能是受到 TypeScript 类型系统的限制而无法精确定义类型的场景.
function getCacheData(key: string): any {
    return (window as any).cache[key];
}

interface Cat {
    name: string;
    run(): void;
}

const tom = getCacheData('tom') as Cat;
tom.run();

// 断言规则:
// - 联合类型可以被断言为其中一个类型
// - 父类可以被断言为子类
// - 任何类型都可以被断言为 any
// - any 可以被断言为任何类型
// - 要使得 A 能够被断言为 B，只需要 A 兼容 B 或 B 兼容 A 即可
// - 其实前四种情况都是最后一个的特例

// 双重断言
cat as any as Fish // 容易导致运行时错误, **不推荐**

// 类型断言只会影响 TypeScript 编译时的类型，类型断言语句在编译结果中会被删除. 因此类型断言不是类型转换，它不会真的影响到变量的类型

interface Animal {
    name: string;
}
interface Cat {
    name: string;
    run(): void;
}

const animal: Animal = {
    name: 'tom'
};
let tom = animal as Cat; // 正常: animal 断言为 Cat，只需要满足 Animal 兼容 Cat 或 Cat 兼容 Animal 即可
let tom1: Cat = animal; // 报错: animal 赋值给 tom，需要满足 Cat 兼容 Animal 才行


// 最佳: 通过给 getCacheData 函数添加了一个泛型 <T>，我们可以更加规范的实现对 getCacheData 返回值的约束，这也同时去除掉了代码中的 any
function getCacheData<T>(key: string): T {
    return (window as any).cache[key];
}

interface Cat {
    name: string;
    run(): void;
}

const tom = getCacheData<Cat>('tom');
tom.run();
```