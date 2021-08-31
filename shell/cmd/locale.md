# locale
参考:
- [glibc locale使用简析](https://openeuler.org/zh/blog/wangshuo/glibc%20locale%E4%BD%BF%E7%94%A8%E7%AE%80%E6%9E%90/glibc+locale%E4%BD%BF%E7%94%A8%E7%AE%80%E6%9E%90.html)

软件在运行时的语言环境.

> localedef : 转化语言环境和字符集源文件以生成语言环境数据库, 由它创建的语言环境对象代码可被`setlocale()`使用. `setlocale()`由glibc实现, 具体会加载哪些locale要看其实现.

> **locale-gen是localedef的脚本封装**. 用户只要去 /etc/locale.gen 文件里, 把想要生成的 locale 取消注释, 然后运行 locale-gen, locale-gen 就会去检查 /etc/locale.gen 里没有被注释的行然后调用 localedef 并传以相应的参数生成相应的 locale.

localedef 是用来生成 locale 的工具. localedef 能够读取 **locale 定义文件**, 以及 charmap 文件, 然后编译它们, 将编译好的二进制的 locale 数据库放在 /usr/lib/locale 目录里 (这个目录不是绝对的, 可以通过 localedef --help 查看编译好的 locale 数据库文件放在哪里), 然后 Linux 系统中一些会使用到 locale 的程序 (基本上是使用 GNU libc 库的程序) 就能够找到并使用这些 locale.

> localedef 可以将编译好的每一个 locale 在 /usr/lib/locale 目录下单独创建一个子目录, 也可以将编译好的 locale 附加到 /usr/lib/locale/locale-archive 文件里 (通过 --add-to-archive 选项控制).

生成了相应的 locale 数据库文件之后, 就可以设置 LC_* 环境变量为那个生成的 locale (一般是直接设置 LANG 变量). 然后, 支持 i18n 的程序就会根据这些环境变量来向用户展示信息, GNOME, KDE 开发的程序一般都支持 i18n, 因为 GTK, Qt 早已经支持 i18n 了.

locale definition file, 即locale 定义文件, 每一个文件会包含某个 locale 的所有信息 (LC_CTYPE 信息, LC_MESSAGE 信息等等), 这文件包含的也是 localedef 命令生成这个 locale 所需要的信息.

locale 定义文件一般位于 /usr/share/i18n/locales 目录中, 不用担心, localedef 自己会找到它们, 可以通过 localedef --help 来查看与 locale 相关的目录.

编译 locale 还需要 charmap 文件, 这种文件一般位于 /usr/share/i18n/charmaps 目录下, 这也是可以通过 localedef --help 命令知晓的.

相关系统文件：
- 在文件/usr/share/i18n/SUPPORTED : 列出了当前系统支持的所有locale与字符集的名字, 是`/usr/share/i18n/locales`和`/usr/share/i18n/charmaps`的组合结果
- 在目录/var/lib/locales/supported.d/下，列出了当前系统已经生成的所有locale的名字
- 在目录/usr/lib/locale/<locale_name>/LC_*，用locale-gen编译出的locale文件
- 在文件/usr/lib/locale/locale-archive中，包含了很多本地已经生成的locale的具体内容，因此这个文件往往很大. 使用命令localedef管理这一文件. 使用locale-gen命令编译出来的locale内容默认写入该文件中
- 在文件/etc/default/locale中，可以手动配置locale环境变量，LC_CTYPE之类
- 在目录/usr/share/i18n/charmaps下，缺省的charmap存放路径
- 在目录/usr/share/i18n/locales下，缺省的locale source file存放路径

相关系统命令：
- locale 列出当前采用的各项本地策略，这些由LC_*环境变量定义

    **en_US.utf8可显示中文(因为它也是utf8, 语言显示策略是en_US)**
- locale charmap 列出系统当前使用的字符集
- locale -a 列出系统中已经安装的所有locale
- locale -m 列出系统中已经安装的所有charmap
- locale-gen --purge 将/usr/lib/locale/里面的locale文件删掉
- localedef --list-archive: 列出文件/usr/lib/locale/locale-archive中所有可用的locale的名字

## example
```bash
$ locale # 查看当前系统的语言环境
LANG=zh_CN.UTF-8
LANGUAGE=zh_CN
LC_CTYPE=zh_CN.UTF-8
LC_NUMERIC="zh_CN.UTF-8"
LC_TIME="zh_CN.UTF-8"
LC_COLLATE="zh_CN.UTF-8"
LC_MONETARY="zh_CN.UTF-8"
LC_MESSAGES="zh_CN.UTF-8"
LC_PAPER="zh_CN.UTF-8"
LC_NAME="zh_CN.UTF-8"
LC_ADDRESS="zh_CN.UTF-8"
LC_TELEPHONE="zh_CN.UTF-8"
LC_MEASUREMENT="zh_CN.UTF-8"
LC_IDENTIFICATION="zh_CN.UTF-8"
LC_ALL=
$ locale -m  # 查看系统支持的所有可用的字符集编码, 所有的字符集都放在 /usr/share/i18n/charmaps
$ locale -a  # 查看系统中所有已配置的区域格式(= `localedef --list-archive` + `C,C.UTF-8,POSIX`), 在`/usr/share/locale`下， 以 zh_CN.UTF-8 locale(`/usr/share/locale/zh_CN.UTF-8`) 为例，该目录中就包含了 LC_MESSAGES
$ sudo localedef -f UTF-8 -i zh_CN zh_XX.UTF-8 # 默认保存在/usr/lib/locale/locale-archive. 当locale_name未以`.UTF-8`作为结尾时`locale -a`会显示`zh_XX`和`zh_XX.utf8`两个locale???.
$ sudo localedef -v -f UTF-8 -i zh_CN /usr/lib/locale/<locale_name> # 根据模板创建locale并放入/usr/lib/locale/<locale_name>, 比如locale_name=zh_XX.UTF-8, locale_name同时会出现在`locale -a`中. 当locale_name未以`.UTF-8`作为结尾时仅创建`zh_XX`.
$ locale -a|grep <name> # 查看已创建的locale
$ env LC_ALL=zh_XX.UTF-8,LANGUAGE=  strace ls # 查看效果
$ sudo localedef --delete-from-archive <name> <name>.utf8 # 清理/usr/lib/locale/locale-archive, name可从`localedef --list-archive`获取
$ sudo rm -rf /usr/lib/locale/xxx.UTF-8 # 清理后`locale -a`不再展示xxx.UTF-8
$  env LANG=C sudo localedef -v --delete-from-archive zh_BB # zh_BB被自定义放在了/usr/lib/locale下, 直接删除/usr/lib/locale/zh_BB即可
locale "zh_BB" not in archive
$ convmv -f iso-8859-1 -t utf-8 <filename> # 转换文件名的编码
$ iconv -f iso-8859-1 -t utf-8 filename > newfile # 转换文件内容的编码方式即转码
```

## 进阶
Linux的glibc中的`setlocale()`:
```c
#include <locale.h>
char* setlocale(int category, const char* locale);
```

category：为locale分类，表达一种locale的领域方面，通常有下面这些预定义常量：LC_ALL、LC_COLLATE、LC_CTYPE、LC_MESSAGES、LC_MONETARY、LC_NUMERIC、LC_TIME，其中 LC_ALL 表示所有其它locale分类的并集.

locale：为期望设定的locale名称字符串，在Linux/Unix环境下，通常以下面格式表示locale名称：`language[_territory][.codeset][@modifier]` from `man setlocale`，language 为 ISO 639 中规定的语言代码，territory 为 ISO 3166 中规定的国家/地区代码(`定义文件放在/usr/share/i18n/locales`)，codeset 为字符集名称(`/usr/share/i18n/charmaps`), modifier为某些 locale 变体的修正值. 它们组合之后的locale信息放在/usr/lib/locale/locale-archive里，是二进制文件. 通过格式所知，locale总是和一定的字符集相联系的, 比如:
1. 我说中文，身处中华人民共和国，使用国标18030字符集来表达字符

    zh_CN.GB18030＝中文_中华人民共和国＋国标18030字符集
1. 我说英语，身处英国，使用ISO-8859-1字符集来表达字符

    en_GB.ISO-8859-1=英语_英国.ISO-8859-1字符集
1. 我说德语，身处德国，使用UTF-8字符集，习惯了欧洲风格

    de_DE.UTF-8@euro＝德语_德国.UTF-8字符集@按照欧洲习惯加以修正


当 locale 为 NULL 时，函数只做取当前 locale 的操作，并通过返回值传回，并不改变当前 locale.
当 locale 为 "" 时，根据环境的设置来设定 locale，检测顺序是：环境变量 LC_ALL，每个单独的locale分类LC_*，最后是 LANG 变量. 为了使程序可以根据环境来改变活动 locale，一般都在程序的初始化阶段加入代码`setlocale(LC_ALL, "")`.

当C语言程序初始化时（刚进入到 main() 时），locale 被初始化为默认的 C locale，其采用的字符编码是所有本地 ANSI 字符集编码的公共部分，是用来书写C语言源程序的最小字符集（所以才起locale名叫：C）.
当用 setlocale() 设置活动 locale 时，如果成功，会返回当前活动 locale 的全名称；如果失败，会返回 NULL. 

locale 是一组 C 程式语言处理自然语言(文字)的显示方式， 也可以简单的说，locale 就是一组地区性语言的显示格式. 由国家语言和各地习俗影响所决定的惯例，或代表一个地理区域的定义所组成，这些惯例包含文字、日期、数字、货币格式和排序等等. 这代表着 locale 可让程式的输出可以直接反应地方区域性的文化.

C 语言的 locale按照将文化传统的各个方面分成12个大类, 这12个大类分别是分为下列各大类:
- [LANGUAGE](https://www.gnu.org/software/gettext/manual/gettext.html#The-LANGUAGE-variable) : 指定个人对语言环境值的主次偏好，比如`zh_CN:en_GB:en`, 作用是如果前面的 locale 缺少翻译，会自动使用后面的 locale 显示界面, 因此它会先使用简体中文,没有翻译时再使用英文.
- LC_ALL : 指定所有的 Locale
- LC_CTYPE : 字元定义即语言符号及其分类 (包含字元分类与转换规则, 比如大写字母，小 写字母，大小写转换，标点符号、可打印字符和其他的字符属性等方面)
- LC_MESSAGES : 信息显示, 包括提示信息, 错误信息, 状态信息, 标题, 标签, 按钮和菜单等
- LC_TIME : 时间显示格式
- LC_NUMERIC : 数字格式
- LC_MONETARY : 货币单位
- LC_COLLATE : 字母顺序与特殊字元比较即**比较和排序习惯**. 比如, mysql会在`order by`中利用其排序规则.
- LC_NAME : 姓名书写方式
- LC_ADDRESS : 地址书写方式
- LC_TELEPHONE : 电话号码书写方式
- LC_MEASUREMENT : 度量衡表达方式
- LC_PAPER : 默认纸张尺寸大小
- LC_IDENTIFICATION : 对locale自身包含信息的概述
- LANG : 语言显示, **用于设定LC_*的默认值**, 因此除非LC_*有指定, 否则都会用LANG

> 优先级的关系： LANGUAGE > LC_ALL>LC_*>LANG, LC_ALL通常为空.

其中与一般使用者息息相关的，是字元定义 (LC_CTYPE) 与语言显示 (LANG). LC_CTYPE 直接关系到某些字元或内码(code point)在目前的 locale 下是否可列印？要如何转换字码？对应到哪一个字？.... 等等.LANG 则关系到软体的信息输出的地域格式.当 LC_MESSAGES、LC_TIME、LC_NUMERIC、 LC_MONETARY 等没有设定的时候，会直接取用 LANG 的环境设定值.

当一个程式启动时，系统会预设给它一个初始 locale，称为 POSIX 或 C locale. 在此 locale 下，程式的表现会与传统的 C 语言中一样， 使用英文做信息输出，只能处理英文等 ASCII 码等等. 如果该程序需要I18N，则它在启动后就会马上调用系统函式来改变它的 locale， 如此它就摇身一变，变成可以处理该 locale 所代表的地区语言了.

## FAQ
### `localedef -v -f UTF-8 -i zh_CN /usr/lib/locale/zh_CN.UTF-8`报`non-symbolic character value should not be used`
与操作系统的glibc有关, 没有其相应的语言包, 比如`/usr/share/locale/zh_CN/LC_MESSAGES/libc.mo`.

```bash
apt-file search libc.mo
libc-l10n: /usr/share/locale/zh_CN/LC_MESSAGES/libc.mo
```

> ubuntu 14.04返回的是`language-pack-zh-hans-base: /usr/share/locale-langpack/zh_CN/LC_MESSAGES/libc.mo`

也可从其他机器的拷贝/usr/share/locale/zh_CN, 再将其LC_MESSAGES清空并替换为libc-l10n包里的LC_MESSAGES即可.

### `locale`输出
：

1、(LC_CTYPE)
6、信息主要是提示信息,错误信息,状态信息,标题,标签,按钮和菜单等(LC_MESSAGES)