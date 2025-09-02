# polars
ref:
- [Polars - 用户指南](https://pola-rs.github.io/polars-book-cn/user-guide/index.html)

## example
```py3
import polars as pl

df = pl.read_csv("https://j.mp/iriscsv")
print(df.filter(pl.col("sepal_length") > 5)
      .group_by("species")
      .agg(pl.all().sum()) # 也可指定列pl.col("another_column").sum()
)
```

- agg()：是 "aggregate" 的缩写，表示对每个分组执行聚合操作, 括号中是具体的聚合逻辑