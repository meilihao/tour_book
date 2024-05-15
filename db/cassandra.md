# Cassandra

## Error
### Out of memory: Kill process
调小启动内存
```sh
# vim /etc/cassandra/conf/jvm.options

#-Xms2g
#-Xmx2g
-Xms512m
-Xmx512m
```