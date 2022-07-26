# opencurve
## 部署


### 单机部署
ref:
- [CurveBS 集群拓扑](https://github.com/opencurve/curveadm/wiki/topology#curvebs-%E9%9B%86%E7%BE%A4%E6%8B%93%E6%89%91)
- [CurveBS 配置项](https://github.com/opencurve/curveadm/wiki/topology#curvebs-%E9%87%8D%E8%A6%81%E9%85%8D%E7%BD%AE%E9%A1%B9)

	s3信息是必填项, 要先部署minio, 否则在`curveadm status`会看到`role=snapshotclone`的容器相继退出.
- [demo config](https://github.com/opencurve/curveadm/tree/master/configs/bs/stand-alone)

模拟disk:
```bash
# modprobe nbd max_part=12 # 卸载用`rmmod nbd`
# qemu-img create -f qcow2 d1.qcow2 16G
# qemu-img create -f qcow2 d1.qcow2 16G
# qemu-img create -f qcow2 d1.qcow2 16G
# qemu-nbd -c /dev/nbd0 d1.qcow2 # 解除用`qemu-nbd -d /dev/nbd0`
# qemu-nbd -c /dev/nbd1 d2.qcow2
# qemu-nbd -c /dev/nbd2 d3.qcow2
```

> 注意: 后续的`curveadm format -f format.yaml`会对所有的disk进行format, 会使用掉整盘(比如这里的`d<N>.qcow2`刚开始很小, 之后会用到近16G)

开始部署:
```bash
# cd ~/opt/opencurve
# ssh-keygen -t ed25519 -f ~/.ssh/test_ed25519  -C "test@my"
# cat format.yaml # chen需要有sudo权限, 因为它会调用mount
user: chen
ssh_port: 22
private_key_file: /home/chen/.ssh/test_ed25519
container_image: opencurvedocker/curvebs:latest
host:
  - 192.168.0.43
disk:
  # support disk range like: /dev/sd[b-e]:/data/chunkserver[1-4]:90
  # if you want to exclude disk /dev/sdc and sdd, use this:
  #   /dev/sd[b-e,^cd]:/data/chunkserver[1-2]:80
  # if you have too many disks on a host(more than 26), use these:
  #   /dev/sd[a-z]:/data/chunkserver[1-26]:70
  #   /dev/sda[a-z]:/data/chunkserver[27-52]:60
  # also support disk and mountpoint with suffix, like:
  #   /dev/nvme[0-9]n1:/data/chunkserver[0-9]n1:90
  - /dev/nbd0:/data/chunkserver0:90  # device:mount_path:format_percent%
  - /dev/nbd1:/data/chunkserver1:90
  - /dev/nbd2:/data/chunkserver2:90
# cat topology.yaml
kind: curvebs
global:
  user: chen
  ssh_port: 22
  private_key_file: /home/chen/.ssh/test_ed25519
  container_image: opencurvedocker/curvebs:latest
  log_dir: /home/${user}/logs/${service_role}${service_host_sequence}
  data_dir: /home/${user}/data/${service_role}${service_host_sequence}
  s3.nos_address: xxx
  s3.snapshot_bucket_name: xxx
  s3.ak: xxx
  s3.sk: xxx
  variable:
    target: 192.168.0.43

etcd_services:
  config:
    listen.ip: ${service_host}
    listen.port: 2380${service_host_sequence}
    listen.client_port: 2379${service_host_sequence}
  deploy:
    - host: ${target}
    - host: ${target}
    - host: ${target}

mds_services:
  config:
    listen.ip: ${service_host}
    listen.port: 670${service_host_sequence}
    listen.dummy_port: 770${service_host_sequence}
  deploy:
    - host: ${target}
    - host: ${target}
    - host: ${target}

chunkserver_services:
  config:
    listen.ip: ${service_host}
    listen.port: 820${service_host_sequence}  # 8200,8201,8202
    data_dir: /data/chunkserver${service_host_sequence}  # /data/chunkserver0, /data/chunksever1
    copysets: 100
  deploy:
    - host: ${target}
    - host: ${target}
    - host: ${target}

snapshotclone_services:
  config:
    listen.ip: ${service_host}
    listen.port: 555${service_host_sequence}
    listen.dummy_port: 810${service_host_sequence}
    listen.proxy_port: 800${service_host_sequence}
  deploy:
    - host: ${target}
    - host: ${target}
    - host: ${target}
# curveadm format -f format.yaml # format disk
# curveadm format --status # 查看进度
# curveadm cluster add my-cluster -f topology.yaml
# curveadm cluster checkout my-cluster
# curveadm cluster ls
# curveadm deploy
# curveadm status [-v]
# cat client.yaml # mds.listen.addr来自`curveadm status`的输出
user: chen
host: 192.168.0.43
ssh_port: 22
private_key_file: /home/chen/.ssh/test_ed25519
container_image: opencurvedocker/curvebs:latest
mds.listen.addr: 192.168.0.43:6700,192.168.0.43:6701,192.168.0.43:6702
log_dir: /home/chen/curvebs/logs/client
# curveadm map chen:/test1 -c client.yaml --create --size 10GB # 大小必须是`N * 10GB`
# curveadm unmap chen:/test1 -c client.yaml
# --- 销毁cluster
# curveadm stop
# curveadm clean
# curveadm cluster rm my-cluster
# --- 清理disk
# umount /data/chunkserver0
# umount /data/chunkserver1
# umount /data/chunkserver2
# qemu-nbd -d /dev/nbd0
# qemu-nbd -d /dev/nbd1
# qemu-nbd -d /dev/nbd2
# rm -rf *.qcow2
# rmmod nbd
# lsblk
```