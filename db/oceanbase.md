# oceanbase

## 源码
ref:
- [大规模分布式数据库是如何实现的 -- 读《OceanBase 数据库源码解析》](https://zhuanlan.zhihu.com/p/655202941)
- [《OceanBase数据库源码解析》面市|社区月报 2023.7](https://open.oceanbase.com/blog/5071467520)
- [万字解析：从 OceanBase 源码剖析 paxos 选举原理](https://zhuanlan.zhihu.com/p/630468476)
- [Oceanbase PaxosStore 源码阅读](https://zhuanlan.zhihu.com/p/395197545)
- [一文讲透 OceanBase 单机版：架构介绍、部署流程、性能测试、MySQL对比、资源配置等等](https://open.oceanbase.com/blog/11260892737)

    主备

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