# chown

## 描述

改变文件拥有者/群组

## 语法格式

```
chown [OPTION]... [OWNER][:[GROUP]] FILE...
```

## 选项

- -v : 显示详细的处理过程
- -R : 递归的处理指定目录下的所有文件(包括子目录)

## 例

    # chown -R chen:chen ./ # 为当前目录下的所有文件修改拥有者和群组
    # chown User a.txt // 仅修改授权用户
    # chown :Group a.txt // 仅修改授权组