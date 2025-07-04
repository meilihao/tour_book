# zed

## FAQ
### 修改theme
`ctrl-k ctrl-t`

推荐主题: Andromeda

### debug
`zed --foreground`

### update
菜单栏->Open Command Palette->`auto update: check`

### [添加proxy](https://github.com/zed-industries/zed/discussions/7525)
```bash
# vim ~/.config/zed/settings.json
{
    ...
    "proxy" : "socks5://localhost:20170"
}
```

### 查看zed extensions
`ctrl-shift-x`

### [failed to spawn command "~/.local/share/zed/languages/pylsp/pylsp-venv/bin/pylsp"](https://github.com/zed-industries/zed/issues/21452)
zed原生支持python, 不知为什么没有`~/.local/share/zed/languages/pylsp/pylsp-venv/bin/pylsp`

当前解决方法, 自行安装pylsp:
```bash
cd ~/.local/share/zed/languages/pylsp
rm -rf pylsp-venv
virtualenv pylsp-venv
source pylsp-venv/bin/activate
pip install "python-lsp-server[all]"
```

### 安装或更新
`curl -f https://zed.dev/install.sh | sh`

安装或更新使用相同的命令.