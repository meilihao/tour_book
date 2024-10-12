# FAQ

## char* argv[]和char** argv的区别
- `char* argv[]`很明显是字符串数组,那`char** argv`是什么?,其实可以反向理解:

```c
char s1[10]; // s1[0] s1[1]等都是char, s1是char*，等同于&s1[0]
char* s2[10]; // s2[0] s2[1]等都是char*, s2是char**，等同于&s2[0]
```

所以,它们是同一个东西.

> C语言中操作字符串是通过它在内存中的存储单元的首地址进行的，这是字符串的终极本质.

扩展:
`char *const *argv`是什么.

**`const修饰规则`**:先忽略类型名（编译器解析的时候也是忽略类型名），看const离哪个近,离谁近就修饰谁,比如:
```
const int *p; // const修饰*p,p 是指针，*p 是指针指向的对象，不可变
int const *p; // const修饰*p,p 是指针，*p 是指针指向的对象，不可变
int *const p; // const修饰p，p 不可变，p 指向的对象可变
const int *const p; // 前一个const 修饰*p,后一个const 修饰p，指针p 和p 指向的对象都不可变
```

因为`const`是修饰`*argv`,argv是指针,所以argv指向的对象,不可变,即`argv[x]`不可变,`argv[x][y]`可变.

> [情景分析“C语言的const关键字”](http://roclinux.cn/?p=557)