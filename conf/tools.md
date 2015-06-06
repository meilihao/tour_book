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

- [goagent](https://github.com/goagent/goagent),可行
- [greatagent](https://github.com/greatagent/greatagent/wiki),未测试
- [lantern](https://getlantern.org/),未测试

# 翻译

## 有道词典

- `ImportError: No module named 'PyQt5.QtWebKitWidgets'`

	  sudo apt-get install python3-pyqt5.qtwebkit

- `ImportError: cannot import name 'QtQuick'`

	  sudo apt-get install python3-pyqt5.qtquick

- `module "QtQuick.Controls" is not installed`

	  sudo apt-get install qml-module-qtquick-controls

# 终端/terminal

## tmux

- `configure: error: "libevent not found"`

    sudo apt-get install libevent-dev

- `configure: error: "curses not found"`

    sudo apt-get install ncurses-dev
    sudo dnf install ncurses-devel