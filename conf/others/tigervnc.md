# tigervnc
ref:
- [Configuring Tigervnc](https://www.linuxfromscratch.org/blfs/view/systemd/xsoft/tigervnc.html)

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
geometry=1024x768
EOF
# vncpasswd
# systemctl start vncserver@:1 # 1为vncserver.users对应的桌面
# systemctl enable vncserver@:1
```

建议配合realvnc的client使用