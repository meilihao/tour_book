# helm

## 添加应用
以[carina](https://github.com/carina-io/carina/blob/main/README_zh.md)举例:
```bash
helm repo add carina-csi-driver https://carina-io.github.io # 添加repo
helm search repo -l carina-csi-driver # 查看支持的version
helm install carina-csi-driver carina-csi-driver/carina-csi-driver --namespace kube-system --version v0.10.0
```