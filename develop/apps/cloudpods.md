# cloudpods
## 概念
- 内置私有云: 管理本地IDC的未云化资源(on-premise), 主要是裸机, KVM虚拟机(Libvirt), VMware ESXi虚拟机(vSphere)
- 纳管私有云(priate cloud): zstack, openstack, 阿里飞天等
- 公有云(public cloud): aliyun, 腾讯云, 华为云, google云等

- 域(domain): 即多租户中的租户. 一个域下有一个完整的用户认证体系和资源和权限体系, 从而允许一个域管理员能够完全自治地管理本域的用户、组、项目、角色和权限策略.

	系统初始化后, 预置一个default域, 作为初始sysadmin账户和system项目所在的域.

## 部署
[部署 Cloudpods 平台的配置要求](https://www.cloudpods.org/zh/docs/setup/config/)

### 私有云
`./run.py -m https://mirrors.aliyun.com/pypi/simple/ virt 192.168.0.43`, 
默认使用了配置config-allinone-current.yml

注意:
1. 需要至少6G内存, 测试过4G但不够

条件调整:
- `ansible_processor_vcpus >= minimal_cpu`
	onecloud/roles/utils/misc-check/tasks/main.yml: minimal_cpu -> 2

## 权限
ref:
- [Openstack租户(项目)、用户、角色的概念与管理](https://juejin.cn/post/7000514148589109285)

## climc
ref:
- [Debug 模式](https://www.cloudpods.org/docs/guides/climc/usage/#debug-%E6%A8%A1%E5%BC%8F)

原理: 构建req再向后端各个服务发送API请求, 实现对资源的操控

子命令: 通过类似`R(&CredentialListOptions{}, "credential-list", "List all credentials",...)`将climc子命令注册进全局的CommandTable, 通过getSubcommandsParser()匹配CommandTable拿到子命令

### 调用链
ref:
- [通过资源名称定位代码](https://www.cloudpods.org/zh/docs/development/api-model/#%E9%80%9A%E8%BF%87%E8%B5%84%E6%BA%90%E5%90%8D%E7%A7%B0%E5%AE%9A%E4%BD%8D%E4%BB%A3%E7%A0%81)

```bash
# grep -rn 'zones' pkg/mcclient/modules # 发现zones属于compute服务,现在到compute服务的模型里面搜索servers关键字
pkg/mcclient/modules/compute/mod_zones.go:27:   Zones = modules.NewComputeManager("zone", "zones",
...
# grep -r 'zones' pkg/compute/models # 发现pkg/compute/models/zones.go比较匹配, 定位此文件
...
pkg/compute/models/zones.go:                    "zones",
...
```

## 源码
- appsrv.SWorkerManager

	- appsrv.SWorkerManager.Run()执行task(IWorkerTask)
		- task会进一步封装成sWorkerTask, 由appsrv.SWorkerManager.schedule()调度执行
			- 当`wm.activeWorker.size() < wm.workerCount && queueSize > 0`时, 新建一个worker, 并追加到wm.activeWorker, 再执行worker.run()
				- Pop一个sWorkerTask, 并execCallback(task)->IWorkerTask.Run()


> 部分routes入口`initHandlers`


组件:
- [ocadm](https://github.com/yunionio/ocadm)

### 定位
cloud create vm:
- [创建vm的接口定义，每个云都是通过这个入口创建vm，可以全局搜索CreateVM关键字查](https://github.com/yunionio/cloudmux/blob/release/3.11/pkg/cloudprovider/resources.go#L315)

kvm create vm:
- 创建虚拟机实例: pkg/compute/models/guests.go#`(manager *SGuestManager) ValidateCreateData`

### tags
- [添加计算节点](https://www.cloudpods.org/zh/docs/setup/host/)

## FAQ
### `Unable to locate package python3-pyyaml`
ubuntu22 python3-pyyaml实际就是python3-yaml, 或`python3 -m pip install PyYAML`