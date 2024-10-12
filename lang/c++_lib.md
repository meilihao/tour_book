# lib
## lock
- [C++ condition variable 和mutex](https://www.cnblogs.com/immortalBlog/p/11680321.html)
- [C++ std::condition_variable 条件变量用法](https://cloud.tencent.com/developer/article/2339238)
- [条件变量condition_variable的使用及陷阱](https://www.cnblogs.com/fenghualong/p/13855360.html)

std::unique_lock vs std::lock_guard:
- lock_guard在构造时或者构造前（std::adopt_lock）就已经获取互斥锁，并且在作用域内保持获取锁的状态，直到作用域结束；而unique_lock在构造时或者构造后（std::defer_lock）获取锁，在作用域范围内可以手动获取锁和释放锁，作用域结束时如果已经获取锁则自动释放锁。
- lock_guard锁的持有只能在lock_guard对象的作用域范围内，作用域范围之外锁被释放，而unique_lock对象支持移动操作，可以将unique_lock对象通过函数返回值返回，这样锁就转移到外部unique_lock对象中，延长锁的持有时间
