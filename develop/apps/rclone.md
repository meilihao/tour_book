# rclone
## 功能
- rclone copy  : 复制
- rclone move  : 移动，如果要在移动后删除空源目录，加上 --delete-empty-src-dirs 参数
- rclone sync  : 同步：将源目录同步到目标目录，只更改目标目录

    **非实时同步**
- rclone size  : 查看网盘文件占用大小
- rclone delete: 删除路径下的文件内容
- rclone purge : 删除路径及其所有文件内容
- rclone mkdir : 创建目录
- rclone rmdir : 删除目录
- rclone rmdirs: 删除指定环境下的空目录。如果加上 --leave-root 参数，则不会删除根目录
- rclone check : 检查源和目的地址数据是否匹配
- rclone ls    : 列出指定路径下的所有的文件以及文件大小和路径
- rclone lsl   : 比上面多一个显示上传时间
- rclone lsd   : 列出指定路径下的目录
- rclone lsf   : 列出指定路径下的目录和文件

### 常用参数
-n = --dry-run  测试运行，查看Rclon在实际运行中会进行哪些操作
-P = --progress 显示实时传输进度，500mS刷新一次，否则默认1分钟刷新一次
--cache-chunk-size SizeSuffi    块的大小，默认5M越大上传越快，占用内存越多，太大可能会导致进程中断
--cache-chunk-total-size SizeSuffix 块可以在本地磁盘上占用的总大小，默认10G
--transfers=N   并行文件数，默认为4
--config string 指定配置文件路径，string为配置文件路径
--ignore-errors 跳过错误
--size-only 根据文件大小校验，不校验hash
--drive-server-side-across-configs  服务端对服务端传输

### filter
--exclude   排除文件或目录
--include   包含文件或目录
--filter    文件过滤规则，相当于上面两个选项的其它使用方式。包含规则以+开头，排除规则以-开头

#### 过滤规则文件
--filter-from <规则文件> : 从文件添加包含/排除规则

过滤规则文件filter-file.txt示例：
```txt
- secret*.jpg
+ *.jpg
+ *.png
+ file2.avi
- /dir/Trash/**
+ /dir/**
- *
```

#### 文件类型过滤
比如--exclude "*.bak" --filter "- *.bak"排除所有bak文件
比如--include "*.{png,jpg}" --filter "+ *.{png,jpg}"包含所有png和jpg文件，排除其他文件
--delete-excluded删除排除的文件。需配合过滤参数使用，否则无效

#### 目录过滤
目录过滤需要在目录名称后面加上/否则会被当做文件进行匹配
以/开头只会匹配根目录（指定目录下），否则匹配所目录，这同样适用于文件
--exclude ".git/"排除所有目录下的.git目录
--exclude "/.git/"只排除根目录下的.git目录
--exclude "{Video,Software}/"排除所有目录下的Video和Software目录
--exclude "/{Video,Software}/"只排除根目录下的Video和Software目录
--include "/{Video,Software}/**"仅包含根目录下的Video和Software目录的所有内容

#### 大小过滤
默认大小单位为kBytes但可以使用k M或G后缀
--min-size过滤小于指定大小的文件. 比如--min-size 50表示不会传输小于50k的文件
--max-size过滤大于指定大小的文件. 比如--max-size 1G表示不会传输大于1G的文件

### log
有4个级别的日志记录：ERROR NOTICE INFO DEBUG
默认情况下Rclon将生成ERROR NOTICE日志

命令  说明
-q  rclone将仅生成ERROR消息
-v  rclone将生成ERROR NOTICE INFO 消息
-vv rclone 将生成ERROR NOTICE INFO DEBUG 消息
--log-level LEVEL   标志控制日志级别

输出日志到文件:
使用--log-file=FILE选项rclone会将Error Info Debug消息以及标准错误重定向到FILE

## 使用
参考:
- [如何使用rclone同步远程云盘](http://zhangzr.com/2018/11/08/rclone/)
- [Rclone结合MinIO Server](http://docs.minio.org.cn/docs/master/rclone-with-minio-server)

```bash
# rclone config  # 添加、删除、管理网盘等操作
# rclone config file # 显示配置文件的路径
# rclone config show # 显示配置文件信息

# rclone config show
[minio-r]
type = s3
provider = Minio
env_auth = false
access_key_id = root
secret_access_key = password
endpoint = http://127.0.0.1:9000

[oss-hwpf]
type = s3
provider = Alibaba
access_key_id = <your-ali-access-key-id>
secret_access_key = <your-ali-secret-access-key>
endpoint = oss-cn-hongkong.aliyuncs.com
acl = public-read

[s3-overseas]
type = s3
provider = AWS
env_auth = false
access_key_id = <your-aws-access-key-id>
secret_access_key = <your-aws-secret-access-key>
region = <your-region-id>
acl = public-read
endpoint = https://s3.<your-region-id>.amazonaws.com

# rclone copy -P /home/SunPma minio-r:/home/SunPma --transfers=8 # 复制本地到网盘，并显示实时传输进度，设置并行上传数为8. bucket path=home/SunPma
# rclone copy 配置名称:网盘路径 配置名称:网盘路径 --drive-server-side-across-configs # 如果需要服务端对服务端的传输可加以下参数（不消耗本地流量）
# rclone lsjson [-R] minio-l: # 以json格式[递归]输出
```

## FAQ
### rclone config位置
`~/.config/rclone/rclone.conf`