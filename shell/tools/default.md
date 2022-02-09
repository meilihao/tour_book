# default
## ubuntu
### 默认编辑器
1. sudo update-alternatives --config editor
1. select-editor
1. echo export EDITOR=/usr/bin/vim >> ~/.bashrc # 终极方法, 需重启terminal

### python3
```bash
$ sudo update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.9 1
$ sudo update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 2
$ sudo update-alternatives --config python3
```

> 语法为: `update-alternatives --install <link> <name> <path> <priority> [--slave link name path]...`

> 语法为: `update-alternatives --remove name path`