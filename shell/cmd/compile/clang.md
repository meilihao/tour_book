# clang
# clang++
```bash
$ llc --version # 支持的target
$ clang++ -print-resource-dir # 打印当前 clang++ 对应的资源目录（resource dir）
$ clang++ -ffreestanding -E -v main.cpp # 看预处理器搜索顺序, 验证头文件搜索路径
```

## FAQ
### clang++ `-nostdinc++, -nostdinc, -nostdlibinc` 三个参数区别
-nostdinc —— 最激进，禁用全部
禁用范围：
标准系统头文件目录（如 /usr/include）
编译器内置头文件目录（如 lib/clang/*/include，包含 <stdint.h>、<stddef.h>、<stdarg.h> 等编译器内置的 C 头文件）
效果：#include <...> 几乎找不到任何东西，除非你手动用 -I 指定路径
典型用途：交叉编译裸机固件时，彻底切断对宿主系统头文件的依赖
2. -nostdlibinc —— 禁用标准库头文件，但保留编译器内置头文件
禁用范围：
标准 C 库头文件目录（如 /usr/include）
标准 C++ 库头文件目录（如 /usr/include/c++/15 或 /usr/lib/llvm-21/include/c++/v1）
不禁用：
编译器内置头文件目录（lib/clang/*/include）仍然保留
效果：#include <stdio.h> 找不到，但 #include <stdint.h> 仍然能找到（因为 <stdint.h> 由编译器内置提供）
典型用途：裸机开发中最常用的选项——既切断了宿主系统的 libc/libc++，又保留了编译器内置的基础类型定义
3. -nostdinc++ —— 只禁用 C++ 标准库头文件
禁用范围：
仅 C++ 标准库头文件目录（如 /usr/include/c++/15 或 /usr/lib/llvm-21/include/c++/v1）
不禁用：
标准 C 库头文件目录（/usr/include）仍然可用
编译器内置头文件目录仍然可用
效果：#include <vector> 找不到，但 #include <stdio.h> 和 #include <stdint.h> 都正常
典型用途：当你想用自定义的 C++ 标准库实现（如交叉编译的 libc++）时，先用 -nostdinc++ 禁用系统自带的，再用 -I 指向你自己的版本

| 参数 | 标准 C 头文件（`/usr/include`） | C++ 标准库头文件 | 编译器内置头文件（`lib/clang/*/include`） |
|------|-------------------------------|-----------------|----------------------------------------|
| 默认（不加任何参数） | ✅ 可用 | ✅ 可用 | ✅ 可用 |
| `-nostdinc` | ❌ 禁用 | ❌ 禁用 | ❌ 禁用 |
| `-nostdlibinc` | ❌ 禁用 | ❌ 禁用 | ✅ 保留 |
| `-nostdinc++` | ✅ 可用 | ❌ 禁用 | ✅ 保留 |
### cstdint区别
/usr/include/c++/15/cstdint : libstdc++ 标准版,属于 GNU libstdc++（GCC 的 C++ 标准库实现
/usr/include/c++/15/tr1/cstdint : libstdc++ TR1 兼容版（已废弃）
/usr/lib/llvm-21/include/c++/v1/cstdint: libc++ 标准版, 属于libc++（LLVM/Clang 的 C++ 标准库实现）
/usr/lib/llvm-21/include/c++/v1/__cxx03/cstdint : libc++ C++03 兼容版

v1 vs __cxx03 是 libc++ 内部为了同时支持 C++03 和 C++11+ 而做的头文件分离，普通开发者无需直接关心
libstdc++ vs libc++ 是两套不同的库实现，二进制不兼容，混用会出问题
TR1 版是 C++11 之前的过渡产物，现代代码不应再使用

### LLVM_ENABLE_THREADS和LIBCXX_ENABLE_THREADS区别
LLVM_ENABLE_THREADS, 用于 LLVM / Clang 编译器本身
LIBCXX_ENABLE_THREADS, 用于	libc++ (标准库)
