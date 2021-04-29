# vue3
## 变化
- template标签下面可以有多个节点
- setup函数可以代替之前的data，methods，computed，watch，Mounted等对象，但是props声明还是在外面
- ref与reactive方法的区别是什么？一个是把值类型添加一层包装，使其变成响应式的引用类型的值. 另一个则是引用类型的值变成响应式的值. 所以两者的区别只是在于是否需要添加一层引用包装？其目的都是对数据添加响应式效果.