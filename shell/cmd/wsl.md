# wsl
```bash
# wsl --list --online #  查看可供安装的可用发行版列表
# wsl --list --verbose # 查看已安装的 WSL 发行版及其状态
# wsl --install ubuntu # 安装在线镜像
# wsl --update    # 更新 WSL 内核
# wsl --shutdown  # 立即终止所有正在运行的 WSL 发行版
# wsl --terminate <Distro>   # 终止指定的发行版
# wsl --unregister <Distro>   # 注销并删除发行版（数据会丢失！）
# wsl --install --from-file Ubuntu2404-250130_x64.wsl # 安装离线镜像
# wsl -d <Distro> # 启动指定发行版
```

## FAQ
### WSL启动ubuntu 24.04.1报`WslRegisterDistribution failed with error: 0x80370114`
控制面板-程序-启用或关闭Windows功能-打开"适用于Linux的Windows子系统", 再重启即可

### 远程到WSL
- 非转发, **推荐**

    "WSL Settings"-网络:
    1. 网络模式: Mirrored

        让wsl共享host ip

        > 默认是nat
    2. 主机地址环回: 启用

        for 容器
    3. 重启wsl

- 转发
    1. 安装openssh-server并配置

        ```conf
        Port 2222                  # 避免与 Windows 的默认 SSH 端口冲突
        ListenAddress 0.0.0.0      # 允许所有 IP 访问
        ```
    1. 将端口转发到WSL

        `netsh interface portproxy add v4tov4 listenport=2222 listenaddress=0.0.0.0 connectport=2222 connectaddress=localhost`, connectxxx是wsl的信息

        > 查看所有portproxy: `netsh interface portproxy show v4tov4`
    1. 设置Windows防火墙入站规则

        `netsh advfirewall firewall add rule name=WSL2 dir=in action=allow protocol=TCP localport=2222`

> `.wsconfig`(在windows当前User目录下)默认不存在, 配置"WSL Settings"后生成