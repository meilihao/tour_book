### angularjs中$http模块POST请求request payload(json)转form data

解决方案：

1. 配置$httpProvider
```js
var myApp = angular.module('app', []);
myApp.config(function($httpProvider) {
    $httpProvider.defaults.transformRequest = function(obj) {
        var str = [];
        for (var p in obj) {
            str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
        }
        return str.join("&");
    }
    $httpProvider.defaults.headers.post = {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
});
//再如下面,在$http中使用data属性聚集参数即可.
```
1. 在post中配置
```js
$http({
    method: 'post',
    url: 'post.php',
    data: {
        name: "aaa",
        id: 1,
        age: 20
    },
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
    transformRequest: function(obj) {
        var str = [];
        for (var p in obj) {
            str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
        }
        return str.join("&");
    }
}).success(function(req) {
    console.log(req);
})
//或者http://www.bennadel.com/blog/2615-posting-form-data-with-http-in-angularjs.htm?utm_source=tuicool
```

### ng-options中as的作用
```json
{
    "myOptions":[
        {
            "id":106,
            "group":"Group 1",
            "label":"Item 1"
        }
    ]
}
```
```html
<select ng-model="myOption" ng-options="value.id as value.label group by value.group for value in myOptions">
    <option>--</option>
</select>
```
value.id将会和模型进行绑定,它的值会被存进ng-model属性里去(也就是option的value值).如果没有写value.id as,而是只写了value.label,那么,整个对象会被作为ng-model的值

## 基本概念

模板(Template) : 带有Angular扩展标记的**HTML** 
指令(Directive) : 用于通过**自定义属性和元素**扩展HTML的行为 
模型(Model) : 用于显示给用户并且与用户互动的**数据** 
作用域(Scope) : 用来存储模型(Model)的上下文(context),模型放在这个context中才能被控制器、指令和表达式等访问到 
表达式(Expression) : **模板中**可以通过它来访问作用域（Scope）中的变量和函数 
视图(View) : 用户看到的内容（即**DOM**） 
数据绑定(Data Binding) : **自动同步**模型(Model)中的数据和视图(View)表现 
控制器(Controller) : 视图(View)背后的**业务逻辑**

> 参考: [angularjs 概念概述](http://docs.ngnice.com/guide/concepts)或[angularjs Conceptual Overview](https://docs.angularjs.org/guide/concepts)

