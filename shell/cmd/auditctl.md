# auditctl

audit是linux系统中用于记录用户底层调用情况的系统,如记录用户执行的open,exit等系统调用, 并会将记录写到日志文件中.

audit可以通过使用auditctl命令来添加或删除audit规则. 设置针对某个用户进行记录,或针对某个进程的进行记录.

audit主要包含2个命令:
- auditd : audit服务进程
- auditctl : audit规则设置工具


## 选项
- s : 查看运行状态
- l : 查看现有audit规则
- a : 添加规则
- d : 删除规则
- D : 清除所有规则
- k <key> : 为一条audit规则设置一个关键字, 关键字可以是31个字节长的字符串,,用于过滤audit记录

## audit规则
`auditctl -a action,filter -S system_call -F field=value -k key_name`:
- action和filter 明确一个事件被记录

	action可以为always或者never，filter明确出对应的匹配过滤，filter可以为：task,exit，user，exclude。

- S(system_call) 明确出系统调用的名字，几个系统调用可以写在一个规则里，如-S xxx -S xxx.

	系统调用的名字可以在/usr/include/asm/unistd_64.h文件中找到

- field=value 作为附加匹配，修改规则以匹配特定架构、GroupID，ProcessID等的事件.

	比如`-F a0=0x6e9`, 这里选择匹配系统调用的第一个参数, 参数内容是要监控进程的PID（这里要用16进制, 也可用10进制）.

	具体有哪些字段，可以参考man linux  https://linux.die.net/man/8/auditctl

	- `b64`是64位sys_call
- k : 附加key, 便于ausearch检索

## example
```bash
# auditctl -a exit,always -F arch=b64 -S kill [-k audit_kill] # 查找who send sigkill. audit.log里可能没有相关进程的killed信息
# auditctl -a always,exit -F arch=b64 -F a1=15 -S kill -k log_kill
# auditctl -a always,exit -F arch=b64 -F a1=9 -S kill -k log_kill
# auditctl -a always,exit -F arch=b64 -F a1!=0x0 -S kill,tkill,tgkill -k log_kill
# service auditd restart
# ausearch -sc kill # 使用ausearch搜索结果
# ausearch -i -k log_kill # `-i`可解析 uid, gid, syscalls, arguments等
```

## FAQ
### 查找收到sigterm的方法(**推荐**)
ref:
- [Who’s killing that process? Who’s dumping prelink files in /tmp? — Linux auditd to the rescue](https://jotdownux.wordpress.com/2016/01/23/whos-killing-that-process-whos-dumping-prelink-files-in-tmp-linux-auditd-to-the-rescue/)
- [我的进程去哪儿了，谁杀了我的进程](https://www.cnblogs.com/xybaby/p/8098229.html)
- [Systemtap无法解决探测点，尽管它已显示在探测列表中 -- kernel编译时gcc选项问题](https://mlog.club/article/3989747)
- [揭开服务程序“被杀”之谜](https://cloud.tencent.com/developer/article/1639080)

- `strace -tt -o Signal_trace.out -p 21021 &`

**The value of si_pid indicates the PID of the process that sent the signal**.

> strace不能trace SIGKILL, 因为它也会被SIGKILL直接干掉.