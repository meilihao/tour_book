# crontab

## 描述

按计划执行程序

每个用户的crontab保存在`/var/spool/cron/crontabs`, 且每个用户至多一个crontab.

> 同时指定weekday和day时, 满足其一就会被选中.
> crontab使用`/etc/crontab`指定的shell执行命令.

## 选项
- l : 打印出当前用户的crontab
- e : 为当前用户编辑crontab
- r : 删除当前用户的crontab

## 例
cron的格式：
```
*    *    *    *    *     <用户>            <要执行的命令>
T    T    T    T    T     (仅仅用于系统
|    |    |    |    |      的 crontab)
|    |    |    |    |
|    |    |    |    +----- 星期几 (0 - 6) (0 是星期天, 或者使用名称)
|    |    |    +---------- 月份 (1 - 12)
|    |    +--------------- 天 (1 - 31)
|    +-------------------- 小时 (0 - 23)
+------------------------- 分钟 (0 - 59)
```

比如:
```sh
47 6	* * 7	root	test -x /usr/sbin/anacron || ( cd / && run-parts --report /etc/cron.weekly )
```
