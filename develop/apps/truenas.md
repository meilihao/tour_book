# truenas scale
version: 21.04-ALPHA.1

> 22.02-RC.2至少需要python 3.9, 推荐系统自带, 否则需要处理很多python相关依赖, 很耗时且可能根本无法处理(比如debian.org bullseye下载的python3.9-minimal有libc版本要求)

## zfs
zfs pool挂在在`/mnt`下.

## 使用
truenas的systemd服务是`middlewared.service`, 程序入口是`/usr/bin/middlewared`(真实入口是`/usr/lib/python3/dist-packages/middlewared/main.py#main()`), 具体代码在`/usr/lib/python3/dist-packages/middlewared`和`/usr/local/lib/middlewared_truenas`.

middlewared设置的主要参数:
- `--debug-level=TRACE` : 使用TRACE log level
- `--log-handler=file`: 日志输出方式, 支持console, file(`/var/log/middlewared.log`). 

### console setup
进入其terminal `console setup`的命令是`/etc/netcli`

配置ip项: `Configure Network Interfaces`
配置gateway项: `Configure Default Route`

> 设置ip和gateway并重启后即可使用webui

> 使用ssh前需要在webui-System Setting-Services中打开ssh服务, truenas ssh默认禁止root登录, 修改勾选其配置`Log in as Root with Password`即可.

> 访问`localhost/api/docs`即可获取api docs(有cdn资源依赖); 访问`localhost/api/v2.0`即可获取openapi json文档, 访问`localhost/api/docs/restful`即可获取已加载openapi json的swagger ui.

### webui
使用nginx作为反向代理, 相关配置在`/etc/nginx/nginx.conf`中

> webui相关文件在: `/usr/share/truenas/webui`

