# vpn
## FAQ
### openvpn登入报`failed to negotiate cipher with server.  Add the server's cipher ('AES-128-CBC') to --data-ciphers (currently 'AES-256-GCM:AES-128-GCM:CHACHA20-POLY1305') if you want to connect to this server`
原因: fedora 40的openvpn太新, server端不支持默认的data-ciphers
追加参数`--data-ciphers 'AES-128-CBC'`