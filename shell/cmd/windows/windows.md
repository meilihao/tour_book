# windows

## FAQ
### 安装驱动
[`pnputil -i -a *.inf`, 需管理员权限](https://help.aliyun.com/document_detail/217543.html#section-1kb-hov-812)

### 查看dll依赖
- [lucasg/Dependencies](https://zhuanlan.zhihu.com/p/395557318)

	需要.net framework>=4.6.2
- `git for windows`里的ldd
### windows版本定义
ref:
- [Update WINVER and _WIN32_WINNT](https://learn.microsoft.com/en-us/cpp/porting/modifying-winver-and-win32-winnt?view=msvc-170)

`x86_64-w64-mingw32-g++ -DHAVE_MINGW -DHAVE_WIN32 -DMINGW64 -DWIN32_VSS -D_WIN32_WINNT=0x500 ...`

或
```c++
#include "stdio.h"
#include "stdint.h"
#define _WIN32_WINNT 0x0600
#include "windows.h"
#include "ws2tcpip.h"

void main(void){
	...
}
```

### 'inet_ntop' was not declared in this scope/'inet_pton' was not declared in this scope
env: [mingw64](http://download.opensuse.org/repositories/windows:/mingw:/win64/openSUSE_Leap_15.4)

source:
```c++
# cat p.c 
#include "stdio.h"
#include "stdint.h"
#include "windows.h"
#include "ws2tcpip.h"

int main(void)
{
  struct sockaddr_in sa;
  char str[INET_ADDRSTRLEN];
  inet_pton(AF_INET, "192.0.2.33", (char *)(&(sa.sin_addr)));
  inet_ntop(AF_INET, &(sa.sin_addr), str, INET_ADDRSTRLEN);
  printf("%s\n", str); // prints "192.0.2.33"
  return 0;
}
# /usr/bin/x86_64-w64-mingw32-g++ -DHAVE_MINGW -DHAVE_WIN32 -DMINGW64 -D_WIN32_WINNT=0x600 -lws2_32 p.c
/usr/lib64/gcc/x86_64-w64-mingw32/12.2.0/../../../../x86_64-w64-mingw32/bin/ld: /tmp/cc8TtleU.o:p.c:(.text+0x2a): undefined reference to `__imp_inet_pton' # 未找到原因, `nm /usr/x86_64-w64-mingw32/sys-root/mingw/lib/libws2_32.a`里存在__imp_inet_pton的符号
/usr/lib64/gcc/x86_64-w64-mingw32/12.2.0/../../../../x86_64-w64-mingw32/bin/ld: /tmp/cc8TtleU.o:p.c:(.text+0x50): undefined reference to `__imp_inet_ntop'
# rpm -qf /usr/x86_64-w64-mingw32/sys-root/mingw/lib/libws2_32.a
mingw64-runtime-10.0.0-lp154.11.4.noarch
# rpm -qf /usr/bin/x86_64-w64-mingw32-g++ 
mingw64-cross-gcc-c++-12.2.0-lp154.34.7.x86_64
```

根据/usr/x86_64-w64-mingw32/sys-root/mingw/include/ws2tcpip.h的定义:
```
#if (_WIN32_WINNT >= 0x0600)
...
#define InetNtopA inet_ntop

WINSOCK_API_LINKAGE LPCWSTR WSAAPI InetNtopW(INT Family, LPCVOID pAddr, LPWSTR pStringBuf, size_t StringBufSIze);
WINSOCK_API_LINKAGE LPCSTR WSAAPI InetNtopA(INT Family, LPCVOID pAddr, LPSTR pStringBuf, size_t StringBufSize);

#define InetNtop __MINGW_NAME_AW(InetNtop)

#define InetPtonA inet_pton

WINSOCK_API_LINKAGE INT WSAAPI InetPtonW(INT Family, LPCWSTR pStringBuf, PVOID pAddr);
WINSOCK_API_LINKAGE INT WSAAPI InetPtonA(INT Family, LPCSTR pStringBuf, PVOID pAddr);

#define InetPton __MINGW_NAME_AW(InetPton)

#endif /*(_WIN32_WINNT >= 0x0600)*/
```

inet_ntop和inet_pton只在`_WIN32_WINNT >= 0x0600`的环境. InetPtonA是asii版, 而InetPtonW是unicode版.