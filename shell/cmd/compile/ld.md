# ld
ref:
- [充分理解Linux GCC 链接生成的Map文件](https://zhuanlan.zhihu.com/p/502051758)

链接器命令, 将`.o`文件转成可执行文件

> ld采用AT&T链接脚本语言

## 选项
- -b : 指定目标代码输入文件的格式, 比如binary
- -Bstatic : 只使用静态库
- -Bdynamic : 只使用动态库
- -defsym : 在输出文件中定义指定的全局符号
- --dynamic-linker /lib/ld-linux.so.2 : 采用32-bit动态连接
- --dynamic-linker /lib64/ld-linux-x86-64.so.2 : 采用64-bit动态连接
- -e : 指定程序入口符号
- -m elf_i386 : 生成 32-bit 的程序
- -l : 将指定的要链接的库文件
- -L : 将比如lib32-glibc的库加入库搜索路径
- -lc : 连接标准 C 语言库， 比如printf
- -Map : 生成符号表文件
- -N : 指定可读写的text和data(section)
- -o : 指定输出文件的名称
- -O <level> : 对于非零level, ld将进行level优化
- -S : 忽略来自输出文件的调试器符号信息
- -s : 忽略来自输出文件的所有符号信息
- -t : 在处理输入文件时显示它们的名称
- -T <scriptflie>: 指定链接脚本

	ld有默认的内置链接脚本, 但也可用`-T`指定
- -Ttext : 指定text段的加载地址, 不指定默认为0
- -Tdata : 指定data段的加载地址
- -Tbss : 指定bss段的加载地址

> lld可使用`-e xxx`指定程序入口; clang不加`-nostdlib`时会连接ld.so即反汇编时main函数前面的部分代码.

## example
```bash
# ld -verbose # 显示当前链接器使用的默认链接脚本的信息, 将输出信息重定向到新的lds文件中,并去掉开头和结尾多余的字符(否则会报错)
# ld [-L/usr/lib64] -luring --verbose # test liburing, like `gcc -lhdf5 --verbose`
```

## 链接脚本
链接过程可由链接脚本(linker script, 扩展名一般是`.ld/.lds`)控制, 其主要用于规定如何把输入obj文件内的section放入elf文件中, 并控制输出文件内各部分在程序地址空间内的布局.

链接脚本具有定制最终生成的二进制文件的作用,它可以定制各种不同的段,定义变量,指定各个段的地址等.

> x86 linux链接脚本见[`vmlinux.lds.S`](https://github.com/torvalds/linux/blob/master/arch/x86/kernel/vmlinux.lds.S)

链接器可生成全局变量`_binary_*_start`和`_binary_*_size`, 以便在程序中定位指令的位置和大小.

常用命令:
- OUTPUT_FORMAT: 输出elf文件的BFD格式
- OUTPUT_ARCH: 指定elf文件头中的机器体系结构
- ENTRY: 指定程序的入口

  设置入口优先规则:
  1. `ld -e`
  1. 链接脚本的ENTRY命令
  1. 通过特定符号(start)指定
  1. 使用代码段的起始地址
  1. 设置为0
- `.`: 当前位置计数器(Location Count, LC)
- ALIGN: 字节对齐调整, 0x1000即4KB对
- SECTION: 定义段的链接分布

	段格式:
	```
	SECTION-NAME [ADDR] [(TYPE)] : [AT(LMA)]
	{
		OUTPUT-SECTION-COMMAND
		OUTPUT-SECTION-COMMAND
		...
	}
	```

	选项:
	- SECTION-NAME: 段名称
	- ADDR: 用于设置虚拟内存地址VMA(Virtual Memory Address)
	- AT: 指定加载内存地址(LMA, Load Memory Address)

	  没有指定AT时, VMA=LMA

	加载内存地址是程序被加载的地址
	虚拟内存地址是程序运行的地址

	`*(.text .rodata)`和`*(.text) *(.rodata)`区别:
	1. `*(.text .rodata)`: 按照输入文件的顺序把相应的代码段和只读数据段加入
	1. `*(.text) *(.rodata)`: 先加入所有文件的代码段, 再加入所有文件的只读数据段
- 内置函数

  - ABSOLUTE(exp) : 返回表达式的绝对地址, 主要用于在段定义中给符号赋绝对值
  - ADDR(section) : 返回段的虚拟地址
  - ALIGN(align) : 返回下一个与align对其的地址
  - SIZEOF(section) : 返回一个段的大小
  - INCLUDE : 引入其他的链接脚本
  - LOADADDR(section) : 返回段的加载地址
  - MAX(exp1, exp2) : 返回最大值
  - MIN(exp1, exp2) : 返回最小值
  - PROVIDE: 定义符号及其值, 这类符号可在目标文件内被引用, 但需要其没有在任何目标文件内被定义过. `PROVIDE(etext = .)`定义了符号etext, 其值是当前地址

  	`readelf -s kernel | grep -E "etext|stack0|end"`: 查看PROVIDE定义的地址

  	在链接脚本中, 符号可以像c一样进行赋值和操作, 允许的操作包括赋值, 加法, 减法, 乘法, 除法, 左移, 右移, 与, 或等

    c中声明一个符号时, 会在符号表中创建一个保存该符号地址的条目. 链接脚本定义的符号仅在符号表中创建一个符号, 并没有分配内存来存储这个符号.

    ```
    # cat xxx.ld
    start_of_ROM = .ROM;
    end_of_ROM = .ROM + sizeof(.ROM);
    # cat yyy.c
    extern char start_of_ROM, end_of_ROM; # 使用时用`&`获取符号地址
    extern char start_of_ROM[], end_of_ROM[]; # 直接使用
    ```

xv6 kernel.ld:
```
/* Simple linker script for the JOS kernel.
   See the GNU ld 'info' manual ("info ld") to learn the syntax. */

OUTPUT_FORMAT("elf32-i386", "elf32-i386", "elf32-i386")
OUTPUT_ARCH(i386)
ENTRY(_start)

SECTIONS
{
	/* Link the kernel at this address: "." means the current address */
        /* Must be equal to KERNLINK */
	. = 0x80100000;

	.text : AT(0x100000) {
		*(.text .stub .text.* .gnu.linkonce.t.*)
	}

	PROVIDE(etext = .);	/* Define the 'etext' symbol to this value */

	.rodata : {
		*(.rodata .rodata.* .gnu.linkonce.r.*)
	}

	/* Include debugging information in kernel memory */
	.stab : {
		PROVIDE(__STAB_BEGIN__ = .);
		*(.stab);
		PROVIDE(__STAB_END__ = .);
	}

	.stabstr : {
		PROVIDE(__STABSTR_BEGIN__ = .);
		*(.stabstr);
		PROVIDE(__STABSTR_END__ = .);
	}

	/* Adjust the address for the data segment to the next page */
	. = ALIGN(0x1000);

	/* Conventionally, Unix linkers provide pseudo-symbols
	 * etext, edata, and end, at the end of the text, data, and bss.
	 * For the kernel mapping, we need the address at the beginning
	 * of the data section, but that's not one of the conventional
	 * symbols, because the convention started before there was a
	 * read-only rodata section between text and data. */
	PROVIDE(data = .);

	/* The data segment */
	.data : {
		*(.data)
	}

	PROVIDE(edata = .);

	.bss : {
		*(.bss)
	}

	PROVIDE(end = .);

	/DISCARD/ : {
		*(.eh_frame .note.GNU-stack)
	}
}
```

说明:
- etext: text段结束地址
- data: data段开始地址
- edata: data段结束地址
- end： 内核结束地址

## FAQ
### os启动起点
由ld链接器决定

### `ld cannot find library but it exists`
```bash
# ninjia
...
/usr/bin/c++ ... -lcurl  -Wl,-Bstatic  -luring  -Wl,-Bdynamic  ../thirdParty/isa-l_crypto/.libs/libisal_crypto.a  bin/libhashtable.a  -lm  -lrt  -ldl  -lrdmacm  -libverbs  -lpthread ...
/usr/bin/ld: cannot find -luring
```

liburing is so, use `-lcurl  -Wl,-Bdynamic  -luring`

### 查看libc版本
`/usr/lib/x86_64-linux-gnu/libc.so.6`

### libc.so.6 is needed by XXX
在相同arm64 iso安装的出来的os上分别构建并安装rpm报libc依赖错误. 打包出的rpm中混入了x86 binary导致依赖了其他版本的libc.
