# vpn

## 连接vpn
参考:
- [How to configure the SSL VPN on Ubuntu](https://community.sophos.com/kb/en-us/125368)
- [OpenVPN Client Sophos](https://www.systemhaus-brandenburg.de/de/Ubuntu_Linux_-_OpenVPN_-_Sophos_1981.html)

步骤:
1. 获取到openvpn配置, 这里是`chen@vpn.example.com.ovpn`
1. 使用配置创建隧道, 通常有两种方法:
    1. 直接使用openvpn: `openvpn --config chen@vpn.example.com.ovpn`
    1. 使用ui工具, 比如`NetworkManager`或deepin 设置中心的vpn.

## FAQ
### openvpn登入报`failed to negotiate cipher with server.  Add the server's cipher ('AES-128-CBC') to --data-ciphers (currently 'AES-256-GCM:AES-128-GCM:CHACHA20-POLY1305') if you want to connect to this server`
原因: fedora 40的openvpn太新, server端不支持默认的data-ciphers
追加参数`--data-ciphers 'AES-128-CBC'`