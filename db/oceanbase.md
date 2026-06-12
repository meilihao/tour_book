# oceanbase

## 源码
ref:
- [大规模分布式数据库是如何实现的 -- 读《OceanBase 数据库源码解析》](https://zhuanlan.zhihu.com/p/655202941)
- [《OceanBase数据库源码解析》面市|社区月报 2023.7](https://open.oceanbase.com/blog/5071467520)
- [万字解析：从 OceanBase 源码剖析 paxos 选举原理](https://zhuanlan.zhihu.com/p/630468476)
- [Oceanbase PaxosStore 源码阅读](https://zhuanlan.zhihu.com/p/395197545)
- [**一文详解，什么是单机分布式一体化？**](https://www.modb.pro/db/623598)
- [一文讲透 OceanBase 单机版：架构介绍、部署流程、性能测试、MySQL对比、资源配置等等](https://open.oceanbase.com/blog/11260892737)

    主备
- [从0到1 OceanBase原生分布式数据库内核实战进阶版](https://obcommunity-private-oss.oceanbase.com/prod/blog/2023-09/%E4%BB%8E0%E5%88%B01%20OceanBase%E5%8E%9F%E7%94%9F%E5%88%86%E5%B8%83%E5%BC%8F%E6%95%B0%E6%8D%AE%E5%BA%93%E5%86%85%E6%A0%B8%E5%AE%9E%E6%88%98%E8%BF%9B%E9%98%B6%E7%89%88.pdf)

    对应的代码较旧, 不是V4.x
- [成为ob贡献者(09):翻译PALF2:如何证明采用了PALF设计就是安全的](https://open.oceanbase.com/blog/16292512341)

  ps:
  1. [多副本日志同步](https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000005683234)

    OceanBase 数据库 V4.0.0 版本参考文件系统，将日志服务抽象为 "Paxos Backed Append Only Log File System"，简称 Palf

## build
ref:
- [OceanBase Developer Guide](https://oceanbase.github.io/oceanbase/zh/toolchain/)
- [oceanbase-4.4.2_CE_BP1编译test_physical_plan_ctx_serialize_compat失败](https://ask.oceanbase.com/t/topic/35644769)

env:
- rocky 9.7

注意:
1. 建议使用dep_create.sh支持的os, 经测试rocky 9, 能成功编译observer(`bash build.sh debug --init --make`), 但编译部分test case会报错, 比如test_physical_plan_ctx_serialize_compat, test_expr_serialize_compat

    dep_create.sh支持的almalinux 8也遇到同样问题

```bash
$ -- 安装依赖
$ dnf install git wget rpm* cpio make glibc-devel glibc-headers binutils m4 libtool libaio python3 # 官方文档缺python3, 如有需要执行`ln -s /usr/bin/python3 /usr/bin/python`(rocky 9不需要, almalinux 8需要)
$ -- fix build
$ vim deps/init/dep_create.sh
...
      rocky)
        version_ge "9.0" && compat_centos9 && return
        version_ge "8.0" && compat_centos8 && return

...
$ bash build.sh debug --init --make -j 8 # 限制编译并发防止耗尽内存, 链接耗尽内存用`-j 1`, 内存16g编译比较勉强建议更高
$ cd build_debug/unittest
$ make -j4 # 构建所有test case
```

查看源码:
```bash
$ sudo dnf install epel-release # for install gtest-devel cmake, only for 查看源码
$ -- 按照官方`OceanBase Developer Guide`安装依赖
$ sudo apt install git wget cmake make build-essential binutils m4 libgtest-dev
$ vim .clangd # 可以用locate或从build_debug/compile_commands.json中查找缺失的头文件
CompileFlags:
  Add:
    [
      "-std=gnu++17",
      "-Wno-reserved-user-defined-literal",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/devtools/include/c++/12",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/devtools/include/c++/12/x86_64-redhat-linux",
      "-I/home/chen/test/oceanbase",
      "-I/home/chen/test/oceanbase/src",
      "-I/home/chen/test/oceanbase/src/plugin/include",
      "-I/home/chen/test/oceanbase/deps/oblib/src",
      "-I/home/chen/test/oceanbase/deps/easy/src",
      "-I/home/chen/test/oceanbase/deps/easy/src/include",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/deps/devel/include",
      "-I/home/chen/test/oceanbase/deps/oblib/src/common",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/deps/devel/include/mariadb",
      "-I/home/chen/test/oceanbase/src/objit/include",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/deps/devel/include/oss_c_sdk",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/deps/devel/include/apr-1",
      "-I/home/chen/test/oceanbase/deps/3rd/usr/local/oceanbase/deps/devel/include/mxml",
    ]
```

执行test case:
```bash
$ cd build_debug
$ make test_ob_election
$ cd unittest/logservice
$ ./test_ob_election # 执行单个test case, 该case有log见test_ob_election.log
```

## FAQ
### 版本CE BP HF区别
ref:
- [OceanBase产品命名规则](https://www.modb.pro/db/1697053342350528512)
- [OceanBase 社区版产品规划及研发进展](https://www.modb.pro/db/1691809021846179840)

OceanBase 社区版发布节奏为每2年一个大版本 release，每3个月一次 feature 版本，每个月一个 bug fix 版本( bp 版本):
- 大版本发布即为架构发生升级， 版本升级类似 MySQL 5.7 升级到 MySQL 8.0, 需要做数据迁移才能完成升级.
- feature 版本即为发布了众多 feature 或大 feature , 本地手动冷升级(本地重启)或者通过 OCP 热升级(不停服务).
- bp 版本即为纯 bug fix 版本, 版本升级直接替换 binary 即可, 可以使用 ODP 升级或使用 OCP 热升级.

具体区别:
- CE, Community Edition即社区版
- bp, 纯 bug fix
- HF, 第X个Bugfix版本的第Y个Hotfix

### 主备
[`不支持原主租户降备后接入成为新主租户的备租户`](https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000001574395)
    要重新做主备???

### docker compose 报`oceanbase dir_scan: failed to make directory /root/demo/etc, because File exists`
删除容器后重新创建

### docker compose 报`No such deploy: demo`
MODE=slim时, 不知为什么没法创建cluster `demo`(`/root/boot/start.sh`的`fastboot()`), 测试了`oceanbase/oceanbase-ce:4.3.5-lts`, `oceanbase/oceanbase-ce:latest(4.3.5.2-102020032025070315)`都不可以

解决: 换MODE=MINI

### 执行单个测试
ref:
- [编写以及运行单测](https://oceanbase.github.io/oceanbase/zh/unittest/)