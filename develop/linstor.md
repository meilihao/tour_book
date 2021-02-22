# linstor

[LINBIT DRBD9 Stack官方PPA](https://launchpad.net/~linbit/+archive/ubuntu/linbit-drbd9-stack)

## 构建
LINBIT未直接提供构建每个linstor组件的方法, 因此构建方法是: 每个项目的`README + .gitlab-ci.yml + Dockerfile` + google

### server

```bash
sdk install gradle 6.8.2
gradle -v
sudo apt install openjdk-11-jdk # base https://github.com/LINBIT/linstor-server/blob/master/build.gradle的sourceCompatibility, 推荐用openjdk-8-jdk, 我这是因为当前已安装openjdk-11-jdk
sudo apt default-jdk-headless dh-systemd python3-all
git clone --depth=1 git@github.com:LINBIT/linstor-server.git -b v1.11.1 # 必须使用git repo否则gradlew会报错
cd linstor-server
make check-submods
```

使用方法:
```bash
./gradlew getProtoc
./gradlew assemble
ll build/distributions # 构建结果, 非deb/rpm
具体使用方法见https://github.com/LINBIT/linstor-server/blob/master/Dockerfile.test-controller
# ./gradlew clean # gradle重新构建时可用于清理

# 继续构建deb:
删除https://github.com/LINBIT/linstor-server/blob/master/debian/control中gradle的限制 # 上面已安装gradle
LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -b -uc # 自定义的LD_LIBRARY_PATH可能影响dpkg-buildpackage的构建
```

### client
```bash
pip3 install python-linstor natsort
git clone --depth=1 git@github.com:LINBIT/linstor-client.git -b v1.6.1
cd linstor-client
NO_DOC="-no-doc" make debrelease # from `.gitlab-ci.yml`
删除https://github.com/LINBIT/linstor-client/blob/master/debian/control中的python-linstor, python3-natsort依赖
make deb
```

### drbd
```bash
到https://launchpad.net/~npalix/+archive/ubuntu/coccinelle下载coccinelle
apt install spatch # spatch需>= 1.0.8, deepin 20.1的是1.0.4
git clone --depth=1 git@github.com:LINBIT/drbd.git -b drbd-9.0.27-1
cd drbd
make clean
make
sudo make install
modprobe drbd
# sudo apt install module-assistant # 构建deb需要
# LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -nc -uc -us # 打出的deb内容不对
```

### drbd-utils
```bash
apt install po4a clitest
git clone --depth=1 git@github.com:LINBIT/drbd-utils.git -b v9.15.1
cd drbd-utils
make -f Makefile.in check-submods
./autogen.sh
./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc --without-83support --without-84support # 根据autogen.sh的输出选择合适的configure参数, 部分参数来自https://github.com/LINBIT/drbd-utils/blob/master/.gitlab-ci.yml
make debrelease
make
sudo make install # 无需deb时, 直接install即可
vim debian/rules
configure_with = --prefix=/usr --localstatedir=/var --sysconfdir=/etc --without-83support --without-84support # 注释文中configure_with操作, 并使用上面的参数替代
# configure_with = --prefix=/usr --localstatedir=/var --sysconfdir=/etc \
#       --sbindir=/usr/sbin --with-udev --with-xen \
#       --with-pacemaker --with-rgmanager --without-bashcompletion

# ifneq ($(WITH_SYSTEMD),)
# configure_with += --with-systemdunitdir=/lib/systemd/system \
#       --with-initscripttype=both
# # alternatively: --with-initscripttype=systemd,
# # if you want to drop the sysv script from the package. 
# # Need to adjust below, in that case.
# endif

# ifeq (,$(filter noprebuiltman,$(DEB_BUILD_OPTIONS)))
# configure_with += --with-prebuiltman
# endif
LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -nc -uc # 使用`-nc`避免重新构建
```

## FAQ
### 
```bash
$ LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -b -uc # for drbd-utils
...
make -C documentation/v9 doc
make[2]: 进入目录“/home/chen/test/drbd-utils/documentation/v9”
test -f drbdsetup.8
test -f drbd.conf.5
test -f drbd.8
make[2]: *** [../../documentation/common/Makefile_v9_com_post:179：drbdsetup.8] 错误 1
make[2]: *** 正在等待未完成的任务....
make[2]: *** [../../documentation/common/Makefile_v9_com_post:49：drbd.conf.5] 错误 1
make[2]: *** [../../documentation/common/Makefile_v9_com_post:49：drbd.8] 错误 1
make[2]: 离开目录“/home/chen/test/drbd-utils/documentation/v9”
make[1]: *** [Makefile:112：doc] 错误 2
make[1]: 离开目录“/home/chen/test/drbd-utils”
dh_auto_build: make -j4 returned exit code 2
make: *** [debian/rules:10：build] 错误 2
dpkg-buildpackage: error: debian/rules build subprocess returned exit status 2
```

`make -C documentation/v9 doc`执行报错, 导致drbdsetup.8等文件未生成, 原因是未执行`make debrelease`. 恢复到git clone时的状态, 按照drbd-utils的流程重新构建即可.

### linstor-client输出api调用信息
client sub cmds在`linstor_client/commands/*.py`里, 它们都是调用了`/usr/lib/python3/dist-packages/linstor/linstorapi.py`里的接口, 因此:
```python
logging.basicConfig(level=logging.WARNING) => logging.basicConfig(level=logging.DEBUG)

...
    def _rest_request_base(self, apicall, method, path, body=None, reconnect=True):
        ....
            self._rest_conn.request(
                method=method,
                url=path,
                body=json.dumps(body) if body is not None else None,
                headers=headers
            )
            self._logger.info("method: {}, url: {}, body: {}, headers: {}".format(method, path, body, headers)) # append

    def __convert_rest_response(self, apicall, response, path):
        ...
            data = json.loads(resp_data)
            self._logger.info("resp.data: {}".format(data)) # append
```

> linstor -m node list # `-m`输出原始resp json

### linstor-server输出api调用信息
server api源码入口在`linstor-server/controller/src/main/java/com/linbit/linstor/api/rest/v1/*.java`, 发现代码未埋点, 改用logger

[修改`/etc/linstor/linstor.toml`](https://github.com/LINBIT/linbit-documentation/blob/master/UG9/en/administration-linstor.adoc#logging):
```toml
[logging]
   level="TRACE"
```