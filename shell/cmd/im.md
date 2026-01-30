# im

## example
```
im-config -l
fcitx-remote
fcitx5-configtool # 配置后如果没问题即可使用fcitx5
```

fcitx配置:
```bash
$ vim ~/.bashrc
export GTK_IM_MODULE=fcitx
export QT_IM_MODULE=fcitx
export XMODIFIERS=@im=fcitx
```