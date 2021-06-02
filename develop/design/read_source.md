# read_source spec
## 规划(以zstack举例)
1. 初始化保存`read_source`信息的repo

    1. 创建git repo: `git@gitee.com:chenhao/hello_zstack.git`
    1. 在源项目创建`__read_source/version.json`

        ```json
        $ cat <<'EOF' > version.json
        {
            "vcs":"git",
            "repo":"git@github.com:zstackio/zstack.git",
            "hash":"2f7d7c0634da00efcccd265ab99429a1349ba3fb",
            "hash_datetime":"2021-04-12T13:19:42Z",
            "tag":"",
            "remark":"阅读zstack v4.1.0"
        }
        EOF
        ```
    1. push到repo

        ```bash
        git remote add origin git@gitee.com:chenhao/hello_zstack.git
        git push -u origin main
        ```



1. 在源repo创建`__read_source`
```bash
$ git clone --depth 1 git@github.com:zstackio/zstack.git
$ cd zstack
$ git submodule add git@gitee.com:chenhao/hello_zstack.git __read_source # 不允许添加empty repo, 因此read_source repo必须有内容, 比如`version.json`.
$ cd __read_source
$ git submodule update --init --recursive # 看情况是否需要更新, 如果是按照上面的步骤直接在源项目创建read_source repo的话, 可忽略
```

## FAQ
### `__read_source`双下划线原因
最大可能避免与其他用户目录重名