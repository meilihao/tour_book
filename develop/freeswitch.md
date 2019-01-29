# freeswitch
v1.8

[安装文档](https://freeswitch.org/confluence/display/FREESWITCH/Installation)
[源码编译](https://freeswitch.org/confluence/display/FREESWITCH/Debian+9+Stretch)

参考:
- [FreeSWITCH中文文档](http://wiki.freeswitch.org.cn/wiki/Mod_Commands.html)
- [FreeSWITCH Event Socket library](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket)
- [闲聊语音编解码](http://www.ctiforum.com/news/guandian/369686.html)

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

1. [You must install libks to build mod_signalwire](https://freeswitch.org/jira/browse/FS-11579),官方暂无修复方法,先注释该mod
```
vim modules.conf # 注释`applications/mod_signalwire`
```

## 中文语音下载地址
https://files.freeswitch.org/releases/sounds/

## voip 客户端
- linux: Zoiper5
- window: xlite

## fs_cli
```
show codec # 查看支持的编码
sofia status profile internal reg # 当前注册的用户
module_exists mod_fy_tools # 检查模块是否存在
load mod_fy_tools # 加载模块
```

### FAQ
1. dialpan添加自定义变量
```xml
    <extension name="test">
      <condition field="destination_number" expression="^666$">
	      <action application="set" data="template_id=5469499f-b434-42af-8852-e89898b4c14f"/> // 自定义变量, 必须在`application="socket"`之前.
	      <action application="set" data="call_type=in_call"/>
        <action application="export" data="RECORD_STEREO=false"/> // 使用单声道(默认是两个, A/B-leg各一)
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

### 编译生成的libfreeswitch.so在哪
```sh
$ make
$ cd .libs
$ ll |grep libfreeswitch.so
lrwxrwxrwx 1 chen chen   22 11月  9 18:28 libfreeswitch.so -> libfreeswitch.so.1.0.0
lrwxrwxrwx 1 chen chen   22 11月  9 18:28 libfreeswitch.so.1 -> libfreeswitch.so.1.0.0
-rwxr-xr-x 1 chen chen  12M 11月  9 18:28 libfreeswitch.so.1.0.0
```

### 编译mod_iblc
由于其依赖的libilbc-dev没找到, 只好到[官方下载代码](https://freeswitch.org/stash/scm/sd/libilbc.git)手工编译.

在`cd src/mod/codecs/mod_ilbc`,`vim Makefile.am`, 把 if HAVE_ILBC 及 else 后面的相关逻辑去掉, Makefile.am变成下面的样子:
```
include $(top_srcdir)/build/modmake.rulesam
MODNAME=mod_ilbc

ILBC_CFLAGS = -I/usr/local/include
ILBC_LIBS = -L/usr/local/lib -lilbc

mod_LTLIBRARIES = mod_ilbc.la
...
mod_ilbc_la_LDFLAGS  = -avoid-version -module -no-undefined -shared
```

然后再执行 make, 系统就会重新生成 Makefile 并编译, 生成的so在`.libs`里.

### 编译mod_bcg729
```
cd freeswitch-1.8.2/src/mod/endpoints
git clone  https://github.com/xadhoom/mod_bcg729.git
cd mod_bcg729
make
cp mod_bcg729.so /usr/lib/freeswitch/mod
```

再编辑`autoload_configs/modules.conf.xml`: 注释`mod_g729`并添加`mod_bcg729`.

修改var.xml,启用G729:
```xml
  <X-PRE-PROCESS cmd="set" data="global_codec_prefs=OPUS,G729,G722,PCMU,PCMA,VP8"/> # 设置FreeSWITCH支持的媒体编码，包括语音和视频(默认仅支持音频编码)
  <X-PRE-PROCESS cmd="set" data="outbound_codec_prefs=OPUS,G729,G722,PCMU,PCMA,VP8"/>
```

最后重启fs即可.

ps:
1. 在vars.xml配置文件中加入`<X-PRE-PROCESScmd="set"data="media_mix_inbound_outbound_codecs=true"/>`,使得B-leg的编解码器列表跟A-LEG一样, 这样操作可以提高系统效率，尽量不转码，可以很大程度上增大系统效率.
1. 注意过长的global_codec_prefs列表可能会超出UDP的MTU(最大传输单元),那将引起呼叫建立失败

### sip.js 报 "unable to acquire streams" + "DOMException: Requested device not found"
```js
let options = {
                media: {
                    local: {
                        audio: document.getElementById('localVideo')
                    },
                    remote: {
                        video: document.getElementById('remoteVideo') // 电脑本身没有摄像头, 注释掉并删除页面中相应的vedio tag即可
                        // This is necessary to do an audio/video call as opposed to just **a video call**
                        audio: document.getElementById('remoteVideo')
                    }
                },
                ua: {
                    uri: _that.form.name + '@' + _that.form.ip
                    wsServers: 'wss://' + _that.form.ip + ':7443'
                    authorizationUser: _that.form.name

                    // FreeSWITCH Default Password
                    password: _that.form.password
                    displayName: _that.form.name
                }
            };
            _that.ua = new SIP.Web.Simple(options)
```

ps:
要根据实际情况决定options.media,并在页面添加相应的audio/video tag.

### internal user之间不能相呼
网络NAT问题.

### sip clinet 呼叫sip.js报"no suitablecandidates found...INCOMPATIBLE_DESTINATION"
`vim sudo deepin-editor sip_profiles/internal.xml`,在`<settings>`节添加如下内容:
```xml
<param name="apply-candidate-acl" value="localnet.auto"/>
<param name="apply-candidate-acl" value="wan_v4.auto"/>
<param name="apply-candidate-acl" value="rfc1918.auto"/>
<param name="apply-candidate-acl" value="any_v4.auto"/>
```

### systemd freeswitch无法启动
由于systemctl 启动一直失败, 复制`freeswitch.service`的命令在terminal运行:
```
# /usr/bin/freeswitch -u freeswitch -g freeswitch -ncwait
15944 Backgrounding.
FreeSWITCH[15943] Error starting system! pid:15944
```
原因未知, 无法找到错误日志.

解决方法:
删除`-u freeswitch -g freeswitch`参数即可.

### Send("api uuid_break xxx") get error "no reply"
```go
// 测试 uuid_break
// "github.com/fiorix/go-eventsocket/eventsocket"
func main() {
	eventsocket.ListenAndServe(":8690", handler)
}

func handler(c *eventsocket.Connection) {
	ev, err := c.Send("connect")
	c.Send("myevents")
    c.Execute("answer", "", false)
    // 播放长音乐
	if _, err = c.Execute("playback", "/home/chen/tmpfs/media/music.wav", false); err != nil {
		log.Fatal(err)
	}

	uuid := ev.Get("Caller-Unique-Id")
	time.Sleep(time.Second * 5)
 
	if _, err := c.Send("api uuid_break " + uuid); err != nil { // api uuid_break xxx 原本没有返回值且fs_cli也不会有相关执行记录,这里是eventsocket库错误报错
		log.Printf("Got error while uuid_break: %s", err)
	} else {
		log.Println("uuid_break is ok")
	}

	//发送app playback  播放开场音频（不支持mp3）
	ev, err = c.Execute("playback", "/home/chen/tmpfs/media/再见.wav", false)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}
```

### mod_sndfile.c:204 Error Opening File [xxx] [No Error.]
通常是权限错误, 因为freeswith.service的`ExecStart=/usr/bin/freeswitch -u root -g root -ncwait $DAEMON_OPTS`包含了用户信息, 修正即可.

### ESL 连接被拒绝
1. 检查`conf/autoload_configs/event_socket.conf.xml`将`<param name="listen-ip" value="127.0.0.1"/>`修改为`<param name="listen-ip" value="0.0.0.0"/>`(ipv4)或者`<param name="listen-ip" value="::"/>`(ipv6)，这样就允许远程ESL控制了.

1. fs_cli日志: `mod_event_socket.c:2659 IP ::ffff:192.168.11.134 Rejected by acl "loopback.auto"`
出现这个问题是因为被服务器拒绝，可以使用添加`<param name="apply-inbound-acl" value="lan"/>`, 其实就是启用了`conf/autoload_configs/acl.conf.xml`里的`<list name="lan" default="allow">`配置.

> `<param name="listen-ip" value="::"/>`,此时如果esl client的host是ipv4时话,系统会自动转成ipv6, 比如`192.168.11.134` -> `::ffff:192.168.11.134`

## 备注
- `conn.Execute("playback", record_file, true)`不是等freeswith播放完录音后执行完毕, 它会在file的PLAYBACK_STOP事件来之前就会完成.
- 播放多录音:　`conn.Execute("playback", "file_string://"+strings.Join(record_files,"!"), true)`