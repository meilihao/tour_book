# [cron](https://github.com/robfig/cron)
ref:
- [Go 每日一库之 cron](https://darjun.github.io/2020/06/25/godailylib/cron/)

## api
- cron.WithSeconds()

    默认的 unix cron 表达式格式为 5 字段（分 时 日 月 周），例如 "30 * * * *" 表示每小时的第 30 分钟.
    
    通过 WithSeconds() 后，表达式扩展为 6 字段（秒 分 时 日 月 周），例如 "0 30 * * * *" 表示每小时的第 30 分钟第 0 秒.