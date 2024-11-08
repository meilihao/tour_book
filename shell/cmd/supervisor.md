# supervisor
不推荐使用, 多次出现其管理的进程出现诡异问题. **建议用systemd service代替**

## 安装
```
$ sudo apt install supervisor # 推荐
$ pip install supervisor # from http://supervisord.org/installing.html
$ sudo systemctl enable supervisor
# 向/etc/supervisor/conf.d加入配置
$ supervisorctl reload
```

## 配置
```conf
# 放入/etc/supervisor/conf.d
[inet_http_server]
port=0.0.0.0:9001
username=admin
password=123456

[program:micro]
directory=/app
command=/app/micro
startsecs=0 ; 启动n秒后没有异常退出，就表示进程正常启动了，默认为1秒
stopwaitsecs=0
autostart=true ; 在supervisord启动的时候也自动启动
autorestart=true ; 程序退出后自动重启,可选值：[unexpected,true,false]，默认为unexpected，表示进程意外杀死后才重启
stdout_logfile=/tmp/micro.log ; 需要注意当指定目录不存在时无法正常启动，所以需要先手动创建目录
stderr_logfile=/tmp/micro.err
```