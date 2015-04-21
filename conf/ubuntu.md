### FAQ

#### 查看cpu/硬盘温度,风扇转速

```shell
$ sudo apt-get install lm-sensors
$ sensors
```

#### ubuntu14.04 风扇狂转

```shell
http://vmlinz.is-programmer.com/posts/25834.html
http://blog.chinaunix.net/uid-521083-id-2109229.html
# 方法1:
# 安装lm-sensors : http://kr.archive.ubuntu.com/ubuntu/pool/universe/l/lm-sensors/lm-sensors_3.3.4-2ubuntu1_amd64.deb
# 安装fancontrol : http://kr.archive.ubuntu.com/ubuntu/pool/universe/l/lm-sensors/fancontrol_3.3.4-2ubuntu1_all.deb

弄不好pwmconfig,fancontrol无法执行,待解决.
调试pwmconfig脚本发现是其中pwmdisable()无法修改"/sys/devices/platform/thinkpad_hwmon/pwm1_enable"内容导致的.

玩thinkpad_acpi的风扇控制接口
cat /proc/acpi/ibm/fan可以得到控制风扇的方法(不知为何内容同样无法修改)

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
