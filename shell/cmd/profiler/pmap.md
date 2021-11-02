# pmap

用于报告进程的内存映射关系，是Linux调试及运维一个很好的工具

参考: [Oracle Solaris 11.2 Information Library (简体中文)/ 用户命令/pmap](https://docs.oracle.com/cd/E56344_01/html/E54075/pmap-1.html#scrolltoc)

## 选项
- -A : 筛选地址范围来显示
- -x : 显示扩展格式
- -d : 显示设备格式
- -q : 不显示header/footer行
- -X : 显示比`-x`更详细的信息. 格式来源`/proc/${PID}/smaps`
- -XX : 显示内核提供的所有信息

## 字段
`-d`:
- mpped : 进程中用于映射到文件的内存总量
- writable/private : 进程私有地址空间的数量
- shared : 进程共享给其他进程的地址空间数量

`-x`:
- Address: 内存开始地址, 以升序显示.
- Kbytes: 映射大小, 即占用内存的字节数（KB）
- RSS: 驻留集(实际被分配的内存)大小, 即保留内存的字节数（KB）, 包括与其他地址空间共享的物理内存
- Dirty: 脏页的字节数（包括共享和私有的）（KB）
- Mode: 内存的权限：read、write、execute、shared、private (写时复制)
- Mapping: 占用内存的文件、或[anon]（分配的内存）、或[stack]（栈）
- Offset: 文件偏移
- Device: 设备名 (major:minor)

`-X`:
- Inode : 内存中所加载的文件所在设备上的inode

> `[anon]` : 在磁盘上没有对应的文件，这些一般都是可执行文件或者动态库里的bss段. 当然有对应文件的mapping也有可能是anonymous，比如文件的数据段.

`pmap <pid>`查看进程的地址空间分布:
1. 第一列: vma的起始地址
1. 第二列: vma的大小
1. 第三列: 属性(r,read; w, write; x, execute; s, shared; p, private)
1. 第四列: 内存映射的文件, [anon]为匿名内存映射

total即virt大小.