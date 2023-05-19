# bochs
ref:
- [bochs一般用法](https://xlem0n.gitee.io/2019/10/23/2019-10-23-%E7%95%AA%E5%A4%96%E7%AF%87%E2%80%94%E2%80%94boch%E8%B0%83%E8%AF%95%E6%96%B9%E6%B3%95-2019/)

**推荐使用qemu-system-i386, FAQ中的`physical memory read error`无法解决**.

## 安装
env: ubuntu 22.04

```bash
$ apt install bochs bochs-sdl bximage # 默认安装的wx在当前环境会崩溃, 因此使用sdl
$ bochs -n # see Bochs 2.7
$ bochs --help cpu
$ bximage # 下面是创建disk image的交互式输入参数example
1
hd
flat
512
test.img
$ bximage # 下面是创建floppy image的交互式输入参数example
1
fd
1.44M
test.img
$ vim boot.asm
org 07C00h                   ; 告诉编译器程序加载到07C00处
       mov ax, cs
       mov ds, ax
       mov es, ax
       call DispStr                    ; 调用显示字符串例程
       jmp $              ; 无限循环
DispStr:
       mov ax, BootMessage
       mov bp, ax                    ; es:bp = 串地址
       mov cx, 16                    ; cx = 串长度
       mov ax, 01301h            ; ah = 13, al = 01h
       mov bx, 000Ch              ; 页号为0(bh = 0) 黑底红字 (bl = 0Ch,高亮)
       mov dl, 0
       int 10h                          ; 10h号中断
       ret
BootMessage:  db "Hello,OS world!"
	times 510-($-$$)   db   0            ; 填充剩下的空间，使生成的二进制代码恰好为512字节
	dw 0xaa55
$ nasm boot.asm -o boot.bin
$ dd if=boot.bin of=test.img bs=512 count=1 conv=notrunc
$ hexdump -C boot.img
$ vim bochsrc # bochsrc example: /usr/share/doc/bochs/examples/bochsrc.gz
# how much memorythe emulated machine will have  
megs: 32  
   
# filename of ROMimages  
romimage:file=/usr/share/bochs/BIOS-bochs-latest  
vgaromimage:file=/usr/share/bochs/VGABIOS-lgpl-latest  
   
# what disk imageswill be used  
floppya:1_44=test.img, status=inserted  
   
# choose the bootdisk.  
boot: floppy  
   
# where do we sendlog messages?  
log: bochsout.txt  
   
# disable themouse  
mouse: enabled=0  
   
# enable keymapping, using Us layout as default  
keyboard:keymap=/usr/share/bochs/keymaps/sdl2-pc-us.map

display_library: sdl2
$ bochs -f bochsrc
```

组件说明:
- bochs-wx: WxWindows plugin for Bochs
- bximage: 制作bochs image工具

自编译(未验证):
```bash
# apt install libx11-dev libxrandr-dev
# LIBS='-lX11'  ./configure \
--enable-debugger \
--enable-iodebug \
--enable-x86-debugger \
--with-x \
--with-x11
# make
# make install
```

自编译参数说明:
1. enable-debugger: 打开bochs自己的调试器
1. ~~enable-disasm~~: 支持反汇编. 2.7没有该option
1. enable-iodebug: 启用io接口调试
1. enable-x86-debugger: 支持x86调试器
1. with-x: 使用x windows
1. with-x11: 使用x11
1. enable-gdb-stub: 支持gdb远程调试, 需要替换enable-debugger

编译可能遇到的错误:
- `error: 'XRRQueryExtension' was not declared in this scope; did you mean 'XQueryExtension'?`

       `vim gui/x.cpp`, 在首行添加`#include <X11/extensions/Xrandr.h>`
- `undefined reference to `XSetForeground'`

     configure时添加`LIBS='-lX11'`  

## rom
ref:
- [Bochs ROM镜像](http://www.bytekits.com/bochs/bochs-rom-images.html)

romimage:
- BIOS-bochs-latest: 默认的rombios, 从地址0xfffe0000开始加载，它的长度正好是128k
- BIOS-bochs-legacy: 旧版本, 不带32位初始化代码的ROM BIOS映像（用于i386和ISA图形卡仿真）. 从地址0xffff0000开始加载，它的长度正好是64k
- VGABIOS-lgpl-latest/vgabios.bin:	Bochs LGPL’d VGA BIOS映像
- vgabios.cirrus.bin:	SeaVGABIOS ROM映像（用于Cirrus适配器）

> vgabios README: /usr/share/doc/vgabios/README.gz

## bochsrc
- display_library: 显示库是显示Bochs VGA屏幕的代码

## FAQ
### ROM: couldn't open ROM imgae file '/usr/share/bochs'
看bochs启动日志提示: "ROM: couldn't open ROM image file '/usr/share/bochs/BIOS-bochs-latest'"

解决:
`apt install bochsbios vgabios`

### bochs启动报`wxWidgets was not used as the configuration interface, so it cannot be used as the display library`
详细log:
```bash
$ bochs -f bochsrc
...
00000000000e[      ] wxWidgets was not used as the configuration interface, so it cannot be used as the display library
00000000000p[      ] >>PANIC<< no alternative display libraries are available
```

没有可用的configuration interface, 编辑bochsrc配置[config_interface](https://github.com/bochs-emu/Bochs/blob/master/bochs/.bochsrc#L42)即可, 在ubuntu 22.04上配置`config_interface: wx`, bochs启动崩溃.

解决方法: 将`config_interface`注释, `display_library`配置为sdl2或x(需安装相应plugin).

### unknown host key name 'XK_0' (wrong keymap ?)
keyboard_mapping与display_library配套

错误配置:
```conf
keyboard:keymap=/usr/share/bochs/keymaps/x11-pc-us.map
display_library:sdl2
```

正确配置:
```conf
keyboard:keymap=/usr/share/bochs/keymaps/sdl2-pc-us.map
display_library:sdl2
```

### 启用"Bochs Enhanced Debugger"
display_library追加`, options="gui_debug"`即可, 已测试x,sdl2

### `bx_dbg_read_linear: physical memory read error (phy=0x0000322f3130, lin=0x00000000322f3130)`
ref:
- [bx_dbg_read_linear: physical memory read error](https://github.com/bochs-emu/Bochs/issues/50)
- [[SeaBIOS] physical memory read error](https://sourceforge.net/p/bochs/discussion/39592/thread/4f3d95a9/)

       可能是bios问题. 将romimage换成/usr/share/seabios/bios.bin后, "physical memory read error"消失, 但bochs还是黑屏. 使用自编译的seabios-1.16.2构建的bin也是黑屏.

0x0000322f3130约在800M位置, 将megs改为1024(2000年时使用megs=32没遇到过该问题), 不再报该错误, 但还是黑屏, 看不到bios界面. 自编译bochs也无法解决该问题, 推荐qemu-system-i386.