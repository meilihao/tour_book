# chmod

## 描述

改变文件的权限

## 语法格式

```
chown [OPTION]... [MODE] FILE...
```

用户类型:
- o : other
- a : u + g + o
操作符:
- - : 取消
- + : 加入
- = : 设置


## 选项

- -v : 显示详细的处理过程
- -R : 递归的处理指定目录下的所有文件(包括子目录)