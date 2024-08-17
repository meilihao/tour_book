# bbc

## tracepoint
ref:
- [eBPF编写避坑指南](https://segmentfault.com/a/1190000041179276)

所有tracepoint: `/sys/kernel/debug/tracing/events`, 比如`tracepoint:skb:kfree_skb`, skb是类型, kfree_skb是具体函数.

具体相关函数的各文件说明:
- enable: 该文件用于启用或禁用 tracepoint 事件的跟踪。通过向该文件写入非零值，可以启用事件的跟踪；写入零值将禁用事件的跟踪。
- filter: 该文件用于设置 tracepoint 事件的过滤器。通过向该文件写入特定的过滤规则，您可以选择性地过滤掉某些事件，以便只记录感兴趣的事件。
- format: 该文件包含了 tracepoint 事件的格式信息。通过读取该文件，您可以查看事件的字段和其对应的数据类型。

    即有函数参数定义
- hist: 该文件包含了 tracepoint 事件的直方图（histogram）信息。通过读取该文件，您可以获取关于事件发生次数的统计信息，例如事件发生的频率、分布等。
- id: 该文件包含了 tracepoint 事件的唯一标识符。通过读取该文件，您可以获取事件的标识符，用于进一步操作或引用该事件。
- inject: 该文件用于向 tracepoint 事件注入数据。通过向该文件写入特定的数据，可以模拟触发事件并将自定义数据注入到事件中。
- trigger: 该文件用于触发 tracepoint 事件。通过向该文件写入非零值，可以手动触发事件的记录enable: 该文件用于启用或禁用 tracepoint 事件的跟踪。通过向该文件写入非零值，可以启用事件的跟踪；写入零值将禁用事件的跟踪。

这些文件提供了与 tracepoint 事件相关的控制和信息访问接口，您可以使用它们来配置和操作 tracepoint 事件的跟踪和记录。请注意，对这些文件的访问需要相应的权限（通常是 root 用户）.
