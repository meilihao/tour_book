# 打包

## 命名
建议: `{project}_{compent}-{version}`, 必须使用compent的原因: 避免起先只有`{project}-{version}`后新增组件, 此时按`rpm -qa|grep <project>`出现多条记录.

## 升级
1. 将升级程序upgrader单独拆成一个rpm并优先升级

	1. 其他组件的升级逻辑可能在新版upgrader中
	2. 假设upgrader在某个组件中, 升级该组件时可能遇到upgrader `is busy`(即在使用中, 因为管理节点可能正在使用该upgrader执行升级)的错误