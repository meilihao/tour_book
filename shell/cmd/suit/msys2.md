# msys2
MSYS2 是一个工具和库的集合, 提供了一个易于使用的环境来构建、安装和运行本机 Windows 软件. 因此msys2中的包都是用于在windows上构建exe的.

[**MSYS2 requires 64 bit Windows 8.1 or newer**](https://www.msys2.org/). msys2-x86_64-20221028.exe是最后一个支持win7的版本.

相比opensuse的[windows:/mingw:/win64](http://download.opensuse.org/repositories/windows:/mingw:/win64)存在mingw cross-toolchain即在linux上构建windows的exe.

[msys2 Environments选择](https://www.msys2.org/docs/environments/):

    - UCRT在构建时和运行时比 MSVC 的兼容性更好, 但它仅默认在 Windows 10 及以上提供

    MSYS2 和 MinGW 都提供 gcc, 但是属于两个完全不同的工具链. 前者属于 msys2-devel, 后者属于 mingw-w64-$arch-toolchain. 使用 mingw-gcc 编译的目标文件是原生的, 而使用 msys2-gcc 编译的目标文件依赖于 msys-2.0.dll 提供的虚拟 POSIX 环境.

## 添加msys2 repo
ref:
- [Repositories and Mirrors](https://www.msys2.org/docs/repos-mirrors/)
- [Repos](https://packages.msys2.org/repos)
- [Mirrors](https://www.msys2.org/dev/mirrors/)
- [MSYS2-Introduction](https://www.msys2.org/wiki/MSYS2-introduction/)

```bash
# vim /etc/pacman.conf
# always include msys!
[msys]
Include = /etc/pacman.d/mirrorlist.msys # for mingw-w64-cross-toolchain

[mingw32]
Include = /etc/pacman.d/mirrorlist.mingw32

[mingw64]
Include = /etc/pacman.d/mirrorlist.mingw64
# vim /etc/pacman.d/mirrorlist.msys
Server = https://repo.msys2.org/msys/$arch
# vim /etc/pacman.d/mirrorlist.mingw32
Server = https://repo.msys2.org/mingw/mingw32
# vim /etc/pacman.d/mirrorlist.mingw64
Server = https://repo.msys2.org/mingw/mingw64
# pacman -Syu
# pacman -S mingw-w64-x86_64-gcc
# pacman -S mingw-w64-x86_64-toolchain # target=host, 在windows上构建exe. windows安装可用[installer: msys2-x86_64-20230318.exe](https://www.msys2.org/)
# pacman -S mingw-w64-cross-toolchain # host={i686,x86_64}-pc-msys and target={i686,x86_64}-w64-mingw32, 用于在windows上构建mingw编译链, 不常用
```

> msys中包的host是{i686,x86_64}-pc-msys, mingw64中的是x86_64-w64-mingw32

## FAQ
### 添加msys2 repo后执行`pacman -Syu`报`signature from "Christoph Reiter (MSYS2 development key) <reiter.christoph@gmail.com>" is unknown trust`
```bash
# pacman-key --init
# pacman-key --populate msys2
# wget https://repo.msys2.org/msys/x86_64/msys2-keyring-1~20230316-1-any.pkg.tar.zst
# pacman -U --config <(echo) msys2-keyring-1_20230316-1-any.pkg.tar.zst
# pacman-key --list-keys
```