# ld
ref:
- [充分理解Linux GCC 链接生成的Map文件](https://zhuanlan.zhihu.com/p/502051758)

链接器命令, 将`.o`文件转成可执行文件

## 选项
- -b : 指定目标代码输入文件的格式, 比如binary
- --dynamic-linker /lib/ld-linux.so.2 : 采用32-bit动态连接
- --dynamic-linker /lib64/ld-linux-x86-64.so.2 : 采用64-bit动态连接
- -e : 指定程序入口符号
- -m elf_i386 : 生成 32-bit 的程序
- -L : 将比如lib32-glibc的库加入库搜索路径
- -lc : 连接标准 C 语言库， 比如printf
- -N : 指定可读写的text和data(section)
- -o : 指定输出文件的名称
- -O <level> : 对于非零level, ld将进行level优化
- -T <scriptflie>: 指定链接脚本
- -Ttext 指定text段的加载地址, 不指定默认为0

> lld可使用`-e xxx`指定程序入口; clang不加`-nostdlib`时会连接ld.so即反汇编时main函数前面的部分代码.

## example
```bash
# ld [-L/usr/lib64] -luring --verbose # test liburing, like `gcc -lhdf5 --verbose`
```

## 链接脚本
链接过程可由链接脚本(linker script, 扩展名一般是`.ld/.lds`)控制, 其主要用于规定如何把输入obj文件内的section放入elf文件中, 并控制输出文件内各部分在程序地址空间内的布局.

链接器可生成全局变量`_binary_*_start`和`_binary_*_size`, 以便在程序中定位指令的位置和大小.

常用命令:
- OUTPUT_FORMAT: 输出elf文件的BFD格式
- OUTPUT_ARCH: 指定elf文件头中的机器体系结构
- ENTRY: 指定程序的入口
- `.`: 当前地址
- PROVIDE: 定义符号及其值, 这类符号在目标文件内被引用, 但没有在任何目标文件内被定义的符号. `PROVIDE(etext = .)`定义了符号etext, 其值是当前地址
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

	加载内存地址是程序被加载的地址
	虚拟内存地址是程序运行的地址

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