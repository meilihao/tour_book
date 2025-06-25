# tigervnc
ref:
- [Configuring Tigervnc](https://www.linuxfromscratch.org/blfs/view/systemd/xsoft/tigervnc.html)
- [Rocky Linux 10.0 Available Now](https://rockylinux.org/zh-CN/news/rocky-linux-10-0-ga-release)

    远程桌面协议（RDP）取代了VNC, 现在是图形远程访问的默认协议

## 部署
```bash
# apt install ./tigervncserver_1.13.1-1ubuntu1_amd64.deb
# bash -c "cat >> /etc/tigervnc/vncserver.users" <<'EOF'
#第一个桌面即5901端口对应ubuntu用户
#第二个桌面即5902端口对应root用户
:1=ubuntu
:2=root
EOF
# bash -c "cat >> /etc/tigervnc/vncserver-config-defaults" <<'EOF'
# session必选,必须与当前桌面匹配,可以在/usr/share/xsessions目录的文件中查看当前桌面
session=gnome
geometry=1920x1024
EOF
# vncpasswd
# systemctl start vncserver@:1 # 1为vncserver.users对应的桌面
# systemctl enable vncserver@:1
# systemctl restart vncserver@:1
```

建议配合realvnc的client使用.

> geometry修改后可能需要重启并重新设置分辨率