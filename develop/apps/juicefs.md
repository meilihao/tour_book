# juicefs
JuiceFS的定位是一个建立在OSS等存储系统之上的一个虚拟文件系统. 当前它使用Redis来做文件的metadata管理，内部有对数据的Cache功能，所以当重复访问数据时，速度也会得到明显提升(但有丢数据的风险).

## [部署](https://github.com/juicedata/juicefs/blob/main/docs/zh_cn/quick_start_guide.md)
1. 准备redis

```bash
mkdir redis-data
sudo docker run -d --name redis \
    -v $PWD/redis-data:/data \
    -p 6379:6379 \
    --restart unless-stopped \
    redis:6.2-alpine redis-server --appendonly yes --requirepass password
```

> 这里redis使用`db 1`.

1. 准备minio
```bash
mkdir minio-data
sudo docker run -d --name minio \
    -e "MINIO_ROOT_USER=root" \
    -e "MINIO_ROOT_PASSWORD=password" \
    -v $PWD/minio-data:/data \
    -p 9000:9000 \
    --restart unless-stopped \
    minio/minio server /data
```

> 使用`http://127.0.0.1:9000`访问 MinIO 管理界面，root 用户初始的 Access Key 和 Secret Key 均为`minioadmin`, 可通过设置MINIO_ROOT_USER和MINIO_ROOT_PASSWORD来修改.

1. 创建并挂载 JuiceFS 文件系统
```
juicefs format \
    --storage minio \
    --bucket http://127.0.0.1:9000/pics \
    --access-key root \
    --secret-key password \
    redis://:password@127.0.0.1:6379/1 \
    pics
sudo juicefs mount -d redis://:password@127.0.0.1:6379/1 /mnt/jfs
df -Th
```

## FAQ
### [JuiceFS 客户端升级](https://github.com/juicedata/juicefs/blob/main/docs/en/client_compile_and_upgrade.md#juicefs-client-upgrade)