## 源码
- [trunas scale构建系统](https://github.com/truenas/scale-build)
- [truenas/middleware - truenas source](https://github.com/truenas/middleware)
- [truenas scale api](https://www.truenas.com/docs/core/api/)

## 源码剖析
middlewared doc: `src/middlewared/middlewared/docs/index.rst`

### api调用
webui通过websocket api进行调用, api req由`src/middlewared/middlewared/main.py#Application.on_message`中的`if message['msg'] == 'method'`逻辑处理.

rest api req处理逻辑由`src/middlewared/middlewared/main.py#Middleware.__initialize()`中的`restful_api.register_resources()`设置, 具体在`Resource(self, self.middleware, ...)`.

```python
class Resource(object):
...
    def __init__(
        self, rest, middleware, name, service_config, parent=None,
        delete=None, get=None, post=None, put=None,
    ):
    ...
    for i in ('delete', 'get', 'post', 'put'):
            operation = getattr(self, i)
            if operation is None:
                continue
            self.rest.app.router.add_route(i.upper(), f'/api/v2.0/{path}', getattr(self, f'on_{i}')) # getattr获取__getattr__设置的`on_xxx`
            self.rest.app.router.add_route(i.upper(), f'/api/v2.0/{path}/', getattr(self, f'on_{i}'))
            self.rest._openapi.add_path(path, i, operation, self.service_config)
            self.__map_method_params(operation)
            ...
    def __getattr__(self, attr): # 为class添加on_{'on_get', 'on_post', 'on_delete', 'on_put'}方法
        if attr in ('on_get', 'on_post', 'on_delete', 'on_put'):
            do = object.__getattribute__(self, 'do')
            method = attr.split('_')[-1]

            if object.__getattribute__(self, method) is None:
                return None

            async def on_method(req, *args, **kwargs):
                resp = web.Response()
                if not self.rest._methods[getattr(self, method)]['no_auth_required']:
                    await authenticate(self.middleware, req)
                kwargs.update(dict(req.match_info))
                return await do(method, req, resp, *args, **kwargs) # 实际处理函数

            return on_method
        return object.__getattribute__(self, attr)
    async def do(self, http_method, req, resp, **kwargs): # 最终处理函数
        ...
```

添加log埋点:
- restful api : `src/middlewared/middlewared/restful.py#Resource.do`中的`result = await self.middleware.call(methodname, *method_args, **method_kwargs)`前添加`self.middleware.logger.info("--- r call: {} {} {}".format(methodname, method_args, method_kwargs))`

    或在`src/middlewared/middlewared/main.py#Middleware.call`中的开头添加`self.logger.info("--- r call: {}".format(locals()))`, 好处是不漏掉`middleware.call`嵌套调用, **推荐**

- websocket api : `src/middlewared/middlewared/main.py#Application.on_message`中的`serviceobj, methodobj = self.middleware._method_lookup(message['method'])`前添加`self.logger.info("--- w call: {} {}".format(message['method'], message.get('params') or []))`

- 可在`middlewared.py#Middleware.call()`里为打印result.

### middlewared处理http逻辑
在`src/middlewared/middlewared/main.py#Middleware.__initialize`

### 加载plugins
根据`src/middlewared/middlewared/main.py#Middleware.overlay_dirs`, 跳转可得, 是通过`src/middlewared/middlewared/utils/plugins.py#LoadPluginsMixin`来加载的:
1. 获得plugins_dirs, plugins_dirs=(main_plugins_dir + Middleware.overlay_dirs)下的plugins目录

    plugins下`.py`结尾的文件即是plugin的入口  
1. 通过load_modules()加载, 并放入services变量

    `load_modules()`原理: 递归遍历文件或目录, 找出匹配的mod, 再用`importlib.import_module`导入即可.    
1. 迭代处理services变量, 再调用`self.add_service(service)`即可载入plugin生成service.

### `middleware.call()`
代码中看到很多`self.middleware.call()`, 那`middleware`来源在哪, 以`src/middlewared/middlewared/alert/source/enclosure_status.py#EnclosureStatusAlertSource`中的`self.middleware.call('enclosure.query')`举例:
1. 点`self.middleware`跳转发现来自`src/middlewared/middlewared/alert/base.py#AlertSource`的`def __init__(self, middleware)`, 同时发现`base.py`中很多类初始化都包含`middleware`, 根据其调用的方法`grep -r "def call_sync("`和`grep -r "def call("`定位到`middlewared/main.py#Middleware.call()`.
1. 根据`grep -r "query(" $(grep -rl "enclosure")`, 发现`enclosure.query`实际是`src/middlewared/middlewared/plugins/enclosure.py#EnclosureService.query()`
1. 最后也调到了`src/middlewared/middlewared/plugins/enclosure_/ses_enclosure_linux.py#EnclosureService.get_ses_enclosures()`

### `src/middlewared/middlewared/service.py#CRUDService.query()`中的`self._config`是什么?
通过`grep -r "\._config = "`查找, 发现CRUDService父类定义`class Service(object, metaclass=ServiceBase)`中的ServiceBase.__new__()设置了`_config`, 其实就是CRUDService子类如DiskService的嵌套类`class Config`.

### table定义
`class xxxModel(sa.Model)`

### 获取middlewared.deb
根据[scale-build/conf/sources.list](https://github.com/truenas/scale-build/blob/master/conf/sources.list)找到[middlewared.deb](https://apt.tn.ixsystems.com/apt-direct/angelfish/{22.02-RC.2,nightlies}/angelfish/pool/main/m/middlewared/)

> middlewared发布RELEASE后会删除了上述url中的`22.02/angelfish`路径即没法下到RELEASE版deb.

配置vscode阅读middlewared.deb提取源码:
```bash
$ cat .vscode/settings.json 
{
    "python.autoComplete.extraPaths": [
        "./usr/lib/python3/dist-packages",
        "/home/chen/.local/lib/python3.9/site-packages",
        "/usr/lib/python3/dist-packages"
    ],
    "python.analysis.extraPaths": [
        "./usr/lib/python3/dist-packages",
        "/home/chen/.local/lib/python3.9/site-packages",
        "/usr/lib/python3/dist-packages"
    ]
}
```

运行代码:
1. 根据`dpkg --info middlewared.deb`获取依赖, 再通过apt/pip3安装依赖

    ```bash
    $ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 31AAC6F320998A97
    $ echo "deb [arch=amd64] http://apt.tn.ixsystems.com/apt-direct/angelfish/22.02-RC.2/angelfish/ truenas main" |sudo tee -a /etc/apt/sources.list.d/truenas.list # for python3-zettarepl
    $ sudo apt install python3-ldap python3-prctl python3-zettarepl
    $ install special pkgs(see FAQ): python3-ldap, python3-prctl, python3-systemd
    $ pip3 install aiohttp aiohttp-wsgi ws4py flask croniter sentry_sdk setproctitle pyjwt pycryptodomex josepy certbot-dns-cloudflare boto3 passlib html2text
    ```
1. 在提取的middlewared.deb数据的根目录中执行代码: `env PYTHONPATH=./usr/lib/python3/dist-packages:/home/chen/.local/lib/python3.9/site-packages:/usr/lib/python3/dist-packages usr/bin/middlewared -h`

    middlewared options:
    - [`--trace-malloc`](https://jira.ixsystems.com/browse/NAS-110712)


### TrueNAS-SCALE-22.02.0.iso里提取middlewared
> iso/live下的squashfs里不包含middlewared

用ncdu统计iso文件大小, 在逐个排查, 最终定位在`TrueNAS-SCALE.update`.


```bash
# binwalk TrueNAS-SCALE.update 

DECIMAL       HEXADECIMAL     DESCRIPTION
--------------------------------------------------------------------------------
0             0x0             Squashfs filesystem, little endian, version 4.0, compression:gzip, size: 1369983003 bytes, 6 inodes, blocksize: 131072 bytes, created: 2022-02-18 16:15:16
# mount -t squashfs -o loop  TrueNAS-SCALE.update  /mnt/squashfs
# cd /mnt/squashfs
# tree .
.
├── manifest.json
├── rootfs.squashfs
└── truenas_install
    ├── __init__.py
    └── __main__.py

1 directory, 4 files
# cat manifest.json |jq .
{
  "date": "2022-02-18T16:15:15.940514",
  "version": "22.02.RELEASE",
  "size": 5450135961,
  "checksums": {
    "rootfs.squashfs": "c1fdaf7032c2c2605e2c9d96e06aba086e06a643",
    "truenas_install/__main__.py": "f6eebffdce4cb8da52ade2bd16b3e9613f8c1048",
    "truenas_install/__init__.py": "da39a3ee5e6b4b0d3255bfef95601890afd80709"
  },
  "kernel_version": "5.10.93+truenas"
}
# mkdir /mnt/squashfs2
# mount -t squashfs -o loop  rootfs.squashfs  /mnt/squashfs2 # 经分析rootfs.squashfs是已安装好middlewared的镜像
```


## FAQ
### db
truenas使用sqlite3, db file在`/data/freenas-v1.db`

### scale-build执行`make checkout`报`module 'functools' has no attribute 'cache'`
根据[functools文档](https://docs.python.org/3/library/functools.html), 需要python 3.9

### `ModuleNotFoundError: No module named '_ldap'`
python3-ldap包含了`/usr/lib/python3/dist-packages/_ldap.cpython-37m-x86_64-linux-gnu.so`, 当前使用python3.9因此无法import python3.7的so.

下载[python3-ldap](https://packages.debian.org/bullseye/python3-ldap), 使用`sudo dpkg -i --ignore-depends=python3 ./python3-ldap_3.2.0-4+b3_amd64.deb`或`sudo apt install -f ./python3-ldap_3.2.0-4+b3_amd64.deb`(`apt install -f`可能不会成功, **推荐使用dpkg**)安装, 再删除`/var/lib/dpkg/status`中python3-ldap的Depends中的python3要求即: `python3 (<< 3.10), python3 (>= 3.9~)`

> 类似的还有python3-prctl, python3-systemd

### [获取升级所需manifest.json](https://update.freenas.org/scale/TrueNAS-SCALE-Angelfish-RC/manifest.json)
### [运行所需manifest.json](https://update.freenas.org/scale/TrueNAS-SCALE-Angelfish-RC/manifest.json)
由[scale_build/packages/build.py](https://github.com/truenas/scale-build/blob/TS-22.02-RC.2/scale_build/packages/build.py#L105)生成.

内容:
```json
{
    "buildtime": 1637616484,
    "train": "TrueNAS-SCALE-Angelfish-RC",
    "version": "22.02-RC.1-2"
}
```

train内容可结合`https://update.freenas.org/scale`获取.

### middlewared执行报`ImportError: cannot import name 'encode' from 'jwt'`
需要使用pyjwt

### apt install python3-zettarepl报`SyntaxError: invalid syntax`
因为安装python3-zettarepl.deb里的脚本会用到py3clean, py3compile, 而它们是python3.7的, 无法处理zettarepl里用的python3.9语法.