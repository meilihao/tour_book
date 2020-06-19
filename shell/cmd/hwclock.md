# hwclock

## 描述

管理硬件时间

## 格式

    alias new_command='command seq'
    unalias new_command

## 例

    # hwclock # 获取硬件时间
    # hwclock -w # 将系统时间写入硬件

## FAQ
### `hwclock -w` 报`hwclock: Timed out waiting for time change.`
env: 飞腾cpu(aarch64)

使用`hwclock -w --debug`发现是系统时间与硬件时间差异太大, 不允许同步导致, 加参数`--update-drift`即可. 执行`hwclock -w --update-drift`成功需要一小段时间, 该时间内执行`hwclock`会报错.