# vue-router

## FAQ

### vue-router路由出口(路由匹配到的组件将渲染在这里)
`children[n].component`的占位符是`parent component`里的`<router-view />`

### `path: ''`
```vue
const router = new VueRouter({
  routes: [
    {
      path: '/user/:id', component: User,
      children: [
        // 当 /user/:id 匹配成功，
        // UserHome 会被渲染在 User 的 <router-view> 中
        { path: '', component: UserHome }, // 如果没有`path: ''`, 你访问 /user/foo 时，User 的出口(`<router-view></router-view>`)是不会渲染任何东西，这是因为没有匹配到合适的子路由

        // ...其他子路由
      ]
    }
  ]
})
```

### 前后路由类似, 页面不变(请求路径被缓存)
如果目的地和当前路由相同，只有参数发生了改变 (比如`path: '/users/:id'`, 从一个用户资料到另一个 /users/1 -> /users/2)，需要使用 beforeRouteUpdate 来响应这个变化 (比如抓取用户信息)

> 从 /users/1 -> /users/2，**原来的组件实例会被复用**. 因为两个路由都渲染同个组件，比起销毁再创建，复用则显得更加高效. **不过，这也意味着组件的生命周期钩子不会再被调用**

### 导航守卫和redirect
导航守卫并没有应用在跳转路由上，而仅仅应用在其目标上. 在`{ path: '/a', redirect: ...`中，为 /a 路由添加一个 beforeEach 或 beforeLeave 守卫并不会有任何效果.