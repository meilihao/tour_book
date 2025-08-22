# pandas
## Timestamp
- `pd.Timestamp(datetime.datetime.now()) + pd.DateOffset(n=-1, hours=1)`:将当前时间 now 往前推一个小时

    pd.Timestamp(now)：这会将now 时间（可以是字符串、datetime 对象或 timestamp）转换为 pandas 的 Timestamp 对象。Timestamp 是 pandas 中用于表示单个时间点的高级数据类型，它比 Python 原生的 datetime 对象功能更强大，尤其是在处理时间序列数据时。

    pd.DateOffset(n=-1, hours=1)：这是 DateOffset 对象。它是一个用来表示时间偏移量的类，例如年、月、日、小时、分钟等. n=-1 表示将时间往前移动一个单位, hours=1 表示这个移动的单位是小时

    将两者结合起来，pd.DateOffset(n=-1, hours=1) 实际上等价于 pd.DateOffset(hours=-1)，都是表示向前移动一个小时

## DataFrame
```py3
new_data = [
    {'name': '张三', 'age': 25},
    {'name': '李四', 'age': 30}
]

# 将新记录转换为 DataFrame
df_new = pd.DataFrame(new_data)

# 使用 pd.concat() 追加记录
df_final = pd.concat([df_empty, df_new], ignore_index=True) # 合并两个DataFrame, ignore_index=True: 第二个DataFrame放弃自身索引, 即索引从df_empty的最后一个索引开始递增
```


