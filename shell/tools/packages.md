# packages.md

## tools
- [command-not-found](https://command-not-found.com/)
- [cmd package search by man](http://manpages.ubuntu.com/manpages/focal/man8/), 比如[cgdisk](http://manpages.ubuntu.com/manpages/focal/man8/cgdisk.8.html)

## 打包(source -> deb/rpm)
- [BuildingTutorial](https://wiki.debian.org/BuildingTutorial#Building_the_modified_package)

## FAQ
### Ubuntu 如何解压 rpm 文件
```
# rpm2cpio xxx.rpm | cpio -div # `apt install rpm2cpio`
```