# list
## 选型
- [神仙打架：PG 和 MySQL 到底哪个更好用？](https://www.tuicool.com/articles/AFJ3YnR)
- [MongoDB复制集原理: 复制集容忍失效数](https://developer.aliyun.com/article/64)

## 原则
参考:
- [亿联银行核心银行系统 TiDB 数据库实践之路](https://www.chainnews.com/articles/164828284690.htm)

数据库的建设思路分为`一上一下`:
- 一上

	将数据库的一些处理工作上升到应用层面去解决，**在数据库层面禁用存储过程、触发器、视图等功能**，让数据库变得轻量化和简单化. 这样做有两个好处：
	1. 减少数据库的压力，让数据库变得更好地满足业务的需求
	1. 让数据库应用层跟数据库层技术进行解耦，让上层的应用不依赖于底层的技术

- 一下

	将数据库的一些分库分表的功能给下放到底层分布式数据库层面来解决，减少一些研发的成本和工作量，让开发者更专注于代码功能的实现，加快项目的上线效率.

## 其他
- [记一次微信数据库解密过程](https://www.freebuf.com/articles/endpoint/195107.html), 也可参考[ppwwyyxx/wechat-dump(**推荐**)](https://github.com/ppwwyyxx/wechat-dump)

## next db
### OceanBase
- [蚂蚁金服庆涛：OceanBase支撑2135亿成交额背后的技术原理](https://blog.51cto.com/u_14164343/2344929)