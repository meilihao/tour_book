# hydra
ORY Hydra是经过强化，经过OpenID认证的OAuth 2.0服务器和OpenID Connect提供商，针对低延迟，高吞吐量和低资源消耗进行了优化. ORY Hydra 不是身份提供者（用户注册，用户登录，密码重置流程）, 而是通过登录和同意应用程序连接到您现有的身份提供者. 以不同的语言实现登录和同意应用程序很容易，并且提供了示例性的同意应用程序（Go，Node）和 SDK.

## 部署
参考:
- [Run your own OAuth2 Server](https://www.ory.sh/run-oauth2-server-open-source-api-security/)
- [**Hydra项目介绍**](https://www.lsdcloud.com/blog/Go/hydra.html)
- [ORY Hydra项目详解](https://blog.csdn.net/qq_37493556/article/details/106699444)

1. 构建hydra

    ```bash
    git clone --depth 1 git@github.com:ory/hydra.git
    cd hydra && go build
    ```
1. 准备database

    1. 在mysql上创建hydra数据库
    1. 初始化db: `./hydra migrate sql mysql://root:admin@tcp\(127.0.0.1:3306\)/hydra`
1. 启动hydra

    [hydra config example](https://github.com/ory/hydra/blob/master/docs/versioned_docs/version-v1.10/reference/configuration.md), 中文注解见[这里](https://blog.csdn.net/qq_37493556/article/details/106699444), 使用配置时需要hydra的`hydra serve -c`参数 

    ```bash
    export SECRETS_SYSTEM=$(export LC_CTYPE=C; cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
    export DSN=mysql://hydra:hydra@tcp\(openhello.net:3306\)/hydra
    env URLS_SELF_ISSUER=http://127.0.0.1:4444 URLS_CONSENT=http://127.0.0.1:3000/consent URLS_LOGIN=http://127.0.0.1:3000/login ./hydra serve all --dangerous-force-http # 3000是hydra-login-consent-node使用的端口
    ```

    检查hydra是否启动成功: `curl http://127.0.0.1:4445/health/ready`, 返回`{"status":"ok"}`即表示成功.
    
    > [查看项目接口文档](https://www.ory.sh/hydra/docs/reference/api/)
1. 部署hydra用户登录和同意流程的参考实现

    ```bash
    git clone --depth 1 git@github.com:ory/hydra-login-consent-node.git
    npm install 
    env HYDRA_ADMIN_URL=http://127.0.0.1:4445 NODE_TLS_REJECT_UNAUTHORIZED=0 npm run start
    ```

## 演示(Authorization Code)
1. Create OAuth2 Consumer App

    ```bash
    ./hydra clients create --endpoint http://127.0.0.1:4445 --id another-consumer --secret consumer-secret -g authorization_code,refresh_token -r token,code,id_token --scope openid,offline --callbacks http://127.0.0.1:9010/callback
    ```
1. 执行流程

    ```bash
    ./hydra token user \
    --port 9010 \
    --auth-url http://127.0.0.1:4444/oauth2/auth \
    --token-url http://127.0.0.1:4444/oauth2/token \
    --client-id another-consumer \
    --client-secret consumer-secret \
    --scope openid,offline \
    --redirect http://127.0.0.1:9010/callback
    ```

    该命令会启动浏览器并打开`http://127.0.0.1:9010/`, 点击页面上的`Authorize application`连接跳转到`http://127.0.0.1:3000`即可进入授权流程.