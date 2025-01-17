# truenas scale
ref:
- [TrueNAS SCALE Clustering Overview](https://www.truenas.com/blog/truenas-scale-clustering/)
- [2024最近，14款免费NAS操作系统及引导！黑群、黑威、Unraid开心版、TrueNAS、国产开源等，先码后看！](https://post.smzdm.com/p/azomrxmn/)
- [fnOS](https://news.mydrivers.com/1/999/999891.htm)
- [truenas scale 下载地址(**全部版本**)](https://download.sys.truenas.net/)

version: 24.04

> `22.02至少需要python 3.9/24.04至少需要python 3.11`. 推荐系统自带, 否则需要处理很多python相关依赖, 很耗时且可能根本无法处理(比如debian.org bullseye下载的python3.9-minimal有libc版本要求)

## 安装
ref:
- [TrueNAS安装及使用教程](https://www.ethanzhang.xyz/2023/05/14/TrueNAS%E5%AE%89%E8%A3%85%E5%8F%8A%E4%BD%BF%E7%94%A8%E6%95%99%E7%A8%8B/)
- [设置TrueNAS-SCALE为中文](https://mi-d.cn/9431)

## zfs
zfs pool挂载在`/mnt`下.

## 使用
truenas的systemd服务是`middlewared.service`, 程序入口是`/usr/bin/middlewared`(真实入口是`/usr/lib/python3/dist-packages/middlewared/main.py#main()`), 具体代码在`/usr/lib/python3/dist-packages/middlewared`.

middlewared设置的主要参数:
- `--debug-level=TRACE` : 使用TRACE log level
- `--log-handler=file`: 日志输出方式, 支持console, file(`/var/log/middlewared.log`). 

### console setup
进入其terminal `console setup`的命令是`/usr/bin/cli --menu --pager`(by python3)

> 以前版本, 比如22.02是`/etc/netcli`

配置ip项: `Configure network interfaces`
配置gateway项: `Configure static routes`

> 设置ip和gateway并重启后即可使用webui

> 使用ssh前需要在`webui->System Setting->Services`中打开ssh服务, truenas ssh默认禁止root登录, 修改勾选其配置`Log in as Root with Password`即可.

> 访问`localhost/api/docs`即可获取api docs(有cdn资源依赖); 访问`localhost/api/v2.0`即可获取openapi json文档, 访问`localhost/api/docs/restful`即可获取已加载openapi json的swagger ui.

### webui
使用nginx作为反向代理, 相关配置在`/etc/nginx/nginx.conf`中

> webui相关文件在: `/usr/share/truenas/webui`

## 源码
- [trunas scale构建系统](https://github.com/truenas/scale-build)
- [truenas/middleware - truenas source](https://github.com/truenas/middleware)
- [truenas scale api](https://www.truenas.com/docs/scale/api/)

## 源码剖析
middlewared doc: `src/middlewared/middlewared/docs/index.rst`

- jobs是在memory中的`self.jobs = JobsQueue(self)`
- event_register是注册event的入口
- `/core/get_services`: 返回所有服务信息, 包括`config`
- `/core/get_methods`:

    example:

    ```json
    {
        "service": "core",
        "cli":"false"
    }
    ```
- `self.middleware.call('pool.scrub.scrub', BOOT_POOL_NAME)`是PoolScrubService的namespace

    ```py
    class PoolScrubService(CRUDService):

    class Config:
        datastore = 'storage.scrub'
        datastore_extend = 'pool.scrub.pool_scrub_extend'
        datastore_prefix = 'scrub_'
        namespace = 'pool.scrub'
        cli_namespace = 'storage.scrub'
    ```

### 远程调试
见api文档的`/core/debug`接口

### api调用
ref:
- [SCALE API Reference](https://www.truenas.com/docs/scale/api/)

rest api到ws api的映射: `/device/get_info` -> `device.get_info()`, 实现是`DeviceService(Service)`的`async def get_info(self, _type)`, 即逻辑实现均是`XXXService(Service)`的`yyy`方法.

没在代码里找到api文档, 推测它可能是具体api实现上的`@accepts()/@returns()`生成的.

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

    或在`src/middlewared/middlewared/main.py#Middleware.call`中的开头添加`self.logger.info("--- r call: {}\n".format(locals()))`, 好处是不漏掉`middleware.call`嵌套调用, **推荐**

- websocket api : `src/middlewared/middlewared/main.py#Application.on_message`中的`serviceobj, methodobj = self.middleware._method_lookup(message['method'])`前添加`self.logger.info("--- w call: {} {}".format(message['method'], message.get('params') or []))`

- 可在`middlewared.py#Middleware.call()`里为打印result.


> log位置: /var/log/middlewared.log

### middlewared 加载plugins
- version: 24.04.1.1

在`src/middlewared/middlewared/main.py#Middleware.__initialize`

入口在`src/middlewared/middlewared/main.py#Middleware.__initialize`的`setup_funcs = await self.__plugins_load()`->跳转可得, 是通过`src/middlewared/middlewared/utils/plugins.py#LoadPluginsMixin`的`_load_plugins`来加载的:
1. 获得plugins_dirs即plugins目录
1. 通过load_modules()加载, 并放入services变量

    `load_modules()`原理: 递归遍历文件或目录, 找出匹配的mod, 再用`importlib.import_module`导入即可.

    `mod.setup`=`async def setup(middleware)`用于初始化mod

    > 使用`--debug-level TRACE`可在`/var/log/middlewared.log`中看到`loaded plugin xxx`. `on_modules_loaded()`是`plugin loaded`后调用的.
1. 迭代处理services变量, 再调用`self.add_service(service)`即可载入plugin生成service.
1. `self.__plugins_setup(setup_funcs)`


```python
# utils/plugins.py
class LoadPluginsMixin(SchemasMixin):

    def __init__(self):
        self._services = {}
        self._services_aliases = {}
        super().__init__()

    def _load_plugins(self, on_module_begin=None, on_module_end=None, on_modules_loaded=None):
        from middlewared.service import Service, CompoundService, ABSTRACT_SERVICES

        services = []
        plugins_dir = os.path.realpath(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..', 'plugins')) # /usr/lib/python3/dist-packages/middlewared/plugins
        if not os.path.exists(plugins_dir):
            raise ValueError(f'plugins dir not found: {plugins_dir}')

        for mod in load_modules(plugins_dir, depth=1):
            if on_module_begin:
                on_module_begin(mod)

            services.extend(load_classes(mod, Service, ABSTRACT_SERVICES)) # [<class 'middlewared.plugins.disk.DiskService'>, ...]

            if on_module_end:
                on_module_end(mod)

        def key(service):
            return service._config.namespace # 部分没有Config的service由ServiceBase添加默认Config, 比如`<class 'middlewared.plugins.usage.UsageService'>`

        for name, parts in itertools.groupby(sorted(set(services), key=key), key=key): # sorted(set(services), key=key) = [<class 'middlewared.plugins.acme_protocol_.acme_svc.ACMEService'>, <class 'middlewared.plugins.acme_protocol_.issue_cert.ACMEService'>, <class 'middlewared.plugins.acme_protocol.DNSAuthenticatorService'>, <class 'middlewared.plugins.acme_protocol_.challenge.DNSAuthenticatorService'>, <class 'middlewared.plugins.acme_protocol_.schema.DNSAuthenticatorService'>, <class 'middlewared.plugins.acme_protocol.ACMERegistrationService'>,...]. `itertools.groupby(sorted(set(services), key=key), key=key)`根据key聚合后的结果是: [{key:"acme", values: [<class 'middlewared.plugins.acme_protocol_.acme_svc.ACMEService'>, <class 'middlewared.plugins.acme_protocol_.issue_cert.ACMEService'>]}]
            parts = list(parts)

            if len(parts) == 1:
                service = parts[0](self)
            else:
                service = CompoundService(self, [part(self) for part in parts]) # 将key(service)相同的聚合为一个CompoundService

            if not service._config.private and not service._config.cli_private and not service._config.cli_namespace:
                raise RuntimeError(f'Service {service!r} does not have CLI namespace set')

            self.add_service(service)

        if on_modules_loaded:
            on_modules_loaded()

        # Now that all plugins have been loaded we can resolve all method params
        # to make sure every schema is patched and references match
        self._resolve_methods(list(self._services.values()), self.get_events())

    def add_service(self, service):
        self._services[service._config.namespace] = service
        if service._config.namespace_alias:
            self._services_aliases[service._config.namespace_alias] = service

    def get_service(self, name):
        service = self._services.get(name)
        if service:
            return service
        return self._services_aliases[name]

    def get_services(self):
        return self._services

# utils/plugins.py
# depth控制importlib.import_module的层数
def load_modules(directory, base=None, depth=0): # 首次进入directory=/usr/lib/python3/dist-packages/middlewared/plugins, 其他比如`/usr/lib/python3/dist-packages/middlewared/plugins/service_`
    directory = os.path.normpath(directory)
    if base is None:
        middlewared_root = os.path.dirname(os.path.dirname(__file__)) # /usr/lib/python3/dist-packages/middlewared
        if os.path.commonprefix((f'{directory}/', f'{middlewared_root}/')) == f'{middlewared_root}/':
            base = '.'.join(
                ['middlewared'] +
                os.path.relpath(directory, middlewared_root).split('/') # os.path.relpath(directory, middlewared_root) = "plugins"
            ) # middlewared.plugins
        else:
            for new_module_path in sys.path:
                if os.path.commonprefix((f'{directory}/', f'{new_module_path}/')) == f'{new_module_path}/':
                    break
            else:
                new_module_path = os.path.dirname(directory)
                logger.debug("Registering new module path %r", new_module_path)
                sys.path.insert(0, new_module_path)

            base = '.'.join(os.path.relpath(directory, new_module_path).split('/'))

    _, dirs, files = next(os.walk(directory)) # dirs = ["service_",...]; files = ["disk.py", ...]
    for f in filter(lambda x: x[-3:] == '.py' and x.find('_freebsd') == -1, files):
        yield importlib.import_module(base if f == '__init__.py' else f'{base}.{f[:-3]}') # middlewared.plugins.disk

    for f in filter(lambda x: x.find('_freebsd') == -1, dirs):
        if depth > 0:
            path = os.path.join(directory, f) # /usr/lib/python3/dist-packages/middlewared/plugins/service_
            yield from load_modules(path, f'{base}.{f}', depth - 1)
```

### `middleware.call()`
代码中看到很多`self.middleware.call()`, 那`middleware`来源在哪, 以`src/middlewared/middlewared/alert/source/enclosure_status.py#EnclosureStatusAlertSource`中的`self.middleware.call('enclosure.query')`举例:
1. 点`self.middleware`跳转发现来自`src/middlewared/middlewared/alert/base.py#AlertSource`的`def __init__(self, middleware)`, 同时发现`base.py`中很多类初始化都包含`middleware`, 根据其调用的方法`grep -r "def call_sync("`和`grep -r "def call("`定位到`middlewared/main.py#Middleware.call()`.
1. 根据`grep -r "query(" $(grep -rl "enclosure")`, 发现`enclosure.query`实际是`src/middlewared/middlewared/plugins/enclosure.py#EnclosureService.query()`
1. 最后也调到了`src/middlewared/middlewared/plugins/enclosure_/ses_enclosure_linux.py#EnclosureService.get_ses_enclosures()`

查找call对于的服务名:
1. plugin可根据其namespace来查找

### `src/middlewared/middlewared/service.py#CRUDService.query()`中的`self._config`是什么?
通过`grep -r "\._config = "`查找, 发现CRUDService父类定义`class Service(object, metaclass=ServiceBase)`中的ServiceBase.__new__()设置了`_config`, 其实就是CRUDService子类如DiskService的嵌套类`class Config`(但经过metaclass修改).

以DiskService的`class Config`举例:
```python
class DiskService(CRUDService):

    class Config:
        datastore = 'storage.disk' # db name storage_disk
        datastore_prefix = 'disk_'
        datastore_extend = 'disk.disk_extend' # DiskService.disk_extend方法
        datastore_extend_context = 'disk.disk_extend_context' # DiskService.disk_extend_context方法
        datastore_primary_key = 'identifier'
        datastore_primary_key_type = 'string'
        event_register = False
        event_send = False
        cli_namespace = 'storage.disk'
```

middleware.call('system.advanced.config'))['kdump_enabled']:
1. SystemAdvancedService(ConfigService)
1. ConfigService.config

    ```python
    options['extend'] = self._config.datastore_extend
    options['extend_context'] = self._config.datastore_extend_context
    options['prefix'] = self._config.datastore_prefix
    ```
1. self.middleware.call('datastore.config', datastore, options)
1. Config.config()
1. `self.query(name, [], options)`

    - columns.extend(list(alias.c)): 执行datastore_exten

### table定义
`class xxxModel(sa.Model)`

### EventSource
`SystemHealthEventSource(EventSource)`初始化代码:
```py
# common/event_source/manager.py
async def subscribe(self, subscriber, ident, name, arg):
            ...
            self.middleware.logger.trace("Creating new instance of event source %r:%r", name, arg)
            self.instances[name][arg] = self.event_sources[name](
                self.middleware, name, arg,
                functools.partial(self._send_event, name, arg),
                functools.partial(self._unsubscribe_all, name, arg),
            )
```

### 开发者模式
truenas 24.04 rootfs只读, 解除需开启该模式[`install-dev-tools`](https://www.truenas.com/docs/scale/scaletutorials/systemsettings/advanced/developermode/)

[也可使用`sudo /usr/local/libexec/disable-rootfs-protection`, 但没有enable-rootfs-protection, 且未验证](https://forums.truenas.com/t/truenas-scale-24-04-0-root-fs-read-only/4235/4)

### ~~获取middlewared.deb~~
根据[scale-build/conf/sources.list](https://github.com/truenas/scale-build/blob/master/conf/sources.list)找到[middlewared.deb](https://apt.tn.ixsystems.com/apt-direct/angelfish/{22.02-RC.2,nightlies}/angelfish/pool/main/m/middlewared/), 从22.02发布后, truenas删除了上述url中的`22.02/angelfish`路径即没法下到RELEASE版deb, 同时该方法获取middlewared还要解决包依赖问题, 因此**应从iso中提前源码**.

配置vscode阅读middlewared.deb提取源码:
```bash
$ cat .env
PYTHONPATH=./usr/lib/python3/dist-packages:/home/chen/.local/lib/python3.9/site-packages:/usr/lib/python3/dist-packages
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


os提取版middlewared的vscode配置(**推荐, 毕竟TrueNAS-SCALE-22.02.0.iso里的镜像已弄好依赖**):
1. 将middlewared.deb解压到/opt/mark/test/truenas_deb, 并删除`/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages/middlewared*`
1. 将iso里的TrueNAS-SCALE.update挂载到/mnt/squashfs, 再将/mnt/squashfs/rootfs.squashfs挂载到/mnt/squashfs2
1. 将`/mnt/squashfs2/usr/lib/python3/dist-packages/middlewared*`拷贝到`/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages`
1. 配置/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages/middlewared

    ```bash
    $ cat .env
    PYTHONPATH=/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages:/mnt/squashfs2/usr/lib/python3/dist-packages:/usr/lib/python3/dist-packages
    $ cat .vscode/settings.json 
    {
        "python.autoComplete.extraPaths": [
            "/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages",
            "/mnt/squashfs2/usr/lib/python3/dist-packages",
            "/usr/lib/python3/dist-packages"
        ],
        "python.analysis.extraPaths": [
            "/opt/mark/test/truenas_deb/usr/lib/python3/dist-packages",
            "/mnt/squashfs2/usr/lib/python3/dist-packages",
            "/usr/lib/python3/dist-packages"
        ]
    }
    ```

### TrueNAS-SCALE-24.04.1.1.iso里提取middlewared
> iso/live下的squashfs里不包含middlewared

用ncdu统计iso文件大小, 在逐个排查, 最终定位在`TrueNAS-SCALE.update`.


```bash
# binwalk TrueNAS-SCALE.update 

DECIMAL       HEXADECIMAL     DESCRIPTION
--------------------------------------------------------------------------------
0             0x0             Squashfs filesystem, little endian, version 4.0, compression:gzip, size: 1369983003 bytes, 6 inodes, blocksize: 131072 bytes, created: 2022-02-18 16:15:16
# 先用文件管理器挂载iso
# mount -t squashfs -o loop  /media/chen/ISOIMAGE/TrueNAS-SCALE.update  /mnt/truenas/squashfs
# cd /mnt/truenas/squashfs
# tree .
.
├── manifest.json
├── rootfs.squashfs
└── truenas_install
    ├── fhs.py
    ├── __init__.py
    └── __main__.py

1 directory, 5 files
# cat manifest.json |jq .
{
  "date": "2024-05-29T16:30:12.599848",
  "version": "24.04.1.1",
  "size": 5232902963,
  "checksums": {
    "rootfs.squashfs": "71f4721ceb9e40ce8a34c38b1326052f0ae4f3d5",
    "truenas_install/__main__.py": "54e3669c57fef1a9714e182ec9a16f0e72e951db",
    "truenas_install/fhs.py": "94bf7b8fe25f6d3c8d95c788c95c664bedb08b33",
    "truenas_install/__init__.py": "da39a3ee5e6b4b0d3255bfef95601890afd80709"
  },
  "kernel_version": "6.6.29-production+truenas"
}
# mkdir /mnt/truenas/squashfs2
# mount -t squashfs -o loop  /mnt/truenas/squashfs/rootfs.squashfs /mnt/truenas/squashfs2 # 经分析rootfs.squashfs是已安装好middlewared的镜像
```

vscode源码阅读配置:
```bash
$ scp -r /mnt/truenas/squashfs2/usr/lib/python3/dist-packages/middlewared ~/test/truenas/deb # 因为iso挂载的squashfs是只读的
$ vim /test/truenas/deb/middlewared/.env
PYTHONPATH=/home/chen/test/truenas/deb:/mnt/truenas/squashfs2/usr/lib/python3/dist-packages:/usr/lib/python3/dist-packages
$ cat /test/truenas/deb/middlewared/.vscode/settings.json
{
    "python.autoComplete.extraPaths": [
        "/home/chen/test/truenas/deb",
        "/mnt/truenas/squashfs2/usr/lib/python3/dist-packages",
        "/usr/lib/python3/dist-packages"
    ],
    "python.analysis.extraPaths": [
        "/home/chen/test/truenas/deb",
        "/mnt/truenas/squashfs2/usr/lib/python3/dist-packages",
        "/usr/lib/python3/dist-packages"
    ]
}
```

### TrueNAS-SCALE 22.02/24.04 iso安装原理
见[scale-build/conf/cd-files/](https://github.com/truenas/scale-build/tree/TS-22.02.0.1/conf/cd-files)

推测:
1. 进入live os后由systemd通过`lib/systemd/system/multi-user.target.wants`即`lib/systemd/system/mount-cd.service`挂载iso到`/cdrom`
1. 之后由`root/.bash_profile`调用live os里的[`/sbin/truenas-install`](https://github.com/truenas/truenas-installer/blob/TS-22.02.0.1/usr/sbin/truenas-install)执行安装

    1. `mount /cdrom/TrueNAS-SCALE.update /mnt -t squashfs -o loop`
    1. `(cd /mnt && echo "$json" | python3 -m truenas_install)`即执行`/mnt/truenas_install`里的代码

可结合iso的fs(`live/filesystem.squashfs`)一起验证.

## FAQ
### db
truenas使用sqlite3, db file在`/data/freenas-v1.db`

> factory-v1.db是初始db, freenas-v1.db 当前使用db,是factory-v1.db复制而来 by `plugins/config.py#reset()`.

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

### middlewared (worker)
`systemctl status middlewared`显示的`middlewared (worker)`进程代码在`worker.py#worker_init`

修改`main.py#__init_procpool`里的max_workers=1便于gdb调试


```bash
# python -m pdb /usr/bin/middlewared --debug-level TRACE --log-handler=file
(Pdb) b /usr/lib/python3/dist-packages/middlewared/main.py:1979 # 进入`self.__plugins_load()`, 因为py是解释型, 因此使用file:lino的形式设置断点
(Pdb) b /usr/lib/python3/dist-packages/middlewared/utils/plugins:24 # `for new_module_path in sys.path`
(Pdb) b /usr/lib/python3/dist-packages/middlewared/utils/plugins:103 # `for mod in load_modules(plugins_dir, depth=1)`
(Pdb) c
(Pdb) s
```

推荐**vscode**:
```bash
# vim .vscode/launch.json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "truenas",
            "type": "debugpy",
            "request": "launch",
            "program": "/usr/bin/middlewared",
            "console": "integratedTerminal",
            "args": [
                "--debug-level=TRACE",
                "--log-handler=file"
            ],
            "justMyCode": false
        }
    ]
}
```

justMyCode=false, 不能遗漏, 否则可能会跳过断点(断点icon变空心圆)

> pycharmce不支持远程调试, 在truenas上安装lxqt后启动发现`/etc/systemd/system/default.target`指向的`/lib/systemd/system/truenas.target`不存在.

pycharm支持`OpenSSH config and authentication agent`来使用ssh, 打开远程项目创建`Run/Debug Configurations`->`Python`:
Name:truenas
script: /usr/bin/middlewared
arg: --debug-level=TRACE --log-handler=file
Environment variables: PYTHONUNBUFFERED=1;PYTHONPATH=/usr/lib/python3.11:/usr/lib/python3/dist-packages # PYTHONPATH缺失会导致启动middlewared失败

> pycharm remote端需pycharm组件, 该步骤由pycharm在创建ssh connections时自动处理.

### 没找到`self.middleware.call('vm.device.query')}`
实际走入了`class VMDeviceService(CRUDService)`中CRUDService的query()方法.


### [构建`github.com/truenas/py-libzfs`, 运行`./configure`报`A working zfs header is required`](https://github.com/truenas/py-libzfs/issues/107)
`./configure --prefix=/usr`

### [构建`github.com/truenas/py-libzfs`, 运行`sudo make install`报`No module named 'Cython'`]
env: Ubuntu 20.04

`apt install cython`仍旧报错, `pip3 install cython`后正常

### [构建`github.com/truenas/py-libzfs`, 运行`libzfs.c:6:10: fatal error: Python.h: No such file or directory`]
`apt install python3-dev`

### sshd_config被还原
truenas scale 24.04具有还原机制, 手动修改后, 过段时间后会被根据db里的配置还原.