# test

## 描述

test用于执行条件检测,有助于避免使用过多的括号.实际使用中**不推荐**.

## 格式

    test EXPRESSION

## 例

```shell
if test $a = $b;then
echo 1
fi
---
test 表达式1 -a 表达式2 # 两个表达式都为真
test 表达式1 -o 表达式2 # 两个表达式有一个为真
```
