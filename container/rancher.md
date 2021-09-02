# rancher
参考:
- [Rancher v2.x 使用手册](https://www.bookstack.cn/books/rancher-v2.x)

## 升级
- [单节点版rancher升级指南](https://blog.maoxianplay.com/posts/rancher-update-2.2.1/)
- [单节点升级(官方)](https://docs.rancher.cn/docs/rancher2/upgrades/upgrades/single-node/_index)

## 清理节点(**重装rancher时必须清理**)
- [清理节点](https://docs.rancher.cn/docs/rancher2/cluster-admin/cleaning-cluster-nodes/_index)
- [Removing Kubernetes Components from Nodes](https://rancher.com/docs/rancher/v2.x/en/cluster-admin/cleaning-cluster-nodes/)
- [卸载 Rancher 2.6](https://docs.rancher.cn/docs/rke2/install/linux_uninstall/_index)

    ```bash
    # --- clean rancher 2.6, 命令有先后顺序
    # ./rke2-killall.sh
    # ./rke2-uninstall.sh
    # systemctl stop /etc/systemd/system/rancher-system-agent.service
    # systemctl disable rancher-system-agent.service
    # rm /etc/systemd/system/rancher*
    # rm /usr/local/bin/rancher*
    # rm -rf /var/lib/rancher*
    # rm -rf /etc/rancher/*
    # systemctl stop kubelet.service
    # systemctl disable kubelet.service
    # rm /etc/systemd/system/kubelet.service
    # rm -rf /var/lib/etcd*
    # rm -rf /var/backups/kube_etcd
    # rm -rf /root/.kube
    ```

## FAQ
### 重置密码
```bash
docker ps -a|grep "rancher/rancher"
docker exec -it 288d7d0668a1 reset-password
New password for default admin user (user-rb2rs):
xxx # `xxx`即为新密码
```

### [rancher2.x升级](https://rancher.com/docs/rancher/v2.x/en/installation/install-rancher-on-linux/upgrades/#upgrading-both-rancher-and-the-underlying-cluster)