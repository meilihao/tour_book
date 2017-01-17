# Datetime时间误差1s

描述: datetime类型的时间存入db有时会出现误差,误差为`+1s`.

mysql 5.7 有该问题
mariadb 10.1 没有

原因:
mysql 5.7 支持[fractional seconds part(fsp,即小数秒)](http://dev.mysql.com/doc/refman/5.7/en/fractional-seconds.html)进行了四舍五入且没有warning或error的提示.
mariadb 10.1 则截断小数部分,但有warning提示.
