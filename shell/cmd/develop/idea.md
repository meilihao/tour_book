# idea

## FAQ
### `Cannot resolve symbol 'java'`/修改Java SDK
File->Project Structure->Project Settings-> Project, 选择SDK版本.

### `cannot find declaration to go to`
File->Project Structure->Project Settings-> Modules, 删除错误的module配置. 在通过`+`->`Import Module`->选中项目根pom.xml所在的目录

其他方法:
1. 删除项目的`.idea`, 再重新打开项目. 如果重新打开后有部分跳转正常, 有部分没有, 再重新打开即可.

### 禁用插件
Flie -> Settings -> Plugins 找到插件, 通过右键菜单禁用

### 源码和sample在一个project, 如何高亮
比如`ocksdb-6.29.5/java`

解决方法: 将其samples和src目录同时`Mark Directory as`为`Sources Root`

### 如何同步maven
窗口右侧选中`m`标签页, 选择第一个`转圈`图标中的`Reload All Maven Projects`, 如果其运行还是报错, 就点击Sync状态窗口右侧的`Try ...`(即强制同步)

### Cannot resolve symbol 'xxx'
选择错误, 右键选择`quick fix`, 尝试让idea自行修复