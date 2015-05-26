# markdown

## haroopad

[haroopad](http://pad.haroopress.com/user.html)是一款覆盖三大主流桌面系统的编辑器，支持 Windows、Mac OS X 和 Linux。 主题样式丰富，语法标亮支持 54 种编程语言。

### 以centos 7 x64为例

1. 下载名为[Linux Binary (64bit)]()的版本
2. 安装haroopad

	```shell
	tar -zxvf haroopad-v0.13.1-x64.tar.gz
	tar -zxvf data.tar.gz  # 程序主体,相关路径都已规划好
	sudo cp -r ./usr /

	tar zxf control.tar.gz
	chmod 755 postinst
	sudo ./postinst # 处理依赖库
	```
3. 创建图标
	```shell
    # 设置icon
    # Icon=/usr/share/icons/hicolor/128x128/apps/haroopad.png
    sudo vi /usr/share/applications/Haroopad.desktop
    ```
    
#### FAQ

1. `error while loading shared libraries: libudev.so.0: cannot open shared object file: No such file or directory`

>运行命令`locate udev.so`,就会发现已存在`/usr/lib64/libudev.so.1`只是名称略有差异(没有时请`sudo yum install systemd-libs`),执行`ln -sf /usr/lib64/libudev.so.1 /usr/lib64/libudev.so.0`即可(其实control.tar.gz中的postinst脚本就是处理这个问题的).

# coding

## liteide

[liteide](https://github.com/visualfc/liteide)是一款简单、开源、跨平台的 Go 语言 IDE.

### 以centos 7 x64为例

1. 下载名为[liteidex27.1.linux-64.tar.bz2](http://sourceforge.net/projects/liteide/files/X27.1/)的版本;如果已安装了qt4.8,那下载[liteidex27.1.linux-64-system-qt4.8.tar.bz2]()即可[**推荐**]
2. 安装liteide

	```shell
	tar -jxvf liteidex27.1.linux-64.tar.bz2
	mv ./liteide ~/opt
    #再运行`~/opt/liteide/bin/liteide`即可
	```
    
#### FAQ

1 . `error while loading shared libraries: libQtGui.so.4: cannot open shared object file: No such file or directory`

>`export LD_LIBRARY_PATH=~/opt/liteide/lib/liteide:$LD_LIBRARY_PATH`

2 . `error while loading shared libraries: libpng12.so.0: cannot open shared object file: No such file or directory`

>`sudo yum install libpng12`,缺少so时可在[pkgs.org](http://pkgs.org)查询

3 . 安装qt4.8

>执行`sudo yum install qt`即可

# editor

## gedit

#### FAQ

- 关闭自动创建备份档(即*~的文件)

> `编辑`菜单－`首选项`－`编辑器`选项卡-`文件保存`节-"在保存前创建备份文件",去钩.

## vim

#### FAQ

- vim配置文件位置

> 进入vim,输入命令`:version`,找到vimrc相关内容即是.

- 关闭自动创建备份档(即*~的文件)

> 检查vimrc是否存在(用户配置文件为~/.vimrc，相关的文件位于~/.vim/；全局配置文件为/etc/vimrc，相关的文件位于/usr/share/vim/),不存在时,`cp /usr/share/vim/vim74/vimrc_example.vim ~/.vimrc`
> 打开配置文件，找到这一句：if has("vms"),将这个判断里的if部分保留，else部分注释（Vim的注释符是"）掉`set backup`即可.