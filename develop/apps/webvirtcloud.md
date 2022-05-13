# webvirtcloud

## 源码
根据README里部署的nginx配置(来自`conf/nginx`), 后端服务使用8000端口, 来自gunicorn.conf.py(结合`conf/requirements.txt`, 它配合Django使用), 结合网上资料, gunicorn的入口是`manage.py`

目录:
- vrtManager : 与libvirt的交互

### 阅读
```bash
$ cat .vscode/settings.json 
{
    "python.autoComplete.extraPaths": [
        "/home/chen/test/webvirtcloud",
        "/usr/lib/python3/dist-packages"
    ],
    "python.analysis.extraPaths": [
        "/home/chen/test/webvirtcloud",
        "/usr/lib/python3/dist-packages"
    ]
}
```