# vscode

## FAQ
### 启动code没反映
env : xubuntu 20.04(xfce) + tightvncserver 

```bash
$ code # 无仍会输出, journalctl也没有输出. 重装vscode无效果
$ /usr/share/code
Xlib: Xlib: extension "XInputExtension" missing on display ":1"
...
$ sudo cp /usr/lib/x86_64-linux-gnu/libxcb.so.1 /usr/lib/x86_64-linux-gnu/libxcb.so.1.bak
$ sudo sed -i 's/BIG-REQUESTS/_IG-REQUESTS/' /usr/lib/x86_64-linux-gnu/libxcb.so.1
$ code # 正常启动, 但编辑时发现"删除"键无效
```

最后网上查到是tightvncserver太旧(2009年后不在维护), 换tigervnc即可, 具体tigervnc安装方法见[vnc.md](vnc.md)