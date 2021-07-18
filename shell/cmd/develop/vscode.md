# vscode
## FAQ
### 插件log
Help -> `Toggle Developer Tools`

### System limit for number of file watchers reached
打开vue项目时报该错, 解决方法:
```conf
sudo vim /etc/sysctl.conf
fs.inotify.max_user_watches=524288
sudo sysctl -p
```