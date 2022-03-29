# vscode
当前使用的插件:
1. C/C++
1. C++ Intellisense
1. Go
1. Python
1. Remote-SSH
1. Rust
1. Python
1. vscode-proto3
1. x86 and x86_64 Assembly
1. Python

## 插件
语言插件:
- C/C++
- C++ Intellisense
- Switcher : 能在头文件和 C/C++ 文件之间跳转
- shellman :  Linux shell 脚本辅助工具, 提供了便捷的shell script 自动补全和联想等功能
- Go
- Rust
- vscode-proto3
- Clang-Format : format proto3 using by google

	`apt install clang-format`

源码标签:
- Bookmarks
- TODO Tree : 显示代码中的`// TODO`

git:
- Git Graph
- GitLens

markdown:
- Markdown Preview Enhanced : Markdown 预览增强插件

增强:
- Remote Development : 远程开发
- Bracket Pair Colorizer : 彩虹括号插件, 用顔色区分括号层次
- koroFileHeader : 自动插入文件头部信息, 快捷键是`Ctrl + Alt +i`
- Chinese (Simplified) Language Pack for Visual Studio Code : vscode中文语言包

## FAQ
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