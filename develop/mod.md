# mod

## lsmod
查看当前kernel已加载的module信息, 信息来自`/proc/modules`.

输出信息:
- Module  : 模块名称
- Size  : 模块大小
- Used by : 依赖其他模块的个数 + 被其他模块依赖的列表

## modprobe
加载指定的模块到内核, 或者载入一组相互依赖的模块. 若在载入过程中发生错误，在modprobe会卸载整组的模块.

modprobe会查看模块 目录`/lib/modules/$(uname -r)`里面的所有模块和文件，除了可选的/etc/modprobe.conf配置文件和/etc/modprobe.d目录外.
modprobe需要一个最新的modules.dep(`/lib/modules/$(uname -r)/modules.dep`)文件，可以用depmod来生成. 该文件列出了每一个模块需要的其他模块，modprobe使用这个去自动添加或删除模块的依赖.

### 格式
```
# modprobe ${mod_name}
```

> modprobe就是调用insmod和rmmod来实现的.

### 选项
- -a : 加载一组匹配的模块
- -c : 输出所有模块的配置信息
- -C : 重载默认配置文件(/etc/modprobe.conf或/etc/modprobe.d)
- -D : 打印模块依赖
- -n : 不实际执行. 可以和-v选项一起使用，调试非常有用
- -r : 选项后指定模块时为卸载指定模块(会清理依赖). 与rmmod功能相同.
- -v : 详细信息
- -q : 不提示任何错误信息

## insmod/rmmod
载入/卸载mod.

## depmod
分析模块的依赖, 以供modprobe使用.

### 选项
- -a : 分析所有可用的模块
- -v : 输出详细信息

## FAQ
### modprobe和insmod区别
insmod不能处理依赖, 而modprobe可以.