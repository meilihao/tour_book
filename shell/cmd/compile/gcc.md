# gcc
gnu编译套件之一.

## 选项:
- -S : 生成汇编

## example
```bash
$ gcc -S hello.c # 生成汇编
```

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
