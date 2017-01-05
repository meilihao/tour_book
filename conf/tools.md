# download/下载

## uget+aria2

[uget](http://ugetdm.com)是一款轻量级的自由开源的下载管理器类似迅雷，可运行Linux、windows和MAC系统上,支持队列下载和恢复下载和通过终端下载的功能.uget采用aria2作为后端，安装aria2插件后可与其进行交互。

`aria2` 是 Linux CLI 界面下的多线程下载工具，与 axel 类似，但比之更强大.它支持 HTTP/HTTPS, FTP, BitTorrent 和 Metalink 协议，支持多线程下断点续传.另外，这里有一个名为 aria2fe 的 aria2 前端 GUI 程序，直接执行里面编译好的二进制程序就可使用.

### 以centos 7 x64为例

1. 下载Fedora分类下的最新版[uget](http://sourceforge.net/projects/urlget/files/uget%20%28stable%29/1.10.4/uget-1.10.4-1.fc21.x86_64.rpm/download)和[aria2](http://sourceforge.net/projects/urlget/files/aria2-plugin/1.18.x/aria2-1.18.2-1.fc21.x86_64.rpm/download)
2. 安装uget和aria2
3. 配置
	```shell
    1. 在uget的`编辑`-`设置`-`插件`选项卡中勾选`启用 aria2 插件`
    2. uget主界面左侧的`分类`窗口-`Home`分类-右键`属性`-`新下载的默认设置1`选项卡-`每台服务器连接数`改为'16'(其它数字也可).
    ```

# 美化

## 主题

[一款很棒的GTK桌面主题：Arc](https://linux.cn/article-5614-1.html)

###  gnome-tweak-tool

gnome-tweak-tool是一款个性化GNOME 3 的图形界面工具,可用于切换gnome-shell主题等.

### Unity Tweak Tool

Unity Tweak Tool 是一个专为Ubuntu Unity 桌面环境的配置工具，允许你对Unity 环境进行各项配置,功能和gnome-tweak-tool类似.

## 爱壁纸/LoveWallpaper

[LoveWallpaper](http://www.lovebizhi.com)是专业的桌面高清壁纸软件，提供万款优质高清壁纸，具有试试手气、按颜色筛选壁纸以及定时切换壁纸功能等功能，充分满足了壁纸达人需求.

### 以centos 7 x64为例

1. 下载Fedora分类下的最新版[LoveWallpaper](http://www.lovebizhi.com/linux_fedora)
2. 安装LoveWallpaper4LinuxFedora.rpm即可

#### FAQ

1 . 依赖`python-pyside`(其实python-pyside依赖qt4)

>方法1(推荐):在`http://www.rpmfind.net/linux/RPM/index.html`或`http://rpm.pbone.net/`搜索`python-pyside`,根据fedora相应路径规划的实际情况下载[python-pyside-1.1.0-4.fc20.x86_64.rpm](http://dl.fedoraproject.org/pub/fedora/linux/releases/20/Everything/x86_64/os/Packages/p/python-pyside-1.1.0-4.fc20.x86_64.rpm)(python-pyside-1.1.0需qt4(x86-64) >= 4.8.5,执行`sudo yum install qt`即可)
>方法2: 下载[python-pyside-1.2.2-2.fc22.src.rpm](http://dl.fedoraproject.org/pub/fedora/linux/development/22/source/SRPMS/p/python-pyside-1.2.2-2.fc22.src.rpm),再用`rpmbuild`处理并安装(python-pyside-1.2.2需qt4(x86-64) >= 4.8.6,请自行解决qt版本)

2 . `python-pyside`依赖`libphonon.so.4()(64bit)`

>在`http://www.rpmfind.net/linux/RPM/index.html`搜索`libphonon.so.4`,根据结果执行`sudo yum install phonon`即可.

3 . `python-pyside`依赖`libQtWebKit.so.4()(64bit)`

>在`http://www.rpmfind.net/linux/RPM/index.html`搜索`libQtWebKit.so.4`,根据结果和对照`http://dl.fedoraproject.org/pub/fedora/linux/releases/20/Everything/x86_64/os/Packages/q`内容,执行`sudo yum install qtwebkit`即可.

4 . `python-pyside`依赖`libshiboken-python2.7.so.1.1()(64bit)`

> 查找方法同上,找到[shiboken-libs-1.1.0-3.fc19.x86_64.rpm ](http://dl.fedoraproject.org/pub/fedora/linux/releases/20/Everything/x86_64/os/Packages/s/shiboken-libs-1.1.0-3.fc19.x86_64.rpm)安装即可.

5 . `from PySide.QtGui import * ImportError: /usr/lib64/python2.7/site-packages/PySide/QtGui.so: undefined symbol`

> 前几天安装liteide时设置了LD_LIBRARY_PATH,现在qt的so库包括了liteide自带的qt动态库和为解决python-pyside依赖而安装的qt库,两者冲突,所以还原LD_LIBRARY_PATH再重启系统即可.

# 梯子

- [shadowsocks-libev](github.com/shadowsocks/shadowsocks-libev)
- [goagent](https://github.com/goagent/goagent),可行
- [greatagent](https://github.com/greatagent/greatagent/wiki),未测试
- [lantern](https://getlantern.org/),未测试

## shadowsocks

- [linux-ubuntu使用shadowsocks客户端配置](https://aitanlu.com/ubuntu-shadowsocks-ke-hu-duan-pei-zhi.html)

# 翻译

## 有道词典

- `ImportError: No module named 'PyQt5.QtWebKitWidgets'`

	  sudo apt-get install python3-pyqt5.qtwebkit

- `ImportError: cannot import name 'QtQuick'`

	  sudo apt-get install python3-pyqt5.qtquick

- `module "QtQuick.Controls" is not installed`

	  sudo apt-get install qml-module-qtquick-controls

# 终端/terminal

## terminix,推荐
- [terminix](https://github.com/gnunn1/terminix)

## tmux

[tmux快捷键](http://blog.csdn.net/hcx25909/article/details/7602935)

- `configure: error: "libevent not found"`

    sudo apt-get install libevent-dev

- `configure: error: "curses not found"`

    sudo apt-get install ncurses-dev
    sudo dnf install ncurses-devel

- 快捷键失效

 按键顺序是`ctrl+b`松开后再按其他键.例如`ctrl+b ？`，应该先同时按`ctrl+b` 松开后，`shift+/（即输入？）`

## Terminator

安装:

    sudo apt-get install terminator

推荐配置：

－　右键-> preferences -> `Profiles` tab -> `Colors` tab -> "Foreground and Background # Build-in schemes"="Ambience" && "Palette # Built-in schemes"="Ambience"

## Bash shell对话框

zenity和whiptail

[如何在Bash Shell脚本中显示对话框](https://linux.cn/article-5558-1.html)

## 通信

### RTX2013

1. [安装Wine1.7](https://www.winehq.org/download/ubuntu)
2. 把wine的环境设置为32位的（实践证明，不这样做的话，安装的RTX是没法用的）

       export WINEARCH="win32"
       sudo rm -rf ~/.wine

3. 安装依赖软件

       winetricks msxml3 msxml6 riched20 riched30 ie6 vcrun6 vcrun2005sp1

 如果64位系统，安装msxml6时，要下载64位的，点击[这里](http://www.microsoft.com/en-us/download/details.aspx?id=3988)

 64位的ubuntu需要把ie6改成ie8(我没改,还是用ie6).

 如果有些软件如果下载不了，可以根据提示使用浏览器下载后，放到～/.cache/winetricks相应目录下，例如vcrun6就放在～/.cache/winetricks/vcrun6/下,再运行winetricks进行安装即可.

4. 打开rtx2013的安装软件安装，这过程会遇到一个组件未注册，忽略即可.安装成功后，就可以使用rtx了.rtx这时还有个问题，就是rtx会一直停留在离开状态，不启用`文件-个人设置-回复设置-自动状态切换`即可.

5. 输入框内容或其他地方的中文乱码

       从windows下拷贝simsun.ttc字体到~/.wine/driver_c/windows/Fonts下

 修改注册表,将如下内容写入rtx.reg:

 ```
REGEDIT4
[HKEY_LOCAL_MACHINE\Software\Microsoft\Windows NT\CurrentVersion\FontSubstitutes]
"Arial"="simsun"
"Arial CE,238"="simsun"
"Arial CYR,204"="simsun"
"Arial Greek,161"="simsun"
"Arial TUR,162"="simsun"
"Courier New"="simsun"
"Courier New CE,238"="simsun"
"Courier New CYR,204"="simsun"
"Courier New Greek,161"="simsun"
"Courier New TUR,162"="simsun"
"FixedSys"="simsun"
"Helv"="simsun"
"Helvetica"="simsun"
"MS Sans Serif"="simsun"
"MS Shell Dlg"="simsun"
"MS Shell Dlg 2"="simsun"
"System"="simsun"
"Tahoma"="simsun"
"Times"="simsun"
"Times New Roman CE,238"="simsun"
"Times New Roman CYR,204"="simsun"
"Times New Roman Greek,161"="simsun"
"Times New Roman TUR,162"="simsun"
"Tms Rmn"="simsun"
```

 最后执行

       $ regedit rtx.reg

# 其他

### flash

    sudo dnf install flash-plugin

前提:设置adobe的repo,即到[官网](https://get.adobe.com/flashplayer/?loc=cn)下载yum文件(rpm文件)再安装即可.

### 字符浏览器

- w3m
	- `w3m -I GBK  www.baidu.com`,解决中文乱码
- lynx
