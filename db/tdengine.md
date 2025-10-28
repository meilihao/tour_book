# TDengine
ref:
- [TDengine 发布历史及下载链接](https://docs.taosdata.com/releases/tdengine/)

## 概念
### datatype(数据类型)
ref:
- [数据类型](https://docs.taosdata.com/reference/taos-sql/datatype/)

## 部署
ref
- [快速体验](https://docs.taosdata.com/get-started/)

docker:
```bash
docker run -d \
  -v $(pwd)/data/taos/dnode/data:/var/lib/taos \
  -v $(pwd)/data/taos/dnode/log:/var/log/taos \
  -p 6030:6030 -p 6041:6041 -p 6043:6043 -p 6060:6060 \
  -p 6044-6049:6044-6049 \
  -p 6044-6045:6044-6045/udp \
  --name tsdb \
  tdengine/tsdb-ee
```