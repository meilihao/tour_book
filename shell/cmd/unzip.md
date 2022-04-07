# unzip
解压zip文件

## 选项
- -l : 查看zip包含的文件
- -v : 类同`-l`, 更详细, 比如压缩率
- -d : 指定解压目录
- -n : 默认覆盖文件, 本选项可不覆盖
- -o : 不必先询问用户，unzip执行后覆盖原有文件
- -t : 检查文件是否已损坏

## example
```bash
# unzip -o test.zip -d /tmp
```

## FAQ
### 解压zip中文乱码
解压时指定字符集`unzip -O CP936 xxx.zip` (用GBK, GB18030也可以)