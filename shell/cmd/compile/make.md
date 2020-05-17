# make
最常用的构建工具

## 选项
- -B : make 命令不会编译那些自从上次编译之后就没有更改的文件, 但此参数会忽略该设定, 全部重新编译
- -C : 将当前工作目录转移到指定的位置, 再还行该目录下的Makefile, 最后返回原目录
- -d : 打印详细信息
- -f : 指定Makefile文件名称
- M= : 当用户需要以某个内核为基础编译一个外部模块的话，需要在make modules 命令中加入`M=dir`, 程序会自动到指定的dir目录中查找模块源码，将其编译，生成ko文件

## makefile
make主要功能就是通过makeflie来实现的. 它定义了各种源文件间的依赖关系, 阐明了源文件如何进行编译.

linux下, 通常用Makefile代替makefile, 通过`configure`来生成. 在命令行执行make时, make默认会在当前目录查找Makefile. 如果使用其他文件作为Makefile则需要用`-f <makefile>`参数明确指明.

### 规则
```makefile
target ... : prerequisites ...
    command
    ...
    ...
```
target也就是一个目标文件，可以是object file，也可以是执行文件, 还可以是一个标签（label）. prerequisites就是，要生成那个target所需要的文件或是目标. command也就是make需要执行的命令（任意的shell命令）. 这是一个文件的依赖关系，也就是说，target这一个或多个的目标文件依赖于prerequisites中的文件，其生成规则定义在 command中. 如果prerequisites中有一个以上的文件比target文件要新，那么command所定义的命令就会被执行. 这就是makefile的规则, 也就是makefile中最核心的内容.

## FAQ
### 了解make时执行了哪些命令
```bash
$ make "V="
```

### make if判断明明正确却没有日志输出
```makefile
pkg:
    @echo "----"
	@if [ -d "$(TOPDIR)/www_replace" ]; then \ # 条件判断明明成功却没输出make执行的命令
		echo "replace www"; \
		rm -rf  $(PKG_DIR)$(PRODUCT_ROOT)/www; \
		cp -r $(TOPDIR)/www_replace $(PKG_DIR)$(PRODUCT_ROOT)/www; \
	fi
```

makefile执行时, `if`即使为true, 里面的命令执行日志也不会输出, 因此建议在if中手动添加`echo`.

### make xxx Is a directory. Stop
Makefile要求每行结尾，一定要确认没有空格，直接是换行.

原因:
```makefile
TOPDIR = $(realpath .) # in docker : `/app/xxx`
```

解决:
```makefile
# in docker : `TOPDIR=/app/xxx
TOPDIR = $(realpath .)
```
