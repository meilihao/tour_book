# myems
ref:
- [MyEMS视频](https://space.bilibili.com/539108162/upload/video)
- [Docker Compose](https://myems.io/docs/installation/docker-compose)

```bash
# git clone --depth 1 -b v5.7.1 git@gitee.com:myems/myems.git
# cd myems
# pushd .
# --- init db
# cd database/install
# mariadb -h 127.0.0.1 -u root -p < myems_billing_db.sql # 借用了同级sago项目的mariadb
# mariadb -h 127.0.0.1 -u root -p < myems_carbon_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_energy_baseline_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_energy_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_energy_model_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_energy_plan_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_energy_prediction_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_fdd_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_historical_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_production_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_reporting_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_system_db.sql
# mariadb -h 127.0.0.1 -u root -p < myems_user_db.sql
# mariadb -h 127.0.0.1 -u root -p -e "show databases;" # 核对dbs
# mariadb -h 127.0.0.1 -u root -p < ../demo-cn/myems_system_db.sql # 插入演示数据, 可选, 见[数据库](https://myems.io/docs/installation/database)
# popd
# pushd .
# --- 修改nginx.conf里的API配置: 都是将`proxy_pass http://127.0.0.1:8000/`改为`proxy_pass http://api:8000/`
# vim myems-admin/nginx.conf
# vim myems-web/nginx.conf
# popd
# --- 配置env: 修改db相关ip, password
# cp myems-aggregation/example.env myems-aggregation/.env
# vim myems-aggregation/.env
# cp myems-api/example.env myems-api/.env
# vim myems-api/.env
# cp myems-cleaning/example.env myems-cleaning/.env
# vim myems-cleaning/.env
# cp myems-modbus-tcp/example.env myems-modbus-tcp/.env
# vim myems-modbus-tcp/.env
# cp myems-normalization/example.env myems-normalization/.env
# vim myems-normalization/.env
# popd
# pushd .
# --- 修改web配置: 未修改
# cd myems-web
# vim src/config.js
# npm i --unsafe-perm=true --allow-root --legacy-peer-deps
# npm run build
# popd
# --- 修改docker-compose-on-linux.yml: 将`/myems-upload`改为`./myems-upload`
# cd others
# mkdir myems-upload
# vim docker-compose-on-linux.yml
# docker compose -f docker-compose-on-linux.yml up -d # 建议先手动pull相关基础镜像再执行该命令, 其次注意与其他docker容器的端口冲突, 比如80,443等
```

[验证](https://myems.io/docs/installation/docker-compose):
- myems-web	127.0.0.1:80	[输入账号密码(`administrator@myems.io/!MyEMS1`)登录成功](https://myems.io/docs/installation/docker-linux#%E9%BB%98%E8%AE%A4%E5%AF%86%E7%A0%81)
- myems-admin	127.0.0.1:8001	输入账号密码(`administrator/!MyEMS1`)登录成功
- myems-api	127.0.0.1:8000/version	返回版本信息

## 文件
- docker-compose-on-linux.yml: from 官方[docker-compose-on-linux.yml](https://gitee.com/myems/myems/blob/master/others/docker-compose-on-linux.yml)