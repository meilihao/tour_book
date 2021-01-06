## nvm

建议使用版本管理工具[nvm](https://github.com/creationix/nvm)，来安装Node.js和io.js。

## 其他安装方式

- [通过包管理器安装](https://github.com/joyent/node/wiki/Installing-Node.js-via-package-manager)
- [Linux 二进制安装](https://github.com/nodejs/help/wiki/Installation)

 1. 在node官网下载[最新版Linux Binaries](http://nodejs.org/dist/v0.12.4/node-v0.12.4-linux-x64.tar.gz).

 2. 将`node-v0.12.4-linux-x64.tar.gz`解压重命名为`node`并移到`/opt`下.

 3. 设置环境变量后重启电脑即可.

         export NODE_HOME=/opt/node
         export PATH=$PATH:$NODE_HOME/bin
         export NODE_PATH=$PATH:$NODE_HOME/lib/node_modules