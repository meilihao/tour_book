# ldconfig

## example
```bash
# ldconfig -v
```

## FAQ
### 设置动态库搜索路径的方式
1. 把库拷贝到动态加载器默认搜索目录: `/lib[64]`或`/usr/lib[64]`
1. 在LD_LIBRARY_PATH环境变量中加上追加库所在路径: `export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/test/lib`
1. 通过`/etc/ld.so.conf`

	较新os都是通过ld.so.conf include `/etc/ld.so.conf.d/*.conf`来引入的

	```bash
	# vim /etc/ld.so.conf.d/test.conf
	/test/lib
	# ldconfig # 必须, 用于更新/etc/ld.so.cache
	```
1. 通过gcc的参数`-Wl,-rpath,`指定: `gcc -o pos main.c -L. -lpos -Wl,-rpath,./`

当系统加载可执行代码时候, 除了需要知道其所依赖的库名字, 还需要知道共享库绝对路径, 此时就需要系统动态载入器(dynamic linker/loader, 对于elf格式的可执行程序, 是`ld-linux.so*`)来进行加载.