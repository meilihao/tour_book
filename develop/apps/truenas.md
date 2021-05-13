# truenas scale
version: 21.04-ALPHA.1

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

## FAQ
### db
truenas使用sqlite3, db file在`/data/freenas-v1.db`