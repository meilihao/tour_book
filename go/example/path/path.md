# path

import "path"

path实现了对斜杠分隔的路径的实用操作函数.

### func Ext(path string) string

功能:

获取path文件扩展名.返回值是路径最后一个斜杠分隔出的路径元素的最后一个'.'起始的后缀（包括'.'）,如果该元素没有'.'会返回空字符串.

参数列表:
- path 表示路径字符串


    fmt.Println(path.Ext("/a/b/c/bar.css")) // .css
    fmt.Println(path.Ext("/a/b/c/bar")) // ""

