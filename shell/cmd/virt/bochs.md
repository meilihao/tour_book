# bochs
ref:
- [bochs一般用法](https://xlem0n.gitee.io/2019/10/23/2019-10-23-%E7%95%AA%E5%A4%96%E7%AF%87%E2%80%94%E2%80%94boch%E8%B0%83%E8%AF%95%E6%96%B9%E6%B3%95-2019/)

**考虑FAQ中的`physical memory read error`, 推荐汇编测试用bochs, 界面输出测试用qemu-system-i386**.

ps: 从使用情况看, qemu-system-i386比bochs更能模拟真实硬件, 便于发现错误.

## 安装
env: ubuntu 22.04

```bash
$ apt install bochs bochs-sdl bximage # 默认安装的wx在当前环境会崩溃, 因此使用sdl
$ bochs -n # see Bochs 2.7
$ bochs --help cpu
$ bximage # 下面是创建disk image的交互式输入参数example. = `bximage -hd -mode="flat" -size=64 -q hd.img`
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

> hd用`ata0-master: type=disk, path="boot.img", mode=flat, sect_size=4096`+`boot: disk`

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
1. enable-gdb-stub: 支持gdb远程调试, 与enable-debugger互斥

    被调试程序需要开启编译选项`-g`, 链接时不能有`-s`选项, 且bochsrc需要设置`gdbstbub:enable=1`
    使bochs在本地1234端口上监听gdb命令, 并向gdb发送命令执行结果

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
- megs: 内存
- romimage: rom bios
- vga : vga显示配置
- vga romimage: VGA rom bios
- floppya : 软驱a
- floppyb : 软驱b
- ata[0-3]: 硬盘或光驱的ata控制器
- ata[0-3]-master: ata设备的主设备
- ata[0-3]-slave: ata设备的从设备
- boot: 启动驱动器
- ips: 模拟的频率
- log: 调试用的log
- panic: 错误的信息
- error: bochs遇到不能模拟的情况, 比如非法指令
- info: 显示一些不常出现的情况
- debug: 主要用来开发bochs软件时的调试信息
- parport1: 并行端口
- vga_update_interval: vga卡刷新频率
- keyboard_serial_delay: 键盘串行延时
- mouse: 鼠标
- private_colormap: gui色彩映射
- keyboard_mapping: 硬盘映射

## 编译seabios
ref:
- [seabios-1.16.2-1.fc39](https://src.fedoraproject.org/rpms/seabios/tree/rawhide)

```bash
$ tree ~/rpmbuild
/home/chen/rpmbuild
├── BUILD
├── BUILDROOT
├── RPMS
├── SOURCES
│  ├── config.seabios-128k
│  ├── config.vga-bochs-display
│  ├── config.vga-stdvga
│  ├── seabios-1.16.2-1.fc39.src.rpm
│  └── seabios-1.16.2.tar.gz
├── SPECS
│  └── seabios.spec
└── SRPMS

6 directories, 6 files
# %if 0%{?fedora:1}
%define cross 1
%endif

Name:           seabios
Version:        1.16.2
Release:        1%{?dist}
Summary:        Open-source legacy BIOS implementation

License:        LGPLv3
URL:            http://www.coreboot.org/SeaBIOS

Source0:        %{name}-%{version}.tar.gz

Source13:       config.vga-stdvga
Source17:       config.seabios-128k
Source21:       config.vga-bochs-display

%if 0%{?cross:1}
BuildRequires: binutils-x86_64-linux-gnu gcc-x86_64-linux-gnu
Buildarch:     noarch
%else
ExclusiveArch: x86_64
%endif

Requires: %{name}-bin = %{version}-%{release}
Requires: seavgabios-bin = %{version}-%{release}

# Seabios is noarch, but required on architectures which cannot build it.
# Disable debuginfo because it is of no use to us.
%global debug_package %{nil}

# Similarly, tell RPM to not complain about x86 roms being shipped noarch
%global _binaries_in_noarch_packages_terminate_build   0

# You can build a debugging version of the BIOS by setting this to a
# value > 1.  See src/config.h for possible values, but setting it to
# a number like 99 will enable all possible debugging.  Note that
# debugging goes to a special qemu port that you have to enable.  See
# the SeaBIOS top-level README file for the magic qemu invocation to
# enable this.
%global debug_level 1


%description
SeaBIOS is an open-source legacy BIOS implementation which can be used as
a coreboot payload. It implements the standard BIOS calling interfaces
that a typical x86 proprietary BIOS implements.


%package bin
Summary: Seabios for x86
Buildarch: noarch


%description bin
SeaBIOS is an open-source legacy BIOS implementation which can be used as
a coreboot payload. It implements the standard BIOS calling interfaces
that a typical x86 proprietary BIOS implements.


%package -n seavgabios-bin
Summary: Seavgabios for x86
Buildarch: noarch

%description -n seavgabios-bin
SeaVGABIOS is an open-source VGABIOS implementation.


%prep
%setup -q
%autopatch -p1

%build
%define _lto_cflags %{nil}
export CFLAGS="$RPM_OPT_FLAGS"
mkdir binaries

build_bios() {
    make clean distclean
    cp $1 .config
    echo "CONFIG_DEBUG_LEVEL=%{debug_level}" >> .config
    make oldnoconfig V=1

    make V=1 \
        EXTRAVERSION="-%{release}" \
        PYTHON=python3 \
%if 0%{?cross:1}
        HOSTCC=gcc \
        CC=x86_64-linux-gnu-gcc \
        AS=x86_64-linux-gnu-as \
        LD=x86_64-linux-gnu-ld \
        OBJCOPY=x86_64-linux-gnu-objcopy \
        OBJDUMP=x86_64-linux-gnu-objdump \
        STRIP=x86_64-linux-gnu-strip \
%endif
        $4

    cp out/$2 binaries/$3
}

# seabios
build_bios %{_sourcedir}/config.seabios-128k bios.bin bios.bin

# seavgabios
%global vgaconfigs bochs-display stdvga
for config in %{vgaconfigs}; do
    build_bios %{_sourcedir}/config.vga-${config} \
               vgabios.bin vgabios-${config}.bin out/vgabios.bin
done


%install
mkdir -p $RPM_BUILD_ROOT%{_datadir}/seabios
mkdir -p $RPM_BUILD_ROOT%{_datadir}/seavgabios
install -m 0644 binaries/bios.bin $RPM_BUILD_ROOT%{_datadir}/seabios/bios.bin
install -m 0644 binaries/vgabios*.bin $RPM_BUILD_ROOT%{_datadir}/seavgabios


%files
%doc COPYING COPYING.LESSER README


%files bin
%dir %{_datadir}/seabios/
%{_datadir}/seabios/bios*.bin

%files -n seavgabios-bin
%dir %{_datadir}/seavgabios/
%{_datadir}/seavgabios/vgabios*.bin


%changelog
* Mon Mar 20 2023 Gerd Hoffmann <kraxel@redhat.com> - 1.16.2-1
- Update to 1.16.2
```

手动配置编译:
```bash
make menuconfig
# General Features: 1. 禁用`Support Xen HVM`; 2. 大小 256, 128不够;
# VGA ROM: `QEMU/Bochs VBE SVGA`
make
```

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

    `when config, append --with-sdl2 --enable-debugger`
- [[SeaBIOS] physical memory read error](https://sourceforge.net/p/bochs/discussion/39592/thread/4f3d95a9/)

       可能是bios问题, 大概率是vgarom. 将romimage换成/usr/share/seabios/bios.bin后, "physical memory read error"消失, 但bochs还是黑屏. 使用自编译的seabios-1.16.2构建的bin也是黑屏.

0x0000322f3130约在800M位置, 将megs改为1024(2000年时使用megs=32没遇到过该问题), 不再报该错误, 但还是黑屏, 看不到bios界面. 自编译bochs也无法解决该问题, 推荐qemu-system-i386.

### `apt source --compile bochsbios`遇到`openjade:/usr/share/sgml/docbook/stylesheet/dsssl/modular/html/../common/dbtable.dsl:224:13:E: 2nd argument for primitive "ancestor" of wrong type: "#node-list()" not a singleton node list`
ref:
- [bochs-2.7 make编译时出错cannot generate system identifier for public text "-//OASIS//DTD DocBook V4.1//EN" ](https://www.cnblogs.com/kendoziyu/p/cannot-generate-system-identifier-for-public-text-OASIS-DTD-DocBook-V4-1-EN.html)

尝试:
1. `apt purge openjade`, 经验证, 无效: dpkg-buildpackage时报缺docbook-dsssl依赖, 而docbook-dsssl依赖openjade.
2. 注释`debian/rules` doc相关命令, 编译后install时报找不到doc

`apt source --compile bochsbios`已编出二进制在`debian/tmp`

### bochs启动报`ata0-0: could not open hard drive image file 'c.img'`
c.img存在c.img.lock, 删除c.img.lock即可