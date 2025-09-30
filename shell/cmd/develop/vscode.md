# vscode

## 插件
语言插件:
- C/C++
- Switcher : 能在头文件和 C/C++ 文件之间跳转
- shellman :  Linux shell 脚本辅助工具, 提供了便捷的shell script 自动补全和联想等功能
- Go
- rust-analyzer
- Python
- Ruff : 基于 Rust 编写的高性能 Python 代码分析工具 (Python linter)
- x86 and x86_64 Assembly
- ASM Code Lens
- vscode-proto3
- Clang-Format : format proto3 using by google

    `apt install clang-format`
- CMake

db:
- database client

源码标签:
- Bookmarks
- TODO Tree : 显示代码中的`// TODO`

git:
- Git Graph
- GitLens

markdown:
- Markdown Preview Enhanced : Markdown 预览增强插件
- markdownlint

增强:
- Remote Development : 远程开发
- Remote-SSH
- Bracket Pair Colorizer : 彩虹括号插件, 用顔色区分括号层次
- koroFileHeader : 自动插入文件头部信息, 快捷键是`Ctrl + Alt +i`
- Chinese (Simplified) Language Pack for Visual Studio Code : vscode中文语言包
- SFTP
- REST Client

## setting.json
```
{
    "go.testFlags": ["-v", "-count=1"],
    "go.toolsEnvVars": {
        "GOEXPERIMENT": "jsonv2" // 开启实验特性
    }
}
```

更新后需重启vscode

## 快捷键
- ctrl + `+` : 放大
- ctrl + `-` : 缩小

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

### 插件log
Help -> `Toggle Developer Tools`

### System limit for number of file watchers reached
打开vue项目时报该错, 解决方法:
```conf
sudo vim /etc/sysctl.conf
fs.inotify.max_user_watches=524288
sudo sysctl -p
```

### vscode-proto3报`Import "xxx.proto" was not found or had errors`
在项目的`.vscode/settings.json`配置:
```json
{
    "protoc": {
        "options": [
            "--proto_path=${workspaceRoot}/pkg/types"
        ]
    }
}
```

即添加相应proto_path即可.

### java没有智能提示
`Java Language Server requires a JDK 17+ to launch itself.`

### extension log
`cmd-shift-p -> Search Show Logs -> Extension Host`

### ssh远程ui
`code --no-sandbox --user-data-dir=/root`

### Remote SSH 连接报错: Could not establish connection to "xxx". XHR failed.
参考:
- [Using “Remote SSH” in VSCode on a target machine that only allows inbound SSH connections](https://stackoverflow.com/questions/56718453/using-remote-ssh-in-vscode-on-a-target-machine-that-only-allows-inbound-ssh-co/56781109#56781109)

Remote SSH连接对端需要先在对端安装`vscode-server-linux-x64.tar.gz`, 它的获取方法是:
`wget -nv -O vscode-server-linux-x64.tar.gz https://update.code.visualstudio.com/commit:$COMMIT_ID/server-linux-${arch}/${quality}`

参数:
- arch : x64|arm64, 这里是x64
- COMMIT_ID : 打开vscode -> Help -> About, 其中的Commit即COMMIT_ID
- quality : 固定为"stable"

上传到对端, 并安装vscode-server-linux-x64.tar.gz:
```bash
$ mkdir -pv ~/.vscode-server/bin/$COMMIT_ID
$ tar -xvzf vscode-server-linux-x64.tar.gz --strip-components 1 -C ~/.vscode-server/bin/$COMMIT_ID
```

保证有$HOME/.vscode-server/bin/$COMMIT_ID/server.sh，这样vscode再连接服务器就不会报错了.

### python env
参考:
- [Using Python environments in VS Code](https://code.visualstudio.com/docs/python/environments#_environment-variable-definitions-file)

`${project}/.vscode/settings.json`:
```json
{
	"python.envFile": "${workspaceFolder}/.vscode/.env",
	"python.pythonPath":"/usr/bin/python3"
}
```

`${workspaceFolder}/.vscode/.env`:
```conf
PYTHONPATH=/xxx1:/xxx2
```

### Remote-SSH报`Failed to set up socket for dynamic port forward to remote port <port>`
ref:
- [VSCode Remote SSH Connection Failed](https://stackoverflow.com/questions/60507713/vscode-remote-ssh-connection-failed)

修改remote端的sshd配置:
```bash
# vim /etc/ssh/sshd_config
...
AllowTcpForwarding  = yes
...
# systemctl restart sshd
```
### python调试时`Copy Value`的值被截断
可在`DEBUG CONSOLE`打印一下再复制

### 正则替换
`^(.*)$`-> `Status$1="$1" //`

### 为c/c++添加heder搜索路径
按Ctrl+Shift+P 输入configuration, 找到`C/C++: Edit Configuration(JSON)` 在c_cpp_properties.json中includePath字段中添加待添加的SDK或者库的头文件路径

### 当前窗口页签(标题栏)不明显
File->Preferences->Settings->Workbench -> Apperance ->`Color Customizations`, 点击`Edit in settings.json`, 添加:
```json
"workbench.colorCustomizations": {
    "tab.activeBackground": "#555555",
    "tab.hoverBackground": "#464646"
}
```

### 顶部菜单栏不明显
```json
"workbench.colorCustomizations": {
    "titleBar.activeBackground": "#555555", // 活动窗口的顶部菜单栏背景色
    "titleBar.inactiveBackground": "#464646", // 非活动窗口的顶部菜单栏背景色
    //"titleBar.activeForeground": "#ffffff", // 活动窗口的顶部菜单栏文字颜色
    //"titleBar.inactiveForeground": "#cccccc" // 非活动窗口的顶部菜单栏文字颜色
},
```

### golint屏蔽警告`receiver name should be a reflection of its identity; don't use generic names such as "this" or "self" (ST1006)`
ref:
- [Configuration](https://staticcheck.dev/docs/configuration/)

1. 创建项目根目录创建staticcheck.conf

    ```json
    checks = ["-ST1006"]
    ```
1. 重启vscode

### SFTP配置
1. 按下 F1（或 Ctrl+Shift+P / Cmd+Shift+P）
1. 输入并选择：SFTP: Config → 会生成 .vscode/sftp.json 文件

    sftp.json:
    ```json
    {
        "name": "env-sit",
        "host": "120.55.x.y",
        "protocol": "sftp",
        "port": 22,
        "username": "root",
        "remotePath": "/root/demo",
        "privateKeyPath": "~/.ssh/dev_ed25519",
        "uploadOnSave": false,
        "useTempFile": false,
        "openSsh": false
    }

    // 配置多个env
    {
        "protocol": "sftp",
        "port": 22,
        "remotePath": "/root/demo",
        "uploadOnSave": false,
        "useTempFile": false,
        "openSsh": false,
        "profiles": {
            "dev": {
                "host": "120.26.1.2",
                "username": "root",
                "privateKeyPath": "~/.ssh/dev_ed25519"
            },
            "sit":{
                "host": "120.55.1.2",
                "username": "root",
                "privateKeyPath": "~/.ssh/dev_ed25519"
            }
        }
    }
    ```