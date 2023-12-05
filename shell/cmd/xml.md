# xml

## xmlstarlet
ref:
- [Modify multiple lines of an XML file using command line](https://unix.stackexchange.com/questions/309676/modify-multiple-lines-of-an-xml-file-using-command-line)

支持修改xml

### example
```bash
xmlstarlet sel -t -v "//element/@attribute" file.xml
```

## xmllint
`xmllint --format XML_FILE`

## xpath
ref:
- [selenium之xpath语法总结](https://learnku.com/articles/50459)

XPath 是一种 XML 路径, 是一种语法或者语言用来查找使用 XML 路径表达的XML中的任意元素.

基本形式, `xpath=//tagname[@attribute='value']`:
- //：选中当前节点
- Tagname：特定节点的标记名
- @：选中属性的标记符
- Attribute：节点的属性名字
- Value：属性值


XPath 有两种类型：
1. 绝对 XPath 路径: 以单个正斜杠（/）开头，这意味着从根节点中选择元素

	缺点是如果元素路径中有一点儿变动的话，XPath 就会获取失败
2. 相对 XPath 路径: 以双正斜杠 // 开始

	相对 XPath 一直让人偏爱的原因就在于不需要从根元素得到一个完整路径

XPath axes 在 XML 文档中从当前上下文节点搜索不同的节点。XPath Axes 是查找动态元素的方法，否则，这是没有 ID、Classname，Name 等常规 XPath 方法无法实现的。

Axes 方法用来查找那些刷新或者执行其他操作而动态改变的元素。Selenium Webdriver 中常用的 Axes 方法很少，例如孩子 (child)，父母 (parent)，祖先 (ancestor)，兄弟姐妹 (sibling)，上一级 (preceding)，自己 (self) 等.