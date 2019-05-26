# css 特效
## 吸顶
1. [position: sticky](https://developer.mozilla.org/zh-CN/docs/Web/CSS/position)
```css
position: sticky;
top: 0px;
```
1. 添加scroll事件，检测页面滚动距离大于一定值（一般为右栏与页面顶部的距离）时切换position为fixed布局