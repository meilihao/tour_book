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