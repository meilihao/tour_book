### 部署go tour

```shell
$ go get github.com/golang/tour
$ go get github.com/golang/net
$ go get github.com/golang/tools
# 因为github上的repo只是golang官方的mirror,因此需要调整相应目录路径
# 将$GOPATH/github.com/golang下的tour,net,tools三个文件夹移到$GOPATH/golang.org/x下面
$ cd $GOPATH/golang.org/x/tour/gotour
$ go build
# 运行go
$ ./gotour

# 当然,也可直接go get golang.org下相应项目的源码,不过下载时推荐使用gopm.io
```

参考:

- [中文版go-tour](https://bitbucket.org/mikespook/go-tour-zh)
- [go-tour源码阅读](http://www.cnblogs.com/yjf512/archive/2012/12/13/2816480.html)