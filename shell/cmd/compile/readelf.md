# readelf
一个非常有用的解析 ELF 二进制文件的工具. 我们可使用 readelf 命令收集
符号、段、节、重定向入口、数据动态链接等相关信息.

## 示例
```
readelf -S <object> # 查询节头表
readelf -l <object> # 查询程序头表
readelf -s <object> # 查询符号表
readelf -e <object> # 查询 ELF 文件头数据
readelf -r <object> # 查询重定位入口
readelf -d <object> # 查询动态段
readelf a.out -x .rodata # 以16进制输出dump的内容
```