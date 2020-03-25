# CouchDB
参考:
- [CouchDB 教程](https://www.w3cschool.cn/couchdb/)
- [CouchDB 让人头痛的十大问题](https://www.linuxidc.com/Linux/2012-02/54134.htm)

存储格式：JSON : **没有表的概念，数据直接以文档的形式存储在数据库中，每个数据库是一个独立的文档集合**, 因此其数据库就类似于sql的table.
查询语言：JavaScript

> 每个文档都有一个全局惟一的标识符（ID）以及一个修订版本号（revision number）.

内部字段解释:
- _id : 全局惟一的标识符，用来惟一标识一个文档
- _rev : 修订版本号，用来实现多版本并发控制（Multiversion concurrency control，MVVC）
- _attachments : 内嵌型附件，以 base64 编码的格式作为文档的一个字段保存.

特点:
- MVCC（Multiversion concurrency control）

    CouchDB一个支持多版本控制的系统，此类系统通常支持多个结点写， 而系统会检测到多个系统的写操作之间的冲突并以一定的算法规则予以解决.
- 水平扩展性

    在扩展性方面，CouchDB使用replication去做. CouchDB的设计基于支持双向的复制（同步）和离线操作. 这意味着多个复制能够对同一数据有其自己的拷贝，可以进行修改，之后将这些变更进行同步.
- REST API(Representational State Transfer，简称REST，表述性状态转移）

    所有的数据都有一个唯一的通过HTTP暴露出来的URI. REST使用HTTP方法 POST，GET，PUT和DELETE来操作对应的四个基本 CRUD(Create，Read，Update，Delete）操作来操作所有的资源.
- 数据查询操作

    CouchDB不支持动态查询，你**必须为你的每一个查询模式建立相应的视图**，并在此视图的基础上进行查询。. 

    **视图是CouchDB中文档的呈现方式，在CouchDB 中保存的是视图的定义**.

    CouchDB 中有两种视图：永久视图和临时视图. 永久视图保存在设计文档的views字段中. 如果需要修改永久视图的定义，只需要通过文档 REST API 来修改设计文档即可. 临时视图是通过发送 POST 请求到 URL/dbName/_temp_view 来执行的. 在POST请求中需要包含视图的定义. **一般来说，临时视图只在开发测试中使用，因为它是即时生成的，性能比较差； 永久视图的运行结果可以被 CouchDB 缓存，因此一般用在生产环境中**.
- 原子性

    支持针对行的原子性修改（concurrent modifications of single documents），但不支持更多的复杂事务操作.
- 数据可靠性

    CouchDB是一个”crash-only”的系统，你可以在任何时候停掉CouchDB并能保证数据的一致性.
- 最终一致性
    
    CouchDB保证最终一致性，使其能够同时提供可用性和分割容忍.
- 离线支持

    CoucbDB能够同步复制到可能会离线的终端设备（比如智能手机），同时当设置再次在线时处理数据同步. CouchDB内置了一个的叫做Futon的通过web访问的管理接口.

## RESTful api
> 单个操作用PUT, 多个用POST.

- 获取数据库列表

    `curl -X GET http://127.0.0.1:5984/_all_dbs -u admin:admin`
- 创建名称为testdb的数据库

    `curl -X PUT http://127.0.0.1:5984/testdb resp: {"ok":true}`
- 删除名称为testdb的数据库

    `curl -X DELETE http://127.0.0.1:5984/testdb resp: {"ok":true}`
- 获取uuid

    1. 获取单个
    ```sh
    curl -X GET http://127.0.0.1:5984/_uuids  -u admin:admin
    {"uuids":["43c23b658fc206f03826ad447200666e"]}
    ```

    1. 获取多个
    ```sh
    curl -X GET 'http://127.0.0.1:5984/_uuids?count=3' -u admin:admin                                                                         20:28:37
    {"uuids":["43c23b658fc206f03826ad4472007240","43c23b658fc206f03826ad447200814b","43c23b658fc206f03826ad4472008d3c"]}
    ```
- 创建doc

    当文档指定`_id`时, 被创建的文档会使用指定的`_id`.

    1. 插入单个doc: 通过PUT请求访问 URL/dbName/doc_id 在testdb中创建ID为doc_id的文档(couchdb 内部用`_id`表示doc_id, 因此不会覆盖文档中的`id`属性); 通过POST请求访问 URL/dbName 也可以创建新文档，不过是由 CouchDB 来生成文档的ID.

    `curl -X PUT http://127.0.0.1:5984/testdb/1925a2a284289df9b55b390525001ca1 -d '{"id":1,"name":"mike"}'  -u admin:admin resp: {"ok":true,"id":"1925a2a284289df9b55b390525001ca1","rev":"1-0c1f72feabb29905ed205d25fbcbf3b3"}`
    
    1. 插入多个doc, 此时必须用POST:

    ```curl
    curl -X POST http://localhost:5984/testdb/_bulk_docs -H 'Content-Type: application/json' -u admin:admin -d '{
        "docs": [
            {"_id":"s1","desc":"bobo","age":15},
            {"_id":"s2","desc":"小红","age":19}
        ]
    }' resp: [{"ok":true,"id":"s1","rev":"5-64964a8e3910a7f17b8bd1b2942811bf"},{"ok":true,"id":"s2","rev":"5-337aef79b97ef2a413f9d2c3d6560917"}]
    ``` 
- 删除

    1. 删除单个doc, 且必须指定版本, 否则会报错

    `curl -X DELETE 'http://127.0.0.1:5984/testdb/2?rev=1-0855a5ac3f1a5dc0d3800d011a5903b9' -u admin:admin  resp: {"ok":true,"id":"2","rev":"2-69a219d56106338377076eef41160636"}`

    1. 删除多个doc
    `curl -X POST http://localhost:5984/testdb/_bulk_docs -H 'Content-Type: application/json' -u admin:admin -d '{
        "docs": [
            {
                "_id": "s1",
                "_rev": "7-d7adfe3c65182edaca7ed84a8b0534b9",
                "_deleted": true
            },
            {
                "_id": "s2",
                "_rev": "7-cf25ed0514fb6eaf261400be9965a536",
                "_deleted": true
            }
        ]
    }' resp: [{"ok":true,"id":"s1","rev":"8-039d129bb5679fa873ca743b682318ea"},{"ok":true,"id":"s2","rev":"8-971697fd6dbff4d7ade911ba2f549852"}]`
- 更新

    CouchDB是当指定_id的文档不存在的时候，就会插入，否则就是更新. 如果更新没有指定_rev参数或参数值不正确，则更新就不会成功，报`Document update conflict`.

    **CouchDB不支持对JSON文档进行部分更新，必须是全更新. 也就是说不能添加或者删除字段，也不能仅更新某些字段的值. 更新文档的时候_rev 参数必须是最新获取的，因为任何修改操作都会引起_rev`值变化**.

    1. 更新单个doc

     `curl -X PUT http://127.0.0.1:5984/testdb/1925a2a284289df9b55b390525001ca1 -d '{"_rev":"2-8631301e81523fb6f58fa99b33f2731f", "id":1,"name":"mike"}'  -u admin:admin resp: {"ok":true,"id":"1925a2a284289df9b55b390525001ca1","rev":"3-fcaf069a73856d697af04d750228ba20"}`

     1. 更新多个doc

     curl -X POST http://localhost:5984/testdb/_bulk_docs -H 'Content-Type: application/json' -u admin:admin -d '{
        "docs": [
            {"_id":"s1","_rev":"10-85d31c96537e6c4e5db128b2260af897","desc":"bobo2","age":56.32},
            {"_id":"s2","_rev":"9-49093809c74d86493969f441efb4a1b7","desc":"小红2","age":19.36}
        ]
    }' resp: [{"ok":true,"id":"s1","rev":"10-85d31c96537e6c4e5db128b2260af897"},{"id":"s2","error":"conflict","reason":"Document update conflict."}]

- 查询

    1. 单条数据查询. 通过GET请求访问 URL/dbName/doc_id 可以获取名称为dbName的数据库中ID为doc_id文档的内容. 文档的内容是一个JSON对象，其中以“ _ ”作为前缀的顶层字段是由CouchDB保留使用的，如_id和_rev.

    `curl -X GET http://127.0.0.1:5984/testdb/1925a2a284289df9b55b390525001ca1  -u admin:admin  resp: {"_id":"1925a2a284289df9b55b390525001ca1","_rev":"1-0c1f72feabb29905ed205d25fbcbf3b3","id":1,"name":"mike"}`

    1. 多条数据查询
        1. _all_docs : 获取全部doc

        ```sh
        $ curl -X GET http://127.0.0.1:5984/testdb/_all_docs  -u admin:admin
        {"total_rows":2,"offset":0,"rows":[
        {"id":"s1","key":"s1","value":{"rev":"11-3662e4e7e79a00c4d1c60b7e5c7e8f7e"}},
        {"id":"s2","key":"s2","value":{"rev":"10-dd1343dfc3fb10ae1069f4311561bec5"}}
        ]}
        ```

        此时返回结果中只包含文档的_id和_rev字段，如果需要获取包含文档所有字段的json，则需在URL后加上请求参数`include_docs=true`.

        ```sh
        curl -X GET 'http://127.0.0.1:5984/testdb/_all_docs?include_docs=true'  -u admin:admin
        {"total_rows":2,"offset":0,"rows":[
        {"id":"s1","key":"s1","value":{"rev":"11-3662e4e7e79a00c4d1c60b7e5c7e8f7e"},"doc":{"_id":"s1","_rev":"11-3662e4e7e79a00c4d1c60b7e5c7e8f7e","desc":"bobo2","age":56.32}},
        {"id":"s2","key":"s2","value":{"rev":"10-dd1343dfc3fb10ae1069f4311561bec5"},"doc":{"_id":"s2","_rev":"10-dd1343dfc3fb10ae1069f4311561bec5","desc":"小红2","age":19.36}}
        ]}
        ```

        1. _find, [selector支持的操作符](http://docs.couchdb.org/en/latest/api/database/find.html#find-selectors)

        curl -X POST --url http://localhost:5984/testdb/_find -H 'Content-Type: application/json' -u admin:admin -d '{
            "selector": {
                "age": {
                    "$gte": 5,
                    "$lte": 20
                }
            },
            "fields": ["_id", "_rev", "age", "desc"],
            "execution_stats": true
        }'

        resp:
        ```json
        {
            "docs": [
                {
                    "_id": "s2",
                    "_rev": "10-dd1343dfc3fb10ae1069f4311561bec5",
                    "age": 19.36,
                    "desc": "小红2"
                }
            ],
            "bookmark": "g1AAAAA0eJzLYWBgYMpgSmHgKy5JLCrJTq2MT8lPzkzJBYkXG4FkOGAyULEsAGqkDgM",
            "execution_stats": {
                "total_keys_examined": 0,
                "total_docs_examined": 3,
                "total_quorum_docs_examined": 0,
                "results_returned": 1,
                "execution_time_ms": 1.325
            },
            "warning": "No matching index found, create an index to optimize query time."
        }
        ```

        1. _bukl_get

        ```sh
        curl -X POST http://localhost:5984/testdb/_bulk_get -H 'Content-Type: application/json' -u admin:admin -d '{
        "docs": [
            {
                "id": "s1",
                "rev": "7-d7adfe3c65182edaca7ed84a8b0534b9"
            },
            {
                "id": "s3"
            },
            {
                "id": "s41"
            }
        ]
        }'
        {
            "results": [
                {
                    "id": "s1",
                    "docs": [
                        {
                            "ok": {
                                "_id": "s1",
                                "_rev": "7-d7adfe3c65182edaca7ed84a8b0534b9",
                                "desc": "bobo",
                                "age": 15
                            }
                        }
                    ]
                },
                {
                    "id": "s3",
                    "docs": [
                        {
                            "ok": {
                                "_id": "s3",
                                "_rev": "4-7b09720f025ae87ad8d48cdd9eb8807f",
                                "_deleted": true
                            }
                        }
                    ]
                },
                {
                    "id": "s41",
                    "docs": [
                        {
                            "error": {
                                "id": "s41",
                                "rev": "undefined",
                                "error": "not_found",
                                "reason": "missing"
                            }
                        }
                    ]
                }
            ]
        }
        ```
        
        1. 通过视图

            ```sh
            $ curl -X POST http://localhost:5984/testdb -H 'Content-Type: application/json' -u admin:admin -d '{ "_id": "_design/example", "language": "javascript", "views": { "getdata": { "map": "function(doc){ if(doc.id>0){emit(doc.id, doc.name)}}"}}}' # 创建视图
            {"ok":true,"id":"_design/example","rev":"1-ef2afac35cc971bd4b623af76c2b8641"}
            $ curl -X GET http://127.0.0.1:5984/testdb/_design/example/_view/getdata -u admin:admin
            {"total_rows":1,"offset":0,"rows":[
            {"id":"1925a2a284289df9b55b390525001ca1","key":1,"value":"mike"}
            ]}
            $ curl -X DELETE 'http://localhost:5984/testdb/_design/example?rev=4-5d38c3d177c10cc93880e663a98f61ae' -H 'Content-Type: application/json' -u admin:admin # 此时需先删除上面的`_design/example`
            {"ok":true,"id":"_design/example","rev":"5-8a7a99ed7ed5790b547970e4afc4c303"}
            $ curl -X POST http://localhost:5984/testdb -H 'Content-Type: application/json' -u admin:admin -d '{ "_id": "_design/example", "language": "javascript", "views": { "getall": {"map": "function(doc){ emit(doc.id, doc.name)}"}}}' # 创建视图
            {"ok":true,"id":"_design/example","rev":"4-5d38c3d177c10cc93880e663a98f61ae"}
            $ curl -X GET http://127.0.0.1:5984/testdb/_design/example/_view/getall -u admin:admin
            {"total_rows":3,"offset":0,"rows":[
            {"id":"s1","key":null,"value":null},
            {"id":"1925a2a284289df9b55b390525001ca1","key":1,"value":"mike"}
            ]}
            ```

- 数据库复制

    1. 本地复制

    curl -X PUT http://127.0.0.1:5984/albums-replica # 创建接收者
    curl -X POST http://127.0.0.1:5984/_replicate -d '{"source":"albums","target":"albums-replica"}' -H "Content-Type: application/json"

    1. 本地到远端复制

    curl -X POST http://127.0.0.1:5984/_replicate -d '{"source":"albums","target":"http://example.org:5984/albums-replica"}' -H "Content-Type:application/json"

    > 如果远端服务器有密码，可以采用这种格式：http://username:pass@remotehost:5984/demo

    1. 远端到本地

    curl -X POST http://127.0.0.1:5984/_replicate -d '{"source":"http://example.org:5984/albums-replica","target":"albums"}' -H "Content-Type:application/json"

    1. 远端到远端

    curl -X POST http://127.0.0.1:5984/_replicate -d '{"source":"http://example.org:5984/albums","target":"http://example.org:5984/albums-replica"}' -H "Content-Type: application/json"

## 设计文档
设计文档是一类特殊的文档，其ID必须以`_design/`开头. 设计文档的存在是使用CouchDB开发Web应用的基础. 在CouchDB中，一个Web应用是与一个设计文档相对应的.

在设计文档中可以包含一些特殊的字段，其中包括：
- views 包含永久的视图定义

    CouchDB中一般将视图称为MapReduce视图，一个MapReduce视图由两个JavaScript函数组成，一个是Map函数，这个函数必须定义；另一个是Reduce函数，这个是可选的，根据程序的需要可以进行选择性定义.
- shows 包含把文档转换成非JSON格式的方法
- lists 包含把视图运行结果转换成非JSON格式的方法
- validate_doc_update 包含验证文档更新是否有效的方法

[设计文档结构](https://github.com/mike-zhang/mikeBlogEssays/blob/master/2014/20141007_couchDB%E6%96%87%E6%A1%A3.md)：
```json
{
    "_id" : "_design/${docName}", "_rev" : "${revision}",
    "language" : "javascript",
    "views": {
        ...
    },
    "shows": {
        ...
    },
    "lists": {
        ...
    },
    "_attachments": {
        ...
    }
}
```

### view函数
`omit()`以kv形式返回.

couchdb内置reduce函数:
- `_count` : 返回map结果集中值的个数, 与查询参数`include_docs=true`冲突
- `_sum` : 返回map结果集中数值的求和结果
- `_stats` : 返回map结果集中数值的统计结果

### key
> v3.0.0 curl的`?key=...`有问题, 新建的view无法匹配到结果(有数据).

参考:
- [Couchbase——查询View（详细版）](https://blog.csdn.net/EntropyArrow/article/details/40858225)

Couchbase Server 支持诸多筛选方式. Key筛选是在view结果产生之后（包括redution函数），也在产生结果排序之后才开始进行.
给key筛选引擎使用的必须是JSON值. 比如说，指定一个单独的key，如果这是个字符串，那么必须包含在引号内, 比如`http://127.0.0.1:5984/testdb/_design/t/_view/all?reduce=false&key="1001"`

**curl参数key匹配的是emit()返回的key, 且两者的格式必须一致**.

#### map
Map方法的参数只有一个，就是当前的文档对象. Map方法的实现需要根据文档对象的内容，确定是否要输出结果. 如果需要输出的话，可以通过emit来完成。 emit方法有两个参数，分别是key和value，分别表示输出结果的键和值.

#### reduce
Reduce方法的参数有三个：key、values和rereduce，分别表示键、值和是否是rereduce. 由于rereduce情况的存在，Reduce 方法一般需要处理两种情况：
1. 传入的参数rereduce的值为false : 表明Reduce方法的输入是 Map方法输出的中间结果.
参数key的值是一个数组，对应于中间结果中的每条记录. 该数组的每个元素都是一个包含两个元素的数组，第一个元素是在Map方法中通过emit输出的键（key），第二个元素是记录所在的文档 ID.
参数values的值是一个数组，对应于 Map 方法中通过emit输出的值（value）.

1. 传入的参数rereduce的值为true : 表明Reduce方法的输入是上次Reduce方法的输出.
参数key的值为null, 参数values的值是一个数组，对应于上次Reduce方法的输出结果.

```sh
$ curl -X POST http://localhost:5984/testdb -H 'Content-Type: application/json' -u admin:admin -d '{
   "_id": "_design/jsTest", 
   "language": "javascript",
   "views": {
       "all": {
           "map": "function(doc) {  emit(doc.id, doc.age); }",
           "reduce": "function (key, values, rereduce) {   return sum(values); }"
       }
   }
}'
$ curl http://127.0.0.1:5984/testdb/_design/jsTest/_view/all?reduce=false -u admin:admin # 不使用reduce
{"total_rows":3,"offset":0,"rows":[
{"id":"s1","key":null,"value":56.32},
{"id":"s2","key":null,"value":19.36},
{"id":"1925a2a284289df9b55b390525001ca1","key":1,"value":null}
]}
$ curl 'http://127.0.0.1:5984/testdb/_design/jsTest/_view/all?group=true' -u admin:admin # 有reduce且没有指定`?reduce=false`时会使用reduce. group指按照emit的key作group by. 
{"rows":[
{"key":null,"value":75.68},
{"key":1,"value":0}
]}
```

## FAQ
### `emit(doc.phoneNumber, doc.billSeconds);`和`emit([doc.phoneNumber], doc.billSeconds);`的区别
`emit(doc.phoneNumber, doc.billSeconds)`　：　返回的key是`"1000"`
`emit([doc.phoneNumber], doc.billSeconds)`:   返回的key是`["1001"]`