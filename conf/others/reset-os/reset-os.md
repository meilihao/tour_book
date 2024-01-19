## ClearOS
1. 环境变量

    1. `~/.ssh`
    1. v2rayA
1. apt/dnf
1. nginx.conf/备份DB
1. 浏览器配置

    1. 书签
    1. SwitchyOmega
    1. better-onetab/OneTab
    1. Bookmarks clean up
1. postman
1. wechat
1. git

## other
- [Ubuntu上安装番茄时钟](https://zhuanlan.zhihu.com/p/350023097)

    `pip3 install tomato-clock`
- [alarm clock](https://alarm-clock-applet.github.io/)

## Dell PowerEdge T630
安装系统选basic graphics mode/字符安装, 估计与其内置显卡次(显存32M, 分辨率最高`1024*768`)有关.

os:
- debian 12

    - kde wayland可运行, 需手动执行startplasma-wayland, ok; startplasma-x11, failed
    - gnome, 启动桌面失败
- fedora worktation 39: basic graphics mode, ok

## repo
```
# from https://www.atzlinux.com/faq.htm#apt-hand-other from install v2raylui
# wget -qO - https://download.sublimetext.com/sublimehq-pub.gpg | sudo apt-key add -
# echo "deb https://apt.atzlinux.com/atzlinux buster main contrib non-free" sudo tee -a /etc/apt/sources.list
```

## softwares
### deepin
- 向日葵

    使用`sudo apt install com.oray.sunlogin.client`, 直接使用`dpkg -i <官方>.deb`会报缺依赖, 且该依赖不再apt repo里.
### llvm
```bash
sudo apt install clang-13 lldb-13 lld-13  llvm-13 llvm-13-dev
```

### 清理/var/spool/postfix/maildrop
- 清理全部: `sudo postsuper -d ALL`
- 清理其中的某个: `postcat -q 165EA1AF2`