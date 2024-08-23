# 深信服HCI
env:
- [SANGFOR HCI 6.10.0](https://support.sangfor.com.cn/productDocument/read?product_id=33&version_id=993&category_id=283660)
- [深信服HCI超融合方案介绍](https://www.yunzhan365.com/basic/94915681.html)
- [深信服aSV服务器虚拟化功能及原理](https://www.sangfor.com.cn/news/202111202111261041)

    kvm+QCOW2
- [深信服超融合平台SDK](https://github.com/ynu/scp)

产品:
- 深信服云管平台SCP(Sangfor Cloud Platform): 支持管理HCI+vmware

## [无代理备份](https://support.sangfor.com.cn/productDocument/read?product_id=36&version_id=1022&category_id=285299)
ref:
- [SCP API接口](https://support.sangfor.com.cn/productDocument/read?product_id=36&version_id=1022&category_id=285689)

无代理备份基于云平台和超融合提供的OpenAPI接口和磁盘数据访问SDK（SFVDDK）, 因此需要SCP的OpenAPI+SFVDDK.

## FAQ
### [开启ssh](https://support.sangfor.com.cn/productDocument/read?product_id=33&version_id=993&category_id=283314&type=1)
没密码.

### 添加本地盘
存储->其他存储->添加存储->本地存储

### 上传iso
存储->其他存储->选中某个`<存储>`记录->管理存储空间