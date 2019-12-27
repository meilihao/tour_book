# c++
## 函数
### 函数重载
函数名相同,  参数列表(个数,类型,顺序)不同

注意:
1. 返回值类型与函数重载无关
1. 调用函数时, 实参的隐式类型转换可能产生二义性.

本质: 采用了name mangling(也叫name decoration)的技术, 即c++编译器默认会对符号名(比如函数名)进行改编,修饰, 简单理解就是重命名, 且不同编译器有不同的重命名规则.

> 可通过编译时选禁止优化和不生成调试信息, 再通过反汇编验证.

### func list
- endl

从定义中看出，endl是一个函数模板，它实例化之后变成一个模板函数，作用是输出一个换行符，并立即刷新缓冲区.

## FAQ
### 编译器支持c++进度
- [C++ Standards Support in GCC](https://gcc.gnu.org/projects/cxx-status.html)
- [C++ Support in Clang](http://clang.llvm.org/cxx_status.html)
- [C++ 20的范围稳了，主流编译器的完整支持会在什么时候？](https://www.zhihu.com/question/313443905)
