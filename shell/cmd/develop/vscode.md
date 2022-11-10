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