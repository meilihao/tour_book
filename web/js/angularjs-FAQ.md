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
