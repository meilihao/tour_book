# glide

[文档](https://deepzz.com/post/glide-package-management-command.html)

## FAQ
- glide设置mirror
```sh
$ glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto --vcs git
$ glide mirror set https://golang.org/x/net https://github.com/golang/net --vcs git
$ glide mirror set https://golang.org/x/tools https://github.com/golang/tools --vcs git
$ glide mirror set https://golang.org/x/text https://github.com/golang/text --vcs git
$ glide mirror set https://golang.org/x/image https://github.com/golang/image --vcs git
$ glide mirror set https://golang.org/x/sys https://github.com/golang/sys --vcs git
```

## FAQ
- Error looking for xxx: Cannot detect VCS
    - 当前路径使用了**链接**,glide不支持
    - 需要自己手动添加部分参数项,比如
        ```
        - package: git.xxx.com/base/lib
        repo:    git@git.xxx.com:base/lib.git
        vcs:     git
        subpackages:
        - log
        - rater
        - pager
        ```
> 当前golang官方的dep也不支持链接.