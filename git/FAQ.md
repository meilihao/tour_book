### 中文文件名是乱码(其实是UNICODE编码)

```shell
# 此时不会对0x80以上的字符进行quote,中文显示正常
git config core.quotepath false
```