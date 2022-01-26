## ClearOS
1. 环境变量

    1. `~/.ssh`
    1. v2rayA
1. apt/dnf
1. nginx.conf/备份DB
1. 浏览器配置

    1. SwitchyOmega
    1. better-onetab/OneTab
    1. Bookmarks clean up
1. postman
1. wechat
1. git

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