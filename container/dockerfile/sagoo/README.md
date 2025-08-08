# README
ref:
- [duruo850/sagooiot-xxx : 参考镜像](https://hub.docker.com/u/duruo850?page=1&search=sagoo)
- [展联科技ZL450边缘网关采集上传及开源物联网平台SagooIoT搭建](https://blog.csdn.net/sym_robot/article/details/139334582)

## 部署
```bash
# --- build sagooiot-ui
WD=/data/tmpfs/sagoo

git clone --depth 1 git@github.com:sagoo-cloud/sagooiot-ui.git
cd sagooiot-ui
pnpm install -g yarn # 使用yarn, npm build会报错
yarn install
cp -rf dist/* $WD/nginx/html

cd $WD
sudo chown -R 999:999 redis/data # 避免redis没有权限操作数据目录, 可`docker compose up -d redis`后进入容器执行`id redis`获取

docker compose up -d emqx 
docker compose up -d redis
docker compose up -d oceanbase
docker compose up -d tdengine
docker compose up -d nginx

# --- init db
mariadb -h172.19.0.2 -P2881 -uroot@test -p --skip_ssl -e "CREATE DATABASE sagoo_iot_open;"
mariadb -h172.19.0.2 -P2881 -uroot@test -Dsagoo_iot_open -p --skip_ssl < init.sql # init.sql from [sagooiot/manifest/docker-compose/mysql/init](https://github.com/sagoo-cloud/sagooiot/tree/main/manifest/docker-compose/mysql/init), 它和[sagooiot/manifest/sql/sql.zip](https://github.com/sagoo-cloud/sagooiot/tree/main/manifest/sql)相同

docker compose up -d sagoo-iot-open
```

[系统访问](https://iotdoc.sagoo.cn/docs/install/sagooiot-install): http://localhost admin/admin123456
