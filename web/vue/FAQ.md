# FAQ

## 组件的data是函数
如果传给组件的 data 是一个原始对象的话，则在建立多个组件实例时它们就会共用这个 data 对象，修改其中一个组件实例的数据就会影响到
其他组件实例的数据，这显然不是我们所期望的.

## 访问404
vue单页应用使用vue-router会有两种配置，即history模式和hash模式，但是hash模式其实会有很多限制，最主要的一点，url地址太丑了，所以我们在生产环境中也希望用到history模式.

但在HTML5 History 模式下, 使用nginx等代理服务器时，会遇到404的问题, 这是因为vue项目编译出来的dist中，并没有真正的本地资源提供给nginx，正确的做法是需要转交给vue-router来做前端路由.

```nginx
...
location / {
    root /root/dist;
    try_files $uri $uri/ /index.html =404; # 原理 404时直接返回index.html,再由ue-router来做前端路由
}
...
```