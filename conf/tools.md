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