# zed

## 插件
1. ruff

## setting
```json
{
  "ui_font_size": 16,
  "buffer_font_size": 16,
  "theme": {
    "mode": "system",
    "light": "One Dark",
    "dark": "One Dark"
  },
  "proxy": "socks5://127.0.0.1:20170",
  "languages": {
    "Python": {
      "language_servers": ["!pylsp", "pyright", "ruff"]
    }
  }
}
```

> 仅用pylsp, 占用cpu太高. `pyright + ruff`首次检查耗时过长, 需等待. ruff目前提示不如pyright好用, 因此排后面.

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
    "proxy" : "socks5://127.0.0.1:20170"
}
```

### 查看zed extensions
`ctrl-shift-x`

### [failed to spawn command "~/.local/share/zed/languages/pylsp/pylsp-venv/bin/pylsp"](https://github.com/zed-industries/zed/issues/21452)
**不推荐使用pylsp, cpu占用太高, 改为pyright/ruff**

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

### log位置
`/home/chen/.local/share/zed/logs`, 包括crash log