### 中文文件名是乱码(其实是UNICODE编码)

```shell
# 此时不会对0x80以上的字符进行quote,中文显示正常
git config --global core.quotepath false
```
### git diff 出现"^M"

Windows用CR LF来定义换行，Linux用LF,Mac用LF(Mac OS 9 及之前是CR，之后换成 UNIX 的 LF). CR全称是Carriage Return ,或者表示为\r, 意思是回车; LF全称是Line Feed，它才是真正意义上的换行表示符.

    git config --global core.whitespace cr-at-eol //让git diff忽略换行符的差异

### 添加/修改配置

    git config --global core.autocrlf input

### 删除配置

    git config --local --unset core.autocrlf

### git crlf换行符

core.autocrlf:

- true : 提交时自动地把行结束符CRLF转换成LF，而在签出代码时把LF转换成CRLF.适用于windows.
- input : 提交时把CRLF转换成LF，签出时不转换.适用于Linux,Mac等类UNIX系统.
- false : 提交和签出时均不转换.

### `.gitignore`无效，不能过滤某些文件

`.gitignore` 文件只能作用于 Untracked Files，也就是那些从来没有被 Git 记录过的文件（即自添加以后，从未 add 及 commit 过的文件).

之所以规则不生效，是因为那些文件曾经被 Git 记录过，因此 .gitignore 对它们完全无效.

解决方法：

    1. 用git rm从 Git 的数据库中删除对于该文件的追踪；
    2. 把对应的规则写入 .gitignore，让忽略真正生效；
    3. 提交＋推送。

最后有一点需要注意的，`git rm –cached filename`或`git rm -rf –cached foldername`删除的是追踪状态，而不是物理文件.