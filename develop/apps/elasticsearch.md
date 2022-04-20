# Elasticsearch

- [安装文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/install-elasticsearch.html)
- [2万字详解，吃透 ES](https://www.tuicool.com/articles/NFjEfeZ)

## Error
### Cannot allocate memory
调小启动内存
```sh
# vim /etc/elasticsearch/jvm.options

#-Xms2g
#-Xmx2g
-Xms512m
-Xmx512m
```

### 验证elasticsearch是否已运行
`curl -X GET "localhost:9200/?pretty"`