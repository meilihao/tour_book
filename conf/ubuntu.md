### FAQ

#### 查看cpu/硬盘温度,风扇转速

```shell
$ sudo apt-get install lm-sensors
$ sensors
```

#### ubuntu14.04 风扇狂转

```shell
# 方法1(**推荐**):
# 安装lm-sensors : http://kr.archive.ubuntu.com/ubuntu/pool/universe/l/lm-sensors/lm-sensors_3.3.4-2ubuntu1_amd64.deb
# 安装fancontrol : http://kr.archive.ubuntu.com/ubuntu/pool/universe/l/lm-sensors/fancontrol_3.3.4-2ubuntu1_all.deb
# 打开thinkpad_acpi的风扇控制支持,否则其参数(/proc/acpi/ibm/fan和/sys/devices/platform/thinkpad_hwmon/pwm1_enable)不允许修改,方法如下:
# 实现开机加载模块时设置，在/etc/modprobe.d/下增加一个配置文件thinkpad-acpi.conf，内容：
options thinkpad_acpi experimental=1 fan_control=1
# 重启电脑

# 启用thinkpad_acpi的风扇控制支持后发现风扇转数下来了,就暂时没用pwmconfig和fancontrol
# 执行pwmconfig
# 执行fancontrol

# ///
# 方法2:
sudo add-apt-repository ppa:linrunner/tlp
sudo apt-get update
sudo apt-get install tlp tlp-rdw
sudo tlp start
# 使用TPL是不需要进行配置的,效果一般.
# ///
# 方法3:
# 安装bumblebee,经测试,无效
```

参考:
- [使用thinkfan控制thinkpad风扇](http://vmlinz.is-programmer.com/posts/25834.html)
- [Ubuntu 10.04风扇声音太大](http://blog.chinaunix.net/uid-521083-id-2109229.html)