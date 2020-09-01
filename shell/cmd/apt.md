# apt

## 描述

debian,ubuntu等发行版的包管理.

## 格式

## 例
```bash
# apt-cache madison pouch # 列出软件包的所有版本
# apt install pouch=1.0.0-0~ubuntu # 安装指定版本的软件包
# apt-get install --reinstall aptitude # 重新安装package
# apt-get install --only-upgrade samba # 仅更新单个package
# apt list -a cifs-utils # package all version
# apt-cache policy cifs-utils # package all version, 推荐
# rmadison cifs-utils # package all version, 推荐
# apt-cache depends -i samba # 查看依赖
# apt install --simulate samba # 仅模拟不安装
# apt install --download-only samba # 仅下载, 不安装
# apt list --installed # 查看已安装的package
# apt-cache show cpio # 查看软件依赖
# apt purge package_name # remove命令卸载指定软件包，但是留下一些包文件. 如果想彻底卸载软件包，包括它的文件，使用purge替换remove
```

## FAQ
### apt install报`Size mismatch`
下载到的deb软件包信息和源信息列表Packages记录(Packages.gz)的数据不相符, 可用`dpkg -i`安装

### apt install 安装的deb的缓存位置
ubuntu中由apt-get获得的文件包保存在/var/cache/apt/archives