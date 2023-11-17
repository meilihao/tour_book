# date

## 描述

显示或设置系统时间

## 选项

- -d,--date : 指定输入的日期与时间,字符串前后必须加上双引号
- -s : 根据字符串内容来设置日期与时间,字符串前后必须加上双引号

格式字符串:

%H 小时(以00-23来表示)
%I 小时(以01-12来表示)
%K 小时(以0-23来表示)
%l 小时(以0-12来表示)
%M 分钟(以00-59来表示)
%p AM或PM
%r 时间(含时分秒，小时以12小时AM/PM来表示)
%s 总秒数起算时间为1970-01-01 00:00:00 UTC
%S 秒(以本地的惯用法来表示)
%T 时间(含时分秒，小时以24小时制来表示)
%X 时间(以本地的惯用法来表示)
%Z 市区
%a 星期的缩写
%A 星期的完整名称
%b 月份英文名的缩写
%B 月份的完整英文名称
%c 日期与时间,和只输入date指令显示同样的结果
%d 日期(以01-31来表示)
%D 日期(含年月日)
%j 该年中的第几天
%m 月份(以01-12来表示)
%q 季度(1~4)
%U 该年中的周数
%w 该周的天数，0代表周日，1代表周一，异词类推
%x 年份日期(以本地的惯用法来表示)
%X = `%H:%M:%S`(以本地的惯用法来表示)
%y 年份(以00-99来表示)
%Y 年份(以四位数来表示)
%n 在显示时，插入新的一行
%t 在显示时，插入tab
MM 月份(必要)
DD 日期(必要)
hh 小时(必要)
mm 分钟(必要)
ss 秒(选择性)

## 例

    # date # 显示系统时间
    # date "+%Y-%m-%d %H:%M:%S" # 按照"年-月-日 小时:分钟:秒"的格式查看当前系统时间
    # date +%s # 显示Unix时间,这里的`+`表示用来启用某些选项
    # date -d @1501124007 # unix时间 -> 本地时间
    # date -du @1501124007 # unix时间 -> utc时间
    # date -d '07/20/2021 06:05:43' +"%s"
    1626732343
    # date -s "21 June 2009 11:01:22" # 设置系统时间
    # clock -w # 把系统时间写入CMOS

### example
```bash
# date -d '20200716 15:43:46.111222333' +'%s.%N' # `2020-07-16 15:43:46.111222333`也是支持的
1594885426.111222333
# date -d @1594885426.111222333
2020年 07月 16日 星期四 15:43:46 CST
# date +%s #显示当前时间的时间戳
1700118443
# date +'%Y-%m-%d %H:%M:%S'
2020-07-16 15:43:46
# date +'%F %T'
2020-07-16 15:43:46
# date -d '2 days ago' +'%F'
2020-07-14
# date -d '+4 hours' +'%T'
19:43:46
# date -d '2020-07-16 +3 days -1 minute' +'%F %T'
2020-07-18 23:59:00
```

## timedatectl
参考:
- [如何设置时间，时区和同步系统时钟使用timedatectl命令](https://www.howtoing.com/set-time-timezone-and-synchronize-time-using-timedatectl-command/)

systemd的时间工具.

### example
```bash
timedatectl list-timezones
timedatectl set-timezone "Asia/Kolkata"
timedatectl set-timezone UTC
timedatectl set-time '16:10:40 2015-11-20'
timedatectl set-time 20151120
timedatectl set-time 15:58:30
timedatectl set-local-rtc 1 # 将硬件时钟的时间标准设置为localtime
timedatectl set-local-rtc 0 # 将硬件时钟的时间标准设置为UTC
timedatectl set-ntp true # 启用远程NTP服务自动时间同步, 前提是安装NTP
timedatectl status # 查看时间
```