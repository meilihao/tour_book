# list

## tools
- [dropwatch - 监听系统丢包信息工具](https://cloud.tencent.com/developer/article/1638140)


## FAQ
### 丢包排查
参考:
- [Linux 系统 UDP 丢包问题分析思路](https://cloud.tencent.com/developer/article/1638140)

    方法1: dropwatch
    方法2: perf

        ```bash
        # perf record -g -a -e skb:kfree_skb
        # perf script
        ```