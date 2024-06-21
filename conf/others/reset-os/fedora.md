# fedora
```bash
$ nvim /etc/dnf/dnf.conf # max_parallel_downloads=5, 加快下载包; 追加`fastestmirror=True`, 即选择最快源, 如果手动换源了, 就忽略它
$ sed -e 's|^metalink=|#metalink=|g' \
    -e 's|^#baseurl=http://download.example/pub/fedora/linux|baseurl=https://mirrors.tuna.tsinghua.edu.cn/fedora|g' \
    -i.bak \
    /etc/yum.repos.d/fedora.repo \
    /etc/yum.repos.d/fedora-modular.repo \
    /etc/yum.repos.d/fedora-updates.repo \
    /etc/yum.repos.d/fedora-updates-modular.repo # use mirror  
$ sudo dnf update && sudo dnf upgrade
$ sudo fwupdmgr refresh --force
$ sudo fwupdmgr get-updates
$ sudo fwupdmgr update
$ sudo dnf autoremove
$ sudo dnf remove --oldinstallonly # 删除旧kernel
$ sudo dnf install python3-pip
$ python -m pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --upgrade pip
$ pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
```
