# cmake
CMake是个一个开源的跨平台自动化建构系统，用来管理软件建置的程序，并不依赖于某特定编译器，并可支持多层目录、多个应用程序与多个库。 它用配置文件控制建构过程（build process）的方式和Unix的make相似，只是CMake的配置文件取名为CMakeLists.txt。CMake并不直接建构出最终的软件，而是产生标准的建构档（如Unix的Makefile; Windows Visual C++的projects/workspaces; google的ninja），然后再依一般的建构方式使用.

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