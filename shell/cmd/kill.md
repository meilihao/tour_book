# kill

## 描述

杀死进程

## 选项

- -l <信号编号> : 若不加<信息编号>选项，则-l参数会列出全部的信号名称
- -s <信号名称或编号> : 发送指定信号
- -<sigal> : 和`-s`相同, 默认15表示正常结束, 9,强制结束; 2, 结束进程但不强制, ctrl+c即是触发了`kill -2`.

## 例

```sh
$ kill PID
$ kill %job # 杀死job工作 (job为job number)
```
