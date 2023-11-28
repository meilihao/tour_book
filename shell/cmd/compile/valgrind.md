# valgrind
ref:
- [valgrind 内存泄漏分析](https://www.cnblogs.com/gmpy/p/14778243.html)

valgrind 是 Linux 业界主流且非常强大的内存泄漏检查工具.

valgrind 将内存泄漏分为 4 类:
- 明确泄漏（definitely lost）：内存还没释放，但已经没有指针指向内存，内存已经不可访问
- 间接泄漏（indirectly lost）：泄漏的内存指针保存在明确泄漏的内存中，随着明确泄漏的内存不可访问，导致间接泄漏的内存也不可访问
- 可能泄漏（possibly lost）：指针并不指向内存头地址，而是指向内存内部的位置
- 仍可访达（still reachable）：指针一直存在且指向内存头部，直至程序退出时内存还没释放
- suppressed: 统计了使用valgrind的某些参数取消了特定库的某些错误，会被归结到这里

## 选项
- log-file : 输出报告文件
- tool=memcheck 做内存检测就是memcheck, 因为valgrind是一个工具集
- leak-check=full 完整检测
- show-reachable=no 是否显示reachable详见内存泄露部分, 通常是no
- track-origins: 查看未初始化的值来自哪里
- workaround-gcc296-bugs=yes 如果你的gcc存在对应的bug，则要设为yes，否则有误报

## 报告
格式:
```
{问题描述}   
at {地址、函数名、模块或代码行} 
by {地址、函数名、代码行}
by ...{逐层依次显示调用堆栈}
Address 0x???????? {描述地址的相对关系}
```

整体格式:
```
1. copyright 版权声明
2. 异常读写报告
2.1 主线程异常读写
2.2 线程A异常读写报告
2.3 线程B异常读写报告
2... 其他线程
3. 堆内存泄露报告
3.1 堆内存使用情况概述(HEAP SUMMARY)
3.2 确信的内存泄露报告(definitely lost)
3.3 可疑内存操作报告 (show-reachable=no关闭)
3.4 泄露情况概述(LEAK SUMMARY)
```

## 使用
```bash
valgrind --tool=memcheck --leak-check=full --show-leak-kinds=all --log-file=leak.log ./valtest
```