faq

### vscode无法跳转报`Couldn't start client Rust Language Server`
vscode 菜单栏 -> File -> Preferences -> Settings, 在编辑器右上角选择`Open Settings(JSON)`按钮, 在settiong.json中追加:
```json
"rust-client.rustupPath": "$HOME/.cargo/bin/rustup"
```

在通过vscode 命令面板启动(`ctrl+shift+p`)输入`rust`, 找到启动`rust server`即可. 通常rls分析project的速度慢(vscode左下角工具栏上会看到icon转圈), 跳转等功能没法立即使用.

### [`rust-toolchain`文件](https://rust-lang.github.io/rustup/concepts/toolchains.html)的例子
1. `nightly-2020-10-19`