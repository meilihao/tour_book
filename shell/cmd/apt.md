# apt

## 描述

debian,ubuntu等发行版的包管理.

## 格式

## 例

    # apt-cache madison pouch # 列出软件包的所有版本
    # apt install pouch=1.0.0-0~ubuntu # 安装指定版本的软件包
    # apt-get install --only-upgrade samba # 仅更新单个package
    # apt list -a cifs-utils # package all version
    # apt-cache policy cifs-utils # package all version, 推荐
    # rmadison cifs-utils # package all version, 推荐

## FAQ
### apt install报`Size mismatch`
下载到的deb软件包信息和源信息列表Packages记录(Packages.gz)的数据不相符, 可用`dpkg -i`安装