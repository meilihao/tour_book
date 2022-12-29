# mssql
## 连接字符串
支持:
- Windows 身份验证
- SQL Server 身份验证: `qlserver://username:password@host/instance?param1=value&param2=value`

## FAQ
### windows server 2012安装sql server 2012报错
报错`安装Microsoft .NET Framework 3.5 时出错\n启用 Windows 功能 NetFx3 时出错, 错误代码:-2146498298...`

这是由于SQL Server 2012数据库系统的运行需要依靠`.NET Framework 3.5`, 但是windows server 2012默认是不安装`.netframework3.5`的, 所以必须先在操作系统上安装`.NET Framework 3.5`

解决方法:
1. 有系统盘, **推荐**
	1. 打开`服务器管理器`, 选择`添加角色和功能`
	1. 一直下一步到功能选项, `选择.NET Framework 3.5`, 在
	1. 在`确认`选择左下角的`指定备用源路径`, 选择iso的`F:\sources\sxs`
	1. 安装

	实际执行的是`dism.exe /online /enable-feature /featurename:netfx3 /Source:D:\WindowsOS\sources\sxs`, 直接执行可能会报错, 因为还需要些依赖, 用`服务器管理器`安装可自动处理依赖.
2. 下载NetFx3.cab, 未找到离线文件, 不推荐

	以管理员cmd执行: `dism.exe /online /add-package /packagepath:C:\WINDOWS\netfx3.cab`

### 获取SQL Server实例名
1. 服务(services.msc) -> `SQL Server(<实例名>)`, 默认实例为(MSSQLSERVER)
2. SQL Server 配置管理器 -> SQL Server 服务 -> `SQL Server(<实例名>)`, 默认实例为(MSSQLSERVER)