# gcc
ref:
- [microarchitecture level](https://developers.redhat.com/blog/2021/01/05/building-red-hat-enterprise-linux-9-for-the-x86-64-v2-microarchitecture-level)

    - rhel9 use x86-64-v2
    - [x86-64 psABI supplement](https://developers.redhat.com/blog/2021/01/05/building-red-hat-enterprise-linux-9-for-the-x86-64-v2-microarchitecture-level#background_of_the_x86_64_microarchitecture_levels)
    - gcc 11/llvm 12 use `-march=` to support microarchitecture level, and please use with glibc 2.33

gnu编译套件之一.

## [选项](https://gcc.gnu.org/onlinedocs/gcc/Option-Index.html#Option-Index):
- -ansi : 只支持ANSI标准的c语言语法
- -c : 只编译不链接, 生成目标文件`.o`
- -Dmacro : 定义指定的宏, 使它能够通过源码中的`#ifdef`进行处理
- -E : 仅预编译
- -g : 生成调试信息, gnu调试器可使用该信息
- -ggdb: 生成gdb专用调试信息, 使用最适合的格式(stabs等), 会有一些gdb专用的扩展
- -gdwarf2: 附带输出调试信息 
- -fPIC: 构建动态库

    在编译阶段, 告诉编译器生成与位置无关代码(position independent code), 即产生的代码中没有绝对地址, 全部使用相对地址, 故而代码被加载器加载到内存的任意位置都可以正确执行
- -fnopic: 关闭fPIC
- -fnobuiltin: 不接受内建函数(builtin function)
- -fnoomiframepointer:省略栈帧指针, backtrace利用调用栈帧信息把函数调用关系层层遍历出来
- -fnostackprotector: 关闭fstackprotector
- -fstackprotector: 提供缓冲区溢出检查机制
- -fnostrictaliasing: 关闭strict aliasing. 当使用strict aliasing时, 编译器会认为在不同类型之间的转换不会发生, 因此执行更激进的编译优化, 比如重新安排执行顺序
- -gN : 生成的可执行程序中含有调试信息, N是1~3, 默认是2, N越大调试信息越多. 与优化选项`-O[N]`冲突
- -I dir : 在头文件的搜索路径中添加dir
- -L dir : 在库文件的搜索路径中添加dir

    LD_LIBRARY_PATH定义了库搜索路径
- -lxxx : 链接名为libxxx.so的库文件
- -m32: 生成32位机器上的代码
- -MD: 生成当前编译程序文件关联的详细信息, 保含目标文件所依赖的所有源码文件, 包括头文件, 但是信息输出将导入`.d`的文件中
- -MM : 自动查找源文件中包含的头文件并输出相关Makefile格式的依赖信息
- -nostdinc: 使编译器不再系统默认的头文件目录中找文件头, 一般和`-I`联合使用, 明确限定头文件的位置
- -o : 指定输出的文件名, 默认是`a.out`
- -O0 : 不进行优化
- -O 或 -O1 : 优化生成代码
- -O2 : 进一步优化
- -O3 : -O2 更进一步优化，包括 inline 函数
- -pedantic : 要求严格遵守ansi标准, 否则发出警告
- -pipe : 编译过程中使用pipe, 加快编译速度
- -P * : 表示在预处理阶段的输出中，阻止生成换行符, 但生成的结果没有换行和排版
    
    `gcc -E -P maco_expand.c -o result.c`
- -shared: 生成共享目标文件, 通常用在建立共享库时
- -static : 链接静态库
- -save-temps : 保留中间文件：如预处理后的结果文件、汇编代码文件与目标文件

    `gcc -c -save-temps main.c`, 会生成预处理后main.i文件
- -S : 只编译不汇编, 生成汇编代码
- -v : 输出详细信息
- -w : 禁止所有的报警
- -Wall : 生成所有告警信息, 在发生警告时停止编译, 即将警告看作错误, **推荐**
- -Werror : 在发生警告时停止编译, 即将警告看作错误
- -Wl, 表示后面的参数将传给 link 程序 ld, 比如`-Wl,rpath=./`可在编译时指定so查找路径.
- --out-implib,dlltest.lib 表示让ld 生成一个名为 dlltest.lib 的导入库

默认情况下, gcc优先使用动态链接, 不存在时才考虑静态链接(需`-static`选项)

主要调试选项:
- -fdump-tree-xxx : 输出gcc编译过程中与ast, gimple等树节点中间表示相关的调试信息
- -fdump-ipa-xxx : 输出与IPA相关的调试信息
- -fdump-rtl-xxx : 输出与RTL(Register transfer language, 寄存器传输语言)中间表示相关的调试信息

## example
```bash
$ gcc -S hello.c # 生成汇编
```

## gcc所遵循的部分约定规则
参考:
- [Linux Makefile 生成 *.d 依赖文件以及 gcc -M -MF -MP 等相关选项说明](https://blog.csdn.net/QQ1452008/article/details/50855810)

- .c : C语言源代码文件
- .a : 由目标文件构成的档案库文件
- .C，.cc或.cxx  : C++源代码文件
- .h : 程序所包含的头文件
- .i  : 已经预处理过的C源代码文件
- .ii : 已经预处理过的C++源代码文件
- .m : Objective-C源代码文件
- .o : 编译后的目标文件
- .s : 汇编语言源代码文件
- .S : 经过预编译的汇编语言源代码文件
- .d : gcc-generated dependency files即gcc生成的描述依赖的文件

## 编译gcc9.1.0
[Building GCC 9 on Ubuntu Linux](https://solarianprogrammer.com/2016/10/07/building-gcc-ubuntu-linux/)

```sh
$ ./configure -v --build=x86_64-linux-gnu --host=x86_64-linux-gnu --target=x86_64-linux-gnu --prefix=/usr/local/gcc-9.1 --enable-checking=release --enable-languages=c,c++ --disable-multilib --program-suffix=-9.1 // 配置选项
```

编译好后更新env:
```sh
export PATH=/usr/local/gcc-9.1/bin:$PATH
export LD_LIBRARY_PATH=/usr/local/gcc-9.1/lib64:$LD_LIBRARY_PATH
```

## FAQ
### [-g、-ggdb、-g3和-ggdb3, -gdwarf-4之间的区别](3.10 Options for Debugging Your Program)
-g和-ggdb之间只有细微的区别:
具体来说，-g产生的debug信息是OS native format， GDB可以使用之, 而-ggdb产生的debug信息更倾向于给GDB使用的. 因此，如果是使用GDB调试器的话，那么使用-ggdb选项. 如果是其他调试器，则使用-g.

3只是包含调试信息的级别(3已是最详细). 这个级别会产生更多的额外debug信息, 比如这个级别可以调试宏.

-gdwarf-<version> : debug信息的格式. 大多数target上的默认版本是4, DWARF5仅是实验性的.

### 编译gcc : `C++ preprocessor "/lib/cpp" fails sanity check`
系统没有c++编译器, 解决方法:
```
# apt install g++
```
### 编译gcc : `cannot run /bin/bash ../.././gmp/config.sub`
或提示:
```
...
configure: error: Building GCC requires GMP 4.2+, MPFR 2.4.0+ and MPC 0.8.0+.
Try the --with-gmp, --with-mpfr and/or --with-mpc options to specify
...
```

依赖的`${gcc_src_root}/gmp-6.1.0/config.sub`不存在
删除`${gcc_src_root}`重新准备`${gcc_src_root}`, 再执行`./contrib/download_prerequisites`, 其提示`All prerequisites downloaded successfully.`即可.
### `./contrib/download_prerequisites`慢
download_prerequisites里的`base_url='ftp://gcc.gnu.org/pub/gcc/infrastructure/'`是使用ftp协议, 换成http即可; 或使用[git mirror](https://www.gnu.org/software/gcc/mirrors.html)提供的`infrastructure`
### `./contrib/download_prerequisites` 校验和不匹配
删除损坏的相应文件, 重新运行`./contrib/download_prerequisites`下载即可.
### 查看gcc的编译选项
```
$ gcc -v
```
### 查看gcc编译优化的具体选项
```
$ gcc -c -Q -O${x} --help=optimizers // 查看 Os/2/3 不同级别优化具体选项
$ gcc -c -Q -march=native  --help=optimizers // 查看目标架构为native 及当前根据cpuid自己选择合适优化选项的具体参数
```

### [ubuntu安装gcc-10](https://launchpad.net/~ubuntu-toolchain-r/+archive/ubuntu/ppa)
```bash
# sudo apt install gcc-10
```

### 编译gcc时重新执行"make check"
error-log:
```conf
Running /tmp/tmp.B0fetZ1B2h.gcc/gcc/testsuite/gcc.c-torture/execute/builtins/builtins.exp ...
FAIL: gcc.c-torture/execute/builtins/fprintf.c execution,  -O0
```

方法:
```
# make check-gcc RUNTESTFLAGS="builtins.exp=fprintf.c -v -v" # `-v` 为输出详细log
```

### ubuntu arm64编译couchdb 3.1.1成功但无法运行, couchdb.log报: undefined ucol_strcollIter_52
排错过程:
1. `grep -r "ucol_strcollIter_52" /opt/couchdb`, 发现couch_icu_driver.so有该symbol, 通过搜索发现`ucol_strcollIter()`由libicuXX提供
1. 通过`locate icu|grep 52`, 发现其他人曾经编过libicu52, 但当前ubuntu 16.04已不再支持libicu52, 果断删除locate输出的相关文件
1. 重新编译发现`ucol_strcollIter_52`还在, 在`/usr`, `/lib*`下搜索`ucol_strcollIter_52`, 没有发现
1. 在couchdb_source搜索`couch_icu_driver`, 发现了`src/couch/compile_commands.json`里面有生成couch_icu_driver.o的命令
1. 将上面发现命令中的`-o priv/icu_driver/couch_icu_driver.o`替换为`-E -o t.i`, 搜索t.i发现从`/usr/local/include/unicode/ucol.h`引入了`ucol_strcollIter_52`, 应该是之前libicu52的残留, 删除`/usr/local/include/unicode`目录, 最后重新编译即可.

> `/usr/local/include/unicode/ucol.h`本身定义的`ucol_strcollIter`没有`_52`后缀, 应该是gcc编译时根据`.h`做了处理, 推测是`.h`中的`# U_ICU_VERSION_SUFFIX _52`+宏展开的原因

### unrecognized command line option '-fmacro-prefix-map=.'
gcc 8才开始支持'-fmacro-prefix-map'

### 获取编译linux内核所用的编译器信息
`cat /proc/version`

### 安装gcc 12
Ubuntu 22.04:
```bash
# apt install gcc-12
# update-alternatives --config gcc # 已有gcc配置项
# update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-12 20 --slave /usr/bin/g++ g++ /usr/bin/g++-12 # 新增gcc配置项
```

### 编译报"undefined reference to `xxx'"
ref:
- ["undefined reference to" 问题解决方法](https://blog.csdn.net/aiwoziji13/article/details/7330333)

原因: 链接时缺少相关的库文件（.a/.so）

解决方法:
1. 追加`-I <lib_header_path> -L <lib_path> -l<lib_name>`
1. 追加` <lib_path>`

    在mingw64环境编译时常遇到该情况, 该情况下优先使用`.dll.a`再考虑`.a`
1. 多个库文件链接顺序问题

    注意库之间的依赖顺序，依赖其他库的库一定要放到被依赖库的前面
1. 定义与实现不一致

    比如在c++代码中链接c语言的库时未添加`extern "C"`

### "undefined reference to `WinMain'"
代码中不存在入口函数即main()函数

### 编译rocksdb 7.10.2报`unrecognized command line optin '-std=c++17'`
ref:
- [g++: error: unrecognized '-std=c++17' (what is g++ version and how to install)](https://stackoverflow.com/questions/60336940/g-error-unrecognized-std-c17-what-is-g-version-and-how-to-install)
- [C++ Standards Support in GCC](https://gcc.gnu.org/projects/cxx-status.html)
- [C++ compiler support](https://en.cppreference.com/w/cpp/compiler_support)

env: gcc-c++-4.8.5

需要gcc-c++需要支持c++17, 至少是v5.0

```bash
yum install centos-release-scl -y # yum -y install oracle-softwarecollection-release-el7
yum install devtoolset-11 -y

# 临时覆盖系统原有的gcc引用
scl enable devtoolset-11 bash
gcc --version

# 永久
echo "source /opt/rh/devtoolset-11/enable" >>/etc/profile
```

C++版本   GCC版本支持情况   GCC版本   指定版本的命令标志
C++98   完全支持    是GCC 6.1之前版本的默认模式   -std=c++98 or -std=gnu++98
C++11   完全支持    从GCC4.8.1版本开始完全支持   -std=c++11 or -std=gnu++11
C++14   完全支持    从GCC 6.1版本开始完全支持，是GCC 6.1到GCC 10 (包括) 的默认模式 -std=c++14 or -std=gnu++14
C++17   完全支持    从GCC 5版本开始，到GCC 7版本，已基本完全支持。 是GCC 11版本的默认模式 -std=c++17 or -std=gnu++17
C++20   未完全支持   从GCC 8版本开始陆续支持C++20特性   -std=c++20 or -std=gnu++20 （GCC9及以前使用-std=c++2a）
C++23   未完全支持（标准还在发展中）  从GCC 11版本开始支持C++23特性    -std=c++2b or -std=gnu++2b

### 查看GCC/GLIB版本信息
查看方法, 前提是没有strip so:
`strings libdhnetsdk.so | grep GLIB`
`strings libdhnetsdk.so | grep GCC`

### `__stack_chk_fail_local`
在makefile CFLAGS中加入`-fno-stack-protector`

注意是在gcc编译时加上参数，不是在ld链接时加上. 需要`make clean + make`

### `skipping incompatible /usr/lib/gcc/x86_64-linux-gnu/12/libgcc.a when searching for -lgcc`
构建使用了`gcc -m32`, os是64位, 没有32位的`lgcc`, 需要`lib32gcc-12-dev`