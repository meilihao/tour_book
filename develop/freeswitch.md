# freeswitch
v1.8

[安装文档](https://freeswitch.org/confluence/display/FREESWITCH/Installation)

参考:
- [FreeSWITCH中文文档](http://wiki.freeswitch.org.cn/wiki/Mod_Commands.html)
- [FreeSWITCH Event Socket library](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket)

## error
> 安装缺失lib后需要清理`make clean && ./configure && make`

1. bootstrap: libtool not found.           You need libtool version 1.5.14 or newer to build FreeSWITCH from source
```
dpkg -L libtool

发现没有/usr/bin/libtool

dpkg -l libtool

libtool 是2.4.6-2版本的

在ubuntu只有libtoolize，修改bootstrap.sh，

libtool=${LIBTOOL:-`${LIBDIR}/apr/build/PrintPath glibtool libtool libtool22 libtool15 libtool14 libtoolize`}
```

1. mod_lua.cpp:37:10: fatal error: lua.h: 没有那个文件或目录
```
$ sudo apt install liblua5.3-dev 
$ cp /usr/include/lua5.3/*.h    src/mod/languages/mod_lua 
```

1. /usr/bin/ld: cannot find -llua
```
$ sudo find -name "liblua.so"
$ /usr/lib/x86_64-linux-gnu
$ sudo ln -sf liblua5.3.so liblua.so
```

1. You must install libopus-dev to build mod_opus
```
$ sudo apt install  libopus-dev
$ make clean && ./configure && make # 需要清理
```

1. You must install libsndfile-dev to build mod_sndfile
```
$ sudo apt install  libsndfile-dev
$ make clean && ./configure && make # 需要清理
```

## 中文语音下载地址
https://files.freeswitch.org/releases/sounds/

## voip 客户端
- linux: Zoiper5
- window: xlite

### FAQ
1. dialpan添加自定义变量
```xml
    <extension name="test">
      <condition field="destination_number" expression="^666$">
	      <action application="set" data="template_id=5469499f-b434-42af-8852-e89898b4c14f"/> // 自定义变量, 必须在`application="socket"`之前.
	      <action application="set" data="call_type=in_call"/>
	      <action application="socket" data="127.0.0.1:8080 async full"/>
      </condition>
    </extension>
```

`application="set"`和`application="export"`区别:

set,export 会出现在`answer, err := conn.ExecuteAnswer()`的answer中.
export会出现在`sm, err := conn.Execute()`的sm中.

出现的格式是`answer.GetHeader("Variable_template_id")`

> conn = github.com/0x19/goesl.SocketConnection

### 使用outbound模式外呼, tcp server丢失事件比如`CHANNEL_ANSWER`
原想将呼入和呼出统一用outbound模式处理, 发现clinet发起`originate {absolute_codec_string=PCMA}{origination_uuid=12345678,template_id=5469499fb43442af8852e89898b4c14f}user/1003 &socket('127.0.0.1:8689 async full')`并断开后, 发现server不能收到`CHANNEL_ANSWER`及其之前的所有事件.

解决方法: 使用inbound模式外呼, 在client端处理event.

### MEDIA_BUG_START
`MEDIA_BUG_START`由`RECORD_START`等操作语音流的指令引发, 可以查看`MEDIA_BUG_START`的header`Variable_current_application_data`和`Variable_current_application`. 

### inbound/outbound/socket
1. outbound呼入: 没提示音,有`CHANNEL_ANSWER`,`CHANNEL_HANGUP`
1. outbound呼出: 没提示音,有`CHANNEL_HANGUP`,没有`CHANNEL_ANSWER`; 调用`"bgapi originate xxx`后继续循环等待事件, 此时`&socket('127.0.0.1:8689 async full')`不起作用
1. inbound呼入(clinet): 有提示音, 有`CHANNEL_ANSWER`,`CHANNEL_HANGUP`

### INCOMPATIBLE_DESTINATION
外呼sip客户端时碰到`INCOMPATIBLE_DESTINATION`, 这通常是编码不兼容导致.

解决方法:
1. `bgapi originate {absolute_codec_string=PCMA,origination_uuid=12345678}user/1002 &park()`
删除`absolute_codec_string`即取消指定编码, 有freeswith自行判断

2. 检查freeswith conf/vars.xml和sip 客户端的编码是否有交集.

### 使用内网ip注册sip, freeswith日志显示是外网ip导致无法呼出
关闭全局的STUN, 比如Zoiper5-Settings-Advanced-Global STUN.

### freeswith高cpu
通常是通过log查找问题,这次不是.

> [Linux下如何定位高CPU/Memory的代码段](https://my.oschina.net/andywang1988/blog/698603)

```sh
$ ps -ef|grep freeswitch
$ top -H -p 20414 # 查看进程中所有线程对应的线程
$ gdb attach 20414
...
(gdb) info threads # 获得所有的线程信息, 与`top -H -p 20414`对照,找出问题线程
...
(gdb) thread 9 # 切换到问题线程的堆栈中
...
(gdb) bt # 查看当前的堆栈信息
...
```

没发现有价值的信息.

再通过strace查找, 哪些调用占用的时间最长，并对着堆栈和源码找出原因
```sh
# sudo strace -c -f -T -p 20414 # 发现errors列有一行数字很大, 推测freeswith 产生了大量error且应该会在log中体现出来
```

使用fs_cli, 发现大量错误日志:
```log
...
2018-11-09 10:42:14.513104 [ERR] switch_core_sqldb.c:587 NATIVE SQL ERR [cannot commit - no transaction is active]
COMMIT
2018-11-09 10:42:14.513104 [ERR] switch_core_sqldb.c:587 NATIVE SQL ERR [attempt to write a readonly database]
BEGIN EXCLUSIVE
2018-11-09 10:42:14.513104 [CRIT] switch_core_sqldb.c:1957 ERROR [attempt to write a readonly database]
2018-11-09 10:42:14.513104 [ERR] switch_core_sqldb.c:587 NATIVE SQL ERR [cannot commit -Content-Type: text/disconnect-noti
...
```

通过google查到[`switch_core_sqldb.c:587 NATIVE SQL ERR [cannot commit - no transaction is active]`](http://lists.freeswitch.org/pipermail/freeswitch-users/2015-October/116452.html), 判断应该是`sudo apt install`安装时db权限不正确导致.

解决:
```
# /var/lib
# chown -R freeswitch:freeswitch freeswitch
# systemctl restart freeswitch
```