# ddb

## oceanbase
参考:
- [OceanBase开源，11张图带你了解分布式数据库的核心知识](https://database.51cto.com/art/202106/664792.htm)
- [OceanBase 4.0 解读：降低分布式数据库使用门槛，谈谈我们对小型化的思考](https://open.oceanbase.com/blog/27200153)
- [如何把数据库系统“做小”？OceanBase 4.0主要解决这两个问题](https://c.m.163.com/news/a/HRU22HF80511G549.html)

OceanBase设计为一个Share-Nothing的架构，所以它是没有任何的共享存储结构的。至少需要部署三个以上的Zone，数据在每个Zone都存储一份。OceanBase的整个设计里面没有任何的单点，每个Zone有多个ObServer节点，这就从架构上解决了高可靠高可用的问题.