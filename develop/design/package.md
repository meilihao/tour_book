# 打包

## 命名
建议: `{project}_{compent}-{version}`, 必须使用compent的原因: 避免起先只有`{project}-{version}`后新增组件, 此时按`rpm -qa|grep <project>`出现多条记录.