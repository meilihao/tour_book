## gogs

### FAQ

#### `ssh: connect to host $host_name port 22: Connection refused`

git clone时碰到,原因是ssh服务未运行.

    sudo apt-get install openssh-server
    sudo systemctl start ssh

#### systemd启动

按照[官方文档](http://gogs.io/docs/intro/faqs.html)编辑systemd服务模板文件(**其中的User,Group,Environment项推荐使用默认的git账号**),再将文件保存到`/etc/systemd/system`后即可用`sudo systemctl start gogs`启动.