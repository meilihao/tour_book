# ingress-nginx
如何部署ingress.

参考:
- [ingress-nginx/examples/deployment/nginx/](https://github.com/kubernetes/ingress-nginx/tree/a3131c569f6b0d6502dda1548f63be4b56e36e42/examples/deployment/nginx)

1. 创建证书的secret

    - rancher: 资源->密文-> 证书列表, 创建名为`openhello-com`的secret.
1. 部署容器

    ```sh
    $ kubectl apply -f service.yaml
    $ kubectl apply -f ingress.yaml
    ```
1. 因为用了标准nginx镜像, 因此部署好后需要修改static-server-nginx容器的nginx配置.

    ```sh
    # apt update
    # apt install nano
    # vim /etc/nginx/conf.d/defaut.conf
    # nginx -s reload
    # rm -rf /var/lib/apt/lists/*
    ```

    > 推荐使用`go+alpine`