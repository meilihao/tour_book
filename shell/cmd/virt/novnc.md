# novnc
noVNC是一个HTML5 VNC客户端，采用HTML5 websockets、Canvas和JavaScript实现，noVNC被普遍应用于各大云计算、虚拟机控制面板中，比如OpenStack Dashboard 和 OpenNebula Sunstone 都用的是 noVNC.

目前大多数 VNC 服务器都不支持 WebSockets, 所以 noVNC 是不能直接连接 VNC 服务器的, 需要一个代理来做 WebSockets 和 TCP sockets 之间的转换. 这个代理在 noVNC 的目录里，叫做 websockify.

> 在 Github 上 noVNC 和 websockify 本来就是独立的两个项目.

## 部署
ref:
- [All in Web | 基于web的远程桌面-noVNC](https://zhuanlan.zhihu.com/p/427144657)

### 直连
只能连接一个vncserver

```
websockify --web=/usr/share/novnc 6080 localhost:5901
```

### token
ref:
- [基于令牌的目标选择](https://github.com/novnc/websockify/wiki/Token-based-target-selection)

允许连接多台vncserver

```bash
websockify --web=/usr/share/novnc 6080 --token-plugin TokenFile --token-source /usr/share/novnc/conf
websockify --web=/usr/share/novnc 6080 --target-config=/usr/share/novnc/conf # 同上, 是上面的旧语法
```

`/usr/share/novnc/conf`目录下的每个文件写入vm vnc信息, 格式是`token1: host1:port1`(token1为自定义token), 因此创建vm时需要提前配置port, 之后访问`http://192.168.1.130:8787/vnc.html?path=websockify/?token=token1`即可. 该方法进入页面需要点击一次右上角的`connect`, 也可使用`http://192.168.1.130:8787/vnc_auto.html?path=/conf?token=token1`直接看到vm屏幕.

> oracle linux 7.9上的novnc只支持websockify旧语法

## FAQ
### 查看websockify version
`vim /usr/bin/websockify`

### 配置tls
ref:
- [Encrypted Connections](https://github.com/novnc/websockify/wiki/Encrypted-Connections)

```bash
openssl req -new -x509 -days 365 -nodes -out self.pem -keyout self.pem # cert和key都在self.pem里
websockify --cert=/usr/share/novnc/self.pem ...
websockify --cert=/usr/share/novnc/self.pem --key=/usr/share/novnc/self.key ... # 当cert和key分开存储时使用该命令
```

> 仅使用tls时可用`ssl_only=true`, 此时也支持配置`--ssl-version tlsv1_2 --ssl-ciphers xxx`

### novnc vm鼠标不可用
测试环境是chrome 64, 为不可用, 但chrome 111是可用的