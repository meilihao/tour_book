# list

## 高可用
- [Linux网卡bond的七种模式详解](https://blog.51cto.com/linuxnote/1680315)

## tools
- [dropwatch - 监听系统丢包信息工具](https://cloud.tencent.com/developer/article/1638140)
- [命令行 DNS 查询工具，支持 DNS-over-TLS 和 DNS-over-HTTPS](https://github.com/mr-karan/doggo)
- [BGP 工具探索](https://linux.cn/article-13857-1.html)

## bgp
- [Syntropy的初创公司提出了一种去中心化自主路由协议（DARP），旨在取代BGP](https://linux.cn/article-13204-1.html)

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