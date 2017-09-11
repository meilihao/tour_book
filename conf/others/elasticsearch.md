# Elasticsearch

- [安装文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/install-elasticsearch.html)

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