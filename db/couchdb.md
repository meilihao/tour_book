# CouchDB
参考:
- [CouchDB 教程](https://www.w3cschool.cn/couchdb/)
- [CouchDB 让人头痛的十大问题](https://www.linuxidc.com/Linux/2012-02/54134.htm)

**强烈推荐不使用它, 坑很深. 个人体验是小众且应用范围窄, 开发麻烦, 查询麻烦等, 总之是麻烦比解决方法多**, 推荐使用tidb cluster代替它:
1. doc频繁变动, couchdb又会保存历史version, 容易耗光存储

存储格式：JSON : **没有表的概念，数据直接以文档的形式存储在数据库中，每个数据库是一个独立的文档集合**, 因此其数据库就类似于sql的table.
查询语言：JavaScript
存储引擎: B-Tree, 复杂度O(n)=log(N). 一个数据库是一个b-tree, 一个view也是b-tree.
支持MVCC, 用`_rev`实现.

> 每个文档都有一个全局惟一的标识符（ID）以及一个修订版本号（revision number）.
> couchdb使用了append-only的文件(仅有追加写)
> 内部数据库以"_"打头, 比如_config 为系统配置数据库，管理员配置也在其中; _users 为用户数据库(authentication)，默认匿名用户可以创建用户.

**注意**：CouchDB 不会彻底删除指定文档，而是会留下一个文档的基本信息，称之为“tombstone”（墓碑），设置的目的是为了实现数据库集群同步复制, 这个问题会导致创建与之前已删除的同名id的doc很麻烦(见`创建doc`), 因此不推荐使用重复doc id.

内部字段解释:
- _id : 全局惟一的标识符，用来惟一标识一个文档
- _rev : 修订版本号(N-后缀, N是修订次数, 后缀是doc的md5值), 用来实现多版本并发控制（Multiversion concurrency control，MVVC）
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

    CouchDB是一个`crash-only`的系统，你可以在任何时候停掉CouchDB并能保证数据的一致性.
- 最终一致性
    
    CouchDB保证最终一致性，使其能够同时提供可用性和分割容忍.
- 离线支持

    CoucbDB能够同步复制到可能会离线的终端设备（比如智能手机），同时当设置再次在线时处理数据同步. CouchDB内置了一个的叫做Futon的通过web访问的管理接口.

冲突处理:
在两个复制的couchdb中修改了同一个doc, 复制时根据一定的规则(由更长修订历史的胜出,相同时由_rev字典序大的胜出)比较后, 胜出的为最新版本, 失败的为最新版本的上一个版本, 复制完成后复制双方会拥有相同的数据, couchdb会把相应的doc标记成冲突状态`"_conflicts": true`. couchdb会将冲突留给应用程序去解决, 而解决一个冲突的通用操作的是首先合并数据到其中一个文档，然后删除旧的数据.

**注意**:
1. 使用自己的uuid来create doc而不是由couchdb生成, 避免重试操作时创建多个uuid不同的doc.

## RESTful api
> 单个操作用PUT, 多个用POST.
> couchdb Futon(自带的web ui): `http://localhost:5984/_utils`.
> CouchDB doc支持附加文件, 会自动对文件进行 base64 解码.
> couchdb doc中以`_`作为前缀的顶层字段是由 CouchDB 保留使用的，如`_id和_rev`.
> CouchDB 目前只支持一种角色(`系统管理员`), 拥有所有权限.
> [CouchDB 2.2 开始支持`POST /{db}/_design/{ddoc}/_view/{view}/queries`](https://docs.couchdb.org/en/stable/api/ddoc/views.html)

- 获取数据库列表

    `curl -X GET http://127.0.0.1:5984/_all_dbs -u admin:admin`
- 创建名称为testdb的数据库

    `curl -X PUT http://127.0.0.1:5984/testdb resp: {"ok":true}`
- 获取数据库信息

    `curl -X GET http://127.0.0.1:5984/testdb resp: {...}`
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

    > 1925a2a284289df9b55b390525001ca1以前存在过(包括已删除)则会变成更新, 因此此时必须提供`_id`和`_rev`. 比如使用`curl -X GET "$host/<db>/<id>?revs=true&open_revs=all" -H "Accept: application/json"`, 可参考[这里](http://garmoncheg.blogspot.com/2013/11/couchdb-restoring-deletedupdated.html).
    
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

    当更新完成之后，返回 HTTP 状态代码 201 ；否则返回 HTTP 状态代码 409，表示有版本冲突.

    1. 更新单个doc

     `curl -X PUT http://127.0.0.1:5984/testdb/1925a2a284289df9b55b390525001ca1 -d '{"_rev":"2-8631301e81523fb6f58fa99b33f2731f", "id":1,"name":"mike"}'  -u admin:admin resp: {"ok":true,"id":"1925a2a284289df9b55b390525001ca1","rev":"3-fcaf069a73856d697af04d750228ba20"}`

     1. 更新多个doc

     curl -X POST http://localhost:5984/testdb/_bulk_docs -H 'Content-Type: application/json' -u admin:admin -d '{
        "docs": [
            {"_id":"s1","_rev":"10-85d31c96537e6c4e5db128b2260af897","desc":"bobo2","age":56.32},
            {"_id":"s2","_rev":"9-49093809c74d86493969f441efb4a1b7","desc":"小红2","age":19.36}
        ]
    }' resp: [{"ok":true,"id":"s1","rev":"10-85d31c96537e6c4e5db128b2260af897"},{"id":"s2","error":"conflict","reason":"Document update conflict."}]
- 附件

    curl -X PUT http://localhost:5984/testdb/1925a2a284289df9b55b390525001ca1/artwork.jpg?rev=3-fcaf069a73856d697af04d750228ba20 --data-binary.jpg -H 'Content-Type: image/jpg'

    附件会出现在doc的`{...,"_attachments":{"artwork.jpg":{...}}`中, 默认仅显示附件的元数据, 获取doc的url请求中加`?attachments=true`会返回包含附件数据(base64编码)的内容.

    curl http://localhost:5984/testdb/1925a2a284289df9b55b390525001ca1/artwork.jpg可看到附件.

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

    **couchdb v2开始, source和target必须是完整的url形式.**

    查看复制进度: `http://username:password@localhost:5984/_active_tasks`的`progress(0~100, 100是完成)`, 新版couchdb用其中的`changes_pending("0/没有该属性"是完成)`判断.
    复制时, **源和目的必须都必须已存在**， 或使用`"create_target":true`选项, 不存在目的时自动创建; 复制仅针对创建复制时的状态, 复制开始后的变更(创建/修改/删除)都不会被复制.

    复制内容包括新增/已修改/已删除的doc. couchdb中每个数据库有一个序列号, 每次变更都会加一, 且会记录哪个变更对应哪个序列号. couchdb支持增量复制.
    couchdb可使用`"continous":true`进行持续复制. 从 CouchDB 1.1.0 开始，可以通过在`_replicator`数据库中插入文档来定义在服务器重启后无需执行任何操作的永久连续复制, 非`_replicator`中定义的复制任务会在couchdb重启后消失. 且对于临时复制，当作业完成时无法查询它们的状态. 删除`_replicator`的定义会导致复制停止.

    > `"continous":true`原理: couchdb使用_changes api, 在任何新文档进入source时自动复制到目标, 不是实时, 有一个复杂的算法来确定复制以获得最佳性能的理想时刻.

    如果想要双向复制，只需要触发两个复制, 源和目标交换.

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
- 获取couchdb配置信息

    ```sh
    $ http://localhost:5984/_config
    $ http://localhost:5984/_config/log
    ```

## 设计文档
设计文档是一类特殊的文档，其ID必须以`_design/`开头. 设计文档的存在是使用CouchDB开发Web应用的基础. 在CouchDB中，一个Web应用是与一个设计文档相对应的, 但一个couchdb数据库允许有多个设计文档.

在设计文档中可以包含一些特殊的字段，其中包括：
- views 包含永久的视图定义

    CouchDB中一般将视图称为MapReduce视图，一个MapReduce视图由两个JavaScript函数组成，一个是Map函数，这个函数必须定义；另一个是Reduce函数，这个是可选的，根据程序的需要可以进行选择性定义.

    view第一次运行时需要遍历所有文档, 结果保持在b-tree中, 之后只有doc有变动才调用该view更新b-tree.
- shows 包含把文档转换成任意格式(比如非JSON)进行输出的方法
- lists 包含把视图运行结果转换成非JSON格式的方法
- validate_doc_update 包含验证文档更新是否有效的方法, 通过函数抛出异常`throw({...})`来取消更新

    couchdb推荐使用内建函数`toJSON`来比较, 因为原生js`[] == []`返回false等行为.

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
`emit()`以kv形式返回. **每个视图的定义中至少需要提供 Map 方法，Reduce 方法是可选的**. 一个doc允许有多个emit()输出, 比如`emit()`在循环中的时候.

couchdb内置reduce函数:
- `_count` : 返回map结果集中值的个数, 与查询参数`include_docs=true`冲突
- `_sum` : 返回map结果集中数值的求和结果
- `_stats` : 返回map结果集中数值的统计结果

运行视图时的可选参数(参数, 说明), 以下仅是部分, 更多的话查看[这里](https://docs.couchdb.org/en/latest/api/ddoc/views.html#db-design-design-doc-view-view-name):
- key : 限定结果中只包含键为该参数值的记录
- startkey : 限定结果中只包含键大于或等于该参数值的记录
- endkey : 限定结果中只包含键小于或等于该参数值的记录

    比如获取all docs except design可使用`/xxx/_all_docs?include_docs=true&endkey="_design"`
- limit : 限定结果中包含的记录的数目
- descending : 指定结果中记录是否按照降序排列

    举例:
    - `startkey=1&descending`: 读取范围是[文件开始,startkey], 且从startkey开始倒读
    - `endkey=1&descending`: 读取范围是[endkey,文件结尾], 且到endkey为止
- skip : 指定结果中需要跳过的记录数目
- group : 指定是否对键进行分组
- reduce : 指定reduce=false可以只返回 Map 方法的运行结果. 获取总条数时需要`reduce=true`
- group_level : 当key为array时, 允许使用group_level指定使用array的前n个元素作为新key来执行mapreduce.
- include_docs=true : 返回doc

> startkey + limit + skip 可实现分页

**注意**:
1. view会对key的结果进行排序
1. 第一次view执行时会在所有doc上执行一次并创建view索引保存到b-tree, 之后仅对变更过的doc执行一次该view, 重新生成key和value.

couchdb查询view的过程:
1. 从最顶上开始读取, 如果startkey存在, 就会从startkey的位置开始读取
1. 一次返回一行, 直至结尾, 如果endkey存在, 就直到endkey指定的位置为止

### key
> v3.0.0 curl的`?key=...`有问题, 新建的view无法匹配到结果(有数据).

**key必须是json格式**.

参考:
- [Couchbase——查询View（详细版）](https://blog.csdn.net/EntropyArrow/article/details/40858225)

Couchbase Server 支持诸多筛选方式. Key筛选是在view结果产生之后（包括redution函数），也在产生结果排序之后才开始进行.
给key筛选引擎使用的必须是JSON值. 比如说，指定一个单独的key，如果这是个字符串，那么必须包含在引号内, 比如`http://127.0.0.1:5984/testdb/_design/t/_view/all?reduce=false&key="1001"`

**curl参数key匹配的是emit()返回的key, 且两者的格式必须一致**.

#### map
Map方法的参数只有一个，就是当前的文档对象. Map方法的实现需要根据文档对象的内容，确定是否要输出结果. 如果需要输出的话，可以通过emit来完成.

emit方法有两个参数, 分别是key和value, 分别表示输出结果的键和值, 且emit可多次调用.

#### reduce
Reduce方法的参数有三个：key、values和rereduce，分别表示键、值和是否是rereduce. 由于rereduce情况的存在，Reduce 方法一般需要处理两种情况：
1. 传入的参数rereduce的值为false : 表明Reduce方法的输入是 Map方法输出的中间结果(kv).
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

## 变更通知

```bash
$ curl -X GET http://localhost/db/_changes
```

返回结果:
- seq : 数据库变更时, couchdb生成的变更序列号
- id : doc._id
- changes : 数组, 默认是doc的_rev, 但可能包含文档冲突及其他信息

支持参数:
- style=all : changes数组包含更多的版本和冲突信息
- since=1 : since起始变更序列号
- feed=continous : 长连接, 连续获取变更
- filter=design_name/filtername : 支持设计文档定义的`filters`过滤函数

## client demo
```python
import couchdb
s = couchdb.Server('http://192.168.0.71:5984')
s.resource.credentials = ('admin', 'password')
print s.version()

db = couchdb.Database('http://192.168.0.71:5984/test')
db.resource.credentials = ('admin', 'password')
print db.name
```

## FAQ
### `emit(doc.phoneNumber, doc.billSeconds);`和`emit([doc.phoneNumber], doc.billSeconds);`的区别
`emit(doc.phoneNumber, doc.billSeconds)`　：　返回的key是`"1000"`
`emit([doc.phoneNumber], doc.billSeconds)`:   返回的key是`["1001"]`

### log
`$ couchdb -c`显示配置位置及加载顺序, 再从配置中找到log的位置.

couchdb log中不含request url, 但可在`journalctl`中看到该url.

### update报409
`_rev`已过时，因为每次文档更新都会对其进行更改.

### 更新丢失
couchdb使用`_rev`更新机制(乐观锁).

a, b同时更新一个文档, 假设a先提交了, b提交报409引发重试, 但未merge a的修改, 此时b提交就会导致a的数据丢失.

### neokylin v10编译couchdb 3.1.1
**编译couchdb时必须在英文环境下即`LANG= make release`, 否则`make release`时会报错**.

#### 1. 安装js-devel 1.8.5
方法1, 自己构建(**推荐**):
```bash
# git clone --depth 1 git@github.com:apache/couchdb-pkg.git
# cd couchdb-pkg
# make couch-js-rpms # 可能需要安装一些依赖 `yum install readline-devel nspr-devel ncurses-devel`
```

> ubuntu打包命令: make couch-js-rpms PLATFORM=ubuntu-xenial # 其他ubuntu版本替换`ubuntu-xenial`即可

方法2, 使用其他源的js-devel:
1. 下载rpm并依次安装
- https://centos.pkgs.org/7/centos-aarch64/js-1.8.5-20.el7.aarch64.rpm.html
- https://centos.pkgs.org/7/centos-x86_64/js-devel-1.8.5-20.el7.x86_64.rpm.html

#### 2. 构建rpm couchdb 3.1.1
参考:
- [如何制作一个标准的 RPM 包](https://gohalo.me/post/linux-create-rpm-package.html)

````bash
# yum install libicu-devel
# # 1. 参考Makefile里的target all, 先根据bin/detect-target.sh获取构建target, neokylin v10选择centos-7
# # 2. 参考README.md里的"rpms or debs from a release tarball", 确定构建参数
# LANG= make copy-couch centos-7 COUCHTARBALL=/root/couchdb-3.1.1.tar.gz PLATFORM=centos-7
# # 上面的命令会输出执行过程并报错, 但rpm构建环境即rpmbuild目录已成功创建, 可根据报错修改rpmbuild/SPECS/couchdb.spec, 再根据上面命令输出的执行过程找到最后的rpmbuild命令, 切换到rpmbuild目录重新执行该命令即可.
````

> ubuntu打包: LANG= make copy-couch ubuntu-xenial COUCHTARBALL=/root/couchdb-3.1.1.tar.gz PLATFORM=ubuntu-xenial

couchdb.spec修改内容:
1. BuildArch 追加 aarch64
1. ExclusiveArch 追加 aarch64
1. BuildRequires: esl-erlang 替换为

    ```conf
    BuildRequires: erlang-asn1
    BuildRequires: erlang-erts
    BuildRequires: erlang-eunit
    BuildRequires: erlang-os_mon
    BuildRequires: erlang-xmerl
    BuildRequires: erlang-erl_interface
    BuildRequires: erlang-reltool
    ```

    > 参考了couchdb.spec里的`0%{?suse_version}`和[官方的`installation-from-source`](https://docs.couchdb.org/en/stable/install/unix.html#installation-from-source)

构建出的rpm是couchdb-3.1.1-1.ky10.ky10.aarch64.rpm.

注意: couchdb-3.1.1没有设置Admin帐号启动会自动退出, 可在`/var/log/couchdb/couchdb.log`里看到该提示, 此时修改`/opt/couchdb/etc/local.ini中的[admin]`取消admin帐号的注释即可.

访问`http://127.0.0.1:5984/_utils#setup`即可.

> couchdb 3.1.1 配置的路径是`/opt/couchdb/etc/`, 不再是`/etc/couchdb`

> neokylin v10默认防火墙是开着的(`firewall-cmd --state`), 因此其他机子访问couchdb 5984时应先打开端口(`firewall-cmd --zone=public --add-port=5984/tcp --permanent && firewall-cmd --reload`).

#### 安装couchdb-3.1.1-1.ky10.ky10.aarch64.rpm 报"file /usr/lib/.build-id/xx/xxx..xxx from install of couchdb-3.1.1-1.ky10.ky10.aarch64.rpm conflicts with file from package erlang-crypto-22.3.4.1-1.ky10.aarch64"
rpmbuild命令追加参数`--define "_build_id_links none"`禁止生成rpm build-id.


### No Admin Account Found, aborting startup
```bash
# vim /etc/couchdb/local.ini
;admin = mysecretpassword # 去掉该注释即可
```

### couchdb listen 0.0.0.0
```
# vim /etc/couchdb/local.ini
bind_address = 0.0.0.0
```

### couchdb partition must not be empty
couchdb 3.1.1支持partition, 新建database时选择`non-partition`即可.

### chttp_auth_cache changes listener diea because the _users database dees not exist
访问`http://127.0.0.1:5984/_utils#setup`, 按照向导设置即可.


### couchdb 1.6/3.1 权限变化
1. 1.6允许不设置admin帐号, 但3.1必须设置否之couchdb启动后会退出
1. 1.6 允许no auth访问couchdb, 3.1设置

    ```ini
    [couchdb]
    default_security=everyone
    [chttpd]
    require_valid_user=false
    [couch_httpd_auth]
    require_valid_user=false
    ```

    部分接口允许匿名访问, 比如创建doc, 查询doc等, 但`/_all_dbs`, 创建view仍需admin权限, 可参考[这里](https://docs.couchdb.org/en/3.1.1/intro/security.html#authentication).

    因此, **推荐带auth访问couchdb**.

### 创建doc报`Invalid rev format`
创建doc时, 请求应不包含rev或rev符号couchdb `doc._rev`(`\d+\-[\w]{32}`)格式

### Document update conflict
doc已删除, 更新其views时报错.

原因: doc仍在couchdb中, 只是rev未知, 此时更新才报该错.
解决方法: 重新创建同_id的doc(**此时PUT的json body中不能有`_rev`属性否则变成了更新, 还报该错**), 获取其rev后再次更新即可.

## 驱动
### [go-kivik/kivik](https://github.com/go-kivik/kivik)
1. `rows.TotalRows()`需要在`rows.Next()`后调用才能获取到值
```go
// curl xxx get `{"total_rows":2, ...}`
rows,err:=db.AllDocs()
defer rows.Close()

fmt.Println(rows.TotalRows()) // 0
for rows.Next{
    rows.ScanDoc()
}
fmt.Println(rows.TotalRows()) // 2
```

### go-kivik/kivik view query
```go
// get only key=`[false,"abc"]`
opts = map[string]interface(){
    "reduce": false,
    "startkey": json.RawMessage(`[false,"abc"]`),
    "endkey": json.RawMessage(`[false,"abc"]`),
}
db:=conn.DB(dbName)
db.Query(ctx, dbName, viewName, opts...)
```

### 仅获取未软删除的doc
query: `key=[false]`

view:
```js
var deleted = false;
if ('deleted' in doc) {
    deleted = doc['deleted']
}
emit([deleted],doc);
```

### db清理
参考:
- [`POST /{db}/_compact`](https://docs.couchdb.org/en/main/api/database/compact.html)
- [Compaction](https://docs.couchdb.org/en/main/maintenance/compaction.html)

存储空间耗光时couchdb的表现: 访问couchdb时好时坏.

> 清理操作也需要空间: 它是先生成新文件再删除旧文件的操作.