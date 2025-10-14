# ossutil
ref:
- [命令行工具ossutil 2.0](https://help.aliyun.com/zh/oss/developer-reference/install-ossutil2)

ossutil 2.0 在多种操作系统中高效管理阿里云对象存储 OSS 资源，实现文件的快速上传、下载、同步和管理，适合开发者、运维人员和企业进行大规模数据迁移和日常运维操作

## example
```
# cat ~/.ossutilconfig 
[default]
accessKeyId=xxx
accessKeySecret=yyy
region=cn-hangzhou
endpoint=oss-cn-hangzhou.aliyuncs.com
# ossutil ls oss://my-input-output/uploadfiles
# ossutil cp xxx.xlsx oss://my-input-output/uploadfiles/xxx.xlsx
# ossutil cp oss://my-input-output/uploadfiles/xxx.xlsx .
```