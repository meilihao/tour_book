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