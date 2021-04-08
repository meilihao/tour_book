## time
### 几天内
```sql
SELECT * FROM 表名 WHERE DATE_ADD(createdate,INTERVAL 7 DAY) >= NOW();
```

### 相差天数
```sql
select datediff('2018-06-26','2018-06-25');
select datediff('2018-06-26 22:00:00','2018-06-25');
```

> 在日期计算中,如果存在时分秒的部分,是会被忽略的只对日期的部分进行计算即只对天计算

## json
### 数组包含
```sql
select * from user where json_contains(a, '["a"]', '$'); -- a = ["a", "b", "c"]
```