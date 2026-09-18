# llvm
llvm-project 仓库顶层目录我列一下，每个目录对应一个子项目：

llvm/ ：LLVM 核心，包含中间表示（IR）、优化器、代码生成器、后端（比如 ARM 后端就在 llvm/lib/Target/ARM ），以及 llc、opt、llvm-objdump 这些基础工具。
clang/ ：C/C++ 编译器前端，负责把源码解析成 AST，再生成 LLVM IR。嵌入式相关的驱动逻辑在 clang/lib/Driver ，比如裸机平台工具链参数处理。
lld/ ：链接器，在嵌入式场景下我们主要用的是 ELF 格式的链接支持，对应 lld/ELF 目录。
compiler-rt/ ：编译器运行时库，提供 __aeabi_* 、 __udivsi3 这类编译器隐含调用的辅助函数，以及 libclc_rt.builtins 相关的裸机支持。
libcxx/ 和 libcxxabi/ ：C++ 标准库和 ABI 层。如果你的目标平台使用 C++，这两个模块是必要的。
libunwind/ ：C++ 异常和栈展开支持。
test/ 和 utils/ ：各模块对应的测试集和测试工具，其中 llvm-lit 是贯穿所有模块的测试驱动。

