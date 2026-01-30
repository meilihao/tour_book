# pandas
ref:
- [Pandas 中文文档](https://pandas.woshinlper.com/)
- [Data Science Cheat Sheet 速查手册](https://s3.amazonaws.com/dq-blog-files/pandas-cheat-sheet.pdf)
- [**Python + Polars + DuckDB：2025年本地大数据分析黄金技术栈**](https://www.toutiao.com/article/7543705763735159346/)

    ```py
    con = duckdb.connect()
    con.register("clean_data", clean_df)
    con.query("select * from clean_data")
    ```
- [性能碾压pandas、polars的数据分析神器来了](https://cloud.tencent.com/developer/article/2428744)
- [Integration with Polars](https://aidoczh.com/duckdb/docs/guides/python/polars.html)

    `pip install -U duckdb 'polars[pyarrow]'`

**Python + Polars + DuckDB技术栈凭借"快、省、轻"三大优势，正在取代传统工具成为数据分析师的新宠. 迁移现有Pandas代码可使用pandas-to-polars自动转换工具**

Pandas 是 Python 的核心数据分析支持库，提供了快速、灵活、明确的数据结构，旨在简单、直观地处理关系型、标记型数据。Pandas 的目标是成为 Python 数据分析实践与实战的必备高级工具. Pandas 基于 NumPy 开发，可以与其它第三方科学计算支持库完美集成.

Pandas 的主要数据结构是 Series（一维数据）与 DataFrame（二维数据），这两种数据结构足以处理金融、统计、社会科学、工程等领域里的大多数典型用例.

在 Pandas 的agg()中，每个聚合操作通常只能直接访问它自己那一列

## 概念
- 标签

    在 pandas 中，label（标签） 通常指用于标识数据的索引名称或列名，是定位和访问数据的关键标识:
    1. 行标签（Row Label）: DataFrame 或 Series 的索引（index），用于标识每一行数据
    1. 列标签（Column Label）: DataFrame 的列名（columns），用于标识每一列数据

## 属性
- [freq](https://blog.csdn.net/qq_41780234/article/details/122536477)

## api
选项:
- pd.set_option('display.float_format', lamba x: '%.2f' % x)：设置浮点数显示格式

导入数据:
- pd.read_csv(filename)：从CSV文件导入数据
- pd.read_table(filename)：从限定分隔符的文本文件导入数据
- pd.read_excel(filename)：从Excel文件导入数据
- pd.read_sql(query, connection_object)：从SQL表/库导入数据
- pd.read_json(json_string)：从JSON格式的字符串导入数据
- pd.read_html(url)：解析URL、字符串或者HTML文件，抽取其中的tables表格
- pd.read_clipboard()：从你的粘贴板获取内容，并传给read_table()
- pd.DataFrame(dict)：从字典对象导入数据，Key是列名，Value是数据

导出数据:
- df.to_csv(filename)：导出数据到CSV文件
- df.to_excel(filename)：导出数据到Excel文件
- df.to_sql(table_name, connection_object)：导出数据到SQL表
- df.to_json(filename)：以Json格式导出数据到文本文件

创建测试对象
- pd.DataFrame(np.random.rand(20,5))：创建20行5列的随机数组成的DataFrame对象
- pd.Series(my_list)：从可迭代对象my_list创建一个Series对象
- df.index = pd.date_range('1900/1/30', periods=df.shape[0])：增加一个日期索引

查看、检查数据:
- df.head([n])：查看DataFrame对象的前n行, 默认n=5
- df.tail([n])：查看DataFrame对象的最后n行, 默认n=5
- df.shape()：查看行数和列数
- df.info()：查看索引、数据类型和内存信息
- df.describe()：查看数值型列的汇总统计
- s.value_counts(dropna=False)：查看Series对象的唯一值和计数
- df.apply(pd.Series.value_counts)：查看DataFrame对象中每一列的唯一值和计数

    groupby和apply通常联用即对分组应用apply
- df.groupby('A')['C1'].apply(min).reset_index() : groupby 分组默认会把分组依据列变成索引，用reset_index 将保留其在df的索引. 使用`reset_index(drop=True)`时就丢弃旧索引，生成从 0 开始的连续整数索引

数据选取:
- df[col]：根据列名，并以Series的形式返回列
- df[[col1, col2]]：以DataFrame形式返回多列
- s.iloc[0]：按位置选取数据
- s.loc['index_one']：按索引选取数据
- df.iloc[0,:]：返回第一行
- df.iloc[:13,:]：返回前12行
- df.iloc[:13,:4]：返回前12行, 前3列
- df.iloc[0,0]：返回第一列的第一个元素
- df[df["A"]==1] = df.loc[df["A"]==1],:]
- df.loc[df["A"]==1],['C1','C2']]
- df.loc[df["A"].isin([1,2]), ['C1']]
- df.loc[(df['A']=1) & (df['B']=2),:] : 查询条件用`&`包裹避免因优先级报错
- df[(df.hvac == 0) & (df.meter == 0)] : 多条件查询

数据清理:
- df.columns = ['a','b','c']：重命名列名
- pd.isnull()：检查DataFrame对象中的空值，并返回一个Boolean数组
- pd.notnull()：检查DataFrame对象中的非空值，并返回一个Boolean数组
- df.dropna()：删除所有包含空值的行
- df.dropna(axis=1)：删除所有包含空值的列, 即只要一行中任意一个字段为空，就会
被删除
- dropna(subset = ['city'])，来指定当一行中的 city 字段为空时，才会被删除
- df.drop_duplicates() : 去重默认会删掉完全重复的行（每个值都一样的行）
- df.drop_duplicates(subset='city')：去重但仅删除了这个字段重复的行，保留了各自不重复的第一行
- df.drop_duplicates(subset='city', keep='last')：类似上行, 但保留的是不重复的最后那行数据, keep默认是first
- df.dropna(axis=1,thresh=n)：删除所有小于n个非空值的行
- df.fillna(x)：用x替换DataFrame对象中所有的空值
- s.astype(float)：将Series中的数据类型更改为float类型
- s.replace(1,'one')：用‘one’代替所有等于1的值
- s.replace([1,3],['one','three'])：用'one'代替1，用'three'代替3
- df.rename(columns=lambda x: x + 1)：批量更改列名
- df.rename(columns={'old_name': 'new_ name'})：选择性更改列名
- df.set_index('column_one')：更改索引列
- df.rename(index=lambda x: x + 1)：批量重命名索引

数据处理：
- schema
    - df['A'] = range(1, 11) # 创建列A
    - df[['A','B']] = range(11, 21) # 创建列A, B
    - df.drop('A', axis=1, inplace=True) # 删除列A, axis=1表示对列操作，inplace=True表示直接修改原数据
- Filter、Sort和GroupBy:
    - df['A'] # 获取列A
    - df[['A','B']] # 获取列A, B
    - df[df[col] > 0.5]：选择col列的值大于0.5的行
    - df.sort_values(col1)：按照列col1排序数据，默认升序排列
    - df.sort_values(col2, ascending=False)：按照列col1降序排列数据
    - df.sort_values([col1,col2], ascending=[True,False])：先按列col1升序排列，后按col2降序排列数据
    - df.groupby(col)：返回一个按列col进行分组的Groupby对象. 流量级别作为汇总的依据列，默认转化为索引列，变成索引后print时该列会消失. 如果不希望它变成索引，向groupby 内传入参数 as_index = False
    - df.groupby([col1,col2])：返回一个按多列进行分组的Groupby对象
    - df.groupby(col1)[col2]：返回按列col1进行分组后，列col2的均值
    - df.pivot_table(index=col1, values=[col2,col3], aggfunc=max)：创建一个按列col1进行分组，并计算col2和col3的最大值的数据透视表
    - df.groupby(col1).agg(np.mean)：返回按列col1分组的所有列的均值
    - data.apply(np.mean)：对DataFrame中的每一列应用函数np.mean
    - data.apply(np.max,axis=1)：对DataFrame中的每一行应用函数np.max
    - pd.cut(x=df['A'], bins=[0,10,20], right=False, labels=['A', 'B']) : 切分（分桶）操作常用于一维数组的分类和打标. 默认是左开右闭, ight=False是左闭右开. labels是给分类打标签

数据合并:
- df1.append(df2)：将df2中的行添加到df1的尾部
- df.concat([df1, df2],axis=1)：将df2中的列添加到df1的尾部
- df1.join(df2,on=col1,how='inner')：对df1和df2执行SQL形式的join. how还支持left, right, outer
- pd.merge(left=df1,right=df2, on=col1, how='inner')：对df1和df2执行SQL形式的join, join的列名不同则用left_on和right_on指定
- pd.merge(left=df1, right=df2, left_index=True, right_index=True, how='inner') : 使用索引进行join

数据统计:
- df.describe()：查看数据值列的汇总统计
- df.mean()：返回所有列的均值
- df.corr()：返回列与列之间的相关系数
- df.count()：返回每一列中的非空值的个数
- df.max()：返回每一列的最大值
- df.min()：返回每一列的最小值
- df.median()：返回每一列的中位数
- df.std()：返回每一列的标准差

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
## example
ref:
- [十分钟入门 Pandas](https://pandas.woshinlper.com/docs/getting_started/10min/)

```py3
import numpy as np
import pandas as pd

# --- 生成对象
s = pd.Series([1, 3, 5, np.nan, 6, 8]) # 创建Series
In [4]: s
Out[4]: 
0    1.0
1    3.0
2    5.0
3    NaN
4    6.0
5    8.0
dtype: float64
dates = pd.date_range('20130101', periods=6)
In [6]: dates
Out[6]: 
DatetimeIndex(['2013-01-01', '2013-01-02', '2013-01-03', '2013-01-04',
               '2013-01-05', '2013-01-06'],
              dtype='datetime64[ns]', freq='D')

In [11]: df
Out[11]: 
                   A         B         C         D
2013-01-01  0.175082 -0.013153  1.019032  1.310196
2013-01-02  1.520397  0.002961  2.385698 -1.221322
2013-01-03 -0.452095 -0.043478 -1.427053  2.002147
2013-01-04 -0.546965 -1.184951  0.075188  0.047707
2013-01-05  0.059853  0.541067  0.032635 -0.083629
2013-01-06 -0.231986 -1.174046  1.206533 -0.792785
# 用 Series 字典对象生成 DataFrame
In [12]: df2 = pd.DataFrame({'A': 1.,
    ...:                   'B': pd.Timestamp('20130102'),
    ...:                   'C': pd.Series(1, index=list(range(4)), dtype='float32'),
    ...:                   'D': np.array([3] * 4, dtype='int32'),
    ...:                   'E': pd.Categorical(["test", "train", "test", "train"]),
    ...:                   'F': 'foo'})

In [13]: df2
Out[13]: 
     A          B    C  D      E    F
0  1.0 2013-01-02  1.0  3   test  foo
1  1.0 2013-01-02  1.0  3  train  foo
2  1.0 2013-01-02  1.0  3   test  foo
3  1.0 2013-01-02  1.0  3  train  foo
In [14]: df2.dtypes
Out[14]: 
A          float64
B    datetime64[s]
C          float32
D            int32
E         category
F           object
dtype: object
# --- 查看数据
In [15]: df.head() # 查看前5行
In [16]: df.tail(3) # 查看尾3行
In [17]: df.index # 查看索引
Out[17]: 
DatetimeIndex(['2013-01-01', '2013-01-02', '2013-01-03', '2013-01-04',
               '2013-01-05', '2013-01-06'],
              dtype='datetime64[ns]', freq='D')
In [19]: df.columns # 查看列名
Out[19]: Index(['A', 'B', 'C', 'D'], dtype='object')
In [20]: df.to_numpy() # 输出不包含行索引和列标签
Out[20]: 
array([[ 0.17508247, -0.01315297,  1.01903216,  1.31019569],
       [ 1.52039728,  0.00296072,  2.3856977 , -1.22132234],
       [-0.4520951 , -0.04347764, -1.42705254,  2.00214738],
       [-0.54696467, -1.18495102,  0.0751884 ,  0.04770738],
       [ 0.05985318,  0.54106671,  0.03263484, -0.08362899],
       [-0.23198584, -1.17404612,  1.20653341, -0.79278493]])
In [21]: df.T # 转置数据
Out[21]: 
   2013-01-01  2013-01-02  2013-01-03  2013-01-04  2013-01-05  2013-01-06
A    0.175082    1.520397   -0.452095   -0.546965    0.059853   -0.231986
B   -0.013153    0.002961   -0.043478   -1.184951    0.541067   -1.174046
C    1.019032    2.385698   -1.427053    0.075188    0.032635    1.206533
D    1.310196   -1.221322    2.002147    0.047707   -0.083629   -0.792785
In [22]: df.sort_index(axis=1, ascending=False) # 按列名降序排. axis(轴): 0, 对行索引排序（默认）;1,对列名排序. ascending(排序方向): False,降序;True, 升序
Out[22]: 
                   D         C         B         A
2013-01-01  1.310196  1.019032 -0.013153  0.175082
...
In [23]: df.sort_values(by='B') # Sort by column B
Out[23]: 
                   A         B         C         D
2013-01-04 -0.546965 -1.184951  0.075188  0.047707
2013-01-06 -0.231986 -1.174046  1.206533 -0.792785
# --- 获取数据
In [25]: df["A"] # 选择单列, 等同df.A
In [27]: df[:3] # 选择前3行
# --- 按标签选择
In [31]: df.loc[dates[0]] # 用标签提取一行数据
In [31]: df.loc['20130102', ['A', 'B']] # 用标签提取单行数据
In [31]: df.loc[:, ['A', 'B']] # 用标签选择多列数据
In [31]: df.loc['20130102':'20130104', ['A', 'B']] # 用标签切片，包含行与列结束点
In [30]: df.loc[dates[0], 'A'] # 提取标量值
Out[30]: 0.17508247095011448
In [31]: df.at[dates[0], 'A'] # 同上
# --- 按位置选择
In [32]: df.iloc[3] # 用整数位置选择
In [33]: df.iloc[3:5, 0:2] # 获取指定行的指定列
In [34]: df.iloc[[1, 2, 4], [0, 2]] # 同上
In [35]: df.iloc[1:3, :] # 获取指定行，所有列
In [36]: df.iloc[:, 1:3] # 获取所有行，指定列
In [30]: df.iloc[1,1] # 提取标量值
Out[30]: 0.17508247095011448
In [31]: df.iat[1,1] # 同上
# --- 布尔索引
In [39]: df[df.A > 0] # 获取A列大于0的行
In [40]: df[df > 0] # 选择 DataFrame 里满足条件的值
Out[40]: 
                   A         B         C         D
2013-01-01  0.469112       NaN       NaN       NaN
...
In [41]: df2 = df.copy()
In [42]: df2['E'] = ['one', 'one', 'two', 'three', 'four', 'three'] # 添加一列
In [44]: df2[df2['E'].isin(['two', 'four'])] # 用 isin() 筛选
# --- 复制
s1 = pd.Series([1, 2, 3, 4, 5, 6], index=pd.date_range('20130102', periods=6)) # 用索引自动对齐新增列的数据
In [48]: df.at[dates[0], 'A'] = 0 # 按标签赋值
In [49]: df.iat[0, 1] = 0 # 按位置赋值
In [52]: df2 = df.copy()
In [53]: df2[df2 > 0] = -df2 # 用 where 条件赋值
# --- 缺失值
In [55]: df1 = df.reindex(index=dates[0:4], columns=list(df.columns) + ['E']) # 重建索引（reindex）可以更改、添加、删除指定轴的索引，并返回数据副本，即不更改原数据
In [56]: df1.loc[dates[0]:dates[1], 'E'] = 1
In [58]: df1.dropna(how='any') # 删除所有含缺失值的行
In [59]: df1.fillna(value=5) # 填充缺失值
In [60]: pd.isna(df1) # 用于检测数据中缺失值（NaN/None） 的函数，返回一个与原 DataFrame 结构相同的布尔型（Boolean）DataFrame
# --- Apply 函数处理数据
In [67]: df.apply(lambda x: x.max() - x.min())
# -- others
for idx, row in df.iterrows(): # 迭代 DataFrame, row是副本
df["item_id"].iloc[0] # 获取第0行item_id的值
int(df.loc[0,"item_id"]) # 同上
df[["used", "cost", "unit_used"]] = 0.0 # 创建多个列, 默认值是0.0
type(df2["D"].iloc[0].astype(int)) # 输出numpy.int64
df2[(df2["E"] == "test") & (df2["F"].str.contains("foo", na=False))] # 不能去掉两个条件的(), 因为&的优先级高于==, 去掉会报"Cannot perform 'rand_' with a dtyped [bool] array and scalar of type [bool]"
df2['Deadline']=pd.to_datetime('2025-8-1') - pd.to_datetime(df2['Date'])
df2['Date']=pd.to_datetime(df2['Date'])
df_items["cost"].astype(float).sum() # 将cost列从string转为float后再求和

df_items = df_items[~df_items["pae"].isna()] # 使用布尔索引过滤 DataFrame, 只保留 "pae" 列不是 NaN 的行

def custom_agg(group):
    hvac_val = group['hvac'].iloc[0]

    if hvac_val == 0:
        return pd.Series({
        'meter': group['meter'].iloc[0],
        'energy': group['energy'].sum(),
        'uptime': group['uptime'].sum(),
        'cap': group['cap'].dropna().iloc[0],
    })
    else:
        unique_havc_dt_uptime = group[['hvac', 'uptime','dt']].drop_duplicates()

        return pd.Series({
            'meter': group['meter'].dropna().iloc[0],
            'energy': group['energy'].sum(),
            'uptime': unique_havc_dt_uptime['uptime'].sum(),
            'cap': group['cap'].dropna().iloc[0],
        })

df2 = df.groupby("hvac", as_index=False).apply(custom_agg).reset_index(drop=True)

df["tag"] = np.select(
    [
        (df.hvac == 0) & (df.meter == 0),
        (df.hvac != 0) & (df.meter == 0),
        (df.hvac != 0) & (df.meter != 0),
    ],
    [
        "Total",
        "HvacTotal",
        "Other",
    ]
) # 添加tag列, 用于辅助groupby

s = pd.Series(["10", "20", "abc", "30"])
pd.to_numeric(s, errors="coerce") # errors="coerce": 如果转换失败，不报错，直接变成 NaN
```