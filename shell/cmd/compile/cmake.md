# cmake
CMake是个一个开源的跨平台自动化建构系统，用来管理软件建置的程序，并不依赖于某特定编译器，并可支持多层目录、多个应用程序与多个库。 它用配置文件控制建构过程（build process）的方式和Unix的make相似，只是CMake的配置文件取名为CMakeLists.txt。CMake并不直接建构出最终的软件，而是产生标准的建构档（如Unix的Makefile; Windows Visual C++的projects/workspaces; google的ninja），然后再依一般的建构方式使用.

## cmake源码安装
```bash
yum install gcc gcc-c++ make automake
wget https://cmake.org/files/v3.23/cmake-3.23.1.tar.gz
tar -xf cmake-3.23.1.tar.gz
./bootstrap && gmake && gmake install
cmake --version
```

## cmake编译
```bash
cd <source> (CMakeLists.txt 所在目录)
cmake .
```

## FAQ
### No CMAKE_CXX_COMPILER could be found
- gcc: `apt install g++`
- llvm: `export CXX=clang++-11 && cmake ..`

> llvm: clang is for C, clang++ is for C++.

`export CXX=clang++-11`可用alias代替:
```conf
# vim ~/.bashrc
alias clang="clang-11"
alias clang++="clang++-11"
```

### Unable to find clang libraries
`apt install libclang-11-dev`

### Could NOT find LuaJIT (missing: LUAJIT_LIBRARIES LUAJIT_INCLUDE_DIR)
`apt install luajit libluajit-5.1-dev`

### XXX_FOUND
ref:
- [Cmake之深入理解find_package()的用法](https://zhuanlan.zhihu.com/p/97369704)

为了方便在项目中引入外部依赖包，cmake官方预定义了许多寻找依赖包的Module，在path_to_your_cmake/share/cmake-<version>/Modules目录下。每个以Find<LibaryName>.cmake命名的文件都可以找到一个包。我们也可以在官方文档中查看到哪些库官方已经定义好了，可以直接使用find_package函数进行引用官方文档：[Find Modules](https://cmake.org/cmake/help/latest/manual/cmake-modules.7.html).

### [构建时报`Unknown CMake command "target_link_options"`](https://github.com/seetafaceengine/SeetaFace2/issues/108)
该功能由3.13开始提供, 升级Cmake至3.17后解决问题

### cmake执行报`Does not match the generator used previously: xxx`
之前编译过CMakeLists.txt后，产生了缓存文件CMakeCache.txt, 执行"rm -f `find -name CMakeCache.txt`"删除即可.

### `all warnings being treated as errors`
`cmake --compile-no-warning-as-error` from cmake 3.24

或`cmake -DCOMPILE_WARNING_AS_ERROR=no`

或在CMakeLists.txt里找`-Werror`并去掉它