# hummingbird
## 部署
```bash
# --- from https://github.com/winc-link/hummingbird
# git clone https://gitee.com/winc-link/hummingbird.git
# cd hummingbird/manifest/docker
# docker-compose up -d
# --- 查看db
# sudo docker exec -it eda220be0cdb sh # cwd:/var/bin
# cat /etc/hummingbird-core/configuration.toml
# apk add sqlite
# cd /
# cp /var/bin/hummingbird/db-data/core-data/core.db .
# sqlite3 core.db
```

访问http://localhost:3000, admin/admin, 首次输入是初始化admin密码, 再输入账号密码即可登入.