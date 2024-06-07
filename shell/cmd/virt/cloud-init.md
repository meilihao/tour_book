# cloud-init
ref:
- [云实例初始化工具cloud-init简介](https://www.cnblogs.com/frankming/p/16281447.html)
- [云实例初始化工具cloud-init源码分析](https://www.cnblogs.com/frankming/p/16281612.html)

cloud-init是一款用于初始化云服务器的工具，它拥有丰富的模块，能够为云服务器提供的能力有：初始化密码、扩容根分区、设置主机名、注入公钥、执行自定义脚本等等，功能十分强大。

目前为止cloud-init是云服务器初始化工具中的事实标准，它几乎适用于所有主流的Linux发行版，也是各大云厂商正在使用的默认工具，社区活跃。基于Python语言使得它能够轻易跨平台、跨架构运行，良好的语法抽象使得它适配新模块、新发行版十分容易。
