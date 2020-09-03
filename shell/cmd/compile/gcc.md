# gcc
gnu编译套件之一.

## 选项:
- -c : 只编译不链接, 生成目标文件`.o`
- -Dmacro : 定义指定的宏, 使它能够通过源码中的`#ifdef`进行处理
- -E : 仅预编译
- -gN : 生成的可执行程序中含有调试信息, N是1~3, 默认是2, N越大调试信息越多. 与优化选项`-O[N]`冲突
- -I dir : 在头文件的搜索路径中添加dir
- -L dir : 在库文件的搜索路径中添加dir
- -lxxx : 链接名为libxxx.so的库文件
- -MM : 自动查找源文件中包含的头文件并输出相关Makefile格式的依赖信息
- -o : 指定输出的文件名, 默认是`a.out`
- -O0 : 不进行优化
- -O 或 -O1 : 优化生成代码
- -O2 : 进一步优化
- -O3 : -O2 更进一步优化，包括 inline 函数
- -pedantic : 要求严格遵守ansi标准, 否则发出警告
- -pipe : 编译过程中使用pipe, 加快编译速度
- -static : 链接静态库
- -S : 只编译不汇编, 生成汇编代码
- -v : 输出详细信息
- -w : 禁止所有的报警
- -Wall : 在发生警告时取消编译, 即将警告看作错误, **推荐**
- -Werror : 在发生警告时取消编译, 即将警告看作错误
- -Wl, 表示后面的参数将传给 link 程序 ld, 比如`-Wl,rpath=./`可在编译时指定so查找路径.

默认情况下, gcc优先使用动态链接, 不存在时才考虑静态链接(需`-static`选项)

## example
```bash
$ gcc -S hello.c # 生成汇编
```

## gcc所遵循的部分约定规则
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

## 交叉编译
参考:
- [深入理解交叉编译(Cross Compile)](https://www.jianshu.com/p/62613863aed0)

交叉编译中的build、host、target选项的含义:
- --build=编译该软件所使用的平台(你现在电脑的平台), config.guess中猜的
- --host=该软件将运行的平台(编译出来的库可以在哪个平台上运行)
- --target=该软件所处理的目标平台, 即告诉该软件编译出来的工具链生成的代码的运行平台, 即编译程序能够为其生成机器码的平台. 这个选项只有在建立交叉编译环境的时候用到, 比如compile cross-compiler, binutils，toolchain时.

build和host比较好理解,但是target就不好办了.

一般来说,我们平时所说的交差编译用不到他target的,比如`./configure --build=i386-linux --host=arm-linux`就可以了,在386的平台上编译可以运行在arm板的程序.但是,一般我们都是编译程序,而不是编译工具,如果我们编译工具,比如gcc,这个target就有用了.如果我们需要在一台机器上为arm开发板编译一个可以处理mips程序的gcc,那么target就是mips了.

`./configure --build=powerpc --host=i686 --target=mips`: 在powerpc编译, 而编译出的工具在i686运行, 且该工具为mips平台生成机器码.

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