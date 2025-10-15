# ubuntu
## FAQ
### ubuntu输入正确用户密码登陆时重新跳转到登陆界面即无法登陆

原因：用户home目录下的.Xauthority文件拥有者变成了root，从而以用户登陆的时候无法都取.Xauthority文件.

解决：删除home目录下的.Xauthority文件，再重启(chown修改文件属性不可行).