# 深信服HCI
env:
- [SANGFOR HCI 6.10.0](https://support.sangfor.com.cn/productDocument/read?product_id=33&version_id=993&category_id=283660)
- [深信服HCI超融合方案介绍](https://www.yunzhan365.com/basic/94915681.html)
- [深信服aSV服务器虚拟化功能及原理](https://www.sangfor.com.cn/news/202111202111261041)

    kvm+QCOW2
- [深信服超融合平台SDK](https://github.com/ynu/scp)

产品:
- 深信服云管平台SCP(Sangfor Cloud Platform): 支持管理HCI+vmware

    SCP是运行在HCI上的vm

## [无代理备份](https://support.sangfor.com.cn/productDocument/read?product_id=36&version_id=1022&category_id=285299)
ref:
- [SCP API接口](https://support.sangfor.com.cn/productDocument/read?product_id=36&version_id=1022&category_id=285689)

    api注意点:
    1. region/serviceg分别固定为`cn-south-1/sdk-api`
    1. 与标准AWS4-HMAC-SHA256的差异:
        1. CanonicalHeaders只包含了`x-amz-date`, 且SignedHeaders只包含了`x-amz-date`
        1. CanonicalQueryString传空, 即不参与签名

> Open-API-zh_CN-2024-05-16.pdf : **文档质量很糟糕**

无代理备份基于云平台和超融合提供的OpenAPI接口和磁盘数据访问SDK（SFVDDK）, 因此需要SCP的OpenAPI+SFVDDK.

scp 6.10.0和hci 6.10.0及以后版本, 且scp需要纳管hci.

配置:
1. 在hci开启`无代理备份数据传输服务`(系统管理->端口管理)
1. 在scp配置`LAN模式备份数据传输网口`(产品与服务->备份与CDP->设置, 需要选择资源池)
1. 通过scp+admin创建权限策略(包含所需备份权限)+关联账户xxx(比如afb, agent free backup) 

## FAQ
### [开启ssh](https://support.sangfor.com.cn/productDocument/read?product_id=33&version_id=993&category_id=283314&type=1)
没密码.

### 添加本地盘
存储->其他存储->添加存储->本地存储

### 上传iso
存储->其他存储->选中某个`<存储>`记录->管理存储空间

### 部署SCP
1. 导入SCP 虚拟机镜像
1. [配置SCP vm所在节点的`物理出口`, 并将scp vm的网络连接到该`物理出口`](https://www.bilibili.com/video/BV1WM4m1C7NT)
1. 开启scp vm, 进入维护模式修改ip, 并测试连通性, 再重启
1. 通过`https://设置的IP:4430`访问scp即可

### scp添加集群后没有原先创建的HCI vm
云资源->资源池->HCI->创建, 创建成功后能看到HCI vms(scp vm除外)