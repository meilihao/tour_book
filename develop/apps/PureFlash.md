# PureFlash
ref:
- [全闪分布式存储之PureFlash](https://cloud.tencent.com/developer/article/2363606)
- [部署](https://github.com/cocalele/PureFlash/blob/master/docker/run-all.sh)

## 组件
1. pfconductor

    集群控制模块

    res:
    - [数据库初始化脚本](https://github.com/cocalele/pfconductor/blob/master/res/init_s5metadb.sql)
    - [pfcli](https://github.com/cocalele/pfconductor/blob/master/pfcli)

## build
ref:
- [build_and_run.txt](https://github.com/cocalele/PureFlash/blob/master/build_and_run.txt)

os: CentOS Stream 9

1. deps
```bash
# wget https://rpmfind.net/linux/centos-stream/9-stream/CRB/x86_64/os/Packages/ninja-build-1.10.2-6.el9.x86_64.rpm # from https://rpmfind.net/linux/rpm2html/search.php?query=ninja-build
# ninja --version
# wget https://dlcdn.apache.org//ant/binaries/apache-ant-1.10.13-bin.tar.bz2
# tar -xf apache-ant-1.10.13-bin.tar.bz2 && mv apache-ant-1.10.13 /opt && ln -s /opt/apache-ant-1.10.13/bin/ant /usr/local/bin/
# ant -version
# dnf install java-17-openjdk java-17-openjdk-devel # java-17-openjdk-devel is for javac
# java --version
# dnf config-manager --set-enabled crb # for cppunit-devel
# dnf install epel-release # for gperftools-devel
# dnf install libuuid libuuid-devel gperftools-devel cppunit-devel nasm  libaio-devel liburing-devel rdma-core-devel libcurl-devel # nasm,libaio-devel,liburing-devel,rdma-core-devel,libcurl-devel for build PureFlash
# dnf install cmake libtool
# cmake --version
```

1. build
```bash
# git clone --depth 1 git@github.com:cocalele/PureFlash.git
# set PFHOME=$(pwd)/PureFlash
# cd PureFlash
# git submodule update --init
# --- build zookeeper
# cd thirdParty/zookeeper
# ant compile_jute
# cd zookeeper-client/zookeeper-client-c
# autoreconf -if
# ./configure --enable-debug --without-cppunit
# make
# --- build PureFlash
# cd PureFlash
# export PUREFLASH_HOME=`pwd`
# mkdir build_deb; cd build_deb
# cmake -GNinja -DCMAKE_BUILD_TYPE=Debug -DCMAKE_MAKE_PROGRAM=ninja ..
# pushd .
# vim ../common/include/pf_fixed_size_queue.h # see FAQ
# vim build.ninja # see FAQ
# cd ../thirdParty/isa-l_crypto && make install # see FAQ
# popd
# ninja
# --- build pfconductor
# git clone --depth 1 git@github.com:cocalele/pfconductor.git
# cd pfconductor
# git submodule update --init
# ant -f jconductor.xml
# ./pfcli --help
```

## FAQ
### `error: package java.net.http does not exist`
[Package java.net.http](https://docs.oracle.com/en/java/javase/17/docs/api/java.net.http/java/net/http/package-summary.html)是java 11开始加入的, 刚开始是安装了jdk 8, 后来安装jdk 17, 检查后JAVA_HOME使用的还是java 8的, 经检查是`~/.bashrc`里设置的, 取消即可.

官方编译用了[jdk 14](https://github.com/cocalele/PureFlash/blob/master/build_and_run.txt).

### build PureFlash get: `error: ‘logic_error’ is not a member of ‘std’`
other like: `error: ‘runtime_error’ is not a member of ‘std’`

add `[stdexcept](https://stackoverflow.com/questions/4861777/missing-stdruntime-error-in-qtmingw)`:
```bash
# vim common/include/pf_fixed_size_queue.h
...
#include <stdlib.h>
#include <errno.h>
#include "pf_utils.h"
#include "pf_lock.h"
#include <stdexcept> # add this line
...
```

### build PureFlash get: `fatal error: isa-l_crypto/aes_cbc.h: No such file or directory`
执行`ninja`前先在`thirdParty/isa-l_crypto`执行`make install`

### build PureFlash get: `ld: cannot find -luring`
ref:
- [`"lib3 dynamically and lib2 statically":gcc program.o -llib1 -Wl,-Bstatic -llib2 -Wl,-Bdynamic -llib3`](https://stackoverflow.com/questions/6578484/telling-gcc-directly-to-link-a-library-statically)

```bash
# ninjia
...
/usr/bin/c++ ... -luuid  bin/libs5common.a  -laio  -lcurl  -Wl,-Bstatic  -luring  -Wl,-Bdynamic  ../thirdParty/isa-l_crypto/.libs/libisal_crypto.a  bin/libhashtable.a  -lm  -lrt  -ldl  -lrdmacm  -libverbs  -lpthread && cd /home/chen/code/PureFlash/build_deb/pfs && cp /home/chen/code/PureFlash/pfs/pfs_template.conf /home/chen/code/PureFlash/build_deb
/usr/bin/ld: cannot find -luring
```

liburing is so, use `-lcurl  -Wl,-Bdynamic  -luring` in `build.ninja`.
